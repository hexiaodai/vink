package template_instance

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/kubevm.io/vink/pkg/k8s/apis/vink/v1alpha1"
	"github.com/samber/lo"
	corev1 "k8s.io/api/core/v1"
	resource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubevirtv1 "kubevirt.io/api/core/v1"
	cdiv1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
	"sigs.k8s.io/yaml"
)

const defaultNetworkAnno = "v1.multus-cni.io/default-network"

func (r *Reconciler) buildVirtualMachineFromTemplate(_ context.Context, tpl *v1alpha1.Template, tplInstance *v1alpha1.TemplateInstance) (*kubevirtv1.VirtualMachine, error) {
	var (
		dvTplResult = buildDataVolumeTemplates(tpl)
		volumes     = buildVolumes(dvTplResult)
		disks       = buildDisks(volumes)
		networks    = buildNetworks(tpl.Spec.Network)
		ifaces      = buildInterfaces(networks)

		memQty = resource.MustParse(tpl.Spec.Compute.Memory.Size)
	)

	cloudInitVolume, cloudInitDisk, err := buildCloudInit(tpl)
	if err != nil {
		return nil, err
	}
	volumes = append(volumes, lo.FromPtr(cloudInitVolume))
	disks = append(disks, lo.FromPtr(cloudInitDisk))

	vm := &kubevirtv1.VirtualMachine{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: tplInstance.Namespace,
			Name:      tplInstance.Name,
			Labels:    map[string]string{"app.kubernetes.io/created-by": "vink.kubevm.io"},
			OwnerReferences: []metav1.OwnerReference{{
				APIVersion:         tplInstance.APIVersion,
				Kind:               tplInstance.Kind,
				Name:               tplInstance.Name,
				UID:                tplInstance.UID,
				Controller:         lo.ToPtr(true),
				BlockOwnerDeletion: lo.ToPtr(true),
			}},
		},
		Spec: kubevirtv1.VirtualMachineSpec{
			DataVolumeTemplates: dvTplResult.Templates,
			RunStrategy:         lo.ToPtr(kubevirtv1.RunStrategyAlways),
			Template: &kubevirtv1.VirtualMachineInstanceTemplateSpec{
				Spec: kubevirtv1.VirtualMachineInstanceSpec{
					Domain: kubevirtv1.DomainSpec{
						CPU: &kubevirtv1.CPU{
							Cores:   uint32(tpl.Spec.Compute.Cpu.Cores),
							Threads: uint32(tpl.Spec.Compute.Cpu.Threads),
						},
						Memory: &kubevirtv1.Memory{
							Guest: &memQty,
						},
						Resources: kubevirtv1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceMemory: memQty,
							},
						},
						Devices: kubevirtv1.Devices{
							Disks:      disks,
							Interfaces: ifaces,
						},
					},
					Networks: networks,
					Volumes:  volumes,
				},
			},
		},
	}

	return vm, nil
}

type DataVolumeTemplateResult struct {
	Templates    []kubevirtv1.DataVolumeTemplateSpec
	VolumeNameBy map[string]string
}

