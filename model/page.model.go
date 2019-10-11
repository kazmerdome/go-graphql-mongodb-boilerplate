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

// PageCollectionName ...
const PageCollectionName = "Page"

// PageSearchFields ...
var PageSearchFields = []string{"name", "name_slug", "title", "slug"}

// Page ...
type Page struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Locale    Locale             `json:"locale" bson:"locale,omitempty"`
	Status    Status             `json:"status" bson:"status,omitempty"`
	Name      string             `json:"name" bson:"name,omitempty"`
	NameSlug  string             `json:"name_slug" bson:"name_slug,omitempty"`
	Title     string             `json:"title" bson:"title,omitempty"`
	Slug      string             `json:"slug" bson:"slug,omitempty"`
	Seo       Seo                `json:"seo" bson:"seo,omitempty"`
	Modules   []Modules          `json:"modules" bson:"modules,omitempty"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

// Modules ...
type Modules struct {
	ID   primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Type string             `json:"type" bson:"type,omitempty"`
	Info interface{}        `json:"info" bson:"info,omitempty"`
	Data interface{}        `json:"data" bson:"data,omitempty"`
}

// PageWhereUniqueInput ...
type PageWhereUniqueInput struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
}

// PageWhereInput ...
type PageWhereInput struct {
	ID       *primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Locale   *Locale             `json:"locale" bson:"locale,omitempty"`
	Status   *Status             `json:"status" bson:"status,omitempty"`
	Name     *string             `json:"name" bson:"name,omitempty"`
	NameSlug *string             `json:"name_slug" bson:"name_slug,omitempty"`
	Title    *string             `json:"title" bson:"title,omitempty"`
	Slug     *string             `json:"slug" bson:"slug,omitempty"`
	OR       []bson.M            `json:"$or,omitempty" bson:"$or,omitempty"`
}

// PageOrderByInput ...
type PageOrderByInput string

// Create ...
func (p *Page) Create() error {
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()

	collection := db.GetCollection(PageCollectionName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	item := new(Page)
	collection.FindOne(ctx, bson.M{"locale": p.Locale, "name_slug": p.NameSlug}).Decode(&item)
	if item.Slug != "" {
		return errors.New("Page name_slug is already exist in this locale")
	}
	collection.FindOne(ctx, bson.M{"locale": p.Locale, "slug": p.Slug}).Decode(&item)
	if item.Slug != "" {
		return errors.New("Page slug is already exist in this locale")
	}

	res, err := collection.InsertOne(ctx, p)

	if err != nil {
		return err
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		p.ID = oid
	}

	return nil
}

// One ...
func (p *Page) One(filter *PageWhereInput) error {
	collection := db.GetCollection(PageCollectionName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	collection.FindOne(ctx, &filter).Decode(&p)
	return nil
}

// List ...
func (p *Page) List(filter *PageWhereInput, orderBy *PageOrderByInput, skip *int, limit *int) ([]*Page, error) {
	var items []*Page
	orderByKey := "created_at"
	orderByValue := -1
	collection := db.GetCollection(PageCollectionName)
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
func (p *Page) Update() error {
	p.UpdatedAt = time.Now()

	collection := db.GetCollection(PageCollectionName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	item := new(Page)
	collection.FindOne(ctx, bson.M{"locale": p.Locale, "name_slug": p.NameSlug, "_id": bson.M{"$ne": p.ID}}).Decode(&item)
	if item.Slug != "" {
		return errors.New("page name_slug is already exist in this locale")
	}
	collection.FindOne(ctx, bson.M{"locale": p.Locale, "slug": p.Slug, "_id": bson.M{"$ne": p.ID}}).Decode(&item)
	if item.Slug != "" {
		return errors.New("page slug is already exist in this locale")
	}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": p.ID}, bson.M{"$set": p})
	collection.FindOne(ctx, bson.M{"_id": p.ID}).Decode(&p)

	if err != nil {
		return err
	}

	return nil
}

// Delete ...
func (p *Page) Delete() error {
	collection := db.GetCollection(PageCollectionName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	collection.FindOne(ctx, bson.M{"_id": p.ID}).Decode(&p)

	if p.Slug == "" {
		return errors.New("item doesn't exist")
	}

	_, err := collection.DeleteOne(ctx, bson.M{"_id": p.ID})
	if err != nil {
		return err
	}

	return nil
}
