
// GENERATED FILE -- DO NOT EDIT

package label

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
    VirtualMachine
)

func (r ResourceTypes) String() string {
	switch r {
	case 1:
		return "DataVolume"
	case 2:
		return "VirtualMachine"
	}
	return "Unknown"
}

// Instance describes a single resource label
type Instance struct {
	// The name of the label.
	Name string

	// Description of the label.
	Description string

	// FeatureStatus of this label.
	FeatureStatus FeatureStatus

	// Hide the existence of this label when outputting usage information.
	Hidden bool

	// Mark this label as deprecated when generating usage information.
	Deprecated bool

	// The types of resources this label applies to.
	Resources []ResourceTypes
}

var (

	IoVinkDisk = Instance {
		Name:          "vink.io/disk",
		Description:   "Defines the boot disk and data disks of the virtual "+
                        "machine, where 'boot' represents the boot disk of the "+
                        "virtual machine, and 'data' represents the data disks of "+
                        "the virtual machine.",
		FeatureStatus: Alpha,
		Hidden:        true,
		Deprecated:    false,
		Resources: []ResourceTypes{
			DataVolume,
		},
	}

	IoVinkOsFamily = Instance {
		Name:          "vink.io/os-family",
		Description:   "Defines the operating system family of the virtual "+
                        "machine, for example, 'windows', 'centos', 'ubuntu', "+
                        "'debian'.",
		FeatureStatus: Alpha,
		Hidden:        true,
		Deprecated:    false,
		Resources: []ResourceTypes{
			DataVolume,
			VirtualMachine,
		},
	}

	IoVinkOsVersion = Instance {
		Name:          "vink.io/os-version",
		Description:   "Defines the operating system version of the virtual "+
                        "machine.",
		FeatureStatus: Alpha,
		Hidden:        true,
		Deprecated:    false,
		Resources: []ResourceTypes{
			DataVolume,
			VirtualMachine,
		},
	}

)

func AllResourceLabels() []*Instance {
	return []*Instance {
		&IoVinkDisk,
		&IoVinkOsFamily,
		&IoVinkOsVersion,
	}
}

func AllResourceTypes() []string {
	return []string {
		"DataVolume",
		"VirtualMachine",
	}
}
