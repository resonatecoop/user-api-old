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

  // FeaturedTrack *Track

  Kvstore map[string]string `pg:",hstore"`
  Followers []uuid.UUID `sql:",type:uuid[]" pg:",array"`

  // AdminUsers map[string]string `pg:",hstore;sql:",notnull""`
  // Members map[string]string `pg:",hstore"`
  // SubGroups map[string]string `pg:",hstore"`

  // Links map[string]string `pg:",hstore"`
  // Tags map[string]string `pg:",hstore"`
  // Tracks map[string]string `pg:",hstore"`
  // TrackGroups map[string]string `pg:",hstore"`

  // artist
  // Labels map[string]string `pg:",hstore"`
  // Payees []*User

  // distributor
  // Distributees []*UserGroup
}
