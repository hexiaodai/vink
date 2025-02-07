#!/usr/bin/env bash
set -euo pipefail

DEBUG_MODE=${DEBUG_MODE:-false}

# Function to build Helm options
build_helm_options() {
    local options=""
    if [ "${DEBUG_MODE}" = "true" ]; then
        options="--debug"
    fi
    echo "${options}"
}

# Function to uninstall a Helm Release if it exists
uninstall_release() {
    local release_name=$1
    local namespace=$2
    local helm_options

    helm_options=$(build_helm_options)

    if helm ls --namespace "${namespace}" | grep -q "${release_name}"; then
        echo "Uninstalling ${release_name} from namespace ${namespace}"
        helm uninstall --namespace "${namespace}" "${release_name}" ${helm_options}
    else
        echo "Release ${release_name} does not exist in namespace ${namespace}, skipping uninstall"
    fi
}

# Uninstall Vink
echo "Checking and uninstalling Vink"
uninstall_release vink vink

# Uninstall KubeVirt
echo "Checking and uninstalling KubeVirt"
uninstall_release kubevirt kubevirt
uninstall_release cdi cdi

# Uninstall Monitor
echo "Checking and uninstalling Monitor"
uninstall_release monitoring monitoring

# Uninstall Ceph
echo "Checking and uninstalling Ceph"
uninstall_release rook-ceph rook-ceph
for CRD in $(kubectl get crd -n rook-ceph | awk '/ceph.rook.io/ {print $1}'); do
    kubectl get -n rook-ceph "$CRD" -o name | \
    xargs -I {} kubectl patch -n rook-ceph {} --type merge -p '{"metadata":{"finalizers": []}}'
done
uninstall_release rook-ceph-cluster rook-ceph

# Uninstall Kube-OVN
echo "Checking and uninstalling Kube-OVN"
uninstall_release kube-ovn kube-system

echo "Uninstallation completed 👏"
