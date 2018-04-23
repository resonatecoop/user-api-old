package main

import (
	"fmt"
	"log"
	"net/http"

	"toy-api/internal/userserver"
	"toy-api/rpc/user"
)

func main() {
	fmt.Printf("Toy User Service on :8080")

	server := userserver.NewServer()
	twirpHandler := user.NewToyUserServer(server, nil)

	log.Fatal(http.ListenAndServe(":8080", twirpHandler))
}
