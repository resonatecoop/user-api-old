package model

import (
  "time"
  // "fmt"
  // "log"

  "github.com/satori/go.uuid"
  "github.com/go-pg/pg"
  "github.com/twitchtv/twirp"

  pb "user-api/rpc/trackgroup"
  trackpb "user-api/rpc/track"
  tagpb "user-api/rpc/tag"

  uuidpkg "user-api/pkg/uuid"
  errorpkg "user-api/pkg/error"
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
  About string

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

func SearchTrackGroups(query string, db *pg.DB,) (*tagpb.SearchResults, twirp.Error) {
  var trackGroups []TrackGroup

  pgerr := db.Model(&trackGroups).
    ColumnExpr("track_group.id, track_group.display_artist, track_group.title, track_group.tracks, track_group.user_group_id, track_group.type, track_group.cover, track_group.about").
    Where("to_tsvector('english'::regconfig, COALESCE(title, '') || ' ' || COALESCE(f_arr2str(tags), '')) @@ (plainto_tsquery('english'::regconfig, ?)) = true", query).
    Where("private = false").
    Select()
  if pgerr != nil {
    return nil, errorpkg.CheckError(pgerr, "track_group")
  }

  var playlists []*tagpb.RelatedTrackGroup
  var albums []*tagpb.RelatedTrackGroup
  for _, trackGroup := range trackGroups {
    userGroups, pgerr := GetRelatedUserGroups([]uuid.UUID{trackGroup.UserGroupId}, db)
    if pgerr != nil {
      return nil, errorpkg.CheckError(pgerr, "user_group")
    }
    searchTrackGroup := &tagpb.RelatedTrackGroup{
      Id: trackGroup.Id.String(),
      Title: trackGroup.Title,
      TotalTracks: int32(len(trackGroup.Tracks)),
      UserGroup: userGroups[0],
      Cover: trackGroup.Cover,
      DisplayArtist: trackGroup.DisplayArtist,
      Type: trackGroup.Type,
      About: trackGroup.About,
    }
    if trackGroup.Type == "playlist" {
      playlists = append(playlists, searchTrackGroup)
    } else {
      albums = append(albums, searchTrackGroup)
    }
  }
  return &tagpb.SearchResults{
    Playlists: playlists,
    Albums: albums,
  }, nil
}

// Get related track groups from ids
func GetTrackGroupsFromIds(ids []uuid.UUID, db *pg.DB, types []string) ([]*tagpb.RelatedTrackGroup, twirp.Error) {
	var trackGroupsResponse []*tagpb.RelatedTrackGroup
	if len(ids) > 0 {
		var t []TrackGroup
		pgerr := db.Model(&t).
			Where("id in (?)", pg.In(ids)).
      Where("type in (?)", pg.In(types)).
			Select()
		if pgerr != nil {
			return nil, errorpkg.CheckError(pgerr, "track_group")
		}
		for _, trackGroup := range t {
			trackGroupsResponse = append(trackGroupsResponse, getTrackGroupResponse(trackGroup))
		}
	}
	return trackGroupsResponse, nil
}

// Get related track groups from TrackGroup models
func GetTrackGroups(t []TrackGroup) ([]*tagpb.RelatedTrackGroup) {
  var trackGroupsResponse []*tagpb.RelatedTrackGroup
  trackGroupIds := map[uuid.UUID]bool{}
  for _, trackGroup := range t {
    if trackGroupIds[trackGroup.Id] == true {
      // do not append duplicate track group
    } else {
      trackGroupIds[trackGroup.Id] = true
      trackGroupsResponse = append(trackGroupsResponse, getTrackGroupResponse(trackGroup))
    }
  }
	return trackGroupsResponse
}

func getTrackGroupResponse(trackGroup TrackGroup) (*tagpb.RelatedTrackGroup) {
  return &tagpb.RelatedTrackGroup {
    Id: trackGroup.Id.String(),
    Title: trackGroup.Title,
    Cover: trackGroup.Cover,
    DisplayArtist: trackGroup.DisplayArtist,
    Type: trackGroup.Type,
    About: trackGroup.About,
    Private: trackGroup.Private,
    TotalTracks: int32(len(trackGroup.Tracks)),
  }
}

func GetTrackGroupIds(t []*pb.TrackGroup, db *pg.Tx) ([]uuid.UUID, error) {
	trackGroups := make([]*TrackGroup, len(t))
	trackGroupIds := make([]uuid.UUID, len(t))
	for i, trackGroup := range t {
		id, twerr := uuidpkg.GetUuidFromString(trackGroup.Id)
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

func userGroupsExist(tx *pg.Tx, trackGroup *pb.TrackGroup, userGroupId, labelId uuid.UUID) (error) {
  var userGroupIds []uuid.UUID
  if trackGroup.UserGroupId != "" {
    userGroupIds = append(userGroupIds, userGroupId)
  }
  if trackGroup.LabelId != "" {
    userGroupIds = append(userGroupIds, labelId)
  }
  userGroupIds = uuidpkg.RemoveDuplicates(userGroupIds)

  for _, id := range userGroupIds {
    u := &UserGroup{Id: id}
    pgerr := tx.Model(u).WherePK().Select()
    if pgerr != nil {
      return pgerr
    }
  }
  return nil
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
  pgerr := userGroupsExist(tx, trackGroup, t.UserGroupId, t.LabelId)
  if pgerr != nil {
    return pgerr, "user_group"
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
      return pgerr, "track"
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

  err, table = t.GetIds(trackGroup)
  if err != nil {
    return err, table
  }

  trackGroupToUpdate := &TrackGroup{Id: t.Id}
  pgerr := tx.Model(trackGroupToUpdate).
    Column("user_group_id", "label_id", "tracks").
    WherePK().
    Select()
  if pgerr != nil {
    return pgerr, "track_group"
  }

  // make sure owner user group (with UserGroupId) and label (LabelId) exists if specified
  pgerr = userGroupsExist(tx, trackGroup, t.UserGroupId, t.LabelId)
  if pgerr != nil {
    return pgerr, "user_group"
  }

  columns := []string{"title", "updated_at", "release_date", "cover", "display_artist", "type",
    "multiple_composers", "private", "about"}

  // Update usergroup if needed
  if trackGroupToUpdate.UserGroupId != t.UserGroupId {
    columns = append(columns, "user_group_id")

    // Update tracks user_group_id if track group is release (lp, ep, single)
    if trackGroup.Type != "playlist" && len(trackGroupToUpdate.Tracks) > 0 {
      var tracks []Track
      _, pgerr = tx.Model(&tracks).
        Set("user_group_id = ?", t.UserGroupId).
        Where("id IN (?)", pg.In(trackGroupToUpdate.Tracks)).
        Update()
      if pgerr != nil {
        return pgerr, "tracks"
      }
    }
  }

  // Update label if needed
  if trackGroupToUpdate.LabelId != t.LabelId {
    columns = append(columns, "label_id")
  }

  // Update tags
  tagIds, pgerr := GetTagIds(trackGroup.Tags, tx)
  if pgerr != nil {
    return pgerr, "tag"
  }
  if !uuidpkg.Equal(trackGroupToUpdate.Tags, tagIds) {
    t.Tags = tagIds
    columns = append(columns, "tags")
  }

  t.UpdatedAt = time.Now()
  _, pgerr = tx.Model(t).
    Column(columns...).
    WherePK().
    Returning("*").
    Update()
  if pgerr != nil {
    return pgerr, "track_group"
  }

  return tx.Commit(), table
}

func (t *TrackGroup) Delete(tx *pg.Tx) (error, string) {
  var table string

  pgerr := tx.Model(t).WherePK().Select()
  if pgerr != nil {
    return pgerr, "track_group"
  }

  // Delete playlist track group from user (CreatorId) playlists and tracks track group
  if t.Type == "playlist" {
    var user User
    _, pgerr := tx.Model(&user).
      Set("playlists = array_remove(playlists, ?)", t.Id).
      Where("id = ?", t.CreatorId).
      Update()
    if pgerr != nil {
      return pgerr, "user"
    }

    if len(t.Tracks) > 0 {
      var tracks []Track
      _, pgerr := tx.Model(&tracks).
        Set("track_groups = array_remove(track_groups, ?)", t.Id).
        Where("id IN (?)", pg.In(t.Tracks)).
        Update()
      if pgerr != nil {
        return pgerr, "track"
      }
    }
  } else { // Delete tracks if track group not of type playlist
    for _, id := range(t.Tracks) {
      track := &Track{Id: id}
      pgerr, table := track.Delete(tx)
      if pgerr != nil {
        return pgerr, table
      }
    }
  }

  pgerr = tx.Delete(t)
  if pgerr != nil {
    return pgerr, "track_group"
  }

  return nil, table
}

func (t *TrackGroup) AddTracks(db *pg.DB, tracks []*trackpb.Track) (error, string) {
  var table string
  tx, err := db.Begin()
  if err != nil {
    return err, table
  }
  defer tx.Rollback()

  trackIds, pgerr, table := GetTrackIds(tracks, tx)
  if pgerr != nil {
    return pgerr, table
  }

  // Add Tracks to Trackgroup tracks array
  // var trackGroup TrackGroup
  res, pgerr := tx.Model(t).
    Set("tracks = (select array_agg(distinct e) from unnest(tracks || ?) e)", pg.Array(trackIds)).
    WherePK().
    Update()
  if res.RowsAffected() == 0 {
    return pg.ErrNoRows, "track_group"
  }
  if pgerr != nil {
    return pgerr, "track_group"
  }

  // Add Trackgroup to Tracks track_groups array
  var tracksToUpdate []Track
  _, pgerr = tx.Model(&tracksToUpdate).
    Set("track_groups = (select array_agg(distinct e) from unnest(track_groups || ?) e)", pg.Array([]uuid.UUID{t.Id})).
    Where("id IN (?)", pg.In(trackIds)).
    Update()

  if pgerr != nil {
    return pgerr, "track"
  }

  return tx.Commit(), table
}

func (t *TrackGroup) RemoveTracks(db *pg.DB, tracks []*trackpb.Track) (error, string) {
  var table string
  tx, err := db.Begin()
  if err != nil {
    return err, table
  }
  defer tx.Rollback()

  trackIds, pgerr, table := GetTrackIds(tracks, tx)
  if pgerr != nil {
    return pgerr, table
  }

  // Remove Tracks from Trackgroup tracks array
  res, pgerr := tx.Model(t).
    Set("tracks = (select array_agg(e) from unnest(tracks) e where e <> all(?))", pg.Array(trackIds)).
    WherePK().
    Update()
  if res.RowsAffected() == 0 {
    return pg.ErrNoRows, "track_group"
  }
  if pgerr != nil {
    return pgerr, "track_group"
  }

  // Remove Trackgroup from Tracks track_groups array
  var tracksToUpdate []Track
  _, pgerr = tx.Model(&tracksToUpdate).
    Set("track_groups = array_remove(track_groups, ?)", t.Id).
    Where("id IN (?)", pg.In(trackIds)).
    Update()
  if pgerr != nil {
    return pgerr, "track"
  }

  return tx.Commit(), table
}

func (t *TrackGroup) GetIds(trackGroup *pb.TrackGroup) (error, string) {
  creatorId, err := uuidpkg.GetUuidFromString(trackGroup.CreatorId)
  if err != nil {
    return err, "owner"
  }

  if trackGroup.UserGroupId != "" {
    userGroupId, err := uuidpkg.GetUuidFromString(trackGroup.UserGroupId)
    if err != nil {
      return err, "user_group"
    }
    t.UserGroupId = userGroupId
  }

  if trackGroup.LabelId != "" {
    labelId, err := uuidpkg.GetUuidFromString(trackGroup.LabelId)
    if err != nil {
      return err, "user_group"
    }
    t.LabelId = labelId
  }

  t.CreatorId = creatorId
  return nil, ""
}

func (t *TrackGroup) UpdateUserGroupTrackGroups(tx *pg.Tx, oldUserGroupId, newUserGroupId uuid.UUID) (error, string) {
  // Verify that new user group exists
  u := &UserGroup{Id: newUserGroupId}
  pgerr := tx.Model(u).WherePK().Select()
  if pgerr != nil {
    return pgerr, "user_group"
  }

  // Remove track group from old user group track_groups array
  var oldUserGroup UserGroup
  _, pgerr = tx.Model(&oldUserGroup).
    Set("track_groups = array_remove(track_groups, ?)", t.Id).
    Where("id = ?", oldUserGroupId).
    Update()
  if pgerr != nil {
    return pgerr, "user_group"
  }

  // Add track group to new user group track_groups array
  var newUserGroup UserGroup
  _, pgerr = tx.Model(&newUserGroup).
    Set("track_groups = (select array_agg(distinct e) from unnest(track_groups || ?) e)", pg.Array([]uuid.UUID{t.Id})).
    Where("id = ?", newUserGroupId).
    Update()
  if pgerr != nil {
    return pgerr, "user_group"
  }
  return nil, ""
}
