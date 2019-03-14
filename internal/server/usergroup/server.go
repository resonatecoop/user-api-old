package usergroupserver

import (
	// "fmt"
	// "reflect"
	// "time"
	"context"
	"net/url"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/twitchtv/twirp"
	"github.com/satori/go.uuid"

  userpb "user-api/rpc/user"
	pb "user-api/rpc/usergroup"
	// trackpb "user-api/rpc/track"
	tagpb "user-api/rpc/tag"
	"user-api/internal"
	"user-api/internal/database/model"
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

	u := &model.UserGroup{
		DisplayName: userGroup.DisplayName,
		Description: userGroup.Description,
		OwnerId: ownerId,
		ShortBio: userGroup.ShortBio,
		Avatar: userGroup.Avatar,
		Banner: userGroup.Banner,
		GroupEmailAddress: userGroup.GroupEmailAddress,
		Publisher: userGroup.Publisher,
		Pro: userGroup.Pro,
	}

	pgerr, table := u.Create(s.db, userGroup)
	if pgerr != nil {
		return nil, internal.CheckError(pgerr, table)
	}

  return &pb.UserGroup{
		Id: u.Id.String(),
		DisplayName: u.DisplayName,
		Avatar: u.Avatar,
		Banner: u.Banner,
		ShortBio: u.ShortBio,
		OwnerId: u.OwnerId.String(),
		Type: userGroup.Type,
		Address: userGroup.Address,
		Links: userGroup.Links,
		Tags: userGroup.Tags,
		RecommendedArtists: userGroup.RecommendedArtists,
		Privacy: userGroup.Privacy,
		Publisher: userGroup.Publisher,
		Pro: userGroup.Pro,
	}, nil
}

func (s *Server) SearchUserGroups(ctx context.Context, q *tagpb.Query) (*tagpb.SearchResults, error) {
  if len(q.Query) < 3 {
    return nil, twirp.InvalidArgumentError("query", "must be a valid search query")
  }

  searchResults, twerr := model.SearchUserGroups(q.Query, s.db)
  if twerr != nil {
    return nil, twerr
  }
  return searchResults, nil
}

