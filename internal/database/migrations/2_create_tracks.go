package main

import (

  "github.com/go-pg/migrations"
	"github.com/go-pg/pg/orm"

  "user-api/internal/database/models"
)

func init() {
	migrations.Register(func(db migrations.DB) error {
    if _, err := db.Exec(`CREATE TYPE status AS ENUM ('paid', 'free');`); err != nil {
      return err
    }

		if _, err := orm.CreateTable(db, &models.Track{}, nil); err != nil {
			return err
		}
		return nil
	}, func(db migrations.DB) error {
    if _, err := orm.DropTable(db, &models.Track{}, nil); err != nil {
      return err
    }
    if _, err := db.Exec(`DROP TYPE IF EXISTS status;`); err != nil {
      return err
    }
    return nil
	})
}
