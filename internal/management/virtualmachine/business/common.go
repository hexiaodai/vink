package business

import (
	"context"
	"fmt"
	"sync"

	dvv1alpha1 "github.com/kubevm.io/vink/apis/management/datavolume/v1alpha1"
	"golang.org/x/sync/errgroup"

	"github.com/kubevm.io/vink/apis/label"
	vmv1alpha1 "github.com/kubevm.io/vink/apis/management/virtualmachine/v1alpha1"
	"github.com/kubevm.io/vink/pkg/proto"
	"github.com/kubevm.io/vink/pkg/clients"
	"github.com/kubevm.io/vink/pkg/clients/gvr"
	"github.com/kubevm.io/vink/pkg/utils"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	virtv1 "kubevirt.io/api/core/v1"
	cdiv1beta1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
)

func getVirtualMachineInstance(ctx context.Context, namespace, name string) (*virtv1.VirtualMachineInstance, error) {
	dcli := clients.GetClients().GetDynamicKubeClient()

	unobj, err := dcli.Resource(gvr.From(virtv1.VirtualMachineInstance{})).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return clients.FromUnstructured[virtv1.VirtualMachineInstance](unobj)
}

func getDataVolume(ctx context.Context, namespace, name string) (*cdiv1beta1.DataVolume, error) {
	dcli := clients.GetClients().GetDynamicKubeClient()

	unobj, err := dcli.Resource(gvr.From(cdiv1beta1.DataVolume{})).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return clients.FromUnstructured[cdiv1beta1.DataVolume](unobj)
}

func getVirtualMachineDataVolumes(ctx context.Context, vm *virtv1.VirtualMachine) (root *cdiv1beta1.DataVolume, data []*cdiv1beta1.DataVolume, err error) {
	eg := errgroup.Group{}
	eg.SetLimit(10)
	dvs := make([]*cdiv1beta1.DataVolume, 0, len(vm.Spec.Template.Spec.Volumes))
	lock := &sync.Mutex{}
	for _, volume := range vm.Spec.Template.Spec.Volumes {
		copy := volume
		if volume.DataVolume == nil {
			continue
		}
		eg.Go(func() error {
			result, err := getDataVolume(ctx, vm.Namespace, copy.DataVolume.Name)
			if errors.IsNotFound(err) {
				return nil
			} else if err != nil {
				return err
			}
			lock.Lock()
			dvs = append(dvs, result)
			lock.Unlock()
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, nil, err
	}

	for _, dv := range dvs {
		if proto.DataVolumeTypeEqual(dv.Labels[label.DatavolumeType.Name], dvv1alpha1.DataVolumeType_ROOT) {
			root = dv
		} else {
			data = append(data, dv)
		}
	}
	return
}

func crdToAPIVirtualMachine(ctx context.Context, vm *virtv1.VirtualMachine) (*vmv1alpha1.VirtualMachine, error) {
	eg := errgroup.Group{}
	eg.SetLimit(10)

	var vmi *virtv1.VirtualMachineInstance
	eg.Go(func() error {
		result, err := getVirtualMachineInstance(ctx, vm.Namespace, vm.Name)
		if errors.IsNotFound(err) {
			return nil
		} else if err != nil {
			return err
		}
		vmi = result
		return nil
	})

	datadisk := make([]*cdiv1beta1.DataVolume, 0, len(vm.Spec.Template.Spec.Volumes)-2)
	var rootdisk *cdiv1beta1.DataVolume
	eg.Go(func() error {
		var err error
		rootdisk, datadisk, err = getVirtualMachineDataVolumes(ctx, vm)
		return err
	})
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	var diskinters []interface{}
	for _, disk := range datadisk {
		diskinters = append(diskinters, disk)
	}
	return &vmv1alpha1.VirtualMachine{
		Namespace:              vm.Namespace,
		Name:                   vm.Name,
		CreationTimestamp:      timestamppb.New(vm.CreationTimestamp.Time),
		VirtualMachine:         utils.MustConvertToProtoStruct(vm),
		VirtualMachineInstance: utils.MustConvertToProtoStruct(vmi),
		VirtualMachineDataVolume: &vmv1alpha1.VirtualMachine_DataVolume{
			Root: utils.MustConvertToProtoStruct(rootdisk),
			Data: utils.MustConvertToProtoStructs(diskinters),
		},
	}, nil
}

func newSampleVirtualMachine(namespace, name string, config *vmv1alpha1.VirtualMachineConfig) *virtv1.VirtualMachine {
	return &virtv1.VirtualMachine{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Labels:      map[string]string{},
			Annotations: map[string]string{},
		},
		Spec: virtv1.VirtualMachineSpec{
			RunStrategy: lo.ToPtr(virtv1.RunStrategyAlways),
			Template: &virtv1.VirtualMachineInstanceTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      map[string]string{},
					Annotations: map[string]string{},
				},
				Spec: virtv1.VirtualMachineInstanceSpec{
					Domain: virtv1.DomainSpec{
						CPU: &virtv1.CPU{
							Cores: config.Compute.CpuCores,
						},
						Resources: virtv1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse(config.Compute.Memory),
							},
						},
						Devices: virtv1.Devices{},
					},
				},
			},
		},
	}
}

