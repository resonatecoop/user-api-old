package userserver

import (
	// "fmt"
	// "reflect"
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

func (s *Server) GetUser(ctx context.Context, user *pb.User) (*pb.User, error) {
	id, err := uuid.FromString(user.Id)
	if err != nil {
		return nil, twirp.InvalidArgumentError("id", "must be a valid uuid")
	}
	u := &models.User{Id: id}
	err = s.db.Select(u)
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, twirp.NotFoundError("user does not exist")
		}
		return nil, twirp.NewError("unknown", err.Error())
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
		return nil, twirp.RequiredArgumentError(argument)
	}

	newuser := &models.User{
		Username: user.Username,
		FullName: user.FullName,
		Email: user.Email,
		DisplayName: user.DisplayName,
	}
	_, err := s.db.Model(newuser).Returning("*").Insert()

	if err != nil {
		pgerr := err.(pg.Error)
		code := pgerr.Field('C')
		name := pgerr.Field('n')
		var message string
		if code == "23505" { // unique_violation // TODO put code in var
			message = strings.TrimPrefix(strings.TrimSuffix(name, "_key"), "users_")
			return nil, twirp.NewError("already_exists", message)
		} else {
			message = pgerr.Field('M')
			return nil, twirp.NewError("unknown", message)
		}
	}

	return &pb.User{
		Id: newuser.Id.String(),
		Username: newuser.Username,
		DisplayName: newuser.DisplayName,
		FullName: newuser.FullName,
		Email: newuser.Email,
	}, nil
}
