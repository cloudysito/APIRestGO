package repository

import (
	"context"
	"errors"

	"github.com/cloudysito/apirestgo/models"
)

type PlayerRepository interface {
	Save(ctx context.Context, player *models.Player) error
	GetAll(ctx context.Context) ([]models.Player, error)
	GetByName(ctx context.Context, name string) (models.Player, error)
	GetRankStats(ctx context.Context) ([]models.RankStats, error)
}

type InMemoryRepository struct {
	players []models.Player
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		players: []models.Player{},
	}
}

func (r *InMemoryRepository) Save(ctx context.Context, player *models.Player) error {
	r.players = append(r.players, *player)
	return nil
}

func (r *InMemoryRepository) GetAll(ctx context.Context) ([]models.Player, error) {
	return r.players, nil
}

func (r *InMemoryRepository) GetByName(ctx context.Context, name string) (models.Player, error) {
	for _, player := range r.players {
		if player.Name == name {
			return player, nil
		}
	}

	return models.Player{}, errors.New("player not found")
}

func (r *InMemoryRepository) GetRankStats(ctx context.Context) ([]models.RankStats, error) {
	return nil, nil
}
