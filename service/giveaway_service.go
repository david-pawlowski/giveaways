package service

import (
	"github.com/david-pawlowski/giveaway/models"
	"github.com/david-pawlowski/giveaway/repository"
)

type GiveawayService interface {
	Add(code models.Code) (models.Code, error)
	Get(id int) (models.Code, error)
	MarkClaimed(id int) bool
	GetRandomCode() (models.Code, error)
}

type giveawayService struct {
	store *repository.InMemoryStore
}

func NewGiveawayService(store *repository.InMemoryStore) GiveawayService {
	return &giveawayService{
		store: store,
	}
}

func (gs *giveawayService) Add(code models.Code) (models.Code, error) {
	return gs.store.Add(code)
}

func (gs *giveawayService) Get(id int) (models.Code, error) {
	return gs.store.Get(id)
}

func (gs *giveawayService) MarkClaimed(id int) bool {
	return gs.store.MarkClaimed(id)
}

func (gs *giveawayService) GetRandomCode() (models.Code, error) {
	return gs.store.GetRandomCode()
}
