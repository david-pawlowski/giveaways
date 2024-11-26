package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/david-pawlowski/giveaway/models"
	"github.com/david-pawlowski/giveaway/repository"
)

// MockGiveawayService implements the GiveawayService interface for testing
type MockGiveawayService struct {
	mockGetRandomCode func() (models.Giveaway, error)
	mockAdd           func(models.Giveaway) error
}

func (m *MockGiveawayService) GetRandomCode() (models.Giveaway, error) {
	return m.mockGetRandomCode()
}

func (m *MockGiveawayService) Add(code models.Giveaway) error {
	if m.mockAdd != nil {
		return m.mockAdd(code)
	}
	return nil
}

func TestGetRandomCode(t *testing.T) {
	tests := []struct {
		name           string
		mockResponse   models.Giveaway
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "successful_retrieval",
			mockResponse: models.Giveaway{
				Game: &models.Game{
					ID:           1,
					Name:         "Cyberpunk 2077",
					Category:     "Test",
					DevelopedBy:  "Valve",
					PrimaryImage: "http://test.cdn/image1",
				},
				Code:    "ABCD-1234-EFGH",
				Claimed: false,
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"game":"Cyberpunk 2077","code":"ABCD-1234-EFGH","claimed":false}`,
		},
		{
			name:           "no_codes_available",
			mockResponse:   models.Giveaway{},
			mockError:      repository.ErrNoCodes,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "We are out of codes.\n",
		},
		{
			name:           "internal_server_error",
			mockResponse:   models.Giveaway{},
			mockError:      errors.New("database connection failed"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Internal server error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockGiveawayService{
				mockGetRandomCode: func() (models.Giveaway, error) {
					return tt.mockResponse, tt.mockError
				},
			}

			handler := NewGiveawayHandler(mockService)

			req := httptest.NewRequest(http.MethodGet, "/giveaway", nil)
			w := httptest.NewRecorder()

			handler.GetRandomCode(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Test %s: expected status %d; got %d", tt.name, tt.expectedStatus, w.Code)
			}

			if tt.mockError == nil {
				var got models.Giveaway
				err := json.Unmarshal(w.Body.Bytes(), &got)
				if err != nil {
					t.Fatalf("Test %s: failed to unmarshal response: %v", tt.name, err)
				}

				var want models.Giveaway
				err = json.Unmarshal([]byte(tt.expectedBody), &want)
				if err != nil {
					t.Fatalf("Test %s: failed to unmarshal expected body: %v", tt.name, err)
				}

				if got != want {
					t.Errorf("Test %s: expected body %v; got %v", tt.name, want, got)
				}
			} else {
				if w.Body.String() != tt.expectedBody {
					t.Errorf("Test %s: expected body %q; got %q", tt.name, tt.expectedBody, w.Body.String())
				}
			}
		})
	}
}

func TestCreateCode(t *testing.T) {
	tests := []struct {
		name           string
		inputCode      models.Giveaway
		mockAddError   error
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "successful_creation",
			inputCode: models.Giveaway{
				Game: &models.Game{
					ID:           1,
					Name:         "Cyberpunk 2077",
					Category:     "Test",
					DevelopedBy:  "Valve",
					PrimaryImage: "http://test.cdn/image1",
				},
				Code:    "HLIF-3333-GAME",
				Claimed: false,
			},
			mockAddError:   nil,
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"game":"Half-Life 3","code":"HLIF-3333-GAME","claimed":false}`,
		},
		{
			name: "invalid_input",
			inputCode: models.Giveaway{
				Game: nil, // Empty game name
				Code: "",
			},
			mockAddError:   errors.New("invalid input"),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid input\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockGiveawayService{
				mockAdd: func(code models.Giveaway) error {
					if tt.mockAddError != nil {
						return tt.mockAddError
					}
					if code != tt.inputCode {
						t.Errorf("Expected code %v, got %v", tt.inputCode, code)
					}
					return nil
				},
			}

			handler := NewGiveawayHandler(mockService)

			body, err := json.Marshal(tt.inputCode)
			if err != nil {
				t.Fatalf("Failed to marshal input code: %v", err)
			}

			req := httptest.NewRequest(http.MethodPost, "/giveaway", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			handler.CreateCode(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Test %s: expected status %d; got %d", tt.name, tt.expectedStatus, w.Code)
			}

			if tt.mockAddError == nil {
				var got models.Giveaway
				err := json.Unmarshal(w.Body.Bytes(), &got)
				if err != nil {
					t.Fatalf("Test %s: failed to unmarshal response: %v", tt.name, err)
				}

				if got != tt.inputCode {
					t.Errorf("Test %s: expected body %v; got %v", tt.name, tt.inputCode, got)
				}
			} else {
				if w.Body.String() != tt.expectedBody {
					t.Errorf("Test %s: expected body %q; got %q", tt.name, tt.expectedBody, w.Body.String())
				}
			}
		})
	}
}
