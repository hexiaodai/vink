package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	kubevirtcorev1 "kubevirt.io/api/core/v1"
	cdiv1beta1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
)

var (
	// EndpointGVR is the name of tthe Endpoint GVR.
	VirtualMachineGVR = schema.GroupVersionResource{
		Group:    GroupVersion.Group,
		Version:  GroupVersion.Version,
		Resource: "virtualmachines",
	}
	// EndpointGVK is the name of tthe Endpoint GVK.
	VirtualMachineGVK = schema.GroupVersionKind{
		Group:   GroupVersion.Group,
		Version: GroupVersion.Version,
		Kind:    "VirtualMachine",
	}

	// EndpointFinalizer is the name of the finalizer for Endpoint.
	VirtualMachineFinalizer = VirtualMachineGVR.GroupResource().String()
)

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:categories=vink,path="virtualmachines",scope=Namespaced

// VirtualMachine is the Schema for the virtualmachines API
type VirtualMachine struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VirtualMachineSpec   `json:"spec,omitempty"`
	Status VirtualMachineStatus `json:"status,omitempty"`
}

type VirtualMachineStatus struct {
	Definition      *kubevirtcorev1.VirtualMachine         `json:"definition,omitempty"`
	Instance        *kubevirtcorev1.VirtualMachineInstance `json:"instance,omitempty"`
	RootDiskVolume  *cdiv1beta1.DataVolume                 `json:"rootDiskVolume,omitempty"`
	DataDiskVolumes []*cdiv1beta1.DataVolume               `json:"dataDiskVolumes,omitempty"`
}

// VirtualMachineSpec defines the desired state of VirtualMachine.
type VirtualMachineSpec struct{}

//+kubebuilder:object:root=true

// VirtualMachineList contains a list of VirtualMachine
type VirtualMachineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VirtualMachine `json:"items"`
}

type DiskVolumeTemplate struct{}

type Networks struct{}

func init() {
	SchemeBuilder.Register(&VirtualMachine{}, &VirtualMachineList{})
}
