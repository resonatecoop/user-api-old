package usergroupserver_test

import (
	// "fmt"
	// "reflect"
	"context"
	"net/url"

	"github.com/go-pg/pg"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/twitchtv/twirp"
	"github.com/satori/go.uuid"

	pb "user-api/rpc/usergroup"
	tagpb "user-api/rpc/tag"
	userpb "user-api/rpc/user"
	// trackpb "user-api/rpc/track"
	// "user-api/internal"
	"user-api/internal/database/models"
)

var _ = Describe("UserGroup server", func() {
  const already_exists_code twirp.ErrorCode = "already_exists"
  const invalid_argument_code twirp.ErrorCode = "invalid_argument"
  const not_found_code twirp.ErrorCode = "not_found"

	Describe("GetUserGroup", func() {
		Context("with valid uuid", func() {
			It("should respond with user_group (artist) if it exists", func() {
				userGroup := &pb.UserGroup{Id: newArtist.Id.String()}
				resp, err := service.GetUserGroup(context.Background(), userGroup)

				Expect(err).NotTo(HaveOccurred())
				Expect(resp.Id).To(Equal(newArtist.Id.String()))
				Expect(resp.DisplayName).To(Equal(newArtist.DisplayName))
				Expect(resp.Description).To(Equal(newArtist.Description))
				Expect(resp.ShortBio).To(Equal(newArtist.ShortBio))
				Expect(resp.Avatar).To(Equal(newArtist.Avatar))
				Expect(resp.Banner).To(Equal(newArtist.Banner))
				Expect(resp.GroupEmailAddress).To(Equal(newArtist.GroupEmailAddress))
				Expect(resp.Publisher).To(Equal(newArtist.Publisher))
				Expect(resp.Pro).To(Equal(newArtist.Pro))
				Expect(resp.OwnerId).To(Equal(newArtist.OwnerId.String()))
				Expect(resp.Type.Id).To(Equal(newArtistGroupTaxonomy.Id.String()))
				Expect(resp.Type.Type).To(Equal("artist"))

				Expect(resp.Address.Id).To(Equal(artistAddress.Id.String()))

				Expect(len(resp.Tags)).To(Equal(1))
				Expect(resp.Tags[0].Id).To(Equal(newGenreTag.Id.String()))
				Expect(resp.Tags[0].Type).To(Equal(newGenreTag.Type))
				Expect(resp.Tags[0].Name).To(Equal(newGenreTag.Name))

				Expect(len(resp.Links)).To(Equal(1))
				Expect(resp.Links[0].Id).To(Equal(newLink.Id.String()))
				Expect(resp.Links[0].Uri).To(Equal(newLink.Uri))
				Expect(resp.Links[0].Platform).To(Equal(newLink.Platform))

				Expect(len(resp.Members)).To(Equal(1))
				Expect(resp.Members[0].Id).To(Equal(newUserProfile.Id.String()))
				Expect(resp.Members[0].DisplayName).To(Equal(newArtistUserGroupMember.DisplayName))
				Expect(resp.Members[0].Avatar).To(Equal(newUserProfile.Avatar))
				Expect(len(resp.Members[0].Tags)).To(Equal(1))
				Expect(resp.Members[0].Tags[0].Id).To(Equal(newRoleTag.Id.String()))
				Expect(resp.Members[0].Tags[0].Type).To(Equal(newRoleTag.Type))
				Expect(resp.Members[0].Tags[0].Name).To(Equal(newRoleTag.Name))

				Expect(len(resp.MemberOfGroups)).To(Equal(1))
				Expect(resp.MemberOfGroups[0].Id).To(Equal(newLabel.Id.String()))
				Expect(resp.MemberOfGroups[0].DisplayName).To(Equal(newArtist.DisplayName))
				Expect(resp.MemberOfGroups[0].Avatar).To(Equal(newArtist.Avatar))

				Expect(len(resp.TrackGroups)).To(Equal(1))
				Expect(resp.TrackGroups[0].Id).To(Equal(newAlbum.Id.String()))
				Expect(resp.TrackGroups[0].Title).To(Equal(newAlbum.Title))
				Expect(resp.TrackGroups[0].Cover).To(Equal(newAlbum.Cover))
				Expect(resp.TrackGroups[0].Type).To(Equal(newAlbum.Type))
				Expect(resp.TrackGroups[0].About).To(Equal(newAlbum.About))
				Expect(resp.TrackGroups[0].DisplayArtist).To(Equal(newAlbum.DisplayArtist))
				Expect(resp.TrackGroups[0].TotalTracks).To(Equal(int32(1)))

				Expect(resp.FeaturedTrackGroup.Id).To(Equal(newAlbum.Id.String()))
				Expect(resp.FeaturedTrackGroup.Title).To(Equal(newAlbum.Title))
				Expect(resp.FeaturedTrackGroup.Cover).To(Equal(newAlbum.Cover))
				Expect(resp.FeaturedTrackGroup.Type).To(Equal(newAlbum.Type))
				Expect(resp.FeaturedTrackGroup.About).To(Equal(newAlbum.About))
				Expect(resp.FeaturedTrackGroup.DisplayArtist).To(Equal(newAlbum.DisplayArtist))
				Expect(resp.FeaturedTrackGroup.TotalTracks).To(Equal(int32(1)))

				Expect(len(resp.HighlightedTracks)).To(Equal(1))
				Expect(resp.HighlightedTracks[0].Id).To(Equal(newTrack.Id.String()))
				Expect(resp.HighlightedTracks[0].Title).To(Equal(newTrack.Title))
				Expect(resp.HighlightedTracks[0].TrackServerId).To(Equal(newTrack.TrackServerId.String()))
				Expect(resp.HighlightedTracks[0].Duration).To(Equal(newTrack.Duration))
				Expect(resp.HighlightedTracks[0].Status).To(Equal(newTrack.Status))
				Expect(resp.HighlightedTracks[0].TrackNumber).To(Equal(newTrack.TrackNumber))
			})
			It("should respond with user_group (label) if it exists", func() {
				userGroup := &pb.UserGroup{Id: newLabel.Id.String()}
				resp, err := service.GetUserGroup(context.Background(), userGroup)

				Expect(err).NotTo(HaveOccurred())

				Expect(len(resp.TrackGroups)).To(Equal(2))
			})
			It("should respond with not_found error if user_group does not exist", func() {
				id := uuid.NewV4()
				for (id == newArtist.Id || id == newRecommendedArtist.Id || id == newLabel.Id || id == newDistributor.Id || id == newUserProfile.Id) {
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

	Describe("SearchUserGroups", func() {
    Context("with valid query", func() {
      It("should respond with public people, artists and labels", func() {
        q := &tagpb.Query{Query: "best"}
        res, err := service.SearchUserGroups(context.Background(), q)

        Expect(err).NotTo(HaveOccurred())
        Expect(res).NotTo(BeNil())

				Expect(len(res.People)).To(Equal(1))
				Expect(res.People[0].Id).To(Equal(newUserProfile.Id.String()))
				Expect(res.People[0].DisplayName).To(Equal(newUserProfile.DisplayName))
				Expect(res.People[0].Avatar).To(Equal(newUserProfile.Avatar))

				Expect(len(res.Artists)).To(Equal(1))
				Expect(res.Artists[0].Id).To(Equal(newArtist.Id.String()))
				Expect(res.Artists[0].DisplayName).To(Equal(newArtist.DisplayName))
				Expect(res.Artists[0].Avatar).To(Equal(newArtist.Avatar))

				Expect(len(res.Labels)).To(Equal(1))
				Expect(res.Labels[0].Id).To(Equal(newLabel.Id.String()))
				Expect(res.Labels[0].DisplayName).To(Equal(newLabel.DisplayName))
				Expect(res.Labels[0].Avatar).To(Equal(newLabel.Avatar))
      })
    })
    Context("with invalid query", func() {
      It("should respond with invalid error", func() {
        q := &tagpb.Query{Query: "ti"}
        resp, err := service.SearchUserGroups(context.Background(), q)

        Expect(resp).To(BeNil())
        Expect(err).To(HaveOccurred())

        twerr := err.(twirp.Error)
        Expect(twerr.Code()).To(Equal(invalid_argument_code))
        Expect(twerr.Meta("argument")).To(Equal("query"))
      })
    })
  })

	Describe("GetLabelUserGroups", func() {
		It("should respond with user_groups of type label", func() {
			emptyReq := &tagpb.Empty{}
			u := url.URL{}
			queryString := u.Query()
			queryString.Set("page", "1")
			queryString.Set("limit", "50")
			ctx := context.WithValue(context.Background(), "query", queryString)
			resp, err := service.GetLabelUserGroups(ctx, emptyReq)

			Expect(err).NotTo(HaveOccurred())
			Expect(resp).NotTo(BeNil())
			Expect(len(resp.Labels)).To(Equal(1))
			Expect(resp.Labels[0].Id).To(Equal(newLabel.Id.String()))
			Expect(resp.Labels[0].DisplayName).To(Equal(newLabel.DisplayName))
			Expect(resp.Labels[0].Avatar).To(Equal(newLabel.Avatar))
		})
	})

	/*Describe("GetTrackAnalytics", func() {
		Context("with valid uuid", func() {
			It("should respond with artist track analytics", func() {
				userGroup := &pb.UserGroup{Id: newArtist.Id.String()}
				res, err := service.GetTrackAnalytics(context.Background(), userGroup)

				Expect(err).NotTo(HaveOccurred())
				Expect(res).NotTo(BeNil())
				Expect(len(res.ArtistTrackAnalytics)).To(Equal(1))
				Expect(res.ArtistTrackAnalytics[0].Id).To(Equal(newTrack.Id.String()))
				Expect(res.ArtistTrackAnalytics[0].Title).To(Equal(newTrack.Title))
				Expect(res.ArtistTrackAnalytics[0].TotalPlays).To(Equal(int32(4)))
				Expect(res.ArtistTrackAnalytics[0].PaidPlays).To(Equal(int32(3)))
				Expect(res.ArtistTrackAnalytics[0].FreePlays).To(Equal(int32(1)))
				Expect(res.ArtistTrackAnalytics[0].TotalCredits).To(Equal(float32(0.06)))
				Expect(res.ArtistTrackAnalytics[0].UserGroupCredits).To(Equal(0.7*float32(0.06)))
				Expect(res.ArtistTrackAnalytics[0].ResonateCredits).To(Equal(0.3*float32(0.06)))
			})

			It("should respond with label track analytics", func() {
				userGroup := &pb.UserGroup{Id: newLabel.Id.String()}
				res, err := service.GetTrackAnalytics(context.Background(), userGroup)

				Expect(err).NotTo(HaveOccurred())
				Expect(res).NotTo(BeNil())
				Expect(len(res.LabelTrackAnalytics)).To(Equal(2))

				Expect(res.LabelTrackAnalytics[0].UserGroup).NotTo(BeNil())
				Expect(res.LabelTrackAnalytics[0].UserGroup.Id).To(Equal(newArtist.Id.String()))
				Expect(res.LabelTrackAnalytics[0].UserGroup.DisplayName).To(Equal(newArtist.DisplayName))
				Expect(res.LabelTrackAnalytics[0].UserGroup.Avatar).To(Equal(newArtist.Avatar))
				Expect(len(res.LabelTrackAnalytics[0].Tracks)).To(Equal(1))
				Expect(res.LabelTrackAnalytics[0].Tracks[0].Id).To(Equal(newTrack.Id.String()))
				Expect(res.LabelTrackAnalytics[0].Tracks[0].Title).To(Equal(newTrack.Title))
				Expect(res.LabelTrackAnalytics[0].Tracks[0].TotalPlays).To(Equal(int32(4)))
				Expect(res.LabelTrackAnalytics[0].Tracks[0].PaidPlays).To(Equal(int32(3)))
				Expect(res.LabelTrackAnalytics[0].Tracks[0].FreePlays).To(Equal(int32(1)))
				Expect(res.LabelTrackAnalytics[0].Tracks[0].TotalCredits).To(Equal(float32(0.06)))
				Expect(res.LabelTrackAnalytics[0].Tracks[0].UserGroupCredits).To(Equal(0.7*float32(0.06)))
				Expect(res.LabelTrackAnalytics[0].Tracks[0].ResonateCredits).To(Equal(0.3*float32(0.06)))

				Expect(res.LabelTrackAnalytics[1].UserGroup).NotTo(BeNil())
				Expect(res.LabelTrackAnalytics[1].UserGroup.Id).To(Equal(newLabel.Id.String()))
				Expect(res.LabelTrackAnalytics[1].UserGroup.DisplayName).To(Equal(newLabel.DisplayName))
				Expect(res.LabelTrackAnalytics[1].UserGroup.Avatar).To(Equal(newLabel.Avatar))
				Expect(len(res.LabelTrackAnalytics[1].Tracks)).To(Equal(1))
				Expect(res.LabelTrackAnalytics[1].Tracks[0].Id).To(Equal(featuringTrack.Id.String()))
				Expect(res.LabelTrackAnalytics[1].Tracks[0].Title).To(Equal(featuringTrack.Title))
				Expect(res.LabelTrackAnalytics[1].Tracks[0].TotalPlays).To(Equal(int32(1)))
				Expect(res.LabelTrackAnalytics[1].Tracks[0].PaidPlays).To(Equal(int32(1)))
				Expect(res.LabelTrackAnalytics[1].Tracks[0].FreePlays).To(Equal(int32(0)))
				Expect(res.LabelTrackAnalytics[1].Tracks[0].TotalCredits).To(Equal(float32(0.02)))
				Expect(res.LabelTrackAnalytics[1].Tracks[0].UserGroupCredits).To(Equal(0.7*float32(0.02)))
				Expect(res.LabelTrackAnalytics[1].Tracks[0].ResonateCredits).To(Equal(0.3*float32(0.02)))
			})

			It("should respond with not_found error if user_group does not exist", func() {
				id := uuid.NewV4()
				for (id == newArtist.Id || id == newRecommendedArtist.Id || id == newLabel.Id || id == newDistributor.Id || id == newUserProfile.Id) {
					id = uuid.NewV4()
				}
				userGroup := &pb.UserGroup{Id: id.String()}
				res, err := service.GetTrackAnalytics(context.Background(), userGroup)

				Expect(res).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(not_found_code))
			})

		})
		Context("with invalid uuid", func() {
			It("should respond with invalid_argument error if user_group is invalid", func() {
				userGroup := &pb.UserGroup{Id: "123"}
				res, err := service.GetTrackAnalytics(context.Background(), userGroup)

				Expect(res).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
				Expect(twerr.Meta("argument")).To(Equal("id"))
			})
		})
	})*/

	Describe("UpdateUserGroup", func() {
		Context("with valid uuid", func() {
			It("should update user_group if it exists", func() {
				tags := []*tagpb.Tag{&tagpb.Tag{Type: "genre", Name: "experimental"}}
				links := []*pb.Link{&pb.Link{Platform: "instagram", Uri: "https://instagram/bestartistever"}}
				userGroup := &pb.UserGroup{
					Id: newArtist.Id.String(),
					DisplayName: "new display name",
					Description: "new description",
					ShortBio: "short bio",
					Avatar: newArtist.Avatar,
					Publisher: map[string]string{
						"name": "new publisher name",
						"number": "new 1E3",
					},
					Pro: map[string]string{
						"name": "new PRO name",
						"number": "new 2BA",
					},
					Address: &userpb.StreetAddress{Id: artistAddress.Id.String(), Data: map[string]string{"some": "new data"}},
					Type: &pb.GroupTaxonomy{Id: newArtistGroupTaxonomy.Id.String(), Type: "artist"},
					Privacy: &pb.Privacy{Id: newArtist.Privacy.Id.String(), Private: true, OwnedTracks: false, SupportedArtists: true},
					OwnerId: newArtist.OwnerId.String(),
					Tags: tags,
					Links: links,
				}
				_, err := service.UpdateUserGroup(context.Background(), userGroup)

				Expect(err).NotTo(HaveOccurred())

				address := new(models.StreetAddress)
				err = db.Model(address).Where("id = ?", artistAddress.Id).Select()
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
				Expect(updatedUserGroup.DisplayName).To(Equal(userGroup.DisplayName))
				Expect(updatedUserGroup.ShortBio).To(Equal(userGroup.ShortBio))
				Expect(updatedUserGroup.Avatar).To(Equal(userGroup.Avatar))
				Expect(updatedUserGroup.Description).To(Equal(userGroup.Description))
				Expect(updatedUserGroup.Banner).To(Equal(userGroup.Banner))
				Expect(updatedUserGroup.Publisher).To(Equal(userGroup.Publisher))
				Expect(updatedUserGroup.Pro).To(Equal(userGroup.Pro))
				Expect(len(updatedUserGroup.Tags)).To(Equal(1))
				Expect(updatedUserGroup.Tags[0]).NotTo(Equal(newGenreTag.Id))

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

				deletedLink := models.Link{Id: newLink.Id}
				err = db.Model(&deletedLink).
					WherePK().
					Select()
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(pg.ErrNoRows))

				// unchanged
				Expect(updatedUserGroup.OwnerId).To(Equal(newArtist.OwnerId))
				Expect(updatedUserGroup.TypeId).To(Equal(newArtist.TypeId))
			})
			It("should update user group type from user to artist", func() {
				userGroup := &pb.UserGroup{
					Id: newUserProfile.Id.String(),
					DisplayName: newUserProfile.DisplayName,
					Description: newUserProfile.Description,
					ShortBio: newUserProfile.ShortBio,
					Avatar: newUserProfile.Avatar,
					Address: &userpb.StreetAddress{Id: newAddress.Id.String(), Data: newAddress.Data},
					Type: &pb.GroupTaxonomy{Id: newArtistGroupTaxonomy.Id.String(), Type: "artist"},
					Privacy: &pb.Privacy{Id: newArtist.Privacy.Id.String(), Private: true, OwnedTracks: false, SupportedArtists: true},
					OwnerId: newUserProfile.OwnerId.String(),
				}
				_, err := service.UpdateUserGroup(context.Background(), userGroup)

				Expect(err).NotTo(HaveOccurred())

				updatedUserGroup := new(models.UserGroup)
				err = db.Model(updatedUserGroup).Where("id = ?", newUserProfile.Id).Select()
				Expect(err).NotTo(HaveOccurred())
				Expect(updatedUserGroup.TypeId).To(Equal(newArtistGroupTaxonomy.Id))
			})
			It("should respond with not_found error if user_group does not exist", func() {
				id := uuid.NewV4()
				for (id == newArtist.Id || id == newRecommendedArtist.Id || id == newLabel.Id || id == newDistributor.Id || id == newUserProfile.Id) {
					id = uuid.NewV4()
				}
				userGroup := &pb.UserGroup{
					Id: id.String(),
					DisplayName: "new display name",
					Description: "new description",
					Avatar: newArtist.Avatar,
					Address: &userpb.StreetAddress{Id: artistAddress.Id.String(), Data: map[string]string{"some": "data"}},
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
		Context("with invalid type", func() {
			It("should not update user group from artist to user", func() {
				userGroup := &pb.UserGroup{
					Id: newArtist.Id.String(),
					DisplayName: "new display name",
					Description: "new description",
					ShortBio: "short bio",
					Avatar: newArtist.Avatar,
					Publisher: map[string]string{
						"name": "new publisher name",
						"number": "new 1E3",
					},
					Pro: map[string]string{
						"name": "new PRO name",
						"number": "new 2BA",
					},
					Address: &userpb.StreetAddress{Id: artistAddress.Id.String(), Data: map[string]string{"some": "new data"}},
					Type: &pb.GroupTaxonomy{Id: newUserGroupTaxonomy.Id.String(), Type: "user"},
					Privacy: &pb.Privacy{Id: newArtist.Privacy.Id.String(), Private: true, OwnedTracks: false, SupportedArtists: true},
					OwnerId: newArtist.OwnerId.String(),
				}
				res, err := service.UpdateUserGroup(context.Background(), userGroup)

				Expect(res).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
				Expect(twerr.Meta("argument")).To(Equal("type"))
			})
			It("should not update user group from label to user", func() {
				userGroup := &pb.UserGroup{
					Id: newLabel.Id.String(),
					DisplayName: newLabel.DisplayName,
					Avatar: newLabel.Avatar,
					Address: &userpb.StreetAddress{Id: labelAddress.Id.String(), Data: labelAddress.Data},
					Type: &pb.GroupTaxonomy{Id: newUserGroupTaxonomy.Id.String(), Type: "user"},
					Privacy: &pb.Privacy{Id: newArtist.Privacy.Id.String(), Private: true, OwnedTracks: false, SupportedArtists: true},
					OwnerId: newLabel.OwnerId.String(),
				}
				res, err := service.UpdateUserGroup(context.Background(), userGroup)

				Expect(res).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
				Expect(twerr.Meta("argument")).To(Equal("type"))
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
					Address: &userpb.StreetAddress{Id: artistAddress.Id.String(), Data: map[string]string{"some": "data"}},
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

	Describe("AddRecommended", func () {
		Context("with user_group_id and recommended_id", func() {
			It("should add recommended and recommended_by to user groups", func() {
				userGroupRecommended := &pb.UserGroupRecommended{
					UserGroupId: newArtist.Id.String(),
					RecommendedId: newRecommendedArtist.Id.String(),
				}
				_, err := service.AddRecommended(context.Background(), userGroupRecommended)

				Expect(err).NotTo(HaveOccurred())

				userGroup := new(models.UserGroup)
				err = db.Model(userGroup).Where("id = ?", newArtist.Id).Select()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(userGroup.RecommendedArtists)).To(Equal(2))
				Expect(userGroup.RecommendedArtists).To(ContainElement(newRecommendedArtist.Id))

				recommended := new(models.UserGroup)
				err = db.Model(recommended).Where("id = ?", newRecommendedArtist.Id).Select()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(recommended.RecommendedBy)).To(Equal(1))
				Expect(recommended.RecommendedBy).To(ContainElement(newArtist.Id))
			})
			It("should respond with not_found error if user group does exist", func() {
				id := uuid.NewV4()
				for (id == newArtist.Id || id == newRecommendedArtist.Id || id == newLabel.Id || id == newDistributor.Id || id == newUserProfile.Id) {
					id = uuid.NewV4()
				}
				userGroupRecommended := &pb.UserGroupRecommended{
					UserGroupId: id.String(),
					RecommendedId: newRecommendedArtist.Id.String(),
				}
				resp, err := service.AddRecommended(context.Background(), userGroupRecommended)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(not_found_code))
			})
			It("should respond with not_found error if recommended user group does exist", func() {
				id := uuid.NewV4()
				for (id == newArtist.Id || id == newRecommendedArtist.Id || id == newLabel.Id || id == newDistributor.Id || id == newUserProfile.Id) {
					id = uuid.NewV4()
				}
				userGroupRecommended := &pb.UserGroupRecommended{
					UserGroupId: newArtist.Id.String(),
					RecommendedId: id.String(),
				}
				resp, err := service.AddRecommended(context.Background(), userGroupRecommended)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(not_found_code))
			})
		})
		Context("with invalid user_group_id", func() {
			It("should respond with invalid_argument_code error", func() {
				userGroupRecommended := &pb.UserGroupRecommended{
					UserGroupId: "1",
					RecommendedId: newRecommendedArtist.Id.String(),
				}
				resp, err := service.AddRecommended(context.Background(), userGroupRecommended)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
			})
		})
		Context("with invalid recommended_id", func() {
			It("should respond with invalid_argument_code error", func() {
				userGroupRecommended := &pb.UserGroupRecommended{
					UserGroupId: newArtist.Id.String(),
					RecommendedId: "",
				}
				resp, err := service.AddRecommended(context.Background(), userGroupRecommended)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
			})
		})
	})

	Describe("RemoveRecommended", func () {
		Context("with user_group_id and recommended_id", func() {
			It("should remove recommended and recommended_by from user groups", func() {
				userGroupRecommended := &pb.UserGroupRecommended{
					UserGroupId: newArtist.Id.String(),
					RecommendedId: newRecommendedArtist.Id.String(),
				}
				_, err := service.RemoveRecommended(context.Background(), userGroupRecommended)

				Expect(err).NotTo(HaveOccurred())

				userGroup := new(models.UserGroup)
				err = db.Model(userGroup).Where("id = ?", newArtist.Id).Select()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(userGroup.RecommendedArtists)).To(Equal(1))
				Expect(userGroup.RecommendedArtists).NotTo(ContainElement(newRecommendedArtist.Id))

				recommended := new(models.UserGroup)
				err = db.Model(recommended).Where("id = ?", newRecommendedArtist.Id).Select()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(recommended.RecommendedBy)).To(Equal(0))
				Expect(recommended.RecommendedBy).NotTo(ContainElement(newArtist.Id))
			})
			It("should respond with not_found error if user group does exist", func() {
				id := uuid.NewV4()
				for (id == newArtist.Id || id == newRecommendedArtist.Id || id == newLabel.Id || id == newDistributor.Id || id == newUserProfile.Id) {
					id = uuid.NewV4()
				}
				userGroupRecommended := &pb.UserGroupRecommended{
					UserGroupId: id.String(),
					RecommendedId: newRecommendedArtist.Id.String(),
				}
				resp, err := service.RemoveRecommended(context.Background(), userGroupRecommended)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(not_found_code))
			})
			It("should respond with not_found error if recommended user group does exist", func() {
				id := uuid.NewV4()
				for (id == newArtist.Id || id == newRecommendedArtist.Id || id == newLabel.Id || id == newDistributor.Id || id == newUserProfile.Id) {
					id = uuid.NewV4()
				}
				userGroupRecommended := &pb.UserGroupRecommended{
					UserGroupId: newArtist.Id.String(),
					RecommendedId: id.String(),
				}
				resp, err := service.RemoveRecommended(context.Background(), userGroupRecommended)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(not_found_code))
			})
		})
		Context("with invalid user_group_id", func() {
			It("should respond with invalid_argument_code error", func() {
				userGroupRecommended := &pb.UserGroupRecommended{
					UserGroupId: "1",
					RecommendedId: newRecommendedArtist.Id.String(),
				}
				resp, err := service.RemoveRecommended(context.Background(), userGroupRecommended)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
			})
		})
		Context("with invalid recommended_id", func() {
			It("should respond with invalid_argument_code error", func() {
				userGroupRecommended := &pb.UserGroupRecommended{
					UserGroupId: newArtist.Id.String(),
					RecommendedId: "",
				}
				resp, err := service.RemoveRecommended(context.Background(), userGroupRecommended)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
			})
		})
	})

	Describe("AddMembers", func() {
		Context("with valid UserGroupId and Members ids", func() {
			It("should add new members with given display_name", func() {
				userGroupMembers := &pb.UserGroupMembers{
					UserGroupId: newRecommendedArtist.Id.String(),
					Members: []*pb.UserGroup{
						&pb.UserGroup{
							Id: newUserProfile.Id.String(),
							DisplayName: "John Doe",
							Tags: []*tagpb.Tag{
								&tagpb.Tag{
									Type: "role",
									Name: "keyboard",
								},
								&tagpb.Tag{
									Type: "role",
									Name: "singer",
								},
							},
						},
					},
				}

				_, err := service.AddMembers(context.Background(), userGroupMembers)

				Expect(err).NotTo(HaveOccurred())

				artist := models.UserGroup{Id: newRecommendedArtist.Id}
				err = db.Model(&artist).
					Column("Members").
					WherePK().
					Select()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(artist.Members)).To(Equal(1))
				Expect(artist.Members[0].Id).To(Equal(newUserProfile.Id))

				userProfile := models.UserGroup{Id: newUserProfile.Id}
				err = db.Model(&userProfile).
					Column("MemberOfGroups").
					WherePK().
					Select()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(userProfile.MemberOfGroups)).To(Equal(2))

				userGroupMember := models.UserGroupMember{UserGroupId: newRecommendedArtist.Id, MemberId: newUserProfile.Id}
				err = db.Model(&userGroupMember).
					Where("user_group_id = ?user_group_id").
					Where("member_id = ?member_id").
					Select()
				Expect(err).NotTo(HaveOccurred())
				Expect(userGroupMember.DisplayName).To(Equal("John Doe"))
				Expect(len(userGroupMember.Tags)).To(Equal(2))
			})
			It("should add new members with default display_name", func() {
				userGroupMembers := &pb.UserGroupMembers{
					UserGroupId: newDistributor.Id.String(),
					Members: []*pb.UserGroup{
						&pb.UserGroup{
							Id: newLabel.Id.String(),
						},
					},
				}

				_, err := service.AddMembers(context.Background(), userGroupMembers)

				Expect(err).NotTo(HaveOccurred())

				label := models.UserGroup{Id: newLabel.Id}
				err = db.Model(&label).
					Column("MemberOfGroups").
					WherePK().
					Select()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(label.MemberOfGroups)).To(Equal(1))
				Expect(label.MemberOfGroups[0].Id).To(Equal(newDistributor.Id))

				userGroupMember := models.UserGroupMember{UserGroupId: newDistributor.Id, MemberId: newLabel.Id}
				err = db.Model(&userGroupMember).
					Where("user_group_id = ?user_group_id").
					Where("member_id = ?member_id").
					Select()
				Expect(err).NotTo(HaveOccurred())
				Expect(userGroupMember.DisplayName).To(Equal(newLabel.DisplayName))
				Expect(len(userGroupMember.Tags)).To(Equal(0))
			})
		})
		Context("with invalid UserGroupId or Members ids", func() {
			It("should respond with invalid argument error if one of Members id not valid", func() {
				userGroupMembers := &pb.UserGroupMembers{
					UserGroupId: newArtist.Id.String(),
					Members: []*pb.UserGroup{
						&pb.UserGroup{
							Id: "1",
						},
					},
				}
				resp, err := service.AddMembers(context.Background(), userGroupMembers)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
				Expect(twerr.Meta("argument")).To(Equal("member id"))
			})
			It("should respond with invalid argument error if UserGroupId not valid", func() {
				userGroupMembers := &pb.UserGroupMembers{
					UserGroupId: "1",
					Members: []*pb.UserGroup{
						&pb.UserGroup{
							Id: newUserProfile.Id.String(),
						},
					},
				}
				resp, err := service.AddMembers(context.Background(), userGroupMembers)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
				Expect(twerr.Meta("argument")).To(Equal("user_group id"))
			})
			It("should respond with not found error if one of Members does not exist", func() {
				id := uuid.NewV4()
				for (id == newArtist.Id || id == newRecommendedArtist.Id || id == newLabel.Id || id == newDistributor.Id || id == newUserProfile.Id) {
					id = uuid.NewV4()
				}
				userGroupMembers := &pb.UserGroupMembers{
					UserGroupId: newArtist.Id.String(),
					Members: []*pb.UserGroup{
						&pb.UserGroup{
							Id: id.String(),
						},
					},
				}

				resp, err := service.AddMembers(context.Background(), userGroupMembers)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(not_found_code))
			})
			It("should respond with not found error if UserGroup does not exist", func() {
				id := uuid.NewV4()
				for (id == newArtist.Id || id == newRecommendedArtist.Id || id == newLabel.Id || id == newDistributor.Id || id == newUserProfile.Id) {
					id = uuid.NewV4()
				}
				userGroupMembers := &pb.UserGroupMembers{
					UserGroupId: id.String(),
					Members: []*pb.UserGroup{
						&pb.UserGroup{
							Id: newUserProfile.Id.String(),
						},
					},
				}

				resp, err := service.AddMembers(context.Background(), userGroupMembers)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(not_found_code))
			})
		})
	})

	Describe("DeleteMembers", func() {
		Context("with valid UserGroupId and Members ids", func() {
			It("should delete members", func() {
				userGroupMembers := &pb.UserGroupMembers{
					UserGroupId: newRecommendedArtist.Id.String(),
					Members: []*pb.UserGroup{
						&pb.UserGroup{
							Id: newUserProfile.Id.String(),
							DisplayName: "John Doe",
							Tags: []*tagpb.Tag{
								&tagpb.Tag{
									Type: "role",
									Name: "keyboard",
								},
								&tagpb.Tag{
									Type: "role",
									Name: "singer",
								},
							},
						},
					},
				}

				_, err := service.DeleteMembers(context.Background(), userGroupMembers)

				Expect(err).NotTo(HaveOccurred())

				artist := models.UserGroup{Id: newRecommendedArtist.Id}
				err = db.Model(&artist).
					Column("Members").
					WherePK().
					Select()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(artist.Members)).To(Equal(0))

				userProfile := models.UserGroup{Id: newUserProfile.Id}
				err = db.Model(&userProfile).
					Column("MemberOfGroups").
					WherePK().
					Select()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(userProfile.MemberOfGroups)).To(Equal(1))

				userGroupMember := models.UserGroupMember{UserGroupId: newRecommendedArtist.Id, MemberId: newUserProfile.Id}
				err = db.Model(&userGroupMember).
					Where("user_group_id = ?user_group_id").
					Where("member_id = ?member_id").
					Select()
				Expect(err).To(HaveOccurred())
			})
		})
		Context("with invalid UserGroupId or Members ids", func() {
			It("should respond with invalid argument error if one of Members id not valid", func() {
				userGroupMembers := &pb.UserGroupMembers{
					UserGroupId: newArtist.Id.String(),
					Members: []*pb.UserGroup{
						&pb.UserGroup{
							Id: "1",
						},
					},
				}
				resp, err := service.DeleteMembers(context.Background(), userGroupMembers)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
				Expect(twerr.Meta("argument")).To(Equal("member id"))
			})
			It("should respond with invalid argument error if UserGroupId not valid", func() {
				userGroupMembers := &pb.UserGroupMembers{
					UserGroupId: "1",
					Members: []*pb.UserGroup{
						&pb.UserGroup{
							Id: newUserProfile.Id.String(),
						},
					},
				}
				resp, err := service.DeleteMembers(context.Background(), userGroupMembers)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
				Expect(twerr.Meta("argument")).To(Equal("user_group id"))
			})
			It("should respond with not found error if one of Members does not exist", func() {
				id := uuid.NewV4()
				for (id == newArtist.Id || id == newRecommendedArtist.Id || id == newLabel.Id || id == newDistributor.Id || id == newUserProfile.Id) {
					id = uuid.NewV4()
				}
				userGroupMembers := &pb.UserGroupMembers{
					UserGroupId: newArtist.Id.String(),
					Members: []*pb.UserGroup{
						&pb.UserGroup{
							Id: id.String(),
						},
					},
				}

				resp, err := service.DeleteMembers(context.Background(), userGroupMembers)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(not_found_code))
			})
			It("should respond with not found error if UserGroup does not exist", func() {
				id := uuid.NewV4()
				for (id == newArtist.Id || id == newRecommendedArtist.Id || id == newLabel.Id || id == newDistributor.Id || id == newUserProfile.Id) {
					id = uuid.NewV4()
				}
				userGroupMembers := &pb.UserGroupMembers{
					UserGroupId: id.String(),
					Members: []*pb.UserGroup{
						&pb.UserGroup{
							Id: newUserProfile.Id.String(),
						},
					},
				}

				resp, err := service.DeleteMembers(context.Background(), userGroupMembers)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(not_found_code))
			})
		})
	})

	Describe("DeleteUserGroup", func() {
		Context("with valid uuid", func() {
			It("should delete user_group (artist) if it exists", func() {
				userGroup := &pb.UserGroup{Id: newArtist.Id.String()}

				userGroupToDelete := new(models.UserGroup)
				err := db.Model(userGroupToDelete).
					Column("OwnerOfTracks", "OwnerOfTrackGroups", "user_group.*").
					Where("id = ?", newArtist.Id).Select()
				Expect(err).NotTo(HaveOccurred())

				_, err = service.DeleteUserGroup(context.Background(), userGroup)
				Expect(err).NotTo(HaveOccurred())

				var links []*models.Link
				err = db.Model(&links).
					Where("id in (?)", pg.In(userGroupToDelete.Links)).
					Select()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(links)).To(Equal(0))

				var privacies []*models.UserGroupPrivacy
				err = db.Model(&privacies).
					Where("id in (?)", pg.In([]uuid.UUID{userGroupToDelete.PrivacyId})).
					Select()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(privacies)).To(Equal(0))

				var addresses []*models.StreetAddress
				err = db.Model(&addresses).
					Where("id in (?)", pg.In([]uuid.UUID{userGroupToDelete.AddressId})).
					Select()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(privacies)).To(Equal(0))

				var userGroupMembers []models.UserGroupMember
				err = db.Model(&userGroupMembers).
					Where("user_group_id = ?", newArtist.Id).
					WhereOr("member_id = ?", newArtist.Id).
					Select()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(userGroupMembers)).To(Equal(0))

				var recommendedBy []models.UserGroup
				err = db.Model(&recommendedBy).
					Where("id in (?)", pg.In(userGroupToDelete.RecommendedBy)).
					Select()
				Expect(err).NotTo(HaveOccurred())
				for _, r := range recommendedBy {
					Expect(r.RecommendedArtists).NotTo(ContainElement(userGroupToDelete.Id))
				}

				var recommendedArtists []models.UserGroup
				err = db.Model(&recommendedArtists).
					Where("id in (?)", pg.In(userGroupToDelete.RecommendedArtists)).
					Select()
				Expect(err).NotTo(HaveOccurred())
				for _, r := range recommendedArtists {
					Expect(r.RecommendedBy).NotTo(ContainElement(userGroupToDelete.Id))
				}

				ownerOfTrackGroupIds := make([]uuid.UUID, len(userGroupToDelete.OwnerOfTrackGroups))
				for i, trackGroup := range(userGroupToDelete.OwnerOfTrackGroups) {
				  ownerOfTrackGroupIds[i] = trackGroup.Id
				}
				var ownerOfTrackGroups []models.TrackGroup
				err = db.Model(&ownerOfTrackGroups).
					Where("id in (?)", pg.In(ownerOfTrackGroupIds)).
					Select()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(ownerOfTrackGroups)).To(Equal(0))

				ownerOfTrackIds := make([]uuid.UUID, len(userGroupToDelete.OwnerOfTracks))
				for i, track := range(userGroupToDelete.OwnerOfTracks) {
				  ownerOfTrackIds[i] = track.Id
				}
				var ownerOfTracks []models.Track
				err = db.Model(&ownerOfTracks).
					Where("id in (?)", pg.In(ownerOfTrackIds)).
					Select()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(ownerOfTracks)).To(Equal(0))

				var artistOfTracks []models.Track
				err = db.Model(&artistOfTracks).
					Where("id in (?)", pg.In(userGroupToDelete.ArtistOfTracks)).
					Select()
				Expect(err).NotTo(HaveOccurred())
				for _, t := range artistOfTracks {
					Expect(t.Artists).NotTo(ContainElement(userGroupToDelete.Id))
				}

				var userGroups []models.UserGroup
				err = db.Model(&userGroups).
					Where("id in (?)", pg.In([]uuid.UUID{userGroupToDelete.Id})).
					Select()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(userGroups)).To(Equal(0))
			})
			It("should delete user_group (label) if it exists", func() {
				userGroup := &pb.UserGroup{Id: newLabel.Id.String()}

				userGroupToDelete := new(models.UserGroup)
				err := db.Model(userGroupToDelete).
					Column("OwnerOfTracks", "OwnerOfTrackGroups", "user_group.*", "LabelOfTrackGroups").
					Where("id = ?", newLabel.Id).Select()
				Expect(err).NotTo(HaveOccurred())

				_, err = service.DeleteUserGroup(context.Background(), userGroup)
				Expect(err).NotTo(HaveOccurred())

				var privacies []*models.UserGroupPrivacy
				err = db.Model(&privacies).
					Where("id in (?)", pg.In([]uuid.UUID{userGroupToDelete.PrivacyId})).
					Select()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(privacies)).To(Equal(0))

				var addresses []*models.StreetAddress
				err = db.Model(&addresses).
					Where("id in (?)", pg.In([]uuid.UUID{userGroupToDelete.AddressId})).
					Select()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(privacies)).To(Equal(0))

				var userGroupMembers []models.UserGroupMember
				err = db.Model(&userGroupMembers).
					Where("user_group_id = ?", newLabel.Id).
					WhereOr("member_id = ?", newLabel.Id).
					Select()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(userGroupMembers)).To(Equal(0))

				ownerOfTrackGroupIds := make([]uuid.UUID, len(userGroupToDelete.OwnerOfTrackGroups))
				for i, trackGroup := range(userGroupToDelete.OwnerOfTrackGroups) {
				  ownerOfTrackGroupIds[i] = trackGroup.Id
				}
				var ownerOfTrackGroups []models.TrackGroup
				err = db.Model(&ownerOfTrackGroups).
					Where("id in (?)", pg.In(ownerOfTrackGroupIds)).
					Select()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(ownerOfTrackGroups)).To(Equal(0))

				labelOfTrackGroupIds := make([]uuid.UUID, len(userGroupToDelete.LabelOfTrackGroups))
				for i, trackGroup := range(userGroupToDelete.LabelOfTrackGroups) {
					labelOfTrackGroupIds[i] = trackGroup.Id
				}
				var labelOfTrackGroups []models.TrackGroup
				err = db.Model(&labelOfTrackGroups).
					Where("id in (?)", pg.In(labelOfTrackGroupIds)).
					Select()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(labelOfTrackGroups)).To(Equal(0))

				ownerOfTrackIds := make([]uuid.UUID, len(userGroupToDelete.OwnerOfTracks))
				for i, track := range(userGroupToDelete.OwnerOfTracks) {
				  ownerOfTrackIds[i] = track.Id
				}
				var ownerOfTracks []models.Track
				err = db.Model(&ownerOfTracks).
					Where("id in (?)", pg.In(ownerOfTrackIds)).
					Select()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(ownerOfTracks)).To(Equal(0))

				var userGroups []models.UserGroup
				err = db.Model(&userGroups).
					Where("id in (?)", pg.In([]uuid.UUID{userGroupToDelete.Id})).
					Select()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(userGroups)).To(Equal(0))
			})
			It("should respond with not_found error if user_group does not exist", func() {
				id := uuid.NewV4()
				for (id == newArtist.Id || id == newRecommendedArtist.Id || id == newLabel.Id || id == newDistributor.Id || id == newUserProfile.Id) {
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
				tags := make([]*tagpb.Tag, 1)
				tags[0] = &tagpb.Tag{Type: "genre", Name: "rock"}
				ownerId := newUser.Id.String()
				userGroup := &pb.UserGroup{
					DisplayName: "group2",
					Avatar: avatar,
					Type: &pb.GroupTaxonomy{Type: "artist"},
					OwnerId: ownerId,
					ShortBio: "short bio",
					Address: &userpb.StreetAddress{Data: map[string]string{"some": "data"}},
					Tags: tags,
					RecommendedArtists: []*tagpb.RelatedUserGroup{
						&tagpb.RelatedUserGroup{Id: newRecommendedArtist.Id.String()},
					},
					Publisher: map[string]string{
						"name": "publisher name",
						"number": "1E3",
					},
					Pro: map[string]string{
						"name": "PRO name",
						"number": "2BA",
					},
				}
				resp, err := service.CreateUserGroup(context.Background(), userGroup)

				Expect(err).NotTo(HaveOccurred())
				Expect(resp.Id).NotTo(Equal(""))
				Expect(resp.DisplayName).To(Equal(userGroup.DisplayName))
				Expect(resp.ShortBio).To(Equal(userGroup.ShortBio))
				Expect(resp.Publisher).To(Equal(userGroup.Publisher))
				Expect(resp.Pro).To(Equal(userGroup.Pro))
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

				id, err := uuid.FromString(resp.Id)
				Expect(err).NotTo(HaveOccurred())
				createdUserGroup := new(models.UserGroup)
				err = db.Model(createdUserGroup).Where("id = ?", id).Select()
				Expect(err).NotTo(HaveOccurred())
				Expect(createdUserGroup.RecommendedArtists).To(ContainElement(newRecommendedArtist.Id))

				recommended := new(models.UserGroup)
				err = db.Model(recommended).Where("id = ?", newRecommendedArtist.Id).Select()
				Expect(err).NotTo(HaveOccurred())
				Expect(recommended.RecommendedBy).To(ContainElement(id))
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
			emptyReq := &tagpb.Empty{}
			groupTaxonomies, err := service.GetUserGroupTypes(context.Background(), emptyReq)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(groupTaxonomies.Types)).To(Equal(3))
		})
	})
})
