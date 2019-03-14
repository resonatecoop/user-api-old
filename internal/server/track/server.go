package trackserver

import (
	// "fmt"
	// "time"
	"context"

	"github.com/go-pg/pg"
	"github.com/twitchtv/twirp"
	"github.com/satori/go.uuid"

	// userpb "user-api/rpc/user"
	pb "user-api/rpc/track"
	tagpb "user-api/rpc/tag"
	"user-api/internal"
	"user-api/internal/database/model"
)

type Server struct {
	db *pg.DB
}

func NewServer(db *pg.DB) *Server {
	return &Server{db: db}
}

func (s *Server) GetTracks(ctx context.Context, req *pb.TracksList) (*pb.TracksList, error) {
	trackIds := make([]uuid.UUID, len(req.Tracks))
	for i, track := range req.Tracks {
		id, twerr := internal.GetUuidFromString(track.Id)
		if twerr != nil {
			return nil, twerr
		}
		trackIds[i] = id
	}
	tracksResponse, twerr := model.GetTracks(trackIds, s.db, true, ctx)
	if twerr != nil {
		return nil, twerr
	}
	return &pb.TracksList{
		Tracks: tracksResponse,
	}, nil
}

func (s *Server) SearchTracks(ctx context.Context, q *tagpb.Query) (*tagpb.SearchResults, error) {
  if len(q.Query) < 3 {
    return nil, twirp.InvalidArgumentError("query", "must be a valid search query")
  }

  searchResults, twerr := model.SearchTracks(q.Query, s.db)
  if twerr != nil {
    return nil, twerr
  }
  return searchResults, nil
}

func (s *Server) CreateTrack(ctx context.Context, track *pb.Track) (*pb.Track, error) {
  // Track is created then added to a TrackGroup on track group creation
  err := checkRequiredAttributes(track)
  if err != nil {
    return nil, err
  }

  t := &model.Track{
    Title: track.Title,
    Status: track.Status,
    Enabled: track.Enabled,
    TrackNumber: track.TrackNumber,
    Duration: track.Duration,
  }

  if pgerr, table := t.Create(s.db, track); pgerr != nil {
    return nil, internal.CheckError(pgerr, table)
  }

  return track, nil
}

func (s *Server) UpdateTrack(ctx context.Context, track *pb.Track) (*tagpb.Empty, error) {
	err := checkRequiredAttributes(track)
	if err != nil {
		return nil, err
	}

	t, err := getTrackModel(track)
	if err != nil {
		return nil, err
	}

	if pgerr, table := t.Update(s.db, track); pgerr != nil {
    return nil, internal.CheckError(pgerr, table)
  }
  return &tagpb.Empty{}, nil
}

func (s *Server) DeleteTrack(ctx context.Context, track *pb.Track) (*tagpb.Empty, error) {
	t, twerr := getTrackModel(track)
	if twerr != nil {
		return nil, twerr
	}

	tx, err := s.db.Begin()
	if err != nil {
		return nil, internal.CheckError(err, "")
	}
	defer tx.Rollback()

	if pgerr, table := t.Delete(tx); pgerr != nil {
		return nil, internal.CheckError(pgerr, table)
	}
	err = tx.Commit()
  if err != nil {
    return nil, internal.CheckError(err, "")
  }
	return &tagpb.Empty{}, nil
}

func getTrackModel(track *pb.Track) (*model.Track, twirp.Error) {
  id, err := internal.GetUuidFromString(track.Id)
  if err != nil {
    return nil, err
  }
  return &model.Track{
    Id: id,
    Title: track.Title,
    Status: track.Status,
    Enabled: track.Enabled,
    TrackNumber: track.TrackNumber,
		Duration: track.Duration,
  }, nil
}

func checkRequiredAttributes(track *pb.Track) (twirp.Error) {
	if track.Title == "" || track.Status == "" || track.TrackNumber == 0 || track.CreatorId == "" || track.UserGroupId == "" { // track.Artists?
		var argument string
		switch {
		case track.Title == "":
			argument = "title"
		case track.Status == "":
			argument = "status"
		case track.CreatorId == "":
			argument = "creator_id"
		case track.UserGroupId == "":
			argument = "user_group_id"
		case track.TrackNumber == 0:
			argument = "track_number"
		}
		return twirp.RequiredArgumentError(argument)
	}
	return nil
}

/*func (s *Server) GetTrack(ctx context.Context, track *pb.Track) (*pb.Track, error) {
	t, err := getTrackModel(track)
	if err != nil {
		return nil, err
	}

	pgerr := s.db.Model(t).
			Column("track.*").
			WherePK().
			Select()
	if pgerr != nil {
		return nil, internal.CheckError(pgerr, "track")
	}
	track.UserGroupId = t.UserGroupId.String()
	track.CreatorId = t.CreatorId.String()
	track.TrackServerId = t.TrackServerId.String()
	track.Title = t.Title
	track.Status = t.Status
	track.Enabled = t.Enabled
	track.TrackNumber = t.TrackNumber
	track.Duration = t.Duration

	// Get tags
	tags, twerr := model.GetTags(t.Tags, s.db)
	if twerr != nil {
		return nil, twerr
	}
	track.Tags = tags

  // Get artists (id, name, avatar)
	artists, pgerr := model.GetRelatedUserGroups(t.Artists, s.db)
	if pgerr != nil {
		return nil, internal.CheckError(pgerr, "user_group")
	}
	track.Artists = artists

  // Get track_groups (id, title, cover) that are not playlists (i.e. LP, EP or Single)
	trackGroups, twerr := model.GetTrackGroupsFromIds(t.TrackGroups, s.db, []string{"lp", "ep", "single"})
	if twerr != nil {
		return nil, twerr
	}
	track.TrackGroups = trackGroups

	return track, nil
}*/
