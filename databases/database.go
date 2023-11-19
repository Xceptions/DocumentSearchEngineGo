package databases

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// initializing function to connect to the database
// returns `client`, the connection object
func ConnectDB() *mongo.Client {
	Mongo_URL := "mongodb://127.0.0.1:27017"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(Mongo_URL))

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to mongoDB")
	return client
}
