package model

import (
	"net/http"

	"go-graphql-mongodb-boilerplate/guard"
	"go-graphql-mongodb-boilerplate/utility"

	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

/**
 * MODEL
 */
type AuthToken struct {
	Token string `json:"token" bson:"token,omitempty"`
}

/**
 * BASIC STRATEGY
 */
// model
type AuthBasicStrategy struct {
	Email    string `json:"email" bson:"email,omitempty"`
	Username string `json:"username" bson:"username,omitempty"`
	Password string `json:"password" bson:"password,omitempty"`
}

// DTO
type AuthBasicStrategySignUpDTO struct {
	Email    string `json:"email" bson:"email,omitempty"`
	Username string `json:"username" bson:"username,omitempty"`
	Password string `json:"password" bson:"password,omitempty"`
}
type AuthBasicStrategySignInDTO struct {
	Username *string `json:"username" bson:"username,omitempty"`
	Password string  `json:"password" bson:"password,omitempty"`
}

// Operations
func (a *AuthBasicStrategy) SignUp() (AuthToken, error) {
	token := AuthToken{}

	// validate empty data
	if a.Email == "" || a.Password == "" {
		return token, echo.NewHTTPError(http.StatusUnauthorized)
	}

	// Salt and hash the password using the bcrypt algorithm
	// The second argument is the cost of hashing, which we arbitrarily set as 8 (this value can be more or less, depending on the computing power you wish to utilize)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(a.Password), 8)
	if err != nil {
		return token, err
	}

	// Create new user
	newUser := User{}
	userCreateDTO := UserCreateDTO{
		Email:    a.Email,
		Username: a.Username,
		Password: string(hashedPassword),
		Role:     guard.User,
	}
	if err := newUser.Create(&userCreateDTO); err != nil {
		return token, err
	}
	if utility.IsZeroVal(newUser.ID) {
		return token, echo.NewHTTPError(http.StatusUnauthorized)
	}

	// Create jwt token
	jwt, err := utility.GenerateJWTToken(newUser.ID.Hex())
	if err != nil {
		return token, echo.NewHTTPError(http.StatusUnauthorized)
	}

	token.Token = jwt
	return token, nil
}

func (a *AuthBasicStrategy) SignIn() (AuthToken, error) {
	token := AuthToken{}

	// Check user by email OR username
	where := UserWhereDTO{}
	storedUser := User{}
	if a.Username == "" {
		return token, echo.NewHTTPError(http.StatusUnauthorized)
	}
	if a.Username != "" {
		where.Username = &a.Username
	}
	err := storedUser.One(&where)
	if err != nil {
		return token, err
	}
	if utility.IsZeroVal(storedUser.ID) {
		return token, echo.NewHTTPError(http.StatusUnauthorized)
	}

	// If user exists, get db hashed password
	password := storedUser.Password
	if password == "" {
		return token, echo.NewHTTPError(http.StatusUnauthorized)
	}

	// Compare the stored hashed password, with the hashed version of the password that was received
	if err = bcrypt.CompareHashAndPassword([]byte(password), []byte(a.Password)); err != nil {
		return token, echo.NewHTTPError(http.StatusUnauthorized)
	}

	// Create jwt token
	jwt, err := utility.GenerateJWTToken(storedUser.ID.Hex())
	if err != nil {
		return token, echo.NewHTTPError(http.StatusUnauthorized)
	}

	token.Token = jwt
	return token, nil
}
