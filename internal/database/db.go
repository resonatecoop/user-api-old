package database

import (
	"github.com/go-pg/pg"
)

var Conn *pg.DB

func init() {
	Conn = connect()
}

func connect() *pg.DB {
	return pg.Connect(&pg.Options{
    Addr: "127.0.0.1:5432",
    User:     "resonate_dev_user",
    Password: "password",
    Database: "resonate_dev",
	})
}
