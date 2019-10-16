package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func isFileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// GetSecret config provider
func GetSecret(secret string) string {
	secret = strings.ToUpper(secret)

	// 1. parse .env file if it exists
	if isFileExists(".env") {
		godotenv.Load()
	}

	// 2. get env variable
	env := os.Getenv(secret)

	// 3. check if env == ***DOCKER_SECRET***
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

	return env
}
