#!/bin/bash

set -euo pipefail
set -x

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"


if helm status cdi --namespace cdi > /dev/null 2>&1; then
    helm delete cdi --namespace cdi --debug
fi

if helm status kubevirt --namespace kubevirt > /dev/null 2>&1; then
    helm delete kubevirt --namespace kubevirt --debug
fi

if helm status monitoring --namespace monitoring > /dev/null 2>&1; then
    helm delete monitoring --namespace monitoring --debug
fi

if helm status rook-ceph --namespace rook-ceph > /dev/null 2>&1; then
    helm delete rook-ceph --namespace rook-ceph --debug
fi
for CRD in $(kubectl get crd -n rook-ceph | awk '/ceph.rook.io/ {print $1}'); do
    kubectl get -n rook-ceph "$CRD" -o name | \
    xargs -I {} kubectl patch -n rook-ceph {} --type merge -p '{"metadata":{"finalizers": []}}'
done
if helm status rook-ceph-cluster --namespace rook-ceph > /dev/null 2>&1; then
    helm delete rook-ceph-cluster --namespace rook-ceph --debug
fi

kubectl delete -f https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/master/client/config/crd/snapshot.storage.k8s.io_volumesnapshotclasses.yaml --ignore-not-found=true
kubectl delete -f https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/master/client/config/crd/snapshot.storage.k8s.io_volumesnapshotcontents.yaml --ignore-not-found=true
kubectl delete -f https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/master/client/config/crd/snapshot.storage.k8s.io_volumesnapshots.yaml --ignore-not-found=true

kubectl delete -f https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/master/deploy/kubernetes/snapshot-controller/rbac-snapshot-controller.yaml --ignore-not-found=true
kubectl delete -f https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/master/deploy/kubernetes/snapshot-controller/setup-snapshot-controller.yaml --ignore-not-found=true

kubectl delete -f https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/master/deploy/kubernetes/csi-snapshotter/rbac-csi-snapshotter.yaml --ignore-not-found=true
kubectl delete -f https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/master/deploy/kubernetes/csi-snapshotter/rbac-external-provisioner.yaml --ignore-not-found=true
kubectl delete -f https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/master/deploy/kubernetes/csi-snapshotter/setup-csi-snapshotter.yaml --ignore-not-found=true

curl -sSL https://raw.githubusercontent.com/kubeovn/kube-ovn/release-1.13/dist/images/cleanup.sh -o ${DIR}/kubeovn-cleanup.sh
bash ${DIR}/kubeovn-cleanup.sh || true
if helm status kube-ovn --namespace kube-system > /dev/null 2>&1; then
    helm delete kube-ovn --namespace kube-system --debug
fi
bash ${DIR}/kubeovn-cleanup.sh || true
rm -f ${DIR}/kubeovn-cleanup.sh
