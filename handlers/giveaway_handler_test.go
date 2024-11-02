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

type MockGiveawayService struct {
	mockGetRandomCode func() (models.Giveaway, error)
	mockAdd           func(models.Giveaway)
}

func (m *MockGiveawayService) GetRandomCode() (models.Giveaway, error) {
	return m.mockGetRandomCode()
}

func (m *MockGiveawayService) Add(code models.Giveaway) {
	if m.mockAdd != nil {
		m.mockAdd(code)
	}
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
				Game:    "Cyberpunk 2077",
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
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "successful_creation",
			inputCode: models.Giveaway{
				Game:    "Half-Life 3",
				Code:    "HLIF-3333-GAME",
				Claimed: false,
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"game":"Half-Life 3","code":"HLIF-3333-GAME","claimed":false}`,
		},
		{
			name: "invalid_input",
			inputCode: models.Giveaway{
				Game: "",
				Code: "",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "game field cannot be empty\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var addCalled bool
			mockService := &MockGiveawayService{
				mockAdd: func(code models.Giveaway) {
					addCalled = true
					if code != tt.inputCode {
						t.Errorf("Expected code %v, got %v", tt.inputCode, code)
					}
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

			if tt.expectedStatus == http.StatusCreated {
				if !addCalled {
					t.Error("Expected Add to be called, but it wasn't")
				}

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
