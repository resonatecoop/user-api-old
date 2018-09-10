package database

import (
	"github.com/go-pg/pg"

	"user-api/config"
)

func Connect(testing bool) *pg.DB {
	dbConfig := config.Configs.Dev
	if (testing) {
		dbConfig = config.Configs.Testing
	}
	return pg.Connect(&pg.Options{
    Addr: dbConfig.Addr,
    User:     dbConfig.User,
    Password: dbConfig.Password,
    Database: dbConfig.Database,
	})
}
