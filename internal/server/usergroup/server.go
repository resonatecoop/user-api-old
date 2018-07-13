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
	createUserGroup := func(db *pg.DB, userGroup *pb.UserGroup, ownerId uuid.UUID) (error, string, *models.UserGroup) {
		var table string
		tx, err := db.Begin()
		if err != nil {
			return err, table, nil
		}
		defer tx.Rollback()

		groupTaxonomy := new(models.GroupTaxonomy)
		pgerr := s.db.Model(groupTaxonomy).Where("type = ?", userGroup.Type.Type).First()
		if pgerr != nil {
			return pgerr, "group_taxonomy", nil
		}

		var newAddress *models.StreetAddress
		if userGroup.Address != nil {
			newAddress = &models.StreetAddress{Data: userGroup.Address.Data}
			_, pgerr = s.db.Model(newAddress).Returning("*").Insert()
			if pgerr != nil {
				return pgerr, "street_address", nil
			}
		}

		linkIds, pgerr := getLinkIds(userGroup, s.db)
		if pgerr != nil {
			return pgerr, "link", nil
		}

		tagIds, pgerr := getTagIds(userGroup, s.db)
		if pgerr != nil {
			return pgerr, "tag", nil
		}

		recommendedArtistIds, pgerr := getRelatedUserGroupIds(userGroup.RecommendedArtists, s.db)
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
			// SubGroups: subGroupIds,
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

		// Building response
		userGroup.Address.Id = newUserGroup.AddressId.String()
		userGroup.Type.Id = newUserGroup.TypeId.String()
		userGroup.Privacy = &pb.Privacy{
			Id: newUserGroup.Privacy.Id.String(),
			Private: newUserGroup.Privacy.Private,
			OwnedTracks: newUserGroup.Privacy.OwnedTracks,
			SupportedArtists: newUserGroup.Privacy.SupportedArtists,
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
		Privacy: userGroup.Privacy,
	}, nil
}

// TODO handle privacy settings
func (s *Server) GetUserGroup(ctx context.Context, userGroup *pb.UserGroup) (*pb.UserGroup, error) {
	id, err := internal.GetUuidFromString(userGroup.Id)
	if err != nil {
		return nil, err
	}
	u := &models.UserGroup{Id: id}

	pgerr := s.db.Model(u).
		Column("user_group.*", "Privacy", "Type", "Address").
		WherePK().
		Select()
	if pgerr != nil {
		return nil, internal.CheckError(pgerr, "user_group")
	}

	// Get user group links
	links := make([]*pb.Link, len(u.Links))
	if len(links) > 0 {
		var groupLinks []models.Link
		pgerr = s.db.Model(&groupLinks).
			Where("id in (?)", pg.In(u.Links)).
			Select()
		if pgerr != nil {
			return nil, internal.CheckError(pgerr, "link")
		}
		for i, link := range groupLinks {
			links[i] = &pb.Link{Id: link.Id.String(), Platform: link.Platform, Uri: link.Uri}
		}
	}

	// Get user group tags
	// TODO refacto using interface
	tags := make([]*pb.Tag, len(u.Tags))
	if len(tags) > 0 {
		var groupTags []models.Tag
		pgerr = s.db.Model(&groupTags).
			Where("id in (?)", pg.In(u.Tags)).
			Select()
		if pgerr != nil {
			return nil, internal.CheckError(pgerr, "tag")
		}
		for i, tag := range groupTags {
			tags[i] = &pb.Tag{Id: tag.Id.String(), Type: tag.Type, Name: tag.Name}
		}
	}

	// Get related user groups
	recommendedArtists, twerr := getRelatedUserGroups(u.RecommendedArtists, s.db)
	if twerr != nil {
		return nil, twerr
	}
	sub_groups, twerr := getRelatedUserGroups(u.SubGroups, s.db)
	if twerr != nil {
		return nil, twerr
	}
	labels, twerr := getRelatedUserGroups(u.Labels, s.db)
	if twerr != nil {
		return nil, twerr
	}

	address := &userpb.StreetAddress{Id: u.Address.Id.String(), Data: u.Address.Data}
	privacy := &pb.Privacy{Id: u.Privacy.Id.String(), Private: u.Privacy.Private, SupportedArtists: u.Privacy.SupportedArtists, OwnedTracks: u.Privacy.OwnedTracks}
	groupType:= &pb.GroupTaxonomy{Id: u.Type.Id.String(), Type: u.Type.Type}

	return &pb.UserGroup{
		Id: u.Id.String(),
		DisplayName: u.DisplayName,
		Description: u.Description,
		ShortBio: u.ShortBio,
		Avatar: u.Avatar,
		Banner: u.Banner,
		GroupEmailAddress: u.GroupEmailAddress,
		OwnerId: u.OwnerId.String(),
		Type: groupType,
		Privacy: privacy,
		Address: address,
		Links: links,
		Tags: tags,
		RecommendedArtists: recommendedArtists,
		SubGroups: sub_groups,
		Labels: labels,

		// Tracks
		// TrackGroups
		// Members
	}, nil
}

