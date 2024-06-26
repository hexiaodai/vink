package business

import (
	"errors"
	"strconv"

	"github.com/kubevm.io/vink/pkg/proto"

	"github.com/kubevm.io/vink/apis/annotation"
	"github.com/kubevm.io/vink/apis/common"
	"github.com/kubevm.io/vink/apis/label"
	dvv1alpha1 "github.com/kubevm.io/vink/apis/management/datavolume/v1alpha1"
	"github.com/kubevm.io/vink/pkg/utils"
	"google.golang.org/protobuf/types/known/timestamppb"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	cdiv1beta1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
)

var (
	defaultAccessMode = corev1.ReadWriteOnce
)

func crdToAPIDataVolume(dv *cdiv1beta1.DataVolume) (*dvv1alpha1.DataVolume, error) {
	datavolume := dvv1alpha1.DataVolume{
		Namespace:         dv.Namespace,
		Name:              dv.Name,
		DataVolume:        utils.MustConvertToProtoStruct(dv),
		CreationTimestamp: timestamppb.New(dv.CreationTimestamp.Time),
	}
	return &datavolume, nil
}

func generateDataVolumeCRD(namespace, name string, config *dvv1alpha1.DataVolumeConfig) (*cdiv1beta1.DataVolume, error) {
	dataVolumeSource := cdiv1beta1.DataVolumeSource{}
	switch v := config.DataSource.DataSource.(type) {
	case *dvv1alpha1.DataVolumeConfig_DataSource_Http_:
		dataVolumeSource.HTTP = &cdiv1beta1.DataVolumeSourceHTTP{
			URL: v.Http.Url,
			// FIXME:
			// ExtraHeaders: ,
		}
	case *dvv1alpha1.DataVolumeConfig_DataSource_Registry_:
		dataVolumeSource.Registry = &cdiv1beta1.DataVolumeSourceRegistry{
			URL: &v.Registry.Url,
		}
	case *dvv1alpha1.DataVolumeConfig_DataSource_S3_:
		dataVolumeSource.S3 = &cdiv1beta1.DataVolumeSourceS3{
			URL: v.S3.Url,
		}
	case *dvv1alpha1.DataVolumeConfig_DataSource_Blank_:
		dataVolumeSource.Blank = &cdiv1beta1.DataVolumeBlankImage{}
	default:
		return nil, errors.New("unsupported data source type")
	}

	dvcrd := cdiv1beta1.DataVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				label.VinkDatavolumeType.Name: proto.DataVolumeTypeFromEnum(config.DataVolumeType),
			},
			Annotations: map[string]string{
				annotation.IoKubevirtCdiStorageBindImmediateRequested.Name: strconv.FormatBool(true),
			},
		},
		Spec: cdiv1beta1.DataVolumeSpec{
			Source: &dataVolumeSource,
			PVC: &corev1.PersistentVolumeClaimSpec{
				AccessModes: []corev1.PersistentVolumeAccessMode{defaultAccessMode},
				Resources: corev1.VolumeResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceStorage: resource.MustParse(config.BoundPvc.Capacity),
					},
				},
				StorageClassName: &config.BoundPvc.StorageClassName,
			},
		},
	}

	if config.DataVolumeType == dvv1alpha1.DataVolumeType_ROOT {
		dvcrd.Labels[label.VinkVirtualmachineOs.Name] = proto.OperatingSystemTypeFromEnum(config.OperatingSystem.Type)
		var version string
		switch config.OperatingSystem.Type {
		case common.OperatingSystemType_WINDOWS:
			version = proto.OperatingSystemWindowsVersionFromEnum(config.OperatingSystem.GetWindows())
		case common.OperatingSystemType_CENTOS:
			version = proto.OperatingSystemCentOSVersionFromEnum(config.OperatingSystem.GetCentos())
		case common.OperatingSystemType_UBUNTU:
			version = proto.OperatingSystemUbuntuVersionFromEnum(config.OperatingSystem.GetUbuntu())
		case common.OperatingSystemType_DEBIAN:
			version = proto.OperatingSystemDebianVersionFromEnum(config.OperatingSystem.GetDebian())
		default:
		}
		dvcrd.Labels[label.VinkVirtualmachineVersion.Name] = version
	}

	return &dvcrd, nil
}
