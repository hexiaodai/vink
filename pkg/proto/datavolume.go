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


var operatingSystemTypeToString = map[v1alpha1.OperatingSystemType]string{
	v1alpha1.OperatingSystemType_LINUX:   "linux",
	v1alpha1.OperatingSystemType_WINDOWS: "windows",
	v1alpha1.OperatingSystemType_CENTOS:  "centos",
	v1alpha1.OperatingSystemType_UBUNTU:  "ubuntu",
	v1alpha1.OperatingSystemType_DEBIAN:  "debian",
}

var operatingSystemTypeToEnum = map[string]v1alpha1.OperatingSystemType{
	"linux":   v1alpha1.OperatingSystemType_LINUX,
	"windows": v1alpha1.OperatingSystemType_WINDOWS,
	"centos":  v1alpha1.OperatingSystemType_CENTOS,
	"ubuntu":  v1alpha1.OperatingSystemType_UBUNTU,
	"debian":  v1alpha1.OperatingSystemType_DEBIAN,
}

func OperatingSystemTypeFromEnum(osType v1alpha1.OperatingSystemType) string {
	return operatingSystemTypeToString[osType]
}

func OperatingSystemTypeFromString(osType string) v1alpha1.OperatingSystemType {
	return operatingSystemTypeToEnum[osType]
}

var operatingSystemWindowsVersionToString = map[v1alpha1.OperatingSystemWindowsVersion]string{
	v1alpha1.OperatingSystemWindowsVersion_WINDOWS_10: "10",
	v1alpha1.OperatingSystemWindowsVersion_WINDOWS_11: "11",
}

var operatingSystemWindowsVersionToEnum = map[string]v1alpha1.OperatingSystemWindowsVersion{
	"10": v1alpha1.OperatingSystemWindowsVersion_WINDOWS_10,
	"11": v1alpha1.OperatingSystemWindowsVersion_WINDOWS_11,
}

func OperatingSystemWindowsVersionFromEnum(osVersion v1alpha1.OperatingSystemWindowsVersion) string {
	return operatingSystemWindowsVersionToString[osVersion]
}

func OperatingSystemWindowsVersionFromString(osVersion string) v1alpha1.OperatingSystemWindowsVersion {
	return operatingSystemWindowsVersionToEnum[osVersion]
}

var operatingSystemCentOSVersionToString = map[v1alpha1.OperatingSystemCentOSVersion]string{
	v1alpha1.OperatingSystemCentOSVersion_CENTOS_7: "7",
	v1alpha1.OperatingSystemCentOSVersion_CENTOS_8: "8",
}

var operatingSystemCentOSVersionToEnum = map[string]v1alpha1.OperatingSystemCentOSVersion{
	"7": v1alpha1.OperatingSystemCentOSVersion_CENTOS_7,
	"8": v1alpha1.OperatingSystemCentOSVersion_CENTOS_8,
}

func OperatingSystemCentOSVersionFromEnum(osVersion v1alpha1.OperatingSystemCentOSVersion) string {
	return operatingSystemCentOSVersionToString[osVersion]
}

func OperatingSystemCentOSVersionFromString(osVersion string) v1alpha1.OperatingSystemCentOSVersion {
	return operatingSystemCentOSVersionToEnum[osVersion]
}

var operatingSystemUbuntuVersionToString = map[v1alpha1.OperatingSystemUbuntuVersion]string{
	v1alpha1.OperatingSystemUbuntuVersion_UBUNTU_18_04: "18.04",
	v1alpha1.OperatingSystemUbuntuVersion_UBUNTU_20_04: "20.04",
	v1alpha1.OperatingSystemUbuntuVersion_UBUNTU_22_04: "22.04",
}

var operatingSystemUbuntuVersionToEnum = map[string]v1alpha1.OperatingSystemUbuntuVersion{
	"18.04": v1alpha1.OperatingSystemUbuntuVersion_UBUNTU_18_04,
	"20.04": v1alpha1.OperatingSystemUbuntuVersion_UBUNTU_20_04,
	"22.04": v1alpha1.OperatingSystemUbuntuVersion_UBUNTU_22_04,
}

func OperatingSystemUbuntuVersionFromEnum(osVersion v1alpha1.OperatingSystemUbuntuVersion) string {
	return operatingSystemUbuntuVersionToString[osVersion]
}

func OperatingSystemUbuntuVersionFromString(osVersion string) v1alpha1.OperatingSystemUbuntuVersion {
	return operatingSystemUbuntuVersionToEnum[osVersion]
}

var operatingSystemDebianVersionToString = map[v1alpha1.OperatingSystemDebianVersion]string{
	v1alpha1.OperatingSystemDebianVersion_DEBIAN_9:  "9",
	v1alpha1.OperatingSystemDebianVersion_DEBIAN_10: "10",
	v1alpha1.OperatingSystemDebianVersion_DEBIAN_11: "11",
}

var operatingSystemDebianVersionToEnum = map[string]v1alpha1.OperatingSystemDebianVersion{
	"9":  v1alpha1.OperatingSystemDebianVersion_DEBIAN_9,
	"10": v1alpha1.OperatingSystemDebianVersion_DEBIAN_10,
	"11": v1alpha1.OperatingSystemDebianVersion_DEBIAN_11,
}

func OperatingSystemDebianVersionFromEnum(osVersion v1alpha1.OperatingSystemDebianVersion) string {
	return operatingSystemDebianVersionToString[osVersion]
}

func OperatingSystemDebianVersionFromString(osVersion string) v1alpha1.OperatingSystemDebianVersion {
	return operatingSystemDebianVersionToEnum[osVersion]
}
