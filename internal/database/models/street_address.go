package models

import "github.com/satori/go.uuid"

type StreetAddress struct {
  Id uuid.UUID  `sql:"type:uuid,default:uuid_generate_v4()"`
  PersonalData bool `sql:",notnull"`
  Data map[string]string `pg:",hstore"`
}
