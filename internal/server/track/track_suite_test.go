package trackserver_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/go-pg/pg"
	"github.com/satori/go.uuid"

	"user-api/internal/database/models"
	"user-api/internal/database"
	trackserver "user-api/internal/server/track"
)

var (
	db *pg.DB
	service *trackserver.Server
	newUser *models.User
	newTrack *models.Track
	newArtistGroupTaxonomy *models.GroupTaxonomy
	newLabelGroupTaxonomy *models.GroupTaxonomy
	newArtistUserGroup *models.UserGroup
	newLabelUserGroup *models.UserGroup
)

func TestTrack(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Track server Suite")
}

var _ = BeforeSuite(func() {
	testing := true
	db = database.Connect(testing)
	service = trackserver.NewServer(db)

	newAddress := &models.StreetAddress{Data: map[string]string{"some": "data"}}
	err := db.Insert(newAddress)
	Expect(err).NotTo(HaveOccurred())

	newArtistGroupTaxonomy = &models.GroupTaxonomy{Type: "artist", Name: "Artist"}
	err = db.Insert(newArtistGroupTaxonomy)
	Expect(err).NotTo(HaveOccurred())

	newLabelGroupTaxonomy = &models.GroupTaxonomy{Type: "label", Name: "Label"}
	err = db.Insert(newLabelGroupTaxonomy)
	Expect(err).NotTo(HaveOccurred())

	newUser = &models.User{Username: "username", FullName: "full name", Email: "email@fake.com"}
	err = db.Insert(newUser)
	Expect(err).NotTo(HaveOccurred())

	// Create a new user_group
	avatar := make([]byte, 5)
	newArtistUserGroup = &models.UserGroup{
		DisplayName: "artist",
		Avatar: avatar,
		OwnerId: newUser.Id,
		TypeId: newArtistGroupTaxonomy.Id,
		AddressId: newAddress.Id,
	}
	err = db.Insert(newArtistUserGroup)
	Expect(err).NotTo(HaveOccurred())


	newLabelUserGroup = &models.UserGroup{
		DisplayName: "label",
		Avatar: avatar,
		OwnerId: newUser.Id,
		TypeId: newLabelGroupTaxonomy.Id,
		AddressId: newAddress.Id,
	}
	err = db.Insert(newLabelUserGroup)
	Expect(err).NotTo(HaveOccurred())

	// Create a new track
	newTrack = &models.Track{CreatorId: newUser.Id, UserGroupId: newArtistUserGroup.Id, Title: "track title", Status: "free"} // TODO add other attr
	err = db.Insert(newTrack)
	Expect(err).NotTo(HaveOccurred())

	newArtistUserGroup.Tracks = []uuid.UUID{newTrack.Id}
	_, err = db.Model(newArtistUserGroup).Column("tracks").WherePK().Update()
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	// Delete all tracks
	var tracks []models.Track
	err := db.Model(&tracks).Select()
	Expect(err).NotTo(HaveOccurred())
	_, err = db.Model(&tracks).Delete()
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

	// Delete all tags
	var tags []models.Tag
	err = db.Model(&tags).Select()
	Expect(err).NotTo(HaveOccurred())
	if len(tags) > 0 {
		_, err = db.Model(&tags).Delete()
		Expect(err).NotTo(HaveOccurred())
	}
})
