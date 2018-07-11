package models

import (
  "time"
  // "fmt"
  "github.com/go-pg/pg/orm"
  "github.com/satori/go.uuid"
)

type UserGroup struct {
  Id uuid.UUID `sql:"type:uuid,default:uuid_generate_v4()"`
  CreatedAt time.Time `sql:"default:now()"`
  UpdatedAt time.Time
  DisplayName string `sql:",unique,notnull"`
  Description string
  ShortBio string
  Avatar []byte `sql:",notnull"`
  Banner []byte
  GroupEmailAddress string

  PrivacyId uuid.UUID `sql:"type:uuid,notnull"`
  Privacy *UserGroupPrivacy

  AddressId uuid.UUID  `sql:"type:uuid,notnull"`
  Address *StreetAddress

  TypeId uuid.UUID `sql:"type:uuid,notnull"`
  Type *GroupTaxonomy

  OwnerId uuid.UUID `sql:"type:uuid,notnull"`
  Owner *User

  Links []uuid.UUID `sql:",type:uuid[]" pg:",array"`
  Tags []uuid.UUID `sql:",type:uuid[]" pg:",array"`
  RecommendedArtists []uuid.UUID `sql:",type:uuid[]" pg:",array"`
  HighlightedTracks []uuid.UUID `sql:",type:uuid[]" pg:",array"`
  // FeaturedTrackGroup uuid.UUID  `sql:"type:uuid"`

  Kvstore map[string]string `pg:",hstore"`
  Followers []uuid.UUID `sql:",type:uuid[]" pg:",array"`

  AdminUsers []uuid.UUID `sql:",type:uuid[]" pg:",array"`
  SubGroups []uuid.UUID `sql:",type:uuid[]" pg:",array"`

  // TODO need classic m2m junction table to store additional info
  // associated to member, e.g. this member plays drums in this band
  // Members

  Tracks []uuid.UUID `sql:",type:uuid[]" pg:",array"`
  TrackGroups []uuid.UUID `sql:",type:uuid[]" pg:",array"`

  // artist
  Labels []uuid.UUID `sql:",type:uuid[]" pg:",array"`
  // Payees []*User

  // distributor
  // Distributees []*UserGroup
}

func (u *UserGroup) BeforeInsert(db orm.DB) error {
  newPrivacy := &UserGroupPrivacy{Private: false, OwnedTracks: true, SupportedArtists: true}
  _, pgerr := db.Model(newPrivacy).Returning("*").Insert()
  if pgerr != nil {
    return pgerr
  }
  u.PrivacyId = newPrivacy.Id

  return nil
}
