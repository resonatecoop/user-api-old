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
	createUserGroup := func(db *pg.DB, u *pb.UserGroup, ownerId uuid.UUID) (error, string, *models.UserGroup) {
		var table string
		tx, err := db.Begin()
		if err != nil {
			return err, table, nil
		}
		defer tx.Rollback()

		groupTaxonomy := new(models.GroupTaxonomy)
		pgerr := s.db.Model(groupTaxonomy).Where("type = ?", userGroup.Type).First()
		if pgerr != nil {
			return pgerr, "group_taxonomy", nil
		}

		var newAddress *models.StreetAddress
		if u.Address != nil {
			newAddress = &models.StreetAddress{Data: u.Address.Data}
			_, pgerr = s.db.Model(newAddress).Returning("*").Insert()
			if pgerr != nil {
				return pgerr, "street_address", nil
			}
		}

		links := make([]*models.Link, len(u.Links))
		linkIds := make([]uuid.UUID, len(u.Links))
		for i, link := range u.Links {
			links[i] = &models.Link{Platform: link.Platform, Uri: link.Uri}
			_, pgerr = s.db.Model(links[i]).Returning("*").Insert()
			if pgerr != nil {
				return pgerr, "link", nil
			}
			linkIds[i] = links[i].Id
		}


		tags := make([]*models.Tag, len(u.Tags))
		tagIds := make([]uuid.UUID, len(u.Tags))
		for i, tag := range u.Tags {
			tags[i] = &models.Tag{Type: tag.Type, Name: tag.Name}
			_, pgerr = s.db.Model(tags[i]).
				Where("type = ?", tags[i].Type).
				Where("lower(name) = lower(?)", tags[i].Name).
				Returning("*").
				SelectOrInsert()
			if pgerr != nil {
				return pgerr, "tag", nil
			}
			tagIds[i] = tags[i].Id
		}

		recommendedArtistIds, pgerr := getRelatedUserGroupIds(u.RecommendedArtists, s.db)
		if pgerr != nil {
			return pgerr, "user_group", nil
		}

		subGroupIds, pgerr := getRelatedUserGroupIds(u.SubGroups, s.db)
		if pgerr != nil {
			return pgerr, "user_group", nil
		}

		newUserGroup := &models.UserGroup{
			DisplayName: userGroup.DisplayName,
			Avatar: userGroup.Avatar,
			Banner: userGroup.Banner,
			ShortBio: userGroup.ShortBio,
			OwnerId: ownerId,
			TypeId: groupTaxonomy.Id,
			AddressId: newAddress.Id,
			Tags: tagIds,
			Links: linkIds,
			RecommendedArtists: recommendedArtistIds,
			SubGroups: subGroupIds,
			// Members
		}
		_, pgerr = s.db.Model(newUserGroup).Returning("*").Insert()
		if pgerr != nil {
			return pgerr, "user_group", nil
		}

		pgerr = db.Model(newUserGroup).
			Column("Privacy").
			WherePK().
			Select()
		if pgerr != nil {
			return pgerr, "user_group", nil
		}

		return tx.Commit(), table, newUserGroup
	}

	requiredErr := checkRequiredAttributes(userGroup)
	if requiredErr != nil {
		return nil, requiredErr
	}

	ownerId, err := internal.GetUuidFromString(userGroup.OwnerId)
	if err != nil {
		return nil, err
	}

	pgerr, table, newUserGroup := createUserGroup(s.db, userGroup, ownerId)
	if pgerr != nil {
		return nil, internal.CheckError(pgerr, table)
	}

	// Building response
	userGroup.Address.Id = newUserGroup.AddressId.String()
	privacy := &pb.Privacy{
		Id: newUserGroup.Privacy.Id.String(),
		Private: newUserGroup.Privacy.Private,
		OwnedTracks: newUserGroup.Privacy.OwnedTracks,
		SupportedArtists: newUserGroup.Privacy.SupportedArtists,
	}

	for i, linkId := range(newUserGroup.Links) {
		userGroup.Links[i].Id = linkId.String()
	}
	for i, tagId := range(newUserGroup.Tags) {
		userGroup.Tags[i].Id = tagId.String()
	}

  return &pb.UserGroup{
		Id: newUserGroup.Id.String(),
		DisplayName: newUserGroup.DisplayName,
		Avatar: newUserGroup.Avatar,
		Banner: newUserGroup.Banner,
		ShortBio: newUserGroup.ShortBio,
		OwnerId: newUserGroup.OwnerId.String(),
		Type: userGroup.Type,
		Address: userGroup.Address,
		Links: userGroup.Links,
		Tags: userGroup.Tags,
		RecommendedArtists: userGroup.RecommendedArtists,
		SubGroups: userGroup.SubGroups,
		Privacy: privacy,
	}, nil
}

