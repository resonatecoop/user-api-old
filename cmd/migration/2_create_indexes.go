package main

import (
  "github.com/go-pg/migrations"
)

func init() {
	migrations.MustRegisterTx(func(db migrations.DB) error {
    if _, err := db.Exec(`
      CREATE FUNCTION f_arr2str(uuid[])
      RETURNS text LANGUAGE sql IMMUTABLE AS $$SELECT array_to_string($1, ',')$$
    `); err != nil {
      return err
    }
    if _, err := db.Exec(`
      CREATE INDEX user_groups_gin_idx ON user_groups
      USING GIN (to_tsvector('english', coalesce(display_name, '') || ' ' || COALESCE(f_arr2str(tags), '')))
    `); err != nil {
      return err
    }

    if _, err := db.Exec(`
      CREATE INDEX tags_gin_idx ON tags USING GIN (to_tsvector('english', coalesce(name, '')))
    `); err != nil {
      return err
    }

    if _, err := db.Exec(`
      CREATE INDEX tracks_gin_idx ON tracks
      USING GIN (to_tsvector('english', COALESCE(title, '') || ' ' || COALESCE(f_arr2str(tags), '')))
    `); err != nil {
      return err
    }

    if _, err := db.Exec(`
      CREATE INDEX track_groups_gin_idx ON track_groups
      USING GIN (to_tsvector('english', COALESCE(title, '') || ' ' || COALESCE(f_arr2str(tags), '')))
    `); err != nil {
      return err
    }

		return nil
	}, func(db migrations.DB) error {
    if _, err := db.Exec(`DROP INDEX IF EXISTS user_groups_gin_idx CASCADE`); err != nil {
      return err
    }
    if _, err := db.Exec(`DROP INDEX IF EXISTS tags_gin_idx CASCADE`); err != nil {
      return err
    }
    if _, err := db.Exec(`DROP INDEX IF EXISTS tracks_gin_idx CASCADE`); err != nil {
      return err
    }
    if _, err := db.Exec(`DROP INDEX IF EXISTS track_groups_gin_idx CASCADE`); err != nil {
      return err
    }

    if _, err := db.Exec(`DROP FUNCTION IF EXISTS f_arr2str`); err != nil {
      return err
    }


    return nil
	})
}
