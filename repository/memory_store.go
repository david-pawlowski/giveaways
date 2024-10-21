package repository

import (
	"errors"
	"github.com/david-pawlowski/giveaway/models"
)

var ErrNoCodes = errors.New("You used all of my codes already!!!")

type InMemoryStore []*models.Giveaway

func (s *InMemoryStore) Add(code models.Giveaway) {
	(*s) = append((*s), &code)
}

func (s *InMemoryStore) GetRandomCode() (models.Giveaway, error) {
	for i := range *s {
		giveaway := (*s)[i]
		if !giveaway.Claimed {
			giveaway.Claimed = true
			return *giveaway, nil
		}
	}
	return models.Giveaway{}, ErrNoCodes
}
