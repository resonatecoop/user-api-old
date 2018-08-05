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
	trackpb "user-api/rpc/track"
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
	createUserGroup := func(userGroup *pb.UserGroup, ownerId uuid.UUID) (error, string, *models.UserGroup) {
		var table string
		tx, err := s.db.Begin()
		if err != nil {
			return err, table, nil
		}
		defer tx.Rollback()

		groupTaxonomy := new(models.GroupTaxonomy)
		pgerr := tx.Model(groupTaxonomy).Where("type = ?", userGroup.Type.Type).First()
		if pgerr != nil {
			return pgerr, "group_taxonomy", nil
		}

		var newAddress *models.StreetAddress
		if userGroup.Address != nil {
			newAddress = &models.StreetAddress{Data: userGroup.Address.Data}
			_, pgerr = tx.Model(newAddress).Returning("*").Insert()
			if pgerr != nil {
				return pgerr, "street_address", nil
			}
		}

		linkIds, pgerr := getLinkIds(userGroup.Links, tx)
		if pgerr != nil {
			return pgerr, "link", nil
		}

		tagIds, pgerr := models.GetTagIds(userGroup.Tags, tx)
		if pgerr != nil {
			return pgerr, "tag", nil
		}

		recommendedArtistIds, pgerr := models.GetRelatedUserGroupIds(userGroup.RecommendedArtists, tx)
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
		}
		_, pgerr = tx.Model(newUserGroup).Returning("*").Insert()
		if pgerr != nil {
			return pgerr, "user_group", nil
		}

		_, pgerr = tx.Exec(`
			UPDATE user_groups
			SET recommended_by = (select array_agg(distinct e) from unnest(recommended_by || ?) e)
			WHERE id IN (?)
		`, pg.Array([]uuid.UUID{newUserGroup.Id}), pg.In(recommendedArtistIds))
		if pgerr != nil {
			return pgerr, "user_group", nil
		}

		pgerr = tx.Model(newUserGroup).
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

	pgerr, table, newUserGroup := createUserGroup(userGroup, ownerId)
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
		Column("user_group.*", "Privacy", "Type", "Address", "Members", "MemberOfGroups").
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
	tags, twerr := models.GetTags(u.Tags, s.db)
	if twerr != nil {
		return nil, twerr
	}

	// Get related user groups
	recommendedArtists, pgerr := models.GetRelatedUserGroups(u.RecommendedArtists, s.db)
	if pgerr != nil {
		return nil, internal.CheckError(pgerr, "user_group")
	}
	members, pgerr, table := getUserGroupMembers(id, u.Members, true, s.db)
	if pgerr != nil {
		return nil, internal.CheckError(pgerr, table)
	}
	memberOfGroups, pgerr, table := getUserGroupMembers(id, u.MemberOfGroups, false, s.db)
	if pgerr != nil {
		return nil, internal.CheckError(pgerr, table)
	}

	// Get related tracks/track groups
	highlightedTracks, pgerr := models.GetTracks(u.HighlightedTracks, s.db, true)
	if pgerr != nil {
		return nil, internal.CheckError(pgerr, "track")
	}
	trackGroups, pgerr := models.GetTrackGroups(u.TrackGroups, s.db, []string{"lp", "ep", "single", "playlist"})
	if pgerr != nil {
		return nil, internal.CheckError(pgerr, "track_group")
	}
	var featuredTrackGroup *trackpb.RelatedTrackGroup
	if (u.FeaturedTrackGroupId != uuid.UUID{}) {
		featuredTrackGroups, pgerr := models.GetTrackGroups([]uuid.UUID{u.FeaturedTrackGroupId}, s.db, []string{"lp", "ep", "single", "playlist"})
		if pgerr != nil {
			return nil, internal.CheckError(pgerr, "track_group")
		}
		featuredTrackGroup = featuredTrackGroups[0]
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
		Members: members,
		MemberOfGroups: memberOfGroups,
		HighlightedTracks: highlightedTracks,
		FeaturedTrackGroup: featuredTrackGroup,
		TrackGroups: trackGroups,
	}, nil
}

