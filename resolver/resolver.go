package resolver

import "go-graphql-mongodb-boilerplate/generated/gqlgen"

// Resolver ...
type Resolver struct {
	Token *string
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
