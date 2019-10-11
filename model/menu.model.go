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

// MenuCollectionName ...
const MenuCollectionName = "Menu"

// MenuSearchFields ...
var MenuSearchFields = []string{"name", "slug"}

// Menu ...
type Menu struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Locale    Locale             `json:"locale" bson:"locale,omitempty"`
	Status    Status             `json:"status" bson:"status,omitempty"`
	Name      string             `json:"name" bson:"name,omitempty"`
	Slug      string             `json:"slug" bson:"slug,omitempty"`
	Links     []MenuLink         `json:"links" bson:"links,omitempty"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

// MenuLink ...
type MenuLink struct {
	ID   primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Type MenuLinkType       `json:"type" bson:"type,omitempty"`
	Name string             `json:"name" bson:"name,omitempty"`
	URL  string             `json:"url" bson:"url,omitempty"`
}

// MenuWhereUniqueInput ...
type MenuWhereUniqueInput struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
}

// MenuWhereInput ...
type MenuWhereInput struct {
	ID     *primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Slug   *string             `json:"slug" bson:"slug,omitempty"`
	Locale *Locale             `json:"locale" bson:"locale,omitempty"`
	Status *Status             `json:"status" bson:"status,omitempty"`
	OR     []bson.M            `json:"$or,omitempty" bson:"$or,omitempty"`
}

// MenuLinkWhereUniqueInput ...
type MenuLinkWhereUniqueInput struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
}

// MenuLinkType ...
type MenuLinkType string

// MenuOrderByInput ...
type MenuOrderByInput string

// Create ...
func (m *Menu) Create() error {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()

	collection := db.GetCollection(MenuCollectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	item := new(Menu)
	collection.FindOne(ctx, bson.M{"locale": m.Locale, "slug": m.Slug}).Decode(&item)

	if item.Name != "" {
		return errors.New("menu slug is already exist in this locale")
	}

	res, err := collection.InsertOne(ctx, m)

	if err != nil {
		return err
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		m.ID = oid
	}

	return nil
}

// One ...
func (m *Menu) One(filter *MenuWhereInput) error {
	collection := db.GetCollection(MenuCollectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection.FindOne(ctx, &filter).Decode(&m)
	return nil
}

// List ...
func (m *Menu) List(filter *MenuWhereInput, orderBy *MenuOrderByInput, skip *int, limit *int) ([]*Menu, error) {
	var items []*Menu
	orderByKey := "created_at"
	orderByValue := -1
	collection := db.GetCollection(MenuCollectionName)
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
func (m *Menu) Update() error {
	m.UpdatedAt = time.Now()

	collection := db.GetCollection(MenuCollectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	item := new(Menu)
	collection.FindOne(ctx, bson.M{"locale": m.Locale, "slug": m.Slug, "_id": bson.M{"$ne": m.ID}}).Decode(&item)

	if item.Slug != "" {
		return errors.New("menu slug is already exist in this locale")
	}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": m.ID}, bson.M{"$set": m})
	collection.FindOne(ctx, bson.M{"_id": m.ID}).Decode(&m)

	if err != nil {
		return err
	}

	return nil
}

// Delete ...
func (m *Menu) Delete() error {
	collection := db.GetCollection(MenuCollectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection.FindOne(ctx, bson.M{"_id": m.ID}).Decode(&m)
	if m.Slug == "" {
		return errors.New("menu doesn't exist")
	}

	_, err := collection.DeleteOne(ctx, bson.M{"_id": m.ID})
	if err != nil {
		return err
	}

	return nil
}
