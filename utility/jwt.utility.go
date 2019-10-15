package utility

import (
	"aery-graphql/config"
	"errors"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// GenerateJWTToken ...
func GenerateJWTToken(ID string) (string, error) {
	sessionExpiration, err := strconv.ParseInt(config.GetSecret("SESSION_EXPIRATION"), 10, 32)
	sessionSecret := config.GetSecret("SESSION_SECRET")
	if err != nil {
		return "", err
	}

	expireToken := time.Now().Add(time.Hour * time.Duration(sessionExpiration)).Unix()

	claims := &jwt.StandardClaims{
		ExpiresAt: expireToken,
		IssuedAt:  time.Now().Unix(),
		Id:        ID,
		Issuer:    "aery-labs",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(sessionSecret))

	if err != nil {
		return "", err
	}

	return tokenString, err
}

// VerifyJWTToken ...
func VerifyJWTToken(tokenString string) (string, error) {
	var id string
	sessionSecret := config.GetSecret("SESSION_SECRET")

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
