package utils

import (
	"encoding/json"

	"google.golang.org/protobuf/types/known/structpb"
)

func ConvertToProtoStruct(obj interface{}) (*structpb.Struct, error) {
	m := make(map[string]interface{})
	bs, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(bs, &m); err != nil {
		return nil, err
	}
	return structpb.NewStruct(m)
}

func MustConvertToProtoStruct(obj interface{}) *structpb.Struct {
	if obj == nil {
		return nil
	}
	m := make(map[string]interface{})
	bs, _ := json.Marshal(obj)
	json.Unmarshal(bs, &m)
	pb, _ := structpb.NewStruct(m)
	return pb
}

func MustConvertToProtoStructs(objs []interface{}) []*structpb.Struct {
	var pbs []*structpb.Struct
	for _, obj := range objs {
		m := make(map[string]interface{})
		bs, _ := json.Marshal(obj)
		json.Unmarshal(bs, &m)
		pb, _ := structpb.NewStruct(m)
		pbs = append(pbs, pb)
	}
	return pbs
}