func buildDataVolumeTemplates(tpl *v1alpha1.Template) *DataVolumeTemplateResult {
	result := &DataVolumeTemplateResult{
		Templates:    make([]kubevirtv1.DataVolumeTemplateSpec, 0, len(tpl.Spec.Storage.DataDisks)+1),
		VolumeNameBy: make(map[string]string, len(tpl.Spec.Storage.DataDisks)+1),
	}

	rootDvName := fmt.Sprintf("%s-root", tpl.Name)
	result.VolumeNameBy[rootDvName] = "root"
	result.Templates = append(result.Templates, kubevirtv1.DataVolumeTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Name: rootDvName,
			Annotations: map[string]string{
				"cdi.kubevirt.io/storage.bind.immediate.requested": "true",
			},
		},
		Spec: cdiv1.DataVolumeSpec{
			// Storage: &cdiv1.StorageSpec{
			// 	AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
			// 	Resources: corev1.ResourceRequirements{
			// 		Requests: corev1.ResourceList{
			// 			corev1.ResourceStorage: resource.MustParse(tpl.Spec.Storage.RootDisk.Size),
			// 		},
			// 	},
			// 	StorageClassName: lo.ToPtr(tpl.Spec.Storage.RootDisk.StorageClass),
			// },
			PVC: &corev1.PersistentVolumeClaimSpec{
				AccessModes: []corev1.PersistentVolumeAccessMode{
					corev1.ReadWriteOnce,
				},
				Resources: corev1.VolumeResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceStorage: resource.MustParse(tpl.Spec.Storage.RootDisk.Size),
					},
				},
				StorageClassName: lo.ToPtr(tpl.Spec.Storage.RootDisk.StorageClass),
			},
			Source: buildImageSource(tpl.Spec.General.Source),
		},
	})

	for _, disk := range tpl.Spec.Storage.DataDisks {
		diskID := generateDiskID()
		dvName := fmt.Sprintf("%s-%s", tpl.Name, diskID)
		result.VolumeNameBy[dvName] = diskID

		result.Templates = append(result.Templates, kubevirtv1.DataVolumeTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Name: dvName,
				Annotations: map[string]string{
					"cdi.kubevirt.io/storage.bind.immediate.requested": "true",
				},
			},
			Spec: cdiv1.DataVolumeSpec{
				// Storage: &cdiv1.StorageSpec{
				// 	Resources: corev1.ResourceRequirements{
				// 		Requests: corev1.ResourceList{
				// 			corev1.ResourceStorage: resource.MustParse(disk.Size),
				// 		},
				// 	},
				// 	StorageClassName: lo.ToPtr(disk.StorageClass),
				// },
				PVC: &corev1.PersistentVolumeClaimSpec{
					AccessModes: []corev1.PersistentVolumeAccessMode{
						corev1.ReadWriteOnce,
					},
					Resources: corev1.VolumeResourceRequirements{
						Requests: corev1.ResourceList{
							corev1.ResourceStorage: resource.MustParse(disk.Size),
						},
					},
					StorageClassName: lo.ToPtr(disk.StorageClass),
				},
				Source: &cdiv1.DataVolumeSource{
					Blank: &cdiv1.DataVolumeBlankImage{},
				},
			},
		})
	}

	return result
}

func buildImageSource(source *v1alpha1.ImageSource) *cdiv1.DataVolumeSource {
	switch {
	case source.Builtin != nil:
		return nil
	case source.Http != nil:
		return &cdiv1.DataVolumeSource{HTTP: &cdiv1.DataVolumeSourceHTTP{URL: source.Http.Url}}
	case source.S3 != nil:
		return &cdiv1.DataVolumeSource{S3: &cdiv1.DataVolumeSourceS3{URL: source.S3.Url}}
	case source.Registry != nil:
		return &cdiv1.DataVolumeSource{Registry: &cdiv1.DataVolumeSourceRegistry{URL: &source.Registry.Url}}
	case source.Pvc != nil:
		return &cdiv1.DataVolumeSource{PVC: &cdiv1.DataVolumeSourcePVC{Name: source.Pvc.Name}}
	case source.DataVolume != nil:
		return nil
	default:
		return nil
	}
}

func buildVolumes(dvTplResult *DataVolumeTemplateResult) []kubevirtv1.Volume {
	volumes := make([]kubevirtv1.Volume, 0, len(dvTplResult.Templates))
	for _, dvTpl := range dvTplResult.Templates {
		volumes = append(volumes, kubevirtv1.Volume{
			Name: dvTplResult.VolumeNameBy[dvTpl.Name],
			VolumeSource: kubevirtv1.VolumeSource{
				DataVolume: &kubevirtv1.DataVolumeSource{
					Name: dvTpl.Name,
				},
			},
		})
	}

	return volumes
}

func buildDisks(volumes []kubevirtv1.Volume) []kubevirtv1.Disk {
	disks := make([]kubevirtv1.Disk, 0, len(volumes))
	for _, volume := range volumes {
		var bootOrder *uint
		if volume.Name == "root" {
			bootOrder = lo.ToPtr[uint](1)
		}
		disks = append(disks, kubevirtv1.Disk{
			BootOrder: bootOrder,
			Name:      volume.Name,
			DiskDevice: kubevirtv1.DiskDevice{
				Disk: &kubevirtv1.DiskTarget{
					Bus: "virtio",
				},
			},
		})
	}
	return disks
}

