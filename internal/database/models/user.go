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
  DisplayName string `sql:",unique,notnull"`
  FirstName string
  LastName string
  Email string `sql:",unique,notnull"`
  Username string `sql:",unique,notnull"`
  Member bool `sql:",notnull"`
  Avatar []byte
  NewsletterNotification bool

  // ResidenceAddressId uuid.UUID `sql:",notnull"`

  FavoriteTracks []uuid.UUID `sql:",type:uuid[]" pg:",array"`
  FollowedGroups []uuid.UUID `sql:",type:uuid[]" pg:",array"`
  // Playlists

  // MemberOfGroups
  // Tags

  // Shares => Membership API?
}
