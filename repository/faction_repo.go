package repository

import (
	"context"

	"github.com/cloudysito/apirestgo/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type FactionRepo struct {
	collection *mongo.Collection
}

func NewFactionRepo(db *mongo.Database) *FactionRepo {
	return &FactionRepo{
		collection: db.Collection("factions"),
	}
}

func (r *FactionRepo) CreateFaction(ctx context.Context, faction *models.Faction) error {
	_, err := r.collection.InsertOne(ctx, faction)
	return err
}
