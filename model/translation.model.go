package model

import (
	"aery-graphql/db"
	"aery-graphql/utility"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TranslationCollectionName ...
const TranslationCollectionName = "Translation"

// TranslationSearchFields ...
var TranslationSearchFields = []string{"key", "value"}

// Translation ...
type Translation struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Locale    Locale             `json:"locale" bson:"locale,omitempty"`
	Key       string             `json:"key" bson:"key,omitempty"`
	Value     string             `json:"value" bson:"value,omitempty"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

// TranslationWhereUniqueInput ...
type TranslationWhereUniqueInput struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
}

// TranslationWhereInput ...
type TranslationWhereInput struct {
	ID     *primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Locale *Locale             `json:"locale" bson:"locale,omitempty"`
	Key    *string             `json:"key" bson:"key,omitempty"`
	Value  *string             `json:"value" bson:"value,omitempty"`
	OR     []bson.M            `json:"$or,omitempty" bson:"$or,omitempty"`
}

// TranslationOrderByInput ...
type TranslationOrderByInput string

// Create ...
func (t *Translation) Create() error {
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()

	collection := db.GetCollection(TranslationCollectionName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	item := new(Translation)
	collection.FindOne(ctx, bson.M{"locale": t.Locale, "key": t.Key}).Decode(&item)
	if item.Key != "" {
		return errors.New("translation key is already exist in this locale")
	}

	res, err := collection.InsertOne(ctx, t)
	if err != nil {
		return err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		t.ID = oid
	}
	return nil
}

// One ...
func (t *Translation) One(filter *TranslationWhereInput) error {
	collection := db.GetCollection(TranslationCollectionName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	collection.FindOne(ctx, &filter).Decode(&t)
	return nil
}

// List ...
func (t *Translation) List(filter *TranslationWhereInput, orderBy *TranslationOrderByInput, skip *int, limit *int) ([]*Translation, error) {
	var items []*Translation
	orderByKey := "created_at"
	orderByValue := -1
	collection := db.GetCollection(TranslationCollectionName)
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
func (t *Translation) Update() error {
	t.UpdatedAt = time.Now()

	collection := db.GetCollection(TranslationCollectionName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	item := new(Translation)
	collection.FindOne(ctx, bson.M{"locale": t.Locale, "key": t.Key, "_id": bson.M{"$ne": t.ID}}).Decode(&item)

	if item.Key != "" {
		return errors.New("Translation key is already exist in this locale")
	}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": t.ID}, bson.M{"$set": t})
	collection.FindOne(ctx, bson.M{"_id": t.ID}).Decode(&t)

	if err != nil {
		return err
	}

	return nil
}

// Delete ...
func (t *Translation) Delete() error {
	collection := db.GetCollection(TranslationCollectionName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	collection.FindOne(ctx, bson.M{"_id": t.ID}).Decode(&t)

	if t.Key == "" {
		return errors.New("item doesn't exist")
	}

	_, err := collection.DeleteOne(ctx, bson.M{"_id": t.ID})
	if err != nil {
		return err
	}

	return nil
}
