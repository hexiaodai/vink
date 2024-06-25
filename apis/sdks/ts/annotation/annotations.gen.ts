
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
  
  IoKubevirtCdiStorageBindImmediateRequested: {
    name: "cdi.kubevirt.io/storage.bind.immediate.requested",
    description: "CDI executes binding requests immediately.",
    featureStatus: FeatureStatus.Alpha,
    hidden: true,
    deprecated: false,
    resources: [
      ResourceTypes.DataVolume,
    ]
  },
};

function allResourceAnnotations(): Instance[] {
  return [
    instances.IoKubevirtCdiStorageBindImmediateRequested,
  ];
}

function allResourceTypes(): string[] {
  return [
    "DataVolume",
  ];
}
