# Vink

[中文文档](./README.zh.md)

Vink (Virtual Machines in Kubernetes) is an open-source cloud-native virtualization platform built on KubeVirt, providing Kubernetes with native full-lifecycle management capabilities for virtual machines. By deeply integrating the Kube-OVN network plugin, Rook-Ceph distributed storage, External-Snapshotter snapshot system, and the Prometheus-Grafana monitoring stack, it constructs a lightweight cloud-native virtualization platform.

## Features

- **[Virtual Machine Management](./docs/vm-management.md):** Supports the creation, deletion, editing, cloning, snapshotting, and live migration of virtual machines.

- **[Image Management](./docs/volume.md#镜像):** Provides image import and sharing capabilities.

- **[Network Management](./docs/network.md):** Adopts a hybrid Underlay/Overlay architecture, supporting VPCs, ACLs, subnets, and virtual subnets.

- **[Storage Management](./docs/volume.md):** Enables disk creation and expansion, with compatibility for Ceph distributed storage.

- **Kubernetes Native Integration:** Enables VM deployment using Kubernetes-native APIs.

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
