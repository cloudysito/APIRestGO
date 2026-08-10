package models

type RankStats struct {
	Rank  string `json:"rank" bson:"_id"`
	Count int    `json:"count" bson:"count"`
}
