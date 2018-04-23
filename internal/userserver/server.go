package userserver

import (
	"context"

	pb "toy-api/rpc/user"
)

// Server implements the ToyUser service
type Server struct {
	Users []*pb.User
}

// NewServer creates an instance of our server
func NewServer() *Server {
	return &Server{Users: []*pb.User{}}
}

func (s *Server) GetUsers(ctx context.Context, empty *pb.Empty) (*pb.Users, error) {
	return &pb.Users{Users: s.Users}, nil
}

func (s *Server) CreateUser(ctx context.Context, user *pb.User) (*pb.User, error) {
	s.Users = append(s.Users, user)
	return &pb.User{Email: user.Email, Name: user.Name}, nil
}
