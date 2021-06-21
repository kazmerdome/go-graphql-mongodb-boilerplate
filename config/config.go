package config

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

var MEMORY_STORE = make(map[string]string)

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

	// from memory strore - if the value is already stored in memory
	if MEMORY_STORE[secret] != "" {
		return MEMORY_STORE[secret]
	}

	// from source
	// 1. parse .env file if exists
	if isFileExists(".env") {
		godotenv.Load()
	}
	// 2. get env variable
	value := os.Getenv(secret)
	// 3. check if docker secret mode
	// in production image, the workplace is the /run folder
	// the docker secrets live in run/secrets folder
	if value == "***DOCKER_SECRET***" {
		path, err := filepath.Abs("./secrets")
		if err != nil {
			log.Fatal(err)
		}

		secretFile := path + "/" + secret
		secretValue, err := ioutil.ReadFile(secretFile)
		if err != nil {
			log.Fatal(err)
		}
		s := strings.TrimSpace(string(secretValue))
		value = string(s)
	}

	// return and store value in memory
	MEMORY_STORE[secret] = value
	return value
}
