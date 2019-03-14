package address_server_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http/httptest"
	"net/http"
	"encoding/json"

	addressServer "user-api/internal/server/address"
)

var (
	service *addressServer.Server
	okResults addressServer.Results
	algolia *httptest.Server
)

const (
	okResponse = `{
		"hits": [
			{
				"country": {"de":"Frankreich","ru":"Франция","pt":"França","it":"Francia","hu":"Franciaország","es":"Francia","zh":"法国","ar":"فرنسا","default":"France","ja":"フランス","pl":"Francja","ro":"Franța","nl":"Frankrijk"},
				"is_country":false,
				"city":{"ar":["باريس"],"default":["Paris"],"ru":["Париж"],"ja":["パリ"],"it":["Parigi"],"pl":["Paryż"],"hu":["Párizs"],"es":["París"],"zh":["巴黎"],"nl":["Parijs"]},
				"is_highway":false,
				"importance":15,
				"_tags":["capital","boundary/administrative","city","place/city","country/fr","source/pristine"],
				"postcode":["75000"],
				"county":{"default":["Paris"]},
				"population":2220445,
				"country_code":"fr",
				"is_city":true,
				"is_popular":true,
				"administrative":["Île-de-France"],
				"admin_level":2,
				"district":"Paris",
				"locale_names":{"default":["Paris"]},
				"_geoloc":{"lat":48.8546,"lng":2.34771},
				"objectID":"af9599b85ad97dcead64bbb7bc500445"
			}
		],
		"nbHits":1,
		"processingTimeMS":16,
		"query":"Paris",
		"params":"query=Paris&hitsPerPage=1",
		"degradedQuery":false
	}`
)

func TestAddress(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Address Suite")
}

var _ = BeforeSuite(func() {
	algolia = AlgoliaResponseStub()
	service = addressServer.NewServer(algolia.URL, "", "")
	json.Unmarshal([]byte(okResponse), &okResults)
})

var _ = AfterSuite(func() {
	algolia.Close()
})

func AlgoliaResponseStub() *httptest.Server {
  return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // switch r.RequestURI {
    //   case "/latest/meta-data/instance-id":
    //     resp = "i-12345"
    //   case "/latest/meta-data/placement/availability-zone":
    //     resp = "us-west-2a"
    //   default:
    //     http.Error(w, "not found", http.StatusNotFound)
    //     return
    // }
		w.WriteHeader(http.StatusOK)
    w.Write([]byte(okResponse))
  }))
}
