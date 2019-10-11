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

// GeneralCollectionName ...
const GeneralCollectionName = "General"

// General ...
type General struct {
	ID         primitive.ObjectID  `json:"_id,omitempty" bson:"_id,omitempty"`
	Locale     Locale              `json:"locale" bson:"locale,omitempty"`
	Seo        Seo                 `json:"seo" bson:"seo,omitempty"`
	Info       []GeneralInfo       `json:"info" bson:"info,omitempty"`
	Social     []GeneralSocial     `json:"social" bson:"social,omitempty"`
	Navigation []GeneralNavigation `json:"navigation" bson:"navigation,omitempty"`
	Homepage   string              `json:"homepage" bson:"homepage,omitempty"`
	CreatedAt  time.Time           `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt  time.Time           `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

// GeneralInfo ...
type GeneralInfo struct {
	Type  GeneralInfoType `json:"type" bson:"type,omitempty"`
	Value string          `json:"value" bson:"value,omitempty"`
}

// GeneralSocial ...
type GeneralSocial struct {
	Type  GeneralSocialType `json:"type" bson:"type,omitempty"`
	Value string            `json:"value" bson:"value,omitempty"`
}

// GeneralNavigation ...
type GeneralNavigation struct {
	Type GeneralNavigationType `json:"type" bson:"type,omitempty"`
	Menu primitive.ObjectID    `json:"menu" bson:"menu,omitempty"`
}

// GeneralWhereUniqueInput ...
type GeneralWhereUniqueInput struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
}

// GeneralWhereInput ...
type GeneralWhereInput struct {
	ID     *primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Locale *Locale             `json:"locale" bson:"locale,omitempty"`
}

// GeneralInfoType ...
type GeneralInfoType string

// GeneralSocialType ...
type GeneralSocialType string

// GeneralNavigationType ...
type GeneralNavigationType string

// GeneralOrderByInput ...
type GeneralOrderByInput string

// Create ...
func (g *General) Create() error {
	g.CreatedAt = time.Now()
	g.UpdatedAt = time.Now()

	collection := db.GetCollection(GeneralCollectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	item := new(General)
	collection.FindOne(ctx, bson.M{"locale": g.Locale}).Decode(&item)
	if item.Locale != "" {
		return errors.New("general is already exist in this locale")
	}

	res, err := collection.InsertOne(ctx, g)
	if err != nil {
		return err
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		g.ID = oid
	}

	return nil
}

// One ...
func (g *General) One(filter *GeneralWhereInput) error {
	collection := db.GetCollection(GeneralCollectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection.FindOne(ctx, &filter).Decode(&g)
	return nil
}

// List ...
func (g *General) List(filter *GeneralWhereInput, orderBy *GeneralOrderByInput, skip *int, limit *int) ([]*General, error) {
	var items []*General
	orderByKey := "created_at"
	orderByValue := -1
	collection := db.GetCollection(GeneralCollectionName)
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

	cursor, err := collection.Find(ctx, filter, options)
	cursor.All(ctx, &items)

	if err != nil {
		return items, err
	}
	return items, nil
}

// Update ...
func (g *General) Update() error {
	g.UpdatedAt = time.Now()

	collection := db.GetCollection(GeneralCollectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	item := new(General)
	collection.FindOne(ctx, bson.M{"locale": g.Locale, "_id": bson.M{"$ne": g.ID}}).Decode(&item)

	if item.Locale != "" {
		return errors.New("general is already exist in this locale")
	}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": g.ID}, bson.M{"$set": g})
	collection.FindOne(ctx, bson.M{"_id": g.ID}).Decode(&g)

	if err != nil {
		return err
	}

	return nil
}

// Delete ...
func (g *General) Delete() error {
	collection := db.GetCollection(GeneralCollectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection.FindOne(ctx, bson.M{"_id": g.ID}).Decode(&g)
	if g.Locale == "" {
		return errors.New("general doesn't exist")
	}

	_, err := collection.DeleteOne(ctx, bson.M{"_id": g.ID})
	if err != nil {
		return err
	}

	return nil
}
