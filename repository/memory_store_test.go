package repository

import (
	"fmt"
	"sync"
	"testing"

	"github.com/david-pawlowski/giveaway/models"
)

func TestInMemoryStore_Add(t *testing.T) {
	tests := []struct {
		name    string
		code    models.Giveaway
		wantErr bool
	}{
		{
			name: "valid giveaway",
			code: models.Giveaway{
				Game:    "TestGame",
				Code:    "1234-1234",
				Claimed: false,
			},
			wantErr: false,
		},
		{
			name: "empty game name",
			code: models.Giveaway{
				Game:    "",
				Code:    "1234-1234",
				Claimed: false,
			},
			wantErr: true,
		},
		{
			name: "empty code",
			code: models.Giveaway{
				Game:    "TestGame",
				Code:    "",
				Claimed: false,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewInMemoryStore()
			err := store.Add(tt.code)

			if (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && len(store.codes) != 1 {
				t.Errorf("Add() didn't store the code, got len = %d, want 1", len(store.codes))
			}
		})
	}
}

func TestInMemoryStore_GetRandomCode(t *testing.T) {
	tests := []struct {
		name      string
		setupFn   func(*InMemoryStore)
		wantErr   bool
		wantCode  string
		checkCode func(*testing.T, models.Giveaway, error)
	}{
		{
			name:    "get code from empty store",
			setupFn: func(s *InMemoryStore) {},
			wantErr: true,
			checkCode: func(t *testing.T, g models.Giveaway, err error) {
				if err != ErrNoCodes {
					t.Errorf("expected ErrNoCodes, got %v", err)
				}
			},
		},
		{
			name: "get code from store with one unclaimed code",
			setupFn: func(s *InMemoryStore) {
				s.Add(models.Giveaway{
					Game: "TestGame",
					Code: "1234-1234",
				})
			},
			wantErr:  false,
			wantCode: "1234-1234",
			checkCode: func(t *testing.T, g models.Giveaway, err error) {
				if !g.Claimed {
					t.Error("code should be marked as claimed")
				}
			},
		},
		{
			name: "get code from store with all claimed codes",
			setupFn: func(s *InMemoryStore) {
				s.Add(models.Giveaway{
					Game:    "TestGame",
					Code:    "1234-1234",
					Claimed: true,
				})
			},
			wantErr: true,
			checkCode: func(t *testing.T, g models.Giveaway, err error) {
				if err != ErrNoCodes {
					t.Errorf("expected ErrNoCodes, got %v", err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewInMemoryStore()
			tt.setupFn(store)

			got, err := store.GetRandomCode()

			if (err != nil) != tt.wantErr {
				t.Errorf("GetRandomCode() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantCode != "" && got.Code != tt.wantCode {
				t.Errorf("GetRandomCode() got = %v, want %v", got.Code, tt.wantCode)
			}

			if tt.checkCode != nil {
				tt.checkCode(t, got, err)
			}
		})
	}
}

func TestInMemoryStore_ConcurrentAccess(t *testing.T) {
	store := NewInMemoryStore()

	// Add some initial codes
	initialCodes := []models.Giveaway{
		{Game: "Game1", Code: "CODE1"},
		{Game: "Game2", Code: "CODE2"},
		{Game: "Game3", Code: "CODE3"},
	}

	for _, code := range initialCodes {
		store.Add(code)
	}

	var wg sync.WaitGroup
	results := make(chan error, len(initialCodes))

	// Test concurrent retrievals
	// Only attempt to retrieve the exact number of codes we have
	for i := 0; i < len(initialCodes); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := store.GetRandomCode()
			results <- err
		}()
	}

	wg.Wait()
	close(results)

	// Count successful retrievals
	successCount := 0
	for err := range results {
		if err == nil {
			successCount++
		}
	}

	// We should have exactly len(initialCodes) successful retrievals
	if successCount != len(initialCodes) {
		t.Errorf("Expected %d successful retrievals, got %d", len(initialCodes), successCount)
	}
}

// Add a separate test for concurrent adds and retrieval behavior
func TestInMemoryStore_ConcurrentAddAndRetrieve(t *testing.T) {
	store := NewInMemoryStore()
	var wg sync.WaitGroup

	// Test concurrent additions
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			code := models.Giveaway{
				Game: fmt.Sprintf("Game%d", i),
				Code: fmt.Sprintf("CODE%d", i),
			}
			if err := store.Add(code); err != nil {
				t.Errorf("Failed to add code: %v", err)
			}
		}(i)
	}

	// Wait for additions to complete
	wg.Wait()

	// Verify the number of codes in the store
	if len(store.codes) != 5 {
		t.Errorf("Expected 5 codes in store, got %d", len(store.codes))
	}

	// Try to retrieve all codes
	retrievedCodes := make(map[string]bool)
	var retrievalMutex sync.Mutex

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			code, err := store.GetRandomCode()
			if err != nil {
				t.Errorf("Unexpected error retrieving code: %v", err)
				return
			}

			retrievalMutex.Lock()
			retrievedCodes[code.Code] = true
			retrievalMutex.Unlock()
		}()
	}

	wg.Wait()

	// Verify we got all unique codes
	if len(retrievedCodes) != 5 {
		t.Errorf("Expected to retrieve 5 unique codes, got %d", len(retrievedCodes))
	}

	// Verify all codes are now claimed
	_, err := store.GetRandomCode()
	if err != ErrNoCodes {
		t.Errorf("Expected ErrNoCodes after all codes claimed, got %v", err)
	}
}

func TestInMemoryStore_GetRandomCodeDistribution(t *testing.T) {
	store := NewInMemoryStore()
	codes := []models.Giveaway{
		{Game: "Game1", Code: "CODE1"},
		{Game: "Game2", Code: "CODE2"},
		{Game: "Game3", Code: "CODE3"},
	}

	iterations := 100
	for i := 0; i < iterations; i++ {
		for _, code := range codes {
			store.Add(code)
		}

		for range codes {
			code, err := store.GetRandomCode()
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !code.Claimed {
				t.Error("Retrieved code should be marked as claimed")
			}
		}

		_, err := store.GetRandomCode()
		if err != ErrNoCodes {
			t.Errorf("Expected ErrNoCodes, got %v", err)
		}
	}
}
