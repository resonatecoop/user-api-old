package usergroupserver

import (
	// "fmt"
	// "reflect"
	// "time"
	"context"

	"github.com/go-pg/pg"
	// "github.com/twitchtv/twirp"
	// "github.com/satori/go.uuid"

  userpb "user-api/rpc/user"
	pb "user-api/rpc/usergroup"
	// "user-api/internal"
	// "user-api/internal/database/models"
)

type Server struct {
	db *pg.DB
}

func NewServer(db *pg.DB) *Server {
	return &Server{db: db}
}

func (s *Server) CreateUserGroup(ctx context.Context, user *pb.UserGroup) (*pb.UserGroup, error) {
  return &pb.UserGroup{}, nil
}

func (s *Server) GetUserGroup(ctx context.Context, user *pb.UserGroup) (*pb.UserGroup, error) {
  return &pb.UserGroup{}, nil
}

func (s *Server) UpdateUserGroup(ctx context.Context, user *pb.UserGroup) (*userpb.Empty, error) {
  return &userpb.Empty{}, nil
}

func (s *Server) DeleteUserGroup(ctx context.Context, user *pb.UserGroup) (*userpb.Empty, error) {
  return &userpb.Empty{}, nil
}

func (s *Server) GetLabelUserGroups(ctx context.Context, empty *userpb.Empty) (*pb.GroupedUserGroups, error) {
  return &pb.GroupedUserGroups{}, nil
}

func (s *Server) GetUserGroupsByGenre(ctx context.Context, tag *pb.Tag) (*pb.GroupedUserGroups, error) {
  return &pb.GroupedUserGroups{}, nil
}
