package usergroupserver_test

import (
	"testing"
	// "time"

	// pb "user-api/rpc/user"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/go-pg/pg"
	// "github.com/satori/go.uuid"

	"user-api/internal/database"
	"user-api/internal/database/models"
	usergroupserver "user-api/internal/server/usergroup"
)

var db *pg.DB
var service *usergroupserver.Server
var newUser *models.User
var newArtist *models.UserGroup
var newLabel *models.UserGroup
var newArtistGroupTaxonomy *models.GroupTaxonomy
var newLabelGroupTaxonomy *models.GroupTaxonomy

func TestUsergroup(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Usergroup Suite")
}

var _ = BeforeSuite(func() {
		testing := true
		db = database.Connect(testing)
		service = usergroupserver.NewServer(db)

		// Create a new user (users table's empty)
		newUser = &models.User{Username: "username", FullName: "full name", DisplayName: "display name", Email: "email@fake.com"}
		err := db.Insert(newUser)
		Expect(err).NotTo(HaveOccurred())

		// Create group taxonomies
		newArtistGroupTaxonomy = &models.GroupTaxonomy{Type: "artist", Name: "Artist"}
		err = db.Insert(newArtistGroupTaxonomy)
		Expect(err).NotTo(HaveOccurred())

		newLabelGroupTaxonomy = &models.GroupTaxonomy{Type: "label", Name: "Label"}
		err = db.Insert(newLabelGroupTaxonomy)
		Expect(err).NotTo(HaveOccurred())


		// Create user groups
		avatar := make([]byte, 5)
		// admins := []uuid.UUID{newUser.Id}
		newArtist = &models.UserGroup{
			DisplayName: "best artist ever",
			Avatar: avatar,
			OwnerId: newUser.Id,
			TypeId: newArtistGroupTaxonomy.Id,
			// AdminUsers: admins,
		}
		err = db.Insert(newArtist)
		Expect(err).NotTo(HaveOccurred())

		newLabel = &models.UserGroup{
			DisplayName: "best label ever",
			Avatar: avatar,
			OwnerId: newUser.Id,
			TypeId: newLabelGroupTaxonomy.Id,
			// AdminUsers: admins,
		}
		err = db.Insert(newLabel)
})

var _ = AfterSuite(func() {
	// Delete all users
	var users []models.User
  err := db.Model(&users).Select()
	Expect(err).NotTo(HaveOccurred())
	_, err = db.Model(&users).Delete()
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
