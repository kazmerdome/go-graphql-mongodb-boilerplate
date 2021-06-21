package main

import (
	"go-graphql-mongodb-boilerplate/config"
	"go-graphql-mongodb-boilerplate/db"
	"go-graphql-mongodb-boilerplate/server"
	"go-graphql-mongodb-boilerplate/utility"
)

func main() {
	/*
	 * Start the program
	 */
	utility.ShowLogo()

	/**
	 * Connect to db
	 */
	defaultDB := db.Database{
		DataBaseRefName: "default",
		URL:             config.GetSecret("MONGO_URI"),
		DataBaseName:    config.GetSecret("MONGO_DATABASE"),
		RetryWrites:     config.GetSecret("MONGO_RETRYWRITES"),
	}
	defaultDB.Init()
	defer defaultDB.Disconnect()

	/**
	 * Create custom mongo indexes
	 */

	/*
	 * Start http server
	 */
	server.New()
}
