package repository

import (
	"context"
	"errors"

	"github.com/cloudysito/apirestgo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type FactionRepository interface {
	CreateFaction(ctx context.Context, faction *models.Faction) error
	GetAllFactions(ctx context.Context) ([]models.Faction, error)
	GetFactionByName(ctx context.Context, name string) (models.Faction, error)
	UpdateFaction(ctx context.Context, name string, updatedFaction *models.Faction) error
	UpdateFactionPartial(ctx context.Context, name string, updates map[string]interface{}) error
	DeleteFaction(ctx context.Context, name string) error
}

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

func (r *FactionRepo) GetFactionByName(ctx context.Context, name string) (models.Faction, error) {
	var faction models.Faction
	filter := bson.M{"name": name}
	err := r.collection.FindOne(ctx, filter).Decode(&faction)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Faction{}, errors.New("Faction not found.")
		}
		return models.Faction{}, err
	}

	return faction, nil
}

func (r *FactionRepo) UpdateFaction(ctx context.Context, name string, updatedFaction *models.Faction) error {
	filter := bson.M{"name": name}
	update := bson.M{"$set": updatedFaction}
	result, err := r.collection.UpdateOne(ctx, filter, update)

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("Faction not found.")
	}

	return nil
}

func (r *FactionRepo) UpdateFactionPartial(ctx context.Context, name string, updates map[string]interface{}) error {
	filter := bson.M{"name": name}
	update := bson.M{"$set": updates}
	result, err := r.collection.UpdateOne(ctx, filter, update)

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("Faction not found.")
	}

	return nil
}

func (r *FactionRepo) DeleteFaction(ctx context.Context, name string) error {
	filter := bson.M{"name": name}
	result, err := r.collection.DeleteOne(ctx, filter)

	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("Faction not found.")
	}

	return nil
}
