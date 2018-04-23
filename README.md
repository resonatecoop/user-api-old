# Toy User API

This is a simple User API that uses Twirp. It allows to list and create users with name and email fields. Learn more about
Twirp at its [website](https://twitchtv.github.io/twirp/docs/intro.html) or
[repo](https://github.com/twitchtv/twirp).

## Installation

Follow Twirp [installation guide](https://twitchtv.github.io/twirp/docs/install.html) in order to install the protobuf compiler `protoc` and Go and Twirp protoc plugins `protoc-gen-go` and `protoc-gen-twirp`

## Try it out

First, download and put this repo into `$GOPATH/src`

Then, run the server
```
go run ./cmd/server/main.go
```

Run the client
```
go run ./cmd/client/main.go
```

## Example curl requests

### CreateUser
```sh
curl --request "POST" \
     --location "http://localhost:8080/twirp/resonate.toyapi.user.ToyUser/CreateUser" \
     --header "Content-Type:application/json" \
     --data '{"name": "john", "email": "john@doe.com"}' \
     --verbose

{"name":"john","email":"john@doe.com"}
```

### GetUsers
```sh
curl --request "POST" \
     --location "http://localhost:8080/twirp/resonate.toyapi.user.ToyUser/GetUsers" \
     --header "Content-Type:application/json" \
     --data '{}' \
     --verbose

{"users":[{"name":"john","email":"john@doe.com"}]}
```

## Code structure

The protobuf definition for the service lives in
`rpc/user/service.proto`.
The generated Twirp and Go protobuf code is in the same directory.

The implementation of the server is in `internal/userserver`.

Finally, `cmd/server` and `cmd/client` wrap things together into executable main
packages.
