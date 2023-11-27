package main

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IdToDoc struct {
	ID       primitive.ObjectID `json:"id"`
	Document string             `json:"document"`
}

// type IdToDoc struct {
// 	ID       primitive.ObjectID `json:"_id" bson:"_id"`
// 	Document string             `json:"title"`
// }
