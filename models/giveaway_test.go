package models

import (
	"testing"
	"time"
)

func TestGiveaway_Validate(t *testing.T) {
	tests := []struct {
		name    string
		code    Giveaway
		wantErr error
	}{
		{
			name: "valid_code",
			code: Giveaway{
				Game:         "Half-Life 3",
				Code:         "HLIF-3333-GAME",
				Claimed:      false,
				CreatedAt:    time.Now(),
				Expires:      time.Now().Add(time.Hour),
				PrimaryImage: "http://test.go",
			},
			wantErr: nil,
		},
		{
			name: "empty_game",
			code: Giveaway{
				Game:         "",
				Code:         "HLIF-3333-GAME",
				Claimed:      false,
				CreatedAt:    time.Now(),
				Expires:      time.Now().Add(time.Hour),
				PrimaryImage: "http://test.go",
			},
			wantErr: ErrEmptyGame,
		},
		{
			name: "empty_code",
			code: Giveaway{
				Game:         "Half-Life 3",
				Code:         "",
				Claimed:      false,
				CreatedAt:    time.Now(),
				Expires:      time.Now().Add(time.Hour),
				PrimaryImage: "http://test.go",
			},
			wantErr: ErrEmptyCode,
		},
		{
			name: "both_empty",
			code: Giveaway{
				Game:         "",
				Code:         "",
				Claimed:      false,
				CreatedAt:    time.Now(),
				Expires:      time.Now().Add(time.Hour),
				PrimaryImage: "http://test.go",
			},
			wantErr: ErrEmptyGame, // Will return first error encountered
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.code.Validate()
			if err != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
