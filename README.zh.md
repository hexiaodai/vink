# Vink

[英文文档](./README.md)

Vink（Virtual Machines in Kubernetes）是基于 KubeVirt 构建的开源云原生虚拟化平台，为 Kubernetes 提供原生虚拟机全生命周期管理能力。通过深度集成 Kube-OVN 网络插件、Rook-Ceph 分布式存储、External-Snapshotter 快照系统，以及 Prometheus-Grafana 监控体系，构建了轻量化云原生虚拟化平台。

## 功能

- **[虚拟机管理](./docs/vm-management.md):** 支持虚拟机的创建、删除、编辑、克隆、快照、实时迁移。

- **[镜像管理](./docs/volume.md#镜像):** 提供镜像导入和共享功能。

- **[网络管理](./docs/network.md):** 采用混合 Underlay/Overlay 架构，支持 VPC、ACL、子网和虚拟子网。

- **[存储管理](./docs/volume.md):** 支持磁盘的创建和扩展，兼容 Ceph 分布式存储。

- **Kubernetes 原生集成:** 支持使用 Kubernetes 原生 API 部署虚拟机。

- **可扩展性与高可用性:** 通过自动扩展和故障转移机制确保可靠性。

## 架构

## 快速开始

[快速开始指南](./docs/index.md)

### 安装

使用 Helm 安装 Vink，运行以下命令：

```bash
helm upgrade --install --create-namespace --namespace vink vink oci://registry-1.docker.io/hejianmin/vink --wait --timeout 1800s --debug
```

### 清理

要卸载 Vink 并删除所有相关资源，请使用：

```bash
helm delete --namespace vink vink
```

<!-- # Vink

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

## RoadMap -->
