package error

import (
  "fmt"
  "strings"

  "github.com/go-pg/pg"
  "github.com/twitchtv/twirp"
)

// CheckError receives error and table name and returns corresponding twirp error
func CheckError(err error, table string) (twirp.Error) {
	if err != nil {
		if err == pg.ErrNoRows {
			return twirp.NotFoundError(fmt.Sprintf("%s does not exist", table))
		}
    twerr, ok := err.(twirp.Error)
    if ok {
      argument := twerr.Meta("argument")
      if table != "" {
        argument = table + " " + argument
      }
      if twerr.Meta("argument") == "id" {
        return twirp.InvalidArgumentError(argument, "must be a valid uuid")
      }
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
      } else if code == "23503" { // foreign_key_violation
        message = pgerr.Field('M')
        return twirp.NotFoundError(message)
      } else {
				message = pgerr.Field('M')
        fmt.Println(twirp.NewError("unknown", message))
				return twirp.NewError("unknown", message)
			}
		}
		return twirp.NewError("unknown", err.Error())
	}
	return nil
}
