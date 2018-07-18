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
var newUser *models.User
var newTrack *models.Track
var newUserGroup *models.UserGroup

func TestUser(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "User server Suite")
}

var _ = BeforeSuite(func() {
		testing := true
		db = database.Connect(testing)
		service = userserver.NewServer(db)

		// Create a new user (users table's empty)
		newUser = &models.User{Username: "username", FullName: "full name", DisplayName: "display name", Email: "email@fake.com"}
		err := db.Insert(newUser)
		Expect(err).NotTo(HaveOccurred())

		// Create a new track
		duration, _ := time.ParseDuration("10m10s")
		cover := make([]byte, 5)
		newTrack = &models.Track{PublishDate: time.Now(), Title: "track title", Duration: duration, Status: "free", Cover: cover}
		err = db.Insert(newTrack)
		Expect(err).NotTo(HaveOccurred())

		// Create a new user_group
		newGroupTaxonomy := &models.GroupTaxonomy{Type: "artist", Name: "Artist"}
		err = db.Insert(newGroupTaxonomy)
		Expect(err).NotTo(HaveOccurred())
		avatar := make([]byte, 5)
		admins := []uuid.UUID{newUser.Id}
		newUserGroup = &models.UserGroup{
			DisplayName: "artist",
			Avatar: avatar,
			OwnerId: newUser.Id,
			TypeId: newGroupTaxonomy.Id,
			AdminUsers: admins,
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
