package resolver

import (
	"aery-graphql/generated/gqlgen"
	"aery-graphql/guard"
	"aery-graphql/model"
	"context"
)

/**
 * Generated Structs
 */
type generalResolver struct{ *Resolver }
type generalNavigationResolver struct{ *Resolver }

/**
 * Generated Funcs
 */

// General ...
func (r *Resolver) General() gqlgen.GeneralResolver {
	return &generalResolver{r}
}

// GeneralNavigation ...
func (r *Resolver) GeneralNavigation() gqlgen.GeneralNavigationResolver {
	return &generalNavigationResolver{r}
}

/**
 * Field Methods
 */
func (r *generalResolver) Seo(ctx context.Context, obj *model.General) (*model.Seo, error) {
	return &obj.Seo, nil
}
func (r *generalResolver) Info(ctx context.Context, obj *model.General) ([]*model.GeneralInfo, error) {
	var length int
	if len(obj.Info) > 0 {
		length = len(obj.Info)
	}
	items := make([]*model.GeneralInfo, length)
	for index, item := range obj.Info {
		actualItem := item
		items[index] = &actualItem
	}
	return items, nil
}
func (r *generalResolver) Social(ctx context.Context, obj *model.General) ([]*model.GeneralSocial, error) {
	var length int
	if len(obj.Social) > 0 {
		length = len(obj.Social)
	}
	items := make([]*model.GeneralSocial, length)
	for index, item := range obj.Social {
		actualItem := item
		items[index] = &actualItem
	}
	return items, nil
}
func (r *generalResolver) Navigation(ctx context.Context, obj *model.General) ([]*model.GeneralNavigation, error) {
	var length int
	if len(obj.Navigation) > 0 {
		length = len(obj.Navigation)
	}
	items := make([]*model.GeneralNavigation, length)
	for index, item := range obj.Navigation {
		actualItem := item
		items[index] = &actualItem
	}
	return items, nil
}
func (r *generalNavigationResolver) Menu(ctx context.Context, obj *model.GeneralNavigation) (*model.Menu, error) {
	menu := model.Menu{}
	filter := model.MenuWhereInput{ID: &obj.Menu}
	menu.One(&filter)
	return &menu, nil
}

/**
 * Queries
 */
func (r *queryResolver) General(ctx context.Context, where *model.GeneralWhereInput) (*model.General, error) {
	if where == nil {
		where = &model.GeneralWhereInput{}
	}
	item := model.General{}
	item.One(where)
	if item.Locale == "" {
		return nil, nil
	}
	return &item, nil
}
func (r *queryResolver) Generals(
	ctx context.Context,
	where *model.GeneralWhereInput,
	orderBy *model.GeneralOrderByInput,
	skip *int,
	limit *int,
) ([]*model.General, error) {
	if where == nil {
		where = &model.GeneralWhereInput{}
	}
	item := model.General{}
	items, err := item.List(where, orderBy, skip, limit)
	if err != nil {
		return nil, err
	}
	return items, nil
}

/**
 * Mutations
 */
func (r *mutationResolver) CreateGeneral(ctx context.Context, data gqlgen.GeneralCreateInput) (*model.General, error) {
	if err := guard.Auth([]guard.Role{guard.Editor, guard.Admin}, *r.Resolver.Token); err != nil {
		return nil, err
	}

	item := model.General{
		Locale: data.Locale,
	}
	if data.Homepage != nil {
		item.Homepage = *data.Homepage
	}
	if data.Seo != nil {
		item.Seo = model.Seo{
			Title: data.Seo.Create.Title,
		}
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
	if data.Info != nil && len(data.Info.Create) > 0 {
		item.Info = make([]model.GeneralInfo, len(data.Info.Create))
		for index, info := range data.Info.Create {
			item.Info[index] = model.GeneralInfo{
				Type:  info.Type,
				Value: info.Value,
			}
		}
	}
	if data.Social != nil && len(data.Social.Create) > 0 {
		item.Social = make([]model.GeneralSocial, len(data.Social.Create))
		for index, social := range data.Social.Create {
			item.Social[index] = model.GeneralSocial{
				Type:  social.Type,
				Value: social.Value,
			}
		}
	}
	if data.Navigation != nil && len(data.Navigation.Create) > 0 {
		item.Navigation = make([]model.GeneralNavigation, len(data.Navigation.Create))
		for index, navigation := range data.Navigation.Create {
			item.Navigation[index] = model.GeneralNavigation{
				Type: navigation.Type,
				Menu: navigation.Menu,
			}
		}
	}

	if err := item.Create(); err != nil {
		return nil, err
	}

	return &item, nil
}
