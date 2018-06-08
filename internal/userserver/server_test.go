package userserver_test

import (
	// "fmt"
	// "reflect"
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/twitchtv/twirp"

	pb "user-api/rpc/user"
	// . "user-api/internal/userserver"
)

var _ = Describe("Userserver", func() {
	Describe("CreateUser", func() {
		const already_exists_code twirp.ErrorCode = "already_exists"
		const invalid_argument_code twirp.ErrorCode = "invalid_argument"

		Context("with all required attributes", func() {
			It("should create a new user", func() {
				user := &pb.User{Username: "janed", FullName: "jane d", DisplayName: "jad", Email: "jane@d.com"}
				resp, err := service.CreateUser(context.Background(), user)

				Expect(err).NotTo(HaveOccurred())
				Expect(resp.Username).To(Equal("janed"))
				Expect(resp.FullName).To(Equal("jane d"))
				Expect(resp.DisplayName).To(Equal("jad"))
				Expect(resp.Email).To(Equal("jane@d.com"))
				Expect(resp.Id).NotTo(Equal(""))
			})

			It("should not create a user with same display_name", func() {
				user := &pb.User{Username: "janedoe", FullName: "jane doe", DisplayName: "jad", Email: "jane@doe.com"}
				resp, err := service.CreateUser(context.Background(), user)
				twerr := err.(twirp.Error)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				Expect(twerr.Code()).To(Equal(already_exists_code))
				Expect(twerr.Msg()).To(Equal("display_name"))
			})

			It("should not create a user with same email", func() {
				user := &pb.User{Username: "janedoe", FullName: "jane doe", DisplayName: "jadoe", Email: "jane@d.com"}
				resp, err := service.CreateUser(context.Background(), user)
				twerr := err.(twirp.Error)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				Expect(twerr.Code()).To(Equal(already_exists_code))
				Expect(twerr.Msg()).To(Equal("email"))
			})

			It("should not create a user with same username", func() {
				user := &pb.User{Username: "janed", FullName: "jane doe", DisplayName: "jadoe", Email: "jane@doe.com"}
				resp, err := service.CreateUser(context.Background(), user)
				twerr := err.(twirp.Error)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				Expect(twerr.Code()).To(Equal(already_exists_code))
				Expect(twerr.Msg()).To(Equal("username"))
			})
		})

		Context("with missing required attributes", func() {
			It("should not create a user without display_name", func() {
				user := &pb.User{Username: "johnd", FullName: "john doe", DisplayName: "", Email: "john@doe.com"}
				resp, err := service.CreateUser(context.Background(), user)
				twerr := err.(twirp.Error)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				Expect(twerr.Code()).To(Equal(invalid_argument_code))
				Expect(twerr.Meta("argument")).To(Equal("display_name"))

			})

			It("should not create a user without email", func() {
				user := &pb.User{Username: "johnd", FullName: "john doe", DisplayName: "john", Email: ""}
				resp, err := service.CreateUser(context.Background(), user)
				twerr := err.(twirp.Error)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				Expect(twerr.Code()).To(Equal(invalid_argument_code))
				Expect(twerr.Meta("argument")).To(Equal("email"))
			})

			It("should not create a user without username", func() {
				user := &pb.User{Username: "", FullName: "john doe", DisplayName: "john", Email: "john@doe.com"}
				resp, err := service.CreateUser(context.Background(), user)
				twerr := err.(twirp.Error)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				Expect(twerr.Code()).To(Equal(invalid_argument_code))
				Expect(twerr.Meta("argument")).To(Equal("username"))
			})

			It("should not create a user without full_name", func() {
				user := &pb.User{Username: "johnd", FullName: "", DisplayName: "john", Email: "john@doe.com"}
				resp, err := service.CreateUser(context.Background(), user)
				twerr := err.(twirp.Error)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				Expect(twerr.Code()).To(Equal(invalid_argument_code))
				Expect(twerr.Meta("argument")).To(Equal("full_name"))
			})
		})
	})
})
