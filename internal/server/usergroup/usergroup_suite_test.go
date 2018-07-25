package usergroupserver_test

import (
	"testing"
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

var (
	db *pg.DB
	service *usergroupserver.Server
	newUser *models.User
	newArtist *models.UserGroup
	newUserProfile *models.UserGroup
	newRecommendedArtist *models.UserGroup
	newLabel *models.UserGroup
	newDistributor *models.UserGroup
	newArtistUserGroupMember *models.UserGroupMember
	newLabelUserGroupMember *models.UserGroupMember
	newUserGroupTaxonomy *models.GroupTaxonomy
	newArtistGroupTaxonomy *models.GroupTaxonomy
	newLabelGroupTaxonomy *models.GroupTaxonomy
	newLink *models.Link
	newGenreTag *models.Tag
	newRoleTag *models.Tag
	newAddress *models.StreetAddress
	artistAddress *models.StreetAddress
)

func TestUsergroup(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Usergroup server Suite")
}

var _ = BeforeSuite(func() {
		testing := true
		db = database.Connect(testing)
		service = usergroupserver.NewServer(db)

		// Create a new user (users table's empty)
		newUser = &models.User{Username: "username", FullName: "full name", Email: "email@fake.com"}
		err := db.Insert(newUser)
		Expect(err).NotTo(HaveOccurred())

		// Create group taxonomies
		newUserGroupTaxonomy = &models.GroupTaxonomy{Type: "user", Name: "User"}
		err = db.Insert(newUserGroupTaxonomy)
		Expect(err).NotTo(HaveOccurred())

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

		// Create tags
		newGenreTag = &models.Tag{Type: "genre", Name: "rock"}
		err = db.Insert(newGenreTag)
		Expect(err).NotTo(HaveOccurred())
		newRoleTag = &models.Tag{Type: "role", Name: "bassist"}
		err = db.Insert(newRoleTag)
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

		newUserProfile = &models.UserGroup{
			DisplayName: "DJ JohnD",
			Avatar: avatar,
			OwnerId: newUser.Id,
			TypeId: newUserGroupTaxonomy.Id,
			AddressId: newAddress.Id,
		}
		_, err = db.Model(newUserProfile).Returning("*").Insert()
		Expect(err).NotTo(HaveOccurred())

		links := []uuid.UUID{newLink.Id}
		genreTags := []uuid.UUID{newGenreTag.Id}
		artistAddress = &models.StreetAddress{Data: map[string]string{"some": "artist data"}}
		err = db.Insert(artistAddress)
		Expect(err).NotTo(HaveOccurred())
		newArtist = &models.UserGroup{
			DisplayName: "best artist ever",
			Avatar: avatar,
			OwnerId: newUser.Id,
			TypeId: newArtistGroupTaxonomy.Id,
			AddressId: artistAddress.Id,
			Links: links,
			Tags: genreTags,
		}
		_, err = db.Model(newArtist).Returning("*").Insert()
		Expect(err).NotTo(HaveOccurred())

		roleTags := []uuid.UUID{newRoleTag.Id}
		newArtistUserGroupMember = &models.UserGroupMember{
			UserGroupId: newArtist.Id,
			MemberId: newUserProfile.Id,
			DisplayName: "John Doe",
			Tags: roleTags,
		}
		_, err = db.Model(newArtistUserGroupMember).Returning("*").Insert()
		Expect(err).NotTo(HaveOccurred())

		newLabelUserGroupMember = &models.UserGroupMember{
			UserGroupId: newLabel.Id,
			MemberId: newArtist.Id,
			DisplayName: newArtist.DisplayName,
		}
		_, err = db.Model(newLabelUserGroupMember).Returning("*").Insert()
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
			Column("Privacy", "Address").
			WherePK().
			Select()
		Expect(err).NotTo(HaveOccurred())

		newDistributor = &models.UserGroup{
			DisplayName: "distributor",
			Avatar: avatar,
			OwnerId: newUser.Id,
			TypeId: newArtistGroupTaxonomy.Id,
			AddressId: newAddress.Id,
		}
		_, err = db.Model(newDistributor).Returning("*").Insert()
		Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	// Delete all userGroupMembers
	var userGroupMembers []models.UserGroupMember
	err := db.Model(&userGroupMembers).Select()
	Expect(err).NotTo(HaveOccurred())
	_, err = db.Model(&userGroupMembers).Delete()
	Expect(err).NotTo(HaveOccurred())

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

	// Delete all user group privacies
	var privacies []models.UserGroupPrivacy
	err = db.Model(&privacies).Select()
	Expect(err).NotTo(HaveOccurred())
	_, err = db.Model(&privacies).Delete()
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

	// Delete all links
	// var links []models.Link
	// err = db.Model(&links).Select()
	// Expect(err).NotTo(HaveOccurred())
	// _, err = db.Model(&links).Delete()
	// Expect(err).NotTo(HaveOccurred())

	// Delete all tags
	var tags []models.Tag
	err = db.Model(&tags).Select()
	Expect(err).NotTo(HaveOccurred())
	_, err = db.Model(&tags).Delete()
	Expect(err).NotTo(HaveOccurred())
})
