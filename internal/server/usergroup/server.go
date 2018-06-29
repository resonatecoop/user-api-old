package usergroupserver

import (
	// "fmt"
	// "reflect"
	"time"
	"context"
	"net/url"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/twitchtv/twirp"
	"github.com/satori/go.uuid"

  userpb "user-api/rpc/user"
	pb "user-api/rpc/usergroup"
	"user-api/internal"
	"user-api/internal/database/models"
)

type Server struct {
	db *pg.DB
}

func NewServer(db *pg.DB) *Server {
	return &Server{db: db}
}

func (s *Server) CreateUserGroup(ctx context.Context, userGroup *pb.UserGroup) (*pb.UserGroup, error) {
	requiredErr := checkRequiredAttributes(userGroup)
	if requiredErr != nil {
		return nil, requiredErr
	}

	ownerId, err := internal.GetUuidFromString(userGroup.OwnerId)
	if err != nil {
		return nil, err
	}

	typeId, err := internal.GetUuidFromString(userGroup.TypeId)
	if err != nil {
		return nil, err
	}

	adminUsers := []uuid.UUID{ownerId}
	newUserGroup := &models.UserGroup{
		DisplayName: userGroup.DisplayName,
		Avatar: userGroup.Avatar,
		OwnerId: ownerId,
		TypeId: typeId,
		AdminUsers: adminUsers,
		// TODO create and add Address
	}
	_, pgerr := s.db.Model(newUserGroup).Returning("*").Insert()

	twerr := internal.CheckError(pgerr, "user_group")
	if twerr != nil {
		return nil, twerr
	}

  return &pb.UserGroup{
		Id: newUserGroup.Id.String(),
		DisplayName: newUserGroup.DisplayName,
		Avatar: newUserGroup.Avatar,
		OwnerId: newUserGroup.OwnerId.String(),
		TypeId: newUserGroup.TypeId.String(),
		AdminUsers: internal.ConvertUuidToStrArray(newUserGroup.AdminUsers),
		// Address: newUserGroup.Address,
	}, nil
}

func (s *Server) GetUserGroup(ctx context.Context, userGroup *pb.UserGroup) (*pb.UserGroup, error) {
	u, err := getUserGroupModel(userGroup)
	if err != nil {
		return nil, err
	}

	pgerr := s.db.Select(u)
	twerr := internal.CheckError(pgerr, "user_group")
	if twerr != nil {
		return nil, twerr
	}

	return &pb.UserGroup{
		Id: u.Id.String(),
		DisplayName: u.DisplayName,
		Description: u.Description,
		Avatar: u.Avatar,
		Banner: u.Banner,
		OwnerId: u.OwnerId.String(),
		TypeId: u.TypeId.String(),
		GroupEmailAddress: u.GroupEmailAddress,
		// Address
		// Links
		// Tags
		// Tracks
		// TrackGroups
		// Members
	}, nil
}

// TODO update address, links, tags
func (s *Server) UpdateUserGroup(ctx context.Context, userGroup *pb.UserGroup) (*userpb.Empty, error) {
	err := checkRequiredAttributes(userGroup)

	if err != nil {
		return nil, err
	}

	u, err := getUserGroupModel(userGroup)
	if err != nil {
		return nil, err
	}

	u.UpdatedAt = time.Now()
	_, pgerr := s.db.Model(u).WherePK().Returning("*").UpdateNotNull()
	twerr := internal.CheckError(pgerr, "user_group")
	if twerr != nil {
		return nil, twerr
	}

  return &userpb.Empty{}, nil
}

