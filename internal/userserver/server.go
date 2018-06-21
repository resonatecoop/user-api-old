package userserver

import (
	// "fmt"
	// "reflect"
	"time"
	"context"

	"github.com/go-pg/pg"
	"github.com/twitchtv/twirp"
	"github.com/satori/go.uuid"

	pb "user-api/rpc/user"
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

/* TODO add to response:
- tags
- residence_address
- member_of_groups
- playlists */
func (s *Server) GetUser(ctx context.Context, user *pb.User) (*pb.User, error) {
	u, err := getUserModel(user)
	if err != nil {
		return nil, err
	}

	pgerr := s.db.Select(u)
	twerr := internal.CheckError(pgerr, "user")
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
		FavoriteTracks: internal.ConvertUuidToStrArray(u.FavoriteTracks),
		FollowedGroups: internal.ConvertUuidToStrArray(u.FollowedGroups),
	}, nil
}

func (s *Server) CreateUser(ctx context.Context, user *pb.User) (*pb.User, error) {
	requiredErr := checkRequiredAttributes(user)
	if requiredErr != nil {
		return nil, requiredErr
	}

	newUser := &models.User{
		Username: user.Username,
		FullName: user.FullName,
		Email: user.Email,
		DisplayName: user.DisplayName,
	}
	_, err := s.db.Model(newUser).Returning("*").Insert()

	twerr := internal.CheckError(err, "user")
	if twerr != nil {
		return nil, twerr
	}

	return &pb.User{
		Id: newUser.Id.String(),
		Username: newUser.Username,
		DisplayName: newUser.DisplayName,
		FullName: newUser.FullName,
		Email: newUser.Email,
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
	twerr := internal.CheckError(pgerr, "user")
	if twerr != nil {
		return nil, twerr
	}
	return &pb.Empty{}, nil
}

func (s *Server) DeleteUser(ctx context.Context, user *pb.User) (*pb.Empty, error) {
	deleteUser := func(db *pg.DB, u *models.User) (error, string) {
		var table string
		tx, err := db.Begin()
		if err != nil {
			return err, table
		}
		defer tx.Rollback()

		user := new(models.User)
		pgerr := tx.Model(user).
	    Column("user.favorite_tracks").
	    Where("id = ?", u.Id).
	    Select()
		if pgerr != nil {
			return pgerr, "user"
		}

		if len(user.FavoriteTracks) > 0 {
			_, pgerr = tx.ExecOne(`
				UPDATE tracks
				SET favorite_of_users = array_remove(favorite_of_users, ?)
				WHERE id IN (?)
			`, u.Id, pg.In(user.FavoriteTracks))
				if pgerr != nil {
					return pgerr, "track"
				}
		}

		if len(user.FollowedGroups) > 0 {
			_, pgerr = tx.ExecOne(`
				UPDATE user_groups
				SET followers = array_remove(followers, ?)
				WHERE id IN (?)
			`, u.Id, pg.In(user.FollowedGroups))
				if pgerr != nil {
					return pgerr, "user_group"
				}
		}

		pgerr = s.db.Delete(u)
		if pgerr != nil {
			return pgerr, "user"
		}

		return tx.Commit(), table
	}

	u, requiredErr := getUserModel(user)
	if requiredErr != nil {
		return nil, requiredErr
	}

	if pgerr, table := deleteUser(s.db, u); pgerr != nil {
		return nil, internal.CheckError(pgerr, table)
	}

	return &pb.Empty{}, nil
}

func (s *Server) ConnectToUserGroup(ctx context.Context, userToUserGroup *pb.UserToUserGroup) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

// TODO: refacto, pretty similar to AddFavoriteTrack
func (s *Server) FollowGroup(ctx context.Context, userToUserGroup *pb.UserToUserGroup) (*pb.Empty, error) {
	followGroup := func(db *pg.DB, userId uuid.UUID, userGroupId uuid.UUID) (error, string) {
		var table string
		tx, err := db.Begin()
		if err != nil {
			return err, table
		}
		defer tx.Rollback()

		// Add userGroupId to user FollowedGroups
		userGroupIdArr := []uuid.UUID{userGroupId}
		_, pgerr := tx.ExecOne(`
			UPDATE users
			SET followed_groups = (select array_agg(distinct e) from unnest(followed_groups || ?) e)
			WHERE id = ?
		`, pg.Array(userGroupIdArr), userId)
		// WHERE NOT favorite_tracks @> ?
		if pgerr != nil {
			table = "user"
			return pgerr, table
		}

		// Add userId to userGroup Followers
		userIdArr := []uuid.UUID{userId}
		_, pgerr = tx.ExecOne(`
			UPDATE user_groups
			SET followers = (select array_agg(distinct e) from unnest(followers || ?) e)
			WHERE id = ?
		`, pg.Array(userIdArr), userGroupId)
		if pgerr != nil {
			table = "user_group"
			return pgerr, table
		}
		return tx.Commit(), table
	}
	userId, err := internal.GetUuidFromString(userToUserGroup.UserId)
	if err != nil {
		return nil, err
	}
	userGroupId, err := internal.GetUuidFromString(userToUserGroup.UserGroupId)
	if err != nil {
		return nil, err
	}

	if pgerr, table := followGroup(s.db, userId, userGroupId); pgerr != nil {
		return nil, internal.CheckError(pgerr, table)
	}

	return &pb.Empty{}, nil
}

func (s *Server) UnfollowGroup(ctx context.Context, userToUserGroup *pb.UserToUserGroup) (*pb.Empty, error) {
	unfollowGroup := func(db *pg.DB, userId uuid.UUID, userGroupId uuid.UUID) (error, string) {
		var table string
		tx, err := db.Begin()
		if err != nil {
			return err, table
		}
		// Rollback tx on error.
		defer tx.Rollback()

		// Remove userGroupId from user FollowedGroups
		_, pgerr := tx.ExecOne(`
			UPDATE users
			SET followed_groups = array_remove(followed_groups, ?)
			WHERE id = ?
		`, userGroupId, userId)
		if pgerr != nil {
			table = "user"
			return pgerr, table
		}

		// Remove userId from track FavoriteOfUsers
		_, pgerr = tx.ExecOne(`
			UPDATE user_groups
			SET followers = array_remove(followers, ?)
			WHERE id = ?
		`, userId, userGroupId)
		if pgerr != nil {
			table = "user_group"
			return pgerr, table
		}
		return tx.Commit(), table
	}

	userId, err := internal.GetUuidFromString(userToUserGroup.UserId)
	if err != nil {
		return nil, err
	}
	userGroupId, err := internal.GetUuidFromString(userToUserGroup.UserGroupId)
	if err != nil {
		return nil, err
	}

	if pgerr, table := unfollowGroup(s.db, userId, userGroupId); pgerr != nil {
		return nil, internal.CheckError(pgerr, table)
	}
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

	userId, err := internal.GetUuidFromString(userToTrack.UserId)
	if err != nil {
		return nil, err
	}
	trackId, err := internal.GetUuidFromString(userToTrack.TrackId)
	if err != nil {
		return nil, err
	}

	if pgerr, table := addFavoriteTrack(s.db, userId, trackId); pgerr != nil {
		return nil, internal.CheckError(pgerr, table)
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
	userId, err := internal.GetUuidFromString(userToTrack.UserId)
	if err != nil {
		return nil, err
	}
	trackId, err := internal.GetUuidFromString(userToTrack.TrackId)
	if err != nil {
		return nil, err
	}

	if pgerr, table := removeFavoriteTrack(s.db, userId, trackId); pgerr != nil {
		return nil, internal.CheckError(pgerr, table)
	}
	return &pb.Empty{}, nil
}



func getUserModel(user *pb.User) (*models.User, twirp.Error) {
	id, err := internal.GetUuidFromString(user.Id)
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
