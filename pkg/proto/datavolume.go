package proto

import (
	"github.com/kubevm.io/vink/apis/management/datavolume/v1alpha1"
)

var dataVolumeTypeToString = map[v1alpha1.DataVolumeType]string{
	v1alpha1.DataVolumeType_IMAGE: "image",
	v1alpha1.DataVolumeType_ROOT:  "root",
	v1alpha1.DataVolumeType_DATA:  "data",
}

var dataVolumeTypeToEnum = map[string]v1alpha1.DataVolumeType{
	"image": v1alpha1.DataVolumeType_IMAGE,
	"root":  v1alpha1.DataVolumeType_ROOT,
	"data":  v1alpha1.DataVolumeType_DATA,
}

func DataVolumeTypeFromEnum(diskType v1alpha1.DataVolumeType) string {
	return dataVolumeTypeToString[diskType]
}

func DataVolumeTypeFromString(diskType string) v1alpha1.DataVolumeType {
	return dataVolumeTypeToEnum[diskType]
}

func DataVolumeTypeEqual(s string, e v1alpha1.DataVolumeType) bool {
	return DataVolumeTypeFromEnum(e) == s
}
