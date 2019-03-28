package trackgroupserver_test

import (
  "context"
  "time"
  // "fmt"
  "github.com/go-pg/pg"
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "github.com/twitchtv/twirp"
  "github.com/satori/go.uuid"
  "github.com/golang/protobuf/ptypes"

  pb "user-api/rpc/trackgroup"
  trackpb "user-api/rpc/track"
  tagpb "user-api/rpc/tag"
  "user-api/internal/model"
)


var _ = Describe("TrackGroup server", func() {
	const already_exists_code twirp.ErrorCode = "already_exists"
	const invalid_argument_code twirp.ErrorCode = "invalid_argument"
	const not_found_code twirp.ErrorCode = "not_found"

  Describe("GetTrackGroup", func() {
    Context("with valid uuid", func() {
      It("should respond with release if it exists", func() {
        release := &pb.TrackGroup{Id: newAlbum.Id.String()}

        res, err := service.GetTrackGroup(context.Background(), release)

        Expect(err).NotTo(HaveOccurred())
        Expect(res.Id).To(Equal(newAlbum.Id.String()))
        Expect(res.Title).To(Equal(newAlbum.Title))
        releaseDate, err := ptypes.TimestampProto(newAlbum.ReleaseDate)
        Expect(err).NotTo(HaveOccurred())
        Expect(res.ReleaseDate.Seconds).To(Equal(releaseDate.Seconds))
        Expect(res.Type).To(Equal(newAlbum.Type))
        Expect(res.Cover).To(Equal(newAlbum.Cover))
        Expect(res.About).To(Equal(newAlbum.About))
        Expect(res.DisplayArtist).To(Equal(newAlbum.DisplayArtist))
        Expect(res.MultipleComposers).To(Equal(newAlbum.MultipleComposers))
        Expect(res.Private).To(Equal(newAlbum.Private))

        Expect(res.UserGroup.Id).To(Equal(newArtistUserGroup.Id.String()))
        Expect(res.UserGroup.DisplayName).To(Equal(newArtistUserGroup.DisplayName))
        Expect(res.UserGroup.Avatar).To(Equal(newArtistUserGroup.Avatar))

        Expect(res.Label.Id).To(Equal(newLabelUserGroup.Id.String()))
        Expect(res.Label.DisplayName).To(Equal(newLabelUserGroup.DisplayName))
        Expect(res.Label.Avatar).To(Equal(newLabelUserGroup.Avatar))

        Expect(len(res.Tags)).To(Equal(1))
        Expect(res.Tags[0].Id).To(Equal(newGenreTag.Id.String()))
        Expect(res.Tags[0].Type).To(Equal(newGenreTag.Type))
        Expect(res.Tags[0].Name).To(Equal(newGenreTag.Name))

        Expect(len(res.Tracks)).To(Equal(1))
        Expect(res.Tracks[0].Id).To(Equal(newTrack.Id.String()))
        Expect(res.Tracks[0].Title).To(Equal(newTrack.Title))
        Expect(res.Tracks[0].TrackServerId).To(Equal(newTrack.TrackServerId.String()))
        Expect(res.Tracks[0].Duration).To(Equal(newTrack.Duration))
        Expect(res.Tracks[0].Status).To(Equal(newTrack.Status))
        Expect(res.Tracks[0].TrackNumber).To(Equal(newTrack.TrackNumber))
        Expect(len(res.Tracks[0].Artists)).To(Equal(1))
        Expect(len(res.Tracks[0].TrackGroups)).To(Equal(0))
        Expect(res.Tracks[0].Artists[0].Id).To(Equal(newArtistUserGroup.Id.String()))
        Expect(res.Tracks[0].Artists[0].DisplayName).To(Equal(newArtistUserGroup.DisplayName))
        Expect(res.Tracks[0].Artists[0].Avatar).To(Equal(newArtistUserGroup.Avatar))
      })

      It("should respond with playlist if it exists", func() {
        playlist := &pb.TrackGroup{Id: newPlaylist.Id.String()}

        res, err := service.GetTrackGroup(context.Background(), playlist)

        Expect(err).NotTo(HaveOccurred())
        Expect(res.Id).To(Equal(newPlaylist.Id.String()))
        Expect(res.Title).To(Equal(newPlaylist.Title))
        releaseDate, err := ptypes.TimestampProto(newPlaylist.ReleaseDate)
        Expect(err).NotTo(HaveOccurred())
        Expect(res.ReleaseDate.Seconds).To(Equal(releaseDate.Seconds))
        Expect(res.Type).To(Equal(newPlaylist.Type))
        Expect(res.Cover).To(Equal(newPlaylist.Cover))
        Expect(res.DisplayArtist).To(Equal(newPlaylist.DisplayArtist))
        Expect(res.MultipleComposers).To(Equal(newPlaylist.MultipleComposers))
        Expect(res.Private).To(Equal(newPlaylist.Private))

        Expect(len(res.Tags)).To(Equal(0))

        Expect(len(res.Tracks)).To(Equal(1))
        Expect(res.Tracks[0].Id).To(Equal(newTrack.Id.String()))
        Expect(res.Tracks[0].Title).To(Equal(newTrack.Title))
        Expect(res.Tracks[0].TrackServerId).To(Equal(newTrack.TrackServerId.String()))
        Expect(res.Tracks[0].Duration).To(Equal(newTrack.Duration))
        Expect(res.Tracks[0].Status).To(Equal(newTrack.Status))
        Expect(res.Tracks[0].TrackNumber).To(Equal(newTrack.TrackNumber))
        Expect(len(res.Tracks[0].Artists)).To(Equal(1))
        Expect(res.Tracks[0].Artists[0].Id).To(Equal(newArtistUserGroup.Id.String()))
        Expect(res.Tracks[0].Artists[0].DisplayName).To(Equal(newArtistUserGroup.DisplayName))
        Expect(res.Tracks[0].Artists[0].Avatar).To(Equal(newArtistUserGroup.Avatar))

        Expect(len(res.Tracks[0].TrackGroups)).To(Equal(1))
        Expect(res.Tracks[0].TrackGroups[0].Id).To(Equal(newAlbum.Id.String()))
        Expect(res.Tracks[0].TrackGroups[0].Title).To(Equal(newAlbum.Title))
        Expect(res.Tracks[0].TrackGroups[0].Cover).To(Equal(newAlbum.Cover))
      })
      It("should respond with not_found error if track group does not exist", func() {
        id := uuid.NewV4()
        for id == newPlaylist.Id || id == newAlbum.Id {
          id = uuid.NewV4()
        }
        trackGroup := &pb.TrackGroup{Id: id.String()}
        resp, err := service.GetTrackGroup(context.Background(), trackGroup)

        Expect(resp).To(BeNil())
        Expect(err).To(HaveOccurred())

        twerr := err.(twirp.Error)
        Expect(twerr.Code()).To(Equal(not_found_code))
      })
    })
    Context("with invalid uuid", func() {
      It("should respond with invalid_argument error", func() {
        id := "45"
        trackGroup := &pb.TrackGroup{Id: id}
        resp, err := service.GetTrackGroup(context.Background(), trackGroup)

        Expect(resp).To(BeNil())
        Expect(err).To(HaveOccurred())

        twerr := err.(twirp.Error)
        Expect(twerr.Code()).To(Equal(invalid_argument_code))
        Expect(twerr.Meta("argument")).To(Equal("id"))
      })
    })
  })

  Describe("SearchTrackGroups", func() {
    Context("with valid query", func() {
      It("should respond with public albums and playlists", func() {
        q := &tagpb.Query{Query: "title"}
        res, err := service.SearchTrackGroups(context.Background(), q)

        Expect(err).NotTo(HaveOccurred())
        Expect(res).NotTo(BeNil())
        Expect(len(res.Playlists)).To(Equal(0))
        Expect(len(res.Albums)).To(Equal(1))
        Expect(res.Albums[0].Id).To(Equal(newAlbum.Id.String()))
        Expect(res.Albums[0].Title).To(Equal(newAlbum.Title))
        Expect(res.Albums[0].Cover).To(Equal(newAlbum.Cover))
        Expect(res.Albums[0].DisplayArtist).To(Equal(newAlbum.DisplayArtist))
        Expect(res.Albums[0].Type).To(Equal(newAlbum.Type))
        Expect(res.Albums[0].About).To(Equal(newAlbum.About))
        Expect(res.Albums[0].TotalTracks).To(Equal(int32(len(newAlbum.Tracks))))
        Expect(res.Albums[0].UserGroup).NotTo(BeNil())
        Expect(res.Albums[0].UserGroup.Id).To(Equal(newArtistUserGroup.Id.String()))
        Expect(res.Albums[0].UserGroup.DisplayName).To(Equal(newArtistUserGroup.DisplayName))
        Expect(res.Albums[0].UserGroup.Avatar).To(Equal(newArtistUserGroup.Avatar))
      })
    })
    Context("with invalid query", func() {
      It("should respond with invalid error", func() {
        q := &tagpb.Query{Query: "ti"}
        resp, err := service.SearchTrackGroups(context.Background(), q)

        Expect(resp).To(BeNil())
        Expect(err).To(HaveOccurred())

        twerr := err.(twirp.Error)
        Expect(twerr.Code()).To(Equal(invalid_argument_code))
        Expect(twerr.Meta("argument")).To(Equal("query"))
      })
    })
  })

  Describe("UpdateTrackGroup", func() {
    Context("with valid uuid", func() {
			It("should update track group if it exists", func() {
        cover := make([]byte, 5)
        cover[0] = '1'
        releaseDate, err := ptypes.TimestampProto(time.Now())
        Expect(err).NotTo(HaveOccurred())
				trackGroup := &pb.TrackGroup{
          Id: newAlbum.Id.String(),
					Title: "best album ever",
					CreatorId: newUser.Id.String(),
					UserGroupId: artist.Id.String(),
          LabelId: label.Id.String(),
          ReleaseDate: releaseDate,
          Type: "ep",
          Cover: cover,
          DisplayArtist: "Various",
          MultipleComposers: true,
          Private: false,
          About: "about",
          Tags: []*tagpb.Tag{
            &tagpb.Tag{
              Type: "genre",
              Name: "rock",
            },
          },
				}
				_, err = service.UpdateTrackGroup(context.Background(), trackGroup)

				Expect(err).NotTo(HaveOccurred())

				t := new(model.TrackGroup)
				err = db.Model(t).Where("id = ?", newAlbum.Id).Select()
				Expect(err).NotTo(HaveOccurred())
				Expect(t.Title).To(Equal(trackGroup.Title))
        updatedReleaseDate, err := ptypes.TimestampProto(t.ReleaseDate)
        Expect(err).NotTo(HaveOccurred())
				Expect(updatedReleaseDate.Seconds).To(Equal(trackGroup.ReleaseDate.Seconds))
				Expect(t.DisplayArtist).To(Equal(trackGroup.DisplayArtist))
        Expect(t.Cover).To(Equal(trackGroup.Cover))
				Expect(t.About).To(Equal(trackGroup.About))
        Expect(t.MultipleComposers).To(Equal(trackGroup.MultipleComposers))
				Expect(t.Private).To(Equal(trackGroup.Private))
        Expect(t.Type).To(Equal(trackGroup.Type))

        Expect(len(t.Tags)).To(Equal(1))
        addedTag := model.Tag{Id: t.Tags[0]}
        err = db.Model(&addedTag).WherePK().Returning("*").Select()
        Expect(addedTag.Type).To(Equal("genre"))
        Expect(addedTag.Name).To(Equal("rock"))

        Expect(t.UserGroupId).To(Equal(artist.Id))
        Expect(t.LabelId).To(Equal(label.Id))

        var tracksInAlbum []model.Track
        err = db.Model(&tracksInAlbum).Where("id IN (?)", pg.In(t.Tracks)).Select()
        Expect(err).NotTo(HaveOccurred())
        for _, track := range tracksInAlbum {
          Expect(track.UserGroupId).To(Equal(artist.Id))
        }

				// unchanged
				Expect(t.CreatorId).To(Equal(newAlbum.CreatorId))
				Expect(len(t.Tracks)).To(Equal(1))
			})
			It("should respond with not_found error if track group does not exist", func() {
				id := uuid.NewV4()
				for id == newAlbum.Id || id == newPlaylist.Id {
					id = uuid.NewV4()
				}
        cover := make([]byte, 5)
        cover[0] = '1'
        releaseDate, err := ptypes.TimestampProto(time.Now())
        Expect(err).NotTo(HaveOccurred())
				trackGroup := &pb.TrackGroup{
          Id: id.String(),
					Title: "best album ever",
					CreatorId: newUser.Id.String(),
					UserGroupId: newArtistUserGroup.Id.String(),
          LabelId: newLabelUserGroup.Id.String(),
          ReleaseDate: releaseDate,
          Type: "lp",
          Cover: cover,
          DisplayArtist: "Various",
          MultipleComposers: true,
          Private: false,
				}
				resp, err := service.UpdateTrackGroup(context.Background(), trackGroup)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(not_found_code))
			})
      It("should respond with not_found error if user_group does not exist", func() {
        userGroupId := uuid.NewV4()
        for userGroupId == newLabelUserGroup.Id || userGroupId == newArtistUserGroup.Id || userGroupId == artist.Id || userGroupId == label.Id {
          userGroupId = uuid.NewV4()
        }
        cover := make([]byte, 5)
        cover[0] = '1'
        releaseDate, err := ptypes.TimestampProto(time.Now())
        Expect(err).NotTo(HaveOccurred())
        trackGroup := &pb.TrackGroup{
          Id: newAlbum.Id.String(),
          Title: "best album ever",
          CreatorId: newUser.Id.String(),
          UserGroupId: userGroupId.String(),
          LabelId: label.Id.String(),
          ReleaseDate: releaseDate,
          Type: "ep",
          Cover: cover,
          DisplayArtist: "Various",
          MultipleComposers: true,
          Private: false,
          About: "about",
          Tags: []*tagpb.Tag{
            &tagpb.Tag{
              Type: "genre",
              Name: "rock",
            },
          },
        }
        resp, err := service.UpdateTrackGroup(context.Background(), trackGroup)

        Expect(resp).To(BeNil())
        Expect(err).To(HaveOccurred())
        twerr := err.(twirp.Error)
        Expect(twerr.Code()).To(Equal(not_found_code))
      })
      It("should respond with not_found error if label does not exist", func() {
        userGroupId := uuid.NewV4()
        for userGroupId == newLabelUserGroup.Id || userGroupId == newArtistUserGroup.Id || userGroupId == artist.Id || userGroupId == label.Id {
          userGroupId = uuid.NewV4()
        }
        cover := make([]byte, 5)
        cover[0] = '1'
        releaseDate, err := ptypes.TimestampProto(time.Now())
        Expect(err).NotTo(HaveOccurred())
        trackGroup := &pb.TrackGroup{
          Id: newAlbum.Id.String(),
          Title: "best album ever",
          CreatorId: newUser.Id.String(),
          UserGroupId: artist.Id.String(),
          LabelId: userGroupId.String(),
          ReleaseDate: releaseDate,
          Type: "ep",
          Cover: cover,
          DisplayArtist: "Various",
          MultipleComposers: true,
          Private: false,
          About: "about",
          Tags: []*tagpb.Tag{
            &tagpb.Tag{
              Type: "genre",
              Name: "rock",
            },
          },
        }
        resp, err := service.UpdateTrackGroup(context.Background(), trackGroup)

        Expect(resp).To(BeNil())
        Expect(err).To(HaveOccurred())
        twerr := err.(twirp.Error)
        Expect(twerr.Code()).To(Equal(not_found_code))
      })
		})
		Context("with invalid uuid", func() {
			It("should respond with invalid_argument error", func() {
				id := "45"
        cover := make([]byte, 5)
        cover[0] = '1'
        releaseDate, err := ptypes.TimestampProto(time.Now())
        Expect(err).NotTo(HaveOccurred())
				trackGroup := &pb.TrackGroup{
          Id: id,
					Title: "best album ever",
					CreatorId: newUser.Id.String(),
					UserGroupId: newArtistUserGroup.Id.String(),
          LabelId: newLabelUserGroup.Id.String(),
          ReleaseDate: releaseDate,
          Type: "lp",
          Cover: cover,
          DisplayArtist: "Various",
          MultipleComposers: true,
          Private: false,
				}
				resp, err := service.UpdateTrackGroup(context.Background(), trackGroup)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
				Expect(twerr.Meta("argument")).To(Equal("id"))
			})
      It("should not update track group if user_group_id is invalid", func() {
        cover := make([]byte, 5)
        cover[0] = '1'
        releaseDate, err := ptypes.TimestampProto(time.Now())
        Expect(err).NotTo(HaveOccurred())
        trackGroup := &pb.TrackGroup{
          Id: newAlbum.Id.String(),
          Title: "best album ever",
          CreatorId: newUser.Id.String(),
          UserGroupId: "2",
          LabelId: label.Id.String(),
          ReleaseDate: releaseDate,
          Type: "ep",
          Cover: cover,
          DisplayArtist: "Various",
          MultipleComposers: true,
          Private: false,
          About: "about",
        }
        resp, err := service.UpdateTrackGroup(context.Background(), trackGroup)

        Expect(resp).To(BeNil())
        Expect(err).To(HaveOccurred())
        twerr := err.(twirp.Error)
        Expect(twerr.Code()).To(Equal(invalid_argument_code))
        Expect(twerr.Meta("argument")).To(Equal("user_group id"))
      })
      It("should not update track group if label_id is invalid", func() {
        cover := make([]byte, 5)
        cover[0] = '1'
        releaseDate, err := ptypes.TimestampProto(time.Now())
        Expect(err).NotTo(HaveOccurred())
        trackGroup := &pb.TrackGroup{
          Id: newAlbum.Id.String(),
          Title: "best album ever",
          CreatorId: newUser.Id.String(),
          UserGroupId: artist.Id.String(),
          LabelId: "1",
          ReleaseDate: releaseDate,
          Type: "ep",
          Cover: cover,
          DisplayArtist: "Various",
          MultipleComposers: true,
          Private: false,
          About: "about",
        }
        resp, err := service.UpdateTrackGroup(context.Background(), trackGroup)

        Expect(resp).To(BeNil())
        Expect(err).To(HaveOccurred())
        twerr := err.(twirp.Error)
        Expect(twerr.Code()).To(Equal(invalid_argument_code))
        Expect(twerr.Meta("argument")).To(Equal("user_group id"))
      })
		})
  })

  Describe("AddTracksToTrackGroup", func() {
    Context("with track_group_id and track ids", func() {
      It("should add tracks to track group", func() {
        tracksToTrackGroup := &pb.TracksToTrackGroup{
          TrackGroupId: newPlaylist.Id.String(),
          Tracks: []*trackpb.Track{
            &trackpb.Track{Id: playlistTrack.Id.String()},
          },
        }
        _, err := service.AddTracksToTrackGroup(context.Background(), tracksToTrackGroup)

        trackGroup := new(model.TrackGroup)
        err = db.Model(trackGroup).Where("id = ?", newPlaylist.Id).Select()
        Expect(err).NotTo(HaveOccurred())
        Expect(trackGroup.Tracks).To(HaveLen(2))
        Expect(trackGroup.Tracks).To(ContainElement(playlistTrack.Id))

        var tracks []model.Track
        err = db.Model(&tracks).
          Where("id in (?)", pg.In([]uuid.UUID{playlistTrack.Id})).
          Select()
        Expect(err).NotTo(HaveOccurred())
        for _, track := range tracks {
          Expect(track.TrackGroups).To(ContainElement(newPlaylist.Id))
        }
      })
      It("should respond with not_found error if track group does not exist", func() {
        id := uuid.NewV4()
        for id == newPlaylist.Id || id == newAlbum.Id {
          id = uuid.NewV4()
        }
        tracksToTrackGroup := &pb.TracksToTrackGroup{
          TrackGroupId: id.String(),
          Tracks: []*trackpb.Track{
            &trackpb.Track{Id: playlistTrack.Id.String()},
          },
        }
        resp, err := service.AddTracksToTrackGroup(context.Background(), tracksToTrackGroup)

        Expect(resp).To(BeNil())
        Expect(err).To(HaveOccurred())

        twerr := err.(twirp.Error)
        Expect(twerr.Code()).To(Equal(not_found_code))
      })
      It("should respond with not_found error if one of the tracks does not exist", func() {
        id := uuid.NewV4()
        for id == newTrack.Id || id == playlistTrack.Id {
          id = uuid.NewV4()
        }
        tracksToTrackGroup := &pb.TracksToTrackGroup{
          TrackGroupId: newPlaylist.Id.String(),
          Tracks: []*trackpb.Track{
            &trackpb.Track{Id: id.String()},
          },
        }
        resp, err := service.AddTracksToTrackGroup(context.Background(), tracksToTrackGroup)

        Expect(resp).To(BeNil())
        Expect(err).To(HaveOccurred())

        twerr := err.(twirp.Error)
        Expect(twerr.Code()).To(Equal(not_found_code))
      })
    })
    Context("with invalid user_group_id", func() {
      It("should respond with invalid_argument_code error", func() {
				tracksToTrackGroup := &pb.TracksToTrackGroup{
          TrackGroupId: "1",
          Tracks: []*trackpb.Track{
            &trackpb.Track{Id: playlistTrack.Id.String()},
          },
        }
				resp, err := service.AddTracksToTrackGroup(context.Background(), tracksToTrackGroup)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
			})
    })
    Context("with invalid track id", func() {
      It("should respond with invalid_argument_code error", func() {
        tracksToTrackGroup := &pb.TracksToTrackGroup{
          TrackGroupId: newPlaylist.Id.String(),
          Tracks: []*trackpb.Track{
            &trackpb.Track{Id: ""},
          },
        }
        resp, err := service.AddTracksToTrackGroup(context.Background(), tracksToTrackGroup)

        Expect(resp).To(BeNil())
        Expect(err).To(HaveOccurred())

        twerr := err.(twirp.Error)
        Expect(twerr.Code()).To(Equal(invalid_argument_code))
      })
    })
  })

  Describe("RemoveTracksFromTrackGroup", func() {
    Context("with track_group_id and track ids", func() {
      It("should add tracks to track group", func() {
        tracksToTrackGroup := &pb.TracksToTrackGroup{
          TrackGroupId: newPlaylist.Id.String(),
          Tracks: []*trackpb.Track{
            &trackpb.Track{Id: playlistTrack.Id.String()},
          },
        }
        _, err := service.RemoveTracksFromTrackGroup(context.Background(), tracksToTrackGroup)

        trackGroup := new(model.TrackGroup)
        err = db.Model(trackGroup).Where("id = ?", newPlaylist.Id).Select()
        Expect(err).NotTo(HaveOccurred())
        Expect(trackGroup.Tracks).To(HaveLen(1))
        Expect(trackGroup.Tracks).NotTo(ContainElement(playlistTrack.Id))

        var tracks []model.Track
        err = db.Model(&tracks).
          Where("id in (?)", pg.In([]uuid.UUID{playlistTrack.Id})).
          Select()
        Expect(err).NotTo(HaveOccurred())
        for _, track := range tracks {
          Expect(track.TrackGroups).NotTo(ContainElement(newPlaylist.Id))
        }
      })
      It("should respond with not_found error if track group does not exist", func() {
        id := uuid.NewV4()
        for id == newPlaylist.Id || id == newAlbum.Id {
          id = uuid.NewV4()
        }

        tracksToTrackGroup := &pb.TracksToTrackGroup{
          TrackGroupId: id.String(),
          Tracks: []*trackpb.Track{
            &trackpb.Track{Id: playlistTrack.Id.String()},
          },
        }
        resp, err := service.RemoveTracksFromTrackGroup(context.Background(), tracksToTrackGroup)

        Expect(resp).To(BeNil())
        Expect(err).To(HaveOccurred())

        twerr := err.(twirp.Error)
        Expect(twerr.Code()).To(Equal(not_found_code))
      })
      It("should respond with not_found error if one of the tracks does not exist", func() {
        id := uuid.NewV4()
        for id == newTrack.Id || id == playlistTrack.Id {
          id = uuid.NewV4()
        }
        tracksToTrackGroup := &pb.TracksToTrackGroup{
          TrackGroupId: newPlaylist.Id.String(),
          Tracks: []*trackpb.Track{
            &trackpb.Track{Id: id.String()},
          },
        }
        resp, err := service.RemoveTracksFromTrackGroup(context.Background(), tracksToTrackGroup)

        Expect(resp).To(BeNil())
        Expect(err).To(HaveOccurred())

        twerr := err.(twirp.Error)
        Expect(twerr.Code()).To(Equal(not_found_code))
      })
    })
    Context("with invalid user_group_id", func() {
      It("should respond with invalid_argument_code error", func() {
        tracksToTrackGroup := &pb.TracksToTrackGroup{
          TrackGroupId: "1",
          Tracks: []*trackpb.Track{
            &trackpb.Track{Id: playlistTrack.Id.String()},
          },
        }
        resp, err := service.RemoveTracksFromTrackGroup(context.Background(), tracksToTrackGroup)

        Expect(resp).To(BeNil())
        Expect(err).To(HaveOccurred())

        twerr := err.(twirp.Error)
        Expect(twerr.Code()).To(Equal(invalid_argument_code))
      })
    })
    Context("with invalid track id", func() {
      It("should respond with invalid_argument_code error", func() {
        tracksToTrackGroup := &pb.TracksToTrackGroup{
          TrackGroupId: newPlaylist.Id.String(),
          Tracks: []*trackpb.Track{
            &trackpb.Track{Id: ""},
          },
        }
        resp, err := service.RemoveTracksFromTrackGroup(context.Background(), tracksToTrackGroup)

        Expect(resp).To(BeNil())
        Expect(err).To(HaveOccurred())

        twerr := err.(twirp.Error)
        Expect(twerr.Code()).To(Equal(invalid_argument_code))
      })
    })
  })

  Describe("CreateTrackGroup", func() {
		Context("with all required attributes", func() {
			It("should create a new track group (not playlist)", func() {
        cover := make([]byte, 5)
        releaseDate, err := ptypes.TimestampProto(time.Now())
        Expect(err).NotTo(HaveOccurred())
				trackGroup := &pb.TrackGroup{
					Title: "best album ever",
					CreatorId: newUser.Id.String(),
					UserGroupId: newArtistUserGroup.Id.String(),
          LabelId: newLabelUserGroup.Id.String(),
          ReleaseDate: releaseDate,
          Type: "lp",
          Cover: cover,
          About: "about",
          DisplayArtist: "Various",
          MultipleComposers: true,
          Private: false,
					Tags: []*tagpb.Tag{
						&tagpb.Tag{
							Type: "genre",
							Name: "rock",
						},
					},
          Tracks: []*trackpb.Track{
            &trackpb.Track{
              Id: newTrack.Id.String(),
            },
          },
				}
				resp, err := service.CreateTrackGroup(context.Background(), trackGroup)

				Expect(err).NotTo(HaveOccurred())
				Expect(resp).NotTo(BeNil())

				Expect(resp.Id).NotTo(BeNil())
				Expect(resp.Title).To(Equal(trackGroup.Title))
				Expect(resp.CreatorId).To(Equal(trackGroup.CreatorId))
        Expect(resp.UserGroupId).To(Equal(trackGroup.UserGroupId))
        Expect(resp.LabelId).To(Equal(trackGroup.LabelId))
        Expect(resp.ReleaseDate).To(Equal(trackGroup.ReleaseDate))
        Expect(resp.Type).To(Equal(trackGroup.Type))
        Expect(resp.About).To(Equal(trackGroup.About))
        Expect(resp.Cover).To(Equal(trackGroup.Cover))
        Expect(resp.DisplayArtist).To(Equal(trackGroup.DisplayArtist))
        Expect(resp.MultipleComposers).To(Equal(trackGroup.MultipleComposers))
				Expect(resp.Private).To(Equal(trackGroup.Private))

				Expect(len(resp.Tags)).To(Equal(1))
				Expect(resp.Tags[0].Id).NotTo(Equal(""))
				Expect(resp.Tags[0].Type).To(Equal("genre"))
				Expect(resp.Tags[0].Name).To(Equal("rock"))

				Expect(len(resp.Tracks)).To(Equal(1))
				Expect(resp.Tracks[0].Id).To(Equal(newTrack.Id.String()))
				Expect(resp.Tracks[0].Title).To(Equal(newTrack.Title))
        Expect(resp.Tracks[0].TrackServerId).To(Equal(newTrack.TrackServerId.String()))
        Expect(resp.Tracks[0].Duration).To(Equal(newTrack.Duration))
        Expect(resp.Tracks[0].TrackNumber).To(Equal(newTrack.TrackNumber))
        Expect(resp.Tracks[0].Status).To(Equal(newTrack.Status))
        Expect(len(resp.Tracks[0].Artists)).To(Equal(1))
        Expect(resp.Tracks[0].Artists[0].Id).To(Equal(newArtistUserGroup.Id.String()))
        Expect(resp.Tracks[0].Artists[0].DisplayName).To(Equal(newArtistUserGroup.DisplayName))
				Expect(resp.Tracks[0].Artists[0].Avatar).To(Equal(newArtistUserGroup.Avatar))

        // artist := new(model.UserGroup)
        // err = db.Model(artist).Where("id = ?", newArtistUserGroup.Id).Select()
        // Expect(err).NotTo(HaveOccurred())
        // label := new(model.UserGroup)
        // err = db.Model(label).Where("id = ?", newLabelUserGroup.Id).Select()
        // Expect(err).NotTo(HaveOccurred())

        // Expect(artist.TrackGroups).To(ContainElement(trackGroupId))
				// Expect(label.TrackGroups).To(ContainElement(trackGroupId))

        trackGroupId, err := uuid.FromString(resp.Id)
				Expect(err).NotTo(HaveOccurred())
        track := new(model.Track)
        err = db.Model(track).Where("id = ?", newTrack.Id).Select()
        Expect(err).NotTo(HaveOccurred())
        Expect(len(track.TrackGroups)).To(Equal(3))
        Expect(track.TrackGroups).To(ContainElement(trackGroupId))

        newTrackGroup := new(model.TrackGroup)
        err = db.Model(newTrackGroup).Where("id = ?", trackGroupId).Select()
        Expect(newTrackGroup.Tracks).To(ContainElement(newTrack.Id))

        createdReleaseDate, err := ptypes.TimestampProto(newTrackGroup.ReleaseDate)
        Expect(err).NotTo(HaveOccurred())
        Expect(createdReleaseDate.Seconds).To(Equal(releaseDate.Seconds))
			})
      It("should create a new track group (playlist)", func() {
        cover := make([]byte, 5)
        releaseDate, err := ptypes.TimestampProto(time.Now())
        Expect(err).NotTo(HaveOccurred())
				trackGroup := &pb.TrackGroup{
					Title: "best playlist ever",
					CreatorId: newUser.Id.String(),
          ReleaseDate: releaseDate,
          Type: "playlist",
          Cover: cover,
          Private: true,
					Tags: []*tagpb.Tag{
						&tagpb.Tag{
							Type: "genre",
							Name: "rock",
						},
					},
          Tracks: []*trackpb.Track{
            &trackpb.Track{
              Id: newTrack.Id.String(),
            },
          },
				}
				resp, err := service.CreateTrackGroup(context.Background(), trackGroup)

				Expect(err).NotTo(HaveOccurred())
				Expect(resp).NotTo(BeNil())

				Expect(resp.Id).NotTo(BeNil())
				Expect(resp.Title).To(Equal(trackGroup.Title))
				Expect(resp.CreatorId).To(Equal(trackGroup.CreatorId))
        Expect(resp.UserGroupId).To(Equal(trackGroup.UserGroupId))
        Expect(resp.LabelId).To(Equal(trackGroup.LabelId))
        Expect(resp.ReleaseDate).To(Equal(trackGroup.ReleaseDate))
        Expect(resp.Type).To(Equal(trackGroup.Type))
        Expect(resp.Cover).To(Equal(trackGroup.Cover))
        Expect(resp.DisplayArtist).To(Equal(trackGroup.DisplayArtist))
        Expect(resp.MultipleComposers).To(Equal(trackGroup.MultipleComposers))
				Expect(resp.Private).To(Equal(trackGroup.Private))

				Expect(len(resp.Tags)).To(Equal(1))
				Expect(resp.Tags[0].Id).NotTo(Equal(""))
				Expect(resp.Tags[0].Type).To(Equal("genre"))
				Expect(resp.Tags[0].Name).To(Equal("rock"))

				Expect(len(resp.Tracks)).To(Equal(1))
				Expect(resp.Tracks[0].Id).To(Equal(newTrack.Id.String()))
				Expect(resp.Tracks[0].Title).To(Equal(newTrack.Title))
        Expect(resp.Tracks[0].TrackServerId).To(Equal(newTrack.TrackServerId.String()))
        Expect(resp.Tracks[0].Duration).To(Equal(newTrack.Duration))
        Expect(resp.Tracks[0].TrackNumber).To(Equal(newTrack.TrackNumber))
        Expect(resp.Tracks[0].Status).To(Equal(newTrack.Status))
        Expect(len(resp.Tracks[0].Artists)).To(Equal(1))
        Expect(resp.Tracks[0].Artists[0].Id).To(Equal(newArtistUserGroup.Id.String()))
        Expect(resp.Tracks[0].Artists[0].DisplayName).To(Equal(newArtistUserGroup.DisplayName))
				Expect(resp.Tracks[0].Artists[0].Avatar).To(Equal(newArtistUserGroup.Avatar))

				user := new(model.User)
				err = db.Model(user).Where("id = ?", newUser.Id).Select()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(user.Playlists)).To(Equal(2))

				trackGroupId, err := uuid.FromString(resp.Id)
				Expect(err).NotTo(HaveOccurred())
        Expect(user.Playlists).To(ContainElement(trackGroupId))

        track := new(model.Track)
        err = db.Model(track).Where("id = ?", newTrack.Id).Select()
        Expect(err).NotTo(HaveOccurred())
        Expect(len(track.TrackGroups)).To(Equal(4))
        Expect(track.TrackGroups).To(ContainElement(trackGroupId))

        newTrackGroup := new(model.TrackGroup)
        err = db.Model(newTrackGroup).Where("id = ?", trackGroupId).Select()
        Expect(newTrackGroup.Tracks).To(ContainElement(newTrack.Id))
			})
		})

		Context("with missing required attributes", func() {
			It("should not create a track group without title", func() {
        cover := make([]byte, 5)
        releaseDate, err := ptypes.TimestampProto(time.Now())
        Expect(err).NotTo(HaveOccurred())
				trackGroup := &pb.TrackGroup{
					Title: "",
					CreatorId: newUser.Id.String(),
          ReleaseDate: releaseDate,
          Type: "playlist",
          Cover: cover,
          Private: true,
					Tags: []*tagpb.Tag{
						&tagpb.Tag{
							Type: "genre",
							Name: "rock",
						},
					},
          Tracks: []*trackpb.Track{
            &trackpb.Track{
              Id: newTrack.Id.String(),
            },
          },
				}
				resp, err := service.CreateTrackGroup(context.Background(), trackGroup)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
				Expect(twerr.Meta("argument")).To(Equal("title"))
			})
			It("should not create a track group without release date", func() {
        cover := make([]byte, 5)
        // releaseDate, err := ptypes.TimestampProto(time.Now())
        // Expect(err).NotTo(HaveOccurred())
				trackGroup := &pb.TrackGroup{
					Title: "best playlist ever",
					CreatorId: newUser.Id.String(),
          // ReleaseDate: releaseDate,
          Type: "playlist",
          Cover: cover,
          Private: true,
					Tags: []*tagpb.Tag{
						&tagpb.Tag{
							Type: "genre",
							Name: "rock",
						},
					},
          Tracks: []*trackpb.Track{
            &trackpb.Track{
              Id: newTrack.Id.String(),
            },
          },
				}
				resp, err := service.CreateTrackGroup(context.Background(), trackGroup)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())
				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
				Expect(twerr.Meta("argument")).To(Equal("release_date"))
			})
			It("should not create a track group without creator_id", func() {
        cover := make([]byte, 5)
        releaseDate, err := ptypes.TimestampProto(time.Now())
        Expect(err).NotTo(HaveOccurred())
				trackGroup := &pb.TrackGroup{
					Title: "best playlist ever",
					CreatorId: "",
          ReleaseDate: releaseDate,
          Type: "playlist",
          Cover: cover,
          Private: true,
					Tags: []*tagpb.Tag{
						&tagpb.Tag{
							Type: "genre",
							Name: "rock",
						},
					},
          Tracks: []*trackpb.Track{
            &trackpb.Track{
              Id: newTrack.Id.String(),
            },
          },
				}
				resp, err := service.CreateTrackGroup(context.Background(), trackGroup)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
				Expect(twerr.Meta("argument")).To(Equal("creator_id"))
			})
			It("should not create a track group without type", func() {
        cover := make([]byte, 5)
        releaseDate, err := ptypes.TimestampProto(time.Now())
        Expect(err).NotTo(HaveOccurred())
				trackGroup := &pb.TrackGroup{
					Title: "best playlist ever",
					CreatorId: newUser.Id.String(),
          ReleaseDate: releaseDate,
          Type: "",
          Cover: cover,
          Private: true,
					Tags: []*tagpb.Tag{
						&tagpb.Tag{
							Type: "genre",
							Name: "rock",
						},
					},
          Tracks: []*trackpb.Track{
            &trackpb.Track{
              Id: newTrack.Id.String(),
            },
          },
				}
				resp, err := service.CreateTrackGroup(context.Background(), trackGroup)
				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())
				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
				Expect(twerr.Meta("argument")).To(Equal("type"))
			})
			It("should not create a track group without cover", func() {
        releaseDate, err := ptypes.TimestampProto(time.Now())
        Expect(err).NotTo(HaveOccurred())
				trackGroup := &pb.TrackGroup{
					Title: "best playlist ever",
					CreatorId: newUser.Id.String(),
          ReleaseDate: releaseDate,
          Type: "playlist",
          // Cover: cover,
          Private: true,
					Tags: []*tagpb.Tag{
						&tagpb.Tag{
							Type: "genre",
							Name: "rock",
						},
					},
          Tracks: []*trackpb.Track{
            &trackpb.Track{
              Id: newTrack.Id.String(),
            },
          },
				}
				resp, err := service.CreateTrackGroup(context.Background(), trackGroup)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())
				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
				Expect(twerr.Meta("argument")).To(Equal("cover"))
			})
      It("should not create a track group (release) without user group id", func() {
        cover := make([]byte, 5)
        releaseDate, err := ptypes.TimestampProto(time.Now())
        Expect(err).NotTo(HaveOccurred())
				trackGroup := &pb.TrackGroup{
					Title: "best album ever",
					CreatorId: newUser.Id.String(),
					// UserGroupId: newArtistUserGroup.Id.String(),
          LabelId: newLabelUserGroup.Id.String(),
          ReleaseDate: releaseDate,
          Type: "lp",
          Cover: cover,
          DisplayArtist: "Various",
          MultipleComposers: true,
          Private: false,
					Tags: []*tagpb.Tag{
						&tagpb.Tag{
							Type: "genre",
							Name: "rock",
						},
					},
          Tracks: []*trackpb.Track{
            &trackpb.Track{
              Id: newTrack.Id.String(),
            },
          },
				}
        resp, err := service.CreateTrackGroup(context.Background(), trackGroup)

        Expect(resp).To(BeNil())
        Expect(err).To(HaveOccurred())
        twerr := err.(twirp.Error)
        Expect(twerr.Code()).To(Equal(invalid_argument_code))
        Expect(twerr.Meta("argument")).To(Equal("user_group_id"))
      })
		})

		Context("with invalid attributes", func() {
			It("should not create a track group if creator_id is invalid", func() {
        cover := make([]byte, 5)
        releaseDate, err := ptypes.TimestampProto(time.Now())
        Expect(err).NotTo(HaveOccurred())
				trackGroup := &pb.TrackGroup{
					Title: "best album ever",
					CreatorId: "123",
					UserGroupId: newArtistUserGroup.Id.String(),
          LabelId: newLabelUserGroup.Id.String(),
          ReleaseDate: releaseDate,
          Type: "lp",
          Cover: cover,
          DisplayArtist: "Various",
          MultipleComposers: true,
          Private: false,
					Tags: []*tagpb.Tag{
						&tagpb.Tag{
							Type: "genre",
							Name: "rock",
						},
					},
          Tracks: []*trackpb.Track{
            &trackpb.Track{
              Id: newTrack.Id.String(),
            },
          },
				}
				resp, err := service.CreateTrackGroup(context.Background(), trackGroup)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())
				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
			})
			It("should not create a track group if user_group_id is invalid", func() {
        cover := make([]byte, 5)
        releaseDate, err := ptypes.TimestampProto(time.Now())
        Expect(err).NotTo(HaveOccurred())
				trackGroup := &pb.TrackGroup{
					Title: "best album ever",
					CreatorId: newUser.Id.String(),
					UserGroupId: "123",
          LabelId: newLabelUserGroup.Id.String(),
          ReleaseDate: releaseDate,
          Type: "lp",
          Cover: cover,
          DisplayArtist: "Various",
          MultipleComposers: true,
          Private: false,
					Tags: []*tagpb.Tag{
						&tagpb.Tag{
							Type: "genre",
							Name: "rock",
						},
					},
          Tracks: []*trackpb.Track{
            &trackpb.Track{
              Id: newTrack.Id.String(),
            },
          },
				}
				resp, err := service.CreateTrackGroup(context.Background(), trackGroup)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())
				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
			})
      It("should not create a track group if label_id is invalid", func() {
        cover := make([]byte, 5)
        releaseDate, err := ptypes.TimestampProto(time.Now())
        Expect(err).NotTo(HaveOccurred())
				trackGroup := &pb.TrackGroup{
					Title: "best album ever",
					CreatorId: newUser.Id.String(),
					UserGroupId: newArtistUserGroup.Id.String(),
          LabelId: "456",
          ReleaseDate: releaseDate,
          Type: "lp",
          Cover: cover,
          DisplayArtist: "Various",
          MultipleComposers: true,
          Private: false,
					Tags: []*tagpb.Tag{
						&tagpb.Tag{
							Type: "genre",
							Name: "rock",
						},
					},
          Tracks: []*trackpb.Track{
            &trackpb.Track{
              Id: newTrack.Id.String(),
            },
          },
				}
				resp, err := service.CreateTrackGroup(context.Background(), trackGroup)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())
				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
			})
			It("should not create a track group if creator does not exist", func() {
				userId := uuid.NewV4()
				for userId == newUser.Id {
					userId = uuid.NewV4()
				}
        cover := make([]byte, 5)
        releaseDate, err := ptypes.TimestampProto(time.Now())
        Expect(err).NotTo(HaveOccurred())
				trackGroup := &pb.TrackGroup{
					Title: "best album ever",
					CreatorId: userId.String(),
					UserGroupId: newArtistUserGroup.Id.String(),
          LabelId: newLabelUserGroup.Id.String(),
          ReleaseDate: releaseDate,
          Type: "lp",
          Cover: cover,
          DisplayArtist: "Various",
          MultipleComposers: true,
          Private: false,
					Tags: []*tagpb.Tag{
						&tagpb.Tag{
							Type: "genre",
							Name: "rock",
						},
					},
          Tracks: []*trackpb.Track{
            &trackpb.Track{
              Id: newTrack.Id.String(),
            },
          },
				}
				resp, err := service.CreateTrackGroup(context.Background(), trackGroup)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())
				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(not_found_code))
			})
			It("should not create a track group if user_group does not exist", func() {
				userGroupId := uuid.NewV4()
				for userGroupId == newLabelUserGroup.Id || userGroupId == newArtistUserGroup.Id || userGroupId == artist.Id || userGroupId == label.Id {
					userGroupId = uuid.NewV4()
				}
        cover := make([]byte, 5)
        releaseDate, err := ptypes.TimestampProto(time.Now())
        Expect(err).NotTo(HaveOccurred())
				trackGroup := &pb.TrackGroup{
					Title: "best album ever",
					CreatorId: newUser.Id.String(),
					UserGroupId: userGroupId.String(),
          LabelId: newLabelUserGroup.Id.String(),
          ReleaseDate: releaseDate,
          Type: "lp",
          Cover: cover,
          DisplayArtist: "Various",
          MultipleComposers: true,
          Private: false,
					Tags: []*tagpb.Tag{
						&tagpb.Tag{
							Type: "genre",
							Name: "rock",
						},
					},
          Tracks: []*trackpb.Track{
            &trackpb.Track{
              Id: newTrack.Id.String(),
            },
          },
				}
				resp, err := service.CreateTrackGroup(context.Background(), trackGroup)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())
				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(not_found_code))
			})
      It("should not create a track group if label does not exist", func() {
				userGroupId := uuid.NewV4()
				for userGroupId == newLabelUserGroup.Id || userGroupId == newArtistUserGroup.Id || userGroupId == artist.Id || userGroupId == label.Id {
					userGroupId = uuid.NewV4()
				}
        cover := make([]byte, 5)
        releaseDate, err := ptypes.TimestampProto(time.Now())
        Expect(err).NotTo(HaveOccurred())
				trackGroup := &pb.TrackGroup{
					Title: "best album ever",
					CreatorId: newUser.Id.String(),
					LabelId: userGroupId.String(),
          UserGroupId: newLabelUserGroup.Id.String(),
          ReleaseDate: releaseDate,
          Type: "lp",
          Cover: cover,
          DisplayArtist: "Various",
          MultipleComposers: true,
          Private: false,
					Tags: []*tagpb.Tag{
						&tagpb.Tag{
							Type: "genre",
							Name: "rock",
						},
					},
          Tracks: []*trackpb.Track{
            &trackpb.Track{
              Id: newTrack.Id.String(),
            },
          },
				}
				resp, err := service.CreateTrackGroup(context.Background(), trackGroup)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())
				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(not_found_code))
			})
		})
	})

  Describe("DeleteTrackGroup", func() {
    Context("with valid uuid", func() {
      It("should delete trackGroup (playlist) if it exists and remove it from user playlists", func() {
        trackGroup := &pb.TrackGroup{Id: newPlaylist.Id.String()}

        trackGroupToDelete := new(model.TrackGroup)
        err := db.Model(trackGroupToDelete).Where("id = ?", newPlaylist.Id).Select()
        Expect(err).NotTo(HaveOccurred())

        _, err = service.DeleteTrackGroup(context.Background(), trackGroup)

        Expect(err).NotTo(HaveOccurred())

        user := new(model.User)
        err = db.Model(user).Where("id = ?", trackGroupToDelete.CreatorId).Select()
        Expect(err).NotTo(HaveOccurred())
        Expect(user.Playlists).NotTo(ContainElement(trackGroupToDelete.Id))

        var tracks []model.Track
        err = db.Model(&tracks).
          Where("id in (?)", pg.In(trackGroupToDelete.Tracks)).
          Select()
        Expect(err).NotTo(HaveOccurred())
        Expect(len(tracks)).To(Equal(1))
        for _, track := range tracks {
          Expect(track.TrackGroups).NotTo(ContainElement(trackGroupToDelete.Id))
        }

        var trackGroups []model.Track
        err = db.Model(&trackGroups).
          Where("id in (?)", pg.In([]uuid.UUID{trackGroupToDelete.Id})).
          Select()
        Expect(err).NotTo(HaveOccurred())
        Expect(len(trackGroups)).To(Equal(0))
      })
      It("should delete trackGroup (release) if it exists and associated tracks", func() {
        trackGroup := &pb.TrackGroup{Id: newAlbum.Id.String()}

        trackGroupToDelete := new(model.TrackGroup)
        err := db.Model(trackGroupToDelete).Where("id = ?", newAlbum.Id).Select()
        Expect(err).NotTo(HaveOccurred())

        _, err = service.DeleteTrackGroup(context.Background(), trackGroup)

        Expect(err).NotTo(HaveOccurred())

        // owner := new(model.UserGroup)
        // err = db.Model(owner).Where("id = ?", trackGroupToDelete.UserGroupId).Select()
        // Expect(err).NotTo(HaveOccurred())
        // Expect(owner.TrackGroups).NotTo(ContainElement(trackGroupToDelete.Id))
        //
        // label := new(model.UserGroup)
        // err = db.Model(label).Where("id = ?", trackGroupToDelete.LabelId).Select()
        // Expect(err).NotTo(HaveOccurred())
        // Expect(label.TrackGroups).NotTo(ContainElement(trackGroupToDelete.Id))

        var tracks []model.Track
        err = db.Model(&tracks).
          Where("id in (?)", pg.In(trackGroupToDelete.Tracks)).
          Select()
        Expect(err).NotTo(HaveOccurred())
        Expect(len(tracks)).To(Equal(0))

        var trackGroups []model.Track
        err = db.Model(&trackGroups).
          Where("id in (?)", pg.In([]uuid.UUID{trackGroupToDelete.Id})).
          Select()
        Expect(err).NotTo(HaveOccurred())
        Expect(len(trackGroups)).To(Equal(0))
      })
      It("should respond with not_found error if track group does not exist", func() {
        id := uuid.NewV4()
        for id == newPlaylist.Id || id == newAlbum.Id {
          id = uuid.NewV4()
        }
        trackGroup := &pb.TrackGroup{Id: id.String()}
        resp, err := service.DeleteTrackGroup(context.Background(), trackGroup)

        Expect(resp).To(BeNil())
        Expect(err).To(HaveOccurred())

        twerr := err.(twirp.Error)
        Expect(twerr.Code()).To(Equal(not_found_code))
      })
    })
    Context("with invalid uuid", func() {
      It("should respond with invalid_argument error", func() {
        id := "45"
        trackGroup := &pb.TrackGroup{Id: id}
        resp, err := service.DeleteTrackGroup(context.Background(), trackGroup)

        Expect(resp).To(BeNil())
        Expect(err).To(HaveOccurred())

        twerr := err.(twirp.Error)
        Expect(twerr.Code()).To(Equal(invalid_argument_code))
        Expect(twerr.Meta("argument")).To(Equal("id"))
      })
    })
  })
})
