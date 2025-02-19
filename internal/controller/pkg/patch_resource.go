package pkg

import (
	"context"
	"encoding/json"
	"fmt"

	"k8s.io/apimachinery/pkg/types"

	"github.com/kubevm.io/vink/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func PatchAnnotations(ctx context.Context, cli client.Client, cr client.Object, anno string, metrics interface{}) error {
	metricsBytes, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("failed to marshal metrics for %s %s: %w", cr.GetObjectKind(), cr.GetName(), err)
	}

	patch := map[string]interface{}{
		"metadata": map[string]interface{}{
			"annotations": map[string]string{
				anno: string(metricsBytes),
			},
		},
	}

	patchBytes, err := json.Marshal(patch)
	if err != nil {
		return fmt.Errorf("failed to marshal patch for %s %s: %w", cr.GetObjectKind(), cr.GetName(), err)
	}

	if err := cli.Patch(ctx, cr, client.RawPatch(types.MergePatchType, patchBytes)); err != nil {
		return fmt.Errorf("failed to patch annotations for %s %s: %w", cr.GetObjectKind(), cr.GetName(), err)
	}

	log.Debugf("Successfully patched annotation %q for %s %s", anno, cr.GetObjectKind(), cr.GetName())
	return nil
}