func (s *Server) UpdateUserGroup(ctx context.Context, userGroup *pb.UserGroup) (*trackpb.Empty, error) {
	updateUserGroup := func(userGroup *pb.UserGroup, u *models.UserGroup) (error, string) {
		var table string
		tx, err := s.db.Begin()
		if err != nil {
			return err, "user_group"
		}
		defer tx.Rollback()

		// Update address
		addressId, twerr := internal.GetUuidFromString(userGroup.Address.Id)
		if twerr != nil {
			return twerr, "street_address"
		}
		address := &models.StreetAddress{Id: addressId, Data: userGroup.Address.Data}
		_, pgerr := tx.Model(address).Column("data").WherePK().Update()
		// _, pgerr := db.Model(address).Set("data = ?", pg.Hstore(userGroup.Address.Data)).Where("id = ?id").Update()
		if pgerr != nil {
			return pgerr, "street_address"
		}

		// Update privacy
		privacyId, twerr := internal.GetUuidFromString(userGroup.Privacy.Id)
		if twerr != nil {
			return twerr, "user_group_privacy"
		}
		privacy := &models.UserGroupPrivacy{Id: privacyId, Private: userGroup.Privacy.Private, OwnedTracks: userGroup.Privacy.OwnedTracks, SupportedArtists: userGroup.Privacy.SupportedArtists}
		_, pgerr = tx.Model(privacy).WherePK().Returning("*").UpdateNotNull()
		if pgerr != nil {
			return pgerr, "user_group_privacy"
		}

		// Update tags
		tagIds, pgerr := models.GetTagIds(userGroup.Tags, tx)
		if pgerr != nil {
			return pgerr, "tag"
		}

		// Update links
		linkIds, pgerr := getLinkIds(userGroup.Links, tx)
		if pgerr != nil {
			return pgerr, "link"
		}
		// Delete links if needed
		pgerr = tx.Model(u).WherePK().Column("links").Select()
		if pgerr != nil {
			return pgerr, "user_group"
		}
		linkIdsToDelete := internal.Difference(u.Links, linkIds)
		if len(linkIdsToDelete) > 0 {
			_, pgerr = tx.Model((*models.Link)(nil)).
				Where("id in (?)", pg.In(linkIdsToDelete)).
				Delete()
			if pgerr != nil {
				return pgerr, "link"
			}
		}

		// Update recommended artists
		// recommendedArtistIds, pgerr := models.GetRelatedUserGroupIds(userGroup.RecommendedArtists, tx)
		// if pgerr != nil {
		// 	return pgerr, "user_group"
		// }

		// Update user group
		u.Tags = tagIds
		u.Links = linkIds
		// u.RecommendedArtists = recommendedArtistIds
		u.UpdatedAt = time.Now()
		_, pgerr = tx.Model(u).WherePK().Returning("*").UpdateNotNull()
		if pgerr != nil {
			return pgerr, "user_group"
		}

		return tx.Commit(), table
	}

	err := checkRequiredAttributes(userGroup)
	if err != nil {
		return nil, err
	}

	u, err := getUserGroupModel(userGroup)
	if err != nil {
		return nil, err
	}

	if pgerr, table := updateUserGroup(userGroup, u); pgerr != nil {
		return nil, internal.CheckError(pgerr, table)
	}

  return &trackpb.Empty{}, nil
}

func (s *Server) DeleteUserGroup(ctx context.Context, userGroup *pb.UserGroup) (*trackpb.Empty, error) {
	id, twerr := internal.GetUuidFromString(userGroup.Id)
	if twerr != nil {
		return nil, twerr
	}
	u := &models.UserGroup{Id: id}

	tx, err := s.db.Begin()
	if err != nil {
		return nil, internal.CheckError(err, "")
	}
	defer tx.Rollback()

	if pgerr, table := u.Delete(tx); pgerr != nil {
		return nil, internal.CheckError(pgerr, table)
	}

	err = tx.Commit()
	if err != nil {
		return nil, internal.CheckError(err, "")
	}

	return &trackpb.Empty{}, nil
}

