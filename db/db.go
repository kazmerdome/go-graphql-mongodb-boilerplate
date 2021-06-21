package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var databases = make(map[string]*mongo.Database)

// GetCollection ...
func GetCollection(collectionName string, dataBaseRefName string) *mongo.Collection {
	db := databases[dataBaseRefName]
	if db == nil {
		log.Fatal("Invalid DatabaseRefName [", dataBaseRefName, "]")
	}
	return db.Collection(collectionName)
}

// DbConfig ...
type Database struct {
	URL             string
	DataBaseName    string
	RetryWrites     string
	DataBaseRefName string
}

func (db Database) Init() {
	fmt.Println("\nconnecting \033[0;36m", db.DataBaseRefName, "\033[0m db...")
	connectionURI := db.URL + "/" + db.DataBaseName + "?retryWrites=" + db.RetryWrites

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionURI))

	if err != nil {
		fmt.Println(err)
		log.Fatal("Mongo connection error!")
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	databases[db.DataBaseRefName] = client.Database(db.DataBaseName)
	fmt.Printf("\033[1;32mconnected successfully!\033[0m\n")
}

func (db Database) Disconnect() {
	fmt.Println("\ndisconnecting \033[0;36m", db.DataBaseRefName, "\033[0m db...")

	closingDB := databases[db.DataBaseRefName]
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	closingDB.Client().Disconnect(ctx)

	fmt.Printf("\033[1;32mdisconnected successfully!\033[0m\n")
}
