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
  "github.com/satori/go.uuid"
  "github.com/go-pg/pg"
  // "github.com/twitchtv/twirp"

  // "time"
  "fmt"
  // "reflect"
)

func main() {
  testing := true
  db := database.Connect(testing)


  newArtistGroupTaxonomy := &models.GroupTaxonomy{Type: "artist", Name: "Artist"}
  _ = db.Insert(newArtistGroupTaxonomy)

  newLabelGroupTaxonomy := &models.GroupTaxonomy{Type: "label", Name: "Label"}
  _ = db.Insert(newLabelGroupTaxonomy)

  newUser := &models.User{Username: "username", FullName: "full name", DisplayName: "display name", Email: "email@fake.com"}
  _ = db.Insert(newUser)
  //
  // // duration, _ := time.ParseDuration("10m10s")
  // // cover := make([]byte, 5)
  // // newTrack := &models.Track{PublishDate: time.Now(), Title: "track title", Duration: duration, Status: "free", Cover: cover}
  // // err= db.Insert(newTrack)
  //

  avatar := make([]byte, 5)
  artist := &models.UserGroup{
    DisplayName: "best artist ever",
    Avatar: avatar,
    OwnerId: newUser.Id,
    TypeId: newArtistGroupTaxonomy.Id,
    AddressId: uuid.NewV4(),
  }
  _, pgerr := db.Model(artist).Returning("*").Insert()

  db.Model(artist).
		Column("user_group.*", "Privacy", "Type", "Address").
		WherePK().
		Select()

  // artistId, _ := uuid.FromString("a85aae16-c1e9-4eab-ac39-3d23d0fe3728")
  // artist := &models.UserGroup{DisplayName: "worst artist ever", Id: artistId}
  // _, pgerr := db.Model(artist).WherePK().Returning("*").UpdateNotNull()
  //
  //
  // //
  // newLabel := &models.UserGroup{
  //   DisplayName: "best label ever",
  //   Avatar: avatar,
  //   OwnerId: newUser.Id,
  //   TypeId: newLabelGroupTaxonomy.Id,
  //   // AdminUsers: admins,
  // }
  // 	_ = db.Insert(newLabel)
  //   pgerr = db.Model(newArtist).
  //     Column("Privacy").
  //     WherePK().
  //     Select()
  //     fmt.Println(pgerr)
  // id, _ := uuid.FromString("b217a9b4-2844-439e-b7cb-50a928bca74")
  // user := models.User{Id: id}
  // // user := new(models.User)
  // pgerr := db.Model(&user).
  //     Column("user.*", "OwnerOfGroups").
  //     Select()
  // if pgerr != nil {
  //     panic(pgerr)
  // }
  //
  // fmt.Println(user.Id, user.Username, reflect.TypeOf(user.OwnerOfGroups))
  // fmt.Printf("err: %+v\n", pgerr)
  // fmt.Printf("artist: %+v\n", artist)
}
