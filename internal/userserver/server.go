package userserver

import (
	// "fmt"
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
	twerr := checkError(pgerr)
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

	twerr := checkError(err)
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
	twerr := checkError(pgerr)
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
	twerr := checkError(pgerr)
	if twerr != nil {
		return nil, twerr
	}
	return &pb.Empty{}, nil
}

func (s *Server) ConnectToUserGroup(ctx context.Context, user *pb.UserToUserGroup) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

func (s *Server) FollowArtist(ctx context.Context, user *pb.UserToUserGroup) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

func (s *Server) UnfollowArtist(ctx context.Context, user *pb.UserToUserGroup) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

func (s *Server) AddFavoriteTrack(ctx context.Context, user_to *pb.UserToTrack) (*pb.Empty, error) {
	// user_id, err := getUuidFromString()


	return &pb.Empty{}, nil
}

func (s *Server) RemoveFavoriteTrack(ctx context.Context, user *pb.UserToTrack) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

func getUuidFromString(id string) (uuid.UUID, twirp.Error) {
	uid, err := uuid.FromString(id)
	if err != nil {
		return nil, twirp.InvalidArgumentError("id", "must be a valid uuid")
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

func checkError(err error) (twirp.Error) {
	if err != nil {
		if err == pg.ErrNoRows {
			return twirp.NotFoundError("user does not exist")
		}
		pgerr, ok := err.(pg.Error)
		if ok {
			code := pgerr.Field('C')
			name := pgerr.Field('n')
			var message string
			if code == "23505" { // unique_violation
				message = strings.TrimPrefix(strings.TrimSuffix(name, "_key"), "users_")
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
