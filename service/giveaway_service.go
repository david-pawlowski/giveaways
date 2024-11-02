package service

import (
	"errors"
	"github.com/david-pawlowski/giveaway/models"
	"github.com/david-pawlowski/giveaway/repository"
	"log"
)

type GiveawayService interface {
	Add(giveaway models.Giveaway) error
	GetRandomCode() (models.Giveaway, error)
}

type giveawayService struct {
	store *repository.InMemoryStore
}

func NewGiveawayService(store *repository.InMemoryStore) (GiveawayService, error) {
	if store == nil {
		return nil, errors.New("store cannot be nil")
	}
	return &giveawayService{
		store: store,
	}, nil
}

func (gs *giveawayService) Add(giveaway models.Giveaway) error {
	if err := giveaway.Validate(); err != nil {
		log.Printf("Invalid giveaway data: %v", err)
		return err
	}
	gs.store.Add(giveaway)
	return nil
}

func (gs *giveawayService) GetRandomCode() (models.Giveaway, error) {
	code, err := gs.store.GetRandomCode()
	if err != nil {
		log.Printf("Error getting random code: %v", err)
		return models.Giveaway{}, err
	}
	return code, nil
}
