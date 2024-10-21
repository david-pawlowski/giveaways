package service

import (
	"github.com/david-pawlowski/giveaway/models"
	"github.com/david-pawlowski/giveaway/repository"
)

type GiveawayService interface {
	Add(giveaway models.Giveaway)
	GetRandomCode() (models.Giveaway, error)
}

type giveawayService struct {
	store *repository.InMemoryStore
}

func NewGiveawayService(store *repository.InMemoryStore) GiveawayService {
	return &giveawayService{
		store: store,
	}
}

func (gs *giveawayService) Add(giveaway models.Giveaway) {
	gs.store.Add(giveaway)
}

func (gs *giveawayService) GetRandomCode() (models.Giveaway, error) {
	return gs.store.GetRandomCode()
}