// TODO handle privacy settings
func (s *Server) GetUserGroup(ctx context.Context, userGroup *pb.UserGroup) (*pb.UserGroup, error) {
	id, err := internal.GetUuidFromString(userGroup.Id)
	if err != nil {
		return nil, err
	}
	u := &model.UserGroup{Id: id}

	pgerr := s.db.Model(u).
		Column("user_group.*", "Privacy", "Type", "Address", "Members", "MemberOfGroups",
			"OwnerOfTrackGroups", "LabelOfTrackGroups").
		WherePK().
		Select()
	if pgerr != nil {
		return nil, internal.CheckError(pgerr, "user_group")
	}

	// Get user group links
	links := make([]*pb.Link, len(u.Links))
	if len(links) > 0 {
		var groupLinks []model.Link
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
	tags, twerr := model.GetTags(u.Tags, s.db)
	if twerr != nil {
		return nil, twerr
	}

	// Get related user groups
	recommendedArtists, pgerr := model.GetRelatedUserGroups(u.RecommendedArtists, s.db)
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
	highlightedTracks, pgerr := model.GetTracks(u.HighlightedTracks, s.db, true, ctx)
	if pgerr != nil {
		return nil, internal.CheckError(pgerr, "track")
	}

	trackGroups := model.GetTrackGroups(append(u.OwnerOfTrackGroups, u.LabelOfTrackGroups...))

	var featuredTrackGroup *tagpb.RelatedTrackGroup
	if (u.FeaturedTrackGroupId != uuid.UUID{}) {
		featuredTrackGroups, pgerr := model.GetTrackGroupsFromIds([]uuid.UUID{u.FeaturedTrackGroupId}, s.db, []string{"lp", "ep", "single", "playlist"})
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
		Publisher: u.Publisher,
		Pro: u.Pro,
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

func (s *Server) UpdateUserGroup(ctx context.Context, userGroup *pb.UserGroup) (*tagpb.Empty, error) {
	err := checkRequiredAttributes(userGroup)
	if err != nil {
		return nil, err
	}

	u, err := getUserGroupModel(userGroup)
	if err != nil {
		return nil, err
	}

	if pgerr, table := u.Update(s.db, userGroup); pgerr != nil {
		return nil, internal.CheckError(pgerr, table)
	}

  return &tagpb.Empty{}, nil
}

func (s *Server) DeleteUserGroup(ctx context.Context, userGroup *pb.UserGroup) (*tagpb.Empty, error) {
	id, twerr := internal.GetUuidFromString(userGroup.Id)
	if twerr != nil {
		return nil, twerr
	}
	u := &model.UserGroup{Id: id}

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

	return &tagpb.Empty{}, nil
}

func (s *Server) GetLabelUserGroups(ctx context.Context, empty *tagpb.Empty) (*pb.GroupedUserGroups, error) {
	var labels []model.UserGroup

	// err := s.db.Model(&labels).
	// 	Join("LEFT JOIN group_taxonomies AS g").
	// 	JoinOn("g.id = user_group.type_id").
	// 	Where("g.type = ?", "label").
	//   Apply(orm.Pagination(ctx.Value("query").(url.Values))).
	//   Select()

	groupTaxonomy := new(model.GroupTaxonomy)
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

func (s *Server) GetUserGroupTypes(ctx context.Context, empty *tagpb.Empty) (*pb.GroupTaxonomies, error) {
	var types []model.GroupTaxonomy
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

func (s *Server) AddMembers(ctx context.Context, userGroupMembers *pb.UserGroupMembers) (*tagpb.Empty, error) {
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
			m := &model.UserGroup{Id: memberId}
			pgerr := tx.Model(m).WherePK().Select()
			if pgerr != nil {
				return pgerr, "user_group"
			}

			userGroupMember := &model.UserGroupMember{UserGroupId: userGroupId, MemberId: memberId}

			// set display_name
			// if not provided, default will be member user_group display_name
			if member.DisplayName != "" {
				userGroupMember.DisplayName = member.DisplayName
			} else {
				userGroupMember.DisplayName = m.DisplayName
			}

			// create tags
			tagIds, pgerr := model.GetTagIds(member.Tags, tx)
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
	return &tagpb.Empty{}, nil
}

func (s *Server) DeleteMembers(ctx context.Context, userGroupMembers *pb.UserGroupMembers) (*tagpb.Empty, error) {
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
			userGroupMember := &model.UserGroupMember{UserGroupId: userGroupId, MemberId: memberId}
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
	return &tagpb.Empty{}, nil
}

func (s *Server) AddRecommended(ctx context.Context, userGroupRecommended *pb.UserGroupRecommended) (*tagpb.Empty, error) {
	userGroupId, twerr := internal.GetUuidFromString(userGroupRecommended.UserGroupId)
	if twerr != nil {
		return nil, twerr
	}
	recommendedId, twerr := internal.GetUuidFromString(userGroupRecommended.RecommendedId)
	if twerr != nil {
		return nil, twerr
	}
	u := &model.UserGroup{Id: userGroupId}

	if pgerr, table := u.AddRecommended(s.db, recommendedId); pgerr != nil {
		return nil, internal.CheckError(pgerr, table)
	}
	return &tagpb.Empty{}, nil
}

func (s *Server) RemoveRecommended(ctx context.Context, userGroupRecommended *pb.UserGroupRecommended) (*tagpb.Empty, error) {
	userGroupId, twerr := internal.GetUuidFromString(userGroupRecommended.UserGroupId)
	if twerr != nil {
		return nil, twerr
	}
	recommendedId, twerr := internal.GetUuidFromString(userGroupRecommended.RecommendedId)
	if twerr != nil {
		return nil, twerr
	}
	u := &model.UserGroup{Id: userGroupId}

	if pgerr, table := u.RemoveRecommended(s.db, recommendedId); pgerr != nil {
		return nil, internal.CheckError(pgerr, table)
	}
	return &tagpb.Empty{}, nil
}

/*func (s *Server) GetTrackAnalytics(ctx context.Context, userGroup *pb.UserGroup) (*pb.UserGroupTrackAnalytics, error) {
	id, err := internal.GetUuidFromString(userGroup.Id)
	if err != nil {
		return nil, err
	}
	u := &model.UserGroup{Id: id}
	pgerr := s.db.Model(u).
		Column("Type").
		WherePK().
		Select()
	if pgerr != nil {
		return nil, internal.CheckError(pgerr, "user_group")
	}

	if u.Type.Type == "artist" {
		artistTrackAnalytics, twerr := u.GetUserGroupTrackAnalytics(s.db)
		if twerr != nil {
			return nil, twerr
		}
		return &pb.UserGroupTrackAnalytics{
	    ArtistTrackAnalytics: artistTrackAnalytics,
	  }, nil
	} else if u.Type.Type == "label" {
		pgerr := s.db.Model(u).
			Column("user_group.id", "user_group.display_name", "user_group.avatar", "Members").
			WherePK().
			Select()
		if pgerr != nil {
			return nil, internal.CheckError(pgerr, "user_group")
		}
		userGroups := make([]*pb.LabelTrackAnalytics, len(u.Members)+1)

		// Track analytics of artists members of label
		for i, member := range u.Members {
			trackAnalytics, twerr := member.GetUserGroupTrackAnalytics(s.db)
			if twerr != nil {
				return nil, twerr
			}
			labelArtistTrackAnalytics := &pb.LabelTrackAnalytics{
				UserGroup: &tagpb.RelatedUserGroup{
					Id: member.Id.String(),
					DisplayName: member.DisplayName,
					Avatar: member.Avatar,
				},
				Tracks: trackAnalytics,
			}
			userGroups[i] = labelArtistTrackAnalytics
		}

		// Track analytics of label
		trackAnalytics, twerr := u.GetUserGroupTrackAnalytics(s.db)
		if twerr != nil {
			return nil, twerr
		}
		userGroups[len(u.Members)] = &pb.LabelTrackAnalytics{
			UserGroup: &tagpb.RelatedUserGroup{
				Id: u.Id.String(),
				DisplayName: u.DisplayName,
				Avatar: u.Avatar,
			},
			Tracks: trackAnalytics,
		}
		return &pb.UserGroupTrackAnalytics{
			LabelTrackAnalytics: userGroups,
		}, nil
	}

	return &pb.UserGroupTrackAnalytics{}, nil
}*/

func getUserGroupModel(userGroup *pb.UserGroup) (*model.UserGroup, twirp.Error) {
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
	// typeId, err := internal.GetUuidFromString(userGroup.Type.Id)
	// if err != nil {
	// 	return nil, err
	// }
	// ownerId, err := internal.GetUuidFromString(userGroup.OwnerId)
	// if err != nil {
	// 	return nil, err
	// }
	return &model.UserGroup{
		Id: id,
		DisplayName: userGroup.DisplayName,
		Description: userGroup.Description,
		AddressId: addressId,
		PrivacyId: privacyId,
		// TypeId: typeId,
		// OwnerId: ownerId,
		ShortBio: userGroup.ShortBio,
		Avatar: userGroup.Avatar,
		Banner: userGroup.Banner,
		GroupEmailAddress: userGroup.GroupEmailAddress,
		Publisher: userGroup.Publisher,
		Pro: userGroup.Pro,
	}, nil
}

func getUserGroupMembers(userGroupId uuid.UUID, userGroups []model.UserGroup, members bool, db *pg.DB) ([]*pb.UserGroup, error, string) {
	userGroupsResponse := make([]*pb.UserGroup, len(userGroups))
	for i, userGroup := range userGroups {
		u := &pb.UserGroup{Avatar: userGroup.Avatar}
		userGroupMember := model.UserGroupMember{}
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
			var tags []*model.Tag
			err = db.Model(&tags).
				Where("id in (?)", pg.In(userGroupMember.Tags)).
				Select()
			if err != nil {
				return nil, err, "tag"
			}
			for _, tag := range tags {
				u.Tags = append(u.Tags, &tagpb.Tag{Id: tag.Id.String(), Type: tag.Type, Name: tag.Name})
			}
		}
		userGroupsResponse[i] = u
	}
	return userGroupsResponse, nil, ""
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
