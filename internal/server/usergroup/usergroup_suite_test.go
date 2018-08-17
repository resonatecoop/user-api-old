package usergroupserver_test

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
	newAlbum *models.TrackGroup
	newTrack *models.Track
	featuringTrack *models.Track
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
			DisplayName: "DJ Best",
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
			RecommendedArtists: []uuid.UUID{newUserProfile.Id},
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
			DisplayName: "recommended",
			Avatar: avatar,
			OwnerId: newUser.Id,
			TypeId: newArtistGroupTaxonomy.Id,
			AddressId: newAddress.Id,
			RecommendedArtists: []uuid.UUID{newArtist.Id},
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

		// Create a new track
		tagIds := []uuid.UUID{newGenreTag.Id}
		newTrack = &models.Track{
			CreatorId: newUser.Id,
			UserGroupId: newArtist.Id,
			Artists: []uuid.UUID{newArtist.Id},
			Title: "track title",
			Status: "free",
			Tags: tagIds,
		}
		err = db.Insert(newTrack)
		Expect(err).NotTo(HaveOccurred())

		featuringTrack = &models.Track{
			CreatorId: newUser.Id,
			UserGroupId: newLabel.Id,
			Artists: []uuid.UUID{newArtist.Id, newRecommendedArtist.Id},
			Title: "featuring track",
			Status: "free",
			Tags: tagIds,
		}
		err = db.Insert(featuringTrack)
		Expect(err).NotTo(HaveOccurred())

		// Create a new album
		tracks := []uuid.UUID{newTrack.Id}
		newAlbum = &models.TrackGroup{
			CreatorId: newUser.Id,
			UserGroupId: newArtist.Id,
			LabelId: newLabel.Id,
			Title: "album title",
			ReleaseDate: time.Now(),
			Type: "lp",
			Cover: avatar,
			Tracks: tracks,
			Tags: tagIds,
		}
		err = db.Insert(newAlbum)
		Expect(err).NotTo(HaveOccurred())

		newUserProfile.RecommendedBy = []uuid.UUID{newArtist.Id}
		_, err = db.Model(newUserProfile).Column("recommended_by").WherePK().Update()
		Expect(err).NotTo(HaveOccurred())

		newArtist.Tracks = []uuid.UUID{newTrack.Id, featuringTrack.Id}
		newArtist.TrackGroups = []uuid.UUID{newAlbum.Id}
		newArtist.RecommendedBy = []uuid.UUID{newRecommendedArtist.Id}
		newArtist.HighlightedTracks = []uuid.UUID{newTrack.Id}
		newArtist.FeaturedTrackGroupId = newAlbum.Id
		_, err = db.Model(newArtist).
			Column("tracks", "track_groups", "recommended_by", "highlighted_tracks", "featured_track_group_id").
			WherePK().Update()
		Expect(err).NotTo(HaveOccurred())

		newRecommendedArtist.Tracks = []uuid.UUID{featuringTrack.Id}
		_, err = db.Model(newRecommendedArtist).Column("tracks").WherePK().Update()
		Expect(err).NotTo(HaveOccurred())

		newLabel.TrackGroups = []uuid.UUID{newAlbum.Id}
		newLabel.Tracks = []uuid.UUID{featuringTrack.Id}
		_, err = db.Model(newLabel).Column("tracks", "track_groups").WherePK().Update()
		Expect(err).NotTo(HaveOccurred())

		newTrack.TrackGroups = []uuid.UUID{newAlbum.Id}
		_, err = db.Model(newTrack).Column("track_groups").WherePK().Update()
		Expect(err).NotTo(HaveOccurred())

		// Create plays
		for i := 1; i <= 3; i++ {
			newTrackPlay := &models.Play{
				UserId: newUser.Id,
				TrackId: newTrack.Id,
				Type: "paid",
				Credits: 0.02, // constant for simplicity
			}
			err = db.Insert(newTrackPlay)
			Expect(err).NotTo(HaveOccurred())
		}
		newFreeTrackPlay := &models.Play{
			UserId: newUser.Id,
			TrackId: newTrack.Id,
			Type: "free",
			Credits: 0.00,
		}
		err = db.Insert(newFreeTrackPlay)
		Expect(err).NotTo(HaveOccurred())

		featuringTrackPlay := &models.Play{
			UserId: newUser.Id,
			TrackId: featuringTrack.Id,
			Type: "paid",
			Credits: 0.02,
		}
		err = db.Insert(featuringTrackPlay)
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

	// Delete all track groups
	var trackGroups []models.TrackGroup
	err = db.Model(&trackGroups).Select()
	Expect(err).NotTo(HaveOccurred())
	if len(trackGroups) > 0 {
		_, err = db.Model(&trackGroups).Delete()
		Expect(err).NotTo(HaveOccurred())
	}

	// Delete all userGroupMembers
	var userGroupMembers []models.UserGroupMember
	err = db.Model(&userGroupMembers).Select()
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
	_, err = db.Model(&tags).Delete()
	Expect(err).NotTo(HaveOccurred())
})
