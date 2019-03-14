package trackgroupserver

import (
  "context"
  // "fmt"
  "github.com/go-pg/pg"
  "github.com/twitchtv/twirp"
  "github.com/satori/go.uuid"
  // "github.com/golang/protobuf/ptypes/timestamp"
  "github.com/golang/protobuf/ptypes"

  // userpb "user-api/rpc/user"
  // trackpb "user-api/rpc/track"
  tagpb "user-api/rpc/tag"
  pb "user-api/rpc/trackgroup"
  "user-api/internal"
  "user-api/internal/database/model"
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

  // trackGroup.UserGroupId = t.UserGroupId.String()
  trackGroup.CreatorId = t.CreatorId.String()
  // trackGroup.LabelId = t.LabelId.String()
  trackGroup.Title = t.Title
  trackGroup.About = t.About
  trackGroup.ReleaseDate = releaseDate
  trackGroup.Type = t.Type
  trackGroup.Cover = t.Cover
  trackGroup.DisplayArtist = t.DisplayArtist
  trackGroup.MultipleComposers = t.MultipleComposers
  trackGroup.Private = t.Private

  // Get tags
  tags, twerr := model.GetTags(t.Tags, s.db)
  if twerr != nil {
    return nil, twerr
  }
  trackGroup.Tags = tags

  // Get UserGroup and Label if exists
  if t.UserGroupId.String() != "" {
    userGroup, pgerr := model.GetRelatedUserGroups([]uuid.UUID{t.UserGroupId}, s.db)
    if pgerr != nil {
      return nil, internal.CheckError(pgerr, "user_group")
    }
    trackGroup.UserGroup = userGroup[0]
  }
  if t.LabelId.String() != "" {
    label, pgerr := model.GetRelatedUserGroups([]uuid.UUID{t.LabelId}, s.db)
    if pgerr != nil {
      return nil, internal.CheckError(pgerr, "user_group")
    }
    trackGroup.Label = label[0]
  }

  // Get tracks
  showTrackGroup := t.Type == "playlist"
  tracks, twerr := model.GetTracks(t.Tracks, s.db, showTrackGroup, ctx)
  if twerr != nil {
    return nil, twerr
  }
  trackGroup.Tracks = tracks

  return trackGroup, nil
}

func (s *Server) SearchTrackGroups(ctx context.Context, q *tagpb.Query) (*tagpb.SearchResults, error) {
  if len(q.Query) < 3 {
    return nil, twirp.InvalidArgumentError("query", "must be a valid search query")
  }

  searchResults, twerr := model.SearchTrackGroups(q.Query, s.db)
  if twerr != nil {
    return nil, twerr
  }
  return searchResults, nil
}

func (s *Server) CreateTrackGroup(ctx context.Context, trackGroup *pb.TrackGroup) (*pb.TrackGroup, error) {
  twerr := checkRequiredAttributes(trackGroup)
  if twerr != nil {
    return nil, twerr
  }

  t := &model.TrackGroup{
    Title: trackGroup.Title,
    Type: trackGroup.Type,
    Cover: trackGroup.Cover,
    DisplayArtist: trackGroup.DisplayArtist,
    MultipleComposers: trackGroup.MultipleComposers,
    Private: trackGroup.Private,
    About: trackGroup.About,
  }

  releaseDate, err := ptypes.Timestamp(trackGroup.ReleaseDate)
  if err != nil {
    return nil, twirp.InvalidArgumentError("release_date", "must be a valid time")
  }
  t.ReleaseDate = releaseDate

  if pgerr, table := t.Create(s.db, trackGroup); pgerr != nil {
    return nil, internal.CheckError(pgerr, table)
  }

  return trackGroup, nil
}

func (s *Server) UpdateTrackGroup(ctx context.Context, trackGroup *pb.TrackGroup) (*tagpb.Empty, error) {
  twerr := checkRequiredAttributes(trackGroup)
  if twerr != nil {
    return nil, twerr
  }

  t, twerr := getTrackGroupModel(trackGroup)
	if twerr != nil {
		return nil, twerr
	}
  releaseDate, err := ptypes.Timestamp(trackGroup.ReleaseDate)
  if err != nil {
    return nil, twirp.InvalidArgumentError("release_date", "must be a valid time")
  }
  t.ReleaseDate = releaseDate

	if pgerr, table := t.Update(s.db, trackGroup); pgerr != nil {
    return nil, internal.CheckError(pgerr, table)
  }
  return &tagpb.Empty{}, nil
}

func (s *Server) DeleteTrackGroup(ctx context.Context, trackGroup *pb.TrackGroup) (*tagpb.Empty, error) {
  t, twerr := getTrackGroupModel(trackGroup)
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

func (s *Server) AddTracksToTrackGroup(ctx context.Context, tracksToTrackGroup *pb.TracksToTrackGroup) (*tagpb.Empty, error) {
  id, twerr := internal.GetUuidFromString(tracksToTrackGroup.TrackGroupId)
  if twerr != nil {
    return nil, twerr
  }
  t := &model.TrackGroup{Id: id}

  if pgerr, table := t.AddTracks(s.db, tracksToTrackGroup.Tracks); pgerr != nil {
		return nil, internal.CheckError(pgerr, table)
	}
  return &tagpb.Empty{}, nil
}

func (s *Server) RemoveTracksFromTrackGroup(ctx context.Context, tracksToTrackGroup *pb.TracksToTrackGroup) (*tagpb.Empty, error) {
  id, twerr := internal.GetUuidFromString(tracksToTrackGroup.TrackGroupId)
  if twerr != nil {
    return nil, twerr
  }
  t := &model.TrackGroup{Id: id}

  if pgerr, table := t.RemoveTracks(s.db, tracksToTrackGroup.Tracks); pgerr != nil {
    return nil, internal.CheckError(pgerr, table)
  }
  return &tagpb.Empty{}, nil
}

func checkRequiredAttributes(trackGroup *pb.TrackGroup) (twirp.Error) {
  if trackGroup.Title == "" || (trackGroup.ReleaseDate == nil) || trackGroup.Type == "" || len(trackGroup.Cover) == 0 || trackGroup.CreatorId == "" {
    var argument string
    switch {
    case trackGroup.Title == "":
      argument = "title"
    case trackGroup.ReleaseDate == nil:
      argument = "release_date"
    case trackGroup.Type == "":
      argument = "type"
    case len(trackGroup.Cover) == 0:
      argument = "cover"
    case trackGroup.CreatorId == "":
      argument = "creator_id"
    }
    return twirp.RequiredArgumentError(argument)
  }
  // A playlist does not have necessarily a owner user group (with id UserGroupId)
  // if it is a private user playlist
  // But other types of track groups (lp, ep, single) have to belong to a user group
  if trackGroup.Type != "playlist" && trackGroup.UserGroupId == "" {
    return twirp.RequiredArgumentError("user_group_id")
  }
  return nil
}

func getTrackGroupModel(trackGroup *pb.TrackGroup) (*model.TrackGroup, twirp.Error) {
  id, twerr := internal.GetUuidFromString(trackGroup.Id)
  if twerr != nil {
    return nil, twerr
  }
  return &model.TrackGroup{
    Id: id,
    Title: trackGroup.Title,
    Type: trackGroup.Type,
    Cover: trackGroup.Cover,
    DisplayArtist: trackGroup.DisplayArtist,
    MultipleComposers: trackGroup.MultipleComposers,
    Private: trackGroup.Private,
    About: trackGroup.About,
  }, nil
}
