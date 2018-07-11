package models

import "github.com/satori/go.uuid"

type Tag struct {
  Id uuid.UUID `sql:"type:uuid,default:uuid_generate_v4()"`
  Type string `sql:",notnull"`
  Name string `sql:",notnull"`
}
