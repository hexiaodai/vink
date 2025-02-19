package node

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/kubevm.io/vink/apis/annotation"
	apitypes "github.com/kubevm.io/vink/apis/types"
	"github.com/kubevm.io/vink/internal/controller/pkg"
	"github.com/kubevm.io/vink/pkg/clients"
	"github.com/kubevm.io/vink/pkg/log"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type CephStorageCollector struct {
	client client.Client
	cache  cache.Cache
}

func NewCephStorageCollector(client client.Client, cache cache.Cache) (*CephStorageCollector, error) {
	return &CephStorageCollector{
		client: client,
		cache:  cache,
	}, nil
}

func (cs *CephStorageCollector) Collector(ctx context.Context, interval time.Duration) {
	wait.Until(func() {
		if err := cs.update(ctx); err != nil {
			log.Errorf("Failed to update node storage: %v", err)
		}
	}, interval, ctx.Done())
}

func (cs *CephStorageCollector) update(ctx context.Context) error {
	osds, err := clients.Clients.Ceph.ListOsds(ctx)
	if err != nil {
		return fmt.Errorf("failed to list osds: %w", err)
	}

	osdMap := make(map[string][]*apitypes.NodeCephStorage)
	for _, osd := range osds {
		info := apitypes.NodeCephStorage{
			Osd:                  int32(osd.OsdMap.Osd),
			Up:                   false,
			BluestoreBdevDevNode: osd.OsdMetadata.BluestoreBdevDevNode,
			BluestoreBdevType:    osd.OsdMetadata.BluestoreBdevType,
		}
		if osd.OsdMap.Up == 1 {
			info.Up = true
		}
		if total, err := cs.queryPrometheusForStorageTotal(ctx, osd.OsdMetadata.Hostname, osd); err != nil {
			log.Errorf("failed to get total storage for OSD %d: %v", osd.OsdMap.Osd, err)
		} else {
			info.StorageTotal = float32(total)
		}
		if usage, err := cs.queryPrometheusForStorageUsage(ctx, osd.OsdMetadata.Hostname, osd); err != nil {
			log.Errorf("failed to get usage storage for OSD %d: %v", osd.OsdMap.Osd, err)
		} else {
			info.StorageUsage = float32(usage)
		}
		osdMap[osd.OsdMetadata.Hostname] = append(osdMap[osd.OsdMetadata.Hostname], &info)
	}

	for nodeName, info := range osdMap {
		node := corev1.Node{}
		if err := cs.client.Get(ctx, types.NamespacedName{Name: nodeName}, &node); err != nil {
			log.Errorf("Failed to get node: %v", err)
			continue
		}
		if err := pkg.PatchAnnotations(ctx, cs.client, &node, annotation.VinkStorage.Name, info); err != nil {
			log.Errorf("Failed to update node annotations: %v", err)
			continue
		}
	}
	return nil
}

func (cs *CephStorageCollector) queryPrometheusForStorageTotal(ctx context.Context, nodeName string, osd *clients.Osd) (float64, error) {
	return pkg.QueryPrometheus(ctx, fmt.Sprintf(`ceph_osd_stat_bytes{ceph_daemon="osd.%s"} * on (pod, namespace) group_left(node) kube_pod_info{node="%s"}`, strconv.Itoa(osd.OsdMap.Osd), nodeName))
}

func (cs *CephStorageCollector) queryPrometheusForStorageUsage(ctx context.Context, nodeName string, osd *clients.Osd) (float64, error) {
	return pkg.QueryPrometheus(ctx, fmt.Sprintf(`ceph_osd_stat_bytes_used{ceph_daemon="osd.%v"} * on (pod, namespace) group_left(node) kube_pod_info{node="%s"}`, strconv.Itoa(osd.OsdMap.Osd), nodeName))
}
