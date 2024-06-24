package business

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	spv2beta1 "github.com/spidernet-io/spiderpool/pkg/k8s/apis/spiderpool.spidernet.io/v2beta1"
	"golang.org/x/sync/errgroup"

	"github.com/kubevm.io/vink/apis/annotation"
	"github.com/kubevm.io/vink/apis/label"
	vmv1alpha1 "github.com/kubevm.io/vink/apis/management/virtualmachine/v1alpha1"
	"github.com/kubevm.io/vink/internal/pkg/virtualmachine"
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

func getVirtualMachine(ctx context.Context, namespace, name string) (*virtv1.VirtualMachine, error) {
	dcli := clients.GetClients().GetDynamicKubeClient()

	unobj, err := dcli.Resource(gvr.From(virtv1.VirtualMachine{})).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return clients.FromUnstructured[virtv1.VirtualMachine](unobj)
}

func getVirtualMachineInstance(ctx context.Context, namespace, name string) (*virtv1.VirtualMachineInstance, error) {
	dcli := clients.GetClients().GetDynamicKubeClient()

	unobj, err := dcli.Resource(gvr.From(virtv1.VirtualMachineInstance{})).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return clients.FromUnstructured[virtv1.VirtualMachineInstance](unobj)
}

func getVirtualMachineNetwork(ctx context.Context, namespace, name string) (*spv2beta1.SpiderEndpoint, error) {
	dcli := clients.GetClients().GetDynamicKubeClient()

	unobj, err := dcli.Resource(gvr.From(spv2beta1.SpiderEndpoint{})).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return clients.FromUnstructured[spv2beta1.SpiderEndpoint](unobj)
}

func getDataVolume(ctx context.Context, namespace, name string) (*cdiv1beta1.DataVolume, error) {
	dcli := clients.GetClients().GetDynamicKubeClient()

	unobj, err := dcli.Resource(gvr.From(cdiv1beta1.DataVolume{})).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return clients.FromUnstructured[cdiv1beta1.DataVolume](unobj)
}

func getVirtualMachineDisks(ctx context.Context, vm *virtv1.VirtualMachine) (boot *cdiv1beta1.DataVolume, data []*cdiv1beta1.DataVolume, err error) {
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
		if dv.Labels[label.IoVinkDisk.Name] == "boot" {
			boot = dv
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

	// var net *spv2beta1.SpiderEndpoint
	// eg.Go(func() error {
	// 	result, err := getVirtualMachineNetwork(ctx, vm.Namespace, vm.Name)
	// 	if errors.IsNotFound(err) {
	// 		return nil
	// 	} else if err != nil {
	// 		return err
	// 	}
	// 	net = result
	// 	return nil
	// })

	datadisk := make([]*cdiv1beta1.DataVolume, 0, len(vm.Spec.Template.Spec.Volumes)-2)
	var rootdisk *cdiv1beta1.DataVolume
	eg.Go(func() error {
		var err error
		rootdisk, datadisk, err = getVirtualMachineDisks(ctx, vm)
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
		// VirtualMachineNetwork:  utils.MustConvertToProtoStruct(net),
		VirtualMachineDisk: &vmv1alpha1.VirtualMachine_Disk{
			Root: utils.MustConvertToProtoStruct(rootdisk),
			Data: utils.MustConvertToProtoStructs(diskinters),
		},
	}, nil
}

func crdsToAPIVirtualMachine(vm *virtv1.VirtualMachine, vmi *virtv1.VirtualMachineInstance, net *spv2beta1.SpiderEndpoint, bootdisk *cdiv1beta1.DataVolume, datadisks []*cdiv1beta1.DataVolume) (*vmv1alpha1.VirtualMachine, error) {
	var interfaces []interface{}
	for _, disk := range datadisks {
		interfaces = append(interfaces, disk)
	}
	return &vmv1alpha1.VirtualMachine{
		Namespace:              vm.Namespace,
		Name:                   vm.Name,
		CreationTimestamp:      timestamppb.New(vm.CreationTimestamp.Time),
		VirtualMachine:         utils.MustConvertToProtoStruct(vm),
		VirtualMachineInstance: utils.MustConvertToProtoStruct(vmi),
		// VirtualMachineNetwork:  utils.MustConvertToProtoStruct(net),
		VirtualMachineDisk: &vmv1alpha1.VirtualMachine_Disk{
			// Boot: utils.MustConvertToProtoStruct(bootdisk),
			Data: utils.MustConvertToProtoStructs(interfaces),
		},
	}, nil
}

func crdToAPIVirtualMachineInstance(in *virtv1.VirtualMachineInstance) (*vmv1alpha1.VirtualMachineInstance, error) {
	pbSpec, err := utils.ConvertToProtoStruct(in.Spec)
	if err != nil {
		return nil, err
	}
	pbStatus, err := utils.ConvertToProtoStruct(in.Status)
	if err != nil {
		return nil, err
	}
	return &vmv1alpha1.VirtualMachineInstance{
		Name:      in.Name,
		Namespace: in.Namespace,
		Spec:      pbSpec,
		Status:    pbStatus,
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
							Cores: config.Resources.CpuCores,
						},
						Resources: virtv1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse(config.Resources.Memory),
							},
						},
						Devices: virtv1.Devices{},
					},
				},
			},
		},
	}
}

