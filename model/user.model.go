package model

import (
	"context"
	"errors"
	"time"

	"go-graphql-mongodb-boilerplate/db"
	"go-graphql-mongodb-boilerplate/guard"
	"go-graphql-mongodb-boilerplate/utility"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/**
 * DB Info
 */
const DB_COLLECTION_NAME__USER = "User"
const DB_REF_NAME__USER = "default"

/**
 * SEARCH regex fields
 */
var SEARCH_FILEDS__USER = []string{"email", "username"}

/**
 * MODEL
 */
type User struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Email     string             `json:"email" bson:"email,omitempty" validate:"email,required"`
	Username  string             `json:"username" bson:"username,omitempty" validate:"required"`
	Password  string             `json:"password" bson:"password,omitempty"`
	Role      guard.Role         `json:"role" bson:"role,omitempty"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

/**
 * ENUM
 */
type UserOrderByENUM string

/**
 * DTO
 */

// Read
type UserWhereUniqueDTO struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
}
type UserWhereDTO struct {
	ID       *primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Email    *string             `json:"email" bson:"email,omitempty"`
	Username *string             `json:"username" bson:"username,omitempty"`
	Role     *guard.Role         `json:"role" bson:"role,omitempty"`
	OR       []bson.M            `json:"$or,omitempty" bson:"$or,omitempty"`
}

// Write
type UserCreateDTO struct {
	ID        *primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Email     string              `json:"email" bson:"email,omitempty"`
	Username  string              `json:"username" bson:"username,omitempty"`
	Password  string              `json:"password" bson:"password,omitempty"`
	Role      guard.Role          `json:"role" bson:"role,omitempty"`
	CreatedAt time.Time           `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time           `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
type UserUpdateDTO struct {
	Email     string      `json:"email" bson:"email,omitempty"`
	Username  string      `json:"username" bson:"username,omitempty"`
	Password  *string     `json:"password" bson:"password,omitempty"`
	Role      *guard.Role `json:"role" bson:"role,omitempty"`
	CreatedAt time.Time   `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time   `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

/**
 * OPERATIONS
 */

// Read

// One
func (u *User) One(filter *UserWhereDTO) error {
	collection := db.GetCollection(DB_COLLECTION_NAME__USER, DB_REF_NAME__USER)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection.FindOne(ctx, &filter).Decode(&u)

	return nil
}

// List
func (u *User) List(filter *UserWhereDTO, orderBy *UserOrderByENUM, skip *int, limit *int, customQuery *bson.M) ([]*User, error) {
	var items []*User
	orderByKey := "created_at"
	orderByValue := -1
	collection := db.GetCollection(DB_COLLECTION_NAME__USER, DB_REF_NAME__USER)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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

	var queryFilter interface{}
	if filter != nil {
		queryFilter = filter
	}
	if !utility.IsZeroVal(customQuery) {
		queryFilter = customQuery
	}

	cursor, err := collection.Find(ctx, &queryFilter, options)
	if err != nil {
		return items, err
	}
	err = cursor.All(ctx, &items)
	if err != nil {
		return items, err
	}

	return items, nil
}

// Count
func (u *User) Count(filter *UserWhereDTO) (int, error) {
	collection := db.GetCollection(DB_COLLECTION_NAME__USER, DB_REF_NAME__USER)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	count, err := collection.CountDocuments(ctx, filter, nil)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

// Write Operations

// Create
func (r *User) Create(data *UserCreateDTO) error {
	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()
	// validate
	if err := utility.ValidateStruct(data); err != nil {
		return err
	}
	// collection
	collection := db.GetCollection(DB_COLLECTION_NAME__USER, DB_REF_NAME__USER)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// check uniqe
	item := new(User)
	f := bson.M{
		"$or": []bson.M{
			{"email": data.Email},
			{"username": data.Username},
		},
	}
	collection.FindOne(ctx, f).Decode(&item)
	if item.Email != "" {
		return errors.New("user is already exist")
	}
	// operation
	res, err := collection.InsertOne(ctx, data)
	if err != nil {
		return err
	}
	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return errors.New("server error")
	}
	collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&r)
	return nil
}

// Update
func (r *User) Update(where primitive.ObjectID, data *UserUpdateDTO) error {
	data.UpdatedAt = time.Now()
	// validate
	if utility.IsZeroVal(where) {
		return errors.New("internal server error")
	}
	if err := utility.ValidateStruct(data); err != nil {
		return err
	}
	// collection
	collection := db.GetCollection(DB_COLLECTION_NAME__USER, DB_REF_NAME__USER)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// check user is exists
	collection.FindOne(ctx, bson.M{"_id": where}).Decode(&r)
	if r.Email == "" {
		return errors.New("item not found")
	}
	// check unique
	item := new(User)
	f := bson.M{
		"$or": []bson.M{
			{"email": data.Email, "_id": bson.M{"$ne": where}},
			{"username": data.Username, "_id": bson.M{"$ne": where}},
		},
	}
	collection.FindOne(ctx, f).Decode(&item)
	if item.Email != "" {
		return errors.New("user is already exist")
	}
	// operation
	_, err := collection.UpdateOne(ctx, bson.M{"_id": where}, bson.M{"$set": data})
	collection.FindOne(ctx, bson.M{"_id": where}).Decode(&r)
	if err != nil {
		return err
	}
	return nil
}

// Delete
func (r *User) Delete() error {
	collection := db.GetCollection(DB_COLLECTION_NAME__USER, DB_REF_NAME__USER)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if utility.IsZeroVal(r.ID) {
		return errors.New("invalid id")
	}
	collection.FindOne(ctx, bson.M{"_id": r.ID}).Decode(&r)
	if r.Email == "" {
		return errors.New("item not found")
	}
	_, err := collection.DeleteOne(ctx, bson.M{"_id": r.ID})
	if err != nil {
		return err
	}
	return nil
}
