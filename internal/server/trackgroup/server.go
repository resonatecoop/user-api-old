package trackgroupserver

import (
  "context"
  // "fmt"
  "github.com/go-pg/pg"
  "github.com/twitchtv/twirp"
  "github.com/satori/go.uuid"
  "github.com/golang/protobuf/ptypes"

  userpb "user-api/rpc/user"
  // trackpb "user-api/rpc/track"
  pb "user-api/rpc/trackgroup"
  "user-api/internal"
  "user-api/internal/database/models"
)

type Server struct {
	db *pg.DB
}

func NewServer(db *pg.DB) *Server {
	return &Server{db: db}
}

// TODO handle private playlist
func (s *Server) GetTrackGroup(ctx context.Context, trackGroup *pb.TrackGroup) (*pb.TrackGroup, error) {
  t, twerr := getTrackGroupModel(trackGroup)
  if twerr != nil {
    return nil, twerr
  }

  pgerr := s.db.Model(t).
      WherePK().
      Select()
  if pgerr != nil {
    return nil, internal.CheckError(pgerr, "track_group")
  }

  releaseDate, err := ptypes.TimestampProto(t.ReleaseDate)
  if err != nil {
    return nil, twirp.InvalidArgumentError("release_date", "must be a valid time")
  }

  trackGroup.UserGroupId = t.UserGroupId.String()
  trackGroup.CreatorId = t.CreatorId.String()
  trackGroup.LabelId = t.LabelId.String()
  trackGroup.Title = t.Title
  trackGroup.ReleaseDate = releaseDate
  trackGroup.Type = t.Type
  trackGroup.Cover = t.Cover
  trackGroup.DisplayArtist = t.DisplayArtist
  trackGroup.MultipleComposers = t.MultipleComposers
  trackGroup.Private = t.Private

  // Get tags
  tags, twerr := models.GetTags(t.Tags, s.db)
  if twerr != nil {
    return nil, twerr
  }
  trackGroup.Tags = tags

  // Get UserGroup and Label if exists
  if trackGroup.UserGroupId != "" {
    userGroup, twerr := models.GetRelatedUserGroups([]uuid.UUID{t.UserGroupId}, s.db)
    if twerr != nil {
      return nil, twerr
    }
    trackGroup.UserGroup = userGroup[0]
  }
  if trackGroup.LabelId != "" {
    label, twerr := models.GetRelatedUserGroups([]uuid.UUID{t.LabelId}, s.db)
    if twerr != nil {
      return nil, twerr
    }
    trackGroup.Label = label[0]
  }

  // Get tracks
  playlist := t.Type == "playlist"
  tracks, twerr := models.GetTracks(t.Tracks, s.db, playlist)
  if twerr != nil {
    return nil, twerr
  }
  trackGroup.Tracks = tracks

  return trackGroup, nil
}

func (s *Server) CreateTrackGroup(ctx context.Context, trackGroup *pb.TrackGroup) (*pb.TrackGroup, error) {
  return &pb.TrackGroup{}, nil
}

func (s *Server) UpdateTrackGroup(ctx context.Context, trackGroup *pb.TrackGroup) (*userpb.Empty, error) {
  return &userpb.Empty{}, nil
}

func (s *Server) DeleteTrackGroup(ctx context.Context, trackGroup *pb.TrackGroup) (*userpb.Empty, error) {
  return &userpb.Empty{}, nil
}

func (s *Server) AddTracksToTrackGroup(ctx context.Context, tracksToTrackGroup *pb.TracksToTrackGroup) (*userpb.Empty, error) {
  return &userpb.Empty{}, nil
}

func (s *Server) DeleteTracksFromTrackGroup(ctx context.Context, tracksToTrackGroup *pb.TracksToTrackGroup) (*userpb.Empty, error) {
  return &userpb.Empty{}, nil
}

func getTrackGroupModel(trackGroup *pb.TrackGroup) (*models.TrackGroup, twirp.Error) {
  id, twerr := internal.GetUuidFromString(trackGroup.Id)
  if twerr != nil {
    return nil, twerr
  }
  return &models.TrackGroup{
    Id: id,
    Title: trackGroup.Title,
    Type: trackGroup.Type,
    Cover: trackGroup.Cover,
    DisplayArtist: trackGroup.DisplayArtist,
    MultipleComposers: trackGroup.MultipleComposers,
    Private: trackGroup.Private,
  }, nil
}
