package trackserver

import (
	// "fmt"
	// "time"
	"context"

	"github.com/go-pg/pg"
	"github.com/twitchtv/twirp"
	// "github.com/satori/go.uuid"

	userpb "user-api/rpc/user"
	pb "user-api/rpc/track"
	"user-api/internal"
	"user-api/internal/database/models"
)

type Server struct {
	db *pg.DB
}

func NewServer(db *pg.DB) *Server {
	return &Server{db: db}
}

func (s *Server) GetTrack(ctx context.Context, track *pb.Track) (*pb.Track, error) {
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
	tags, twerr := models.GetTags(t.Tags, s.db)
	if twerr != nil {
		return nil, twerr
	}
	track.Tags = tags

  // Get artists (id, name, avatar)
	artists, pgerr := models.GetRelatedUserGroups(t.Artists, s.db)
	if pgerr != nil {
		return nil, internal.CheckError(pgerr, "user_group")
	}
	track.Artists = artists

  // Get track_groups (id, title, cover) that are not playlists (i.e. LP, EP or Single)
	trackGroups, twerr := models.GetTrackGroups(t.TrackGroups, s.db, false)
	if twerr != nil {
		return nil, twerr
	}
	track.TrackGroups = trackGroups

	return track, nil
}

func (s *Server) CreateTrack(ctx context.Context, track *pb.Track) (*pb.Track, error) {
  // Track is created then added to a TrackGroup on track group creation
  err := checkRequiredAttributes(track)
  if err != nil {
    return nil, err
  }

  t := &models.Track{
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

func (s *Server) UpdateTrack(ctx context.Context, track *pb.Track) (*userpb.Empty, error) {
	t, err := getTrackModel(track)
	if err != nil {
		return nil, err
	}

	if pgerr, table := t.Update(s.db, track); pgerr != nil {
    return nil, internal.CheckError(pgerr, table)
  }
  return &userpb.Empty{}, nil
}

func (s *Server) DeleteTrack(ctx context.Context, track *pb.Track) (*userpb.Empty, error) {
	t, err := getTrackModel(track)
	if err != nil {
		return nil, err
	}

	if pgerr, table := t.Delete(s.db, track); pgerr != nil {
		return nil, internal.CheckError(pgerr, table)
	}
  return &userpb.Empty{}, nil
}

func getTrackModel(track *pb.Track) (*models.Track, twirp.Error) {
  id, err := internal.GetUuidFromString(track.Id)
  if err != nil {
    return nil, err
  }
  return &models.Track{
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
