// Code generated by protoc-gen-jsonshim. DO NOT EDIT.
package types

import (
	bytes "bytes"
	jsonpb "github.com/golang/protobuf/jsonpb"
)

// MarshalJSON is a custom marshaler for NamespaceName
func (this *NamespaceName) MarshalJSON() ([]byte, error) {
	str, err := TypesMarshaler.MarshalToString(this)
	return []byte(str), err
}

// UnmarshalJSON is a custom unmarshaler for NamespaceName
func (this *NamespaceName) UnmarshalJSON(b []byte) error {
	return TypesUnmarshaler.Unmarshal(bytes.NewReader(b), this)
}

var (
	TypesMarshaler   = &jsonpb.Marshaler{}
	TypesUnmarshaler = &jsonpb.Unmarshaler{AllowUnknownFields: true}
)
