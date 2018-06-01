package userserver

import (
	// "fmt"
	"context"
	"github.com/satori/go.uuid"
	"github.com/go-pg/pg"
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
	id := uuid.NewV4()
	newuser := &models.User{
		Id: id,
		Username: user.Username,
		DisplayName: user.DisplayName,
		FullName: user.FullName,
		Email: user.Email,
	}
	_, err := s.db.Model(newuser).Returning("*").Insert()
	if err != nil {
    return nil, err
	}
	// fmt.Printf("%+v\n", newuser)
	return &pb.User{
		Id: newuser.Id.String(),
		Username: newuser.Username,
		DisplayName: newuser.DisplayName,
		FullName: newuser.FullName,
		Email: newuser.Email,
		FirstName: newuser.FirstName,
		LastName: newuser.LastName,
		Member: newuser.Member,
		Avatar: newuser.Avatar,
		NewsletterNotification: newuser.NewsletterNotification,
	}, nil
}