// TODO handle privacy settings
func (s *Server) GetUserGroup(ctx context.Context, userGroup *pb.UserGroup) (*pb.UserGroup, error) {
	u, err := getUserGroupModel(userGroup)
	if err != nil {
		return nil, err
	}

	// pgerr := s.db.Select(u)
	pgerr := s.db.Model(u).
		Column("user_group.*", "Privacy", "Type", "Address").
		WherePK().
		Select()
	if pgerr != nil {
		return nil, internal.CheckError(pgerr, "user_group")
	}

	// Get user group links
	var groupLinks []models.Link
	pgerr = s.db.Model(&groupLinks).
		Where("id in (?)", pg.In(u.Links)).
		Select()
	if pgerr != nil {
		return nil, internal.CheckError(pgerr, "link")
	}
	links := make([]*pb.Link, len(groupLinks))
	for i, link := range groupLinks {
		links[i] = &pb.Link{Id: link.Id.String(), Platform: link.Platform, Uri: link.Uri}
	}

	// Get user group tags
	var groupTags []models.Tag
	pgerr = s.db.Model(&groupTags).
		Where("id in (?)", pg.In(u.Tags)).
		Select()
	if pgerr != nil {
		return nil, internal.CheckError(pgerr, "tag")
	}
	tags := make([]*pb.Tag, len(groupTags))
	for i, tag := range groupTags {
		tags[i] = &pb.Tag{Id: tag.Id.String(), Type: tag.Type, Name: tag.Name}
	}

	address := &userpb.StreetAddress{Id: u.Address.Id.String(), Data: u.Address.Data}
	privacy := &pb.Privacy{Id: u.Privacy.Id.String(), Private: u.Privacy.Private, SupportedArtists: u.Privacy.SupportedArtists, OwnedTracks: u.Privacy.OwnedTracks}

	return &pb.UserGroup{
		Id: u.Id.String(),
		DisplayName: u.DisplayName,
		Description: u.Description,
		ShortBio: u.ShortBio,
		Avatar: u.Avatar,
		Banner: u.Banner,
		GroupEmailAddress: u.GroupEmailAddress,
		OwnerId: u.OwnerId.String(),
		Type: u.Type.Type,
		Privacy: privacy,
		Address: address,
		Links: links,
		Tags: tags,
		// RecommendedArtists
		// Subgroups

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
	// Delete related Links, Tracks, TrackGroups
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
		ShortBio: userGroup.ShortBio,
		Avatar: userGroup.Avatar,
		Banner: userGroup.Banner,
		GroupEmailAddress: userGroup.GroupEmailAddress,
		// Links
		// Tags
		// Address
	}, nil
}

func getRelatedUserGroupIds(userGroups []*pb.UserGroup, db *pg.DB) ([]uuid.UUID, error) {
	relatedUserGroups := make([]*models.UserGroup, len(userGroups))
	relatedUserGroupIds := make([]uuid.UUID, len(userGroups))
	for i, userGroup := range userGroups {
		id, twerr := internal.GetUuidFromString(userGroup.Id)
		if twerr != nil {
			return nil, twerr.(error)
		}
		relatedUserGroups[i] = &models.UserGroup{Id: id}
		pgerr := db.Model(relatedUserGroups[i]).
			Returning("id", "display_name", "avatar").
			Select()
		if pgerr != nil {
			return nil, pgerr
		}
		userGroup.DisplayName = relatedUserGroups[i].DisplayName
		userGroup.Avatar = relatedUserGroups[i].Avatar
		relatedUserGroupIds[i] = relatedUserGroups[i].Id
	}
	return relatedUserGroupIds, nil
}

func getRelatedUserGroup() () {

}

func checkRequiredAttributes(userGroup *pb.UserGroup) (twirp.Error) {
	if userGroup.DisplayName == "" || len(userGroup.Avatar) == 0 || userGroup.Address.Data == nil || userGroup.Type == "" || userGroup.OwnerId == "" {
		var argument string
		switch {
		case userGroup.DisplayName == "":
			argument = "display_name"
		case len(userGroup.Avatar) == 0:
			argument = "avatar"
		case userGroup.Address.Data == nil:
			argument = "address"
		case userGroup.Type == "":
			argument = "type"
		case userGroup.OwnerId == "":
			argument = "owner"
		}
		return twirp.RequiredArgumentError(argument)
	}
	return nil
}
