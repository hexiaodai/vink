package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:categories=vink,path="templateinstances",scope="Namespaced",shortName="tpli",singular="templateinstance"
// +kubebuilder:printcolumn:name="Template",type=string,JSONPath=".spec.template",description="Referenced template name"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

type TemplateInstance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TemplateInstanceSpec   `json:"spec"`
	Status TemplateInstanceStatus `json:"status,omitempty"`
}

type TemplateInstanceSpec struct {
	// +kubebuilder:validation:Required
	Template string `json:"template"`

	// Parameters to override defaults in the template.
	// +optional
	Parameters map[string]string `json:"parameters,omitempty"`
}

type TemplateInstanceStatus struct {
	// Applied indicates whether the template has been successfully rendered and the VM created.
	// +optional
	Applied bool `json:"applied,omitempty"`

	// Reason provides any error or status message.
	// +optional
	Reason string `json:"reason,omitempty"`

	// VirtualMachineName is the name of the generated VM object.
	// +optional
	VirtualMachineName string `json:"virtualMachineName,omitempty"`
}

//+kubebuilder:object:root=true

type TemplateInstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TemplateInstance `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TemplateInstance{}, &TemplateInstanceList{})
}
