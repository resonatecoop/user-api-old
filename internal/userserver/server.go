package userserver

import (
	"fmt"
	// "reflect"
	"time"
	"context"
	"strings"
	"github.com/go-pg/pg"
	"github.com/twitchtv/twirp"
	"github.com/satori/go.uuid"

	pb "user-api/rpc/user"
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

/* TODO add to response:
- tags
- residence_address
- member_of_groups
- followed_artists
- favorite_tracks
- playlists */
func (s *Server) GetUser(ctx context.Context, user *pb.User) (*pb.User, error) {
	// id, err := uuid.FromString(user.Id)
	// if err != nil {
	// 	return nil, twirp.InvalidArgumentError("id", "must be a valid uuid")
	// }
	// u := &models.User{Id: id}
	u, err := getUserModel(user)
	if err != nil {
		return nil, err
	}

	pgerr := s.db.Select(u)
	twerr := checkError(pgerr, "user")
	if twerr != nil {
		return nil, twerr
	}

	return &pb.User{
		Id: u.Id.String(),
		Username: u.Username,
		DisplayName: u.DisplayName,
		FullName: u.FullName,
		Email: u.Email,
		FirstName: u.FirstName,
		LastName: u.LastName,
		Member: u.Member,
		Avatar: u.Avatar,
		NewsletterNotification: u.NewsletterNotification,
	}, nil
}

func (s *Server) CreateUser(ctx context.Context, user *pb.User) (*pb.User, error) {
	requiredErr := checkRequiredAttributes(user)
	if requiredErr != nil {
		return nil, requiredErr
	}

	newuser := &models.User{
		Username: user.Username,
		FullName: user.FullName,
		Email: user.Email,
		DisplayName: user.DisplayName,
	}
	_, err := s.db.Model(newuser).Returning("*").Insert()

	twerr := checkError(err, "user")
	if twerr != nil {
		return nil, twerr
	}

	return &pb.User{
		Id: newuser.Id.String(),
		Username: newuser.Username,
		DisplayName: newuser.DisplayName,
		FullName: newuser.FullName,
		Email: newuser.Email,
	}, nil
}

func (s *Server) UpdateUser(ctx context.Context, user *pb.User) (*pb.Empty, error) {
	err := checkRequiredAttributes(user)

	if err != nil {
		return nil, err
	}

	u, err := getUserModel(user)
	if err != nil {
		return nil, err
	}

	u.UpdatedAt = time.Now()
	_, pgerr := s.db.Model(u).WherePK().Returning("*").UpdateNotNull()
	twerr := checkError(pgerr, "user")
	if twerr != nil {
		return nil, twerr
	}
	return &pb.Empty{}, nil
}

func (s *Server) DeleteUser(ctx context.Context, user *pb.User) (*pb.Empty, error) {
	u, requiredErr := getUserModel(user)
	if requiredErr != nil {
		return nil, requiredErr
	}

	pgerr := s.db.Delete(u)
	twerr := checkError(pgerr, "user")
	if twerr != nil {
		return nil, twerr
	}
	return &pb.Empty{}, nil
}

func (s *Server) ConnectToUserGroup(ctx context.Context, userToUserGroup *pb.UserToUserGroup) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

func (s *Server) FollowArtist(ctx context.Context, userToUserGroup *pb.UserToUserGroup) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

func (s *Server) UnfollowArtist(ctx context.Context, userToUserGroup *pb.UserToUserGroup) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

func (s *Server) AddFavoriteTrack(ctx context.Context, userToTrack *pb.UserToTrack) (*pb.Empty, error) {
	addFavoriteTrack := func(db *pg.DB, userId uuid.UUID, trackId uuid.UUID) (error, string) {
		var table string
	  tx, err := db.Begin()
	  if err != nil {
			return err, table
	  }
	  // Rollback tx on error.
	  defer tx.Rollback()

		// Add trackId to user FavoriteTracks
		trackIdArr := []uuid.UUID{trackId}
		_, pgerr := tx.ExecOne(`
			UPDATE users
			SET favorite_tracks = (select array_agg(distinct e) from unnest(favorite_tracks || ?) e)
			WHERE id = ?
		`, pg.Array(trackIdArr), userId)
		// WHERE NOT favorite_tracks @> ?
		if pgerr != nil {
			table = "user"
			return pgerr, table
		}

		// Add userId to track FavoriteOfUsers
		userIdArr := []uuid.UUID{userId}
		_, pgerr = tx.ExecOne(`
			UPDATE tracks
			SET favorite_of_users = (select array_agg(distinct e) from unnest(favorite_of_users || ?) e)
			WHERE id = ?
		`, pg.Array(userIdArr), trackId)
		if pgerr != nil {
			table = "track"
			return pgerr, table
		}
	  return tx.Commit(), table
	}

	userId, err := getUuidFromString(userToTrack.UserId)
	if err != nil {
		return nil, err
	}
	trackId, err := getUuidFromString(userToTrack.TrackId)
	if err != nil {
		return nil, err
	}

	if pgerr, table := addFavoriteTrack(s.db, userId, trackId); pgerr != nil {
		return nil, checkError(pgerr, table)
	}

	return &pb.Empty{}, nil
}

func (s *Server) RemoveFavoriteTrack(ctx context.Context, userToTrack *pb.UserToTrack) (*pb.Empty, error) {
	removeFavoriteTrack := func(db *pg.DB, userId uuid.UUID, trackId uuid.UUID) (error, string) {
		var table string
	  tx, err := db.Begin()
	  if err != nil {
			return err, table
	  }
	  // Rollback tx on error.
	  defer tx.Rollback()

		// Remove trackId from user FavoriteTracks
		_, pgerr := tx.ExecOne(`
			UPDATE users
			SET favorite_tracks = array_remove(favorite_tracks, ?)
			WHERE id = ?
		`, trackId, userId)
		if pgerr != nil {
			table = "user"
			return pgerr, table
		}

		// Remove userId from track FavoriteOfUsers
		_, pgerr = tx.ExecOne(`
			UPDATE tracks
			SET favorite_of_users = array_remove(favorite_of_users, ?)
			WHERE id = ?
		`, userId, trackId)
		if pgerr != nil {
			table = "track"
			return pgerr, table
		}
	  return tx.Commit(), table
	}

	// TODO refacto
	userId, err := getUuidFromString(userToTrack.UserId)
	if err != nil {
		return nil, err
	}
	trackId, err := getUuidFromString(userToTrack.TrackId)
	if err != nil {
		return nil, err
	}

	if pgerr, table := removeFavoriteTrack(s.db, userId, trackId); pgerr != nil {
		return nil, checkError(pgerr, table)
	}
	return &pb.Empty{}, nil
}

func getUuidFromString(id string) (uuid.UUID, twirp.Error) {
	uid, err := uuid.FromString(id)
	if err != nil {
		return uuid.UUID{}, twirp.InvalidArgumentError("id", "must be a valid uuid")
	}
	return uid, nil
}

func getUserModel(user *pb.User) (*models.User, twirp.Error) {
	id, err := getUuidFromString(user.Id)
	if err != nil {
		return nil, err
	}
	return &models.User{
		Id: id,
		Username: user.Username,
		DisplayName: user.DisplayName,
		FullName: user.FullName,
		Email: user.Email,
		FirstName: user.FirstName,
		LastName: user.LastName,
		Member: user.Member,
		Avatar: user.Avatar,
		NewsletterNotification: user.NewsletterNotification,
	}, nil
}

func checkRequiredAttributes(user *pb.User) (twirp.Error) {
	if user.Username == "" || user.FullName == "" || user.Email == "" || user.DisplayName == "" {
		var argument string
		switch {
		case user.Username == "":
			argument = "username"
		case user.Email == "":
			argument = "email"
		case user.DisplayName == "":
			argument = "display_name"
		case user.FullName == "":
			argument = "full_name"
		}
		return twirp.RequiredArgumentError(argument)
	}
	return nil
}

// TODO move to utils package
func checkError(err error, table string) (twirp.Error) {
	if err != nil {
		if err == pg.ErrNoRows {
			return twirp.NotFoundError(fmt.Sprintf("%s does not exist", table))
		}
		pgerr, ok := err.(pg.Error)
		if ok {
			code := pgerr.Field('C')
			name := pgerr.Field('n')
			var message string
			if code == "23505" { // unique_violation
				message = strings.TrimPrefix(strings.TrimSuffix(name, "_key"), fmt.Sprintf("%ss_", table))
				return twirp.NewError("already_exists", message)
			} else {
				message = pgerr.Field('M')
				return twirp.NewError("unknown", message)
			}
		}
		return twirp.NewError("unknown", err.Error())
	}
	return nil
}
