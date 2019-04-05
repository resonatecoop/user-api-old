package maps

import (
  "github.com/golang/protobuf/ptypes/struct"
)

// GetMapListValue converts a map of []string to protobuf type *structpb.ListValue
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
