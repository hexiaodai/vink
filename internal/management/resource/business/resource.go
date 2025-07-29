package business

import (
	"context"
	"sync"

	"github.com/kubevm.io/vink/pkg/clients"
	"github.com/kubevm.io/vink/pkg/utils"
	"golang.org/x/sync/errgroup"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"

	apiextensions_v1alpha1 "github.com/kubevm.io/vink/apis/apiextensions/v1alpha1"
	"github.com/kubevm.io/vink/apis/types"
)

func List(ctx context.Context, gvr schema.GroupVersionResource, opts *types.ListOptions) ([]*apiextensions_v1alpha1.CustomResourceDefinition, []*types.ObjectMeta, error) {
	cli := clients.GetClients().GetDynamicKubeClient().Resource(gvr)

	items := make([]unstructured.Unstructured, 0)
	lock := sync.Mutex{}
	if len(opts.NamespaceNames) > 0 {
		eg := errgroup.Group{}
		eg.SetLimit(10)
		for _, nn := range opts.NamespaceNames {
			eg.Go(func() error {
				result, err := cli.Namespace(nn.Namespace).Get(ctx, nn.Name, metav1.GetOptions{})
				if err != nil && !errors.IsNotFound(err) {
					return err
				}
				if errors.IsNotFound(err) {
					return nil
				}
				lock.Lock()
				items = append(items, *result)
				lock.Unlock()
				return nil
			})
		}
		if err := eg.Wait(); err != nil {
			return nil, nil, err
		}
	} else {
		res, err := cli.List(ctx, metav1.ListOptions{
			LabelSelector: opts.LabelSelector, FieldSelector: opts.FieldSelector,
			Limit: int64(opts.Limit), Continue: opts.Continue,
		})
		if err != nil {
			return nil, nil, err
		}
		items = res.Items
	}

	crds := make([]*apiextensions_v1alpha1.CustomResourceDefinition, 0, len(items))
	metadatas := make([]*types.ObjectMeta, 0, len(items))
	for _, item := range items {
		crd, err := utils.ConvertUnstructuredToCRD(item)
		if err != nil {
			return nil, nil, err
		}
		crds = append(crds, crd)
		metadatas = append(metadatas, crd.Metadata)
	}

	return crds, metadatas, nil
}
