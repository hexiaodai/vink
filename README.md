# Vink

Virtual Machines in Kubernetes

## Quickstart Installation Guide

```bash
helm upgrade --install --create-namespace --namespace vink vink oci://registry-1.docker.io/hejianmin/vink --wait --timeout 1800s --debug
```

## Cleanup

```bash
helm delete --namespace vink vink
```

## Quickstart Usage Guide

[Dashboard](./docs/dashboard.md)

[Host](./docs/host.md)

[Virtual Machine Management](./docs/vm-management.md)

[Volume](./docs/volume.md)

[Network](./docs/network.md)

## RoadMap
