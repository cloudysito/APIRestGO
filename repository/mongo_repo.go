package repository

import (
	"context"
	"time"

	"github.com/cloudysito/apirestgo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(client *mongo.Client) *MongoRepository {
	return &MongoRepository{
		collection: client.Database("apirestgo").Collection("players"),
	}
}

func (r *MongoRepository) Save(player models.Player) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.InsertOne(ctx, player)
	return err
}

func (r *MongoRepository) GetAll() ([]models.Player, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var players []models.Player

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var p models.Player
		if err := cursor.Decode(&p); err != nil {
			return nil, err
		}
		players = append(players, p)
	}

	return players, nil
}

func (r *MongoRepository) GetByName(name string) (models.Player, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var player models.Player

	filter := bson.M{"name": name}
	err := r.collection.FindOne(ctx, filter).Decode(&player)

	return player, err
}
