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

func (r *MongoRepository) Save(ctx context.Context, player *models.Player) error {
	cctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.collection.InsertOne(cctx, *player)
	return err
}

func (r *MongoRepository) GetAll(ctx context.Context) ([]models.Player, error) {
	cctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var players []models.Player

	cursor, err := r.collection.Find(cctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(cctx)

	for cursor.Next(cctx) {
		var p models.Player
		if err := cursor.Decode(&p); err != nil {
			return nil, err
		}
		players = append(players, p)
	}

	return players, nil
}

func (r *MongoRepository) GetByName(ctx context.Context, name string) (models.Player, error) {
	cctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var player models.Player

	filter := bson.M{"name": name}
	err := r.collection.FindOne(cctx, filter).Decode(&player)

	return player, err
}

func (r *MongoRepository) GetRankStats(ctx context.Context) ([]models.RankStats, error) {
	cctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$rank"}, {Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}}}}},
	}

	cursor, err := r.collection.Aggregate(cctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(cctx)

	var stats []models.RankStats
	if err := cursor.All(cctx, &stats); err != nil {
		return nil, err
	}

	return stats, nil
}

func (r *MongoRepository) UpdateRank(ctx context.Context, name string, newRank string) error {
	cctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"name": name}
	update := bson.M{"$set": bson.M{"rank": newRank}}

	_, err := r.collection.UpdateOne(cctx, filter, update)
	return err
}

func (r *MongoRepository) DeletePlayer(ctx context.Context, name string) error {
	cctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"name": name}
	_, err := r.collection.DeleteOne(cctx, filter)
	return err
}
