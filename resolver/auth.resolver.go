package resolver

import (
	"context"
	"go-graphql-mongodb-boilerplate/model"
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

/**
 * Mutations
 */
func (r *mutationResolver) AuthBasicStrategySignUp(ctx context.Context, data model.AuthBasicStrategySignUpDTO) (*model.AuthToken, error) {
	basicAuth := model.AuthBasicStrategy{
		Email:    data.Email,
		Username: data.Username,
		Password: data.Password,
	}

	token, err := basicAuth.SignUp()
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *mutationResolver) AuthBasicStrategySignIn(ctx context.Context, data model.AuthBasicStrategySignInDTO) (*model.AuthToken, error) {
	basicAuth := model.AuthBasicStrategy{
		Password: data.Password,
	}
	if data.Username != nil {
		basicAuth.Username = *data.Username
	}

	token, err := basicAuth.SignIn()
	if err != nil {
		return nil, err
	}

	return &token, nil
}
