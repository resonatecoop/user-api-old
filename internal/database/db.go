package database

import (
	"github.com/go-pg/pg"

	"user-api/config"
)

func Connect(testing bool) *pg.DB {
	db_config := config.Config.Dev
	if (testing) {
		db_config = config.Config.Testing
	}
	return pg.Connect(&pg.Options{
    Addr: db_config.Addr,
    User:     db_config.User,
    Password: db_config.Password,
    Database: db_config.Database,
	})
}
