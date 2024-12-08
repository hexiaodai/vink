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

// MarshalJSON is a custom marshaler for FieldSelector
func (this *FieldSelector) MarshalJSON() ([]byte, error) {
	str, err := TypesMarshaler.MarshalToString(this)
	return []byte(str), err
}

// UnmarshalJSON is a custom unmarshaler for FieldSelector
func (this *FieldSelector) UnmarshalJSON(b []byte) error {
	return TypesUnmarshaler.Unmarshal(bytes.NewReader(b), this)
}

// MarshalJSON is a custom marshaler for FieldSelectorGroup
func (this *FieldSelectorGroup) MarshalJSON() ([]byte, error) {
	str, err := TypesMarshaler.MarshalToString(this)
	return []byte(str), err
}

// UnmarshalJSON is a custom unmarshaler for FieldSelectorGroup
func (this *FieldSelectorGroup) UnmarshalJSON(b []byte) error {
	return TypesUnmarshaler.Unmarshal(bytes.NewReader(b), this)
}

var (
	TypesMarshaler   = &jsonpb.Marshaler{}
	TypesUnmarshaler = &jsonpb.Unmarshaler{AllowUnknownFields: true}
)
