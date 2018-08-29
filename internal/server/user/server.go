package userserver

import (
	// "fmt"
	"time"
	"context"

	"github.com/go-pg/pg"
	"github.com/twitchtv/twirp"

	pb "user-api/rpc/user"
	tagpb "user-api/rpc/tag"
	"user-api/internal"
	"user-api/internal/database/models"
)

// Server implements the UserService
type Server struct {
	db *pg.DB
}

// NewServer creates an instance of our server
func NewServer(db *pg.DB) *Server {
	return &Server{db: db}
}

func (s *Server) GetUser(ctx context.Context, user *pb.User) (*pb.User, error) {
	u, err := getUserModel(user)
	if err != nil {
		return nil, err
	}

	pgerr := s.db.Model(u).
      Column("user.*", "OwnerOfGroups").
			Where("id = ?", u.Id).
      Select()
	twerr := internal.CheckError(pgerr, "user")
	if twerr != nil {
		return nil, twerr
	}

	return &pb.User{
		Id: u.Id.String(),
		Username: u.Username,
		// DisplayName: u.DisplayName,
		FullName: u.FullName,
		Email: u.Email,
		FirstName: u.FirstName,
		LastName: u.LastName,
		Member: u.Member,
		NewsletterNotification: u.NewsletterNotification,
		FavoriteTracks: internal.ConvertUuidToStrArray(u.FavoriteTracks),
		FollowedGroups: internal.ConvertUuidToStrArray(u.FollowedGroups),
		OwnerOfGroups: getUserGroupResponse(u.OwnerOfGroups),
	}, nil
}

func (s *Server) GetPlaylists(ctx context.Context, user *pb.User) (*pb.Playlists, error) {
	u, twerr := getUserModel(user)
	if twerr != nil {
		return nil, twerr
	}

	pgerr := s.db.Model(u).Column("user.playlists", "OwnerOfGroups").WherePK().Select()
	if pgerr != nil {
		return nil, internal.CheckError(pgerr, "user")
	}

	userPlaylists, twerr := models.GetTrackGroups(u.Playlists, s.db, []string{"playlist"})
	if twerr != nil {
		return nil, twerr
	}

	playlists := userPlaylists
	return &pb.Playlists{
		Playlists: playlists,
	}, nil
}

func (s *Server) GetFavoriteTracks(ctx context.Context, user *pb.User) (*pb.Tracks, error) {
	u, twerr := getUserModel(user)
	if twerr != nil {
		return nil, twerr
	}

	pgerr := s.db.Model(u).Column("user.favorite_tracks").WherePK().Select()
	if pgerr != nil {
		return nil, internal.CheckError(pgerr, "user")
	}

	favoriteTracks, twerr := models.GetTracks(u.FavoriteTracks, s.db, true, ctx) // will return release info (to display cover)
	if twerr != nil {
		return nil, twerr
	}

	return &pb.Tracks{
		Tracks: favoriteTracks,
	}, nil
}

/*func (s *Server) GetOwnedTracks(ctx context.Context, user *pb.User) (*pb.Tracks, error) {
	u, twerr := getUserModel(user)
	if twerr != nil {
		return nil, twerr
	}
	pgerr := s.db.Model(u).WherePK().Select()
	if pgerr != nil {
		return nil, internal.CheckError(pgerr, "user")
	}

	var ownedTrackIds []uuid.UUID
	_, pgerr = s.db.Query(&ownedTrackIds, `
    SELECT track_id FROM plays
		WHERE user_id = ? and type = 'paid'
		GROUP BY track_id
		HAVING COUNT(DISTINCT plays.id) >= 9
		ORDER BY max(created_at) desc
	`, u.Id)
	if pgerr != nil {
		return nil, internal.CheckError(pgerr, "play")
	}

	ownedTracks, twerr := models.GetTracks(ownedTrackIds, s.db, true, ctx)
	if twerr != nil {
		return nil, twerr
	}

	return &pb.Tracks{
		Tracks: ownedTracks,
	}, nil
}

func (s *Server) GetTrackHistory(ctx context.Context, user *pb.User) (*pb.Tracks, error) {
	u, twerr := getUserModel(user)
	if twerr != nil {
		return nil, twerr
	}
	pgerr := s.db.Model(u).WherePK().Select()
	if pgerr != nil {
		return nil, internal.CheckError(pgerr, "user")
	}

	var trackIds []uuid.UUID
	_, pgerr = s.db.Query(&trackIds, `
		SELECT track_id FROM plays
		WHERE user_id = ? and type = 'paid'
		GROUP BY track_id
		ORDER BY max(created_at) desc
	`, u.Id)
	if pgerr != nil {
		return nil, internal.CheckError(pgerr, "play")
	}


	tracks, twerr := models.GetTracks(trackIds, s.db, true, ctx)
	if twerr != nil {
		return nil, twerr
	}

	return &pb.Tracks{
		Tracks: tracks,
	}, nil
}

func (s *Server) GetSupportedArtists(ctx context.Context, user *pb.User) (*pb.Artists, error) {
	u, twerr := getUserModel(user)
	if twerr != nil {
		return nil, twerr
	}
	pgerr := s.db.Model(u).WherePK().Select()
	if pgerr != nil {
		return nil, internal.CheckError(pgerr, "user")
	}

	// JOIN plays AS play ON (play.track_id = track.id AND play.user_id = ? AND play.type = 'paid')

	var userGroupIds []uuid.UUID
	_, pgerr = s.db.Query(&userGroupIds, `
		WITH owned_tracks
		AS (
			SELECT play.track_id FROM plays AS play
			WHERE play.type = 'paid' AND play.user_id = ?
			GROUP BY play.track_id
			HAVING COUNT(DISTINCT play.id) >= 9
		)
		SELECT track.user_group_id FROM tracks AS track
		JOIN user_groups AS g ON g.id = track.user_group_id
		JOIN group_taxonomies AS t ON t.id = g.type_id AND t.type = 'artist'
		WHERE track.id IN (
			SELECT track_id FROM owned_tracks
		)
		GROUP BY track.user_group_id
		HAVING COUNT(DISTINCT track.id) >= 5
	`, u.Id)

	if pgerr != nil {
		return nil, internal.CheckError(pgerr, "")
	}
	artists, pgerr := models.GetRelatedUserGroups(userGroupIds, s.db)
	if pgerr != nil {
		return nil, internal.CheckError(pgerr, "user_group")
	}
	return &pb.Artists{
		Artists: artists,
	}, nil
}*/

