package proto

import (
	"github.com/kubevm.io/vink/apis/common"
)

var operatingSystemTypeToString = map[common.OperatingSystemType]string{
	common.OperatingSystemType_LINUX:   "linux",
	common.OperatingSystemType_WINDOWS: "windows",
	common.OperatingSystemType_CENTOS:  "centos",
	common.OperatingSystemType_UBUNTU:  "ubuntu",
	common.OperatingSystemType_DEBIAN:  "debian",
}

var operatingSystemTypeToEnum = map[string]common.OperatingSystemType{
	"linux":   common.OperatingSystemType_LINUX,
	"windows": common.OperatingSystemType_WINDOWS,
	"centos":  common.OperatingSystemType_CENTOS,
	"ubuntu":  common.OperatingSystemType_UBUNTU,
	"debian":  common.OperatingSystemType_DEBIAN,
}

func OperatingSystemTypeFromEnum(osType common.OperatingSystemType) string {
	return operatingSystemTypeToString[osType]
}

func OperatingSystemTypeFromString(osType string) common.OperatingSystemType {
	return operatingSystemTypeToEnum[osType]
}

var operatingSystemWindowsVersionToString = map[common.OperatingSystemWindowsVersion]string{
	common.OperatingSystemWindowsVersion_WINDOWS_10: "10",
	common.OperatingSystemWindowsVersion_WINDOWS_11: "11",
}

var operatingSystemWindowsVersionToEnum = map[string]common.OperatingSystemWindowsVersion{
	"10": common.OperatingSystemWindowsVersion_WINDOWS_10,
	"11": common.OperatingSystemWindowsVersion_WINDOWS_11,
}

func OperatingSystemWindowsVersionFromEnum(osVersion common.OperatingSystemWindowsVersion) string {
	return operatingSystemWindowsVersionToString[osVersion]
}

func OperatingSystemWindowsVersionFromString(osVersion string) common.OperatingSystemWindowsVersion {
	return operatingSystemWindowsVersionToEnum[osVersion]
}

var operatingSystemCentOSVersionToString = map[common.OperatingSystemCentOSVersion]string{
	common.OperatingSystemCentOSVersion_CENTOS_7: "7",
	common.OperatingSystemCentOSVersion_CENTOS_8: "8",
}

var operatingSystemCentOSVersionToEnum = map[string]common.OperatingSystemCentOSVersion{
	"7": common.OperatingSystemCentOSVersion_CENTOS_7,
	"8": common.OperatingSystemCentOSVersion_CENTOS_8,
}

func OperatingSystemCentOSVersionFromEnum(osVersion common.OperatingSystemCentOSVersion) string {
	return operatingSystemCentOSVersionToString[osVersion]
}

func OperatingSystemCentOSVersionFromString(osVersion string) common.OperatingSystemCentOSVersion {
	return operatingSystemCentOSVersionToEnum[osVersion]
}

var operatingSystemUbuntuVersionToString = map[common.OperatingSystemUbuntuVersion]string{
	common.OperatingSystemUbuntuVersion_UBUNTU_18_04: "18.04",
	common.OperatingSystemUbuntuVersion_UBUNTU_20_04: "20.04",
	common.OperatingSystemUbuntuVersion_UBUNTU_22_04: "22.04",
}

var operatingSystemUbuntuVersionToEnum = map[string]common.OperatingSystemUbuntuVersion{
	"18.04": common.OperatingSystemUbuntuVersion_UBUNTU_18_04,
	"20.04": common.OperatingSystemUbuntuVersion_UBUNTU_20_04,
	"22.04": common.OperatingSystemUbuntuVersion_UBUNTU_22_04,
}

func OperatingSystemUbuntuVersionFromEnum(osVersion common.OperatingSystemUbuntuVersion) string {
	return operatingSystemUbuntuVersionToString[osVersion]
}

func OperatingSystemUbuntuVersionFromString(osVersion string) common.OperatingSystemUbuntuVersion {
	return operatingSystemUbuntuVersionToEnum[osVersion]
}

var operatingSystemDebianVersionToString = map[common.OperatingSystemDebianVersion]string{
	common.OperatingSystemDebianVersion_DEBIAN_9:  "9",
	common.OperatingSystemDebianVersion_DEBIAN_10: "10",
	common.OperatingSystemDebianVersion_DEBIAN_11: "11",
}

var operatingSystemDebianVersionToEnum = map[string]common.OperatingSystemDebianVersion{
	"9":  common.OperatingSystemDebianVersion_DEBIAN_9,
	"10": common.OperatingSystemDebianVersion_DEBIAN_10,
	"11": common.OperatingSystemDebianVersion_DEBIAN_11,
}

func OperatingSystemDebianVersionFromEnum(osVersion common.OperatingSystemDebianVersion) string {
	return operatingSystemDebianVersionToString[osVersion]
}

func OperatingSystemDebianVersionFromString(osVersion string) common.OperatingSystemDebianVersion {
	return operatingSystemDebianVersionToEnum[osVersion]
}
