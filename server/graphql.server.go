package server

import (
	"context"
	"fmt"

	c "go-graphql-mongodb-boilerplate/config"
	"go-graphql-mongodb-boilerplate/dataloader"
	"go-graphql-mongodb-boilerplate/generated/gqlgen"
	"go-graphql-mongodb-boilerplate/guard"
	"go-graphql-mongodb-boilerplate/resolver"
	"go-graphql-mongodb-boilerplate/utility"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/labstack/echo"
)

var token string
var playgroundPassword string
var submissionToken string

// Get header value and add to gql resolvers
func getHeaders(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		headers := ctx.Request().Header
		token = utility.GetHeaderString("Authorization", headers)
		playgroundPassword = utility.GetHeaderString("Playground-Password", headers)

		return next(ctx)
	}
}

func InitGraphql(e *echo.Echo) {
	// Resolvers && Directives
	resolver := resolver.Resolver{Token: &token}
	config := gqlgen.Config{Resolvers: &resolver}
	config.Directives.Auth = func(ctx context.Context, obj interface{}, next graphql.Resolver, role []guard.Role) (interface{}, error) {
		if err := guard.Auth(role, *resolver.Token); err != nil {
			return nil, fmt.Errorf("Access denied")
		}
		return next(ctx)
	}
	e.Use(getHeaders)

	// new custom handler based on gqlgen version 0.11.3
	queryHandler := handler.New(gqlgen.NewExecutableSchema(config))
	// queryHandler.Use(&debug.Tracer{})
	queryHandler.AddTransport(transport.POST{})
	queryHandler.AddTransport(transport.MultipartForm{})
	queryHandler.SetQueryCache(lru.New(1000))
	queryHandler.Use(extension.AutomaticPersistedQuery{Cache: lru.New(100)})
	queryHandler.AroundOperations(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
		rc := graphql.GetOperationContext(ctx)
		if playgroundPassword == c.GetSecret("GRAPHQL_PLAYGROUND_PASS") {
			rc.DisableIntrospection = false
		} else {
			rc.DisableIntrospection = true
		}
		return next(ctx)
	})

	e.GET("/", echo.WrapHandler(playground.Handler("GraphQL Playground", c.GetSecret("GRAPHQL_ENDPOINT"))))
	e.POST("/query", echo.WrapHandler(dataloader.DataLoaderMiddleware(queryHandler)))
}
