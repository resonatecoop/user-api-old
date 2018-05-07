package userserver

import (
	// "fmt"
	"context"
	"database/sql"
	"github.com/satori/go.uuid"
	pb "user-api/rpc/user"
	models "user-api/internal/database/models"
	"gopkg.in/src-d/go-kallax.v1"
)

// Server implements the ToyUser service
type Server struct {
	// DB *sql.DB
	Store *models.UserStore
}

// NewServer creates an instance of our server
func NewServer(db *sql.DB) *Server {
	store := models.NewUserStore(db)
	return &Server{Store: store}
}

func (s *Server) GetUsers(ctx context.Context, empty *pb.Empty) (*pb.Users, error) {
	q := models.NewUserQuery()

	users, err := s.Store.FindAll(q)
	if err != nil {
		return nil, err
	}
	u := make([]*pb.User, len(users))
	for i := range u {
		u[i] = &pb.User{Id: users[i].ID.String(), Email: users[i].Email, Username: users[i].Username, Address: users[i].Address}
	}
	return &pb.Users{User: u}, nil
}

func (s *Server) CreateUser(ctx context.Context, user *pb.User) (*pb.User, error) {
	u := uuid.NewV4()
	newuser, _ := models.NewUser(kallax.UUID(u), user.Username, user.Email, user.Address)
	err := s.Store.Insert(newuser)
	if err != nil {
    return nil, err
	}
	return &pb.User{Id: newuser.ID.String(), Email: user.Email, Username: user.Username, Address: user.Address}, nil
}
