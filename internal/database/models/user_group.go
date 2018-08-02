package models

import (
  "time"
  // "fmt"
  "github.com/go-pg/pg"
  "github.com/go-pg/pg/orm"
  "github.com/satori/go.uuid"
  // "github.com/twitchtv/twirp"

  // pb "user-api/rpc/usergroup"
  trackpb "user-api/rpc/track"

  "user-api/internal"
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
  RecommendedBy []uuid.UUID `sql:",type:uuid[]" pg:",array"`

  HighlightedTracks []uuid.UUID `sql:",type:uuid[]" pg:",array"`
  FeaturedTrackGroupId uuid.NullUUID `sql:"type:uuid,default:uuid_nil()"`

  Kvstore map[string]string `pg:",hstore"`
  Followers []uuid.UUID `sql:",type:uuid[]" pg:",array"`

  AdminUsers []uuid.UUID `sql:",type:uuid[]" pg:",array"`

  Members []UserGroup `pg:"many2many:user_group_members,fk:user_group_id,joinFK:member_id"`
  MemberOfGroups []UserGroup `pg:"many2many:user_group_members,fk:member_id,joinFK:user_group_id"`

  OwnerOfTracks []Track `pg:"fk:user_group_id"` // user group gets paid for these tracks
  OwnerOfTrackGroups []TrackGroup `pg:"fk:user_group_id"`
  Tracks []uuid.UUID `sql:",type:uuid[]" pg:",array"` // user group owner or displayed as artist for these tracks
  TrackGroups []uuid.UUID `sql:",type:uuid[]" pg:",array"` // user group owner or label for these track groups

  // SubGroups []uuid.UUID `sql:",type:uuid[]" pg:",array"`
  // artist
  // Labels []uuid.UUID `sql:",type:uuid[]" pg:",array"`
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

// TODO delete  track group
// delete from recommended artists
func (u *UserGroup) Delete(tx *pg.Tx) (error, string) {
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
    _, pgerr = tx.Exec(`
      UPDATE users
      SET followed_groups = array_remove(followed_groups, ?)
      WHERE id IN (?)
    `, u.Id, pg.In(u.Followers))
    if pgerr != nil {
      return pgerr, "user"
    }
  }

  var userGroupMembers []UserGroupMember
  _, pgerr = tx.Model(&userGroupMembers).
    Where("user_group_id = ?", u.Id).
    WhereOr("member_id = ?", u.Id).
    Delete()
  if pgerr != nil {
    return pgerr, "user_group_member"
  }

  _, pgerr = tx.Model(u).WherePK().Delete()
  if pgerr != nil {
    return pgerr, "user_group"
  }

  _, pgerr = tx.Model(u.Address).WherePK().Delete()
  if pgerr != nil {
    return pgerr, "street_address"
  }

  _, pgerr = tx.Model(u.Privacy).WherePK().Delete()
  if pgerr != nil {
    return pgerr, "user_group_privacy"
  }

  return nil, ""
}

func (u *UserGroup) AddRecommended(db *pg.DB, recommendedId uuid.UUID) (error, string) {
  var table string
  tx, err := db.Begin()
  if err != nil {
    return err, table
  }
  defer tx.Rollback()

  res, pgerr := tx.Exec(`
    UPDATE user_groups
    SET recommended_artists = (select array_agg(distinct e) from unnest(recommended_artists || ?) e)
    WHERE id = ?
  `, pg.Array([]uuid.UUID{recommendedId}), u.Id)
  if res.RowsAffected() == 0 {
    return pg.ErrNoRows, "user_group"
  }
  if pgerr != nil {
    return pgerr, "user_group"
  }

  res, pgerr = tx.Exec(`
    UPDATE user_groups
    SET recommended_by = (select array_agg(distinct e) from unnest(recommended_by || ?) e)
    WHERE id = ?
  `, pg.Array([]uuid.UUID{u.Id}), recommendedId)
  if res.RowsAffected() == 0 {
    return pg.ErrNoRows, "recommended"
  }
  if pgerr != nil {
    return pgerr, "user_group"
  }

  return tx.Commit(), table
}

func (u *UserGroup) RemoveRecommended(db *pg.DB, recommendedId uuid.UUID) (error, string) {
  var table string
  tx, err := db.Begin()
  if err != nil {
    return err, table
  }
  defer tx.Rollback()

  res, pgerr := tx.Exec(`
    UPDATE user_groups
    SET recommended_artists = array_remove(recommended_artists, ?)
    WHERE id = ?
  `, recommendedId, u.Id)
  if res.RowsAffected() == 0 {
    return pg.ErrNoRows, "user_group"
  }
  if pgerr != nil {
    return pgerr, "user_group"
  }

  res, pgerr = tx.Exec(`
    UPDATE user_groups
    SET recommended_by = array_remove(recommended_by, ?)
    WHERE id = ?
  `, u.Id, recommendedId)
  if res.RowsAffected() == 0 {
    return pg.ErrNoRows, "recommended"
  }
  if pgerr != nil {
    return pgerr, "user_group"
  }

  return tx.Commit(), table
}

// Select user groups in db with given 'ids'
// Return slice of UserGroup response
func GetRelatedUserGroups(ids []uuid.UUID, db *pg.DB) ([]*trackpb.RelatedUserGroup, error) {
	groupsResponse := make([]*trackpb.RelatedUserGroup, len(ids))
	if len(ids) > 0 {
		var groups []UserGroup
		pgerr := db.Model(&groups).
			Where("id in (?)", pg.In(ids)).
			Select()
		if pgerr != nil {
			return nil, pgerr
		}
		for i, group := range groups {
			groupsResponse[i] = &trackpb.RelatedUserGroup{
        Id: group.Id.String(),
        DisplayName: group.DisplayName,
        Avatar: group.Avatar,
      }
		}
	}

	return groupsResponse, nil
}

// Select user groups in db with given ids in 'userGroups'
// Return ids slice
// Used in CreateUserGroup/UpdateUserGroup to add/update ids slice to recommended Artists
func GetRelatedUserGroupIds(userGroups []*trackpb.RelatedUserGroup, db *pg.Tx) ([]uuid.UUID, error) {
	relatedUserGroups := make([]*UserGroup, len(userGroups))
	relatedUserGroupIds := make([]uuid.UUID, len(userGroups))
	for i, userGroup := range userGroups {
		id, twerr := internal.GetUuidFromString(userGroup.Id)
		if twerr != nil {
			return nil, twerr.(error)
		}
		relatedUserGroups[i] = &UserGroup{Id: id}
		pgerr := db.Model(relatedUserGroups[i]).
			WherePK().
			Returning("id", "display_name", "avatar").
			Select()
		if pgerr != nil {
			return nil, pgerr
		}
		userGroup.DisplayName = relatedUserGroups[i].DisplayName
		userGroup.Avatar = relatedUserGroups[i].Avatar
		relatedUserGroupIds[i] = relatedUserGroups[i].Id
	}
	return relatedUserGroupIds, nil
}
