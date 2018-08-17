package models

import (
  "time"
  "context"
  // "fmt"
  "net/url"
  // "log"

  "github.com/satori/go.uuid"
  pb "user-api/rpc/track"
  // tagpb "user-api/rpc/tag"
  "github.com/go-pg/pg"
  "github.com/go-pg/pg/orm"
  "github.com/twitchtv/twirp"

  "user-api/internal"
)


type Track struct {
  Id uuid.UUID `sql:"type:uuid,default:uuid_generate_v4()"`
  CreatedAt time.Time `sql:"default:now()"`
	UpdatedAt time.Time

  Title string `sql:",notnull"`
  Status string  `sql:"type:track_status,notnull"`
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

  Plays []User `pg:"many2many:plays"`
  // Composers with IPI
  // Performers with IPI
}

// TODO add user fav and play count
func GetTracks(ids []uuid.UUID, db *pg.DB, showTrackGroup bool, ctx context.Context) ([]*pb.Track, twirp.Error) {
	var tracksResponse []*pb.Track
	if len(ids) > 0 {
		var t []Track

    pagination := func(q *orm.Query) (*orm.Query, error) {
      if ctx.Value("query") != nil {
        q = q.Apply(orm.Pagination(ctx.Value("query").(url.Values)))
      }
      return q, nil
    }

		pgerr := db.Model(&t).
			Where("id in (?)", pg.In(ids)).
      Join("JOIN unnest(?::uuid[]) with ordinality t(id, ord) using (id)", pg.Array(ids)).
      Order("t.ord").
      Apply(pagination).
			Select()
		if pgerr != nil {
			return nil, internal.CheckError(pgerr, "track")
		}
		for _, track := range t {
      trackResponse := &pb.Track{
        Id: track.Id.String(),
        Title: track.Title,
        TrackServerId: track.TrackServerId.String(),
        Duration: track.Duration,
        Status: track.Status,
        TrackNumber: track.TrackNumber,
      }
      artists, pgerr := GetRelatedUserGroups(track.Artists, db)
      if pgerr != nil {
        return  nil, internal.CheckError(pgerr, "user_group")
      }
      trackResponse.Artists = artists
      if showTrackGroup == true {
        trackGroups, twerr := GetTrackGroups(track.TrackGroups, db, []string{"lp", "ep", "single"})
        if twerr != nil {
          return  nil, twerr
        }
        trackResponse.TrackGroups = trackGroups
      }
			tracksResponse = append(tracksResponse, trackResponse)
		}
	}

	return tracksResponse, nil
}

func GetTrackIds(t []*pb.Track, tx *pg.Tx) ([]uuid.UUID, error, string) {
  tracks := make([]*Track, len(t))
  trackIds := make([]uuid.UUID, len(t))
  for i, track := range t {
    id, twerr := internal.GetUuidFromString(track.Id)
    if twerr != nil {
      return nil, twerr.(error), "track"
    }
    tracks[i] = &Track{Id: id}
    pgerr := tx.Model(tracks[i]).
      WherePK().
      Select()
    if pgerr != nil {
      return nil, pgerr, "track"
    }
    artists, pgerr := GetRelatedUserGroups(tracks[i].Artists, tx.DB())
    if pgerr != nil {
      return  nil, pgerr, "user_group"
    }
    track.Artists = artists
    track.Title = tracks[i].Title
    track.TrackServerId = tracks[i].TrackServerId.String()
    track.Duration = tracks[i].Duration
    track.TrackNumber = tracks[i].TrackNumber
    track.Status = tracks[i].Status
    trackIds[i] = tracks[i].Id
  }
  return trackIds, nil, ""
}

func (t *Track) Update(db *pg.DB, track *pb.Track) (error, string) {
  // Update tags? artists? might not need tx here if not
  // tracks can be added to a track group from dedicated endpoint
  // in TrackGroup Service AddTracksToTrackGroup
  // or on TrackGroup creation (TrackGroup Service CreateTrackGroup)
  var table string
  tx, err := db.Begin()
  if err != nil {
    return err, table
  }
  defer tx.Rollback()

  err, table = t.GetIds(track)
  if err != nil {
    return err, table
  }

  t.UpdatedAt = time.Now()
  _, pgerr := tx.Model(t).
    Column("title", "updated_at", "status", "track_number", "duration", "track_server_id").
    WherePK().
    Returning("*").
    Update()
  if pgerr != nil {
    return pgerr, "track"
  }

  return tx.Commit(), table
}

