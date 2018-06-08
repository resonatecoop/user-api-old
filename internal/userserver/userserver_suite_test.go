package userserver_test

import (
	"testing"

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

func TestUserserver(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Userserver Suite")
}

var _ = BeforeSuite(func() {
		testing := true
		db = database.Connect(testing)
		service = userserver.NewServer(db)
})

var _ = AfterSuite(func() {
	var users []models.User
  err := db.Model(&users).Select()
	Expect(err).NotTo(HaveOccurred())
	_, err = db.Model(&users).Delete()
	Expect(err).NotTo(HaveOccurred())
})
