package utility

import (
	"errors"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ParseBearerToken Bearer token to userID
func ParseBearerToken(bearer string) (primitive.ObjectID, error) {
	rawTokenParts := strings.Split(bearer, "Bearer ")
	if len(rawTokenParts) < 2 {
		return primitive.ObjectID{}, errors.New("invalid header token")
	}
	// verify jwt token string
	userIDHex, err := VerifyJWTToken(rawTokenParts[1])
	if err != nil {
		return primitive.ObjectID{}, err
	}
	// craate primitive.ObjectID
	userID, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		return primitive.ObjectID{}, errors.New("unauthorized")
	}
	return userID, nil
}
