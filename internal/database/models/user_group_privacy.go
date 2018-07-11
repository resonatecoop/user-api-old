package models

import "github.com/satori/go.uuid"

type UserGroupPrivacy struct {
  Id uuid.UUID `sql:"type:uuid,default:uuid_generate_v4()"`
  Private bool
  OwnedTracks bool
  SupportedArtists bool
}
