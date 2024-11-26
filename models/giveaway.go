package models

import "errors"

var (
	ErrEmptyGame = errors.New("game field cannot be empty")
	ErrEmptyCode = errors.New("code field cannot be empty")
)

type Giveaway struct {
	Game    *Game  `json:"game"`
	Code    string `json:"code"`
	Claimed bool   `json:"claimed"`
}

func (g *Giveaway) Validate() error {
	if g.Game == nil {
		return ErrEmptyGame
	}
	if g.Code == "" {
		return ErrEmptyCode
	}
	return nil
}
