package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// GetSecret config provider
func GetSecret(secret string) string {
	secret = strings.ToUpper(secret)
	mode := "dotenv"

	// @TODO
	// secret mode for docker secrets
	// for swarm production
	// if mode == "secret" {
	// 	path, err := filepath.Abs(config.GetString("secret.docker-secret-path"))
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	secretValue, err := ioutil.ReadFile(path + "/" + secret)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	return string(secretValue)
	// }

	// dotenv mode with .env file
	// in project root (Developmenmt only)
	if mode == "dotenv" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	// env mode with global environment variables
	// for production
	return os.Getenv(secret)
}
