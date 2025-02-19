package virtualmachine

import (
	"context"
	"fmt"
	"time"

	"github.com/kubevm.io/vink/apis/annotation"
	apitypes "github.com/kubevm.io/vink/apis/types"
	"github.com/kubevm.io/vink/internal/controller/pkg"
	"github.com/kubevm.io/vink/pkg/log"
	"k8s.io/apimachinery/pkg/util/wait"
	kubevirtv1 "kubevirt.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Collector struct {
	client client.Client
	cache  cache.Cache
}

func NewCollector(client client.Client, cache cache.Cache) (*Collector, error) {
	return &Collector{client: client, cache: cache}, nil
}

func (c *Collector) Collector(ctx context.Context, interval time.Duration) {
	wait.Until(func() {
		if err := c.update(ctx); err != nil {
			log.Errorf("Failed to update virtual machine metrics: %v", err)
		}
	}, interval, ctx.Done())
}

func (c *Collector) update(ctx context.Context) error {
	list := &kubevirtv1.VirtualMachineList{}
	if err := c.client.List(ctx, list); err != nil {
		return fmt.Errorf("failed to list virtual machines: %w", err)
	}

	for _, vm := range list.Items {
		metrics, err := c.getMetrics(ctx, &vm)
		if err != nil {
			log.Errorf("Failed to get metrics for virtual machine %s/%s: %v", vm.Namespace, vm.Name, err)
			continue
		}

		if err := pkg.PatchAnnotations(ctx, c.client, &vm, annotation.VinkMonitor.Name, metrics); err != nil {
			log.Errorf("Failed to update annotations for virtual machine %s/%s: %v", vm.Namespace, vm.Name, err)
			continue
		}
	}
	return nil
}

func (c *Collector) getMetrics(ctx context.Context, vm *kubevirtv1.VirtualMachine) (*apitypes.VirtualMachineResourceMetrics, error) {
	metrics := &apitypes.VirtualMachineResourceMetrics{}

	if cpuUsage, err := c.queryPrometheusForCPUUsage(ctx, vm); err != nil {
		return nil, fmt.Errorf("failed to get CPU usage: %v", err)
	} else {
		metrics.CpuUsage = float32(cpuUsage)
	}

	if memUsage, err := c.queryPrometheusForMemUsage(ctx, vm); err != nil {
		return nil, fmt.Errorf("failed to get memory usage: %v", err)
	} else {
		metrics.MemeryUsage = float32(memUsage)
	}

	return metrics, nil
}

func (c *Collector) queryPrometheusForCPUUsage(ctx context.Context, vm *kubevirtv1.VirtualMachine) (float64, error) {
	return pkg.QueryPrometheus(ctx, fmt.Sprintf(`avg(irate(kubevirt_vmi_cpu_usage_seconds_total{namespace="%s", name="%s"}[1m])) by (name, namespace)`, vm.Namespace, vm.Name))
}

func (c *Collector) queryPrometheusForMemUsage(ctx context.Context, vm *kubevirtv1.VirtualMachine) (float64, error) {
	return pkg.QueryPrometheus(ctx, fmt.Sprintf(`avg((1 - (avg_over_time(kubevirt_vmi_memory_usable_bytes{namespace="%s", name="%s"}[1m]) / avg_over_time(kubevirt_vmi_memory_available_bytes{namespace="%s", name="%s"}[1m])))) by (name, namespace)`, vm.Namespace, vm.Name, vm.Namespace, vm.Name))
}