func (s *Server) UpdateUserGroup(ctx context.Context, userGroup *pb.UserGroup) (*userpb.Empty, error) {
	updateUserGroup := func(db *pg.DB, userGroup *pb.UserGroup, u *models.UserGroup) (twirp.Error) {
		tx, err := db.Begin()
		if err != nil {
			return internal.CheckError(err, "user_group")
		}
		defer tx.Rollback()

		// Update address
		addressId, twerr := internal.GetUuidFromString(userGroup.Address.Id)
		if twerr != nil {
			return twerr
		}
		address := &models.StreetAddress{Id: addressId, Data: userGroup.Address.Data}
		_, pgerr := db.Model(address).Column("data").WherePK().Update()
		// _, pgerr := db.Model(address).Set("data = ?", pg.Hstore(userGroup.Address.Data)).Where("id = ?id").Update()
		if pgerr != nil {
			return internal.CheckError(pgerr, "street_address")
		}

		// Update privacy
		privacyId, twerr := internal.GetUuidFromString(userGroup.Privacy.Id)
		if twerr != nil {
			return twerr
		}
		privacy := &models.UserGroupPrivacy{Id: privacyId, Private: userGroup.Privacy.Private, OwnedTracks: userGroup.Privacy.OwnedTracks, SupportedArtists: userGroup.Privacy.SupportedArtists}
		_, pgerr = s.db.Model(privacy).WherePK().Returning("*").UpdateNotNull()

		// Update tags
		tagIds, pgerr := getTagIds(userGroup, s.db)
		if pgerr != nil {
			return internal.CheckError(pgerr, "tag")
		}

		// Update links
		linkIds, pgerr := getLinkIds(userGroup, s.db)
		if pgerr != nil {
			return internal.CheckError(pgerr, "link")
		}
		// Delete links if needed
		pgerr = db.Model(u).WherePK().Column("links").Select()
		linkIdsToDelete := internal.Difference(u.Links, linkIds)
		_, pgerr = db.Model((*models.Link)(nil)).
    	Where("id in (?)", pg.In(linkIdsToDelete)).
			Delete()

		// Update recommended artists
		recommendedArtistIds, pgerr := getRelatedUserGroupIds(userGroup.RecommendedArtists, s.db)
		if pgerr != nil {
			return internal.CheckError(pgerr, "user_group")
		}

		// Update user group
		u.Tags = tagIds
		u.Links = linkIds
		u.RecommendedArtists = recommendedArtistIds
		u.UpdatedAt = time.Now()
		// _, pgerr := db.Model(u).Where("id = ?", u.Id).Returning("*").Update()
		_, pgerr = db.Model(u).WherePK().Returning("*").Update()

		if pgerr != nil {
			return internal.CheckError(pgerr, "user_group")
		}
		return nil
	}

	err := checkRequiredAttributes(userGroup)
	if err != nil {
		return nil, err
	}

	u, err := getUserGroupModel(userGroup)
	if err != nil {
		return nil, err
	}

	if twerr := updateUserGroup(s.db, userGroup, u); twerr != nil {
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
	id, err := internal.GetUuidFromString(userGroup.Id)
	if err != nil {
		return nil, err
	}
	u := &models.UserGroup{Id: id}

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

	groupTaxonomy := new(models.GroupTaxonomy)
	err := s.db.Model(groupTaxonomy).
		Where("type = ?", "label").
		First()
	twerr := internal.CheckError(err, "group_taxonomy")
	if twerr != nil {
		return nil, twerr
	}

	err = s.db.Model(&labels).
		Where("user_group.type_id = ?", groupTaxonomy.Id).
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

func (s *Server) GetUserGroupTypes(ctx context.Context, empty *userpb.Empty) (*pb.GroupTaxonomies, error) {
	var types []models.GroupTaxonomy
	var groupTaxonomies pb.GroupTaxonomies
	err := s.db.Model(&types).
		Where("group_taxonomy.type != ?", "distributor"). // except distributors, internally added by staff
		Select()
	if err != nil {
		return nil, internal.CheckError(err, "group_taxonomy")
	}

	for _, groupType := range(types) {
		groupTaxonomies.Types = append(groupTaxonomies.Types, &pb.GroupTaxonomy{Id: groupType.Id.String(), Type: groupType.Type, Name: groupType.Name})
	}

  return &groupTaxonomies, nil
}




func getUserGroupModel(userGroup *pb.UserGroup) (*models.UserGroup, twirp.Error) {
	id, err := internal.GetUuidFromString(userGroup.Id)
	if err != nil {
		return nil, err
	}
	addressId, err := internal.GetUuidFromString(userGroup.Address.Id)
	if err != nil {
		return nil, err
	}
	privacyId, err := internal.GetUuidFromString(userGroup.Privacy.Id)
	if err != nil {
		return nil, err
	}
	typeId, err := internal.GetUuidFromString(userGroup.Type.Id)
	if err != nil {
		return nil, err
	}
	ownerId, err := internal.GetUuidFromString(userGroup.OwnerId)
	if err != nil {
		return nil, err
	}
	return &models.UserGroup{
		Id: id,
		DisplayName: userGroup.DisplayName,
		Description: userGroup.Description,
		AddressId: addressId,
		PrivacyId: privacyId,
		TypeId: typeId,
		OwnerId: ownerId,
		ShortBio: userGroup.ShortBio,
		Avatar: userGroup.Avatar,
		Banner: userGroup.Banner,
		GroupEmailAddress: userGroup.GroupEmailAddress,
	}, nil
}

// Select user groups in db with given ids in 'userGroups'
// Return ids slice
// Used in CreateUserGroup/UpdateUserGroup to add/update ids slice to recommended Artists
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
			WherePK().
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

// Select user groups in db with given 'ids'
// Return slice of UserGroup response
// Used in GetUserGroup to respond with info about related user groups: recommended_artists, sub_groups, labels
func getRelatedUserGroups(ids []uuid.UUID, db *pg.DB) ([]*pb.UserGroup, twirp.Error) {
	groupsResponse := make([]*pb.UserGroup, len(ids))

	if len(ids) > 0 {
		var groups []models.UserGroup
		pgerr := db.Model(&groups).
			Where("id in (?)", pg.In(ids)).
			Select()
		if pgerr != nil {
			return nil, internal.CheckError(pgerr, "user_group")
		}
		for i, group := range groups {
			groupsResponse[i] = &pb.UserGroup{Id: group.Id.String(), DisplayName: group.DisplayName, Avatar: group.Avatar}
		}
	}

	return groupsResponse, nil
}

func getTagIds(userGroup *pb.UserGroup, db *pg.DB) ([]uuid.UUID, error) {
	tags := make([]*models.Tag, len(userGroup.Tags))
	tagIds := make([]uuid.UUID, len(userGroup.Tags))
	for i, tag := range(userGroup.Tags) {
		if tag.Id == "" { // new tag to create and add
			tags[i] = &models.Tag{Type: tag.Type, Name: tag.Name}
			_, pgerr := db.Model(tags[i]).
				Where("type = ?", tags[i].Type).
				Where("lower(name) = lower(?)", tags[i].Name).
				Returning("*").
				SelectOrInsert()
			if pgerr != nil {
				return nil, pgerr
			}
			tagIds[i] = tags[i].Id
			tag.Id = tags[i].Id.String()
		} else {
			tagId, twerr := internal.GetUuidFromString(tag.Id)
			if twerr != nil {
				return nil, twerr.(error)
			}
			tagIds[i] = tagId
		}
	}
	return tagIds, nil
}

func getLinkIds(userGroup *pb.UserGroup, db *pg.DB) ([]uuid.UUID, error) {
	links := make([]*models.Link, len(userGroup.Links))
	linkIds := make([]uuid.UUID, len(userGroup.Links))
	for i, link := range userGroup.Links {
		if link.Id == "" {
			links[i] = &models.Link{Platform: link.Platform, Uri: link.Uri}
			_, pgerr := db.Model(links[i]).Returning("*").Insert()
			if pgerr != nil {
				return nil, pgerr
			}
			linkIds[i] = links[i].Id
			link.Id = links[i].Id.String()
		} else {
			linkId, twerr := internal.GetUuidFromString(link.Id)
			if twerr != nil {
				return nil, twerr.(error)
			}
			linkIds[i] = linkId
		}
	}
	return linkIds, nil
}

func checkRequiredAttributes(userGroup *pb.UserGroup) (twirp.Error) {
	if userGroup.DisplayName == "" || len(userGroup.Avatar) == 0 || userGroup.Address.Data == nil || userGroup.Type.Type == "" || userGroup.OwnerId == "" {
		var argument string
		switch {
		case userGroup.DisplayName == "":
			argument = "display_name"
		case len(userGroup.Avatar) == 0:
			argument = "avatar"
		case userGroup.Address.Data == nil:
			argument = "address"
		case userGroup.Type.Type == "":
			argument = "type"
		case userGroup.OwnerId == "":
			argument = "owner"
		}
		return twirp.RequiredArgumentError(argument)
	}
	return nil
}
