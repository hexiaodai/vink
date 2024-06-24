package business

import (
	"errors"
	"fmt"

	"github.com/kubevm.io/vink/apis/annotation"
	"github.com/kubevm.io/vink/apis/label"
	dvv1alpha1 "github.com/kubevm.io/vink/apis/management/datavolume/v1alpha1"
	"github.com/kubevm.io/vink/pkg/utils"
	"google.golang.org/protobuf/types/known/timestamppb"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	cdiv1beta1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
)

const (
	bootDisk = "boot"
	dataDisk = "data"
)

const (
	centos  = "centos"
	ubuntu  = "ubuntu"
	debian  = "debian"
	windows = "windows"
)

var (
	defaultStorageClass = "local-storage"
	defaultAccessMode   = corev1.ReadWriteOnce
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

// func crdToAPIDataVolume(in *cdiv1beta1.DataVolume) (*dvv1alpha1.DataVolume, error) {
// 	config := dvv1alpha1.DataVolumeConfig{}

// 	switch in.Labels[label.IoVinkDisk.Name] {
// 	case bootDisk:
// 		config.Disk = dvv1alpha1.DataVolumeConfig_BOOT
// 	case dataDisk:
// 		config.Disk = dvv1alpha1.DataVolumeConfig_DATA
// 	}

// 	var osFamily dvv1alpha1.DataVolumeConfig_OSFamily
// 	switch in.Labels[label.IoVinkOsFamily.Name] {
// 	case centos:
// 		osFamily.OsFamily = &dvv1alpha1.DataVolumeConfig_OSFamily_Centos_{
// 			Centos: &dvv1alpha1.DataVolumeConfig_OSFamily_Centos{
// 				Version: in.Labels[label.IoVinkOsFamily.Name],
// 			},
// 		}
// 	case ubuntu:
// 		osFamily.OsFamily = &dvv1alpha1.DataVolumeConfig_OSFamily_Ubuntu_{
// 			Ubuntu: &dvv1alpha1.DataVolumeConfig_OSFamily_Ubuntu{
// 				Version: in.Labels[label.IoVinkOsFamily.Name],
// 			},
// 		}
// 	case debian:
// 		osFamily.OsFamily = &dvv1alpha1.DataVolumeConfig_OSFamily_Debian_{
// 			Debian: &dvv1alpha1.DataVolumeConfig_OSFamily_Debian{
// 				Version: in.Labels[label.IoVinkOsFamily.Name],
// 			},
// 		}
// 	case windows:
// 		osFamily.OsFamily = &dvv1alpha1.DataVolumeConfig_OSFamily_Windows_{
// 			Windows: &dvv1alpha1.DataVolumeConfig_OSFamily_Windows{
// 				Version: in.Labels[label.IoVinkOsFamily.Name],
// 			},
// 		}
// 	}
// 	config.OsFamily = &osFamily

// 	var dataSource dvv1alpha1.DataVolumeConfig_DataSource
// 	switch {
// 	case in.Spec.Source.HTTP != nil:
// 		dataSource.DataSource = &dvv1alpha1.DataVolumeConfig_DataSource_Http_{
// 			Http: &dvv1alpha1.DataVolumeConfig_DataSource_Http{
// 				Url: in.Spec.Source.HTTP.URL,
// 				// FIXME:
// 				// Headers: ,
// 			},
// 		}
// 	case in.Spec.Source.S3 != nil:
// 		dataSource.DataSource = &dvv1alpha1.DataVolumeConfig_DataSource_S3_{
// 			S3: &dvv1alpha1.DataVolumeConfig_DataSource_S3{
// 				Url: in.Spec.Source.S3.URL,
// 				// Bucket: in.Spec.Source.S3.Bucket,
// 			},
// 		}
// 	case in.Spec.Source.Registry != nil:
// 		dataSource.DataSource = &dvv1alpha1.DataVolumeConfig_DataSource_Registry_{
// 			Registry: &dvv1alpha1.DataVolumeConfig_DataSource_Registry{
// 				Url: lo.FromPtr(in.Spec.Source.Registry.URL),
// 			},
// 		}
// 	case in.Spec.Source.Upload != nil:
// 		dataSource.DataSource = &dvv1alpha1.DataVolumeConfig_DataSource_Upload_{
// 			Upload: &dvv1alpha1.DataVolumeConfig_DataSource_Upload{},
// 		}
// 	case in.Spec.Source.Blank != nil:
// 		dataSource.DataSource = &dvv1alpha1.DataVolumeConfig_DataSource_Blank_{
// 			Blank: &dvv1alpha1.DataVolumeConfig_DataSource_Blank{},
// 		}
// 	}
// 	config.DataSource = &dataSource

// 	if in.Spec.PVC != nil {
// 		config.BoundPvc = &dvv1alpha1.DataVolumeConfig_BoundPVC{
// 			StorageClassName: lo.FromPtr(in.Spec.PVC.StorageClassName),
// 		}
// 		if quantity, ok := in.Spec.PVC.Resources.Requests[corev1.ResourceStorage]; ok {
// 			config.BoundPvc.Capacity = quantity.String()
// 		}
// 	}

// 	pbSpec, err := utils.ConvertToProtoStruct(in.Spec)
// 	if err != nil {
// 		return nil, err
// 	}
// 	pbStatus, err := utils.ConvertToProtoStruct(in.Status)
// 	if err != nil {
// 		return nil, err
// 	}
// 	datavolume := dvv1alpha1.DataVolume{
// 		Namespace:         in.Namespace,
// 		Name:              in.Name,
// 		Labels:            in.Labels,
// 		Annotations:       in.Annotations,
// 		Spec:              pbSpec,
// 		Status:            pbStatus,
// 		CreationTimestamp: timestamppb.New(in.CreationTimestamp.Time),
// 	}
// 	return &datavolume, nil
// }

func generateDataVolumeCRD(namespace, name string, config *dvv1alpha1.DataVolumeConfig) (*cdiv1beta1.DataVolume, error) {
	var disk string
	switch config.Disk {
	case dvv1alpha1.DataVolumeConfig_BOOT:
		disk = bootDisk
	case dvv1alpha1.DataVolumeConfig_DATA:
		disk = dataDisk
	default:
		return nil, fmt.Errorf("unsupported disk type: %v", config.Disk)
	}

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
				label.IoVinkDisk.Name: disk,
			},
			Annotations: map[string]string{
				annotation.IoKubevirtCdiStorageBindImmediateRequested.Name: "true",
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

	if disk == bootDisk {
		var osfamily, osversion string
		switch v := config.OsFamily.OsFamily.(type) {
		case *dvv1alpha1.DataVolumeConfig_OSFamily_Centos_:
			osfamily = centos
			osversion = v.Centos.Version
		case *dvv1alpha1.DataVolumeConfig_OSFamily_Ubuntu_:
			osfamily = ubuntu
			osversion = v.Ubuntu.Version
		case *dvv1alpha1.DataVolumeConfig_OSFamily_Debian_:
			osfamily = debian
			osversion = v.Debian.Version
		case *dvv1alpha1.DataVolumeConfig_OSFamily_Windows_:
			osfamily = windows
			osversion = v.Windows.Version
		}
		dvcrd.Labels[label.IoVinkOsFamily.Name] = osfamily
		dvcrd.Labels[label.IoVinkOsVersion.Name] = osversion
	}

	return &dvcrd, nil
}
