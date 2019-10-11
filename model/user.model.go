package model

import (
	"aery-graphql/db"
	"aery-graphql/guard"
	"aery-graphql/utility"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UserCollectionName ...
const UserCollectionName = "User"

// UserSearchFields ...
var UserSearchFields = []string{"firstname", "lastname", "email"}

// User ...
type User struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Firstname string             `json:"firstname" bson:"firstname,omitempty"`
	Lastname  string             `json:"lastname" bson:"lastname,omitempty"`
	GoogleID  string             `json:"google_id" bson:"google_id,omitempty"`
	Email     string             `json:"email" bson:"email,omitempty"`
	AppPolicy []AppPolicy        `json:"app_policy" bson:"app_policy,omitempty"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

// AppPolicy ...
type AppPolicy struct {
	Resource AppPolicyResource `json:"resource" bson:"resource,omitempty"`
	Role     guard.Role        `json:"role" bson:"role,omitempty"`
}

// UserWhereUniqueInput ...
type UserWhereUniqueInput struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
}

// UserWhereInput ...
type UserWhereInput struct {
	ID        *primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Firstname *string             `json:"firstname" bson:"firstname,omitempty"`
	Lastname  *string             `json:"lastname" bson:"lastname,omitempty"`
	Email     *string             `json:"email" bson:"email,omitempty"`
	AppPolicy *AppPolicy          `json:"app_policy" bson:"app_policy,omitempty"`
	OR        []bson.M            `json:"$or,omitempty" bson:"$or,omitempty"`
}

// AppPolicyResource ...
type AppPolicyResource string

// UserOrderByInput ...
type UserOrderByInput string

// Create ...
func (u *User) Create() error {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()

	collection := db.GetCollection(UserCollectionName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	item := new(User)
	collection.FindOne(ctx, bson.M{"email": u.Email}).Decode(&item)
	if item.Email != "" {
		return errors.New("user name_slug is already exist")
	}

	res, err := collection.InsertOne(ctx, u)
	if err != nil {
		return err
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		u.ID = oid
	}

	return nil
}

// One ...
func (u *User) One(filter *UserWhereInput) error {
	collection := db.GetCollection(UserCollectionName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	collection.FindOne(ctx, &filter).Decode(&u)
	return nil
}

// List ...
func (u *User) List(filter *UserWhereInput, orderBy *UserOrderByInput, skip *int, limit *int) ([]*User, error) {
	var items []*User
	orderByKey := "created_at"
	orderByValue := -1
	collection := db.GetCollection(UserCollectionName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	options := options.Find()
	if limit != nil {
		options.SetLimit(int64(*limit))
	}
	if skip != nil {
		options.SetSkip(int64(*skip))
	}
	if orderBy != nil {
		orderByKey, orderByValue = utility.GetOrderByKeyAndValue(string(*orderBy))
	}
	options.SetSort(map[string]int{orderByKey: orderByValue})

	cursor, err := collection.Find(ctx, filter, options)
	cursor.All(ctx, &items)

	if err != nil {
		return items, err
	}
	return items, nil
}

// Update ...
func (u *User) Update() error {
	u.UpdatedAt = time.Now()

	collection := db.GetCollection(UserCollectionName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	item := new(User)
	collection.FindOne(ctx, bson.M{"email": u.Email, "_id": bson.M{"$ne": u.ID}}).Decode(&item)
	if item.Email != "" {
		return errors.New("user email is already exist")
	}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": u.ID}, bson.M{"$set": u})
	collection.FindOne(ctx, bson.M{"_id": u.ID}).Decode(&u)

	if err != nil {
		return err
	}

	return nil
}

// Delete ...
func (u *User) Delete() error {
	collection := db.GetCollection(UserCollectionName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	collection.FindOne(ctx, bson.M{"_id": u.ID}).Decode(&u)

	if u.Email == "" {
		return errors.New("user doesn't exist")
	}

	_, err := collection.DeleteOne(ctx, bson.M{"_id": u.ID})
	if err != nil {
		return err
	}

	return nil
}