func (s *Server) DeleteUserGroup(ctx context.Context, userGroup *pb.UserGroup) (*userpb.Empty, error) {
	// Delete related Links, Tracks, TrackGroups, Tags (?)
	// member_of_groups
	deleteUserGroup := func(db *pg.DB, u *models.UserGroup) (error, string) {
		var table string
		tx, err := db.Begin()
		if err != nil {
			return err, table
		}
		defer tx.Rollback()

		userGroup := new(models.UserGroup)
		pgerr := tx.Model(userGroup).
			Column("user_group.followers"). // TODO add other columns
			Where("id = ?", u.Id).
			Select()
		if pgerr != nil {
			return pgerr, "user_group"
		}

		if len(userGroup.Followers) > 0 {
			_, pgerr = tx.ExecOne(`
				UPDATE users
				SET followed_groups = array_remove(followed_groups, ?)
				WHERE id IN (?)
			`, u.Id, pg.In(userGroup.Followers))
			if pgerr != nil {
				return pgerr, "user"
			}
		}

		pgerr = s.db.Delete(u)
		if pgerr != nil {
			return pgerr, "user_group"
		}

		return tx.Commit(), table
	}
	u, requiredErr := getUserGroupModel(userGroup)
	if requiredErr != nil {
		return nil, requiredErr
	}

	if pgerr, table := deleteUserGroup(s.db, u); pgerr != nil {
		return nil, internal.CheckError(pgerr, table)
	}

	return &userpb.Empty{}, nil
}

func (s *Server) GetLabelUserGroups(ctx context.Context, empty *userpb.Empty) (*pb.GroupedUserGroups, error) {
	var labels []models.UserGroup

	// err := s.db.Model(&labels).
	// 	Join("LEFT JOIN group_taxonomies AS g").
	// 	JoinOn("g.id = user_group.type_id").
	// 	Where("g.type = ?", "label").
	//   Apply(orm.Pagination(ctx.Value("query").(url.Values))).
	//   Select()

	groupTaxomony := new(models.GroupTaxonomy)
	err := s.db.Model(groupTaxomony).
		Where("type = ?", "label").
		First()
	twerr := internal.CheckError(err, "group_taxonomy")
	if twerr != nil {
		return nil, twerr
	}

	err = s.db.Model(&labels).
		Where("user_group.type_id = ?", groupTaxomony.Id).
		Apply(orm.Pagination(ctx.Value("query").(url.Values))).
		Select()
	twerr = internal.CheckError(err, "user_group")
	if twerr != nil {
		return nil, twerr
	}

	userGroups := make([]*pb.UserGroup, len(labels))
	for i := range userGroups {
		userGroups[i] = &pb.UserGroup{Id: labels[i].Id.String(), DisplayName: labels[i].DisplayName, Avatar: labels[i].Avatar}
	}

  return &pb.GroupedUserGroups{Labels: userGroups}, nil
}

func (s *Server) GetUserGroupsByGenre(ctx context.Context, tag *pb.Tag) (*pb.GroupedUserGroups, error) {
  return &pb.GroupedUserGroups{}, nil
}

func getUserGroupModel(userGroup *pb.UserGroup) (*models.UserGroup, twirp.Error) {
	id, err := internal.GetUuidFromString(userGroup.Id)
	if err != nil {
		return nil, err
	}
	return &models.UserGroup{
		Id: id,
		DisplayName: userGroup.DisplayName,
		Description: userGroup.Description,
		Avatar: userGroup.Avatar,
		Banner: userGroup.Banner,
		GroupEmailAddress: userGroup.GroupEmailAddress,
		// Links
		// Tags
		// Address
	}, nil
}

func checkRequiredAttributes(userGroup *pb.UserGroup) (twirp.Error) {
	// +admins?
	if userGroup.DisplayName == "" || len(userGroup.Avatar) == 0 || userGroup.Address == "" || userGroup.TypeId == "" || userGroup.OwnerId == "" {
		var argument string
		switch {
		case userGroup.DisplayName == "":
			argument = "display_name"
		case len(userGroup.Avatar) == 0:
			argument = "avatar"
		case userGroup.Address == "":
			argument = "address"
		case userGroup.TypeId == "":
			argument = "type"
		case userGroup.OwnerId == "":
			argument = "owner"
		}
		return twirp.RequiredArgumentError(argument)
	}
	return nil
}
