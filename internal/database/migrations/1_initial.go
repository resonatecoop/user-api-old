package main

import (

  "github.com/go-pg/migrations"
	"github.com/go-pg/pg/orm"

  "user-api/internal/database/models"
)

func init() {
	migrations.Register(func(db migrations.DB) error {
		// if _, err := db.Exec( /* language=sql */ `CREATE EXTENSION IF NOT EXISTS "hstore"`); err != nil {
		// 	return err
		// }
    //
    // if _, err := db.Exec( /* language=sql */ `CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`); err != nil {
    //   return err
    // }
    for _, model := range []interface{}{
      &models.StreetAddress{},
      &models.Tag{},
      &models.User{},
    } {
      if _, err := orm.CreateTable(db, model, nil); err != nil {
        return err
      }
    }
		return nil
	}, func(db migrations.DB) error {
    for _, model := range []interface{}{
      &models.StreetAddress{},
      &models.Tag{},
      &models.User{},
      } {
      if _, err := orm.DropTable(db, model, nil); err != nil {
        return err
      }
    }

    return nil
	})
}
