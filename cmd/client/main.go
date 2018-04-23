package main

import (
    "context"
    "net/http"
    "os"
    "fmt"

    pb "user-api/rpc/user"
)

func main() {
    client := pb.NewToyUserProtobufClient("http://localhost:8080", &http.Client{})

    u1, err := client.CreateUser(context.Background(), &pb.User{Name: "john", Email: "john@doe.com"})
    u2, err := client.CreateUser(context.Background(), &pb.User{Name: "marie", Email: "marie@doe.com"})
    users, err := client.GetUsers(context.Background(), &pb.Empty{})
    if err != nil {
        fmt.Printf("oh no: %v", err)
        os.Exit(1)
    }
    fmt.Printf("New user created: %+v\n", u1)
    fmt.Printf("New user created: %+v\n", u2)
    fmt.Printf("Users: %+v\n", users)
}
