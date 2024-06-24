
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
  DataVolume,VirtualMachine,
}

function resourceTypesToString(type: ResourceTypes): string {
  switch (type) {
    case 1:
      return "DataVolume";
    case 2:
      return "VirtualMachine";
    
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
  
  IoVinkDisk: {
    name: "vink.io/disk",
    description: "Defines the boot disk and data disks of the virtual "+
                        "machine, where 'boot' represents the boot disk of the "+
                        "virtual machine, and 'data' represents the data disks of "+
                        "the virtual machine.",
    featureStatus: FeatureStatus.Alpha,
    hidden: true,
    deprecated: false,
    resources: [
      ResourceTypes.DataVolume,
    ]
  },
  IoVinkOsFamily: {
    name: "vink.io/os-family",
    description: "Defines the operating system family of the virtual "+
                        "machine, for example, 'windows', 'centos', 'ubuntu', "+
                        "'debian'.",
    featureStatus: FeatureStatus.Alpha,
    hidden: true,
    deprecated: false,
    resources: [
      ResourceTypes.DataVolume,ResourceTypes.VirtualMachine,
    ]
  },
  IoVinkOsVersion: {
    name: "vink.io/os-version",
    description: "Defines the operating system version of the virtual "+
                        "machine.",
    featureStatus: FeatureStatus.Alpha,
    hidden: true,
    deprecated: false,
    resources: [
      ResourceTypes.DataVolume,ResourceTypes.VirtualMachine,
    ]
  },
};

function allResourceLabels(): Instance[] {
  return [
    instances.IoVinkDisk,instances.IoVinkOsFamily,instances.IoVinkOsVersion,
  ];
}

function allResourceTypes(): string[] {
  return [
    "DataVolume","VirtualMachine",
  ];
}
