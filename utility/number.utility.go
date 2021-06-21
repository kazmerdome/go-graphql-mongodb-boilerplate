package utility

import (
	"math/rand"
	"time"
)

// CreateRandomNumber ...
func CreateRandomNumber(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(max-min+1) + min
	return randomNumber
}