func buildCloudInit(tpl *v1alpha1.Template) (*kubevirtv1.Volume, *kubevirtv1.Disk, error) {
	cloudInitVolume := kubevirtv1.Volume{
		Name:         "cloud-init",
		VolumeSource: kubevirtv1.VolumeSource{},
	}
	cloudInitDisk := kubevirtv1.Disk{
		Name: "cloud-init",
		DiskDevice: kubevirtv1.DiskDevice{
			Disk: &kubevirtv1.DiskTarget{
				Bus: "virtio",
			},
		},
	}

	var init = tpl.Spec.Initialization
	if init == nil || init.CloudInit == nil || len(init.CloudInit.UserDataBase64) == 0 || len(init.CloudInit.UserData) == 0 {
		defaultInit, err := generateDefaultCloudInit2(tpl)
		if err != nil {
			return nil, nil, err
		}
		encoded := base64.StdEncoding.EncodeToString([]byte(defaultInit))
		cloudInitVolume.VolumeSource.CloudInitNoCloud = &kubevirtv1.CloudInitNoCloudSource{
			UserDataBase64: encoded,
		}
		return &cloudInitVolume, &cloudInitDisk, nil
	}

	if len(init.CloudInit.UserDataBase64) > 0 {
		cloudInitVolume.VolumeSource.CloudInitNoCloud = &kubevirtv1.CloudInitNoCloudSource{
			UserDataBase64: init.CloudInit.UserDataBase64,
		}
	} else {
		cloudInitVolume.VolumeSource.CloudInitNoCloud = &kubevirtv1.CloudInitNoCloudSource{
			UserData: init.CloudInit.UserData,
		}
	}

	return &cloudInitVolume, &cloudInitDisk, nil
}

func buildChpasswdList(user *v1alpha1.UserSpec) string {
	var b strings.Builder
	if user.Password != "" {
		b.WriteString(fmt.Sprintf("%s:%s\n", user.Name, user.Password))
	}
	return b.String()
}

const defaultCloudInit = `
#cloud-config
ssh_pwauth: true
disable_root: false
chpasswd: {"list": "root:dangerous", expire: False}

runcmd:
- dhclient -r && dhclient
- sed -i "/#\?PermitRootLogin/s/^.*$/PermitRootLogin yes/g" /etc/ssh/sshd_config
- systemctl restart sshd.service
`

func generateDefaultCloudInit2(_ *v1alpha1.Template) (string, error) {
	return defaultCloudInit, nil
}

func generateDefaultCloudInit(tpl *v1alpha1.Template) (string, error) {
	cfg := map[string]any{
		"ssh_pwauth":   tpl.Spec.Access.Ssh.Enabled,
		"disable_root": false,
		// "chpasswd":     {"list": "root:dangerous", expire: False},
		"chpasswd": map[string]any{
			"list":   buildChpasswdList(tpl.Spec.General.User),
			"expire": false,
		},
		"runcmd": []string{
			"dhclient -r && dhclient",
		},
	}
	if tpl.Spec.Access.Ssh.Enabled {
		cfg["runcmd"] = append(cfg["runcmd"].([]string), `sed -i '/^#\?PermitRootLogin/s/.*/PermitRootLogin yes/' /etc/ssh/sshd_config`)
		cfg["runcmd"] = append(cfg["runcmd"].([]string), "systemctl restart sshd")
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return "", err
	}

	return "#cloud-config\n" + string(data), nil
}

func buildNetworks(network *v1alpha1.NetworkSpec) []kubevirtv1.Network {
	networks := make([]kubevirtv1.Network, 0, len(network.Interfaces))
	for idx, iface := range network.Interfaces {
		networks = append(networks, kubevirtv1.Network{
			Name: toDashQualifiedName(iface.Nad),
			NetworkSource: kubevirtv1.NetworkSource{
				Multus: &kubevirtv1.MultusNetwork{
					NetworkName: iface.Nad,
					Default:     idx == 0,
				},
			},
		})
	}
	return networks
}

func buildInterfaces(networks []kubevirtv1.Network) []kubevirtv1.Interface {
	interfaces := make([]kubevirtv1.Interface, 0, len(networks))
	for _, iface := range networks {
		interfaces = append(interfaces, kubevirtv1.Interface{
			Name: iface.Name,
			InterfaceBindingMethod: kubevirtv1.InterfaceBindingMethod{
				Bridge: &kubevirtv1.InterfaceBridge{},
			},
		})
	}
	return interfaces
}

func toDashQualifiedName(s string) string {
	return strings.ReplaceAll(s, "/", "-")
}
