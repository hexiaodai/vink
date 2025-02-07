### Install Multus CNI

```bash
kubectl apply -f https://raw.githubusercontent.com/k8snetworkplumbingwg/multus-cni/master/deployments/multus-daemonset.yml
```

### Install Kube-OVN

```bash
helm repo add kubeovn https://kubeovn.github.io/kube-ovn

helm repo update kubeovn

kubectl label node -lbeta.kubernetes.io/os=linux kubernetes.io/os=linux --overwrite

kubectl label node -lnode-role.kubernetes.io/control-plane kube-ovn/role=master --overwrite

MASTER_NODE_IPS=$(kubectl get nodes -l node-role.kubernetes.io/control-plane -o=jsonpath='{.items[*].status.addresses[?(@.type=="InternalIP")].address}' | tr ' ' ',')

helm upgrade --install kube-ovn kubeovn/kube-ovn --namespace kube-system --set MASTER_NODES=${MASTER_NODE_IPS} -f examples/kube-ovn/kube-ovn-override-values.yaml --wait

helm delete kube-ovn --namespace kube-system
```

### Install Monitor

```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts

helm repo update prometheus-community

helm install --create-namespace --namespace monitoring monitoring prometheus-community/kube-prometheus-stack -f examples/monitoring/monitoring-override-values.yaml

helm upgrade --namespace monitoring monitoring prometheus-community/kube-prometheus-stack -f examples/monitoring/monitoring-override-values.yaml

helm delete monitoring --namespace monitoring
```

### Install Ceph

```bash
helm repo add rook-release https://charts.rook.io/release

helm repo update rook-release

helm upgrade --install --create-namespace --namespace rook-ceph rook-ceph rook-release/rook-ceph -f examples/ceph/rook-ceph-override-values.yaml --wait

helm upgrade --install --create-namespace --namespace rook-ceph rook-ceph-cluster rook-release/rook-ceph-cluster -f examples/ceph/rook-ceph-cluster-override-values.yaml --wait

helm delete rook-ceph --namespace rook-ceph
for CRD in $(kubectl get crd -n rook-ceph | awk '/ceph.rook.io/ {print $1}'); do
    kubectl get -n rook-ceph "$CRD" -o name | \
    xargs -I {} kubectl patch -n rook-ceph {} --type merge -p '{"metadata":{"finalizers": []}}'
done

helm delete rook-ceph-cluster --namespace rook-ceph
```

### Install Local Path Storage (For Development and Testing Only)

```bash
helm install --create-namespace --namespace local-path-storage local-path-storage examples/local-path-storage/local-path-storage

helm delete local-path-storage --namespace local-path-storage
```

### Install KubeVirt

```bash
helm install --create-namespace --namespace kubevirt kubevirt examples/kubevirt/kubevirt

helm install --create-namespace --namespace kubevirt kubevirt examples/kubevirt/kubevirt-dev

helm install --create-namespace --namespace cdi cdi examples/kubevirt/cdi

helm delete kubevirt --namespace kubevirt

helm delete cdi --namespace cdi
```

### Install Vink

```bash
helm install --create-namespace --namespace vink vink examples/vink/vink

helm delete vink --namespace vink
```
