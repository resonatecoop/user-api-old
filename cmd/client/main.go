package main
//
// import (
//     "context"
//     "net/http"
//     "os"
//     "fmt"
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
  // "time"
  "fmt"
)

func main() {
  var testing bool
  db := database.Connect(testing)

  newArtistGroupTaxonomy := &models.GroupTaxonomy{Type: "artist", Name: "Artist"}
  _ = db.Insert(newArtistGroupTaxonomy)

  newLabelGroupTaxonomy := &models.GroupTaxonomy{Type: "label", Name: "Label"}
  _ = db.Insert(newLabelGroupTaxonomy)

  newUser := &models.User{Username: "username", FullName: "full name", DisplayName: "display name", Email: "email@fake.com"}
  err := db.Insert(newUser)

  // duration, _ := time.ParseDuration("10m10s")
  // cover := make([]byte, 5)
  // newTrack := &models.Track{PublishDate: time.Now(), Title: "track title", Duration: duration, Status: "free", Cover: cover}
  // err= db.Insert(newTrack)

  avatar := make([]byte, 5)
  // admins := []uuid.UUID{newUser.Id}
  newArtist := &models.UserGroup{
    DisplayName: "best artist ever",
    Avatar: avatar,
    OwnerId: newUser.Id,
    TypeId: newArtistGroupTaxonomy.Id,
    // AdminUsers: admins,
  }
  err = db.Insert(newArtist)

  newLabel := &models.UserGroup{
    DisplayName: "best label ever",
    Avatar: avatar,
    OwnerId: newUser.Id,
    TypeId: newLabelGroupTaxonomy.Id,
    // AdminUsers: admins,
  }
  	err = db.Insert(newLabel)

  fmt.Println(err)
}
