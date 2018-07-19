package userserver_test

import (
	"testing"
	"time"

	// pb "user-api/rpc/user"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/go-pg/pg"
	"github.com/satori/go.uuid"

	"user-api/internal/database"
	"user-api/internal/database/models"
	userserver "user-api/internal/server/user"
)

var db *pg.DB
var service *userserver.Server
var ownerOfFollowedUserGroup *models.User
var newUser *models.User
var newTrack *models.Track
var newUserGroup *models.UserGroup
var newFavoriteTrack *models.Track
var newFollowedUserGroup *models.UserGroup

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

		// Create a new tracks
		duration, _ := time.ParseDuration("10m10s")
		cover := make([]byte, 5)
		newTrack = &models.Track{PublishDate: time.Now(), Title: "track title", Duration: duration, Status: "free", Cover: cover}
		err = db.Insert(newTrack)
		Expect(err).NotTo(HaveOccurred())

		// Create a new users
		ownerOfFollowedUserGroup = &models.User{Username: "ownerOfFollowedUserGroup", FullName: "Owner", DisplayName: "ownerOfFollowedUserGroup", Email: "owner@fake.com"}
		err = db.Insert(ownerOfFollowedUserGroup)
		Expect(err).NotTo(HaveOccurred())

		newGroupTaxonomy := &models.GroupTaxonomy{Type: "artist", Name: "Artist"}
		err = db.Insert(newGroupTaxonomy)
		Expect(err).NotTo(HaveOccurred())
		avatar := make([]byte, 5)

		newUser = &models.User{Username: "username", FullName: "full name", DisplayName: "display name", Email: "email@fake.com"}
		err = db.Insert(newUser)
		Expect(err).NotTo(HaveOccurred())

		newFavoriteTrack = &models.Track{PublishDate: time.Now(), Title: "fav track title", Duration: duration, Status: "free", Cover: cover, FavoriteOfUsers: []uuid.UUID{newUser.Id}}
		err = db.Insert(newFavoriteTrack)
		Expect(err).NotTo(HaveOccurred())

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

		newUser.FollowedGroups = []uuid.UUID{newFollowedUserGroup.Id}
		newUser.FavoriteTracks = []uuid.UUID{newFavoriteTrack.Id}
		_, err = db.Model(newUser).Column("followed_groups", "favorite_tracks").WherePK().Update()
		Expect(err).NotTo(HaveOccurred())

		// Create a new user_group
		newUserGroup = &models.UserGroup{
			DisplayName: "artist",
			Avatar: avatar,
			OwnerId: newUser.Id,
			TypeId: newGroupTaxonomy.Id,
			AddressId: newAddress.Id,
		}
		err = db.Insert(newUserGroup)
		Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	// Delete all users
	var users []models.User
  err := db.Model(&users).Select()
	Expect(err).NotTo(HaveOccurred())
	_, err = db.Model(&users).Delete()
	Expect(err).NotTo(HaveOccurred())

	// Delete all tracks
	var tracks []models.Track
	err = db.Model(&tracks).Select()
	Expect(err).NotTo(HaveOccurred())
	_, err = db.Model(&tracks).Delete()
	Expect(err).NotTo(HaveOccurred())

	// Delete all userGroups
	var userGroups []models.UserGroup
	err = db.Model(&userGroups).Select()
	Expect(err).NotTo(HaveOccurred())
	_, err = db.Model(&userGroups).Delete()
	Expect(err).NotTo(HaveOccurred())

	// Delete all groupTaxonomies
	var groupTaxonomies []models.GroupTaxonomy
	err = db.Model(&groupTaxonomies).Select()
	Expect(err).NotTo(HaveOccurred())
	_, err = db.Model(&groupTaxonomies).Delete()
	Expect(err).NotTo(HaveOccurred())
})
