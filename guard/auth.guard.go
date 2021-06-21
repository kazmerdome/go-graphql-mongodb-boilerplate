package guard

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go-graphql-mongodb-boilerplate/db"
	"go-graphql-mongodb-boilerplate/utility"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Role type
type Role string

// Resource ...
type Resource string

// User and other roles
var (
	User   Role = "USER"
	Editor Role = "EDITOR"
	Admin  Role = "ADMIN"
)

// AuthUser ...
type AuthUser struct {
	ID    primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Email string             `json:"email" bson:"email,omitempty"`
	Role  Role               `json:"role" bson:"role,omitempty"`
}

/**
 * HELPERS
 */

// GetUserData from db user collection
const DB_COLLECTION_NAME__USER = "User"
const DB_REF_NAME__USER = "default"

func (a *AuthUser) getUserData(userID *primitive.ObjectID) error {
	collection := db.GetCollection(DB_COLLECTION_NAME__USER, DB_REF_NAME__USER)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": userID}
	err := collection.FindOne(ctx, &filter).Decode(&a)

	if err != nil {
		return err
	}
	return nil
}

// role checker
func checkRole(userRole Role, requiredRoles []Role) bool {
	isContains := false
	for _, reqiredRole := range requiredRoles {
		if reqiredRole == userRole {
			isContains = true
		}
	}
	return isContains
}

// GetUserCurrentRole get user current role on a specific resource the default role is User
func getUserCurrentRole(userID primitive.ObjectID) Role {
	userCurrentRole := User

	// fetch user from db
	u := AuthUser{}
	if err := u.getUserData(&userID); err != nil {
		return User
	}

	// check Role is exist, if not it means user is also not exist || invalid stored user
	if u.Role == "" {
		return User
	}

	// get current role from user role
	userCurrentRole = u.Role

	return userCurrentRole
}

/**
 * GUARDS
 */

// GetRole Guard
func GetRole(bearerToken string) Role {
	if bearerToken == "" {
		return User
	}
	userID, err := utility.GetObjectIDFromBearerToken(bearerToken)
	if err != nil {
		return User
	}
	userCurrentRole := getUserCurrentRole(userID)
	return userCurrentRole
}

// GetUserID Guard
func GetUserID(bearerToken string) (primitive.ObjectID, error) {
	if bearerToken == "" {
		return primitive.ObjectID{}, fmt.Errorf("Access denied")
	}
	userID, err := utility.GetObjectIDFromBearerToken(bearerToken)
	if err != nil {
		return primitive.ObjectID{}, fmt.Errorf("Access denied")
	}
	return userID, nil
}

// GetUserData Guard
func GetUserData(bearerToken string) (AuthUser, error) {
	u := AuthUser{}
	if bearerToken == "" {
		return u, fmt.Errorf("Access denied")
	}
	// get object ID from token
	userID, err := utility.GetObjectIDFromBearerToken(bearerToken)
	if err != nil {
		return u, fmt.Errorf("Access denied")
	}
	// fetch user from db

	if err := u.getUserData(&userID); err != nil {
		return u, fmt.Errorf("Access denied")
	}
	return u, nil
}

// Auth Guard
func Auth(requiredRoles []Role, bearerToken string) error {
	if bearerToken == "" && checkRole(User, requiredRoles) {
		return nil
	}
	userID, err := utility.GetObjectIDFromBearerToken(bearerToken)
	if err != nil {
		return err
	}
	userCurrentRole := getUserCurrentRole(userID)
	if checkRole(userCurrentRole, requiredRoles) {
		return nil
	}
	return errors.New("unauthorized")
}
