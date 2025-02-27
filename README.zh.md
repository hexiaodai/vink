# Vink

Virtual Machines in Kubernetes

## Overview

- **虚拟机管理：** 支持虚拟机的创建、编辑、克隆、快照、实时迁移与删除。

- **镜像管理：** 提供镜像导入与共享功能。

- **网络管理：** 采用 Underlay/Overlay 混合架构，支持 VPC、ACL、子网及虚拟子网。

- **磁盘管理：** 支持磁盘创建、扩容与删除，并兼容 Ceph 分布式存储。

## Architecture

## Quickstart

### Installation

```bash
helm upgrade --install --create-namespace --namespace vink vink oci://registry-1.docker.io/hejianmin/vink --wait --timeout 1800s --debug
```

### Cleanup

```bash
helm delete --namespace vink vink
```

## RoadMap
