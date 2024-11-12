package business

import (
	"context"

	"github.com/kubevm.io/vink/pkg/clients"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ListEvents(ctx context.Context, clients clients.Clients) ([]string, error) {
	_, err := clients.GetKubeVirtClient().CoreV1().Events("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return nil, nil
}
