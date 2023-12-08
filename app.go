package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gin-gonic/gin"
)

// gets called from main.go as a POST request
// receives the gin context with the request data
// binds the request data to textToSave variable
// generates the id, calls the save_document and
// save_words function async
func addDocumentToDB(c *gin.Context) {
	var textToSave IdToDoc
	saveDocumentChannel := make(chan *mongo.InsertOneResult)
	saveWordsChannel := make(chan *mongo.InsertOneResult)

	if err := c.BindJSON(&textToSave); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	textToSave.ID = primitive.NewObjectID()

	// insert data in IdToDoc
	// it will have its document which is the text gotten from the
	// front end
	go func() {
		result, err := db.Collection("IdToDoc").InsertOne(c, &textToSave)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		saveDocumentChannel <- result
	}()

	// insert word to id
	go func() {
		// var bulkWriteOperations []string
		textToSaveDoc := strings.ToLower(textToSave.Document)
		words := strings.Split(textToSaveDoc, " ")
		wordsMap := make(map[string]string)
		for i := 0; i < len(wordsMap); i += 2 {
			wordsMap[words[i]] = words[i+1]
		}

		// return all collections containing the words in `words`
		WordToIdCursor, err := db.Collection("WordToId").Find(c, bson.M{"word": bson.M{"$in": words}})
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		var WordToIdCollection []WordToId
		if err = WordToIdCursor.All(c, &WordToIdCollection); err != nil {
			panic(err)
		}
		fmt.Println(WordToIdCollection)

		var bulkWriteOperations []mongo.WriteModel
		for _, values := range WordToIdCollection {
			values.IDs = append(values.IDs, textToSave.ID)
		}
		for _, values := range WordToIdCollection {
			toWrite := mongo.NewUpdateOneModel().SetUpdate(bson.M{values.Word: values.IDs})
			bulkWriteOperations = append(bulkWriteOperations, toWrite)
			delete(wordsMap, values.Word)
		}

		// saveWordsChannel <- WordToIdCollection
	}()

	select {
	case SaveDocumentChannelResult := <-saveDocumentChannel:
		c.JSON(http.StatusOK, SaveDocumentChannelResult)
	case SaveWordsChannelResult := <-saveWordsChannel:
		c.JSON(http.StatusOK, SaveWordsChannelResult)
	}
}

func searchForDocumentsContainingTerm(c *gin.Context) {
	var searchTerm SearchTerm
	if err := c.BindJSON(&searchTerm); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	inputText := strings.ToLower(searchTerm.Search)
	inputTextSplit := strings.Split(inputText, " ")

	allOccurrences := []primitive.ObjectID{}

	WordToIdCursor, err := db.Collection("WordToId").Find(c, bson.M{"word": bson.M{"$in": inputTextSplit}})
	if err != nil {
		fmt.Println("fmt - collection returned error")
		log.Println("log - collection returned error")
		panic(err)
	}

	var WordToIdCollection []WordToId
	if err = WordToIdCursor.All(c, &WordToIdCollection); err != nil {
		panic(err)
	}
	fmt.Println(WordToIdCollection)

	for _, values := range WordToIdCollection {
		allOccurrences = append(allOccurrences, values.IDs...)
	}

	fmt.Println(allOccurrences)
	c.JSON(http.StatusOK, allOccurrences)
}
