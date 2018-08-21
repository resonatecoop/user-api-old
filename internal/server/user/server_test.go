package userserver_test

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
				Expect(resp.FullName).To(Equal(newUser.FullName))
				Expect(resp.Email).To(Equal(newUser.Email))
				Expect(len(resp.OwnerOfGroups)).To(Equal(1))
				Expect(resp.OwnerOfGroups[0].Id).To(Equal(newUserGroup.Id.String()))
				Expect(resp.OwnerOfGroups[0].DisplayName).To(Equal(newUserGroup.DisplayName))
				Expect(resp.OwnerOfGroups[0].Avatar).To(Equal(newUserGroup.Avatar))
			})
			It("should respond with not_found error if user does not exist", func() {
				id := uuid.NewV4()
				for id == newUser.Id || id == ownerOfFollowedUserGroup.Id {
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

	Describe("CreatePlay", func() {
		Context("with valid track id and user id", func() {
			It("should create a play and respond with updated play count and credits", func() {
				playRequest := &pb.CreatePlayRequest{
					Play: &pb.Play{
						UserId: newUser.Id.String(),
						TrackId: newFavoriteTrack.Id.String(),
						Type: "paid",
						Credits: 0.04,
					},
					UpdatedCredits: 1.00,
				}
				resp, err := service.CreatePlay(context.Background(), playRequest)

				Expect(err).NotTo(HaveOccurred())
				Expect(resp.UpdatedCredits).To(Equal(playRequest.UpdatedCredits))
				Expect(resp.UpdatedPlayCount).To(Equal(int32(2)))
			})
			It("should respond with not_found error if user does not exist", func() {
				id := uuid.NewV4()
				for id == newUser.Id || id == ownerOfFollowedUserGroup.Id {
					id = uuid.NewV4()
				}

				playRequest := &pb.CreatePlayRequest{
					Play: &pb.Play{
						UserId: id.String(),
						TrackId: newFavoriteTrack.Id.String(),
						Type: "paid",
						Credits: 0.04,
					},
					UpdatedCredits: 1.00,
				}
				resp, err := service.CreatePlay(context.Background(), playRequest)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(not_found_code))
			})
			It("should respond with not_found error if track does not exist", func() {
				id := uuid.NewV4()
				for id == newTrack.Id || id == newFavoriteTrack.Id {
					id = uuid.NewV4()
				}

				playRequest := &pb.CreatePlayRequest{
					Play: &pb.Play{
						UserId: newUser.Id.String(),
						TrackId: id.String(),
						Type: "paid",
						Credits: 0.04,
					},
					UpdatedCredits: 1.00,
				}
				resp, err := service.CreatePlay(context.Background(), playRequest)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(not_found_code))
			})
		})
		Context("with invalid request", func() {
			Context("with invalid uuid", func() {
				It("should respond with invalid_argument error if track_id invalid", func() {
					id := "45"
					playRequest := &pb.CreatePlayRequest{
						Play: &pb.Play{
							UserId: newUser.Id.String(),
							TrackId: id,
							Type: "paid",
							Credits: 0.04,
						},
						UpdatedCredits: 1.00,
					}
					resp, err := service.CreatePlay(context.Background(), playRequest)

					Expect(resp).To(BeNil())
					Expect(err).To(HaveOccurred())

					twerr := err.(twirp.Error)
					Expect(twerr.Code()).To(Equal(invalid_argument_code))
					Expect(twerr.Meta("argument")).To(Equal("id"))
				})
				It("should respond with invalid_argument error if user_id invalid", func() {
					id := "45"
					playRequest := &pb.CreatePlayRequest{
						Play: &pb.Play{
							UserId: id,
							TrackId: newTrack.Id.String(),
							Type: "paid",
							Credits: 0.04,
						},
						UpdatedCredits: 1.00,
					}
					resp, err := service.CreatePlay(context.Background(), playRequest)

					Expect(resp).To(BeNil())
					Expect(err).To(HaveOccurred())

					twerr := err.(twirp.Error)
					Expect(twerr.Code()).To(Equal(invalid_argument_code))
					Expect(twerr.Meta("argument")).To(Equal("id"))
				})
				It("should respond with invalid_argument error if type invalid", func() {
					playRequest := &pb.CreatePlayRequest{
						Play: &pb.Play{
							UserId: newUser.Id.String(),
							TrackId: newTrack.Id.String(),
							Type: "",
							Credits: 0.04,
						},
						UpdatedCredits: 1.00,
					}
					resp, err := service.CreatePlay(context.Background(), playRequest)

					Expect(resp).To(BeNil())
					Expect(err).To(HaveOccurred())
					twerr := err.(twirp.Error)
					Expect(twerr.Code()).To(Equal(invalid_argument_code))
					Expect(twerr.Meta("argument")).To(Equal("type"))
				})
			})
		})
	})

	Describe("GetPlaylists", func() {
		Context("with valid uuid", func() {
			It("should respond with playlists", func() {
				user := &pb.User{Id: newUser.Id.String()}
				resp, err := service.GetPlaylists(context.Background(), user)

				Expect(err).NotTo(HaveOccurred())
				Expect(len(resp.Playlists)).To(Equal(1))
				Expect(resp.Playlists[0].Id).To(Equal(newUserPlaylist.Id.String()))
				Expect(resp.Playlists[0].Title).To(Equal(newUserPlaylist.Title))
				Expect(resp.Playlists[0].Cover).To(Equal(newUserPlaylist.Cover))
				Expect(resp.Playlists[0].Type).To(Equal(newUserPlaylist.Type))
				Expect(resp.Playlists[0].About).To(Equal(newUserPlaylist.About))
				Expect(resp.Playlists[0].Private).To(Equal(newUserPlaylist.Private))
				Expect(resp.Playlists[0].DisplayArtist).To(Equal(newUserPlaylist.DisplayArtist))
				Expect(resp.Playlists[0].TotalTracks).To(Equal(int32(1)))
			})
			It("should respond with not_found error if user does not exist", func() {
				id := uuid.NewV4()
				for id == newUser.Id || id == ownerOfFollowedUserGroup.Id {
					id = uuid.NewV4()
				}
				user := &pb.User{Id: id.String()}
				resp, err := service.GetPlaylists(context.Background(), user)

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
				resp, err := service.GetPlaylists(context.Background(), user)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
				Expect(twerr.Meta("argument")).To(Equal("id"))
			})
		})
	})


	Describe("GetSupportedArtists", func() {
		Context("with valid uuid", func() {
			It("should respond with owned tracks", func() {
				user := &pb.User{Id: newUser.Id.String()}
				res, err := service.GetSupportedArtists(context.Background(), user)

				Expect(err).NotTo(HaveOccurred())
				Expect(len(res.Artists)).To(Equal(1))
				Expect(res.Artists[0].Id).To(Equal(newUserGroup.Id.String()))
				Expect(res.Artists[0].DisplayName).To(Equal(newUserGroup.DisplayName))
				Expect(res.Artists[0].Avatar).To(Equal(newUserGroup.Avatar))
			})
			It("should respond with not_found error if user does not exist", func() {
				id := uuid.NewV4()
				for id == newUser.Id || id == ownerOfFollowedUserGroup.Id {
					id = uuid.NewV4()
				}
				user := &pb.User{Id: id.String()}
				resp, err := service.GetSupportedArtists(context.Background(), user)

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
				resp, err := service.GetSupportedArtists(context.Background(), user)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
				Expect(twerr.Meta("argument")).To(Equal("id"))
			})
		})
	})

	Describe("GetFavoriteTracks", func() {
		Context("with valid uuid", func() {
			It("should respond with favorite tracks", func() {
				user := &pb.User{Id: newUser.Id.String()}
				u := url.URL{}
				queryString := u.Query()
				queryString.Set("page", "1")
    		queryString.Set("limit", "50")
				ctx := context.WithValue(context.Background(), "query", queryString)
				resp, err := service.GetFavoriteTracks(ctx, user)

				Expect(err).NotTo(HaveOccurred())
				Expect(len(resp.Tracks)).To(Equal(1))
				Expect(resp.Tracks[0].Id).To(Equal(newFavoriteTrack.Id.String()))
				Expect(resp.Tracks[0].Title).To(Equal(newFavoriteTrack.Title))
				Expect(resp.Tracks[0].TrackServerId).To(Equal(newFavoriteTrack.TrackServerId.String()))
				Expect(resp.Tracks[0].Duration).To(Equal(newFavoriteTrack.Duration))
				Expect(resp.Tracks[0].Status).To(Equal(newFavoriteTrack.Status))
				Expect(resp.Tracks[0].TrackNumber).To(Equal(newFavoriteTrack.TrackNumber))

				Expect(len(resp.Tracks[0].TrackGroups)).To(Equal(1))
				Expect(resp.Tracks[0].TrackGroups[0].Id).To(Equal(newAlbum.Id.String()))
				Expect(resp.Tracks[0].TrackGroups[0].Title).To(Equal(newAlbum.Title))
				Expect(resp.Tracks[0].TrackGroups[0].Cover).To(Equal(newAlbum.Cover))
				Expect(resp.Tracks[0].TrackGroups[0].Type).To(Equal(newAlbum.Type))
				Expect(resp.Tracks[0].TrackGroups[0].About).To(Equal(newAlbum.About))
				Expect(resp.Tracks[0].TrackGroups[0].Private).To(Equal(newAlbum.Private))
				Expect(resp.Tracks[0].TrackGroups[0].DisplayArtist).To(Equal(newAlbum.DisplayArtist))
				Expect(resp.Tracks[0].TrackGroups[0].TotalTracks).To(Equal(int32(2)))

				Expect(len(resp.Tracks[0].Artists)).To(Equal(1))
				Expect(resp.Tracks[0].Artists[0].Id).To(Equal(newFollowedUserGroup.Id.String()))
				Expect(resp.Tracks[0].Artists[0].DisplayName).To(Equal(newFollowedUserGroup.DisplayName))
				Expect(resp.Tracks[0].Artists[0].Avatar).To(Equal(newFollowedUserGroup.Avatar))
			})
			It("should respond with not_found error if user does not exist", func() {
				id := uuid.NewV4()
				for id == newUser.Id || id == ownerOfFollowedUserGroup.Id {
					id = uuid.NewV4()
				}
				user := &pb.User{Id: id.String()}
				resp, err := service.GetFavoriteTracks(context.Background(), user)

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
				resp, err := service.GetFavoriteTracks(context.Background(), user)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
				Expect(twerr.Meta("argument")).To(Equal("id"))
			})
		})
	})

	Describe("GetOwnedTracks", func() {
		Context("with valid uuid", func() {
			It("should respond with owned tracks", func() {
				user := &pb.User{Id: newUser.Id.String()}
				u := url.URL{}
				queryString := u.Query()
				queryString.Set("page", "1")
    		queryString.Set("limit", "50")
				ctx := context.WithValue(context.Background(), "query", queryString)
				resp, err := service.GetOwnedTracks(ctx, user)

				Expect(err).NotTo(HaveOccurred())
				Expect(len(resp.Tracks)).To(Equal(10))
				Expect(resp.Tracks[0].Id).To(Equal(newTrack.Id.String()))
				Expect(resp.Tracks[0].Title).To(Equal(newTrack.Title))
				Expect(resp.Tracks[0].TrackServerId).To(Equal(newTrack.TrackServerId.String()))
				Expect(resp.Tracks[0].Duration).To(Equal(newTrack.Duration))
				Expect(resp.Tracks[0].Status).To(Equal(newTrack.Status))
				Expect(resp.Tracks[0].TrackNumber).To(Equal(newTrack.TrackNumber))

				Expect(len(resp.Tracks[0].TrackGroups)).To(Equal(1))
				Expect(resp.Tracks[0].TrackGroups[0].Id).To(Equal(newAlbum.Id.String()))
				Expect(resp.Tracks[0].TrackGroups[0].Title).To(Equal(newAlbum.Title))
				Expect(resp.Tracks[0].TrackGroups[0].Cover).To(Equal(newAlbum.Cover))
				Expect(resp.Tracks[0].TrackGroups[0].Type).To(Equal(newAlbum.Type))
				Expect(resp.Tracks[0].TrackGroups[0].About).To(Equal(newAlbum.About))
				Expect(resp.Tracks[0].TrackGroups[0].Private).To(Equal(newAlbum.Private))
				Expect(resp.Tracks[0].TrackGroups[0].TotalTracks).To(Equal(int32(2)))

				Expect(len(resp.Tracks[0].Artists)).To(Equal(2))
				Expect(resp.Tracks[0].Artists[0].Id).To(Equal(newFollowedUserGroup.Id.String()))
				Expect(resp.Tracks[0].Artists[0].DisplayName).To(Equal(newFollowedUserGroup.DisplayName))
				Expect(resp.Tracks[0].Artists[0].Avatar).To(Equal(newFollowedUserGroup.Avatar))
				Expect(resp.Tracks[0].Artists[1].Id).To(Equal(newUserGroup.Id.String()))
				Expect(resp.Tracks[0].Artists[1].DisplayName).To(Equal(newUserGroup.DisplayName))
				Expect(resp.Tracks[0].Artists[1].Avatar).To(Equal(newUserGroup.Avatar))
			})
			It("should respond with not_found error if user does not exist", func() {
				id := uuid.NewV4()
				for id == newUser.Id || id == ownerOfFollowedUserGroup.Id {
					id = uuid.NewV4()
				}
				user := &pb.User{Id: id.String()}
				resp, err := service.GetOwnedTracks(context.Background(), user)

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
				resp, err := service.GetOwnedTracks(context.Background(), user)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
				Expect(twerr.Meta("argument")).To(Equal("id"))
			})
		})
	})

	Describe("GetTrackHistory", func() {
		Context("with valid uuid", func() {
			It("should respond with track history", func() {
				user := &pb.User{Id: newUser.Id.String()}
				u := url.URL{}
				queryString := u.Query()
				queryString.Set("page", "1")
    		queryString.Set("limit", "50")
				ctx := context.WithValue(context.Background(), "query", queryString)
				resp, err := service.GetTrackHistory(ctx, user)

				Expect(err).NotTo(HaveOccurred())
				Expect(len(resp.Tracks)).To(Equal(11))

				Expect(resp.Tracks[0].Id).To(Equal(newFavoriteTrack.Id.String()))
				Expect(resp.Tracks[0].Title).To(Equal(newFavoriteTrack.Title))
				Expect(resp.Tracks[0].TrackServerId).To(Equal(newFavoriteTrack.TrackServerId.String()))
				Expect(resp.Tracks[0].Duration).To(Equal(newFavoriteTrack.Duration))
				Expect(resp.Tracks[0].Status).To(Equal(newFavoriteTrack.Status))
				Expect(resp.Tracks[0].TrackNumber).To(Equal(newFavoriteTrack.TrackNumber))

				Expect(resp.Tracks[1].Id).To(Equal(newTrack.Id.String()))
				Expect(resp.Tracks[1].Title).To(Equal(newTrack.Title))
				Expect(resp.Tracks[1].TrackServerId).To(Equal(newTrack.TrackServerId.String()))
				Expect(resp.Tracks[1].Duration).To(Equal(newTrack.Duration))
				Expect(resp.Tracks[1].Status).To(Equal(newTrack.Status))
				Expect(resp.Tracks[1].TrackNumber).To(Equal(newTrack.TrackNumber))
			})
			It("should respond with not_found error if user does not exist", func() {
				id := uuid.NewV4()
				for id == newUser.Id || id == ownerOfFollowedUserGroup.Id {
					id = uuid.NewV4()
				}
				user := &pb.User{Id: id.String()}
				resp, err := service.GetTrackHistory(context.Background(), user)

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
				resp, err := service.GetTrackHistory(context.Background(), user)

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
				Expect(user.FavoriteTracks).To(HaveLen(2))
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
				Expect(user.FavoriteTracks).To(HaveLen(1))
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
				Expect(user.FollowedGroups).To(HaveLen(2))
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
				Expect(user.FollowedGroups).To(HaveLen(1))
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
					Email: "email@fake.com",
					FirstName: "first name",
				}
				_, err := service.UpdateUser(context.Background(), user)

				Expect(err).NotTo(HaveOccurred())
			})
			It("should respond with not_found error if user does not exist", func() {
				id := uuid.NewV4()
				for id == newUser.Id || id == ownerOfFollowedUserGroup.Id {
					id = uuid.NewV4()
				}
				user := &pb.User{
					Id: id.String(),
					Username: "username",
					FullName: "fullname",
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
				Expect(resp.Email).To(Equal("jane@d.com"))
				Expect(resp.Id).NotTo(Equal(""))
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

				userToDelete := new(models.User)
				err := db.Model(userToDelete).Column("user.*", "OwnerOfGroups").Where("id = ?", newUser.Id).Select()
				Expect(err).NotTo(HaveOccurred())

				_, err = service.DeleteUser(context.Background(), user)

				Expect(err).NotTo(HaveOccurred())

				var followedGroups []models.UserGroup
				err = db.Model(&followedGroups).
					Where("id in (?)", pg.In(userToDelete.FollowedGroups)).
					Select()
				Expect(err).NotTo(HaveOccurred())
				for _, group := range followedGroups {
					Expect(group.Followers).NotTo(ContainElement(userToDelete.Id))
				}

				var favoriteTracks []models.Track
				err = db.Model(&favoriteTracks).
					Where("id in (?)", pg.In(userToDelete.FavoriteTracks)).
					Select()
				Expect(err).NotTo(HaveOccurred())
				for _, track := range favoriteTracks {
					Expect(track.FavoriteOfUsers).NotTo(ContainElement(userToDelete.Id))
				}

				ownerOfGroupIds := make([]uuid.UUID, len(userToDelete.OwnerOfGroups))
				for i, group := range userToDelete.OwnerOfGroups {
					ownerOfGroupIds[i] = group.Id
				}
				var ownerOfGroups []models.UserGroup
				err = db.Model(&ownerOfGroups).
					Where("id in (?)", pg.In(ownerOfGroupIds)).
					Select()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(ownerOfGroups)).To(Equal(0))

				var playlists []models.TrackGroup
				err = db.Model(&playlists).
					Where("id in (?)", pg.In(userToDelete.Playlists)).
					Select()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(playlists)).To(Equal(0))

				var users []models.User
				err = db.Model(&users).
					Where("id in (?)", pg.In([]uuid.UUID{userToDelete.Id})).
					Select()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(users)).To(Equal(0))
			})
			It("should respond with not_found error if user does not exist", func() {
				id := uuid.NewV4()
				for id == newUser.Id || id == ownerOfFollowedUserGroup.Id {
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
