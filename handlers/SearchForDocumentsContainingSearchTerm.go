package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Xceptions/DocumentSearchEngineGo/database"
	"github.com/Xceptions/DocumentSearchEngineGo/models"
)

func SearchForDocumentsContainingTerm(c *gin.Context) {
	var toSearch models.SearchTerm
	var db *mongo.Database

	db = database.ConnectDB()
	
	if err := c.BindJSON(&toSearch); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	inputText := strings.ToLower(toSearch.Search)
	inputTextSplit := strings.Split(inputText, " ")

	allOccurrences := []primitive.ObjectID{}

	WordToIdCursor, err := db.Collection("WordToId").Find(c, bson.M{"word": bson.M{"$in": inputTextSplit}})
	if err != nil {
		log.Println("log - collection returned error")
		panic(err)
	}

	var WordToIdCollection []models.WordToId
	if err = WordToIdCursor.All(c, &WordToIdCollection); err != nil {
		panic(err)
	}

	for _, values := range WordToIdCollection {
		allOccurrences = append(allOccurrences, values.IDs...)
	}

	IdToDocCursor, err := db.Collection("IdToDoc").Find(c, bson.M{"id": bson.M{"$in": allOccurrences}})
	if err != nil {
		log.Println("log - collection returned error")
		panic(err)
	}

	var IdToDocCollection []models.IdToDoc
	if err = IdToDocCursor.All(c, &IdToDocCollection); err != nil {
		panic(err)
	}

	var documents []string

	for _, elements := range IdToDocCollection {
		documents = append(documents, elements.Document)
	}
	c.JSON(http.StatusOK, documents)
}
