//go:generate kallax gen
package models

import "gopkg.in/src-d/go-kallax.v1"

type User struct {
        kallax.Model         `table:"users"`
        ID       kallax.UUID  `pk:""`
        Username string
        Email    string
        Address string
}


func newUser(id kallax.UUID, username string, email string, address string) (*User, error) {
  return &User{ID: id, Username: username, Email: email, Address: address}, nil
}
