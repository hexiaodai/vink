
// GENERATED FILE -- DO NOT EDIT

package annotation

type FeatureStatus int

const (
	Alpha FeatureStatus = iota
	Beta
	Stable
)

func (s FeatureStatus) String() string {
	switch s {
	case Alpha:
		return "Alpha"
	case Beta:
		return "Beta"
	case Stable:
		return "Stable"
	}
	return "Unknown"
}

type ResourceTypes int

const (
	Unknown ResourceTypes = iota
    DataVolume
    Node
    VirtualMachine
)

func (r ResourceTypes) String() string {
	switch r {
	case 1:
		return "DataVolume"
	case 2:
		return "Node"
	case 3:
		return "VirtualMachine"
	}
	return "Unknown"
}

// Instance describes a single resource annotation
type Instance struct {
	// The name of the annotation.
	Name string

	// Description of the annotation.
	Description string

	// FeatureStatus of this annotation.
	FeatureStatus FeatureStatus

	// Hide the existence of this annotation when outputting usage information.
	Hidden bool

	// Mark this annotation as deprecated when generating usage information.
	Deprecated bool

	// The types of resources this annotation applies to.
	Resources []ResourceTypes
}

var (

	IoKubevirtCdiStorageBindImmediateRequested = Instance {
		Name:          "cdi.kubevirt.io/storage.bind.immediate.requested",
		Description:   "CDI executes binding requests immediately.",
		FeatureStatus: Alpha,
		Hidden:        true,
		Deprecated:    false,
		Resources: []ResourceTypes{
			DataVolume,
		},
	}

	IoSpidernetIpamIppool = Instance {
		Name:          "ipam.spidernet.io/ippool",
		Description:   "Define the ippool for the virtual machine. eg "+
                        "'{`ipv4`:[`ippool-v4`]}'",
		FeatureStatus: Alpha,
		Hidden:        true,
		Deprecated:    false,
		Resources: []ResourceTypes{
			VirtualMachine,
		},
	}

	IoSpidernetIpamSubnet = Instance {
		Name:          "ipam.spidernet.io/subnet",
		Description:   "Define the subnet for the virtual machine. eg "+
                        "'{`ipv4`:[`subnet-v4`]}'",
		FeatureStatus: Alpha,
		Hidden:        true,
		Deprecated:    false,
		Resources: []ResourceTypes{
			VirtualMachine,
		},
	}

	IoCniMultusV1DefaultNetwork = Instance {
		Name:          "v1.multus-cni.io/default-network",
		Description:   "Define the default network for the virtual machine. eg "+
                        "'namespace/name'",
		FeatureStatus: Alpha,
		Hidden:        true,
		Deprecated:    false,
		Resources: []ResourceTypes{
			VirtualMachine,
		},
	}

	IoVinkNodeNetworkInterface = Instance {
		Name:          "vink.io/node-network-interface",
		Description:   "Record information of node network interfaces.",
		FeatureStatus: Alpha,
		Hidden:        true,
		Deprecated:    false,
		Resources: []ResourceTypes{
			Node,
		},
	}

)

func AllResourceAnnotations() []*Instance {
	return []*Instance {
		&IoKubevirtCdiStorageBindImmediateRequested,
		&IoSpidernetIpamIppool,
		&IoSpidernetIpamSubnet,
		&IoCniMultusV1DefaultNetwork,
		&IoVinkNodeNetworkInterface,
	}
}

func AllResourceTypes() []string {
	return []string {
		"DataVolume",
		"Node",
		"VirtualMachine",
	}
}
