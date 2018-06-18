package userserver_test

import (
	"testing"
	"time"

	// pb "user-api/rpc/user"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/go-pg/pg"

	"user-api/internal/database"
	"user-api/internal/database/models"
	"user-api/internal/userserver"
)

var db *pg.DB
var service *userserver.Server
var newuser *models.User
var newtrack *models.Track

func TestUserserver(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Userserver Suite")
}

var _ = BeforeSuite(func() {
		testing := true
		db = database.Connect(testing)
		service = userserver.NewServer(db)

		// Create a new user (users table's empty)
		newuser = &models.User{Username: "username", FullName: "full name", DisplayName: "display name", Email: "email@fake.com"}
		err := db.Insert(newuser)
		Expect(err).NotTo(HaveOccurred())

		// Create a new track
		duration, _ := time.ParseDuration("10m10s")
		cover := make([]byte, 5)
		newtrack = &models.Track{PublishDate: time.Now(), Title: "track title", Duration: duration, Status: "free", Cover: cover}
		err = db.Insert(newtrack)
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
})
