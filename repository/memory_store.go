package repository

import (
	"errors"
	"sync"

	"github.com/david-pawlowski/giveaway/models"
)

var ErrNoCodes = errors.New("you used all of my codes already!!!")

type InMemoryStore struct {
	codes []*models.Giveaway
	mutex sync.RWMutex
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		codes: make([]*models.Giveaway, 0),
	}
}

func (s *InMemoryStore) Add(code models.Giveaway) error {
	if err := code.Validate(); err != nil {
		return err
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.codes = append(s.codes, &code)
	return nil
}

func (s *InMemoryStore) GetRandomCode() (models.Giveaway, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, giveaway := range s.codes {
		if !giveaway.Claimed {
			giveaway.Claimed = true
			return *giveaway, nil
		}
	}
	return models.Giveaway{}, ErrNoCodes
}
