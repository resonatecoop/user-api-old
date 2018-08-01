package models

import (
  "time"

  "github.com/satori/go.uuid"
  "user-api/internal"
  "github.com/go-pg/pg"
  "github.com/twitchtv/twirp"
  pb "user-api/rpc/trackgroup"
  trackpb "user-api/rpc/track"
)

type TrackGroup struct {
  Id uuid.UUID `sql:"type:uuid,default:uuid_generate_v4()"`
  CreatedAt time.Time `sql:"default:now()"`
  UpdatedAt time.Time

  Title string `sql:",notnull"`
  ReleaseDate time.Time `sql:",notnull"`
  Type string `sql:"type:track_group_type,notnull"` // EP, LP, Single, Playlist
  Cover []byte `sql:",notnull"`
  DisplayArtist string // for display purposes, e.g. "Various" for compilation
  MultipleComposers bool `sql:",notnull"`
  Private bool `sql:",notnull"`

  CreatorId uuid.UUID `sql:"type:uuid,notnull"`
  Creator *User

  UserGroupId uuid.UUID `sql:"type:uuid,default:uuid_nil()"` // track group belongs to user group, can be null if user playlist
  LabelId uuid.UUID `sql:"type:uuid,default:uuid_nil()"`

  Tracks []uuid.UUID `sql:",type:uuid[]" pg:",array"`
  Tags []uuid.UUID `sql:",type:uuid[]" pg:",array"`

  // TerritoriesIncl []string `pg:",array"`
  // CLineYear time.Time
  // PLineYear time.Time
  // CLineText string
  // PLineText string
  // RightExpiryDate time.Time
  // TotalVolumes int
  // CatalogNumber string
}

func GetTrackGroups(ids []uuid.UUID, db *pg.DB, playlists bool) ([]*trackpb.RelatedTrackGroup, twirp.Error) {
	var trackGroupsResponse []*trackpb.RelatedTrackGroup
	if len(ids) > 0 {
		var t []TrackGroup
    var types []string
    if playlists == true {
      types = []string{"playlist"}
    } else {
      types = []string{"lp", "ep", "single"}
    }
		pgerr := db.Model(&t).
			Where("id in (?)", pg.In(ids)).
      Where("type in (?)", pg.In(types)).
			Select()
		if pgerr != nil {
			return nil, internal.CheckError(pgerr, "track_group")
		}
		for _, trackGroup := range t {
			trackGroupsResponse = append(trackGroupsResponse, &trackpb.RelatedTrackGroup{
        Id: trackGroup.Id.String(),
        Title: trackGroup.Title,
        Cover: trackGroup.Cover,
      })
		}
	}

	return trackGroupsResponse, nil
}

func GetTrackGroupIds(t []*pb.TrackGroup, db *pg.Tx) ([]uuid.UUID, error) {
	trackGroups := make([]*TrackGroup, len(t))
	trackGroupIds := make([]uuid.UUID, len(t))
	for i, trackGroup := range t {
		id, twerr := internal.GetUuidFromString(trackGroup.Id)
		if twerr != nil {
			return nil, twerr.(error)
		}
		trackGroups[i] = &TrackGroup{Id: id}
		pgerr := db.Model(trackGroups[i]).
			WherePK().
			Returning("id", "title", "cover").
			Select()
		if pgerr != nil {
			return nil, pgerr
		}
		trackGroup.Title = trackGroups[i].Title
		trackGroup.Cover = trackGroups[i].Cover
		trackGroupIds[i] = trackGroups[i].Id
	}
	return trackGroupIds, nil
}