func setupVirtualMachineDataDisks(vm *virtv1.VirtualMachine, dataDisksCfg []*vmv1alpha1.VirtualMachineConfig_Storage_DataDisk) error {
	ctx := context.Background()
	dcli := clients.GetClients().GetDynamicKubeClient()

	datadisks := make([]*cdiv1beta1.DataVolume, 0, len(dataDisksCfg))
	for _, item := range dataDisksCfg {
		unboot, err := dcli.Resource(gvr.From(cdiv1beta1.DataVolume{})).Namespace(vm.Namespace).Get(ctx, item.DataVolumeRef, metav1.GetOptions{})
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
					Bus: "virtio",
				},
			},
		})
	}

	vm.Spec.Template.Spec.Volumes = append(vm.Spec.Template.Spec.Volumes, volumes...)
	vm.Spec.Template.Spec.Domain.Devices.Disks = append(vm.Spec.Template.Spec.Domain.Devices.Disks, disks...)

	return nil
}

// func setupVirtualMachineDataDisks(vm *virtv1.VirtualMachine, dataDisksCfg []*vmv1alpha1.VirtualMachineConfig_Storage_DataDisk) error {
// 	ctx := context.Background()
// 	dcli := clients.GetClients().GetDynamicKubeClient()

// 	datadisks := make([]*cdiv1beta1.DataVolume, 0, len(dataDisksCfg))
// 	for _, item := range dataDisksCfg {
// 		unboot, err := dcli.Resource(gvr.From(cdiv1beta1.DataVolume{})).Namespace(vm.Namespace).Get(ctx, item.DataVolumeRef, metav1.GetOptions{})
// 		if err != nil {
// 			return err
// 		}
// 		data, err := clients.FromUnstructured[cdiv1beta1.DataVolume](unboot)
// 		if err != nil {
// 			return err
// 		}
// 		datadisks = append(datadisks, data)
// 	}

// 	volumes := make([]virtv1.Volume, 0, len(datadisks))
// 	disks := make([]virtv1.Disk, 0, len(datadisks))

// 	for _, datadisk := range datadisks {
// 		volumes = append(volumes, virtv1.Volume{
// 			Name: datadisk.Name,
// 			VolumeSource: virtv1.VolumeSource{
// 				DataVolume: &virtv1.DataVolumeSource{
// 					Name: datadisk.Status.ClaimName,
// 				},
// 			},
// 		})
// 	}

// 	for _, volume := range volumes {
// 		disks = append(disks, virtv1.Disk{
// 			Name: volume.Name,
// 			DiskDevice: virtv1.DiskDevice{
// 				Disk: &virtv1.DiskTarget{
// 					Bus: "virtio",
// 				},
// 			},
// 		})
// 	}

// 	vm.Spec.Template.Spec.Volumes = append(vm.Spec.Template.Spec.Volumes, volumes...)
// 	vm.Spec.Template.Spec.Domain.Devices.Disks = append(vm.Spec.Template.Spec.Domain.Devices.Disks, disks...)

// 	for _, datadisk := range datadisks {
// 		owner := metav1.OwnerReference{
// 			APIVersion:         vm.APIVersion,
// 			Kind:               vm.Kind,
// 			Name:               vm.Name,
// 			UID:                vm.UID,
// 			BlockOwnerDeletion: lo.ToPtr(true),
// 		}
// 		datadisk.SetOwnerReferences([]metav1.OwnerReference{owner})

// 		un, _ := clients.Unstructured(datadisk)
// 		if _, err := dcli.Resource(gvr.From(cdiv1beta1.DataVolume{})).Update(ctx, un, metav1.UpdateOptions{}); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

