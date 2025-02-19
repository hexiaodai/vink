
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
		Deprecated:    false,
		Resources: []ResourceTypes{
			DataVolume,
		},
	}

	VinkDatavolumeOwner = Instance {
		Name:          "vink.kubevm.io/datavolume.owner",
		Description:   "Indicates that this DataVolume is being used by a "+
                        "specific virtual machine.",
		FeatureStatus: Alpha,
		Deprecated:    false,
		Resources: []ResourceTypes{
			DataVolume,
		},
	}

	VinkDisks = Instance {
		Name:          "vink.kubevm.io/disks",
		Description:   "",
		FeatureStatus: Alpha,
		Deprecated:    false,
		Resources: []ResourceTypes{
			VirtualMachine,
		},
	}

	VinkHost = Instance {
		Name:          "vink.kubevm.io/host",
		Description:   "Specifies the host machine where the virtual machine is "+
                        "scheduled to run.",
		FeatureStatus: Alpha,
		Deprecated:    false,
		Resources: []ResourceTypes{
			VirtualMachine,
		},
	}

	VinkMonitor = Instance {
		Name:          "vink.kubevm.io/monitor",
		Description:   "",
		FeatureStatus: Alpha,
		Deprecated:    false,
		Resources: []ResourceTypes{
			VirtualMachine,
		},
	}

	VinkNetworks = Instance {
		Name:          "vink.kubevm.io/networks",
		Description:   "",
		FeatureStatus: Alpha,
		Deprecated:    false,
		Resources: []ResourceTypes{
			VirtualMachine,
		},
	}

	VinkOperatingSystem = Instance {
		Name:          "vink.kubevm.io/operating-system",
		Description:   "",
		FeatureStatus: Alpha,
		Deprecated:    false,
		Resources: []ResourceTypes{
			VirtualMachine,
		},
	}

	VinkStorage = Instance {
		Name:          "vink.kubevm.io/storage",
		Description:   "",
		FeatureStatus: Alpha,
		Deprecated:    false,
		Resources: []ResourceTypes{
			Node,
		},
	}

)

func AllResourceAnnotations() []*Instance {
	return []*Instance {
		&IoKubevirtCdiStorageBindImmediateRequested,
		&VinkDatavolumeOwner,
		&VinkDisks,
		&VinkHost,
		&VinkMonitor,
		&VinkNetworks,
		&VinkOperatingSystem,
		&VinkStorage,
	}
}

func AllResourceTypes() []string {
	return []string {
		"DataVolume",
		"Node",
		"VirtualMachine",
	}
}
