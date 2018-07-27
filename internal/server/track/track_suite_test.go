package trackserver_test

import (
	"testing"
	"time"

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
	newAlbum *models.TrackGroup
	newPlaylist *models.TrackGroup
	newArtistGroupTaxonomy *models.GroupTaxonomy
	newLabelGroupTaxonomy *models.GroupTaxonomy
	newArtistUserGroup *models.UserGroup
	newLabelUserGroup *models.UserGroup
	newGenreTag *models.Tag
)

func TestTrack(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Track server Suite")
}

var _ = BeforeSuite(func() {
	testing := true
	db = database.Connect(testing)
	service = trackserver.NewServer(db)

	newGenreTag = &models.Tag{Type: "genre", Name: "pop"}
	err := db.Insert(newGenreTag)
	Expect(err).NotTo(HaveOccurred())

	newAddress := &models.StreetAddress{Data: map[string]string{"some": "data"}}
	err = db.Insert(newAddress)
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
	tagIds := []uuid.UUID{newGenreTag.Id}
	newTrack = &models.Track{
		CreatorId: newUser.Id,
		UserGroupId: newArtistUserGroup.Id,
		Artists: []uuid.UUID{newArtistUserGroup.Id},
		Title: "track title",
		Status: "free",
		Tags: tagIds,
	}
	err = db.Insert(newTrack)
	Expect(err).NotTo(HaveOccurred())

	favoritingUser := &models.User{Username: "fav", FullName: "fav name", Email: "fav@fake.com", FavoriteTracks: []uuid.UUID{newTrack.Id}}
	err = db.Insert(favoritingUser)
	Expect(err).NotTo(HaveOccurred())

	// Create track groups
	// tracks := map[string]string{
	// 	"1": newTrack.Id.String(),
	// }
	tracks := []uuid.UUID{newTrack.Id}
	newAlbum = &models.TrackGroup{
		CreatorId: newUser.Id,
		UserGroupId: newArtistUserGroup.Id,
		LabelId: newLabelUserGroup.Id,
		Title: "album title",
		ReleaseDate: time.Now(),
		Type: "lp",
		Cover: avatar,
		Tracks: tracks,
	}
	err = db.Insert(newAlbum)
	Expect(err).NotTo(HaveOccurred())

	newPlaylist = &models.TrackGroup{
		CreatorId: newUser.Id,
		// UserGroupId: uuid.UUID{},
		// LabelId: uuid.UUID{},
		Title: "playlist title",
		ReleaseDate: time.Now(),
		Type: "playlist",
		Cover: avatar,
		Tracks: tracks,
	}
	err = db.Insert(newPlaylist)
	Expect(err).NotTo(HaveOccurred())

	newTrack.TrackGroups = []uuid.UUID{newAlbum.Id, newPlaylist.Id}
	newTrack.FavoriteOfUsers = []uuid.UUID{favoritingUser.Id}
	_, err = db.Model(newTrack).Column("track_groups").WherePK().Update()
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

	// Delete all track groups
	var trackGroups []models.TrackGroup
	err = db.Model(&trackGroups).Select()
	Expect(err).NotTo(HaveOccurred())
	_, err = db.Model(&trackGroups).Delete()
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
