package model

import "github.com/satori/go.uuid"

// Tenant table
type Tenant struct {
	Id uuid.UUID `sql:"type:uuid,default:uuid_generate_v4()"`
	Name   string `json:"name"`
	Active bool   `json:"active"`
}
