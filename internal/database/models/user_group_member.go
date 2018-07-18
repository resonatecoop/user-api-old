package models

import (
  "time"
  // "fmt"
  "github.com/satori/go.uuid"
)

type UserGroupMember struct {
  Id uuid.UUID `sql:"type:uuid,default:uuid_generate_v4()"`
  CreatedAt time.Time `sql:"default:now()"`
  UpdatedAt time.Time
  UserGroupId uuid.UUID `sql:",type:uuid,notnull"`
  MemberId uuid.UUID `sql:",type:uuid,notnull"`
  DisplayName string
  Tags []uuid.UUID `sql:",type:uuid[]" pg:",array"`
}
