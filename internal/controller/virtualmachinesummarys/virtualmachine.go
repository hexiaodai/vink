package virtualmachinesummarys

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	netv1 "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"
	kubeovn "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	"github.com/kubevm.io/vink/pkg/k8s/apis/vink/v1alpha1"
	"github.com/kubevm.io/vink/pkg/log"
	corev1 "k8s.io/api/core/v1"
	apierr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubevirtv1 "kubevirt.io/api/core/v1"
	cdiv1beta1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type VirtualMachineReconciler struct {
	Client client.Client
	Cache  cache.Cache
}

func (reconciler *VirtualMachineReconciler) Reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	var vm kubevirtv1.VirtualMachine
	if err := reconciler.Client.Get(ctx, request.NamespacedName, &vm); err != nil {
		if apierr.IsNotFound(err) {
			log.Debug("Resource not found. ignoring since object must be deleted")
			return ctrl.Result{}, nil
		} else {
			return ctrl.Result{}, fmt.Errorf("failed to get VirtualMachine: %w", err)
		}
	}

	var summary v1alpha1.VirtualMachineSummary
	err := reconciler.Client.Get(ctx, request.NamespacedName, &summary)
	if err != nil && !apierr.IsNotFound(err) {
		return ctrl.Result{}, fmt.Errorf("failed to list VirtualMachineSummarys: %w", err)
	}
	if apierr.IsNotFound(err) {
		result, err := generateVirtualMachineSummarySpce(&vm)
		if err != nil {
			return ctrl.Result{}, err
		}
		if err := reconciler.Client.Create(ctx, result); err != nil {
			return ctrl.Result{}, err
		}
		summary = *result
	}

	if err := updateVirtualMachineSummary(ctx, reconciler.Client, &vm, &summary); err != nil {
		return ctrl.Result{}, err
	}

	summary.Status.VirtualMachine = v1alpha1.VirtualMachineFromKubeVirt(&vm)
	if err := reconciler.Client.Status().Update(ctx, &summary); err != nil {
		if apierr.IsConflict(err) {
			log.Debugf("VirtualMachineSummary %s/%s conflict", request.Namespace, request.Name)
			return ctrl.Result{Requeue: true}, nil
		}
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (reconciler *VirtualMachineReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kubevirtv1.VirtualMachine{}).
		Complete(reconciler)
}

func setupSummaryDataVolume(ctx context.Context, cli client.Client, vm *kubevirtv1.VirtualMachine, summary *v1alpha1.VirtualMachineSummary) error {
	dvs := []*v1alpha1.DataVolume{}
	for _, volume := range vm.Spec.Template.Spec.Volumes {
		if volume.DataVolume == nil {
			continue
		}
		var dv cdiv1beta1.DataVolume
		err := cli.Get(ctx, client.ObjectKey{Namespace: vm.Namespace, Name: volume.DataVolume.Name}, &dv)
		if err != nil && !apierr.IsNotFound(err) {
			return fmt.Errorf("failed to get DataVolume: %w", err)
		}
		if apierr.IsNotFound(err) {
			continue
		}
		dvs = append(dvs, v1alpha1.DataVolumeFromKubeVirt(&dv))
	}

	summary.Status.DataVolumes = dvs
	return nil
}

func setupSummaryVirtualMachineInstance(ctx context.Context, cli client.Client, vm *kubevirtv1.VirtualMachine, summary *v1alpha1.VirtualMachineSummary) error {
	var vmi kubevirtv1.VirtualMachineInstance
	err := cli.Get(ctx, client.ObjectKey{Namespace: vm.Namespace, Name: vm.Name}, &vmi)
	if err != nil && !apierr.IsNotFound(err) {
		return fmt.Errorf("failed to get VirtualMachineInstance: %w", err)
	}
	if apierr.IsNotFound(err) {
		summary.Status.VirtualMachineInstance = nil
	} else {
		summary.Status.VirtualMachineInstance = v1alpha1.VirtualMachineInstanceFromKubeVirt(&vmi)
	}

	return nil
}