func (s *Server) CreateUser(ctx context.Context, user *pb.User) (*pb.User, error) {
	requiredErr := checkRequiredAttributes(user)
	if requiredErr != nil {
		return nil, requiredErr
	}

	newUser := &models.User{
		Username: user.Username,
		FullName: user.FullName,
		Email: user.Email,
		// DisplayName: user.DisplayName,
	}
	_, err := s.db.Model(newUser).Returning("*").Insert()

	twerr := internal.CheckError(err, "user")
	if twerr != nil {
		return nil, twerr
	}

	return &pb.User{
		Id: newUser.Id.String(),
		Username: newUser.Username,
		// DisplayName: newUser.DisplayName,
		FullName: newUser.FullName,
		Email: newUser.Email,
	}, nil
}

/*func (s *Server) CreatePlay(ctx context.Context, playRequest *pb.CreatePlayRequest) (*pb.CreatePlayResponse, error) {
	if playRequest.Play == nil {
		return nil, twirp.RequiredArgumentError("play")
	}

	if playRequest.Play.Type == "" {
		return nil, twirp.RequiredArgumentError("type")
	}

	userId, twerr := internal.GetUuidFromString(playRequest.Play.UserId)
	if twerr != nil {
		return nil, twerr
	}
	user := &models.User{Id: userId}
	pgerr := s.db.Model(user).WherePK().Select()
	if pgerr != nil {
		return nil, internal.CheckError(pgerr, "user")
	}
	trackId, twerr := internal.GetUuidFromString(playRequest.Play.TrackId)
	if twerr != nil {
		return nil, twerr
	}
	track := &models.Track{Id: trackId}
	pgerr = s.db.Model(track).WherePK().Select()
	if pgerr != nil {
		return nil, internal.CheckError(pgerr, "track")
	}

	newPlay := &models.Play{
		UserId: userId,
		TrackId: trackId,
		Type: playRequest.Play.Type,
		Credits: playRequest.Play.Credits,
	}

	_, pgerr = s.db.Model(newPlay).Returning("*").Insert()
	if pgerr != nil {
		return nil, internal.CheckError(pgerr, "play")
	}

	updatedPlayCount, pgerr := models.CountPlays(trackId, userId, s.db)
	if pgerr != nil {
		return nil, internal.CheckError(pgerr, "play")
	}
	return &pb.CreatePlayResponse{
		UpdatedPlayCount: updatedPlayCount,
		UpdatedCredits: playRequest.UpdatedCredits,
	}, nil
}*/

func (s *Server) UpdateUser(ctx context.Context, user *pb.User) (*tagpb.Empty, error) {
	err := checkRequiredAttributes(user)

	if err != nil {
		return nil, err
	}

	u, err := getUserModel(user)
	if err != nil {
		return nil, err
	}

	u.UpdatedAt = time.Now()
	_, pgerr := s.db.Model(u).
		Column("updated_at", "username", "full_name", "email", "member", "newsletter_notification").
		WherePK().
		Returning("*").
		Update()
	twerr := internal.CheckError(pgerr, "user")
	if twerr != nil {
		return nil, twerr
	}
	return &tagpb.Empty{}, nil
}

