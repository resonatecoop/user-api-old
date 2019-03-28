package usergroupserver_test

import (
	"testing"
	"time"
	"path/filepath"

	// "fmt"

	// pb "user-api/rpc/user"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/go-pg/pg"
	"github.com/satori/go.uuid"

	"user-api/pkg/config"
	"user-api/pkg/postgres"

	"user-api/internal/model"
	usergroupserver "user-api/internal/server/usergroup"
)

var (
	db *pg.DB
	service *usergroupserver.Server
	newUser *model.User
	newArtist *model.UserGroup
	newUserProfile *model.UserGroup
	newRecommendedArtist *model.UserGroup
	newLabel *model.UserGroup
	newDistributor *model.UserGroup
	newArtistUserGroupMember *model.UserGroupMember
	newLabelUserGroupMember *model.UserGroupMember
	newUserGroupTaxonomy *model.GroupTaxonomy
	newArtistGroupTaxonomy *model.GroupTaxonomy
	newLabelGroupTaxonomy *model.GroupTaxonomy
	newLink *model.Link
	newGenreTag *model.Tag
	newRoleTag *model.Tag
	newAddress *model.StreetAddress
	labelAddress *model.StreetAddress
	artistAddress *model.StreetAddress
	newAlbum *model.TrackGroup
	newTrack *model.Track
	featuringTrack *model.Track
)

func TestUsergroup(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Usergroup server Suite")
}

