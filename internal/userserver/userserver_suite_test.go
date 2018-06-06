package userserver_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestUserserver(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Userserver Suite")
}
