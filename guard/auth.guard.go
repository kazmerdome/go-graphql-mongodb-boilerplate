package guard

import (
	"errors"
)

// Role type
type Role string

// User role
const User Role = "USER"

// Partner role
const Partner Role = "PARTNER"

// Editor role
const Editor Role = "EDITOR"

// Admin role
const Admin Role = "ADMIN"

// GetRole ...
func GetRole(token string) Role {
	// TODO
	if token == "ADMIN" {
		return Admin
	}
	if token == "EDITOR" {
		return Editor
	}

	return User
}

func checkRole(usersRole Role, whiteListedRoles []Role) bool {
	conatins := false
	for _, role := range whiteListedRoles {
		if role == usersRole {
			conatins = true
		}
	}
	return conatins
}

// Auth ...
func Auth(roles []Role, token string) error {
	// TODO

	// check the !!token

	// decrypt token

	// fetch user from db

	// check user's roles

	// authenticate w provided infos

	// TEMPORARY SHIT
	if token == "ADMIN" && checkRole(Admin, roles) {
		return nil
	}
	if token == "EDITOR" && checkRole(Editor, roles) {
		return nil
	}
	if token != "ADMIN" && token != "EDITOR" && checkRole(User, roles) {
		return nil
	}

	return errors.New("unauthorized")
}
