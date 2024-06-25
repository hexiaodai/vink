
// GENERATED FILE -- DO NOT EDIT

enum FeatureStatus {
  Alpha,
  Beta,
  Stable
}

function featureStatusToString(status: FeatureStatus): string {
  switch (status) {
    case FeatureStatus.Alpha:
      return "Alpha";
    case FeatureStatus.Beta:
      return "Beta";
    case FeatureStatus.Stable:
      return "Stable";
    default:
      return "Unknown";
  }
}

enum ResourceTypes {
  Unknown,
  DataVolume,
}

function resourceTypesToString(type: ResourceTypes): string {
  switch (type) {
    case 1:
      return "DataVolume";
    
    default:
      return "Unknown";
  }
}

interface Instance {
  name: string;
  description: string;
  featureStatus: FeatureStatus;
  hidden: boolean;
  deprecated: boolean;
  resources: ResourceTypes[];
}

const instances: { [key: string]: Instance } = {
  
  DatavolumeType: {
    name: "datavolume.vink.io/type",
    description: "Defines the type of datavolume, such as root for the "+
                        "system datavolume, image for the system image, and data "+
                        "for the data datavolume.",
    featureStatus: FeatureStatus.Alpha,
    hidden: true,
    deprecated: false,
    resources: [
      ResourceTypes.DataVolume,
    ]
  },
  VirtualmachineOs: {
    name: "virtualmachine.vink.io/os",
    description: "Defines the operating system of the virtual machine, "+
                        "where 'windows' represents the Windows operating system, "+
                        "and 'linux' represents the Linux operating system.",
    featureStatus: FeatureStatus.Alpha,
    hidden: true,
    deprecated: false,
    resources: [
      ResourceTypes.DataVolume,
    ]
  },
  VirtualmachineVersion: {
    name: "virtualmachine.vink.io/version",
    description: "Defines the operating system version of the virtual "+
                        "machine.",
    featureStatus: FeatureStatus.Alpha,
    hidden: true,
    deprecated: false,
    resources: [
      ResourceTypes.DataVolume,
    ]
  },
};

function allResourceLabels(): Instance[] {
  return [
    instances.DatavolumeType,instances.VirtualmachineOs,instances.VirtualmachineVersion,
  ];
}

function allResourceTypes(): string[] {
  return [
    "DataVolume",
  ];
}
