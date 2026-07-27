package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Item struct {
	Name   string `json:"name" bson:"name"`
	Type   string `json:"type" bson:"type"`
	Rarity string `json:"rarity" bson:"rarity"`
}

type Player struct {
	Name      string             `json:"name" bson:"name"`
	Rank      string             `json:"rank" bson:"rank"`
	Inventory []Item             `json:"inventory,omitempty" bson:"inventory,omitempty"`
	FactionID primitive.ObjectID `json:"faction_id,omitempty" bson:"faction_id,omitempty"`
}
