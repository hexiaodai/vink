
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
  DataVolume,Node,VirtualMachine,
}

function resourceTypesToString(type: ResourceTypes): string {
  switch (type) {
    case 1:
      return "DataVolume";
    case 2:
      return "Node";
    case 3:
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
  IoSpidernetIpamIppool: {
    name: "ipam.spidernet.io/ippool",
    description: "Define the ippool for the virtual machine. eg "+
                        "'{`ipv4`:[`ippool-v4`]}'",
    featureStatus: FeatureStatus.Alpha,
    hidden: true,
    deprecated: false,
    resources: [
      ResourceTypes.VirtualMachine,
    ]
  },
  IoSpidernetIpamSubnet: {
    name: "ipam.spidernet.io/subnet",
    description: "Define the subnet for the virtual machine. eg "+
                        "'{`ipv4`:[`subnet-v4`]}'",
    featureStatus: FeatureStatus.Alpha,
    hidden: true,
    deprecated: false,
    resources: [
      ResourceTypes.VirtualMachine,
    ]
  },
  IoCniMultusV1DefaultNetwork: {
    name: "v1.multus-cni.io/default-network",
    description: "Define the default network for the virtual machine. eg "+
                        "'namespace/name'",
    featureStatus: FeatureStatus.Alpha,
    hidden: true,
    deprecated: false,
    resources: [
      ResourceTypes.VirtualMachine,
    ]
  },
  IoVinkNodeNetworkInterface: {
    name: "vink.io/node-network-interface",
    description: "Record information of node network interfaces.",
    featureStatus: FeatureStatus.Alpha,
    hidden: true,
    deprecated: false,
    resources: [
      ResourceTypes.Node,
    ]
  },
};

function allResourceAnnotations(): Instance[] {
  return [
    instances.IoKubevirtCdiStorageBindImmediateRequested,instances.IoSpidernetIpamIppool,instances.IoSpidernetIpamSubnet,instances.IoCniMultusV1DefaultNetwork,instances.IoVinkNodeNetworkInterface,
  ];
}

function allResourceTypes(): string[] {
  return [
    "DataVolume","Node","VirtualMachine",
  ];
}
