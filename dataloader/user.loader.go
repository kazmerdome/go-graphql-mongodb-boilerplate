package dataloader

import (
	"go-graphql-mongodb-boilerplate/generated/dataloaden"
	"go-graphql-mongodb-boilerplate/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// step 1:
//go:generate go run github.com/vektah/dataloaden UserLoader string *go-graphql-mongodb-boilerplate/model.User

// step 2:
// copy generated file to generated/dataloaden and rename generated package name to dataloaden

func getUserLoader() *dataloaden.UserLoader {
	maxLimit := 150

	return dataloaden.NewUserLoader(
		dataloaden.UserLoaderConfig{
			MaxBatch: maxLimit,
			Wait:     1 * time.Millisecond,
			Fetch: func(keys []string) ([]*model.User, []error) {
				var filter model.UserWhereDTO

				item := model.User{}

				objectIds := make([]primitive.ObjectID, len(keys))
				for i, k := range keys {
					oid, _ := primitive.ObjectIDFromHex(k)
					objectIds[i] = oid
				}
				customQuery := bson.M{"_id": bson.M{"$in": objectIds}}

				items, err := item.List(&filter, nil, nil, &maxLimit, &customQuery)
				if err != nil {
					return nil, []error{err}
				}

				w := make(map[string]*model.User, len(items))
				for _, item := range items {
					w[item.ID.Hex()] = item
				}

				result := make([]*model.User, len(keys))
				for i, key := range keys {
					result[i] = w[key]
				}

				return result, nil
			},
		},
	)
}
