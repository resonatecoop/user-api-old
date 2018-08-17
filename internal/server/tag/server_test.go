package tagserver_test

import (
  "context"
  // "time"
  // "github.com/go-pg/pg"
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "github.com/twitchtv/twirp"
  // "github.com/satori/go.uuid"

  pb "user-api/rpc/tag"
  // "user-api/internal/database/models"
)

var _ = Describe("Tag server", func() {
	const already_exists_code twirp.ErrorCode = "already_exists"
	const invalid_argument_code twirp.ErrorCode = "invalid_argument"
	const not_found_code twirp.ErrorCode = "not_found"

  Describe("SearchGenres", func() {
    Context("with valid query", func() {
      It("should respond with tracks and track groups (playlists and albums)", func() {
        q := &pb.Query{Query: "pop"}
        res, err := service.SearchGenres(context.Background(), q)

        Expect(err).NotTo(HaveOccurred())
        Expect(res).NotTo(BeNil())

        Expect(len(res.Playlists)).To(Equal(1))
        Expect(res.Playlists[0].Id).To(Equal(newPlaylist.Id.String()))
        Expect(res.Playlists[0].Title).To(Equal(newPlaylist.Title))
        Expect(res.Playlists[0].TotalTracks).To(Equal(int32(len(newPlaylist.Tracks))))
        Expect(res.Playlists[0].UserGroup).NotTo(BeNil())
        Expect(res.Playlists[0].UserGroup.Id).To(Equal(newArtistUserGroup.Id.String()))
        Expect(res.Playlists[0].UserGroup.DisplayName).To(Equal(newArtistUserGroup.DisplayName))
        Expect(res.Playlists[0].UserGroup.Avatar).To(Equal(newArtistUserGroup.Avatar))

        Expect(len(res.Albums)).To(Equal(1))
        Expect(res.Albums[0].Id).To(Equal(newAlbum.Id.String()))
        Expect(res.Albums[0].Title).To(Equal(newAlbum.Title))
        Expect(res.Albums[0].TotalTracks).To(Equal(int32(len(newAlbum.Tracks))))
        Expect(res.Albums[0].UserGroup).NotTo(BeNil())
        Expect(res.Albums[0].UserGroup.Id).To(Equal(newArtistUserGroup.Id.String()))
        Expect(res.Albums[0].UserGroup.DisplayName).To(Equal(newArtistUserGroup.DisplayName))
        Expect(res.Albums[0].UserGroup.Avatar).To(Equal(newArtistUserGroup.Avatar))
      })
    })
    Context("with invalid query", func() {
      It("should respond with invalid error", func() {
        q := &pb.Query{Query: "po"}
        resp, err := service.SearchGenres(context.Background(), q)

        Expect(resp).To(BeNil())
        Expect(err).To(HaveOccurred())

        twerr := err.(twirp.Error)
        Expect(twerr.Code()).To(Equal(invalid_argument_code))
        Expect(twerr.Meta("argument")).To(Equal("query"))
      })
    })
  })
})
