package repository

import "github.com/david-pawlowski/giveaway/models"

type GiveawayRepository interface {
	Add(giveaway models.Giveaway) error
	GetRandomCode() (models.Giveaway, error)
}
