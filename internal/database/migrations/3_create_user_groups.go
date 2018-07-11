package main

import (

  "github.com/go-pg/migrations"
	"github.com/go-pg/pg/orm"

  "user-api/internal/database/models"
)

func init() {
	migrations.Register(func(db migrations.DB) error {
    if _, err := orm.CreateTable(db, &models.Link{}, nil); err != nil {
      return err
    }
    if _, err := orm.CreateTable(db, &models.UserGroupPrivacy{}, nil); err != nil {
      return err
    }
		if _, err := orm.CreateTable(db, &models.GroupTaxonomy{}, nil); err != nil {
			return err
		}
    if _, err := orm.CreateTable(db, &models.UserGroup{}, nil); err != nil {
      return err
    }
		return nil
	}, func(db migrations.DB) error {
    if _, err := orm.DropTable(db, &models.UserGroup{}, nil); err != nil {
      return err
    }
    if _, err := orm.DropTable(db, &models.GroupTaxonomy{}, nil); err != nil {
      return err
    }
    if _, err := orm.DropTable(db, &models.UserGroupPrivacy{}, nil); err != nil {
      return err
    }
    if _, err := orm.DropTable(db, &models.Link{}, nil); err != nil {
      return err
    }


    return nil
	})
}
