package cmd

import (
	netv1 "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"
	kubeovn "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	"github.com/kubevm.io/vink/internal/controller/virtualmachinesummarys"
	"github.com/kubevm.io/vink/pkg/k8s/apis/vink/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	uruntime "k8s.io/apimachinery/pkg/util/runtime"
	virtv1 "kubevirt.io/api/core/v1"
	cdiv1beta1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/metrics/server"
)

var scheme = runtime.NewScheme()

func init() {
	uruntime.Must(kubeovn.AddToScheme(scheme))
	uruntime.Must(v1alpha1.AddToScheme(scheme))
	uruntime.Must(corev1.AddToScheme(scheme))
	uruntime.Must(virtv1.AddToScheme(scheme))
	uruntime.Must(cdiv1beta1.AddToScheme(scheme))
	uruntime.Must(netv1.AddToScheme(scheme))
}

func NewCRDManager() (ctrl.Manager, error) {
	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&zap.Options{
		Development: true,
	})))

	config := ctrl.GetConfigOrDie()
	config.Burst = 200
	config.QPS = 100

	// mgr, err := ctrl.NewManager(clients.GetClients().GetKubeConfig(), ctrl.Options{
	mgr, err := ctrl.NewManager(config, ctrl.Options{
		Scheme:                  scheme,
		LeaderElectionID:        "vink.kubevm.io",
		LeaderElectionNamespace: "vink",
		LeaderElection:          false,
		Metrics: server.Options{
			BindAddress: "0",
		},
	})
	if err != nil {
		return nil, err
	}

	return mgr, nil
}

func Register(mgr ctrl.Manager) error {
	if err := (&virtualmachinesummarys.VirtualMachineReconciler{
		Client: mgr.GetClient(),
		Cache:  mgr.GetCache(),
	}).SetupWithManager(mgr); err != nil {
		return err
	}

	if err := (&virtualmachinesummarys.VirtualMachineInstanceReconciler{
		Client: mgr.GetClient(),
		Cache:  mgr.GetCache(),
	}).SetupWithManager(mgr); err != nil {
		return err
	}

	if err := (&virtualmachinesummarys.NetworkReconciler{
		Client: mgr.GetClient(),
		Cache:  mgr.GetCache(),
	}).SetupWithManager(mgr); err != nil {
		return err
	}

	if err := (&virtualmachinesummarys.HostReconciler{
		Client: mgr.GetClient(),
		Cache:  mgr.GetCache(),
	}).SetupWithManager(mgr); err != nil {
		return err
	}

	if err := (&virtualmachinesummarys.DataVolumeReconciler{
		Client: mgr.GetClient(),
		Cache:  mgr.GetCache(),
	}).SetupWithManager(mgr); err != nil {
		return err
	}

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		return err
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		return err
	}

	return nil
}
