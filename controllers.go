package main

import (
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gin-gonic/gin"
)

func addDocumentToDB(c *gin.Context) {
	var textToSave IdToDoc
	if err := c.BindJSON(&textToSave); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	textToSave.ID = primitive.NewObjectID()
	fmt.Println("got here")
	result, err := db.Collection("IdToDoc").InsertOne(c, &textToSave)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	fmt.Println("the result is")
	fmt.Println(result)
	c.JSON(http.StatusOK, textToSave)
}
