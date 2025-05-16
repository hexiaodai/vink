package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubevirtcorev1 "kubevirt.io/api/core/v1"
)

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:categories=vink,path="templates",scope="Namespaced",shortName="tpl",singular="template"
// +kubebuilder:printcolumn:JSONPath=".spec.general.os.name",description="os",name="OS",type=string
// +kubebuilder:printcolumn:JSONPath=".spec.general.os.version",description="version",name="VERSION",type=string
// +kubebuilder:printcolumn:JSONPath=".spec.compute.cpu.cores",description="cpu",name="CPU",type=number
// +kubebuilder:printcolumn:JSONPath=".spec.compute.memory.size",description="memory",name="MEMORY",type=string
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

type Template struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TemplateSpec   `json:"spec"`
	Status TemplateStatus `json:"status,omitempty"`
}

type TemplateStatus struct {
}

type TemplateSpec struct {
	// +kubebuilder:validation:Required
	General *GeneralSpec `json:"general"`

	// +kubebuilder:validation:Required
	Compute *ComputeSpec `json:"compute"`

	// +kubebuilder:validation:Required
	Storage *StorageSpec `json:"storage"`

	// +kubebuilder:validation:Optional
	Network *NetworkSpec `json:"network"`

	// +kubebuilder:validation:Optional
	Initialization *InitializationSpec `json:"initialization,omitempty"`

	// +kubebuilder:validation:Optional
	Access *AccessSpec `json:"access,omitempty"`

	// +kubebuilder:validation:Optional
	Advanced *AdvancedSpec `json:"advanced,omitempty"`
}

type GeneralSpec struct {
	// +kubebuilder:validation:Optional
	Os *OperatingSystemSpec `json:"os,omitempty"`

	// +kubebuilder:validation:Required
	Source *ImageSource `json:"source"`

	// +kubebuilder:validation:Required
	User *UserSpec `json:"user"`
}

type OperatingSystemSpec struct {
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// +kubebuilder:validation:Required
	Version string `json:"version"`
}

type ImageSource struct {
	// +kubebuilder:validation:Optional
	Builtin *ImageBuiltinSource `json:"builtin,omitempty"`

	// +kubebuilder:validation:Optional
	Http *ImageHTTPSource `json:"http,omitempty"`

	// +kubebuilder:validation:Optional
	S3 *ImageS3Source `json:"s3,omitempty"`

	// +kubebuilder:validation:Optional
	Registry *ImageRegistrySource `json:"registry,omitempty"`

	// +kubebuilder:validation:Optional
	Pvc *ImagePVCSource `json:"pvc,omitempty"`

	// +kubebuilder:validation:Optional
	DataVolume *ImageDataVolumeSource `json:"dataVolume,omitempty"`
}

type ImageBuiltinSource struct {
	// +kubebuilder:validation:Required
	Distribution string `json:"distribution"`

	// +kubebuilder:validation:Required
	Version string `json:"version"`
}

type ImageHTTPSource struct {
	// +kubebuilder:validation:Required
	Url string `json:"url,omitempty"`
}

type ImageS3Source struct {
	// +kubebuilder:validation:Required
	Url string `json:"url"`
}

type ImageRegistrySource struct {
	// +kubebuilder:validation:Required
	Url string `json:"url"`
}

type ImagePVCSource struct {
	// +kubebuilder:validation:Required
	Name string `json:"name"`
}

type ImageDataVolumeSource struct {
	// +kubebuilder:validation:Required
	Name string `json:"name"`
}

type UserSpec struct {
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// +kubebuilder:validation:Optional
	Password string `json:"password,omitempty"`

	// +kubebuilder:validation:Optional
	PasswordBase64 string `json:"passwordBase64,omitempty"`

	// +kubebuilder:validation:Optional
	PasswordSecretRef string `json:"passwordSecretRef,omitempty"`

	// +kubebuilder:validation:Optional
	SshKey string `json:"sshKey,omitempty"`

	// +kubebuilder:validation:Optional
	SshKeyBase64 string `json:"sshKeyBase64,omitempty"`

	// +kubebuilder:validation:Optional
	SshKeySecretRef string `json:"sshKeySecretRef,omitempty"`
}

type ComputeSpec struct {
	// +kubebuilder:validation:Required
	Cpu *CpuSpec `json:"cpu"`

	// +kubebuilder:validation:Required
	Memory *MemorySpec `json:"memory"`
}

type CpuSpec struct {
	// +kubebuilder:validation:Required
	Cores int `json:"cores"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=1
	Threads int `json:"threads,omitempty"`
}

type MemorySpec struct {
	// +kubebuilder:validation:Required
	Size string `json:"size"`
}

type StorageSpec struct {
	// +kubebuilder:validation:Required
	RootDisk *DiskSpec `json:"rootDisk"`

	// +kubebuilder:validation:Optional
	DataDisks []DiskSpec `json:"dataDisks,omitempty"`
}

type DiskSpec struct {
	// +kubebuilder:validation:Required
	// Name string `json:"name"`

	// // +kubebuilder:validation:Required
	Size string `json:"size"`

	// +kubebuilder:validation:Optional
	StorageClass string `json:"storageClass"`

	// // +kubebuilder:validation:Optional
	// MountPath string `json:"mountPath,omitempty"`
}

type NetworkSpec struct {
	// +kubebuilder:validation:Optional
	Interfaces []NetworkInterface `json:"interfaces"`
}

type NetworkInterface struct {
	// +kubebuilder:validation:Required
	Type string `json:"type"`

	// // +kubebuilder:validation:Required
	// Name string `json:"name"`

	// // +kubebuilder:validation:Required
	// NetworkName string `json:"networkName"`

	// +kubebuilder:validation:Required
	Nad string `json:"nad"`

	// +kubebuilder:validation:Required
	Subnet string `json:"subnet"`

	// +kubebuilder:validation:Optional
	IpPool string `json:"ipPool,omitempty"`
}

// +kubebuilder:validation:Optional
type InitializationSpec struct {
	// +kubebuilder:validation:Optional
	CloudInit *CloudInitSpec `json:"cloudInit"`
}

type CloudInitSpec struct {
	// +kubebuilder:validation:Optional
	UserData string `json:"userData,omitempty"`

	// +kubebuilder:validation:Optional
	UserDataBase64 string `json:"userDataBase64,omitempty"`
}

type AccessSpec struct {
	// +kubebuilder:validation:Optional
	Ssh *SshAccessSpec `json:"ssh,omitempty"`

	// +kubebuilder:validation:Optional
	Console *ConsoleAccessSpec `json:"console,omitempty"`
}

type SshAccessSpec struct {
	// +kubebuilder:validation:Optional
	Enabled bool `json:"enabled,omitempty"`

	// +kubebuilder:validation:Optional
	Port int `json:"port,omitempty"`
}

type ConsoleAccessSpec struct {
	// +kubebuilder:validation:Optional
	Vnc bool `json:"vnc,omitempty"`

	// +kubebuilder:validation:Optional
	Serial bool `json:"serial,omitempty"`
}

type AdvancedSpec struct {
	// +kubebuilder:validation:Optional
	RawVMOverrides *kubevirtcorev1.VirtualMachineSpec `json:"rawVMOverrides,omitempty"`
}

//+kubebuilder:object:root=true

type TemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Template `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Template{}, &TemplateList{})
}
