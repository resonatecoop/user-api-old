package userserver_test

import (
	// "fmt"
	// "reflect"
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/twitchtv/twirp"
	"github.com/satori/go.uuid"

	pb "user-api/rpc/user"
	"user-api/internal/database/models"
)

var _ = Describe("User server", func() {
	const already_exists_code twirp.ErrorCode = "already_exists"
	const invalid_argument_code twirp.ErrorCode = "invalid_argument"
	const not_found_code twirp.ErrorCode = "not_found"

	Describe("GetUser", func() {
		Context("with valid uuid", func() {
			It("should respond with user if it exists", func() {
				user := &pb.User{Id: newUser.Id.String()}
				resp, err := service.GetUser(context.Background(), user)

				Expect(err).NotTo(HaveOccurred())
				Expect(resp.Id).To(Equal(newUser.Id.String()))
				Expect(resp.Username).To(Equal(newUser.Username))
				Expect(resp.FullName).To(Equal(newUser.FullName))
				Expect(resp.DisplayName).To(Equal(newUser.DisplayName))
				Expect(resp.Email).To(Equal(newUser.Email))
			})
			It("should respond with not_found error if user does not exist", func() {
				id := uuid.NewV4()
				for id == newUser.Id {
					id = uuid.NewV4()
				}
				user := &pb.User{Id: id.String()}
				resp, err := service.GetUser(context.Background(), user)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(not_found_code))
			})
		})
		Context("with invalid uuid", func() {
			It("should respond with invalid_argument error", func() {
				id := "45"
				user := &pb.User{Id: id}
				resp, err := service.GetUser(context.Background(), user)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
				Expect(twerr.Meta("argument")).To(Equal("id"))
			})
		})
	})

	Describe("AddFavoriteTrack", func() {
		Context("with user_id and track_id", func() {
			It("should add favorite track", func() {
				userToTrack := &pb.UserToTrack{UserId: newUser.Id.String(), TrackId: newTrack.Id.String()}
				_, err := service.AddFavoriteTrack(context.Background(), userToTrack)

				Expect(err).NotTo(HaveOccurred())

				user := new(models.User)
				err = db.Model(user).Where("id = ?", newUser.Id).Select()
				Expect(err).NotTo(HaveOccurred())
				Expect(user.FavoriteTracks).To(HaveLen(1))
				Expect(user.FavoriteTracks).To(ContainElement(newTrack.Id))

				track := new(models.Track)
				err = db.Model(track).Where("id = ?", newTrack.Id).Select()
				Expect(err).NotTo(HaveOccurred())
				Expect(track.FavoriteOfUsers).To(HaveLen(1))
				Expect(track.FavoriteOfUsers).To(ContainElement(newUser.Id))
			})
			It("should respond with not_found error if user does not exist", func() {
				userId := uuid.NewV4()
				for userId == newUser.Id {
					userId = uuid.NewV4()
				}
				userToTrack := &pb.UserToTrack{UserId: userId.String(), TrackId: newTrack.Id.String()}
				resp, err := service.AddFavoriteTrack(context.Background(), userToTrack)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(not_found_code))
			})
			It("should respond with not_found error if track does not exist", func() {
				trackId := uuid.NewV4()
				for trackId == newTrack.Id {
					trackId = uuid.NewV4()
				}
				userToTrack := &pb.UserToTrack{UserId: newUser.Id.String(), TrackId: trackId.String()}
				resp, err := service.AddFavoriteTrack(context.Background(), userToTrack)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(not_found_code))
			})
		})
		Context("with invalid track_id", func() {
			It("should respond with invalid_argument_code error", func() {
				userToTrack := &pb.UserToTrack{UserId: newUser.Id.String(), TrackId: ""}
				resp, err := service.AddFavoriteTrack(context.Background(), userToTrack)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
			})
		})
		Context("with invalid user_id", func() {
			It("should respond with invalid_argument_code error", func() {
				userToTrack := &pb.UserToTrack{UserId: "", TrackId: newTrack.Id.String()}
				resp, err := service.AddFavoriteTrack(context.Background(), userToTrack)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
			})
		})
	})

	Describe("RemoveFavoriteTrack", func() {
		Context("with user_id and track_id", func() {
			It("should remove favorite track", func() {
				userToTrack := &pb.UserToTrack{UserId: newUser.Id.String(), TrackId: newTrack.Id.String()}
				_, err := service.RemoveFavoriteTrack(context.Background(), userToTrack)

				Expect(err).NotTo(HaveOccurred())

				user := new(models.User)
				err = db.Model(user).Where("id = ?", newUser.Id).Select()
				Expect(err).NotTo(HaveOccurred())
				Expect(user.FavoriteTracks).To(HaveLen(0))
				Expect(user.FavoriteTracks).NotTo(ContainElement(newTrack.Id))

				track := new(models.Track)
				err = db.Model(track).Where("id = ?", newTrack.Id).Select()
				Expect(err).NotTo(HaveOccurred())
				Expect(track.FavoriteOfUsers).To(HaveLen(0))
				Expect(track.FavoriteOfUsers).NotTo(ContainElement(newUser.Id))
			})
			It("should respond with not_found error if user does not exist", func() {
				userId := uuid.NewV4()
				for userId == newUser.Id {
					userId = uuid.NewV4()
				}
				userToTrack := &pb.UserToTrack{UserId: userId.String(), TrackId: newTrack.Id.String()}
				resp, err := service.RemoveFavoriteTrack(context.Background(), userToTrack)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(not_found_code))
			})
			It("should respond with not_found error if track does not exist", func() {
				trackId := uuid.NewV4()
				for trackId == newTrack.Id {
					trackId = uuid.NewV4()
				}
				userToTrack := &pb.UserToTrack{UserId: newUser.Id.String(), TrackId: trackId.String()}
				resp, err := service.RemoveFavoriteTrack(context.Background(), userToTrack)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(not_found_code))
			})
		})
		Context("with invalid track_id", func() {
			It("should respond with invalid_argument_code error", func() {
				userToTrack := &pb.UserToTrack{UserId: newUser.Id.String(), TrackId: ""}
				resp, err := service.RemoveFavoriteTrack(context.Background(), userToTrack)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
			})
		})
		Context("with invalid user_id", func() {
			It("should respond with invalid_argument_code error", func() {
				userToTrack := &pb.UserToTrack{UserId: "", TrackId: newTrack.Id.String()}
				resp, err := service.RemoveFavoriteTrack(context.Background(), userToTrack)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
			})
		})
	})

	Describe("FollowGroup", func() {
		Context("with user_id and user_group_id", func() {
			It("should add followed group", func() {
				userToUserGroup := &pb.UserToUserGroup{UserId: newUser.Id.String(), UserGroupId: newUserGroup.Id.String()}
				_, err := service.FollowGroup(context.Background(), userToUserGroup)

				Expect(err).NotTo(HaveOccurred())

				user := new(models.User)
				err = db.Model(user).Where("id = ?", newUser.Id).Select()
				Expect(err).NotTo(HaveOccurred())
				Expect(user.FollowedGroups).To(HaveLen(1))
				Expect(user.FollowedGroups).To(ContainElement(newUserGroup.Id))

				userGroup := new(models.UserGroup)
				err = db.Model(userGroup).Where("id = ?", newUserGroup.Id).Select()
				Expect(err).NotTo(HaveOccurred())
				Expect(userGroup.Followers).To(HaveLen(1))
				Expect(userGroup.Followers).To(ContainElement(newUser.Id))
			})
			It("should respond with not_found error if user does not exist", func() {
				userId := uuid.NewV4()
				for userId == newUser.Id {
					userId = uuid.NewV4()
				}
				userToUserGroup := &pb.UserToUserGroup{UserId: userId.String(), UserGroupId: newUserGroup.Id.String()}
				resp, err := service.FollowGroup(context.Background(), userToUserGroup)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(not_found_code))
			})
			It("should respond with not_found error if user_group does not exist", func() {
				userGroupId := uuid.NewV4()
				for userGroupId == newUserGroup.Id {
					userGroupId = uuid.NewV4()
				}
				userToUserGroup := &pb.UserToUserGroup{UserId: newUser.Id.String(), UserGroupId: userGroupId.String()}
				resp, err := service.FollowGroup(context.Background(), userToUserGroup)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(not_found_code))
			})
		})
		Context("with invalid user_group_id", func() {
			It("should respond with invalid_argument_code error", func() {
				userToUserGroup := &pb.UserToUserGroup{UserId: newUser.Id.String(), UserGroupId: ""}
				resp, err := service.FollowGroup(context.Background(), userToUserGroup)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
			})
		})
		Context("with invalid user_id", func() {
			It("should respond with invalid_argument_code error", func() {
				userToUserGroup := &pb.UserToUserGroup{UserId: "", UserGroupId: newUserGroup.Id.String()}
				resp, err := service.FollowGroup(context.Background(), userToUserGroup)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
			})
		})
	})

	Describe("UnfollowGroup", func() {
		Context("with user_id and user_group_id", func() {
			It("should remove followed group", func() {
				userToUserGroup := &pb.UserToUserGroup{UserId: newUser.Id.String(), UserGroupId: newUserGroup.Id.String()}
				_, err := service.UnfollowGroup(context.Background(), userToUserGroup)

				Expect(err).NotTo(HaveOccurred())

				user := new(models.User)
				err = db.Model(user).Where("id = ?", newUser.Id).Select()
				Expect(err).NotTo(HaveOccurred())
				Expect(user.FollowedGroups).To(HaveLen(0))
				Expect(user.FollowedGroups).NotTo(ContainElement(newUserGroup.Id))

				userGroup := new(models.UserGroup)
				err = db.Model(userGroup).Where("id = ?", newUserGroup.Id).Select()
				Expect(err).NotTo(HaveOccurred())
				Expect(userGroup.Followers).To(HaveLen(0))
				Expect(userGroup.Followers).NotTo(ContainElement(newUser.Id))
			})
			It("should respond with not_found error if user does not exist", func() {
				userId := uuid.NewV4()
				for userId == newUser.Id {
					userId = uuid.NewV4()
				}
				userToUserGroup := &pb.UserToUserGroup{UserId: userId.String(), UserGroupId: newUserGroup.Id.String()}
				resp, err := service.UnfollowGroup(context.Background(), userToUserGroup)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(not_found_code))
			})
			It("should respond with not_found error if track does not exist", func() {
				userGroupId := uuid.NewV4()
				for userGroupId == newUserGroup.Id {
					userGroupId = uuid.NewV4()
				}
				userToUserGroup := &pb.UserToUserGroup{UserId: newUser.Id.String(), UserGroupId: userGroupId.String()}
				resp, err := service.UnfollowGroup(context.Background(), userToUserGroup)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(not_found_code))
			})
		})
		Context("with invalid user_group_id", func() {
			It("should respond with invalid_argument_code error", func() {
				userToUserGroup := &pb.UserToUserGroup{UserId: newUser.Id.String(), UserGroupId: ""}
				resp, err := service.UnfollowGroup(context.Background(), userToUserGroup)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
			})
		})
		Context("with invalid user_id", func() {
			It("should respond with invalid_argument_code error", func() {
				userToUserGroup := &pb.UserToUserGroup{UserId: "", UserGroupId: newUserGroup.Id.String()}
				resp, err := service.UnfollowGroup(context.Background(), userToUserGroup)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
			})
		})
	})

	Describe("UpdateUser", func() {
		Context("with valid uuid", func() {
			It("should update user if it exists", func() {
				user := &pb.User{
					Id: newUser.Id.String(),
					Username: "new username",
					FullName: "full name",
					DisplayName: "display name",
					Email: "email@fake.com",
					FirstName: "first name",
				}
				_, err := service.UpdateUser(context.Background(), user)

				Expect(err).NotTo(HaveOccurred())
			})
			It("should respond with not_found error if user does not exist", func() {
				id := uuid.NewV4()
				for id == newUser.Id {
					id = uuid.NewV4()
				}
				user := &pb.User{
					Id: id.String(),
					Username: "username",
					FullName: "fullname",
					DisplayName: "displayname",
					Email: "email@fake.comm",
					FirstName: "firstname",
				}
				resp, err := service.UpdateUser(context.Background(), user)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(not_found_code))
			})
		})
		Context("with invalid uuid", func() {
			It("should respond with invalid_argument error", func() {
				id := "45"
				user := &pb.User{
					Id: id,
					Username: "new username",
					FullName: "full name",
					DisplayName: "display name",
					Email: "email@fake.com",
					FirstName: "first name",
				}
				resp, err := service.UpdateUser(context.Background(), user)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
				Expect(twerr.Meta("argument")).To(Equal("id"))
			})
		})
	})

	Describe("CreateUser", func() {
		Context("with all required attributes", func() {
			It("should create a new user", func() {
				user := &pb.User{Username: "janed", FullName: "jane d", DisplayName: "jad", Email: "jane@d.com"}
				resp, err := service.CreateUser(context.Background(), user)

				Expect(err).NotTo(HaveOccurred())
				Expect(resp.Username).To(Equal("janed"))
				Expect(resp.FullName).To(Equal("jane d"))
				Expect(resp.DisplayName).To(Equal("jad"))
				Expect(resp.Email).To(Equal("jane@d.com"))
				Expect(resp.Id).NotTo(Equal(""))
			})

			It("should not create a user with same display_name", func() {
				user := &pb.User{Username: "janedoe", FullName: "jane doe", DisplayName: "jad", Email: "jane@doe.com"}
				resp, err := service.CreateUser(context.Background(), user)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(already_exists_code))
				Expect(twerr.Msg()).To(Equal("display_name"))
			})

			It("should not create a user with same email", func() {
				user := &pb.User{Username: "janedoe", FullName: "jane doe", DisplayName: "jadoe", Email: "jane@d.com"}
				resp, err := service.CreateUser(context.Background(), user)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(already_exists_code))
				Expect(twerr.Msg()).To(Equal("email"))
			})

			It("should not create a user with same username", func() {
				user := &pb.User{Username: "janed", FullName: "jane doe", DisplayName: "jadoe", Email: "jane@doe.com"}
				resp, err := service.CreateUser(context.Background(), user)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(already_exists_code))
				Expect(twerr.Msg()).To(Equal("username"))
			})
		})

		Context("with missing required attributes", func() {
			It("should not create a user without display_name", func() {
				user := &pb.User{Username: "johnd", FullName: "john doe", DisplayName: "", Email: "john@doe.com"}
				resp, err := service.CreateUser(context.Background(), user)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
				Expect(twerr.Meta("argument")).To(Equal("display_name"))

			})

			It("should not create a user without email", func() {
				user := &pb.User{Username: "johnd", FullName: "john doe", DisplayName: "john", Email: ""}
				resp, err := service.CreateUser(context.Background(), user)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
				Expect(twerr.Meta("argument")).To(Equal("email"))
			})

			It("should not create a user without username", func() {
				user := &pb.User{Username: "", FullName: "john doe", DisplayName: "john", Email: "john@doe.com"}
				resp, err := service.CreateUser(context.Background(), user)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
				Expect(twerr.Meta("argument")).To(Equal("username"))
			})

			It("should not create a user without full_name", func() {
				user := &pb.User{Username: "johnd", FullName: "", DisplayName: "john", Email: "john@doe.com"}
				resp, err := service.CreateUser(context.Background(), user)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
				Expect(twerr.Meta("argument")).To(Equal("full_name"))
			})
		})
	})

	Describe("DeleteUser", func() {
		Context("with valid uuid", func() {
			It("should delete user if it exists", func() {
				user := &pb.User{Id: newUser.Id.String()}
				_, err := service.DeleteUser(context.Background(), user)

				Expect(err).NotTo(HaveOccurred())
			})
			It("should respond with not_found error if user does not exist", func() {
				id := uuid.NewV4()
				for id == newUser.Id {
					id = uuid.NewV4()
				}
				user := &pb.User{Id: id.String()}
				resp, err := service.DeleteUser(context.Background(), user)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(not_found_code))
			})
		})
		Context("with invalid uuid", func() {
			It("should respond with invalid_argument error", func() {
				id := "45"
				user := &pb.User{Id: id}
				resp, err := service.DeleteUser(context.Background(), user)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
				Expect(twerr.Meta("argument")).To(Equal("id"))
			})
		})
	})
})
