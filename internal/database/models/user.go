package models

import (
  "time"

  "github.com/satori/go.uuid"
)

type User struct {
  Id uuid.UUID `sql:"type:uuid,default:uuid_generate_v4()"`
  CreatedAt time.Time `sql:"default:now()"`
  UpdatedAt time.Time
  FullName string `sql:",notnull"`
  DisplayName string `sql:",unique"`
  FirstName string
  LastName string
  Email string `sql:",unique,notnull"`
  Username string `sql:",unique"`
  Member bool `sql:",notnull"`
  Avatar []byte
  NewsletterNotification bool

  ResidenceAddressId uuid.UUID  `sql:"type:uuid,notnull"`
  ResidenceAddress *StreetAddress

  FavoriteTracks []uuid.UUID `sql:",type:uuid[]" pg:",array"`
  FollowedGroups []uuid.UUID `sql:",type:uuid[]" pg:",array"`
  // Playlists

  OwnerOfGroups []UserGroup `pg:"fk:owner_id"`
  // MemberOfGroups
}