func (s *Server) DeleteUser(ctx context.Context, user *pb.User) (*tagpb.Empty, error) {
	u, requiredErr := getUserModel(user)
	if requiredErr != nil {
		return nil, requiredErr
	}

	tx, err := s.db.Begin()
	if err != nil {
		return nil, internal.CheckError(err, "")
	}
	defer tx.Rollback()

	if pgerr, table := u.Delete(tx); pgerr != nil {
		return nil, internal.CheckError(pgerr, table)
	}

	err = tx.Commit()
	if err != nil {
		return nil, internal.CheckError(err, "")
	}

	return &tagpb.Empty{}, nil
}

func (s *Server) FollowGroup(ctx context.Context, userToUserGroup *pb.UserToUserGroup) (*tagpb.Empty, error) {
	userId, err := internal.GetUuidFromString(userToUserGroup.UserId)
	if err != nil {
		return nil, err
	}
	userGroupId, err := internal.GetUuidFromString(userToUserGroup.UserGroupId)
	if err != nil {
		return nil, err
	}

	u := &models.User{Id: userId}
	if pgerr, table := u.FollowGroup(s.db, userGroupId); pgerr != nil {
		return nil, internal.CheckError(pgerr, table)
	}

	return &tagpb.Empty{}, nil
}

func (s *Server) UnfollowGroup(ctx context.Context, userToUserGroup *pb.UserToUserGroup) (*tagpb.Empty, error) {
	userId, err := internal.GetUuidFromString(userToUserGroup.UserId)
	if err != nil {
		return nil, err
	}
	userGroupId, err := internal.GetUuidFromString(userToUserGroup.UserGroupId)
	if err != nil {
		return nil, err
	}

	u := &models.User{Id: userId}
	if pgerr, table := u.UnfollowGroup(s.db, userGroupId); pgerr != nil {
		return nil, internal.CheckError(pgerr, table)
	}
	return &tagpb.Empty{}, nil
}

func (s *Server) AddFavoriteTrack(ctx context.Context, userToTrack *pb.UserToTrack) (*tagpb.Empty, error) {
	userId, err := internal.GetUuidFromString(userToTrack.UserId)
	if err != nil {
		return nil, err
	}
	trackId, err := internal.GetUuidFromString(userToTrack.TrackId)
	if err != nil {
		return nil, err
	}

	u := &models.User{Id: userId}
	if pgerr, table := u.AddFavoriteTrack(s.db, trackId); pgerr != nil {
		return nil, internal.CheckError(pgerr, table)
	}

	return &tagpb.Empty{}, nil
}

func (s *Server) RemoveFavoriteTrack(ctx context.Context, userToTrack *pb.UserToTrack) (*tagpb.Empty, error) {
	userId, err := internal.GetUuidFromString(userToTrack.UserId)
	if err != nil {
		return nil, err
	}
	trackId, err := internal.GetUuidFromString(userToTrack.TrackId)
	if err != nil {
		return nil, err
	}

	u := &models.User{Id: userId}
	if pgerr, table := u.RemoveFavoriteTrack(s.db, trackId); pgerr != nil {
		return nil, internal.CheckError(pgerr, table)
	}
	return &tagpb.Empty{}, nil
}



func getUserModel(user *pb.User) (*models.User, twirp.Error) {
	id, err := internal.GetUuidFromString(user.Id)
	if err != nil {
		return nil, err
	}
	return &models.User{
		Id: id,
		Username: user.Username,
		// DisplayName: user.DisplayName,
		FullName: user.FullName,
		Email: user.Email,
		FirstName: user.FirstName,
		LastName: user.LastName,
		Member: user.Member,
		NewsletterNotification: user.NewsletterNotification,
	}, nil
}

func checkRequiredAttributes(user *pb.User) (twirp.Error) {
	if user.Email == ""	|| user.Username == "" || user.FullName == "" {
		var argument string
		switch {
		case user.Email == "":
			argument = "email"
		case user.Username == "":
			argument = "username"
		case user.FullName == "":
			argument = "full_name"
		}
		return twirp.RequiredArgumentError(argument)
	}
	return nil
}

func getUserGroupResponse(ownerOfGroup []models.UserGroup) ([]*tagpb.RelatedUserGroup) {
	groups := make([]*tagpb.RelatedUserGroup, len(ownerOfGroup))
	for i, group := range ownerOfGroup {
		groups[i] = &tagpb.RelatedUserGroup{Id: group.Id.String(), DisplayName: group.DisplayName, Avatar: group.Avatar}
	}
	return groups
}