func setupVirtualMachineDataVolumes(vm *virtv1.VirtualMachine, cfg []*vmv1alpha1.VirtualMachineConfig_Storage_DataVolume) error {
	ctx := context.Background()
	dcli := clients.GetClients().GetDynamicKubeClient()

	datadisks := make([]*cdiv1beta1.DataVolume, 0, len(cfg))
	for _, item := range cfg {
		if len(item.Capacity) > 0 {
			return fmt.Errorf("not implemented")
		}
		unboot, err := dcli.Resource(gvr.From(cdiv1beta1.DataVolume{})).Namespace(item.Ref.Namespace).Get(ctx, item.Ref.Name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		data, err := clients.FromUnstructured[cdiv1beta1.DataVolume](unboot)
		if err != nil {
			return err
		}
		datadisks = append(datadisks, data)
	}

	volumes := make([]virtv1.Volume, 0, len(datadisks))
	disks := make([]virtv1.Disk, 0, len(datadisks))

	for _, datadisk := range datadisks {
		volumes = append(volumes, virtv1.Volume{
			Name: datadisk.Name,
			VolumeSource: virtv1.VolumeSource{
				DataVolume: &virtv1.DataVolumeSource{
					Name: datadisk.Status.ClaimName,
				},
			},
		})
	}

	for _, volume := range volumes {
		disks = append(disks, virtv1.Disk{
			Name: volume.Name,
			DiskDevice: virtv1.DiskDevice{
				Disk: &virtv1.DiskTarget{
					Bus: virtv1.DiskBusVirtio,
				},
			},
		})
	}

	vm.Spec.Template.Spec.Volumes = append(vm.Spec.Template.Spec.Volumes, volumes...)
	vm.Spec.Template.Spec.Domain.Devices.Disks = append(vm.Spec.Template.Spec.Domain.Devices.Disks, disks...)

	return nil
}

func getRootDiskName(vmName string) string {
	return fmt.Sprintf("%s-root", vmName)
}

func setupVirtualMachineRootVolume(vm *virtv1.VirtualMachine, cfg *vmv1alpha1.VirtualMachineConfig_Storage_DataVolume) error {
	ctx := context.Background()
	dcli := clients.GetClients().GetDynamicKubeClient()

	unboot, err := dcli.Resource(gvr.From(cdiv1beta1.DataVolume{})).Namespace(cfg.Ref.Namespace).Get(ctx, cfg.Ref.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	bootobj, err := clients.FromUnstructured[cdiv1beta1.DataVolume](unboot)
	if err != nil {
		return err
	}

	vm.Spec.DataVolumeTemplates = append(vm.Spec.DataVolumeTemplates, virtv1.DataVolumeTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getRootDiskName(vm.Name),
			Namespace: vm.Namespace,
			Labels: map[string]string{
				label.VirtualmachineVersion.Name: bootobj.Labels[label.VirtualmachineVersion.Name],
				label.VirtualmachineOs.Name:      bootobj.Labels[label.VirtualmachineOs.Name],
			},
		},
		Spec: cdiv1beta1.DataVolumeSpec{
			Source: &cdiv1beta1.DataVolumeSource{
				PVC: &cdiv1beta1.DataVolumeSourcePVC{
					Namespace: cfg.Ref.Namespace,
					Name:      cfg.Ref.Name,
				},
			},
			PVC: &corev1.PersistentVolumeClaimSpec{
				AccessModes: []corev1.PersistentVolumeAccessMode{
					corev1.ReadWriteOnce,
				},
				Resources: corev1.VolumeResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceStorage: resource.MustParse(cfg.Capacity),
					},
				},
				StorageClassName: &cfg.StorageClassName,
			},
		},
	})

	vm.Spec.Template.Spec.Volumes = append(vm.Spec.Template.Spec.Volumes, virtv1.Volume{
		Name: proto.DataVolumeTypeFromEnum(dvv1alpha1.DataVolumeType_ROOT),
		VolumeSource: virtv1.VolumeSource{
			DataVolume: &virtv1.DataVolumeSource{
				Name: getRootDiskName(vm.Name),
			},
		},
	})

	vm.Spec.Template.Spec.Domain.Devices.Disks = append(vm.Spec.Template.Spec.Domain.Devices.Disks, virtv1.Disk{
		Name: proto.DataVolumeTypeFromEnum(dvv1alpha1.DataVolumeType_ROOT),
		DiskDevice: virtv1.DiskDevice{
			Disk: &virtv1.DiskTarget{
				Bus: virtv1.DiskBusVirtio,
			},
		},
		BootOrder: lo.ToPtr[uint](1),
	})

	return nil
}

func setupVirtualMachineNetwork(vm *virtv1.VirtualMachine, cfg *vmv1alpha1.VirtualMachineConfig_Network) error {
	return nil
}

func setupVirtualMachineUserConfig(vm *virtv1.VirtualMachine, cfg *vmv1alpha1.VirtualMachineConfig_UserConfig) error {
	if len(cfg.CloudInitBase64) == 0 {
		return nil
	}

	// TODO: ssh keys
	vm.Spec.Template.Spec.Volumes = append(vm.Spec.Template.Spec.Volumes, virtv1.Volume{
		Name: "cloudinit",
		VolumeSource: virtv1.VolumeSource{
			CloudInitNoCloud: &virtv1.CloudInitNoCloudSource{
				UserDataBase64: cfg.CloudInitBase64,
			},
		},
	})

	vm.Spec.Template.Spec.Domain.Devices.Disks = append(vm.Spec.Template.Spec.Domain.Devices.Disks, virtv1.Disk{
		Name: "cloudinit",
		DiskDevice: virtv1.DiskDevice{
			Disk: &virtv1.DiskTarget{
				Bus: virtv1.DiskBusVirtio,
			},
		},
	})

	return nil
}
