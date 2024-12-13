
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
)

func (r ResourceTypes) String() string {
	switch r {
	case 1:
		return "DataVolume"
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

	// Mark this label as deprecated when generating usage information.
	Deprecated bool

	// The types of resources this label applies to.
	Resources []ResourceTypes
}

var (

	VinkDatavolumeType = Instance {
		Name:          "vink.kubevm.io/datavolume.type",
		Description:   "Specifies the type of data volume associated with the "+
                        "virtual machine.",
		FeatureStatus: Alpha,
		Deprecated:    false,
		Resources: []ResourceTypes{
			DataVolume,
		},
	}

	VinkOperatingSystem = Instance {
		Name:          "vink.kubevm.io/operating-system",
		Description:   "Defines the operating system of the virtual machine, "+
                        "where 'windows' represents the Windows operating system, "+
                        "and 'linux' represents the Linux operating system.",
		FeatureStatus: Alpha,
		Deprecated:    false,
		Resources: []ResourceTypes{
			DataVolume,
		},
	}

	VinkOperatingSystemVersion = Instance {
		Name:          "vink.kubevm.io/operating-system.version",
		Description:   "Defines the operating system version of the virtual "+
                        "machine.",
		FeatureStatus: Alpha,
		Deprecated:    false,
		Resources: []ResourceTypes{
			DataVolume,
		},
	}

)

func AllResourceLabels() []*Instance {
	return []*Instance {
		&VinkDatavolumeType,
		&VinkOperatingSystem,
		&VinkOperatingSystemVersion,
	}
}

func AllResourceTypes() []string {
	return []string {
		"DataVolume",
	}
}
