package uuid

import (
	"github.com/satori/go.uuid"
  "github.com/twitchtv/twirp"
)

// ConvertUuidToStrArray returns a slice of uuids for given slive of strings
func ConvertUuidToStrArray(uuids []uuid.UUID) ([]string) {
  strArray := make([]string, len(uuids))
  for i := range uuids {
    strArray[i] = uuids[i].String()
  }
  return strArray
}

// GetUuidFromString returns id as string and returns error if not a valid uuid
func GetUuidFromString(id string) (uuid.UUID, twirp.Error) {
	uid, err := uuid.FromString(id)
	if err != nil {
		return uuid.UUID{}, twirp.InvalidArgumentError("id", "must be a valid uuid")
	}
	return uid, nil
}

// Difference returns difference between two slices of uuids
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

// Equal compares two uuid slices and returns true if equal
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

// RemoveDuplicates returns slices of uuid without duplicates
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
