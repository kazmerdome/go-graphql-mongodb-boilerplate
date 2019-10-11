package server

import (
	"aery-graphql/generated/gqlgen"
	"aery-graphql/model"
	"aery-graphql/resolver"
	"aery-graphql/utility"

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
		locale = model.Locales[utility.GetHeaderString("Locale", headers)]
		return next(c)
	}
}

// GetRoutes ...
func GetRoutes(e *echo.Echo) {
	resolver := resolver.Resolver{
		Token:  &token,
		Locale: &locale,
	}

	e.Use(getHeaders)
	e.GET("/", echo.WrapHandler(handler.Playground("GraphQL Playground", "/query")))
	e.POST("/query", echo.WrapHandler(handler.GraphQL(gqlgen.NewExecutableSchema(gqlgen.Config{Resolvers: &resolver}))))
}
