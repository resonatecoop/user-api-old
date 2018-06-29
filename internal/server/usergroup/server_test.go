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
	// "user-api/internal/database/models"
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
				Expect(resp.Avatar).To(Equal(newArtist.Avatar))
				Expect(resp.Banner).To(Equal(newArtist.Banner))
				Expect(resp.OwnerId).To(Equal(newArtist.OwnerId.String()))
				Expect(resp.TypeId).To(Equal(newArtist.TypeId.String()))
				Expect(resp.GroupEmailAddress).To(Equal(newArtist.GroupEmailAddress))
			})
			It("should respond with not_found error if user_group does not exist", func() {
				id := uuid.NewV4()
				for id == newArtist.Id {
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
				userGroup := &pb.UserGroup{
					Id: newArtist.Id.String(),
					DisplayName: "new display name",
					Description: "new description",
					Avatar: newArtist.Avatar,
					Address: "new address",
					TypeId: newArtist.TypeId.String(),
					OwnerId: newArtist.OwnerId.String(),
					// AdminUsers: internal.ConvertUuidToStrArray(newArtist.AdminUsers),
				}
				_, err := service.UpdateUserGroup(context.Background(), userGroup)

				Expect(err).NotTo(HaveOccurred())
			})
			It("should respond with not_found error if user_group does not exist", func() {
				id := uuid.NewV4()
				for id == newArtist.Id {
					id = uuid.NewV4()
				}
				userGroup := &pb.UserGroup{
					Id: id.String(),
					DisplayName: "new display name",
					Description: "new description",
					Avatar: newArtist.Avatar,
					Address: "new address",
					TypeId: newArtist.TypeId.String(),
					OwnerId: newArtist.OwnerId.String(),
					// AdminUsers: internal.ConvertUuidToStrArray(newArtist.AdminUsers),
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
					Address: "new address",
					TypeId: newArtist.TypeId.String(),
					OwnerId: newArtist.OwnerId.String(),
					// AdminUsers: internal.ConvertUuidToStrArray(newArtist.AdminUsers),
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

  Describe("CreateUserGroup", func() {
		Context("with all required attributes", func() {
			It("should create a new user_group", func() {
				avatar := make([]byte, 5)
				typeId := newArtistGroupTaxonomy.Id.String()
				ownerId := newUser.Id.String()
				userGroup := &pb.UserGroup{
					DisplayName: "group2",
					Avatar: avatar,
					Address: "address",
					TypeId: typeId,
					OwnerId: ownerId,
				}
				resp, err := service.CreateUserGroup(context.Background(), userGroup)

				Expect(err).NotTo(HaveOccurred())
				Expect(resp.Id).NotTo(Equal(""))
				Expect(resp.DisplayName).To(Equal("group2"))
				Expect(resp.Avatar).To(Equal(avatar))
				Expect(resp.TypeId).To(Equal(typeId))
				Expect(resp.OwnerId).To(Equal(ownerId))
				// Expect(resp.Address).To(Equal("address"))
				// Expect(resp.AdminUsers).To(Equal())
			})

			It("should not create a user_group with same display_name", func() {
				avatar := make([]byte, 5)
				typeId := newArtistGroupTaxonomy.Id.String()
				ownerId := newUser.Id.String()
				userGroup := &pb.UserGroup{
					DisplayName: "group2",
					Avatar: avatar,
					Address: "address",
					TypeId: typeId,
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
				typeId := newArtistGroupTaxonomy.Id.String()
				ownerId := newUser.Id.String()
				userGroup := &pb.UserGroup{
					DisplayName: "",
					Avatar: avatar,
					Address: "address",
					TypeId: typeId,
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
				typeId := newArtistGroupTaxonomy.Id.String()
				ownerId := newUser.Id.String()
				userGroup := &pb.UserGroup{
					DisplayName: "group3",
					Address: "address",
					TypeId: typeId,
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
				typeId := newArtistGroupTaxonomy.Id.String()
				ownerId := newUser.Id.String()
				userGroup := &pb.UserGroup{
					DisplayName: "group4",
					Address: "",
					Avatar: avatar,
					TypeId: typeId,
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
					Address: "address",
					Avatar: avatar,
					TypeId: "",
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
				typeId := newArtistGroupTaxonomy.Id.String()
				userGroup := &pb.UserGroup{
					DisplayName: "group5",
					Address: "address",
					Avatar: avatar,
					TypeId: typeId,
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

	Describe("DeleteUserGroup", func() {
		Context("with valid uuid", func() {
			It("should delete user_group if it exists", func() {
				userGroup := &pb.UserGroup{Id: newArtist.Id.String()}
				_, err := service.DeleteUserGroup(context.Background(), userGroup)

				Expect(err).NotTo(HaveOccurred())
			})
			It("should respond with not_found error if user_group does not exist", func() {
				id := uuid.NewV4()
				for id == newArtist.Id {
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
})
