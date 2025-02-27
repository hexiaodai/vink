# Vink

Vink (Virtual Machines in Kubernetes) is a solution for running and managing virtual machines (VMs) natively within Kubernetes. It integrates VMs seamlessly into the Kubernetes ecosystem, allowing users to leverage Kubernetes' orchestration capabilities while maintaining the flexibility of traditional virtualization.

## Features

- **[Virtual Machine Management](./docs/vm-management.md):** Supports VM creation, editing, cloning, snapshots, live migration, and deletion.

- **[Image Management](./docs/volume.md#镜像):** Provides image import and sharing capabilities.

- **[Network Management](./docs/network.md):** Adopts a hybrid Underlay/Overlay architecture, supporting VPCs, ACLs, subnets, and virtual subnets.

- **[Storage Management](./docs/volume.md):** Supports disk creation, expansion, and deletion, with compatibility for Ceph distributed storage.

- **Kubernetes Native Integration:** Enables VM deployment using Kubernetes-native APIs.

- **Multi-Tenant Support:** Provides isolation and resource quotas for different users or workloads.

- **Scalability & High Availability:** Ensures reliability with automated scaling and failover mechanisms.

## Architecture

## Quickstart

[Quickstart Guide](./docs/index.md)

### Installation

To install Vink using Helm, run the following command:

```bash
helm upgrade --install --create-namespace --namespace vink vink oci://registry-1.docker.io/hejianmin/vink --wait --timeout 1800s --debug
```

### Cleanup

To uninstall Vink and remove all associated resources, use:

```bash
helm delete --namespace vink vink
```

## RoadMap
