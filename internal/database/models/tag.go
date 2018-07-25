package models

import (
  "user-api/internal"
  "github.com/satori/go.uuid"
  pb "user-api/rpc/usergroup"
  "github.com/go-pg/pg"
)

type Tag struct {
  Id uuid.UUID `sql:"type:uuid,default:uuid_generate_v4()"`
  Type string `sql:",notnull"`
  Name string `sql:",notnull"`
}

func GetTagIds(t []*pb.Tag, db *pg.Tx) ([]uuid.UUID, error) {
	tags := make([]*Tag, len(t))
	tagIds := make([]uuid.UUID, len(t))
	for i, tag := range(t) {
		if tag.Id == "" { // new tag to create and add
			tags[i] = &Tag{Type: tag.Type, Name: tag.Name}
			_, pgerr := db.Model(tags[i]).
				Where("type = ?", tags[i].Type).
				Where("lower(name) = lower(?)", tags[i].Name).
				Returning("*").
				SelectOrInsert()
			if pgerr != nil {
				return nil, pgerr
			}
			tagIds[i] = tags[i].Id
			tag.Id = tags[i].Id.String()
		} else {
			tagId, twerr := internal.GetUuidFromString(tag.Id)
			if twerr != nil {
				return nil, twerr.(error)
			}
			tagIds[i] = tagId
		}
	}
	return tagIds, nil
}
