package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IdToDoc struct {
	ID       primitive.ObjectID `json:"id"`
	Document string             `json:"document"`
}

type WordToId struct {
	Word string `json:"word"`
	IDs  []primitive.ObjectID
}

type SearchTerm struct {
	Search string `json:"search"`
}
