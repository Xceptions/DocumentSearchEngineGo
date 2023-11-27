package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.POST("/add", addDocumentToDB)

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Unable to start the server")
	}
}
