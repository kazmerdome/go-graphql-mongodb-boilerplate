package guard

import (
	"aery-graphql/db"
	"aery-graphql/utility"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

// Role type
type Role string

// Resource ...
type Resource string

// User and other roles
const (
	resource Resource = "BLACKDOME_SERVER"
	User     Role     = "USER"
	Partner  Role     = "PARTNER"
	Editor   Role     = "EDITOR"
	Admin    Role     = "ADMIN"
)

// AuthUser ...
type AuthUser struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	GoogleID  string             `json:"google_id" bson:"google_id,omitempty"`
	Email     string             `json:"email" bson:"email,omitempty"`
	AppPolicy []AppPolicy        `json:"app_policy" bson:"app_policy,omitempty"`
}

// AppPolicy ...
type AppPolicy struct {
	Resource Resource `json:"resource" bson:"resource,omitempty"`
	Role     Role     `json:"role" bson:"role,omitempty"`
}

/**
 * HELPERS
 */

// GetUserData from db user collection
func (a *AuthUser) getUserData(userID *primitive.ObjectID) error {
	collection := db.GetCollection("User")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := collection.FindOne(ctx, bson.M{"_id": userID}).Decode(&a)
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

	// check []AppPolicy is exist, if not it means user is also not exist || invalid stored user
	if u.AppPolicy == nil && len(u.AppPolicy) < 1 {
		return User
	}

	// get current role from user []AppPolicy
	for _, policy := range u.AppPolicy {
		if policy.Resource == resource {
			userCurrentRole = policy.Role
		}
	}

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
