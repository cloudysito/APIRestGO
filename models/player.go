package models

type Item struct {
	Name   string `json:"name" bson:"name"`
	Type   string `json:"type" bson:"type"`
	Rarity string `json:"rarity" bson:"rarity"`
}

type Player struct {
	Name      string `json:"name" bson:"name"`
	Rank      string `json:"rank" bson:"rank"`
	Inventory []Item `json:"inventory" bson:"inventory"`
}
