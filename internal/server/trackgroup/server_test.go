package trackgroupserver_test

import (
  "context"

  // "github.com/go-pg/pg"
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "github.com/twitchtv/twirp"
  // "github.com/satori/go.uuid"
  "github.com/golang/protobuf/ptypes"

  pb "user-api/rpc/trackgroup"
  // trackpb "user-api/rpc/track"
  // "user-api/internal/database/models"
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
        Expect(res.DisplayArtist).To(Equal(newAlbum.DisplayArtist))
        Expect(res.MultipleComposers).To(Equal(newAlbum.MultipleComposers))
        Expect(res.Private).To(Equal(newAlbum.Private))
        Expect(res.LabelId).To(Equal(newAlbum.LabelId.String()))
        Expect(res.UserGroupId).To(Equal(newAlbum.UserGroupId.String()))

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
        Expect(res.LabelId).To(Equal(newPlaylist.LabelId.String()))
        Expect(res.UserGroupId).To(Equal(newPlaylist.UserGroupId.String()))

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
    })
  })

  Describe("UpdateTrackGroup", func() {

  })

  Describe("AddTracksToTrackGroup", func() {

  })

  Describe("DeleteTracksFromTrackGroup", func() {

  })

  Describe("CreateTrackGroup", func() {

  })

  Describe("DeleteTrackGroup", func() {

  })
})
