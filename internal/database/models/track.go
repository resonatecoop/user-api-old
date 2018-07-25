package models

import (
  "time"
  // "fmt"
  // "log"
  "github.com/satori/go.uuid"
  pb "user-api/rpc/track"
  "github.com/go-pg/pg"
  "user-api/internal"
)


type Track struct {
  Id uuid.UUID `sql:"type:uuid,default:uuid_generate_v4()"`
  CreatedAt time.Time `sql:"default:now()"`
	UpdatedAt time.Time

  Title string `sql:",notnull"`
  Status string  `sql:"type:status,notnull"`
  Enabled bool `sql:",notnull"`
  TrackNumber int32 `sql:",notnull"`
  Duration float32

  TrackGroups []uuid.UUID `sql:",type:uuid[]" pg:",array"`
  FavoriteOfUsers []uuid.UUID `sql:",type:uuid[]" pg:",array"`

  TrackServerId uuid.UUID `sql:"type:uuid,notnull"`

  CreatorId uuid.UUID `sql:"type:uuid,notnull"`
  Creator *User

  UserGroupId uuid.UUID `sql:"type:uuid,notnull"`
  UserGroup *UserGroup // track belongs to user group (the one who gets paid)

  Artists []uuid.UUID `sql:",type:uuid[]" pg:",array"` // for display purposes
  Tags []uuid.UUID `sql:",type:uuid[]" pg:",array"`

  // Composers with IPI
  // Performers with IPI
}

func (t *Track) Create(db *pg.DB, track *pb.Track) (error, string) {
  var table string
  tx, err := db.Begin()
  if err != nil {
    return err, table
  }
  defer tx.Rollback()

  creatorId, err := internal.GetUuidFromString(track.CreatorId)
  if err != nil {
    return err, "user"
  }

  userGroupId, err := internal.GetUuidFromString(track.UserGroupId)
  if err != nil {
    return err, "user_group"
  }
  artistIds, pgerr := GetRelatedUserGroupIds(track.Artists, tx)
  if pgerr != nil {
    return pgerr, "user_group"
  }

  if track.TrackServerId != "" {
    trackServerId, err := internal.GetUuidFromString(track.TrackServerId)
    if err != nil {
      return err, "track_server"
    }
    t.TrackServerId = trackServerId
  }

  // Select or create tags
  tagIds, pgerr := GetTagIds(track.Tags, tx)
  if pgerr != nil {
    return pgerr, "tag"
  }

  t.Tags = tagIds
  t.UserGroupId = userGroupId
  t.CreatorId = creatorId
  t.Artists = artistIds

  _, pgerr = tx.Model(t).Returning("*").Insert()

  if pgerr != nil {
    return pgerr, "track"
  }
  track.Id = t.Id.String()

  // Add track to involved user groups
  // userGroupId can be part of artistIds (artist adding his/her track)
  // or not (label adding a track for one or more artists)
  userGroupIds := internal.RemoveDuplicates(append(artistIds, userGroupId))
  trackIdArr := []uuid.UUID{t.Id}
  _, pgerr = tx.ExecOne(`
    UPDATE user_groups
    SET tracks = (select array_agg(distinct e) from unnest(tracks || ?) e)
    WHERE id IN (?)
  `, pg.Array(trackIdArr), pg.In(userGroupIds))
  if pgerr != nil {
    return pgerr, "user_group"
  }

  return tx.Commit(), table
}
