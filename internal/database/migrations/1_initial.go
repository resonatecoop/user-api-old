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
    if _, err := db.Exec(`CREATE TYPE status AS ENUM ('paid', 'free');`); err != nil {
      return err
    }
    for _, model := range []interface{}{
      &models.StreetAddress{},
      &models.Tag{},
      &models.User{},
      &models.Link{},
      &models.UserGroupPrivacy{},
      &models.GroupTaxonomy{},
      &models.UserGroup{},
      &models.Track{},
      &models.TrackGroup{},
      &models.UserGroupMember{},
    } {
      _, err := orm.CreateTable(db, model, &orm.CreateTableOptions{
        FKConstraints: true,
        IfNotExists: true,
      })
      if err != nil {
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
    if _, err := db.Exec(`DROP TYPE IF EXISTS status CASCADE;`); err != nil {
      return err
    }
    for _, model := range []interface{}{
      &models.Tag{},
      &models.TrackGroup{},
      &models.Track{},
      &models.GroupTaxonomy{},
      &models.UserGroupMember{},
      &models.StreetAddress{},
      &models.UserGroupPrivacy{},
      &models.UserGroup{},
      &models.User{},
      &models.Link{},

      } {
      _, err := orm.DropTable(db, model, &orm.DropTableOptions{
        IfExists: true,
        Cascade:  true,
      })
      if err != nil {
        return err
      }
    }


    return nil
	})
}
