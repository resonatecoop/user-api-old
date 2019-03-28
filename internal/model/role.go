package model

import "context"
import "github.com/satori/go.uuid"

// AccessRole represents access role type
type AccessRole int8

const (
	// SuperAdminRole has all permissions
	SuperAdminRole AccessRole = iota + 1 // 1

	// AdminRole has admin specific permissions
	AdminRole // 2

	// TenantAdminRole can edit tenant specific things
	TenantAdminRole // 3

	// UserRole is a standard user
	UserRole // 4
)

// RBACService represents role-based access control service interface
type RBACService interface {
	EnforceRole(context.Context, AccessRole) bool
	EnforceUser(context.Context, uuid.UUID) bool
	EnforceTenant(context.Context, uuid.UUID) bool
	EnforceTenantAdmin(context.Context, int32) bool
	EnforceTenantAndRole(context.Context, AccessRole, int32) bool
	IsLowerRole(context.Context, AccessRole) bool
}

// Role entity
type Role struct {
	Id int `json:"id"`
	Name string `json:"name"`
}