func getBootDiskName(vmName string) string {
	return fmt.Sprintf("%s-boot", vmName)
}

func setupVirtualMachineBootDisk(vm *virtv1.VirtualMachine, cfg *vmv1alpha1.VirtualMachineConfig_Storage_BootDisk) error {
	ctx := context.Background()
	dcli := clients.GetClients().GetDynamicKubeClient()

	unboot, err := dcli.Resource(gvr.From(cdiv1beta1.DataVolume{})).Namespace(cfg.DataVolumeRef.Namespace).Get(ctx, cfg.DataVolumeRef.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	bootobj, err := clients.FromUnstructured[cdiv1beta1.DataVolume](unboot)
	if err != nil {
		return err
	}

	vm.Spec.DataVolumeTemplates = append(vm.Spec.DataVolumeTemplates, virtv1.DataVolumeTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getBootDiskName(vm.Name),
			Namespace: vm.Namespace,
			Labels: map[string]string{
				label.IoVinkOsVersion.Name: bootobj.Labels[label.IoVinkOsVersion.Name],
				label.IoVinkOsFamily.Name:  bootobj.Labels[label.IoVinkOsFamily.Name],
			},
		},
		Spec: cdiv1beta1.DataVolumeSpec{
			Source: &cdiv1beta1.DataVolumeSource{
				PVC: &cdiv1beta1.DataVolumeSourcePVC{
					Namespace: cfg.DataVolumeRef.Namespace,
					Name:      cfg.DataVolumeRef.Name,
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
		Name: "boot",
		VolumeSource: virtv1.VolumeSource{
			DataVolume: &virtv1.DataVolumeSource{
				Name: getBootDiskName(vm.Name),
			},
		},
	})

	vm.Spec.Template.Spec.Domain.Devices.Disks = append(vm.Spec.Template.Spec.Domain.Devices.Disks, virtv1.Disk{
		Name: "boot",
		DiskDevice: virtv1.DiskDevice{
			Disk: &virtv1.DiskTarget{
				Bus: "virtio",
			},
		},
		BootOrder: lo.ToPtr[uint](1),
	})

	// vm.Labels[label.IoVinkOsVersion.Name] = bootobj.Labels[label.IoVinkOsVersion.Name]
	// vm.Labels[label.IoVinkOsFamily.Name] = bootobj.Labels[label.IoVinkOsFamily.Name]

	return nil
}

func setupVirtualMachineNetwork(vm *virtv1.VirtualMachine, cfg *vmv1alpha1.VirtualMachineConfig_Network) error {
	if cfg == nil {
		return nil
	}

	ippoolCfg := virtualmachine.SubnetConfiguration{
		Interface: "eth0",
		IPv4:      []string{cfg.IppoolRef},
	}
	ippoolStr, err := json.Marshal(ippoolCfg)
	if err != nil {
		return err
	}

	vm.Spec.Template.ObjectMeta.Annotations[annotation.IoCniMultusV1DefaultNetwork.Name] = fmt.Sprintf("%v/%v", metav1.NamespaceSystem, cfg.MultusConfigRef)
	vm.Spec.Template.ObjectMeta.Annotations[annotation.IoSpidernetIpamIppool.Name] = string(ippoolStr)

	vm.Spec.Template.Spec.Networks = []virtv1.Network{lo.FromPtr(virtv1.DefaultPodNetwork())}
	vm.Spec.Template.Spec.Domain.Devices.Interfaces = []virtv1.Interface{lo.FromPtr(virtv1.DefaultMasqueradeNetworkInterface())}
	// vm.Spec.Template.Spec.Domain.Devices.Interfaces = []virtv1.Interface{{
	// 	Name: "default",
	// 	InterfaceBindingMethod: virtv1.InterfaceBindingMethod{
	// 		Passt: &virtv1.InterfacePasst{},
	// 	},
	// }}
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
				// UserDataBase64: base64.StdEncoding.EncodeToString([]byte(cfg.CloudInit)),
				UserDataBase64: cfg.CloudInitBase64,
			},
		},
	})

	vm.Spec.Template.Spec.Domain.Devices.Disks = append(vm.Spec.Template.Spec.Domain.Devices.Disks, virtv1.Disk{
		Name: "cloudinit",
		DiskDevice: virtv1.DiskDevice{
			Disk: &virtv1.DiskTarget{
				Bus: "virtio",
			},
		},
	})

	return nil
}
