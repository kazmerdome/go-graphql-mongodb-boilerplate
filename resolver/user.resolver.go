package resolver

import (
	"aery-graphql/generated/gqlgen"
	"aery-graphql/model"
	"aery-graphql/utility"
	"context"
)

/**
 * Generated Structs
 */
type userResolver struct{ *Resolver }

/**
 * Generated Funcs
 */

// User ...
func (r *Resolver) User() gqlgen.UserResolver {
	return &userResolver{r}
}

/**
 * Field Methods
 */
func (r *userResolver) AppPolicy(ctx context.Context, obj *model.User) ([]*model.AppPolicy, error) {
	var length int
	if len(obj.AppPolicy) > 0 {
		length = len(obj.AppPolicy)
	}

	items := make([]*model.AppPolicy, length)
	for index, item := range obj.AppPolicy {
		actualItem := item
		items[index] = &actualItem
	}

	return items, nil
}

/**
 * Queries
 */
func (r *queryResolver) User(ctx context.Context, where *model.UserWhereInput, searchInFirstnameAndLastnameAndEmail *string) (*model.User, error) {
	if where == nil {
		where = &model.UserWhereInput{}
	}
	if searchInFirstnameAndLastnameAndEmail != nil {
		where.OR = utility.MongoSearchFieldParser(model.UserSearchFields, *searchInFirstnameAndLastnameAndEmail)
	}

	item := model.User{}
	item.One(where)
	if item.Email == "" {
		return nil, nil
	}

	return &item, nil
}

func (r *queryResolver) Users(
	ctx context.Context,
	where *model.UserWhereInput,
	searchInFirstnameAndLastnameAndEmail *string,
	orderBy *model.UserOrderByInput,
	skip *int,
	limit *int,
) ([]*model.User, error) {
	if where == nil {
		where = &model.UserWhereInput{}
	}
	if searchInFirstnameAndLastnameAndEmail != nil {
		where.OR = utility.MongoSearchFieldParser(model.UserSearchFields, *searchInFirstnameAndLastnameAndEmail)
	}

	item := model.User{}
	items, err := item.List(where, orderBy, skip, limit)

	if err != nil {
		return nil, err
	}
	return items, nil
}

/**
 * Mutations
 */
func (r *mutationResolver) CreateUser(ctx context.Context, data gqlgen.UserCreateInput) (*model.User, error) {
	item := model.User{
		Firstname: data.Firstname,
		Lastname:  data.Lastname,
		Email:     data.Email,
		AppPolicy: []model.AppPolicy{},
	}
	if data.GoogleID != nil {
		item.GoogleID = *data.GoogleID
	}
	if data.AppPolicy != nil && len(data.AppPolicy.Create) > 0 {
		item.AppPolicy = make([]model.AppPolicy, len(data.AppPolicy.Create))
		for index, policy := range data.AppPolicy.Create {
			item.AppPolicy[index] = model.AppPolicy{
				Resource: policy.Resource,
				Role:     policy.Role,
			}
		}
	}

	if err := item.Create(); err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *mutationResolver) UpdateUser(ctx context.Context, where model.UserWhereUniqueInput, data gqlgen.UserUpdateInput) (*model.User, error) {
	item := model.User{}
	item.ID = where.ID
	if data.Firstname != nil {
		item.Firstname = *data.Firstname
	}
	if data.Lastname != nil {
		item.Firstname = *data.Lastname
	}
	if data.GoogleID != nil {
		item.Firstname = *data.GoogleID
	}
	if data.AppPolicy != nil && len(data.AppPolicy.Create) > 0 {
		item.AppPolicy = make([]model.AppPolicy, len(data.AppPolicy.Create))
		for index, policy := range data.AppPolicy.Create {
			item.AppPolicy[index] = model.AppPolicy{
				Resource: policy.Resource,
				Role:     policy.Role,
			}
		}
	}

	if err := item.Update(); err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *mutationResolver) DeleteUser(ctx context.Context, where model.UserWhereUniqueInput) (*model.User, error) {
	item := model.User{ID: where.ID}

	if err := item.Delete(); err != nil {
		return nil, err
	}

	return &item, nil
}
