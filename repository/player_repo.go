package repository

import (
	"errors"

	"github.com/cloudysito/apirestgo/models"
)

type PlayerRepository interface {
	Save(player models.Player) error
	GetAll() ([]models.Player, error)
	GetByName(name string) (models.Player, error)
}

type InMemoryRepository struct {
	players []models.Player
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		players: []models.Player{},
	}
}

func (r *InMemoryRepository) Save(player models.Player) error {
	r.players = append(r.players, player)
	return nil
}

func (r *InMemoryRepository) GetAll() ([]models.Player, error) {
	return r.players, nil
}

func (r *InMemoryRepository) GetByName(name string) (models.Player, error) {
	for _, player := range r.players {
		if player.Name == name {
			return player, nil
		}
	}

	return models.Player{}, errors.New("player not found")
}
