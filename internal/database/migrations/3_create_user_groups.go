package main

import (

  "github.com/go-pg/migrations"
	"github.com/go-pg/pg/orm"

  "user-api/internal/database/models"
)

func init() {
	migrations.Register(func(db migrations.DB) error {
    for _, model := range []interface{}{
      &models.Link{},
      &models.UserGroupPrivacy{},
      &models.GroupTaxonomy{},
      &models.UserGroup{},
      &models.UserGroupMember{},
    } {
      if _, err := orm.CreateTable(db, model, nil); err != nil {
        return err
      }
    }
    orm.RegisterTable((*models.UserGroupMember)(nil))
    if _, err := db.Exec(`alter table user_group_members add foreign key (user_group_id) references user_groups(id)`); err != nil {
      return err
    }
    if _, err := db.Exec(`alter table user_group_members add foreign key (member_id) references user_groups(id)`); err != nil {
      return err
    }
		return nil
	}, func(db migrations.DB) error {
    for _, model := range []interface{}{
      &models.Link{},
      &models.UserGroupPrivacy{},
      &models.GroupTaxonomy{},
      &models.UserGroup{},
      &models.UserGroupMember{},
    } {
      if _, err := orm.DropTable(db, model, nil); err != nil {
        return err
      }
    }
    return nil
	})
}
