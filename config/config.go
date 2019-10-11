package config

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

var config *viper.Viper

// Init is an exported method that takes the environment starts the viper
// (external lib) and returns the configuration struct.
func Init(env string) {
	var err error
	config = viper.New()
	config.Set("env", env)
	config.SetConfigName(env)
	config.SetConfigType("yaml")
	config.AddConfigPath("../config/")
	config.AddConfigPath("config/")

	err = config.ReadInConfig()
	if err != nil {
		log.Fatal("error on parsing configuration file")
	}
}

// GetConfig config provider
func GetConfig() *viper.Viper {
	return config
}

// GetSecret config provider
func GetSecret(secret string) string {
	secret = strings.ToUpper(secret)
	mode := config.GetString("secret.mode")

	// @TODO
	// secret mode for docker secrets
	// for swarm production
	if mode == "secret" {
		path, err := filepath.Abs(config.GetString("secret.docker-secret-path"))
		if err != nil {
			log.Fatal(err)
		}

		secretValue, err := ioutil.ReadFile(path + "/" + secret)
		if err != nil {
			log.Fatal(err)
		}

		return string(secretValue)
	}

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
