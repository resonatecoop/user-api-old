package models

import (
  "time"
  // "log"
  // "fmt"
  "github.com/go-pg/pg"
  "github.com/go-pg/pg/orm"
  "github.com/satori/go.uuid"
  "github.com/twitchtv/twirp"

  pb "user-api/rpc/usergroup"
  // trackpb "user-api/rpc/track"
  tagpb "user-api/rpc/tag"

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
  FeaturedTrackGroupId uuid.UUID `sql:"type:uuid,default:uuid_nil()"`

  Followers []uuid.UUID `sql:",type:uuid[]" pg:",array"`

  AdminUsers []uuid.UUID `sql:",type:uuid[]" pg:",array"`

  Members []UserGroup `pg:"many2many:user_group_members,fk:user_group_id,joinFK:member_id"`
  MemberOfGroups []UserGroup `pg:"many2many:user_group_members,fk:member_id,joinFK:user_group_id"`

  OwnerOfTracks []Track `pg:"fk:user_group_id"` // user group gets paid for these tracks
  ArtistOfTracks []uuid.UUID `sql:",type:uuid[]" pg:",array"` // user group displayed as artist for these tracks
  OwnerOfTrackGroups []TrackGroup `pg:"fk:user_group_id"` // user group owner of these track groups
  LabelOfTrackGroups []TrackGroup `pg:"fk:label_id"` // label of these track groups

  Kvstore map[string]string `pg:",hstore"`

  Publisher map[string]string `pg:",hstore"`
  Pro map[string]string `pg:",hstore"`
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

func (u *UserGroup) Create(db *pg.DB, userGroup *pb.UserGroup) (error, string) {
  var table string
  tx, err := db.Begin()
  if err != nil {
    return err, table
  }
  defer tx.Rollback()

  groupTaxonomy := new(GroupTaxonomy)
  pgerr := tx.Model(groupTaxonomy).Where("type = ?", userGroup.Type.Type).First()
  if pgerr != nil {
    return pgerr, "group_taxonomy"
  }
  u.TypeId = groupTaxonomy.Id

  var newAddress *StreetAddress
  if userGroup.Address != nil {
    newAddress = &StreetAddress{Data: userGroup.Address.Data}
    _, pgerr = tx.Model(newAddress).Returning("*").Insert()
    if pgerr != nil {
      return pgerr, "street_address"
    }
  }
  u.AddressId = newAddress.Id

  linkIds, pgerr := getLinkIds(userGroup.Links, tx)
  if pgerr != nil {
    return pgerr, "link"
  }
  u.Links = linkIds

  tagIds, pgerr := GetTagIds(userGroup.Tags, tx)
  if pgerr != nil {
    return pgerr, "tag"
  }
  u.Tags = tagIds

  recommendedArtistIds, pgerr := GetRelatedUserGroupIds(userGroup.RecommendedArtists, tx)
  if pgerr != nil {
    return pgerr, "user_group"
  }
  u.RecommendedArtists = recommendedArtistIds

  _, pgerr = tx.Model(u).Returning("*").Insert()
  if pgerr != nil {
    return pgerr, "user_group"
  }

  _, pgerr = tx.Exec(`
    UPDATE user_groups
    SET recommended_by = (select array_agg(distinct e) from unnest(recommended_by || ?) e)
    WHERE id IN (?)
  `, pg.Array([]uuid.UUID{u.Id}), pg.In(recommendedArtistIds))
  if pgerr != nil {
    return pgerr, "user_group"
  }

  pgerr = tx.Model(u).
    Column("Privacy").
    WherePK().
    Select()
  if pgerr != nil {
    return pgerr, "user_group"
  }

  // Building response
  userGroup.Address.Id = u.AddressId.String()
  userGroup.Type.Id = u.TypeId.String()
  userGroup.Privacy = &pb.Privacy{
    Id: u.Privacy.Id.String(),
    Private: u.Privacy.Private,
    OwnedTracks: u.Privacy.OwnedTracks,
    SupportedArtists: u.Privacy.SupportedArtists,
  }

  return tx.Commit(), table
}

func (u *UserGroup) Update(db *pg.DB, userGroup *pb.UserGroup) (error, string) {
  var table string
  tx, err := db.Begin()
  if err != nil {
    return err, "user_group"
  }
  defer tx.Rollback()

  // Update address
  addressId, twerr := internal.GetUuidFromString(userGroup.Address.Id)
  if twerr != nil {
    return twerr, "street_address"
  }
  address := &StreetAddress{Id: addressId, Data: userGroup.Address.Data}
  _, pgerr := tx.Model(address).Column("data").WherePK().Update()
  // _, pgerr := db.Model(address).Set("data = ?", pg.Hstore(userGroup.Address.Data)).Where("id = ?id").Update()
  if pgerr != nil {
    return pgerr, "street_address"
  }

  // Update privacy
  privacyId, twerr := internal.GetUuidFromString(userGroup.Privacy.Id)
  if twerr != nil {
    return twerr, "user_group_privacy"
  }
  privacy := &UserGroupPrivacy{Id: privacyId, Private: userGroup.Privacy.Private, OwnedTracks: userGroup.Privacy.OwnedTracks, SupportedArtists: userGroup.Privacy.SupportedArtists}
  _, pgerr = tx.Model(privacy).WherePK().Returning("*").UpdateNotNull()
  if pgerr != nil {
    return pgerr, "user_group_privacy"
  }

  // Update tags
  tagIds, pgerr := GetTagIds(userGroup.Tags, tx)
  if pgerr != nil {
    return pgerr, "tag"
  }

  // Update links
  linkIds, pgerr := getLinkIds(userGroup.Links, tx)
  if pgerr != nil {
    return pgerr, "link"
  }
  // Delete links if needed
  pgerr = tx.Model(u).WherePK().Column("links").Select()
  if pgerr != nil {
    return pgerr, "user_group"
  }
  linkIdsToDelete := internal.Difference(u.Links, linkIds)
  if len(linkIdsToDelete) > 0 {
    _, pgerr = tx.Model((*Link)(nil)).
      Where("id in (?)", pg.In(linkIdsToDelete)).
      Delete()
    if pgerr != nil {
      return pgerr, "link"
    }
  }

  // Update user group
  u.Tags = tagIds
  u.Links = linkIds
  // u.RecommendedArtists = recommendedArtistIds
  u.UpdatedAt = time.Now()
  _, pgerr = tx.Model(u).
    Column("updated_at", "pro", "publisher", "links", "tags", "display_name", "avatar", "description", "short_bio", "banner", "group_email_address").
    WherePK().
    Returning("*").
    Update()
  if pgerr != nil {
    return pgerr, "user_group"
  }

  return tx.Commit(), table
}

func SearchUserGroups(query string, db *pg.DB) (*tagpb.SearchResults, twirp.Error) {
  var userGroups []UserGroup

  pgerr := db.Model(&userGroups).
    Column("user_group.id", "user_group.display_name", "user_group.avatar", "Privacy", "Type").
    Where("to_tsvector('english'::regconfig, COALESCE(display_name, '') || ' ' || COALESCE(f_arr2str(tags), '')) @@ (plainto_tsquery('english'::regconfig, ?)) = true", query).
    Where("privacy.private = false").
    Select()
  if pgerr != nil {
    return nil, internal.CheckError(pgerr, "user_group")
  }

  var people []*tagpb.RelatedUserGroup
  var artists []*tagpb.RelatedUserGroup
  var labels []*tagpb.RelatedUserGroup
  for _, userGroup := range userGroups {
    searchUserGroup := &tagpb.RelatedUserGroup{
      Id: userGroup.Id.String(),
      DisplayName: userGroup.DisplayName,
      Avatar: userGroup.Avatar,
    }
    switch userGroup.Type.Type {
    case "user":
       people = append(people, searchUserGroup)
    case "artist":
      artists = append(artists, searchUserGroup)
    case "label":
      labels = append(labels, searchUserGroup)
    }
  }
  return &tagpb.SearchResults{
    People: people,
    Artists: artists,
    Labels: labels,
  }, nil
}

func (u *UserGroup) Delete(tx *pg.Tx) (error, string) {
  pgerr := tx.Model(u).
    Column("user_group.links","user_group.followers", "user_group.recommended_by", "user_group.recommended_artists", "Address", "Privacy",
      "OwnerOfTrackGroups", "LabelOfTrackGroups", "user_group.artist_of_tracks").
    WherePK().
    Select()
  if pgerr != nil {
    return pgerr, "user_group"
  }

  // These tracks contain the user group to delete as artist
  // so we have to remove it from the tracks' artists list
  if len(u.ArtistOfTracks) > 0 {
    _, pgerr = tx.Exec(`
      UPDATE tracks
      SET artists = array_remove(artists, ?)
      WHERE id IN (?)
    `, u.Id, pg.In(u.ArtistOfTracks))
    if pgerr != nil {
      return pgerr, "track"
    }
  }

  // These track groups contain the user group to delete as label
  // so we have to set their label_id as null
  if len(u.LabelOfTrackGroups) > 0 {
    _, pgerr = tx.Model(&u.LabelOfTrackGroups).
      Set("label_id = uuid_nil()").
      Update()
    if pgerr != nil {
      return pgerr, "track"
    }
  }

  // Delete track groups owned by user group to delete
  // if a track is a release (lp, ep, single), its tracks are owned by the same user group
  // and they'll be deleted as well
  for _, trackGroup := range(u.OwnerOfTrackGroups) {
    pgerr, table := trackGroup.Delete(tx)
    if pgerr != nil {
      return pgerr, table
    }
  }

  if len(u.Links) > 0 {
    _, pgerr = tx.Model((*Link)(nil)).
      Where("id in (?)", pg.In(u.Links)).
      Delete()
    if pgerr != nil {
      return pgerr, "link"
    }
  }

  if len(u.RecommendedBy) > 0 {
    _, pgerr = tx.Exec(`
      UPDATE user_groups
      SET recommended_artists = array_remove(recommended_artists, ?)
      WHERE id IN (?)
    `, u.Id, pg.In(u.RecommendedBy))
    if pgerr != nil {
      return pgerr, "user_group"
    }
  }

  if len(u.RecommendedArtists) > 0 {
    _, pgerr = tx.Exec(`
      UPDATE user_groups
      SET recommended_by = array_remove(recommended_by, ?)
      WHERE id IN (?)
    `, u.Id, pg.In(u.RecommendedArtists))
    if pgerr != nil {
      return pgerr, "user_group"
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
func GetRelatedUserGroups(ids []uuid.UUID, db *pg.DB) ([]*tagpb.RelatedUserGroup, error) {
	groupsResponse := make([]*tagpb.RelatedUserGroup, len(ids))
	if len(ids) > 0 {
		var groups []UserGroup
		pgerr := db.Model(&groups).
			Where("id in (?)", pg.In(ids)).
			Select()
		if pgerr != nil {
			return nil, pgerr
		}
		for i, group := range groups {
			groupsResponse[i] = &tagpb.RelatedUserGroup{
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
func GetRelatedUserGroupIds(userGroups []*tagpb.RelatedUserGroup, db *pg.Tx) ([]uuid.UUID, error) {
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

func getLinkIds(l []*pb.Link, db *pg.Tx) ([]uuid.UUID, error) {
	links := make([]*Link, len(l))
	linkIds := make([]uuid.UUID, len(l))
	for i, link := range l {
		if link.Id == "" {
			links[i] = &Link{Platform: link.Platform, Uri: link.Uri}
			_, pgerr := db.Model(links[i]).Returning("*").Insert()
			if pgerr != nil {
				return nil, pgerr
			}
			linkIds[i] = links[i].Id
			link.Id = links[i].Id.String()
		} else {
			linkId, twerr := internal.GetUuidFromString(link.Id)
			if twerr != nil {
				return nil, twerr.(error)
			}
			linkIds[i] = linkId
		}
	}
	return linkIds, nil
}

/*type TrackAnalytics struct {
  Id uuid.UUID
  Title string
  PaidPlays int32
  FreePlays int32
  TotalCredits float32
}

// DEPRECATED - moved to Payment API
func (u *UserGroup) GetUserGroupTrackAnalytics(db *pg.DB) ([]*pb.TrackAnalytics, twirp.Error) {
  pgerr := db.Model(u).
    Column("OwnerOfTracks").
    WherePK().
    Select()
  if pgerr != nil {
    return nil, internal.CheckError(pgerr, "user_group")
  }
  tracks := make([]TrackAnalytics, len(u.OwnerOfTracks))
  trackIds := make([]uuid.UUID, len(u.OwnerOfTracks))
  for i, track := range(u.OwnerOfTracks) {
    tracks[i] = TrackAnalytics{
      Title: track.Title,
    }
    trackIds[i] = track.Id
  }
  artistTrackAnalytics := make([]*pb.TrackAnalytics, len(tracks))

  if len(u.OwnerOfTracks) > 0 {
    _, pgerr := db.Query(&tracks, `
      SELECT play.track_id AS id,
        count(case when play.type = 'paid' then 1 else null end) AS paid_plays,
        count(case when play.type = 'free' then 1 else null end) AS free_plays,
        SUM(play.credits) AS total_credits
      FROM plays AS play
      WHERE play.track_id IN (?)
      GROUP BY play.track_id
    `, pg.In(trackIds))
    if pgerr != nil {
      return nil, internal.CheckError(pgerr, "play")
    }
    for i, track := range(tracks) {
      artistTrackAnalytics[i] = &pb.TrackAnalytics{
        Id: track.Id.String(),
        Title: track.Title,
        TotalPlays: track.PaidPlays + track.FreePlays,
        PaidPlays: track.PaidPlays,
        FreePlays: track.FreePlays,
        TotalCredits: float32(track.TotalCredits),
        UserGroupCredits: 0.7*float32(track.TotalCredits),
        ResonateCredits: 0.3*float32(track.TotalCredits),
      }
    }
  }

  return artistTrackAnalytics, nil
}*/
