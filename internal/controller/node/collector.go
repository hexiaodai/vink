package node

import (
	"context"
	"fmt"
	"time"

	"github.com/kubevm.io/vink/apis/annotation"
	apitypes "github.com/kubevm.io/vink/apis/types"
	"github.com/kubevm.io/vink/internal/controller/pkg"
	"github.com/kubevm.io/vink/pkg/log"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/wait"
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
			log.Errorf("Failed to update node metrics: %v", err)
		}
	}, interval, ctx.Done())
}

func (c *Collector) update(ctx context.Context) error {
	list := &corev1.NodeList{}
	if err := c.client.List(ctx, list); err != nil {
		return fmt.Errorf("failed to list nodes: %w", err)
	}

	for _, node := range list.Items {
		metrics, err := c.getMetrics(ctx, &node)
		if err != nil {
			log.Errorf("Failed to get metrics for node %s: %v", node.Name, err)
			continue
		}

		if err := pkg.PatchAnnotations(ctx, c.client, &node, annotation.VinkMonitor.Name, metrics); err != nil {
			log.Errorf("Failed to update annotations for node %s: %v", node.Name, err)
			continue
		}
	}
	return nil
}

func (c *Collector) getMetrics(ctx context.Context, node *corev1.Node) (*apitypes.NodeResourceMetrics, error) {
	metrics := &apitypes.NodeResourceMetrics{}

	if cpuUsage, err := c.queryPrometheusForCPUUsage(ctx, node); err != nil {
		return nil, fmt.Errorf("failed to get CPU usage: %v", err)
	} else {
		metrics.CpuUsage = float32(cpuUsage)
	}

	if cpuTotal, err := c.queryPrometheusForCPUTotal(ctx, node); err != nil {
		return nil, fmt.Errorf("failed to get total CPU: %v", err)
	} else {
		metrics.CpuTotal = float32(cpuTotal)
	}

	if memUsage, err := c.queryPrometheusForMemUsage(ctx, node); err != nil {
		return nil, fmt.Errorf("failed to get memory usage: %v", err)
	} else {
		metrics.MemeryUsage = float32(memUsage)
	}

	if memTotal, err := c.queryPrometheusForMemTotal(ctx, node); err != nil {
		return nil, fmt.Errorf("failed to get total memory: %v", err)
	} else {
		metrics.MemeryTotal = float32(memTotal)
	}

	if storageTotal, err := c.queryPrometheusForStorageTotal(ctx, node); err != nil {
		return nil, fmt.Errorf("failed to get total storage: %v", err)
	} else {
		metrics.StorageTotal = float32(storageTotal)
	}

	if storageUsage, err := c.queryPrometheusForStorageUsage(ctx, node); err != nil {
		return nil, fmt.Errorf("failed to get storage usage: %v", err)
	} else {
		metrics.StorageUsage = float32(storageUsage)
	}

	return metrics, nil
}

func (c *Collector) queryPrometheusForCPUUsage(ctx context.Context, node *corev1.Node) (float64, error) {
	return pkg.QueryPrometheus(ctx, fmt.Sprintf(`((instance:node_cpu_utilisation:rate5m{job="node-exporter", nodename="%s"} * instance:node_num_cpu:sum{job="node-exporter", nodename="%s"}) != 0 ) / scalar(sum(instance:node_num_cpu:sum{job="node-exporter", nodename="%s"}))`, node.Name, node.Name, node.Name))
}

func (c *Collector) queryPrometheusForCPUTotal(ctx context.Context, node *corev1.Node) (float64, error) {
	return pkg.QueryPrometheus(ctx, fmt.Sprintf(`count(node_cpu_seconds_total{mode="system", nodename="%s"}) by (node)`, node.Name))
}

func (c *Collector) queryPrometheusForMemUsage(ctx context.Context, node *corev1.Node) (float64, error) {
	return pkg.QueryPrometheus(ctx, fmt.Sprintf(`(instance:node_memory_utilisation:ratio{job="node-exporter", nodename="%s"} / scalar(count(instance:node_memory_utilisation:ratio{job="node-exporter", nodename="%s"}))) != 0`, node.Name, node.Name))
}

func (c *Collector) queryPrometheusForMemTotal(ctx context.Context, node *corev1.Node) (float64, error) {
	return pkg.QueryPrometheus(ctx, fmt.Sprintf(`node_memory_MemTotal_bytes{nodename="%s"} / 1024 / 1024 / 1024`, node.Name))
}

func (c *Collector) queryPrometheusForStorageTotal(ctx context.Context, node *corev1.Node) (float64, error) {
	return pkg.QueryPrometheus(ctx, fmt.Sprintf(`sum by (node) (ceph_osd_stat_bytes * on (pod, namespace) group_left(node) kube_pod_info{node="%s"})`, node.Name))
}

func (c *Collector) queryPrometheusForStorageUsage(ctx context.Context, node *corev1.Node) (float64, error) {
	return pkg.QueryPrometheus(ctx, fmt.Sprintf(`sum by (node) (ceph_osd_stat_bytes_used * on (pod, namespace) group_left(node) kube_pod_info{node="%s"})`, node.Name))
}
