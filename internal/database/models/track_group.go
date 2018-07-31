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
