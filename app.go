package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

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
	saveWordsChannel := make(chan *mongo.BulkWriteResult)

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
		wordsMap := make(map[string]int)
		for i := 0; i < len(words); i += 1 {
			wordsMap[words[i]] = 0
		}
		fmt.Println(wordsMap)

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

		// iterate through the documents in the collections, add the new id
		// to the IDs, then prepare them in an updateone statement and add to
		// bulkwrite operation, waiting to write
		var bulkWriteOperations []mongo.WriteModel

		if len(WordToIdCollection) > 0 {
			for _, values := range WordToIdCollection {
				newValuesIDs := append(values.IDs, textToSave.ID)
				toWrite := mongo.NewUpdateOneModel().SetUpdate(bson.M{values.Word: newValuesIDs})
				bulkWriteOperations = append(bulkWriteOperations, toWrite)
				delete(wordsMap, values.Word)
			}
		}
		for word := range wordsMap {
			newValuesIDs := []primitive.ObjectID{textToSave.ID}
			toWrite := mongo.NewInsertOneModel().SetDocument(WordToId{Word: word, IDs: newValuesIDs})
			bulkWriteOperations = append(bulkWriteOperations, toWrite)
		}
		fmt.Println(bulkWriteOperations)

		opts := options.BulkWrite().SetOrdered(false)

		results, err := db.Collection("WordToId").BulkWrite(c, bulkWriteOperations, opts)

		if err != nil {
			panic(err)
		}

		saveWordsChannel <- results
	}()

	select {
	case SaveDocumentChannelResult := <-saveDocumentChannel:
		c.JSON(http.StatusOK, SaveDocumentChannelResult)
	case SaveWordsChannelResult := <-saveWordsChannel:
		c.JSON(http.StatusOK, SaveWordsChannelResult)
	}
}

func searchForDocumentsContainingTerm(c *gin.Context) {
	var toSearch SearchTerm
	if err := c.BindJSON(&toSearch); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	fmt.Println(toSearch)
	inputText := strings.ToLower(toSearch.Search)
	fmt.Println(inputText)
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

	for _, values := range WordToIdCollection {
		allOccurrences = append(allOccurrences, values.IDs...)
	}

	IdToDocCursor, err := db.Collection("IdToDoc").Find(c, bson.M{"id": bson.M{"$in": allOccurrences}})
	if err != nil {
		log.Println("log - collection returned error")
		panic(err)
	}

	var IdToDocCollection []IdToDoc
	if err = IdToDocCursor.All(c, &IdToDocCollection); err != nil {
		panic(err)
	}

	fmt.Println(IdToDocCollection)
	var documents []string

	for _, elements := range IdToDocCollection {
		documents = append(documents, elements.Document)
	}
	c.JSON(http.StatusOK, documents)
}
