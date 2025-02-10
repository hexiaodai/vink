#!/usr/bin/env bash
set -euo pipefail

DEBUG_MODE=${DEBUG_MODE:-false}
WAIT_TIMEOUT=${WAIT_TIMEOUT:-10m}

# Function to build Helm options
build_helm_options() {
    local options="--wait --timeout ${WAIT_TIMEOUT}"
    if [ "${DEBUG_MODE}" = "true" ]; then
        options="${options} --debug"
    fi
    echo "${options}"
}

# Install Multus
kubectl apply -f https://raw.githubusercontent.com/k8snetworkplumbingwg/multus-cni/master/deployments/multus-daemonset.yml

# Install Kube-OVN
echo "Installing Kube-OVN"
helm repo add kubeovn https://kubeovn.github.io/kube-ovn --force-update
helm repo update kubeovn

kubectl label node -l beta.kubernetes.io/os=linux kubernetes.io/os=linux --overwrite
kubectl label node -l node-role.kubernetes.io/control-plane kube-ovn/role=master --overwrite

MASTER_NODE_IPS=$(kubectl get nodes -l node-role.kubernetes.io/control-plane -o=jsonpath='{.items[*].status.addresses[?(@.type=="InternalIP")].address}' | tr ' ' ',')

helm upgrade --install kube-ovn kubeovn/kube-ovn \
    --namespace kube-system \
    --set MASTER_NODES=${MASTER_NODE_IPS} \
    -f examples/kube-ovn/kube-ovn-override-values.yaml \
    $(build_helm_options)

# Install Ceph
helm repo add rook-release https://charts.rook.io/release --force-update
helm repo update rook-release

helm upgrade --install --create-namespace --namespace rook-ceph rook-ceph \
    rook-release/rook-ceph \
    -f examples/ceph/rook-ceph-override-values.yaml \
    $(build_helm_options)

helm upgrade --install --create-namespace --namespace rook-ceph rook-ceph-cluster \
    rook-release/rook-ceph-cluster \
    -f examples/ceph/rook-ceph-cluster-override-values.yaml \
    $(build_helm_options)

# Install Snapshotter
#echo "Installing Snapshotter"
#kubectl kustomize https://github.com/kubernetes-csi/external-snapshotter/client/config/crd | kubectl create -f -
#kubectl -n kube-system kustomize https://github.com/kubernetes-csi/external-snapshotter/deploy/kubernetes/snapshot-controller | kubectl create -f -
#kubectl kustomize https://github.com/kubernetes-csi/external-snapshotter/deploy/kubernetes/csi-snapshotter | kubectl create -f -

# Install Monitor
#echo "Installing Monitor"
#helm repo add prometheus-community https://prometheus-community.github.io/helm-charts --force-update
#helm repo update prometheus-community

# helm upgrade --install --create-namespace --namespace monitoring monitoring \
#     prometheus-community/kube-prometheus-stack \
#     -f examples/monitoring/monitoring-override-values.yaml \
#     $(build_helm_options)

# Install KubeVirt
echo "Installing KubeVirt"
helm upgrade --install --create-namespace --namespace kubevirt kubevirt \
    examples/kubevirt/kubevirt-dev \
    $(build_helm_options)

helm upgrade --install --create-namespace --namespace cdi cdi \
    examples/kubevirt/cdi \
    $(build_helm_options)

# Install Vink
echo "Installing Vink"
helm upgrade --install --create-namespace --namespace vink vink \
    examples/vink/vink \
    $(build_helm_options)

echo "Installation completed 👏"
