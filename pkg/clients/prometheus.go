package clients

import (
	"fmt"

	"github.com/prometheus/client_golang/api"
	prometheusv1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

type Prometheus struct {
	prometheusv1.API
}

func NewPrometheus(addr string) (*Prometheus, error) {
	client, err := api.NewClient(api.Config{
		Address: addr,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Prometheus client: %w", err)
	}

	return &Prometheus{
		prometheusv1.NewAPI(client),
	}, nil
}
