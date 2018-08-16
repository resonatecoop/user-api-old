package userserver_test

import (
	"testing"
	"time"
	// "fmt"

	// pb "user-api/rpc/user"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/go-pg/pg"
	"github.com/satori/go.uuid"

	"user-api/internal/database"
	"user-api/internal/database/models"
	userserver "user-api/internal/server/user"
)

var (
	db *pg.DB
	service *userserver.Server
	ownerOfFollowedUserGroup *models.User
	newUser *models.User
	newTrack *models.Track
	newUserGroup *models.UserGroup
	newFavoriteTrack *models.Track
	newFollowedUserGroup *models.UserGroup
	newAlbum *models.TrackGroup
	newUserPlaylist *models.TrackGroup
	newTrackPlay *models.Play
	newFavoriteTrackPlay *models.Play
)

func TestUser(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "User server Suite")
}

var _ = BeforeSuite(func() {
		testing := true
		db = database.Connect(testing)
		service = userserver.NewServer(db)

		newAddress := &models.StreetAddress{Data: map[string]string{"some": "data"}}
		err := db.Insert(newAddress)
		Expect(err).NotTo(HaveOccurred())

		// Create a new user
		ownerOfFollowedUserGroup = &models.User{Username: "ownerOfFollowedUserGroup", FullName: "Owner", Email: "owner@fake.com"}
		err = db.Insert(ownerOfFollowedUserGroup)
		Expect(err).NotTo(HaveOccurred())

		newGroupTaxonomy := &models.GroupTaxonomy{Type: "artist", Name: "Artist"}
		err = db.Insert(newGroupTaxonomy)
		Expect(err).NotTo(HaveOccurred())

		newLabelGroupTaxonomy := &models.GroupTaxonomy{Type: "label", Name: "Label"}
		err = db.Insert(newLabelGroupTaxonomy)
		Expect(err).NotTo(HaveOccurred())

		avatar := make([]byte, 5)

		newUser = &models.User{Username: "username", FullName: "full name", Email: "email@fake.com"}
		err = db.Insert(newUser)
		Expect(err).NotTo(HaveOccurred())

		/*newLabel := &models.UserGroup{
			DisplayName: "label",
			Avatar: avatar,
			OwnerId: newUser.Id,
			TypeId: newLabelGroupTaxonomy.Id,
			AddressId: newAddress.Id,
		}
		err = db.Insert(newLabel)
		Expect(err).NotTo(HaveOccurred())

		newLabelTrackGroup := &models.TrackGroup{
			CreatorId: newUser.Id,
			UserGroupId: newLabel.Id,
			Title: "label compilation",
			ReleaseDate: time.Now(),
			Type: "lp",
			Cover: avatar,
			Private: false,
		}
		err = db.Insert(newLabelTrackGroup)
		Expect(err).NotTo(HaveOccurred())

		newLabelTracks := []uuid.UUID{}
		for i := 1; i <= 6; i++ {
			newLabelTrack := &models.Track{
				CreatorId: newUser.Id,
				UserGroupId: newLabel.Id,
				TrackGroups: []uuid.UUID{newLabelTrackGroup.Id},
				Title: "compilation track title",
				Status: "paid",
			}
			err = db.Insert(newLabelTrack)
			Expect(err).NotTo(HaveOccurred())
			newLabelTracks = append(newLabelTracks, newLabelTrack.Id)
			for i := 1; i <= 9; i++ {
				newLabelTrackPlay := &models.Play{
					UserId: newUser.Id,
					TrackId: newLabelTrack.Id,
					Type: "paid",
					Credits: 0.02, // constant for simplicity
				}
				err = db.Insert(newLabelTrackPlay)
				Expect(err).NotTo(HaveOccurred())
			}
		}
		newLabelTrackGroup.Tracks = newLabelTracks
		_, err = db.Model(newLabelTrackGroup).Column("tracks").WherePK().Update()
		Expect(err).NotTo(HaveOccurred())*/

		followers := []uuid.UUID{newUser.Id}
		newFollowedUserGroup = &models.UserGroup{
			DisplayName: "followed group",
			Avatar: avatar,
			OwnerId: ownerOfFollowedUserGroup.Id,
			TypeId: newGroupTaxonomy.Id,
			Followers: followers,
			AddressId: newAddress.Id,
		}
		err = db.Insert(newFollowedUserGroup)
		Expect(err).NotTo(HaveOccurred())

		newFavoriteTrack = &models.Track{
			CreatorId: ownerOfFollowedUserGroup.Id,
			UserGroupId: newFollowedUserGroup.Id,
			Artists: []uuid.UUID{newFollowedUserGroup.Id},
			Title: "fav track title",
			Status: "paid",
			FavoriteOfUsers: []uuid.UUID{newUser.Id},
		}
		err = db.Insert(newFavoriteTrack)
		Expect(err).NotTo(HaveOccurred())

		newTrack = &models.Track{
			CreatorId: ownerOfFollowedUserGroup.Id,
			UserGroupId: newFollowedUserGroup.Id,
			Title: "track title",
			Status: "paid",
		}
		err = db.Insert(newTrack)
		Expect(err).NotTo(HaveOccurred())

		albumTracks := []uuid.UUID{newFavoriteTrack.Id, newTrack.Id}
		newAlbum = &models.TrackGroup{
			CreatorId: ownerOfFollowedUserGroup.Id,
			UserGroupId: newFollowedUserGroup.Id,
			Title: "album title",
			ReleaseDate: time.Now(),
			Type: "lp",
			Cover: avatar,
			Tracks: albumTracks,
		}
		err = db.Insert(newAlbum)
		Expect(err).NotTo(HaveOccurred())

		newTrackPlay = &models.Play{
			UserId: newUser.Id,
			TrackId: newTrack.Id,
			Type: "free",
			Credits: 0.00,
		}
		err = db.Insert(newTrackPlay)
		Expect(err).NotTo(HaveOccurred())

		tracks := []uuid.UUID{newTrack.Id}
		newUserPlaylist = &models.TrackGroup{
			CreatorId: newUser.Id,
			Title: "user playlist",
			ReleaseDate: time.Now(),
			Type: "playlist",
			Cover: avatar,
			Tracks: tracks,
			Private: true,
		}
		err = db.Insert(newUserPlaylist)
		Expect(err).NotTo(HaveOccurred())

		newUser.Playlists = []uuid.UUID{newUserPlaylist.Id}
		newUser.FollowedGroups = []uuid.UUID{newFollowedUserGroup.Id}
		newUser.FavoriteTracks = []uuid.UUID{newFavoriteTrack.Id}
		_, err = db.Model(newUser).Column("playlists", "followed_groups", "favorite_tracks").WherePK().Update()
		Expect(err).NotTo(HaveOccurred())

		userGroupAddress := &models.StreetAddress{Data: map[string]string{"some": "data"}}
		err = db.Insert(userGroupAddress)
		Expect(err).NotTo(HaveOccurred())
		// Create a new user_group
		newUserGroup = &models.UserGroup{
			DisplayName: "artist",
			Avatar: avatar,
			OwnerId: newUser.Id,
			TypeId: newGroupTaxonomy.Id,
			AddressId: userGroupAddress.Id,
			Tracks: []uuid.UUID{newTrack.Id},
		}
		err = db.Insert(newUserGroup)
		Expect(err).NotTo(HaveOccurred())

		newUserGroupTrackGroup := &models.TrackGroup{
			CreatorId: newUser.Id,
			UserGroupId: newUserGroup.Id,
			Title: "user group album",
			ReleaseDate: time.Now(),
			Type: "lp",
			Cover: avatar,
			Private: false,
		}
		err = db.Insert(newUserGroupTrackGroup)
		Expect(err).NotTo(HaveOccurred())

		newUserGroupTracks := []uuid.UUID{newTrack.Id}
		for i := 1; i <= 9; i++ {
			newUserGroupTrack := &models.Track{
				CreatorId: newUser.Id,
				UserGroupId: newUserGroup.Id,
				Artists: []uuid.UUID{newUserGroup.Id},
				TrackGroups: []uuid.UUID{newUserGroupTrackGroup.Id},
				Title: "album track title",
				Status: "paid",
			}
			err = db.Insert(newUserGroupTrack)
			Expect(err).NotTo(HaveOccurred())
			newUserGroupTracks = append(newUserGroupTracks, newUserGroupTrack.Id)
			for i := 1; i <= 9; i++ {
				newUserGroupTrackPlay := &models.Play{
					UserId: newUser.Id,
					TrackId: newUserGroupTrack.Id,
					Type: "paid",
					Credits: 0.02, // constant for simplicity
				}
				err = db.Insert(newUserGroupTrackPlay)
				Expect(err).NotTo(HaveOccurred())
			}
		}

		for i := 1; i <= 9; i++ {
			newTrackPaidPlay := &models.Play{
				UserId: newUser.Id,
				TrackId: newTrack.Id,
				Type: "paid",
				Credits: 0.02, // constant for simplicity
			}
			err = db.Insert(newTrackPaidPlay)
			Expect(err).NotTo(HaveOccurred())
		}

		newFavoriteTrackPlay = &models.Play{
			UserId: newUser.Id,
			TrackId: newFavoriteTrack.Id,
			Type: "paid",
			Credits: 0.02,
		}
		err = db.Insert(newFavoriteTrackPlay)
		Expect(err).NotTo(HaveOccurred())

		newFavoriteTrack.TrackGroups = []uuid.UUID{newAlbum.Id}
		_, err = db.Model(newFavoriteTrack).Column("track_groups").WherePK().Update()
		Expect(err).NotTo(HaveOccurred())

		newUserGroup.Tracks = newUserGroupTracks
		_, err = db.Model(newUserGroup).Column("tracks").WherePK().Update()
		Expect(err).NotTo(HaveOccurred())

		newUserGroupTrackGroup.Tracks = newUserGroupTracks
		_, err = db.Model(newUserGroupTrackGroup).Column("tracks").WherePK().Update()
		Expect(err).NotTo(HaveOccurred())

		newTrack.TrackGroups = []uuid.UUID{newAlbum.Id, newUserPlaylist.Id}
		newTrack.Artists = []uuid.UUID{newFollowedUserGroup.Id, newUserGroup.Id}
		_, err = db.Model(newTrack).Column("artists", "track_groups").WherePK().Update()
		Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	// Delete all tracks
	var tracks []models.Track
	err := db.Model(&tracks).Select()
	Expect(err).NotTo(HaveOccurred())
	if len(tracks) > 0 {
		_, err = db.Model(&tracks).Delete()
		Expect(err).NotTo(HaveOccurred())
	}

	var trackGroups []models.TrackGroup
	err = db.Model(&trackGroups).Select()
	Expect(err).NotTo(HaveOccurred())
	if len(trackGroups) > 0 {
		_, err = db.Model(&trackGroups).Delete()
		Expect(err).NotTo(HaveOccurred())
	}

	// Delete all userGroups
	var userGroups []models.UserGroup
	err = db.Model(&userGroups).Select()
	Expect(err).NotTo(HaveOccurred())
	_, err = db.Model(&userGroups).Delete()
	Expect(err).NotTo(HaveOccurred())

	// Delete all streetAddresses
	var streetAddresses []models.StreetAddress
	err = db.Model(&streetAddresses).Select()
	Expect(err).NotTo(HaveOccurred())
	_, err = db.Model(&streetAddresses).Delete()
	Expect(err).NotTo(HaveOccurred())

	// Delete all users
	var users []models.User
	err = db.Model(&users).Select()
	Expect(err).NotTo(HaveOccurred())
	_, err = db.Model(&users).Delete()
	Expect(err).NotTo(HaveOccurred())

	// Delete all groupTaxonomies
	var groupTaxonomies []models.GroupTaxonomy
	err = db.Model(&groupTaxonomies).Select()
	Expect(err).NotTo(HaveOccurred())
	_, err = db.Model(&groupTaxonomies).Delete()
	Expect(err).NotTo(HaveOccurred())

	// Delete all plays
	var plays []models.Play
	err = db.Model(&plays).Select()
	Expect(err).NotTo(HaveOccurred())
	_, err = db.Model(&plays).Delete()
	Expect(err).NotTo(HaveOccurred())

	// Delete all tags
	var tags []models.Tag
	err = db.Model(&tags).Select()
	Expect(err).NotTo(HaveOccurred())
	if len(tags) > 0 {
		_, err = db.Model(&tags).Delete()
		Expect(err).NotTo(HaveOccurred())
	}
})