var _ = BeforeSuite(func() {
	var err error

	cfgPath, err := filepath.Abs("./../../../conf.local.yaml")
	Expect(err).NotTo(HaveOccurred())

	cfg, err := config.Load(cfgPath)
	Expect(err).NotTo(HaveOccurred())

	db, err = pgsql.New(cfg.DB.Test.PSN, cfg.DB.Test.LogQueries, cfg.DB.Test.TimeoutSeconds)
	Expect(err).NotTo(HaveOccurred())
	service = usergroupserver.NewServer(db)

	// Create a new user (users table's empty)
	newUser = &model.User{Username: "username", FullName: "full name", Email: "email@fake.com"}
	err = db.Insert(newUser)
	Expect(err).NotTo(HaveOccurred())

	// Create group taxonomies
	newUserGroupTaxonomy = &model.GroupTaxonomy{Type: "user", Name: "User"}
	err = db.Insert(newUserGroupTaxonomy)
	Expect(err).NotTo(HaveOccurred())

	newArtistGroupTaxonomy = &model.GroupTaxonomy{Type: "artist", Name: "Artist"}
	err = db.Insert(newArtistGroupTaxonomy)
	Expect(err).NotTo(HaveOccurred())

	newLabelGroupTaxonomy = &model.GroupTaxonomy{Type: "label", Name: "Label"}
	err = db.Insert(newLabelGroupTaxonomy)
	Expect(err).NotTo(HaveOccurred())

	newDistributorGroupTaxonomy := &model.GroupTaxonomy{Type: "distributor", Name: "Distributor"}
	err = db.Insert(newDistributorGroupTaxonomy)
	Expect(err).NotTo(HaveOccurred())

	// Create link
	newLink = &model.Link{Platform: "fakebook", Uri: "https://fakebook.com/bestartist"}
	err = db.Insert(newLink)
	Expect(err).NotTo(HaveOccurred())

	// Create tags
	newGenreTag = &model.Tag{Type: "genre", Name: "rock"}
	err = db.Insert(newGenreTag)
	Expect(err).NotTo(HaveOccurred())
	newRoleTag = &model.Tag{Type: "role", Name: "bassist"}
	err = db.Insert(newRoleTag)
	Expect(err).NotTo(HaveOccurred())

	newAddress = &model.StreetAddress{Data: map[string]string{"some": "data"}}
	err = db.Insert(newAddress)
	Expect(err).NotTo(HaveOccurred())

	// Create user groups
	avatar := make([]byte, 5)
	labelAddress = &model.StreetAddress{Data: map[string]string{"some": "label data"}}
	err = db.Insert(labelAddress)
	newLabel = &model.UserGroup{
		DisplayName: "best label ever",
		Avatar: avatar,
		OwnerId: newUser.Id,
		TypeId: newLabelGroupTaxonomy.Id,
		AddressId: labelAddress.Id,
	}
	_, err = db.Model(newLabel).Returning("*").Insert()
	Expect(err).NotTo(HaveOccurred())

	newUserProfile = &model.UserGroup{
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
	artistAddress = &model.StreetAddress{Data: map[string]string{"some": "artist data"}}
	err = db.Insert(artistAddress)
	Expect(err).NotTo(HaveOccurred())
	newArtist = &model.UserGroup{
		DisplayName: "best artist ever",
		Avatar: avatar,
		OwnerId: newUser.Id,
		TypeId: newArtistGroupTaxonomy.Id,
		AddressId: artistAddress.Id,
		Links: links,
		Tags: genreTags,
		RecommendedArtists: []uuid.UUID{newUserProfile.Id},
		Publisher: map[string]string{
			"name": "publisher name",
			"number": "1E3",
		},
		Pro: map[string]string{
			"name": "PRO name",
			"number": "2BA",
		},
	}
	_, err = db.Model(newArtist).Returning("*").Insert()
	Expect(err).NotTo(HaveOccurred())

	roleTags := []uuid.UUID{newRoleTag.Id}
	newArtistUserGroupMember = &model.UserGroupMember{
		UserGroupId: newArtist.Id,
		MemberId: newUserProfile.Id,
		DisplayName: "John Doe",
		Tags: roleTags,
	}
	_, err = db.Model(newArtistUserGroupMember).Returning("*").Insert()
	Expect(err).NotTo(HaveOccurred())

	newLabelUserGroupMember = &model.UserGroupMember{
		UserGroupId: newLabel.Id,
		MemberId: newArtist.Id,
		DisplayName: newArtist.DisplayName,
	}
	_, err = db.Model(newLabelUserGroupMember).Returning("*").Insert()
	Expect(err).NotTo(HaveOccurred())

	newRecommendedArtist = &model.UserGroup{
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

	newDistributor = &model.UserGroup{
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
	newTrack = &model.Track{
		CreatorId: newUser.Id,
		UserGroupId: newArtist.Id,
		Artists: []uuid.UUID{newArtist.Id},
		Title: "track title",
		Status: "free",
		Tags: tagIds,
	}
	err = db.Insert(newTrack)
	Expect(err).NotTo(HaveOccurred())

	featuringTrack = &model.Track{
		CreatorId: newUser.Id,
		UserGroupId: newLabel.Id,
		Artists: []uuid.UUID{newArtist.Id, newRecommendedArtist.Id},
		Title: "featuring track",
		Status: "free",
		Tags: tagIds,
	}
	err = db.Insert(featuringTrack)
	Expect(err).NotTo(HaveOccurred())

	// Create a label compilation
	newCompilation := &model.TrackGroup{
		CreatorId: newUser.Id,
		UserGroupId: newLabel.Id,
		LabelId: newLabel.Id,
		Title: "compil title",
		ReleaseDate: time.Now(),
		Type: "lp",
		Cover: avatar,
		Tracks: []uuid.UUID{featuringTrack.Id},
	}
	err = db.Insert(newCompilation)
	Expect(err).NotTo(HaveOccurred())

	// Create a new album
	tracks := []uuid.UUID{newTrack.Id}
	newAlbum = &model.TrackGroup{
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

	newArtist.ArtistOfTracks = []uuid.UUID{newTrack.Id, featuringTrack.Id}
	newArtist.RecommendedBy = []uuid.UUID{newRecommendedArtist.Id}
	newArtist.HighlightedTracks = []uuid.UUID{newTrack.Id}
	newArtist.FeaturedTrackGroupId = newAlbum.Id
	_, err = db.Model(newArtist).
		Column("artist_of_tracks", "recommended_by", "highlighted_tracks", "featured_track_group_id").
		WherePK().Update()
	Expect(err).NotTo(HaveOccurred())

	newRecommendedArtist.ArtistOfTracks = []uuid.UUID{featuringTrack.Id}
	_, err = db.Model(newRecommendedArtist).Column("artist_of_tracks").WherePK().Update()
	Expect(err).NotTo(HaveOccurred())

	newTrack.TrackGroups = []uuid.UUID{newAlbum.Id}
	_, err = db.Model(newTrack).Column("track_groups").WherePK().Update()
	Expect(err).NotTo(HaveOccurred())

	featuringTrack.TrackGroups = []uuid.UUID{newCompilation.Id}
	_, err = db.Model(featuringTrack).Column("track_groups").WherePK().Update()
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	// Delete all tracks
	var tracks []model.Track
	err := db.Model(&tracks).Select()
	Expect(err).NotTo(HaveOccurred())
	if len(tracks) > 0 {
		_, err = db.Model(&tracks).Delete()
		Expect(err).NotTo(HaveOccurred())
	}

	// Delete all track groups
	var trackGroups []model.TrackGroup
	err = db.Model(&trackGroups).Select()
	Expect(err).NotTo(HaveOccurred())
	if len(trackGroups) > 0 {
		_, err = db.Model(&trackGroups).Delete()
		Expect(err).NotTo(HaveOccurred())
	}

	// Delete all userGroupMembers
	var userGroupMembers []model.UserGroupMember
	err = db.Model(&userGroupMembers).Select()
	Expect(err).NotTo(HaveOccurred())
	if len(userGroupMembers) > 0 {
		_, err = db.Model(&userGroupMembers).Delete()
		Expect(err).NotTo(HaveOccurred())
	}

	// Delete all userGroups
	var userGroups []model.UserGroup
	err = db.Model(&userGroups).Select()
	Expect(err).NotTo(HaveOccurred())
	_, err = db.Model(&userGroups).Delete()
	Expect(err).NotTo(HaveOccurred())

	// Delete all streetAddresses
	var streetAddresses []model.StreetAddress
	err = db.Model(&streetAddresses).Select()
	Expect(err).NotTo(HaveOccurred())
	_, err = db.Model(&streetAddresses).Delete()
	Expect(err).NotTo(HaveOccurred())

	// Delete all user group privacies
	var privacies []model.UserGroupPrivacy
	err = db.Model(&privacies).Select()
	Expect(err).NotTo(HaveOccurred())
	if len(privacies) > 0 {
		_, err = db.Model(&privacies).Delete()
		Expect(err).NotTo(HaveOccurred())
	}

	// Delete all users
	var users []model.User
	err = db.Model(&users).Select()
	Expect(err).NotTo(HaveOccurred())
	_, err = db.Model(&users).Delete()
	Expect(err).NotTo(HaveOccurred())

	// Delete all groupTaxonomies
	var groupTaxonomies []model.GroupTaxonomy
	err = db.Model(&groupTaxonomies).Select()
	Expect(err).NotTo(HaveOccurred())
	_, err = db.Model(&groupTaxonomies).Delete()
	Expect(err).NotTo(HaveOccurred())

	// Delete all links
	var links []model.Link
	err = db.Model(&links).Select()
	Expect(err).NotTo(HaveOccurred())
	if len(links) > 0 {
		_, err = db.Model(&links).Delete()
		Expect(err).NotTo(HaveOccurred())
	}

	// Delete all tags
	var tags []model.Tag
	err = db.Model(&tags).Select()
	Expect(err).NotTo(HaveOccurred())
	_, err = db.Model(&tags).Delete()
	Expect(err).NotTo(HaveOccurred())
})
