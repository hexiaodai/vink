package utils

import (
	"encoding/json"
	"unsafe"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// StringToBytes converts string to byte slice.
func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

func StructToString(s interface{}) string {
	data, err := json.Marshal(s)
	if err != nil {
		return err.Error()
	}
	return BytesToString(data)
}

func MarshalObj(obj client.Object) string {
	if obj == nil {
		return ""
	}
	obj.SetManagedFields(nil)
	return StructToString(obj)
}
