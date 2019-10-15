package db

import (
	"aery-graphql/config"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var database *mongo.Database

// GetCollection ...
func GetCollection(collectionName string) *mongo.Collection {
	return database.Collection(collectionName)
}

// Init ...
func Init() {
	url := config.GetSecret("MONGO_URI")
	dataBaseName := config.GetSecret("MONGO_COLLECTION")
	retryWrites := config.GetSecret("MONGO_RETRYWRITES")
	connectionURI := url + "/" + dataBaseName + "?retryWrites=" + retryWrites

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionURI))

	if err != nil {
		fmt.Println(err)
		log.Fatal("Mongo connection error!")
	}

	database = client.Database(dataBaseName)
}
