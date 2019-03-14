package model

import (
  "time"
  // "fmt"
  "github.com/satori/go.uuid"
)

type UserGroupMember struct {
  CreatedAt time.Time `sql:"default:now()"`
  UpdatedAt time.Time
  UserGroupId uuid.UUID `sql:",pk,type:uuid,notnull"`
  MemberId uuid.UUID `sql:",pk,type:uuid,notnull"`
  DisplayName string
  Tags []uuid.UUID `sql:",type:uuid[]" pg:",array"`
}
