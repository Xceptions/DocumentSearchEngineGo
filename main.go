package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	router := gin.Default()

	router.POST("/add", addDocumentToDB)
	router.POST("/search", searchForDocumentsContainingTerm)

	return router
}

func main() {
	router := setupRouter()

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Unable to start the server")
	}
}
