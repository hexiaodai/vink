
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
    VirtualMachineInstance
)

func (r ResourceTypes) String() string {
	switch r {
	case 1:
		return "DataVolume"
	case 2:
		return "VirtualMachineInstance"
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

	VinkHost = Instance {
		Name:          "vink.kubevm.io/host",
		Description:   "Specifies the host machine where the virtual machine "+
                        "instance is scheduled to run.",
		FeatureStatus: Alpha,
		Deprecated:    false,
		Resources: []ResourceTypes{
			VirtualMachineInstance,
		},
	}

)

func AllResourceAnnotations() []*Instance {
	return []*Instance {
		&IoKubevirtCdiStorageBindImmediateRequested,
		&VinkDatavolumeOwner,
		&VinkHost,
	}
}

func AllResourceTypes() []string {
	return []string {
		"DataVolume",
		"VirtualMachineInstance",
	}
}
