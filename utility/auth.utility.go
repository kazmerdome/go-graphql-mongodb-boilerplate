package utility

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"go-graphql-mongodb-boilerplate/config"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/**
 * BEARER TOKEN
 */
func GetObjectIDFromBearerToken(bearer string) (primitive.ObjectID, error) {
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

/**
 * JWT
 */
func GenerateJWTToken(ID string) (string, error) {
	sessionExpiration, err := strconv.ParseInt(config.GetSecret("JWT_SESSION_EXPIRATION"), 10, 32)
	sessionSecret := config.GetSecret("JWT_SESSION_SECRET")
	if err != nil {
		return "", err
	}

	expireToken := time.Now().Add(time.Hour * time.Duration(sessionExpiration)).Unix()

	claims := &jwt.StandardClaims{
		ExpiresAt: expireToken,
		IssuedAt:  time.Now().Unix(),
		Id:        ID,
		Issuer:    "kazmerdome",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(sessionSecret))

	if err != nil {
		return "", err
	}

	return tokenString, err
}
func VerifyJWTToken(tokenString string) (string, error) {
	var id string
	sessionSecret := config.GetSecret("JWT_SESSION_SECRET")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(sessionSecret), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["jti"].(string) == "" {
			return "", errors.New("invalid stored token")
		}
		id = claims["jti"].(string)
	} else {
		return "", errors.New("invalid stored token")
	}

	return id, nil
}
