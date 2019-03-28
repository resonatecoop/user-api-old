package context

import (
	"context"

	"user-api/internal/model"
)

// KeyString should be used when setting and fetching context values
type KeyString string

// JWTKey is a context key for storing token
var JWTKey = "http_jwt_key"

// Service represents context service
type Service struct{}

// GetUser fetches auth user from context
func (s *Service) GetUser(c context.Context) *model.AuthUser {
	return c.Value(KeyString("_authuser")).(*model.AuthUser)
}
