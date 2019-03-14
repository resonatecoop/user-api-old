package main

import (
	"fmt"
	"context"
	"net/http"
	"github.com/rs/cors"

	userServer "user-api/internal/server/user"
	userGroupServer "user-api/internal/server/usergroup"
	trackServer "user-api/internal/server/track"
	tagServer "user-api/internal/server/tag"
	trackGroupServer "user-api/internal/server/trackgroup"
	addressServer "user-api/internal/server/address"

	userRpc "user-api/rpc/user"
	userGroupRpc "user-api/rpc/usergroup"
	trackRpc "user-api/rpc/track"
	tagRpc "user-api/rpc/tag"
	trackGroupRpc "user-api/rpc/trackgroup"
	addressRpc "user-api/rpc/address"
	"user-api/internal/database"
)

func WithURLQuery(base http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		query := r.URL.Query()
		ctx = context.WithValue(ctx, "query", query)
		r = r.WithContext(ctx)

		base.ServeHTTP(w, r)
	})
}

func main() {
	fmt.Printf("User Service on :8080")

	// psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, username, password, dbname)
	// db, err := sql.Open("postgres", psqlInfo)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	db := database.Connect(false)

	newUserServer := userServer.NewServer(db)
	userTwirpHandler := WithURLQuery(userRpc.NewUserServiceServer(newUserServer, nil))

	newUserGroupServer := userGroupServer.NewServer(db)
	userGroupTwirpHandler := WithURLQuery(userGroupRpc.NewUserGroupServiceServer(newUserGroupServer, nil))

	newTrackServer := trackServer.NewServer(db)
	trackTwirpHandler := trackRpc.NewTrackServiceServer(newTrackServer, nil)

	newTagServer := tagServer.NewServer(db)
	tagTwirpHandler := tagRpc.NewTagServiceServer(newTagServer, nil)

	newTrackGroupServer := trackGroupServer.NewServer(db)
	trackGroupTwirpHandler := trackGroupRpc.NewTrackGroupServiceServer(newTrackGroupServer, nil)

	newAddressServer := addressServer.NewServer("https://places-dsn.algolia.net/1/places/query", "", "")
	addressTwirpHandler := addressRpc.NewAddressServiceServer(newAddressServer, nil)

	mux := http.NewServeMux()
	mux.Handle(userRpc.UserServicePathPrefix, userTwirpHandler)
	mux.Handle(userGroupRpc.UserGroupServicePathPrefix, userGroupTwirpHandler)
	mux.Handle(trackRpc.TrackServicePathPrefix, trackTwirpHandler)
	mux.Handle(tagRpc.TagServicePathPrefix, tagTwirpHandler)
	mux.Handle(trackGroupRpc.TrackGroupServicePathPrefix, trackGroupTwirpHandler)
	mux.Handle(addressRpc.AddressServicePathPrefix, addressTwirpHandler)

	// cors.Default() setup the middleware with default options being
	// all origins accepted with simple methods (GET, POST).
	handler := cors.Default().Handler(mux)

	http.ListenAndServe(":8080", handler)
	defer db.Close()
}
