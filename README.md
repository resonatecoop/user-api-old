# User/Track API

This is the main User/Track API.
It uses Twirp RPC framework. Learn more about
Twirp at its [website](https://twitchtv.github.io/twirp/docs/intro.html) or
[repo](https://github.com/twitchtv/twirp).
It also uses [go-pg](https://github.com/go-pg/pg) PostgreSQL ORM.

## Dev database setup

* Make sure you have latest PostgreSQL installed.
* Create user and database as follow:

username = "resonate-dev-user"

password = "password"

dbname = "resonate-dev"

Add following postgres extensions: "hstore" and "uuid-ossp"

* Run migrations from `./internal/database/migrations`

```sh
$ go run *.go
```

## Dependencies

[Dep](https://github.com/golang/dep) is used as dependency management tool.
`vendor/` folder contains project dependencies and should be in sync with `Gopkg.toml` and `Gopkg.lock`.

## Various tools installation for development

* [Install Protocol Buffers v3](https://developers.google.com/protocol-buffers/docs/gotutorial),
the `protoc` compiler that is used to auto-generate code. The simplest way to do
this is to download pre-compiled binaries for your platform from here:
https://github.com/google/protobuf/releases

It is also available in MacOS through Homebrew:

```sh
$ brew install protobuf
```

* Install [retool](https://github.com/twitchtv/retool). It helps manage go tools like commands and linters.
protoc-gen-go and protoc-gen-twirp plugins were installed into `_tools` folder using retool.

Build the generators and tool dependencies:
```sh
$ retool build
```

Then, to run the `protoc` command and autogenerate Go code for the server interface, make sure to prefix with `retool do`, for example:
```sh
$ retool do protoc --proto_path=$GOPATH/src:. --twirp_out=. --go_out=. ./rpc/user/service.proto
```

## Running/building the server

First, put this repo into `$GOPATH/src`

Then, run the server
```sh
$ go run ./cmd/server/main.go
```

Alternatively, you can build and run an executable binary
```sh
$ cd ./cmd/server/
$ go build
$ ./server
```

## Example curl requests

### CreateUser
```sh
curl --request "POST" \
     --location "http://localhost:8080/twirp/resonate.api.user.UserService/CreateUser" \
     --header "Content-Type:application/json" \
     --data '{"display_name": "john", "full_name": "john doe", "email": "john@doe.com", "username": "johnd"}' \
     --verbose

{"id":"9ef71770-7a1b-4a11-a81e-1b6d177a3598","username":"johnd","email":"john@doe.com","display_name":"john","full_name":"john doe"}
```

## Code structure

The protobuf definition for the service lives in
`rpc/user/service.proto`.
The generated Twirp and Go protobuf code is in the same directory.

The implementation of the server is in `internal/server`.
Database related stuff (migrations, model definitions) can be found in `internal/database`.

Finally, `cmd/server` and `cmd/client` wrap things together into executable main
packages.

## Documentation

Check out `doc.apib` for API documentation.

## Testing

We use Ginkgo and Gomega for testing.

At the moment, you need to create the testing database and run migrations manually before running tests.

* Create user and database as follow:

username = "resonate-testing-user"

password = ""

dbname = "resonate-testing"

Add following extensions: "hstore" and "uuid-ossp" (TODO: add them on initial migration)

* Run migrations from `./internal/database/migrations`

```sh
$ go run *.go testing
```

* Run tests from `./internal/server/user` or `./internal/server/usergroup`

```sh
$ go test
```

Or run all tests using ginkgo CLI from `./`

```sh
$ ginkgo -r
```
