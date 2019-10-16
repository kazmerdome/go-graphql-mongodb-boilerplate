package server

import (
	"aery-graphql/generated/gqlgen"
	"aery-graphql/guard"
	"aery-graphql/model"
	"aery-graphql/resolver"
	"aery-graphql/service"
	"aery-graphql/utility"
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/handler"
	"github.com/labstack/echo"
)

var token string
var locale model.Locale

// Get header value and add to gql resolvers
func getHeaders(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		headers := c.Request().Header
		token = utility.GetHeaderString("Aerylabs-Auth", headers)
		return next(c)
	}
}

// GetRoutes ...
func GetRoutes(e *echo.Echo) {
	auth := new(service.AuthController)
	resolver := resolver.Resolver{Token: &token}
	config := gqlgen.Config{Resolvers: &resolver}
	config.Directives.Auth = func(ctx context.Context, obj interface{}, next graphql.Resolver, role []guard.Role) (interface{}, error) {
		if err := guard.Auth(role, *resolver.Token); err != nil {
			return nil, fmt.Errorf("Access denied")
		}
		return next(ctx)
	}

	e.Use(getHeaders)

	e.GET("/", echo.WrapHandler(handler.Playground("Aery Labs GraphQL Playground", "/query")))
	e.POST(
		"/query",
		echo.WrapHandler(
			handler.GraphQL(
				gqlgen.NewExecutableSchema(config),
				// handler.IntrospectionEnabled(false),
			),
		),
	)
	e.GET("/auth/google/login", auth.OauthGoogle)
	e.GET("/auth/google/callback", auth.OauthGoogleCallback)
}
