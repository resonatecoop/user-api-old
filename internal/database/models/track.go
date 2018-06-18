package models

import (
  "time"

  "github.com/satori/go.uuid"
)

type Track struct {
  Id uuid.UUID `sql:"type:uuid,default:uuid_generate_v4()"`
  PublishDate time.Time `sql:",notnull"`
  CreatedAt time.Time `sql:"default:now()"`
	UpdatedAt time.Time
  Title string `sql:",notnull"`
  Duration time.Duration `sql:",notnull"`
  Status string  `sql:"type:status,notnull"`
  Cover []byte `sql:",notnull"`
  // TotalTrackPlayCount int

  FavoriteOfUsers []uuid.UUID `sql:",type:uuid[]" pg:",array"`

  // Performers map[string]string `pg:",hstore"`
  // Labels map[string]string `pg:",hstore"`
  // Tags map[string]string `pg:",hstore";sql:",notnull"`
  // TrackGroups map[string]string `pg:",hstore"`
  // Links map[string]string `pg:",hstore"`
  //
  // Rights map[string]string `pg:",hstore;sql:",notnull""`

  // Audiofiles []uuid.UUID `sql:",type:uuid[]" pg:",array"`
}
