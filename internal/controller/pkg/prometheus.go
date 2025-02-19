package pkg

import (
	"context"
	"fmt"
	"time"

	"github.com/kubevm.io/vink/pkg/clients"
	"github.com/kubevm.io/vink/pkg/log"
	"github.com/prometheus/common/model"
)

func QueryPrometheus(ctx context.Context, query string) (float64, error) {
	result, warnings, err := clients.Clients.Prometheus.Query(ctx, query, time.Now().UTC())
	if err != nil {
		return 0, err
	}

	if len(warnings) > 0 {
		log.Warnf("warnings from Prometheus: %v", warnings)
	}

	if result.Type() != model.ValVector {
		return 0, fmt.Errorf("unexpected result type: %v", result.Type())
	}

	metrics := result.(model.Vector)
	if len(metrics) == 0 {
		return 0, nil
	}

	return float64(metrics[0].Value), nil
}
