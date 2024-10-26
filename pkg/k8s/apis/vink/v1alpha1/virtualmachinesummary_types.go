package v1alpha1

import (
	kubeovn "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	kubevirtcorev1 "kubevirt.io/api/core/v1"
	cdiv1beta1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
)

var (
	VirtualMachineSummaryGVR = schema.GroupVersionResource{
		Group:    GroupVersion.Group,
		Version:  GroupVersion.Version,
		Resource: "virtualmachinesummarys",
	}
	VirtualMachineSummaryGVK = schema.GroupVersionKind{
		Group:   GroupVersion.Group,
		Version: GroupVersion.Version,
		Kind:    "VirtualMachineSummary",
	}
)

//+genclient
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:categories=vink,path="virtualmachinesummarys",scope=Namespaced

type VirtualMachineSummary struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VirtualMachineSummarySpec   `json:"spec,omitempty"`
	Status VirtualMachineSummaryStatus `json:"status,omitempty"`
}

type VirtualMachineSummaryStatus struct {
	VirtualMachine         *VirtualMachine               `json:"virtualMachine,omitempty"`
	VirtualMachineInstance *VirtualMachineInstance       `json:"virtualMachineInstance,omitempty"`
	DataVolumes            []*DataVolume                 `json:"dataVolumes,omitempty"`
	Network                *VirtualMachineSummaryNetwork `json:"network,omitempty"`
}

func VirtualMachineFromKubeVirt(vm *kubevirtcorev1.VirtualMachine) *VirtualMachine {
	copy := vm.DeepCopy()
	copy.ObjectMeta.ManagedFields = nil
	return &VirtualMachine{
		ObjectMeta: copy.ObjectMeta,
		Spec:       &copy.Spec,
		Status:     &copy.Status,
	}
}

type VirtualMachine struct {
	// +kubebuilder:pruning:PreserveUnknownFields
	ObjectMeta metav1.ObjectMeta                    `json:"metadata,omitempty"`
	Spec       *kubevirtcorev1.VirtualMachineSpec   `json:"spec,omitempty"`
	Status     *kubevirtcorev1.VirtualMachineStatus `json:"status,omitempty"`
}

func VirtualMachineInstanceFromKubeVirt(vmi *kubevirtcorev1.VirtualMachineInstance) *VirtualMachineInstance {
	copy := vmi.DeepCopy()
	copy.ObjectMeta.ManagedFields = nil
	return &VirtualMachineInstance{
		ObjectMeta: copy.ObjectMeta,
		// Spec:       &copy.Spec,
		Status: &copy.Status,
	}
}

type VirtualMachineInstance struct {
	// +kubebuilder:pruning:PreserveUnknownFields
	ObjectMeta metav1.ObjectMeta `json:"metadata,omitempty"`
	// Spec       *kubevirtcorev1.VirtualMachineInstanceSpec   `json:"spec,omitempty"`
	Status *kubevirtcorev1.VirtualMachineInstanceStatus `json:"status,omitempty"`
}

func DataVolumeFromKubeVirt(dv *cdiv1beta1.DataVolume) *DataVolume {
	copy := dv.DeepCopy()
	copy.ObjectMeta.ManagedFields = nil
	return &DataVolume{
		ObjectMeta: copy.ObjectMeta,
		Spec:       &copy.Spec,
		Status:     &copy.Status,
	}
}

type DataVolume struct {
	// +kubebuilder:pruning:PreserveUnknownFields
	ObjectMeta metav1.ObjectMeta            `json:"metadata,omitempty"`
	Spec       *cdiv1beta1.DataVolumeSpec   `json:"spec,omitempty"`
	Status     *cdiv1beta1.DataVolumeStatus `json:"status,omitempty"`
}

func IPFromKubeOVN(ip *kubeovn.IP) *IP {
	copy := ip.DeepCopy()
	copy.ObjectMeta.ManagedFields = nil
	return &IP{
		ObjectMeta: copy.ObjectMeta,
		Spec:       &copy.Spec,
	}
}

type IP struct {
	// +kubebuilder:pruning:PreserveUnknownFields
	ObjectMeta metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec       *kubeovn.IPSpec   `json:"spec,omitempty"`
}

type VirtualMachineSummaryNetwork struct {
	IPs []*IP `json:"ips,omitempty"`
}

type VirtualMachineSummarySpec struct{}

//+kubebuilder:object:root=true

type VirtualMachineSummaryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VirtualMachineSummary `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VirtualMachineSummary{}, &VirtualMachineSummaryList{})
}
