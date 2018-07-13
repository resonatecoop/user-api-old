package usergroupserver_test

import (
	"testing"
	// "time"
	// "fmt"

	// pb "user-api/rpc/user"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/go-pg/pg"
	"github.com/satori/go.uuid"

	"user-api/internal/database"
	"user-api/internal/database/models"
	usergroupserver "user-api/internal/server/usergroup"
)

var db *pg.DB
var service *usergroupserver.Server
var newUser *models.User
var newArtist *models.UserGroup
var newRecommendedArtist *models.UserGroup
var newLabel *models.UserGroup
var newArtistGroupTaxonomy *models.GroupTaxonomy
var newLabelGroupTaxonomy *models.GroupTaxonomy
var newLink *models.Link
var newTag *models.Tag
var newAddress *models.StreetAddress

func TestUsergroup(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Usergroup Suite")
}

var _ = BeforeSuite(func() {
		testing := true
		db = database.Connect(testing)
		service = usergroupserver.NewServer(db)

		// Create a new user (users table's empty)
		newUser = &models.User{Username: "username", FullName: "full name", DisplayName: "displayname", Email: "email@fake.com"}
		err := db.Insert(newUser)
		Expect(err).NotTo(HaveOccurred())

		// Create group taxonomies
		newArtistGroupTaxonomy = &models.GroupTaxonomy{Type: "artist", Name: "Artist"}
		err = db.Insert(newArtistGroupTaxonomy)
		Expect(err).NotTo(HaveOccurred())

		newLabelGroupTaxonomy = &models.GroupTaxonomy{Type: "label", Name: "Label"}
		err = db.Insert(newLabelGroupTaxonomy)
		Expect(err).NotTo(HaveOccurred())

		newDistributorGroupTaxonomy := &models.GroupTaxonomy{Type: "distributor", Name: "Distributor"}
		err = db.Insert(newDistributorGroupTaxonomy)
		Expect(err).NotTo(HaveOccurred())

		// Create link
		newLink = &models.Link{Platform: "fakebook", Uri: "https://fakebook.com/bestartist"}
		err = db.Insert(newLink)
		Expect(err).NotTo(HaveOccurred())

		// Create tag
		newTag = &models.Tag{Type: "genre", Name: "rock"}
		err = db.Insert(newTag)
		Expect(err).NotTo(HaveOccurred())

		newAddress = &models.StreetAddress{Data: map[string]string{"some": "data"}}
		err = db.Insert(newAddress)
		Expect(err).NotTo(HaveOccurred())

		// Create user groups
		avatar := make([]byte, 5)
		newLabel = &models.UserGroup{
			DisplayName: "best label ever",
			Avatar: avatar,
			OwnerId: newUser.Id,
			TypeId: newLabelGroupTaxonomy.Id,
			AddressId: newAddress.Id,
		}
		_, err = db.Model(newLabel).Returning("*").Insert()
		Expect(err).NotTo(HaveOccurred())

		labels := []uuid.UUID{newLabel.Id}
		links := []uuid.UUID{newLink.Id}
		tags := []uuid.UUID{newTag.Id}
		newArtist = &models.UserGroup{
			DisplayName: "best artist ever",
			Avatar: avatar,
			OwnerId: newUser.Id,
			TypeId: newArtistGroupTaxonomy.Id,
			AddressId: newAddress.Id,
			Labels: labels,
			Links: links,
			Tags: tags,
		}
		_, err = db.Model(newArtist).Returning("*").Insert()
		Expect(err).NotTo(HaveOccurred())

		newRecommendedArtist = &models.UserGroup{
			DisplayName: "recommended by best artist ever",
			Avatar: avatar,
			OwnerId: newUser.Id,
			TypeId: newArtistGroupTaxonomy.Id,
			AddressId: newAddress.Id,
		}
		_, err = db.Model(newRecommendedArtist).Returning("*").Insert()
		Expect(err).NotTo(HaveOccurred())

		err = db.Model(newArtist).
			Column("Privacy").
			WherePK().
			Select()
		Expect(err).NotTo(HaveOccurred())
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
	// Expect(err).NotTo(HaveOccurred())
	_, err = db.Model(&userGroups).Delete()
	Expect(err).NotTo(HaveOccurred())

	// Delete all groupTaxonomies
	var groupTaxonomies []models.GroupTaxonomy
	err = db.Model(&groupTaxonomies).Select()
	Expect(err).NotTo(HaveOccurred())
	_, err = db.Model(&groupTaxonomies).Delete()
	Expect(err).NotTo(HaveOccurred())

	// Delete all streetAddresses
	var streetAddresses []models.StreetAddress
	err = db.Model(&streetAddresses).Select()
	Expect(err).NotTo(HaveOccurred())
	_, err = db.Model(&streetAddresses).Delete()
	Expect(err).NotTo(HaveOccurred())

	// Delete all links
	var links []models.Link
	err = db.Model(&links).Select()
	Expect(err).NotTo(HaveOccurred())
	_, err = db.Model(&links).Delete()
	Expect(err).NotTo(HaveOccurred())

	// Delete all tags
	var tags []models.Tag
	err = db.Model(&tags).Select()
	Expect(err).NotTo(HaveOccurred())
	_, err = db.Model(&tags).Delete()
	Expect(err).NotTo(HaveOccurred())

	// Delete all user group privacies
	var privacies []models.UserGroupPrivacy
	err = db.Model(&privacies).Select()
	Expect(err).NotTo(HaveOccurred())
	_, err = db.Model(&privacies).Delete()
	Expect(err).NotTo(HaveOccurred())
})
