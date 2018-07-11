package internal

import (
  "fmt"
  "strings"

	"github.com/satori/go.uuid"
  "github.com/go-pg/pg"
  "github.com/twitchtv/twirp"
)

func ConvertUuidToStrArray(uuids []uuid.UUID) ([]string) {
  strArray := make([]string, len(uuids))
  for i := range uuids {
    strArray[i] = uuids[i].String()
  }
  return strArray
}

func CheckError(err error, table string) (twirp.Error) {
	if err != nil {
		if err == pg.ErrNoRows {
			return twirp.NotFoundError(fmt.Sprintf("%s does not exist", table))
		}
    twerr, ok := err.(twirp.Error)
    if ok && twerr.Meta("argument") == "id" {
      return twerr
    }
		pgerr, ok := err.(pg.Error)
		if ok {
			code := pgerr.Field('C')
			name := pgerr.Field('n')
			var message string
			if code == "23505" { // unique_violation
				message = strings.TrimPrefix(strings.TrimSuffix(name, "_key"), fmt.Sprintf("%ss_", table))
				return twirp.NewError("already_exists", message)
			} else {
				message = pgerr.Field('M')
				return twirp.NewError("unknown", message)
			}
		}
		return twirp.NewError("unknown", err.Error())
	}
	return nil
}

func GetUuidFromString(id string) (uuid.UUID, twirp.Error) {
	uid, err := uuid.FromString(id)
	if err != nil {
		return uuid.UUID{}, twirp.InvalidArgumentError("id", "must be a valid uuid")
	}
	return uid, nil
}
