package repository

import (
	"context"
	"time"

	"github.com/cloudysito/apirestgo/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MongoRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(client *mongo.Client) *MongoRepository {
	return &MongoRepository{
		collection: client.Database("gaming_db").Collection("jugadores"),
	}
}

func (r *MongoRepository) Guardar(jugador models.Jugador) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.InsertOne(ctx, jugador)
	return err
}

func (r *MongoRepository) ObtenerTodos() ([]models.Jugador, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var jugadores []models.Jugador

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var j models.Jugador
		if err := cursor.Decode(&j); err != nil {
			return nil, err
		}
		jugadores = append(jugadores, j)
	}

	return jugadores, nil
}

func (r *MongoRepository) ObtenerPorNombre(nombre string) (models.Jugador, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var jugador models.Jugador

	filtro := bson.M{"nombre": nombre}
	err := r.collection.FindOne(ctx, filtro).Decode(&jugador)

	return jugador, err
}
