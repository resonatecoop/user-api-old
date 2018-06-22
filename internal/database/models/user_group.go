package models

import (
  "time"

  "github.com/satori/go.uuid"
)

type UserGroup struct {
  Id uuid.UUID `sql:"type:uuid,default:uuid_generate_v4()"`
  CreatedAt time.Time `sql:"default:now()"`
  UpdatedAt time.Time
  DisplayName string `sql:",unique,notnull"`
  Description string
  Avatar []byte `sql:",notnull"`
  Banner []byte
  GroupEmailAddress string

  AddressId uuid.UUID  `sql:"type:uuid,notnull"`
  Address *StreetAddress

  TypeId uuid.UUID `sql:"type:uuid,notnull"`
  Type *GroupTaxonomy

  OwnerId uuid.UUID `sql:"type:uuid,notnull"`
  Owner *User

  // FeaturedTrack *Track or multiple tracks?

  Kvstore map[string]string `pg:",hstore"`
  Followers []uuid.UUID `sql:",type:uuid[]" pg:",array"`

  AdminUsers []uuid.UUID `sql:",type:uuid[],notnull" pg:",array"`
  SubGroups []uuid.UUID `sql:",type:uuid[]" pg:",array"`

  // TODO need classic m2m junction table to store Tags
  // associated to member, e.g. this member plays drums in this band
  // Members

  Links []uuid.UUID `sql:",type:uuid[]" pg:",array"`
  Tags []uuid.UUID `sql:",type:uuid[]" pg:",array"`
  Tracks []uuid.UUID `sql:",type:uuid[]" pg:",array"`
  TrackGroups []uuid.UUID `sql:",type:uuid[]" pg:",array"`

  // artist
  Labels []uuid.UUID `sql:",type:uuid[]" pg:",array"`
  // Payees []*User

  // distributor
  // Distributees []*UserGroup
}
