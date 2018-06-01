package main

import (
	"fmt"
	"net/http"
	"github.com/rs/cors"
	"user-api/internal/userserver"
	"user-api/rpc/user"
	"user-api/internal/database"
)

func main() {
	fmt.Printf("User Service on :8080")

	// psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, username, password, dbname)
	// db, err := sql.Open("postgres", psqlInfo)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	db := database.Conn
	server := userserver.NewServer(db)
	twirpHandler := user.NewUserServiceServer(server, nil)

	mux := http.NewServeMux()
	mux.Handle(user.UserServicePathPrefix, twirpHandler)

	// cors.Default() setup the middleware with default options being
	// all origins accepted with simple methods (GET, POST).
	handler := cors.Default().Handler(mux)

	http.ListenAndServe(":8080", handler)
	defer db.Close()
}
