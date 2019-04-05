package address_server

import (
  "context"
  "fmt"
  "encoding/json"
  "bytes"
  "net/http"
  "time"
  "github.com/twitchtv/twirp"

  pb "user-api/rpc/address"
  mapspkg "user-api/internal/pkg/maps"
)

type Server struct {
	// db *pg.DB
  url, appid, apikey string
}

func NewServer(url, appid, apikey string) *Server {
	return &Server{url: url, appid: appid, apikey: apikey}
}

type Error struct {
  Message string
}

type Result struct {
  ObjectId string `json:"objectID"`
  Country map[string]string
  Postcode []string
  Administrative []string
  CountryCode string
  Geoloc map[string]float32 `json:"_geoloc"`
  LocaleNames map[string][]string `json:"locale_names"`
  City map[string][]string
}

type Results struct {
  Hits []*Result
  NbHits int32 `json:"nbHits"`
}

func (s *Server) SearchAddress(ctx context.Context, q *pb.AddressQuery) (*pb.AddressResults, error) {
  if len(q.Query) < 3 {
    return nil, twirp.InvalidArgumentError("query", "must be a valid search query")
  }
  var hitsPerPage string
  if q.HitsPerPage > 0 {
    hitsPerPage = fmt.Sprintf("%d", q.HitsPerPage)
  }

  // TODO Implement retry strategy https://community.algolia.com/places/rest.html#rest-api
  url := s.url
  reqBody := fmt.Sprintf(`{"query":"%s", "type": "%s", "hitsPerPage": "%s"}`, q.Query, q.Type, hitsPerPage)
  jsonStr := []byte(reqBody)
  req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
  // req.Header.Set("X-Algolia-Application-Id", s.appid)
  // req.Header.Set("X-Algolia-API-Key", s.apikey)

  client := &http.Client{Timeout: time.Second * 10}
  resp, err := client.Do(req)
  if err != nil {
    return nil, err
  }

  defer resp.Body.Close()

  if resp.StatusCode != 200 {
    code := twirp.ErrorCode(fmt.Sprintf("%d", resp.StatusCode))
    errorMessage := Error{}
    if err := json.NewDecoder(resp.Body).Decode(&errorMessage); err != nil {
      return nil, err
    }
    fmt.Printf("%+v\n", resp.StatusCode)
    fmt.Printf("%+v\n", code)
    fmt.Printf("%+v\n", errorMessage)
    fmt.Printf("%+v\n", errorMessage.Message)
    fmt.Printf("%+v\n", twirp.NewError(code, errorMessage.Message))

    return nil, twirp.NewError(code, errorMessage.Message)
  }

  results := Results{}
  if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
    return nil, err
  }

  hits := make([]*pb.AddressResult, len(results.Hits))
  for i, hit := range results.Hits {
    hits[i] = &pb.AddressResult{
      ObjectId: hit.ObjectId,
      Country: hit.Country,
      Postcode: hit.Postcode,
      Administrative: hit.Administrative,
      CountryCode: hit.CountryCode,
      Geoloc: hit.Geoloc,
      City: mapspkg.GetMapListValue(hit.City),
      LocaleNames: mapspkg.GetMapListValue(hit.LocaleNames),
    }
  }


  return &pb.AddressResults{
    Hits: hits,
    NbHits: results.NbHits,
  }, nil
}
