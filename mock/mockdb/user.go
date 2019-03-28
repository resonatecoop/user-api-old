package mockdb

import (
	"github.com/go-pg/pg/orm"
	"user-api/internal/model"
)

// User database mock
type User struct {
	// CreateFn          func(orm.DB, model.User) (*model.User, error)
	// ViewFn            func(orm.DB, int64) (*model.User, error)
	// ListFn            func(orm.DB, string, int, int) ([]model.User, error)
	// DeleteFn          func(orm.DB, *model.User) error
	// UpdateFn          func(orm.DB, *model.User) (*model.User, error)
	FindByAuthFn      func(orm.DB, string) (*model.User, error)
	FindByTokenFn     func(orm.DB, string) (*model.User, error)
	UpdateLastLoginFn func(orm.DB, *model.User) error
}

// Create mock
// func (u *User) Create(db orm.DB, usr model.User) (*model.User, error) {
// 	return u.CreateFn(db, usr)
// }
//
// // View mock
// func (u *User) View(db orm.DB, id int64) (*model.User, error) {
// 	return u.ViewFn(db, id)
// }
//
// // List mock
// func (u *User) List(db orm.DB, q string, limit, page int) ([]model.User, error) {
// 	return u.ListFn(db, q, limit, page)
// }
//
// // Delete mock
// func (u *User) Delete(db orm.DB, usr *model.User) error {
// 	return u.DeleteFn(db, usr)
// }

// Update mock
// func (u *User) Update(db orm.DB, usr *model.User) (*model.User, error) {
// 	return u.UpdateFn(db, usr)
// }

// FindByAuth mock
func (u *User) FindByAuth(db orm.DB, auth string) (*model.User, error) {
	return u.FindByAuthFn(db, auth)
}

// FindByToken mock
func (u *User) FindByToken(db orm.DB, auth string) (*model.User, error) {
	return u.FindByTokenFn(db, auth)
}

// UpdateLastLogin mock
func (u *User) UpdateLastLogin(db orm.DB, usr *model.User) error {
	return u.UpdateLastLoginFn(db, usr)
}
