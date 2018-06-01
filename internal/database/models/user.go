package models

import (
  "time"

  "github.com/satori/go.uuid"
)


type User struct {
  Id uuid.UUID `sql:"type:uuid"`
  FullName string `sql:",notnull"`
  DisplayName string `sql:",unique,notnull"`
  FirstName string
  LastName string
  Email string `sql:",unique,notnull"`
  Username string `sql:",unique,notnull"`
  Member bool `sql:",notnull"`
  Avatar []byte
  NewsletterNotification bool
  CreatedAt time.Time `sql:"default:now()"`
	UpdatedAt time.Time

  // ResidenceAddress
  // Tags
  // MemberOfGroups
  // Shares
  // FavouriteTracks
  // Playlists
  // FollowedArtists
  // PaymentMechanisms
}
