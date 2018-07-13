package main
//
// import (
//     "context"
//     "net/http"
//     "os"
    // "fmt"
//
//     pb "user-api/rpc/user"
// )
//
// func main() {
//     client := pb.NewUserServiceProtobufClient("http://localhost:8080", &http.Client{})
//
//     u1, err := client.CreateUser(context.Background(), &pb.User{Username: "janed", FullName: "jane doe", DisplayName: "jad", Email: "jane@d.com"})
//     // u2, err := client.CreateUser(context.Background(), &pb.User{Name: "marie", Email: "marie@doe.com"})
//     // users, err := client.GetUsers(context.Background(), &pb.Empty{})
//     if err != nil {
//         fmt.Printf("oh no: %v", err)
//         os.Exit(1)
//     }
//     fmt.Printf("New user created: %+v\n", u1)
//     // fmt.Printf("New user created: %+v\n", u2)
//     // fmt.Printf("Users: %+v\n", users)
// }

import (
  "user-api/internal/database"
  "user-api/internal/database/models"
  pb "user-api/rpc/usergroup"
  "github.com/satori/go.uuid"
  // "github.com/go-pg/pg"
  // "github.com/twitchtv/twirp"

  // "time"
  // "fmt"
  // "reflect"
)

func main() {
  testing := true
  db := database.Connect(testing)

  privacyId, _ := uuid.FromString("b5d58ac6-dbed-413a-9ca5-6c683ec52063")
  p := &pb.Privacy{Private: false}
  privacy := &models.UserGroupPrivacy{Id: privacyId, Private: true, OwnedTracks: p.OwnedTracks, SupportedArtists: p.SupportedArtists}
  _, _ = db.Model(privacy).WherePK().Returning("*").UpdateNotNull()
  // newArtistGroupTaxonomy := &models.GroupTaxonomy{Type: "artist", Name: "Artist"}
  // _ = db.Insert(newArtistGroupTaxonomy)
  //
  // newLabelGroupTaxonomy := &models.GroupTaxonomy{Type: "label", Name: "Label"}
  // _ = db.Insert(newLabelGroupTaxonomy)
  //
  // newUser := &models.User{Username: "username", FullName: "full name", DisplayName: "display name", Email: "email@fake.com"}
  // _ = db.Insert(newUser)
  //
  // newLink := &models.Link{Platform: "fakebook", Uri: "https://fakebook.com/bestartist"}
  // _ = db.Insert(newLink)
  //
  // // Create tag
  // newTag := &models.Tag{Type: "genre", Name: "rock"}
  // _ = db.Insert(newTag)
  //
  // newAddress := &models.StreetAddress{Data: map[string]string{"some": "data"}}
  // _ = db.Insert(newAddress)

  // // duration, _ := time.ParseDuration("10m10s")
  // // cover := make([]byte, 5)
  // // newTrack := &models.Track{PublishDate: time.Now(), Title: "track title", Duration: duration, Status: "free", Cover: cover}
  // // err= db.Insert(newTrack)

  // links := []uuid.UUID{newLink.Id}
  // tags := []uuid.UUID{newTag.Id}
  // avatar := make([]byte, 5)
  // artist := &models.UserGroup{
  //   DisplayName: "best artist ever",
  //   Avatar: avatar,
  //   OwnerId: newUser.Id,
  //   TypeId: newArtistGroupTaxonomy.Id,
  //   AddressId: newAddress.Id,
  //   Links: links,
  //   Tags: tags,
  // }
  // _, _ = db.Model(artist).Returning("*").Insert()
  //
  // db.Model(artist).
	// 	Column("user_group.*", "Privacy", "Type", "Address").
	// 	WherePK().
	// 	Select()
}
