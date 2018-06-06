package userserver

import (
	// "fmt"
	"context"
	"strings"
	"github.com/go-pg/pg"
	"github.com/twitchtv/twirp"
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

func (s *Server) GetUsers(ctx context.Context, empty *pb.Empty) (*pb.Users, error) {
	// q := models.NewUserQuery()
	//
	// users, err := s.Store.FindAll(q)
	// if err != nil {
	// 	return nil, err
	// }
	u := make([]*pb.User, 3)
	// for i := range u {
	// 	u[i] = &pb.User{Id: users[i].ID.String(), Email: users[i].Email, Username: users[i].Username, Address: users[i].Address}
	// }
	return &pb.Users{User: u}, nil
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
