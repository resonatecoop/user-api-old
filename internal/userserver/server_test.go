package userserver_test

import (
	// "fmt"
	// "reflect"
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/twitchtv/twirp"
	"github.com/satori/go.uuid"

	pb "user-api/rpc/user"
	"user-api/internal/database/models"
	// . "user-api/internal/userserver"
)

var _ = Describe("Userserver", func() {
	const already_exists_code twirp.ErrorCode = "already_exists"
	const invalid_argument_code twirp.ErrorCode = "invalid_argument"
	const not_found_code twirp.ErrorCode = "not_found"

	Describe("GetUser", func() {
		Context("with valid uuid", func() {
			It("should respond with user if it exists", func() {
				// First create a new user (users table's empty)
				newuser := &models.User{Username: "username", FullName: "full name", DisplayName: "display name", Email: "email@fake.com"}
				err := db.Insert(newuser)
				Expect(err).NotTo(HaveOccurred())

				// Then get it
				user := &pb.User{Id: newuser.Id.String()}
				resp, err := service.GetUser(context.Background(), user)

				Expect(err).NotTo(HaveOccurred())
				Expect(resp.Id).To(Equal(newuser.Id.String()))
				Expect(resp.Username).To(Equal(newuser.Username))
				Expect(resp.FullName).To(Equal(newuser.FullName))
				Expect(resp.DisplayName).To(Equal(newuser.DisplayName))
				Expect(resp.Email).To(Equal(newuser.Email))

				err = db.Delete(newuser)
				Expect(err).NotTo(HaveOccurred())
			})
			It("should respond with not_found error if user does not exist", func() {
				id := uuid.NewV4()
				user := &pb.User{Id: id.String()}
				resp, err := service.GetUser(context.Background(), user)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(not_found_code))
			})
		})
		Context("with invalid uuid", func() {
			It("should respond with invalid_argument error", func() {
				id := "45"
				user := &pb.User{Id: id}
				resp, err := service.GetUser(context.Background(), user)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
				Expect(twerr.Meta("argument")).To(Equal("id"))
			})
		})
	})

	Describe("CreateUser", func() {
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

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(already_exists_code))
				Expect(twerr.Msg()).To(Equal("display_name"))
			})

			It("should not create a user with same email", func() {
				user := &pb.User{Username: "janedoe", FullName: "jane doe", DisplayName: "jadoe", Email: "jane@d.com"}
				resp, err := service.CreateUser(context.Background(), user)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(already_exists_code))
				Expect(twerr.Msg()).To(Equal("email"))
			})

			It("should not create a user with same username", func() {
				user := &pb.User{Username: "janed", FullName: "jane doe", DisplayName: "jadoe", Email: "jane@doe.com"}
				resp, err := service.CreateUser(context.Background(), user)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(already_exists_code))
				Expect(twerr.Msg()).To(Equal("username"))
			})
		})

		Context("with missing required attributes", func() {
			It("should not create a user without display_name", func() {
				user := &pb.User{Username: "johnd", FullName: "john doe", DisplayName: "", Email: "john@doe.com"}
				resp, err := service.CreateUser(context.Background(), user)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
				Expect(twerr.Meta("argument")).To(Equal("display_name"))

			})

			It("should not create a user without email", func() {
				user := &pb.User{Username: "johnd", FullName: "john doe", DisplayName: "john", Email: ""}
				resp, err := service.CreateUser(context.Background(), user)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
				Expect(twerr.Meta("argument")).To(Equal("email"))
			})

			It("should not create a user without username", func() {
				user := &pb.User{Username: "", FullName: "john doe", DisplayName: "john", Email: "john@doe.com"}
				resp, err := service.CreateUser(context.Background(), user)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
				Expect(twerr.Meta("argument")).To(Equal("username"))
			})

			It("should not create a user without full_name", func() {
				user := &pb.User{Username: "johnd", FullName: "", DisplayName: "john", Email: "john@doe.com"}
				resp, err := service.CreateUser(context.Background(), user)

				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())

				twerr := err.(twirp.Error)
				Expect(twerr.Code()).To(Equal(invalid_argument_code))
				Expect(twerr.Meta("argument")).To(Equal("full_name"))
			})
		})
	})
})