func setupSummaryNetwork(ctx context.Context, cli client.Client, vm *kubevirtv1.VirtualMachine, summary *v1alpha1.VirtualMachineSummary) error {
	ips := make([]*v1alpha1.IP, 0, len(vm.Spec.Template.Spec.Networks))
	for _, net := range vm.Spec.Template.Spec.Networks {
		var networkName string
		if net.Multus != nil {
			networkName = net.Multus.NetworkName
		} else if net.Pod != nil {
			networkName = vm.Spec.Template.ObjectMeta.Annotations["v1.multus-cni.io/default-network"]
		}
		parts := strings.Split(networkName, "/")
		if len(parts) != 2 {
			log.Warnf("invalid multus network name: %s", networkName)
			continue
		}

		multusCR := netv1.NetworkAttachmentDefinition{}
		err := cli.Get(ctx, client.ObjectKey{Namespace: parts[0], Name: parts[1]}, &multusCR)
		if err != nil && !apierr.IsNotFound(err) {
			return fmt.Errorf("failed to get NetworkAttachmentDefinition: %w", err)
		}
		if apierr.IsNotFound(err) {
			continue
		}

		config := map[string]interface{}{}
		if err := json.Unmarshal([]byte(multusCR.Spec.Config), &config); err != nil {
			log.Warnf("failed to unmarshal NetworkAttachmentDefinition config: %v", err)
			continue
		}
		var provider string
		if cniType, ok := config["type"].(string); ok && cniType == "kube-ovn" {
			if p, ok := config["provider"].(string); ok {
				provider = p
			}
		} else if ipamConfig, ok := config["ipam"].(map[string]interface{}); ok {
			if cniType, ok := ipamConfig["type"].(string); ok && cniType == "kube-ovn" {
				if p, ok := ipamConfig["provider"].(string); ok {
					provider = p
				}
			}
		}
		if len(provider) == 0 {
			log.Warnf("invalid provider for NetworkAttachmentDefinition: %s", net.Multus.NetworkName)
			continue
		}

		ipName := fmt.Sprintf("%s.%s.%s", vm.Name, vm.Namespace, provider)
		var ip kubeovn.IP
		err = cli.Get(ctx, client.ObjectKey{Name: ipName}, &ip)
		if err != nil && !apierr.IsNotFound(err) {
			return fmt.Errorf("failed to get IP: %w", err)
		}
		if apierr.IsNotFound(err) {
			continue
		}
		if ip.Spec.PodName != vm.Name || ip.Spec.Namespace != vm.Namespace || ip.Spec.PodType != "VirtualMachine" {
			continue
		}
		ips = append(ips, v1alpha1.IPFromKubeOVN(&ip))
	}

	summary.Status.Network = &v1alpha1.VirtualMachineSummaryNetwork{IPs: ips}
	return nil
}

func setupSummaryHost(ctx context.Context, cli client.Client, vm *kubevirtv1.VirtualMachine, summary *v1alpha1.VirtualMachineSummary) error {
	var vmi kubevirtv1.VirtualMachineInstance
	err := cli.Get(ctx, client.ObjectKey{Namespace: vm.Namespace, Name: vm.Name}, &vmi)
	if err != nil && !apierr.IsNotFound(err) {
		return fmt.Errorf("failed to get VirtualMachineInstance: %w", err)
	}
	if apierr.IsNotFound(err) {
		summary.Status.Host = nil
		return nil
	}

	var node corev1.Node
	err = cli.Get(ctx, client.ObjectKey{Name: vmi.Status.NodeName}, &node)
	if err != nil && !apierr.IsNotFound(err) {
		return fmt.Errorf("failed to get Node: %w", err)
	}
	if apierr.IsNotFound(err) {
		summary.Status.Host = nil
	} else {
		summary.Status.Host = v1alpha1.NodeFromKube(&node)
	}

	return nil
}

func generateVirtualMachineSummarySpce(vm *kubevirtv1.VirtualMachine) (*v1alpha1.VirtualMachineSummary, error) {
	summary := v1alpha1.VirtualMachineSummary{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: vm.Namespace,
			Name:      vm.Name,
		},
		Spec: v1alpha1.VirtualMachineSummarySpec{},
	}

	return &summary, nil
}

func updateVirtualMachineSummary(ctx context.Context, cli client.Client, vm *kubevirtv1.VirtualMachine, summary *v1alpha1.VirtualMachineSummary) error {
	summary.Status.VirtualMachine = v1alpha1.VirtualMachineFromKubeVirt(vm)

	if err := setupSummaryDataVolume(ctx, cli, vm, summary); err != nil {
		return err
	}

	if err := setupSummaryVirtualMachineInstance(ctx, cli, vm, summary); err != nil {
		return err
	}

	if err := setupSummaryNetwork(ctx, cli, vm, summary); err != nil {
		return err
	}

	if err := setupSummaryHost(ctx, cli, vm, summary); err != nil {
		return err
	}

	return nil
}
