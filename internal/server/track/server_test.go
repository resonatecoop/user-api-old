package trackserver_test

import (
	// "fmt"
	// "reflect"
	"context"

	// "github.com/go-pg/pg"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/twitchtv/twirp"
	"github.com/satori/go.uuid"

	pb "user-api/rpc/track"
	usergrouppb "user-api/rpc/usergroup"
	"user-api/internal/database/models"
)

var _ = Describe("Track server", func() {
	const already_exists_code twirp.ErrorCode = "already_exists"
	const invalid_argument_code twirp.ErrorCode = "invalid_argument"
	const not_found_code twirp.ErrorCode = "not_found"

	Describe("GetTrack", func() {
		Context("with valid uuid", func() {
			It("should respond with track if it exists", func() {
				track := &pb.Track{Id: newTrack.Id.String()}

				res, err := service.GetTrack(context.Background(), track)
				Expect(err).NotTo(HaveOccurred())

				Expect(res.Title).To(Equal(newTrack.Title))
				Expect(res.Status).To(Equal(newTrack.Status))
				Expect(res.Duration).To(Equal(newTrack.Duration))
				Expect(res.Enabled).To(Equal(newTrack.Enabled))
				Expect(res.TrackNumber).To(Equal(newTrack.TrackNumber))
				Expect(res.CreatorId).To(Equal(newTrack.CreatorId.String()))
				Expect(res.UserGroupId).To(Equal(newTrack.UserGroupId.String()))
				Expect(res.TrackServerId).To(Equal(newTrack.TrackServerId.String()))

				Expect(len(res.Tags)).To(Equal(1))
				Expect(res.Tags[0].Id).To(Equal(newGenreTag.Id.String()))
				Expect(res.Tags[0].Type).To(Equal(newGenreTag.Type))
				Expect(res.Tags[0].Name).To(Equal(newGenreTag.Name))

				Expect(len(res.Artists)).To(Equal(1))
				Expect(res.Artists[0].Id).To(Equal(newArtistUserGroup.Id.String()))
				Expect(res.Artists[0].Avatar).To(Equal(newArtistUserGroup.Avatar))
				Expect(res.Artists[0].DisplayName).To(Equal(newArtistUserGroup.DisplayName))

				Expect(len(res.TrackGroups)).To(Equal(1))
				Expect(res.TrackGroups[0].Id).To(Equal(newAlbum.Id.String()))
				Expect(res.TrackGroups[0].Cover).To(Equal(newAlbum.Cover))
				Expect(res.TrackGroups[0].Title).To(Equal(newAlbum.Title))
			})
			It("should respond with not_found error if track does not exist", func() {
				id := uuid.NewV4()
				for id == newTrack.Id {
					id = uuid.NewV4()
				}
				track := &pb.Track{Id: id.String()}
				resp, err := service.GetTrack(context.Background(), track)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(not_found_code))
			})
		})
		Context("with invalid uuid", func() {
			It("should respond with invalid_argument error", func() {
				id := "45"
				track := &pb.Track{Id: id}
				resp, err := service.GetTrack(context.Background(), track)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
				Expect(twerr.Meta("argument")).To(Equal("id"))
			})
		})
	})

	Describe("UpdateTrack", func() {
		Context("with valid uuid", func() {
			It("should update track if it exists", func() {
				trackServerId := uuid.NewV4()
				track := &pb.Track{
					Id: newTrack.Id.String(),
					Title: "best track ever",
					TrackNumber: 1,
					Status: "paid",
					Duration: 250.12,
					CreatorId: newUser.Id.String(),
					UserGroupId: newArtistUserGroup.Id.String(),
					TrackServerId: trackServerId.String(),
				}
				_, err := service.UpdateTrack(context.Background(), track)

				Expect(err).NotTo(HaveOccurred())

				t := new(models.Track)
				err = db.Model(t).Where("id = ?", newTrack.Id).Select()
				Expect(err).NotTo(HaveOccurred())
				Expect(t.Title).To(Equal(track.Title))
				Expect(t.TrackNumber).To(Equal(track.TrackNumber))
				Expect(t.Status).To(Equal(track.Status))
				Expect(t.Duration).To(Equal(track.Duration))
				Expect(t.TrackServerId.String()).To(Equal(track.TrackServerId))

				// unchanged
				Expect(t.CreatorId.String()).To(Equal(track.CreatorId))
				Expect(t.UserGroupId.String()).To(Equal(track.UserGroupId))
				Expect(len(t.Tags)).To(Equal(1))
				Expect(len(t.Artists)).To(Equal(1))
				Expect(len(t.TrackGroups)).To(Equal(2))
			})
			It("should respond with not_found error if track does not exist", func() {
				id := uuid.NewV4()
				for id == newTrack.Id {
					id = uuid.NewV4()
				}
				track := &pb.Track{
					Id: id.String(),
					Title: "best track ever",
					TrackNumber: 1,
					Status: "paid",
					Duration: 250.12,
					CreatorId: newUser.Id.String(),
					UserGroupId: newArtistUserGroup.Id.String(),
				}
				resp, err := service.UpdateTrack(context.Background(), track)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(not_found_code))
			})
		})
		Context("with invalid uuid", func() {
			It("should respond with invalid_argument error", func() {
				id := "45"
				track := &pb.Track{
					Id: id,
					Title: "best track ever",
					TrackNumber: 1,
					Status: "paid",
					Duration: 250.12,
					CreatorId: newUser.Id.String(),
					UserGroupId: newArtistUserGroup.Id.String(),
				}
				resp, err := service.UpdateTrack(context.Background(), track)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
				Expect(twerr.Meta("argument")).To(Equal("id"))
			})
		})
	})

	Describe("CreateTrack", func() {
		Context("with all required attributes", func() {
			It("should create a new track", func() {
				track := &pb.Track{
					Title: "best track ever",
					TrackNumber: 1,
					Status: "paid",
					CreatorId: newUser.Id.String(),
					UserGroupId: newArtistUserGroup.Id.String(),
					Artists: []*usergrouppb.UserGroup{
						&usergrouppb.UserGroup{
							Id: newArtistUserGroup.Id.String(),
						},
					},
					Tags: []*usergrouppb.Tag{
						&usergrouppb.Tag{
							Type: "genre",
							Name: "rock",
						},
					},
				}
				resp, err := service.CreateTrack(context.Background(), track)

				Expect(err).NotTo(HaveOccurred())
				Expect(resp).NotTo(BeNil())

				Expect(resp.Id).NotTo(BeNil())
				Expect(resp.Title).To(Equal(track.Title))
				Expect(resp.TrackNumber).To(Equal(track.TrackNumber))
				Expect(resp.Status).To(Equal(track.Status))
				Expect(resp.CreatorId).To(Equal(track.CreatorId))
				Expect(resp.UserGroupId).To(Equal(track.UserGroupId))
				Expect(len(resp.Tags)).To(Equal(1))
				Expect(resp.Tags[0].Id).NotTo(Equal(""))
				Expect(resp.Tags[0].Type).To(Equal("genre"))
				Expect(resp.Tags[0].Name).To(Equal("rock"))
				Expect(len(resp.Artists)).To(Equal(1))
				Expect(resp.Artists[0].Id).To(Equal(newArtistUserGroup.Id.String()))
				Expect(resp.Artists[0].DisplayName).To(Equal(newArtistUserGroup.DisplayName))
				Expect(resp.Artists[0].Avatar).To(Equal(newArtistUserGroup.Avatar))

				artist := new(models.UserGroup)
				err = db.Model(artist).Where("id = ?", newArtistUserGroup.Id).Select()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(artist.Tracks)).To(Equal(2))

				trackId, err := uuid.FromString(resp.Id)
				Expect(err).NotTo(HaveOccurred())
				Expect(artist.Tracks).To(ContainElement(trackId))
			})
		})

		Context("with missing required attributes", func() {
			It("should not create a track without title", func() {
				track := &pb.Track{
					Title: "",
					TrackNumber: 1,
					Status: "paid",
					CreatorId: newUser.Id.String(),
					UserGroupId: newArtistUserGroup.Id.String(),
					Artists: []*usergrouppb.UserGroup{
						&usergrouppb.UserGroup{
							Id: newArtistUserGroup.Id.String(),
						},
					},
				}
				resp, err := service.CreateTrack(context.Background(), track)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
				Expect(twerr.Meta("argument")).To(Equal("title"))
			})
			It("should not create a track without status", func() {
				track := &pb.Track{
					Title: "title",
					TrackNumber: 1,
					Status: "",
					CreatorId: newUser.Id.String(),
					UserGroupId: newArtistUserGroup.Id.String(),
					Artists: []*usergrouppb.UserGroup{
						&usergrouppb.UserGroup{
							Id: newArtistUserGroup.Id.String(),
						},
					},
				}
				resp, err := service.CreateTrack(context.Background(), track)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())
				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
				Expect(twerr.Meta("argument")).To(Equal("status"))
			})
			It("should not create a track without creator_id", func() {
				track := &pb.Track{
					Title: "title",
					TrackNumber: 1,
					Status: "paid",
					CreatorId: "",
					UserGroupId: newArtistUserGroup.Id.String(),
					Artists: []*usergrouppb.UserGroup{
						&usergrouppb.UserGroup{
							Id: newArtistUserGroup.Id.String(),
						},
					},
				}
				resp, err := service.CreateTrack(context.Background(), track)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
				Expect(twerr.Meta("argument")).To(Equal("creator_id"))
			})
			It("should not create a track without user_group_id", func() {
				track := &pb.Track{
					Title: "title",
					TrackNumber: 1,
					Status: "paid",
					CreatorId: newUser.Id.String(),
					UserGroupId: "",
					Artists: []*usergrouppb.UserGroup{
						&usergrouppb.UserGroup{
							Id: newArtistUserGroup.Id.String(),
						},
					},
				}
				resp, err := service.CreateTrack(context.Background(), track)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())
				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
				Expect(twerr.Meta("argument")).To(Equal("user_group_id"))
			})
			It("should not create a track without track_number", func() {
				track := &pb.Track{
					Title: "title",
					TrackNumber: 0,
					Status: "free",
					CreatorId: newUser.Id.String(),
					UserGroupId: newArtistUserGroup.Id.String(),
					Artists: []*usergrouppb.UserGroup{
						&usergrouppb.UserGroup{
							Id: newArtistUserGroup.Id.String(),
						},
					},
				}
				resp, err := service.CreateTrack(context.Background(), track)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())
				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
				Expect(twerr.Meta("argument")).To(Equal("track_number"))
			})
		})

		Context("with invalid attributes", func() {
			It("should not create a track if creator_id is invalid", func() {
				track := &pb.Track{
					Title: "title",
					TrackNumber: 1,
					Status: "free",
					CreatorId: "12a",
					UserGroupId: newArtistUserGroup.Id.String(),
					Artists: []*usergrouppb.UserGroup{
						&usergrouppb.UserGroup{
							Id: newArtistUserGroup.Id.String(),
						},
					},
				}
				resp, err := service.CreateTrack(context.Background(), track)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())
				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
			})
			It("should not create a track if user_group_id is invalid", func() {
				track := &pb.Track{
					Title: "title",
					TrackNumber: 1,
					Status: "free",
					CreatorId: newUser.Id.String(),
					UserGroupId: "abc1",
					Artists: []*usergrouppb.UserGroup{
						&usergrouppb.UserGroup{
							Id: newArtistUserGroup.Id.String(),
						},
					},
				}
				resp, err := service.CreateTrack(context.Background(), track)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())
				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
			})
			It("should not create a track if creator does not exist", func() {
				userId := uuid.NewV4()
				for userId == newUser.Id {
					userId = uuid.NewV4()
				}
				track := &pb.Track{
					Title: "title",
					TrackNumber: 1,
					Status: "free",
					CreatorId: userId.String(),
					UserGroupId: newArtistUserGroup.Id.String(),
					Artists: []*usergrouppb.UserGroup{
						&usergrouppb.UserGroup{
							Id: newArtistUserGroup.Id.String(),
						},
					},
				}
				resp, err := service.CreateTrack(context.Background(), track)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())
				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(not_found_code))
			})
			It("should not create a track if user_group does not exist", func() {
				userGroupId := uuid.NewV4()
				for userGroupId == newLabelUserGroup.Id || userGroupId == newArtistUserGroup.Id {
					userGroupId = uuid.NewV4()
				}
				track := &pb.Track{
					Title: "title",
					TrackNumber: 1,
					Status: "free",
					CreatorId: newUser.Id.String(),
					UserGroupId: userGroupId.String(),
					Artists: []*usergrouppb.UserGroup{
						&usergrouppb.UserGroup{
							Id: newArtistUserGroup.Id.String(),
						},
					},
				}
				resp, err := service.CreateTrack(context.Background(), track)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())
				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(not_found_code))
			})
			It("should not create a track if one of the artists does not exist", func() {
				userGroupId := uuid.NewV4()
				for userGroupId == newLabelUserGroup.Id || userGroupId == newArtistUserGroup.Id {
					userGroupId = uuid.NewV4()
				}
				track := &pb.Track{
					Title: "title",
					TrackNumber: 1,
					Status: "free",
					CreatorId: newUser.Id.String(),
					UserGroupId: userGroupId.String(),
					Artists: []*usergrouppb.UserGroup{
						&usergrouppb.UserGroup{
							Id: userGroupId.String(),
						},
					},
				}
				resp, err := service.CreateTrack(context.Background(), track)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())
				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(not_found_code))
			})
		})
	})

	XDescribe("DeleteTrack", func() {
		Context("with valid uuid", func() {
			It("should delete track if it exists", func() {
				track := &pb.Track{Id: newTrack.Id.String()}

				trackToDelete := new(models.Track)
				err := db.Model(trackToDelete).Where("id = ?", newTrack.Id).Select()
				Expect(err).NotTo(HaveOccurred())

				_, err = service.DeleteTrack(context.Background(), track)

				Expect(err).NotTo(HaveOccurred())
			})
			It("should respond with not_found error if user does not exist", func() {
				id := uuid.NewV4()
				for id == newTrack.Id {
					id = uuid.NewV4()
				}
				track := &pb.Track{Id: id.String()}
				resp, err := service.DeleteTrack(context.Background(), track)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(not_found_code))
			})
		})
		Context("with invalid uuid", func() {
			It("should respond with invalid_argument error", func() {
				id := "45"
				track := &pb.Track{Id: id}
				resp, err := service.DeleteTrack(context.Background(), track)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
				Expect(twerr.Meta("argument")).To(Equal("id"))
			})
		})
	})
})