package mock

import (
	"context"

	"user-api/internal/model"
)

// RBAC Mock
type RBAC struct {
	EnforceRoleFn          func(context.Context, model.AccessRole) bool
	EnforceUserFn          func(context.Context, int64) bool
	EnforceTenantFn        func(context.Context, int32) bool
	EnforceTenantAdminFn   func(context.Context, int32) bool
	IsLowerRoleFn          func(context.Context, model.AccessRole) bool
	EnforceTenantAndRoleFn func(context.Context, model.AccessRole, int32) bool
}

// EnforceRole mock
func (a *RBAC) EnforceRole(c context.Context, role model.AccessRole) bool {
	return a.EnforceRoleFn(c, role)
}

// EnforceUser mock
func (a *RBAC) EnforceUser(c context.Context, id int64) bool {
	return a.EnforceUserFn(c, id)
}

// EnforceTenant mock
func (a *RBAC) EnforceTenant(c context.Context, id int32) bool {
	return a.EnforceTenantFn(c, id)
}

// EnforceTenantAdmin mock
func (a *RBAC) EnforceTenantAdmin(c context.Context, id int32) bool {
	return a.EnforceTenantAdminFn(c, id)
}

// EnforceTenantAndRole mock
func (a *RBAC) EnforceTenantAndRole(c context.Context, role model.AccessRole, id int32) bool {
	return a.EnforceTenantAndRoleFn(c, role, id)
}

// IsLowerRole mock
func (a *RBAC) IsLowerRole(c context.Context, role model.AccessRole) bool {
	return a.IsLowerRoleFn(c, role)
}
