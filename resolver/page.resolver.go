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
type pageResolver struct{ *Resolver }

/**
 * Generated Funcs
 */

// Page ...
func (r *Resolver) Page() gqlgen.PageResolver {
	return &pageResolver{r}
}

/**
 * Field Methods
 */
func (r *pageResolver) Seo(ctx context.Context, obj *model.Page) (*model.Seo, error) {
	return &obj.Seo, nil
}

func (r *pageResolver) Modules(ctx context.Context, obj *model.Page) ([]*model.Modules, error) {
	var length int
	if len(obj.Modules) > 0 {
		length = len(obj.Modules)
	}

	items := make([]*model.Modules, length)
	for index, item := range obj.Modules {
		actualItem := item
		items[index] = &actualItem
	}

	return items, nil
}

/**
 * Queries
 */
func (r *queryResolver) Page(ctx context.Context, where *model.PageWhereInput, searchInNameAndNameSlugAndTitleAndSlug *string) (*model.Page, error) {
	if where == nil {
		where = &model.PageWhereInput{}
	}

	role := guard.GetRole(*r.Resolver.Token)
	if role == guard.User {
		where.Status = &model.StatusActive
	}

	if searchInNameAndNameSlugAndTitleAndSlug != nil {
		where.OR = utility.MongoSearchFieldParser(model.PageSearchFields, *searchInNameAndNameSlugAndTitleAndSlug)
	}

	item := model.Page{}
	item.One(where)
	if item.NameSlug == "" {
		return nil, nil
	}

	return &item, nil
}

func (r *queryResolver) Pages(
	ctx context.Context,
	where *model.PageWhereInput,
	searchInNameAndNameSlugAndTitleAndSlug *string,
	orderBy *model.PageOrderByInput,
	skip *int,
	limit *int,
) ([]*model.Page, error) {
	role := guard.GetRole(*r.Resolver.Token)

	if where == nil {
		where = &model.PageWhereInput{}
	}
	if role == guard.User {
		where.Status = &model.StatusActive
	}
	if searchInNameAndNameSlugAndTitleAndSlug != nil {
		where.OR = utility.MongoSearchFieldParser(model.PageSearchFields, *searchInNameAndNameSlugAndTitleAndSlug)
	}

	item := model.Page{}
	items, err := item.List(where, orderBy, skip, limit)

	if err != nil {
		return nil, err
	}
	return items, nil
}

/**
 * Mutations
 */
func (r *mutationResolver) CreatePage(ctx context.Context, data gqlgen.PageCreateInput) (*model.Page, error) {
	item := model.Page{
		Locale:   data.Locale,
		Status:   data.Status,
		Name:     data.Name,
		NameSlug: data.NameSlug,
		Title:    data.Title,
		Slug:     data.Slug,
		Seo:      model.Seo{},
		Modules:  []model.Modules{},
	}

	if data.Seo != nil && data.Seo.Create != nil {
		item.Seo.Title = data.Seo.Create.Title
		if data.Seo.Create.Description != nil {
			item.Seo.Description = *data.Seo.Create.Description
		}
		if data.Seo.Create.Image != nil {
			item.Seo.Image = *data.Seo.Create.Image
		}
		if data.Seo.Create.Keywords != nil {
			item.Seo.Keywords = *data.Seo.Create.Keywords
		}
	}

	if data.Modules != nil && len(data.Modules.Create) > 0 {
		item.Modules = make([]model.Modules, len(data.Modules.Create))
		for index, module := range data.Modules.Create {
			item.Modules[index] = model.Modules{
				ID:   primitive.NewObjectID(),
				Type: module.Type,
				Info: module.Info,
				Data: module.Data,
			}
		}
	}

	if err := item.Create(); err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *mutationResolver) UpdatePage(ctx context.Context, where model.PageWhereUniqueInput, data gqlgen.PageUpdateInput) (*model.Page, error) {
	item := model.Page{}
	item.ID = where.ID
	item.Locale = data.Locale
	item.NameSlug = data.NameSlug
	item.Slug = data.Slug

	if data.Status != nil {
		item.Status = *data.Status
	}
	if data.Name != nil {
		item.Name = *data.Name
	}
	if data.Title != nil {
		item.Title = *data.Title
	}
	if data.Seo != nil && data.Seo.Create != nil {
		item.Seo.Title = data.Seo.Create.Title
		if data.Seo.Create.Description != nil {
			item.Seo.Description = *data.Seo.Create.Description
		}
		if data.Seo.Create.Image != nil {
			item.Seo.Image = *data.Seo.Create.Image
		}
		if data.Seo.Create.Keywords != nil {
			item.Seo.Keywords = *data.Seo.Create.Keywords
		}
	}
	if data.Modules != nil && len(data.Modules.Create) > 0 {
		item.Modules = make([]model.Modules, len(data.Modules.Create))
		for index, module := range data.Modules.Create {
			item.Modules[index] = model.Modules{
				ID:   primitive.NewObjectID(),
				Type: module.Type,
				Info: module.Info,
				Data: module.Data,
			}
		}
	}

	if err := item.Update(); err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *mutationResolver) DeletePage(ctx context.Context, where model.PageWhereUniqueInput) (*model.Page, error) {
	item := model.Page{ID: where.ID}

	if err := item.Delete(); err != nil {
		return nil, err
	}

	return &item, nil
}