func (t *TrackGroup) Create(db *pg.DB, trackGroup *pb.TrackGroup) (error, string) {
  var table string
  tx, err := db.Begin()
  if err != nil {
    return err, table
  }
  defer tx.Rollback()

  err, table = t.GetIds(trackGroup)
  if err != nil {
    return err, table
  }

  // make sure owner user group (with UserGroupId) and label (LabelId) exists if specified
  var userGroupIds []uuid.UUID
  if trackGroup.UserGroupId != "" {
    userGroupIds = append(userGroupIds, t.UserGroupId)
  }
  if trackGroup.LabelId != "" {
    userGroupIds = append(userGroupIds, t.LabelId)
  }
  userGroupIds = internal.RemoveDuplicates(userGroupIds)

  for _, id := range userGroupIds {
    u := &UserGroup{Id: id}
    pgerr := tx.Model(u).WherePK().Select()
    if pgerr != nil {
      return pgerr, "user_group"
    }
  }

  // Select or create tags
  tagIds, pgerr := GetTagIds(trackGroup.Tags, tx)
  if pgerr != nil {
    return pgerr, "tag"
  }
  t.Tags = tagIds

  // Add tracks
  trackIds, pgerr, table := GetTrackIds(trackGroup.Tracks, tx)
  if pgerr != nil {
    return pgerr, table
  }
  t.Tracks = trackIds

  // Insert track group
  _, pgerr = tx.Model(t).Returning("*").Insert()
  if pgerr != nil {
    return pgerr, "track_group"
  }
  trackGroup.Id = t.Id.String()

  trackGroupIdArr := []uuid.UUID{t.Id}
  // Add track group to owner user group/label track_groups if exist
  if trackGroup.UserGroupId != "" || trackGroup.LabelId != "" {
    _, pgerr = tx.Exec(`
      UPDATE user_groups
      SET track_groups = (select array_agg(distinct e) from unnest(track_groups || ?) e)
      WHERE id IN (?)
    `, pg.Array(trackGroupIdArr), pg.In(userGroupIds))
    if pgerr != nil {
      return pgerr, "user_group"
    }
  }

  // Add track group to user playlists if of type playlist
  if trackGroup.Type == "playlist" {
    _, pgerr = tx.Exec(`
      UPDATE users
      SET playlists = (select array_agg(distinct e) from unnest(playlists || ?) e)
      WHERE id = ?
    `, pg.Array(trackGroupIdArr), t.CreatorId)
    if pgerr != nil {
      return pgerr, "user"
    }
  }

  // Add track group to tracks
  if len(trackGroup.Tracks) > 0 {
    _, pgerr = tx.Exec(`
      UPDATE tracks
      SET track_groups = (select array_agg(distinct e) from unnest(track_groups || ?) e)
      WHERE id IN (?)
    `, pg.Array(trackGroupIdArr), pg.In(trackIds))
    if pgerr != nil {
      return pgerr, "user"
    }
  }

  return tx.Commit(), table
}

func (t *TrackGroup) Update(db *pg.DB, trackGroup *pb.TrackGroup) (error, string) {
  var table string
  tx, err := db.Begin()
  if err != nil {
    return err, table
  }
  defer tx.Rollback()

  return tx.Commit(), table
}

func (t *TrackGroup) Delete(db *pg.DB, trackGroup *pb.TrackGroup) (error, string) {
  var table string
  tx, err := db.Begin()
  if err != nil {
    return err, table
  }
  defer tx.Rollback()

  return tx.Commit(), table
}

func (t *TrackGroup) GetIds(trackGroup *pb.TrackGroup) (error, string) {
  creatorId, err := internal.GetUuidFromString(trackGroup.CreatorId)
  if err != nil {
    return err, "owner"
  }

  if trackGroup.UserGroupId != "" {
    userGroupId, err := internal.GetUuidFromString(trackGroup.UserGroupId)
    if err != nil {
      return err, "user_group"
    }
    t.UserGroupId = userGroupId
  }

  if trackGroup.LabelId != "" {
    labelId, err := internal.GetUuidFromString(trackGroup.LabelId)
    if err != nil {
      return err, "user_group"
    }
    t.LabelId = labelId
  }

  t.CreatorId = creatorId
  return nil, ""
}
