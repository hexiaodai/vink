# Vink

[英文文档](./README.en.md)

Vink（Virtual Machines in Kubernetes）是基于 KubeVirt 构建的开源云原生虚拟化平台，为 Kubernetes 提供原生虚拟机全生命周期管理能力。通过深度集成 Kube-OVN 网络插件、Rook-Ceph 分布式存储、External-Snapshotter 快照系统，以及 Prometheus & Grafana 监控体系，构建了轻量化云原生虚拟化平台。

## 功能

- **易于安装:** 您可以直接在 Kubernetes 集群上，使用 Helm 安装它。

- **虚拟机管理:** 轻松创建、删除、编辑、克隆、快照、实时迁移虚拟机。

- **镜像管理:** 提供镜像导入和共享功能。

- **网络管理:** 采用 Kube-OVN 与 Multus CNI，实现 VLAN Underlay/Overlay 混合网络（支持多网络接口）与弹性 IP（EIP）。如果外部网络需要连接到虚拟机，建议使用 VLAN 或未标记的 Underlay 网络；否则推荐使用 Overlay 网络，它能提供更多功能（如 VPC、ACL 和虚拟子网）。

- **存储管理:** 采用 Ceph 分布式存储，支持快照、克隆以及虚拟机的实时扩容和实时迁移。

- **监控:** 采用 Prometheus 和 Grafana 监控系统，提供虚拟机、节点、Ceph 存储和 Kube-OVN 组件的监控数据。

- **Kubernetes 原生集成:** 使用 Kubernetes 原生 API 部署虚拟机。

- **可扩展性与高可用性:** 通过自动扩展和故障转移机制确保虚拟机的可靠性。

## 架构

![Vink](./docs/images/vink.png)

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