func (s *Server) GetLabelUserGroups(ctx context.Context, empty *trackpb.Empty) (*pb.GroupedUserGroups, error) {
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

func (s *Server) GetUserGroupTypes(ctx context.Context, empty *trackpb.Empty) (*pb.GroupTaxonomies, error) {
	var types []models.GroupTaxonomy
	var groupTaxonomies pb.GroupTaxonomies
	err := s.db.Model(&types).
		Where("group_taxonomy.type != ?", "distributor"). // except distributors, internally added by staff
		Select()
	if err != nil {
		return nil, internal.CheckError(err, "group_taxonomy")
	}

	for _, groupType := range(types) {
		groupTaxonomies.Types = append(groupTaxonomies.Types, &pb.GroupTaxonomy{
			Id: groupType.Id.String(),
			Type: groupType.Type,
			Name: groupType.Name,
		})
	}

  return &groupTaxonomies, nil
}

func (s *Server) AddMembers(ctx context.Context, userGroupMembers *pb.UserGroupMembers) (*trackpb.Empty, error) {
	addMembers := func() (error, string) {
		var table string
		tx, err := s.db.Begin()
		if err != nil {
			return err, table
		}
		defer tx.Rollback()

		userGroupId, err := internal.GetUuidFromString(userGroupMembers.UserGroupId)
		if err != nil {
			return err, "user_group"
		}

		for _, member := range(userGroupMembers.Members) {
			// verify uuid
			memberId, err := internal.GetUuidFromString(member.Id)
			if err != nil {
				return err, "member"
			}

			// get user_group (should exist)
			m := &models.UserGroup{Id: memberId}
			pgerr := tx.Model(m).WherePK().Select()
			if pgerr != nil {
				return pgerr, "user_group"
			}

			userGroupMember := &models.UserGroupMember{UserGroupId: userGroupId, MemberId: memberId}

			// set display_name
			// if not provided, default will be member user_group display_name
			if member.DisplayName != "" {
				userGroupMember.DisplayName = member.DisplayName
			} else {
				userGroupMember.DisplayName = m.DisplayName
			}

			// create tags
			tagIds, pgerr := models.GetTagIds(member.Tags, tx)
			if pgerr != nil {
				return pgerr, "tag"
			}
			userGroupMember.Tags = tagIds

			// create UserGroup/Member relation
			_, pgerr = tx.Model(userGroupMember).Insert()
			if pgerr != nil {
				return pgerr, "user_group_member"
			}
		}

		return tx.Commit(), table
	}
	if err, table := addMembers(); err != nil {
		return nil, internal.CheckError(err, table)
	}
	return &trackpb.Empty{}, nil
}

func (s *Server) DeleteMembers(ctx context.Context, userGroupMembers *pb.UserGroupMembers) (*trackpb.Empty, error) {
	deleteMembers := func() (error, string) {
		var table string
		tx, err := s.db.Begin()
		if err != nil {
			return err, table
		}
		defer tx.Rollback()

		userGroupId, err := internal.GetUuidFromString(userGroupMembers.UserGroupId)
		if err != nil {
			return err, "user_group"
		}

		for _, member := range(userGroupMembers.Members) {
			// verify uuid
			memberId, err := internal.GetUuidFromString(member.Id)
			if err != nil {
				return err, "member"
			}

			// delete UserGroup/Member relation
			// we don't delete tags because they could be used for other members
			// avoiding having multiple tags that represent the same thing
			userGroupMember := &models.UserGroupMember{UserGroupId: userGroupId, MemberId: memberId}
			res, pgerr := tx.Model(userGroupMember).
				Where("user_group_id = ?user_group_id").
				Where("member_id = ?member_id").
				Delete()
			if res.RowsAffected() == 0 {
				return pg.ErrNoRows, "user_group"
			}
			if pgerr != nil {
				return pgerr, "user_group_member"
			}
		}

		return tx.Commit(), table
	}
	if err, table := deleteMembers(); err != nil {
		return nil, internal.CheckError(err, table)
	}
	return &trackpb.Empty{}, nil
}

func (s *Server) AddRecommended(ctx context.Context, userGroupRecommended *pb.UserGroupRecommended) (*trackpb.Empty, error) {
	userGroupId, twerr := internal.GetUuidFromString(userGroupRecommended.UserGroupId)
	if twerr != nil {
		return nil, twerr
	}
	recommendedId, twerr := internal.GetUuidFromString(userGroupRecommended.RecommendedId)
	if twerr != nil {
		return nil, twerr
	}
	u := &models.UserGroup{Id: userGroupId}

	if pgerr, table := u.AddRecommended(s.db, recommendedId); pgerr != nil {
		return nil, internal.CheckError(pgerr, table)
	}
	return &trackpb.Empty{}, nil
}

func (s *Server) RemoveRecommended(ctx context.Context, userGroupRecommended *pb.UserGroupRecommended) (*trackpb.Empty, error) {
	userGroupId, twerr := internal.GetUuidFromString(userGroupRecommended.UserGroupId)
	if twerr != nil {
		return nil, twerr
	}
	recommendedId, twerr := internal.GetUuidFromString(userGroupRecommended.RecommendedId)
	if twerr != nil {
		return nil, twerr
	}
	u := &models.UserGroup{Id: userGroupId}

	if pgerr, table := u.RemoveRecommended(s.db, recommendedId); pgerr != nil {
		return nil, internal.CheckError(pgerr, table)
	}
	return &trackpb.Empty{}, nil
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

func getUserGroupMembers(userGroupId uuid.UUID, userGroups []models.UserGroup, members bool, db *pg.DB) ([]*pb.UserGroup, error, string) {
	userGroupsResponse := make([]*pb.UserGroup, len(userGroups))
	for i, userGroup := range userGroups {
		u := &pb.UserGroup{Avatar: userGroup.Avatar}
		userGroupMember := models.UserGroupMember{}
		if members { // userGroups are members of user group with userGroupId
			userGroupMember.UserGroupId = userGroupId
			userGroupMember.MemberId = userGroup.Id
			u.Id = userGroupMember.MemberId.String()
		} else { // user group with userGroupId is member of userGroups
			userGroupMember.UserGroupId = userGroup.Id
			userGroupMember.MemberId = userGroupId
			u.Id = userGroupMember.UserGroupId.String()
		}
		err := db.Model(&userGroupMember).
			WherePK().
			// Where("user_group_id = ?user_group_id").
			// Where("member_id = ?member_id").
			Select()
		if err != nil {
			return nil, err, "user_group_member"
		}

		u.DisplayName = userGroupMember.DisplayName

		// get tags
		if len(userGroupMember.Tags) > 0 {
			var tags []*models.Tag
			err = db.Model(&tags).
				Where("id in (?)", pg.In(userGroupMember.Tags)).
				Select()
			if err != nil {
				return nil, err, "tag"
			}
			for _, tag := range tags {
				u.Tags = append(u.Tags, &trackpb.Tag{Id: tag.Id.String(), Type: tag.Type, Name: tag.Name})
			}
		}
		userGroupsResponse[i] = u
	}
	return userGroupsResponse, nil, ""
}

func getLinkIds(l []*pb.Link, db *pg.Tx) ([]uuid.UUID, error) {
	links := make([]*models.Link, len(l))
	linkIds := make([]uuid.UUID, len(l))
	for i, link := range l {
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
