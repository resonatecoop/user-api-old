package model

import "github.com/satori/go.uuid"

type Link struct {
  Id uuid.UUID `sql:"type:uuid,default:uuid_generate_v4()"`
  Uri string `sql:",notnull"`
  Type string
  Platform string `sql:",notnull"`
  PersonalData bool `sql:",notnull"`
}
