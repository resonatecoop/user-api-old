package models

import (
  "time"

  "github.com/satori/go.uuid"
)

type User struct {
  Id uuid.UUID `sql:"type:uuid,default:uuid_generate_v4()"`
  CreatedAt time.Time `sql:"default:now()"`
  UpdatedAt time.Time
  Username string `sql:",notnull,unique"`
  FullName string `sql:",notnull"`
  FirstName string
  LastName string
  Email string `sql:",unique,notnull"`
  Member bool `sql:",notnull"`
  NewsletterNotification bool

  FavoriteTracks []uuid.UUID `sql:",type:uuid[]" pg:",array"`
  FollowedGroups []uuid.UUID `sql:",type:uuid[]" pg:",array"`
  Playlists []uuid.UUID `sql:",type:uuid[]" pg:",array"`

  OwnerOfGroups []UserGroup `pg:"fk:owner_id"`
}
