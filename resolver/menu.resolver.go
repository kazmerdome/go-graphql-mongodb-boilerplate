package resolver

import (
	"aery-graphql/generated/gqlgen"
	"aery-graphql/guard"
	"aery-graphql/model"
	"aery-graphql/utility"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/**
 * Generated Structs
 */
type menuResolver struct{ *Resolver }

/**
 * Generated Funcs
 */

// Menu ...
func (r *Resolver) Menu() gqlgen.MenuResolver {
	return &menuResolver{r}
}

/**
 * Field Methods
 */
func (r *menuResolver) Links(ctx context.Context, obj *model.Menu) ([]*model.MenuLink, error) {
	var length int
	if len(obj.Links) > 0 {
		length = len(obj.Links)
	}
	items := make([]*model.MenuLink, length)
	for index, item := range obj.Links {
		actualItem := item
		items[index] = &actualItem
	}
	return items, nil
}

/**
 * Queries
 */
func (r *queryResolver) Menu(ctx context.Context, where *model.MenuWhereInput, searchInNameAndSlug *string) (*model.Menu, error) {
	if where == nil {
		where = &model.MenuWhereInput{}
	}

	role := guard.GetRole(*r.Resolver.Token)
	if role == guard.User {
		where.Status = &model.StatusActive
	}

	if searchInNameAndSlug != nil {
		where.OR = utility.SearchParser(model.MenuSearchFields, *searchInNameAndSlug)
	}

	item := model.Menu{}
	item.One(where)
	if item.Slug == "" {
		return nil, nil
	}

	return &item, nil
}

func (r *queryResolver) Menus(
	ctx context.Context,
	where *model.MenuWhereInput,
	searchInNameAndSlug *string,
	orderBy *model.MenuOrderByInput,
	skip *int,
	limit *int,
) ([]*model.Menu, error) {
	role := guard.GetRole(*r.Resolver.Token)

	if where == nil {
		where = &model.MenuWhereInput{}
	}
	if role == guard.User {
		where.Status = &model.StatusActive
	}
	if searchInNameAndSlug != nil {
		where.OR = utility.SearchParser(model.MenuSearchFields, *searchInNameAndSlug)
	}

	item := model.Menu{}
	items, err := item.List(where, orderBy, skip, limit)

	if err != nil {
		return nil, err
	}
	return items, nil
}

/**
 * Mutations
 */
func (r *mutationResolver) CreateMenu(ctx context.Context, data *gqlgen.MenuCreateInput) (*model.Menu, error) {
	item := model.Menu{
		Locale: data.Locale,
		Status: data.Status,
		Name:   data.Name,
		Slug:   data.Slug,
	}

	if data.Links != nil && len(data.Links.Create) > 0 {
		item.Links = make([]model.MenuLink, len(data.Links.Create))

		for index, link := range data.Links.Create {
			item.Links[index] = model.MenuLink{
				ID:   primitive.NewObjectID(),
				Type: link.Type,
				Name: link.Name,
				URL:  link.URL,
			}
		}
	}

	if err := item.Create(); err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *mutationResolver) UpdateMenu(ctx context.Context, where model.MenuWhereUniqueInput, data gqlgen.MenuUpdateInput) (*model.Menu, error) {
	item := model.Menu{}
	item.ID = where.ID
	item.Locale = data.Locale
	item.Slug = data.Slug

	if data.Status != nil {
		item.Status = *data.Status
	}
	if data.Name != nil {
		item.Name = *data.Name
	}
	if data.Links != nil && len(data.Links.Create) > 0 {
		item.Links = make([]model.MenuLink, len(data.Links.Create))

		for index, link := range data.Links.Create {
			item.Links[index] = model.MenuLink{
				ID:   primitive.NewObjectID(),
				Type: link.Type,
				Name: link.Name,
				URL:  link.URL,
			}
		}
	}

	if err := item.Update(); err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *mutationResolver) DeleteMenu(ctx context.Context, where model.MenuWhereUniqueInput) (*model.Menu, error) {
	item := model.Menu{ID: where.ID}

	if err := item.Delete(); err != nil {
		return nil, err
	}

	return &item, nil
}