func (t *Track) Create(db *pg.DB, track *pb.Track) (error, string) {
  var table string
  tx, err := db.Begin()
  if err != nil {
    return err, table
  }
  defer tx.Rollback()

  err, table = t.GetIds(track)
  if err != nil {
    return err, table
  }

  artistIds, pgerr := GetRelatedUserGroupIds(track.Artists, tx)
  if pgerr != nil {
    return pgerr, "user_group"
  }

  // Select or create tags
  tagIds, pgerr := GetTagIds(track.Tags, tx)
  if pgerr != nil {
    return pgerr, "tag"
  }

  t.Tags = tagIds
  t.Artists = artistIds

  _, pgerr = tx.Model(t).Returning("*").Insert()

  if pgerr != nil {
    return pgerr, "track"
  }
  track.Id = t.Id.String()

  // Add track to involved user groups
  // userGroupId can be part of artistIds (artist adding his/her track)
  // or not (label adding a track for one or more artists)
  userGroupIds := internal.RemoveDuplicates(append(artistIds, t.UserGroupId))
  trackIdArr := []uuid.UUID{t.Id}
  _, pgerr = tx.Exec(`
    UPDATE user_groups
    SET tracks = (select array_agg(distinct e) from unnest(tracks || ?) e)
    WHERE id IN (?)
  `, pg.Array(trackIdArr), pg.In(userGroupIds))
  if pgerr != nil {
    return pgerr, "user_group"
  }

  return tx.Commit(), table
}

func (t *Track) Delete(tx *pg.Tx) (error, string) {
  var table string
  // Delete from track server?
  pgerr := tx.Model(t).WherePK().Select()
  if pgerr != nil {
    return pgerr, "track"
  }

  // Delete track from user_group/artists tracks array
  userGroupIds := internal.RemoveDuplicates(append(t.Artists, t.UserGroupId))
  _, pgerr = tx.Exec(`
    UPDATE user_groups
    SET tracks = array_remove(tracks, ?)
    WHERE id IN (?)
  `, t.Id, pg.In(userGroupIds))
  if pgerr != nil {
    return pgerr, "user_group"
  }

  // Delete track from track_groups tracks array
  if len(t.TrackGroups) > 0 {
    _, pgerr = tx.Exec(`
      UPDATE track_groups
      SET tracks = array_remove(tracks, ?)
      WHERE id IN (?)
    `, t.Id, pg.In(t.TrackGroups))
    if pgerr != nil {
      return pgerr, "track_group"
    }
  }

  // Delete track from user favorite_tracks array
  if len(t.FavoriteOfUsers) > 0 {
    _, pgerr = tx.Exec(`
      UPDATE users
			SET favorite_tracks = array_remove(favorite_tracks, ?)
			WHERE id IN (?)
    `, t.Id, pg.In(t.FavoriteOfUsers))
    if pgerr != nil {
      return pgerr, "user_group"
    }
  }

  pgerr = tx.Delete(t)
  if pgerr != nil {
    return pgerr, "track"
  }

  return nil, table
}

// Get track OwnerId, UserGroupId and TrackServerId as UUID
func (t *Track) GetIds(track *pb.Track) (error, string) {
  creatorId, err := internal.GetUuidFromString(track.CreatorId)
  if err != nil {
    return err, "owner"
  }

  userGroupId, err := internal.GetUuidFromString(track.UserGroupId)
  if err != nil {
    return err, "user_group"
  }

  if track.TrackServerId != "" {
    trackServerId, err := internal.GetUuidFromString(track.TrackServerId)
    if err != nil {
      return err, "track_server"
    }
    t.TrackServerId = trackServerId
  }

  t.UserGroupId = userGroupId
  t.CreatorId = creatorId

  return nil, ""
}
