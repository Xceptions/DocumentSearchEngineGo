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


// Returns a list of the strings that contain the searched word.
// The function receives a Search term client
// It then retrieves the documents in our DB collection (WordToId),
// containing that word, gets their corresponding ids and returns
// the document (from IdToDoc) using the ids.
// The function performs the following steps:
//     - receives a user input (str) from the client
//     - creates an allOccurrence list to hold a list of the ids
//         that contain the words in Search term
//     - splits the Search term into a list of lowercase words called
//         `words` for searching the db
//     - retrieves a collection of documents containing words in `words`
//         from the WordToId collection
//     - for all the documents in the collection above, it extends their
//         `ids` field into the allOccurrence. This `ids` field is a
//         a list containing the _id of the documents that contain these
//         words
//     - creates an `IdToDocCollection` to convert the allOccurrence
//         List[str] to List[str]. This is needed because that
//         is the datatype stored in the MongoDB collection
//     - retrieves a collection of documents containing ObjectId(ids) in
//         `IdToDocCollections`
//     - returns a list of document(str) contained in documents(MongoDB)
//         above or empty list
// Args:
//     Search term(str) - a string sent from the client through the gin context
//         with the intention of searching our DB for it.
// Returns
//     List[str]: a list of the documents that contain the words in Search term
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
