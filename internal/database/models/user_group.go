package models

import (
  "time"
  // "fmt"
  "github.com/go-pg/pg"
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

  Members []UserGroup `pg:"many2many:user_group_members,fk:user_group_id,joinFK:member_id"`
  MemberOfGroups []UserGroup `pg:"many2many:user_group_members,fk:member_id,joinFK:user_group_id"`

  Tracks []uuid.UUID `sql:",type:uuid[]" pg:",array"`
  TrackGroups []uuid.UUID `sql:",type:uuid[]" pg:",array"`

  // TODO remove
  SubGroups []uuid.UUID `sql:",type:uuid[]" pg:",array"`
  // artist
  Labels []uuid.UUID `sql:",type:uuid[]" pg:",array"`
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

// TODO delete track and track group
func (u *UserGroup) Delete(db *pg.DB) (error, string) {
  var table string
  tx, err := db.Begin()
  if err != nil {
    return err, table
  }
  defer tx.Rollback()

  pgerr := tx.Model(u).
    Column("user_group.links","user_group.followers", "Address", "Privacy").
    WherePK().
    Select()
  if pgerr != nil {
    return pgerr, "user_group"
  }

  if len(u.Links) > 0 {
    _, pgerr = tx.Model((*Link)(nil)).
      Where("id in (?)", pg.In(u.Links)).
      Delete()
    if pgerr != nil {
      return pgerr, "link"
    }
  }

  if len(u.Followers) > 0 {
    _, pgerr = tx.ExecOne(`
      UPDATE users
      SET followed_groups = array_remove(followed_groups, ?)
      WHERE id IN (?)
    `, u.Id, pg.In(u.Followers))
    if pgerr != nil {
      return pgerr, "user"
    }
  }

  _, pgerr = tx.Model(u.Address).WherePK().Delete()
  if pgerr != nil {
    return pgerr, "street_address"
  }

  _, pgerr = tx.Model(u.Privacy).WherePK().Delete()
  if pgerr != nil {
    return pgerr, "user_group_privacy"
  }

  var userGroupMembers []UserGroupMember
  _, pgerr = tx.Model(&userGroupMembers).
    Where("user_group_id = ?", u.Id).
    WhereOr("member_id = ?", u.Id).
    Delete()
  if pgerr != nil {
    return pgerr, "user_group_member"
  }

  pgerr = tx.Delete(u)
  if pgerr != nil {
    return pgerr, "user_group"
  }

  return tx.Commit(), table
}
