package address_server_test

import (
  "context"
  // "fmt"
  // "reflect"
  // "time"
  // "github.com/go-pg/pg"
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "github.com/twitchtv/twirp"

  pb "user-api/rpc/address"
  // "user-api/internal/database/model"
)


var _ = Describe("Address server", func() {
  const invalid_argument_code twirp.ErrorCode = "invalid_argument"

  Describe("SearchAddress", func() {
    Context("with valid query", func() {
      It("should respond with address search results", func() {
        q := &pb.AddressQuery{Query: "Paris"}
        resp, err := service.SearchAddress(context.Background(), q)

        Expect(err).NotTo(HaveOccurred())
        Expect(resp).NotTo(BeNil())
        Expect(resp.Hits).NotTo(BeNil())
        Expect(len(resp.Hits)).To(Equal(1))
        Expect(resp.Hits[0].ObjectId).To(Equal(okResults.Hits[0].ObjectId))
        Expect(resp.Hits[0].Country).To(Equal(okResults.Hits[0].Country))
        Expect(resp.Hits[0].Postcode).To(Equal(okResults.Hits[0].Postcode))
        Expect(resp.Hits[0].Administrative).To(Equal(okResults.Hits[0].Administrative))
        Expect(resp.Hits[0].CountryCode).To(Equal(okResults.Hits[0].CountryCode))
        Expect(resp.Hits[0].Geoloc).To(Equal(okResults.Hits[0].Geoloc))
        Expect(resp.NbHits).To(Equal(okResults.NbHits))
      })
    })
    Context("with invalid query", func() {
      It("should respond with invalid error", func() {
        q := &pb.AddressQuery{Query: "Pa"}
        resp, err := service.SearchAddress(context.Background(), q)

        Expect(resp).To(BeNil())
        Expect(err).To(HaveOccurred())

        twerr := err.(twirp.Error)
        Expect(twerr.Code()).To(Equal(invalid_argument_code))
        Expect(twerr.Meta("argument")).To(Equal("query"))
      })
    })
  })
})
