package main

import (
	"context"
	"net/url"

	"github.com/lucasvillarinho/restql"
)

// UserService handles the business logic for users
type UserService struct{}

// NewUserService creates a new instance of UserService
func NewUserService() *UserService {
	return &UserService{}
}

// GetUsers returns filtered users using RestQL + GORM
func (s *UserService) GetUsers(ctx context.Context, queryParams url.Values) ([]User, string, any, error) {
	qb, err := restql.Parse(queryParams, UsersSchema)
	if err != nil {
		return nil, "", nil, err
	}
	sql, args := qb.ToSQL()

	var users []User
	result := db.WithContext(ctx).Raw(sql, args...).Scan(&users)
	if result.Error != nil {
		return nil, sql, args, result.Error
	}

	return users, sql, args, nil
}
