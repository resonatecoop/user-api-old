package main

import (
	"fmt"
	"log"
	"net/http"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"user-api/internal/userserver"
	"user-api/rpc/user"
)

const (
	host = "localhost"
	port = 5432
	username = "toy-api-dev-user"
	password = "password"
	dbname = "toy-api-dev"
)

func main() {
	fmt.Printf("Toy User Service on :8080")

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, username, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	server := userserver.NewServer(db)
	twirpHandler := user.NewToyUserServer(server, nil)

	mux := http.NewServeMux()
	mux.Handle(user.ToyUserPathPrefix, twirpHandler)

	// cors.Default() setup the middleware with default options being
	// all origins accepted with simple methods (GET, POST).
	handler := cors.Default().Handler(mux)

	http.ListenAndServe(":8080", handler)
	defer db.Close()
}
