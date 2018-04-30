package main

import (
	"fmt"
	// "log"
	"net/http"
	"github.com/rs/cors"
	"user-api/internal/userserver"
	"user-api/rpc/user"
)

func main() {
	fmt.Printf("Toy User Service on :8080")

	server := userserver.NewServer()
	twirpHandler := user.NewToyUserServer(server, nil)

	mux := http.NewServeMux()
	mux.Handle(user.ToyUserPathPrefix, twirpHandler)

	// cors.Default() setup the middleware with default options being
	// all origins accepted with simple methods (GET, POST).
	handler := cors.Default().Handler(mux)

	http.ListenAndServe(":8080", handler)
}
