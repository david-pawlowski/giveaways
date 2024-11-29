package models

import (
	"errors"
	"time"
)

var (
	ErrEmptyGame = errors.New("game field cannot be empty")
	ErrEmptyCode = errors.New("code field cannot be empty")
)

// Should I use constructor to make defaults, and autoincrement id?
type Giveaway struct {
	ID           int       `json:"id"`
	CreatedAt    time.Time `json:"createdAt"`
	Expires      time.Time `json:"expires"`
	PrimaryImage string    `json:"primaryImage"`
	Game         string    `json:"game"`
	Code         string    `json:"code"`
	Claimed      bool      `json:"claimed"`
}

func (g *Giveaway) Validate() error {
	if g.Game == "" {
		return ErrEmptyGame
	}
	if g.Code == "" {
		return ErrEmptyCode
	}
	return nil
}
