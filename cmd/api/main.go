package main

import (
	"flag"
	"log"
	"net/http"
	// "context"

	"user-api/pkg/zerolog"

	"user-api/internal/openapi"

	"user-api/internal/iam"

	"github.com/justinas/alice"

	"user-api/pkg/mw"

	// "user-api/pkg/hooks"

	// contextpkg "user-api/pkg/context"

	"github.com/gorilla/mux"

	// "user-api/internal/iam/rbac"
	"user-api/internal/iam/secure"

	"user-api/pkg/config"
	"user-api/pkg/jwt"
	"user-api/pkg/postgres"

	iamdb "user-api/internal/iam/platform/postgres"
	iampb "user-api/rpc/iam"

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
)

func main() {
	cfgPath := flag.String("p", "./conf.local.yaml", "Path to config file")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	checkErr(err)

	router := mux.NewRouter().StrictSlash(true)
	registerRoutes(router, cfg)

	openapi.New(router, cfg.OpenAPI.Username, cfg.OpenAPI.Password)

	mws := alice.New(mw.CORS, mw.AuthContext, mw.WithURLQuery)

	http.ListenAndServe(cfg.Server.Port, mws.Then(router))
}

func registerRoutes(r *mux.Router, cfg *config.Configuration) {
	db, err := pgsql.New(cfg.DB.Dev.PSN, cfg.DB.Dev.LogQueries, cfg.DB.Dev.TimeoutSeconds)
	checkErr(err)

	// rbacSvc := new(rbac.Service)
	// ctxSvc := new(contextpkg.Service)
	log := zerolog.New()

	secureSvc := secure.New(cfg.App.MinPasswordStrength)

	j := jwt.New(cfg.JWT.Secret, cfg.JWT.Duration, cfg.JWT.Algorithm)

	newUserServer := userServer.NewServer(db)
	userTwirpHandler := userRpc.NewUserServiceServer(newUserServer, nil)
	// userTwirpHandler := userRpc.NewUserServiceServer(newUserServer, hooks.WithJWTAuth(j))
	// userSvc := user.NewLoggingService(
	// 	user.New(db, userdb.NewUser(), rbacSvc, secureSvc, ctxSvc), log)
	r.PathPrefix(userRpc.UserServicePathPrefix).Handler(userTwirpHandler)

	newUserGroupServer := userGroupServer.NewServer(db)
	userGroupTwirpHandler := userGroupRpc.NewUserGroupServiceServer(newUserGroupServer, nil)
	r.PathPrefix(userGroupRpc.UserGroupServicePathPrefix).Handler(userGroupTwirpHandler)

	newTrackServer := trackServer.NewServer(db)
	trackTwirpHandler := trackRpc.NewTrackServiceServer(newTrackServer, nil)
	r.PathPrefix(trackRpc.TrackServicePathPrefix).Handler(trackTwirpHandler)

	newTagServer := tagServer.NewServer(db)
	tagTwirpHandler := tagRpc.NewTagServiceServer(newTagServer, nil)
	r.PathPrefix(tagRpc.TagServicePathPrefix).Handler(tagTwirpHandler)

	newTrackGroupServer := trackGroupServer.NewServer(db)
	trackGroupTwirpHandler := trackGroupRpc.NewTrackGroupServiceServer(newTrackGroupServer, nil)
	r.PathPrefix(trackGroupRpc.TrackGroupServicePathPrefix).Handler(trackGroupTwirpHandler)

	// TODO add algolia variables to config file
	newAddressServer := addressServer.NewServer("https://places-dsn.algolia.net/1/places/query", "", "")
	addressTwirpHandler := addressRpc.NewAddressServiceServer(newAddressServer, nil)
	r.PathPrefix(addressRpc.AddressServicePathPrefix).Handler(addressTwirpHandler)

	iamSvc := iam.NewLoggingService(iam.New(db, j, iamdb.NewUser(), secureSvc), log)

	r.PathPrefix(iampb.IAMPathPrefix).Handler(
		iampb.NewIAMServer(iamSvc, nil))
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
