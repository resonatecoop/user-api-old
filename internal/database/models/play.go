package models

import (
  "time"

  "github.com/go-pg/pg"
  "github.com/satori/go.uuid"
)

type Play struct {
  Id uuid.UUID `sql:"type:uuid,default:uuid_generate_v4()"`
  CreatedAt time.Time `sql:"default:now()"`
  UserId uuid.UUID `sql:",type:uuid,notnull"`
  TrackId uuid.UUID `sql:",type:uuid,notnull"`
  Type string `sql:"type:play_type,notnull"`
  Credits float32 `sql:",notnull"`
}

// Count number of times a track has been played (and paid) by a user
func CountPlays(trackId uuid.UUID, userId uuid.UUID, db *pg.DB) (int32, error) {
  count, err := db.Model((*Play)(nil)).
    Where("user_id = ?", userId).
    Where("track_id = ?", trackId).
    Where("type = 'paid'").
    Count()
  if err != nil {
    return 0, err
  }
  return int32(count), nil
}
