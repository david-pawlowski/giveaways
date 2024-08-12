package repository

import (
	"errors"
	"github.com/david-pawlowski/giveaway/models"
	"math/rand"
	"sync"
)

var ErrNotFound = errors.New("code not found")
var ErrNoCodes = errors.New("no more codes available")

type InMemoryStore struct {
	codes []models.Code
	mu    sync.Mutex
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		codes: make([]models.Code, 0),
	}
}

func (s *InMemoryStore) Add(code models.Code) (models.Code, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.codes = append(s.codes, code)
	return code, nil
}

func (s *InMemoryStore) Get(id int) (models.Code, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if id < 0 || id >= len(s.codes) {
		return models.Code{}, ErrNotFound
	}
	return s.codes[id], nil
}

func (s *InMemoryStore) MarkClaimed(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if id < 0 || id >= len(s.codes) {
		return false
	}
	s.codes[id].Claimed = true
	return true
}

func (s *InMemoryStore) GetRandomCode() (models.Code, error) {
	unclaimedCodes := []models.Code{}
	for _, code := range s.codes {
		if !code.Claimed {
			unclaimedCodes = append(unclaimedCodes, code)
		}
	}
	if len(unclaimedCodes) == 0 {
		return models.Code{}, ErrNoCodes
	}
	randomIndex := rand.Intn(len(unclaimedCodes))
	s.MarkClaimed(randomIndex)
	return unclaimedCodes[randomIndex], nil
}
