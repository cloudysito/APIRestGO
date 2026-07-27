package repository

import (
	"context"

	"github.com/cloudysito/apirestgo/models"
	"go.mongodb.org/mongo-driver/bson"
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

func (r *FactionRepo) GetAllFactions(ctx context.Context) ([]models.Faction, error) {
	var factions []models.Faction
	cursor, err := r.collection.Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var faction models.Faction
		if err := cursor.Decode(&faction); err != nil {
			return nil, err
		}
		factions = append(factions, faction)
	}

	return factions, nil
}
