package model

import (
  "user-api/internal"
  "github.com/satori/go.uuid"
  pb "user-api/rpc/tag"
  "github.com/go-pg/pg"
  "github.com/twitchtv/twirp"
)

type Tag struct {
  Id uuid.UUID `sql:"type:uuid,default:uuid_generate_v4()"`
  Type string `sql:",notnull"`
  Name string `sql:",notnull"`
}

func SearchTags(query string, tagType string, db *pg.DB,) ([]*Tag, twirp.Error) {
  var tags []*Tag

  pgerr := db.Model(&tags).
    ColumnExpr("tag.id").
    Where("to_tsvector('english'::regconfig, COALESCE(name, '')) @@ (plainto_tsquery('english'::regconfig, ?)) = true", query).
    Where("type = ?", tagType).
    Select()
  if pgerr != nil {
    return nil, internal.CheckError(pgerr, "tag")
  }
  return tags, nil
}

func GetTags(tagIds []uuid.UUID, db *pg.DB) ([]*pb.Tag, twirp.Error) {
  tags := make([]*pb.Tag, len(tagIds))
  if len(tags) > 0 {
    var t []Tag
    pgerr := db.Model(&t).
      Where("id in (?)", pg.In(tagIds)).
      Select()
    if pgerr != nil {
      return nil, internal.CheckError(pgerr, "tag")
    }
    for i, tag := range t {
      tags[i] = &pb.Tag{Id: tag.Id.String(), Type: tag.Type, Name: tag.Name}
    }
  }
  return tags, nil
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
