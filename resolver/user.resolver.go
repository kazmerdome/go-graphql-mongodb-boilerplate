package resolver

import (
	"context"
	"fmt"

	"go-graphql-mongodb-boilerplate/guard"
	"go-graphql-mongodb-boilerplate/model"
	"go-graphql-mongodb-boilerplate/utility"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

/**
 * Generated Structs
 */

/**
 * Generated Funcs
 */

/**
 * Field Methods
 */

/**
 * Queries
 */
func (r *queryResolver) UserMe(ctx context.Context) (*model.User, error) {
	userID, err := guard.GetUserID(*r.Resolver.Token)
	if err != nil {
		return nil, err
	}
	where := model.UserWhereDTO{ID: &userID}
	item := model.User{}
	item.One(&where)
	if item.Email == "" {
		return nil, fmt.Errorf("Access denied")
	}
	return &item, nil
}

func (r *queryResolver) User(ctx context.Context, where *model.UserWhereDTO, search *string) (*model.User, error) {
	if where == nil {
		where = &model.UserWhereDTO{}
	}
	if search != nil {
		where.OR = utility.MongoSearchFieldParser(model.SEARCH_FILEDS__USER, *search)
	}
	item := model.User{}
	item.One(where)
	if item.Email == "" {
		return nil, nil
	}
	return &item, nil
}

func (r *queryResolver) Users(ctx context.Context, where *model.UserWhereDTO, search *string, in []*primitive.ObjectID, orderBy *model.UserOrderByENUM, skip *int, limit *int) ([]*model.User, error) {
	if where == nil {
		where = &model.UserWhereDTO{}
	}
	if search != nil {
		where.OR = utility.MongoSearchFieldParser(model.SEARCH_FILEDS__USER, *search)
	}

	// "in" operation for cherrypicking by ids
	var customQuery *primitive.M
	if in != nil {
		q := bson.M{"_id": bson.M{"$in": in}}
		customQuery = &q
	}

	item := model.User{}
	items, err := item.List(where, orderBy, skip, limit, customQuery)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *queryResolver) UserCount(ctx context.Context, where *model.UserWhereDTO, search *string) (*int, error) {
	u := model.User{}
	if where == nil {
		where = &model.UserWhereDTO{}
	}
	count, err := u.Count(where)
	if err != nil {
		return nil, err
	}
	return &count, nil
}

/**
 * Mutations
 */
func (r *mutationResolver) CreateUser(ctx context.Context, data model.UserCreateDTO) (*model.User, error) {
	item := model.User{}
	data.Role = guard.User // locked
	// Salt and hash the password using the bcrypt algorithm
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), 8)
	if err != nil {
		return nil, err
	}
	data.Password = string(hashedPassword)
	// Validate
	if err := utility.ValidateStruct(data); err != nil {
		return nil, err
	}
	// model operation
	if err := item.Create(&data); err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *mutationResolver) UpdateUser(ctx context.Context, where model.UserWhereUniqueDTO, data model.UserUpdateDTO) (*model.User, error) {
	item := model.User{}
	whereID := where.ID
	// Validate
	if err := utility.ValidateStruct(data); err != nil {
		return nil, err
	}
	if r.Resolver.Token != nil {
		role := guard.GetRole(*r.Resolver.Token)
		// admin level operations
		if role == guard.Admin {
			item.ID = where.ID
		}
		// contractor level operations
		if role != guard.Admin {
			userID, err := guard.GetUserID(*r.Resolver.Token)
			if err != nil {
				return nil, err
			}
			whereID = userID
		}
	}
	// Hash new password if provided
	if data.Password != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*data.Password), 8)
		if err != nil {
			return nil, err
		}
		pw := string(hashedPassword)
		data.Password = &pw
	}

	if err := item.Update(whereID, &data); err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *mutationResolver) DeleteUser(ctx context.Context, where model.UserWhereUniqueDTO) (*model.User, error) {
	item := model.User{ID: where.ID}
	if err := item.Delete(); err != nil {
		return nil, err
	}
	return &item, nil
}
