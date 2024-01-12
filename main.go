package main

import (
	"log"

	"github.com/Xceptions/DocumentSearchEngineGo/handlers"
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	router := gin.Default()

	router.POST("/add", handlers.AddDocumentToDB)
	router.POST("/search", handlers.SearchForDocumentsContainingTerm)

	return router
}

// Input: None
// is the entrypoint for the application. It
// sets up a gin router and maps to various
// views, then runs the server on port 8080
func main() {
	router := setupRouter()

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Unable to start the server")
	}
}
