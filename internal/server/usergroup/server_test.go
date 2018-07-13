package usergroupserver_test

import (
	// "fmt"
	// "reflect"
	"context"
	"net/url"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/twitchtv/twirp"
	"github.com/satori/go.uuid"

	pb "user-api/rpc/usergroup"
	userpb "user-api/rpc/user"
	// "user-api/internal"
	"user-api/internal/database/models"
)

var _ = Describe("UserGroup server", func() {
  const already_exists_code twirp.ErrorCode = "already_exists"
  const invalid_argument_code twirp.ErrorCode = "invalid_argument"
  const not_found_code twirp.ErrorCode = "not_found"

	Describe("GetUserGroup", func() {
		Context("with valid uuid", func() {
			It("should respond with user_group if it exists", func() {
				userGroup := &pb.UserGroup{Id: newArtist.Id.String()}
				resp, err := service.GetUserGroup(context.Background(), userGroup)

				Expect(err).NotTo(HaveOccurred())
				Expect(resp.Id).To(Equal(newArtist.Id.String()))
				Expect(resp.DisplayName).To(Equal(newArtist.DisplayName))
				Expect(resp.Description).To(Equal(newArtist.Description))
				Expect(resp.ShortBio).To(Equal(newArtist.ShortBio))
				Expect(resp.Avatar).To(Equal(newArtist.Avatar))
				Expect(resp.Banner).To(Equal(newArtist.Banner))
				Expect(resp.OwnerId).To(Equal(newArtist.OwnerId.String()))
				Expect(resp.Type.Id).To(Equal(newArtistGroupTaxonomy.Id.String()))
				Expect(resp.Type.Type).To(Equal("artist"))
				Expect(len(resp.Labels)).To(Equal(1))
				Expect(resp.Labels[0].Id).To(Equal(newLabel.Id.String()))
				Expect(resp.Labels[0].DisplayName).To(Equal(newLabel.DisplayName))
				Expect(resp.Labels[0].Avatar).To(Equal(newLabel.Avatar))
				Expect(len(resp.Tags)).To(Equal(1))
				Expect(resp.Tags[0].Id).To(Equal(newTag.Id.String()))
				Expect(resp.Tags[0].Type).To(Equal(newTag.Type))
				Expect(resp.Tags[0].Name).To(Equal(newTag.Name))
				Expect(len(resp.Links)).To(Equal(1))
				Expect(resp.Links[0].Id).To(Equal(newLink.Id.String()))
				Expect(resp.Links[0].Uri).To(Equal(newLink.Uri))
				Expect(resp.Links[0].Platform).To(Equal(newLink.Platform))
				Expect(resp.GroupEmailAddress).To(Equal(newArtist.GroupEmailAddress))
			})
			It("should respond with not_found error if user_group does not exist", func() {
				id := uuid.NewV4()
				for (id == newArtist.Id || id == newRecommendedArtist.Id || id == newLabel.Id || id == newDistributor.Id) {
					id = uuid.NewV4()
				}
				userGroup := &pb.UserGroup{Id: id.String()}
				resp, err := service.GetUserGroup(context.Background(), userGroup)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(not_found_code))
			})
		})
		Context("with invalid uuid", func() {
			It("should respond with invalid_argument error", func() {
				id := "45"
				userGroup := &pb.UserGroup{Id: id}
				resp, err := service.GetUserGroup(context.Background(), userGroup)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
				Expect(twerr.Meta("argument")).To(Equal("id"))
			})
		})
	})

	Describe("GetLabelUserGroups", func() {
		It("should respond with user_groups of type label", func() {
			emptyReq := &userpb.Empty{}
			u := url.URL{}
			ctx := context.WithValue(context.Background(), "query", u.Query())
			resp, err := service.GetLabelUserGroups(ctx, emptyReq)

			Expect(err).NotTo(HaveOccurred())
			Expect(resp).NotTo(BeNil())
			Expect(len(resp.Labels)).To(Equal(1))
			Expect(resp.Labels[0].Id).To(Equal(newLabel.Id.String()))
			Expect(resp.Labels[0].DisplayName).To(Equal(newLabel.DisplayName))
			Expect(resp.Labels[0].Avatar).To(Equal(newLabel.Avatar))
		})
	})

	Describe("UpdateUserGroup", func() {
		Context("with valid uuid", func() {
			It("should update user_group if it exists", func() {
				tags := []*pb.Tag{&pb.Tag{Type: "genre", Name: "experimental"}}
				links := []*pb.Link{&pb.Link{Platform: "instagram", Uri: "https://instagram/bestartistever"}}
				recommendedArtists := []*pb.UserGroup{&pb.UserGroup{Id: newRecommendedArtist.Id.String()}}
				userGroup := &pb.UserGroup{
					Id: newArtist.Id.String(),
					DisplayName: "new display name",
					Description: "new description",
					Avatar: newArtist.Avatar,
					Address: &userpb.StreetAddress{Id: newAddress.Id.String(), Data: map[string]string{"some": "new data"}},
					Type: &pb.GroupTaxonomy{Id: newArtistGroupTaxonomy.Id.String(), Type: "artist"},
					Privacy: &pb.Privacy{Id: newArtist.Privacy.Id.String(), Private: true, OwnedTracks: false, SupportedArtists: true},
					OwnerId: newArtist.OwnerId.String(),
					Tags: tags,
					Links: links,
					RecommendedArtists: recommendedArtists,
				}
				_, err := service.UpdateUserGroup(context.Background(), userGroup)

				Expect(err).NotTo(HaveOccurred())

				address := new(models.StreetAddress)
				err = db.Model(address).Where("id = ?", newAddress.Id).Select()
				Expect(err).NotTo(HaveOccurred())
				Expect(address.Data).To(Equal(map[string]string{"some": "new data"}))

				privacy := new(models.UserGroupPrivacy)
				err = db.Model(privacy).Where("id = ?", newArtist.Privacy.Id).Select()
				Expect(err).NotTo(HaveOccurred())
				Expect(privacy.Private).To(Equal(true))
				Expect(privacy.OwnedTracks).To(Equal(false))
				Expect(privacy.SupportedArtists).To(Equal(true))

				updatedUserGroup := new(models.UserGroup)
				err = db.Model(updatedUserGroup).Where("id = ?", newArtist.Id).Select()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(updatedUserGroup.Tags)).To(Equal(1))
				Expect(updatedUserGroup.Tags[0]).NotTo(Equal(newTag.Id))

				addedTag := models.Tag{Id: updatedUserGroup.Tags[0]}
				err = db.Model(&addedTag).WherePK().Returning("*").Select()
				Expect(addedTag.Type).To(Equal("genre"))
				Expect(addedTag.Name).To(Equal("experimental"))

				Expect(len(updatedUserGroup.Links)).To(Equal(1))
				addedLink := models.Link{Id: updatedUserGroup.Links[0]}
				err = db.Model(&addedLink).WherePK().Returning("*").Select()
				Expect(addedLink.Platform).To(Equal("instagram"))
				Expect(addedLink.Uri).To(Equal("https://instagram/bestartistever"))
				err = db.Model(newLink).WherePK().Returning("*").Select()
				Expect(err).To(HaveOccurred())

				Expect(len(updatedUserGroup.RecommendedArtists)).To(Equal(1))
				Expect(updatedUserGroup.RecommendedArtists[0]).To(Equal(newRecommendedArtist.Id))

			})
			It("should respond with not_found error if user_group does not exist", func() {
				id := uuid.NewV4()
				for (id == newArtist.Id || id == newRecommendedArtist.Id || id == newLabel.Id || id == newDistributor.Id) {
					id = uuid.NewV4()
				}
				userGroup := &pb.UserGroup{
					Id: id.String(),
					DisplayName: "new display name",
					Description: "new description",
					Avatar: newArtist.Avatar,
					Address: &userpb.StreetAddress{Id: newAddress.Id.String(), Data: map[string]string{"some": "data"}},
					Type: &pb.GroupTaxonomy{Id: newArtistGroupTaxonomy.Id.String(), Type: "artist"},
					Privacy: &pb.Privacy{Id: newArtist.Privacy.Id.String()},
					OwnerId: newArtist.OwnerId.String(),
				}
				resp, err := service.UpdateUserGroup(context.Background(), userGroup)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(not_found_code))
			})
		})
		Context("with invalid uuid", func() {
			It("should respond with invalid_argument error", func() {
				id := "45"
				userGroup := &pb.UserGroup{
					Id: id,
					DisplayName: "new display name",
					Description: "new description",
					Avatar: newArtist.Avatar,
					Address: &userpb.StreetAddress{Id: newAddress.Id.String(), Data: map[string]string{"some": "data"}},
					Type: &pb.GroupTaxonomy{Id: newArtistGroupTaxonomy.Id.String(), Type: "artist"},
					OwnerId: newArtist.OwnerId.String(),
					Privacy: &pb.Privacy{Id: newArtist.Privacy.Id.String()},
				}
				resp, err := service.UpdateUserGroup(context.Background(), userGroup)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
				Expect(twerr.Meta("argument")).To(Equal("id"))
			})
		})
	})

	Describe("AddSubGroups", func() {
		Context("with valid label userGroupId and subGroupIds", func() {
			It("should add subGroups to userGroup's SubGroups and userGroup to subGroups's Labels", func() {
				userGroupToUserGroups := &pb.UserGroupToUserGroups{
					UserGroupId: newLabel.Id.String(),
					UserGroups: []*pb.UserGroup{
						&pb.UserGroup{Id: newArtist.Id.String()},
					},
				}
				_, err := service.AddSubGroups(context.Background(), userGroupToUserGroups)

				Expect(err).NotTo(HaveOccurred())

				label := models.UserGroup{Id: newLabel.Id}
				err = db.Model(&label).
					Column("user_group.sub_groups").
					WherePK().
					Select()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(label.SubGroups)).To(Equal(1))
				Expect(label.SubGroups[0]).To(Equal(newArtist.Id))

				artist := models.UserGroup{Id: newArtist.Id}
				err = db.Model(&artist).
					Column("user_group.labels").
					WherePK().
					Select()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(artist.Labels)).To(Equal(1))
				Expect(artist.Labels[0]).To(Equal(newLabel.Id))
			})
		})

		Context("with valid distributor userGroupId and subGroupIds", func() {
			It("should add subGroups to userGroup's SubGroups", func() {
				userGroupToUserGroups := &pb.UserGroupToUserGroups{
					UserGroupId: newDistributor.Id.String(),
					UserGroups: []*pb.UserGroup{
						&pb.UserGroup{Id: newLabel.Id.String()},
					},
				}
				_, err := service.AddSubGroups(context.Background(), userGroupToUserGroups)

				Expect(err).NotTo(HaveOccurred())

				distributor := models.UserGroup{Id: newDistributor.Id}
				err = db.Model(&distributor).
					Column("user_group.sub_groups").
					WherePK().
					Select()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(distributor.SubGroups)).To(Equal(1))
				Expect(distributor.SubGroups[0]).To(Equal(newLabel.Id))

				label := models.UserGroup{Id: newLabel.Id}
				err = db.Model(&label).
					Column("user_group.labels").
					WherePK().
					Select()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(label.Labels)).To(Equal(0))
			})
		})

		Context("with invalid userGroupId or subGroupIds", func() {
			It("should respond with invalid argument error if userGroupId not valid", func() {
				userGroupToUserGroups := &pb.UserGroupToUserGroups{
					UserGroupId: "",
					UserGroups: []*pb.UserGroup{
						&pb.UserGroup{Id: newLabel.Id.String()},
					},
				}
				resp, err := service.AddSubGroups(context.Background(), userGroupToUserGroups)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
				Expect(twerr.Meta("argument")).To(Equal("id"))
			})
			It("should respond with invalid argument error if one of subGroupIds not valid", func() {
				userGroupToUserGroups := &pb.UserGroupToUserGroups{
					UserGroupId: newLabel.Id.String(),
					UserGroups: []*pb.UserGroup{
						&pb.UserGroup{Id: ""},
					},
				}
				resp, err := service.AddSubGroups(context.Background(), userGroupToUserGroups)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
				Expect(twerr.Meta("argument")).To(Equal("id"))
			})
			It("should respond with not found error if userGroup does not exist", func() {
				id := uuid.NewV4()
				for (id == newArtist.Id || id == newRecommendedArtist.Id || id == newLabel.Id || id == newDistributor.Id) {
					id = uuid.NewV4()
				}
				userGroupToUserGroups := &pb.UserGroupToUserGroups{
					UserGroupId: id.String(),
					UserGroups: []*pb.UserGroup{
						&pb.UserGroup{Id: newLabel.Id.String()},
					},
				}
				resp, err := service.AddSubGroups(context.Background(), userGroupToUserGroups)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(not_found_code))
			})
			It("should respond with not found error if one of subGroup does not exist", func() {
				id := uuid.NewV4()
				for (id == newArtist.Id || id == newRecommendedArtist.Id || id == newLabel.Id || id == newDistributor.Id) {
					id = uuid.NewV4()
				}
				userGroupToUserGroups := &pb.UserGroupToUserGroups{
					UserGroupId: newLabel.Id.String(),
					UserGroups: []*pb.UserGroup{
						&pb.UserGroup{Id: id.String()},
					},
				}
				resp, err := service.AddSubGroups(context.Background(), userGroupToUserGroups)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(not_found_code))
			})
		})
	})

	Describe("RemoveSubGroups", func() {
		Context("with valid label userGroupId and subGroupIds", func() {
			It("should remove subGroups from userGroup's SubGroups and userGroup from subGroups's Labels", func() {
				userGroupToUserGroups := &pb.UserGroupToUserGroups{
					UserGroupId: newLabel.Id.String(),
					UserGroups: []*pb.UserGroup{
						&pb.UserGroup{Id: newArtist.Id.String()},
					},
				}
				_, err := service.RemoveSubGroups(context.Background(), userGroupToUserGroups)

				Expect(err).NotTo(HaveOccurred())

				label := models.UserGroup{Id: newLabel.Id}
				err = db.Model(&label).
					Column("user_group.sub_groups").
					WherePK().
					Select()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(label.SubGroups)).To(Equal(0))

				artist := models.UserGroup{Id: newArtist.Id}
				err = db.Model(&artist).
					Column("user_group.labels").
					WherePK().
					Select()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(artist.Labels)).To(Equal(0))
			})
		})
		Context("with invalid userGroupId or subGroupIds", func() {
			It("should respond with invalid argument error if userGroupId not valid", func() {
				userGroupToUserGroups := &pb.UserGroupToUserGroups{
					UserGroupId: "",
					UserGroups: []*pb.UserGroup{
						&pb.UserGroup{Id: newLabel.Id.String()},
					},
				}
				resp, err := service.RemoveSubGroups(context.Background(), userGroupToUserGroups)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
				Expect(twerr.Meta("argument")).To(Equal("id"))
			})
			It("should respond with invalid argument error if one of subGroupIds not valid", func() {
				userGroupToUserGroups := &pb.UserGroupToUserGroups{
					UserGroupId: newLabel.Id.String(),
					UserGroups: []*pb.UserGroup{
						&pb.UserGroup{Id: ""},
					},
				}
				resp, err := service.RemoveSubGroups(context.Background(), userGroupToUserGroups)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
				Expect(twerr.Meta("argument")).To(Equal("id"))
			})
			It("should respond with not found error if userGroup does not exist", func() {
				id := uuid.NewV4()
				for (id == newArtist.Id || id == newRecommendedArtist.Id || id == newLabel.Id || id == newDistributor.Id) {
					id = uuid.NewV4()
				}
				userGroupToUserGroups := &pb.UserGroupToUserGroups{
					UserGroupId: id.String(),
					UserGroups: []*pb.UserGroup{
						&pb.UserGroup{Id: newLabel.Id.String()},
					},
				}
				resp, err := service.RemoveSubGroups(context.Background(), userGroupToUserGroups)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(not_found_code))
			})
			It("should respond with not found error if one of subGroup does not exist", func() {
				id := uuid.NewV4()
				for (id == newArtist.Id || id == newRecommendedArtist.Id || id == newLabel.Id || id == newDistributor.Id) {
					id = uuid.NewV4()
				}
				userGroupToUserGroups := &pb.UserGroupToUserGroups{
					UserGroupId: newLabel.Id.String(),
					UserGroups: []*pb.UserGroup{
						&pb.UserGroup{Id: id.String()},
					},
				}
				resp, err := service.RemoveSubGroups(context.Background(), userGroupToUserGroups)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(not_found_code))
			})
		})
	})

	Describe("DeleteUserGroup", func() {
		Context("with valid uuid", func() {
			It("should delete user_group if it exists", func() {
				userGroup := &pb.UserGroup{Id: newArtist.Id.String()}
				_, err := service.DeleteUserGroup(context.Background(), userGroup)

				Expect(err).NotTo(HaveOccurred())
			})
			It("should respond with not_found error if user_group does not exist", func() {
				id := uuid.NewV4()
				for (id == newArtist.Id || id == newRecommendedArtist.Id || id == newLabel.Id || id == newDistributor.Id) {
					id = uuid.NewV4()
				}
				userGroup := &pb.UserGroup{Id: id.String()}
				resp, err := service.DeleteUserGroup(context.Background(), userGroup)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(not_found_code))
			})
		})
		Context("with invalid uuid", func() {
			It("should respond with invalid_argument error", func() {
				id := "45"
				userGroup := &pb.UserGroup{Id: id}
				resp, err := service.DeleteUserGroup(context.Background(), userGroup)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
				Expect(twerr.Meta("argument")).To(Equal("id"))
			})
		})
	})

  Describe("CreateUserGroup", func() {
		Context("with all required attributes", func() {
			It("should create a new user_group", func() {
				avatar := make([]byte, 5)
				tags := make([]*pb.Tag, 1)
				tags[0] = &pb.Tag{Type: "genre", Name: "rock"}
				ownerId := newUser.Id.String()
				userGroup := &pb.UserGroup{
					DisplayName: "group2",
					Avatar: avatar,
					Type: &pb.GroupTaxonomy{Type: "artist"},
					OwnerId: ownerId,
					ShortBio: "short bio",
					Address: &userpb.StreetAddress{Data: map[string]string{"some": "data"}},
					Tags: tags,
				}
				resp, err := service.CreateUserGroup(context.Background(), userGroup)

				Expect(err).NotTo(HaveOccurred())
				Expect(resp.Id).NotTo(Equal(""))
				Expect(resp.DisplayName).To(Equal("group2"))
				Expect(resp.ShortBio).To(Equal("short bio"))
				Expect(resp.Avatar).To(Equal(avatar))
				Expect(resp.Type.Id).To(Equal(newArtistGroupTaxonomy.Id.String()))
				Expect(resp.Type.Type).To(Equal("artist"))
				Expect(resp.OwnerId).To(Equal(ownerId))
				Expect(len(resp.Tags)).To(Equal(1))
				Expect(resp.Tags[0].Id).NotTo(Equal(""))
				Expect(resp.Tags[0].Type).To(Equal("genre"))
				Expect(resp.Tags[0].Name).To(Equal("rock"))
				Expect(resp.Address.Id).NotTo(Equal(""))
				Expect(resp.Privacy.Id).NotTo(Equal(""))
			})

			It("should not create a user_group with same display_name", func() {
				avatar := make([]byte, 5)
				// typeId := newArtistGroupTaxonomy.Id.String()
				ownerId := newUser.Id.String()
				userGroup := &pb.UserGroup{
					DisplayName: "group2",
					Avatar: avatar,
					Address: &userpb.StreetAddress{Data: map[string]string{"some": "data"}},
					Type: &pb.GroupTaxonomy{Type: "artist"},
					OwnerId: ownerId,
				}
				resp, err := service.CreateUserGroup(context.Background(), userGroup)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(already_exists_code))
				Expect(twerr.Msg()).To(Equal("display_name"))
			})
		})

		Context("with missing required attributes", func() {
			It("should not create a user_group without display_name", func() {
				avatar := make([]byte, 5)
				// typeId := newArtistGroupTaxonomy.Id.String()
				ownerId := newUser.Id.String()
				userGroup := &pb.UserGroup{
					DisplayName: "",
					Avatar: avatar,
					Address: &userpb.StreetAddress{Data: map[string]string{"some": "data"}},
					Type: &pb.GroupTaxonomy{Type: "artist"},
					OwnerId: ownerId,
				}
				resp, err := service.CreateUserGroup(context.Background(), userGroup)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
				Expect(twerr.Meta("argument")).To(Equal("display_name"))

			})

			It("should not create a user_group without avatar", func() {
				// typeId := newArtistGroupTaxonomy.Id.String()
				ownerId := newUser.Id.String()
				userGroup := &pb.UserGroup{
					DisplayName: "group3",
					Address: &userpb.StreetAddress{Data: map[string]string{"some": "data"}},
					Type: &pb.GroupTaxonomy{Type: "artist"},
					OwnerId: ownerId,
				}
				resp, err := service.CreateUserGroup(context.Background(), userGroup)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
				Expect(twerr.Meta("argument")).To(Equal("avatar"))
			})

			It("should not create a user without address", func() {
				avatar := make([]byte, 5)
				ownerId := newUser.Id.String()
				userGroup := &pb.UserGroup{
					DisplayName: "group4",
					Address: &userpb.StreetAddress{},
					Avatar: avatar,
					Type: &pb.GroupTaxonomy{Type: "artist"},
					OwnerId: ownerId,
				}
				resp, err := service.CreateUserGroup(context.Background(), userGroup)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
				Expect(twerr.Meta("argument")).To(Equal("address"))
			})

			It("should not create a user without type", func() {
				avatar := make([]byte, 5)
				ownerId := newUser.Id.String()
				userGroup := &pb.UserGroup{
					DisplayName: "group5",
					Address: &userpb.StreetAddress{Data: map[string]string{"some": "data"}},
					Avatar: avatar,
					Type: &pb.GroupTaxonomy{},
					OwnerId: ownerId,
				}
				resp, err := service.CreateUserGroup(context.Background(), userGroup)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
				Expect(twerr.Meta("argument")).To(Equal("type"))
			})

			It("should not create a user without owner", func() {
				avatar := make([]byte, 5)
				userGroup := &pb.UserGroup{
					DisplayName: "group5",
					Address: &userpb.StreetAddress{Data: map[string]string{"some": "data"}},
					Avatar: avatar,
					Type: &pb.GroupTaxonomy{Type: "artist"},
					OwnerId: "",
				}
				resp, err := service.CreateUserGroup(context.Background(), userGroup)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
				Expect(twerr.Meta("argument")).To(Equal("owner"))
			})
		})
  })

	Describe("GetUserGroupTypes", func() {
		It("should respond with group_taxonomies except distributor", func() {
			emptyReq := &userpb.Empty{}
			groupTaxonomies, err := service.GetUserGroupTypes(context.Background(), emptyReq)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(groupTaxonomies.Types)).To(Equal(2))
		})
	})
})
