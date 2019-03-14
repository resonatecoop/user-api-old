package internal

import (
  "fmt"
  "strings"

  "github.com/golang/protobuf/ptypes/struct"

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

func GetUuidFromString(id string) (uuid.UUID, twirp.Error) {
	uid, err := uuid.FromString(id)
	if err != nil {
		return uuid.UUID{}, twirp.InvalidArgumentError("id", "must be a valid uuid")
	}
	return uid, nil
}

func Difference(a, b []uuid.UUID) []uuid.UUID {
    mb := map[uuid.UUID]bool{}
    for _, x := range b {
        mb[x] = true
    }
    ab := []uuid.UUID{}
    for _, x := range a {
        if _, ok := mb[x]; !ok {
            ab = append(ab, x)
        }
    }
    return ab
}

func RemoveDuplicates(elements []uuid.UUID) []uuid.UUID {
  // Use map to record duplicates as we find them
  encountered := map[uuid.UUID]bool{}
  result := []uuid.UUID{}

  for v := range elements {
    if encountered[elements[v]] == true {
      // Do not add duplicate
    } else {
      // Record this element as an encountered element
      encountered[elements[v]] = true
      // Append to result slice
      result = append(result, elements[v])
    }
  }
  // Return the new slice.
  return result
}

// Compare two uuid slices
func Equal(a, b []uuid.UUID) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func GetMapListValue(m map[string][]string) (map[string]*structpb.ListValue) {
  mapListValue := make(map[string]*structpb.ListValue)
  for k, v := range m {
    mapListValue[k] = getListValue(v)
  }
  return mapListValue
}

func getListValue(strArr []string) (*structpb.ListValue) {
  values := make([]*structpb.Value, len(strArr))
  for i, v := range strArr {
    values[i] = &structpb.Value{Kind: &structpb.Value_StringValue{v}}
  }
  return &structpb.ListValue{
    Values: values,
  }
}
