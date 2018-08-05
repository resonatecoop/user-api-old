package models

import (
  "time"

  "github.com/satori/go.uuid"
)

type Play struct {
  CreatedAt time.Time `sql:"default:now()"`
  UserId uuid.UUID `sql:",pk,type:uuid,notnull"`
  TrackId uuid.UUID `sql:",pk,type:uuid,notnull"`
  Type string `sql:"type:play_type,notnull"`
  Credits float32 `sql:",notnull"`
}
