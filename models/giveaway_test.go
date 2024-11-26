package models

import "testing"

func TestGiveaway_Validate(t *testing.T) {
	tests := []struct {
		name    string
		code    Giveaway
		wantErr error
	}{
		{
			name: "valid_code",
			code: Giveaway{
				Game: &Game{
					ID:           1,
					Name:         "Half-Life 3",
					Category:     "Test",
					DevelopedBy:  "Valve",
					PrimaryImage: "http://test.cdn/image1",
				},
				Code:    "HLIF-3333-GAME",
				Claimed: false,
			},
			wantErr: nil,
		},
		{
			name: "empty_game",
			code: Giveaway{
				Game:    nil,
				Code:    "HLIF-3333-GAME",
				Claimed: false,
			},
			wantErr: ErrEmptyGame,
		},
		{
			name: "empty_code",
			code: Giveaway{
				Game: &Game{
					ID:           1,
					Name:         "Half-Life 3",
					Category:     "Test",
					DevelopedBy:  "Valve",
					PrimaryImage: "http://test.cdn/image1",
				},
				Code:    "",
				Claimed: false,
			},
			wantErr: ErrEmptyCode,
		},
		{
			name: "both_empty",
			code: Giveaway{
				Game:    nil,
				Code:    "",
				Claimed: false,
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
