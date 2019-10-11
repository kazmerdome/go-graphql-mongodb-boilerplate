package resolver

import (
	"aery-graphql/generated/gqlgen"
	"aery-graphql/model"
)

// Resolver ...
type Resolver struct {
	Token  *string
	Locale *model.Locale
}

// Mutation ...
func (r *Resolver) Mutation() gqlgen.MutationResolver {
	return &mutationResolver{r}
}

// Query ...
func (r *Resolver) Query() gqlgen.QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
