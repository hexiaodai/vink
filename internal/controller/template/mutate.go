package template

import (
	"context"
	"fmt"

	"github.com/kubevm.io/vink/pkg/k8s/apis/vink/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

const (
	defaultStorageClass = "ceph-block"

	defaultNetworkType = "bridge"

	defaultNetworkNad = "vink/vink"

	defaultNetworkSubnet = "vink"
)

func (webhook *Webhook) mutate(ctx context.Context, obj runtime.Object) error {
	tpl, ok := obj.(*v1alpha1.Template)
	if !ok {
		return fmt.Errorf("object is not a template")
	}

	var allErrs field.ErrorList

	if err := webhook.mutateOs(ctx, tpl); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhook.mutateStorage(ctx, tpl); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhook.mutateNetwork(ctx, tpl); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhook.mutateAccess(ctx, tpl); err != nil {
		allErrs = append(allErrs, err)
	}

	if len(allErrs) == 0 {
		return nil
	}

	return errors.NewInvalid(schema.GroupKind{Group: v1alpha1.GroupVersion.Group, Kind: "Template"}, tpl.Name, allErrs)
}

func (webhook *Webhook) mutateOs(_ context.Context, tpl *v1alpha1.Template) *field.Error {
	if tpl.Spec.General.Source.Builtin == nil {
		return nil
	}

	if tpl.Spec.General.Os == nil {
		tpl.Spec.General.Os = &v1alpha1.OperatingSystemSpec{}
	}
	tpl.Spec.General.Os.Name = tpl.Spec.General.Source.Builtin.Distribution
	tpl.Spec.General.Os.Version = tpl.Spec.General.Source.Builtin.Version

	return nil
}

func (webhook *Webhook) mutateStorage(_ context.Context, tpl *v1alpha1.Template) *field.Error {
	if len(tpl.Spec.Storage.RootDisk.StorageClass) == 0 {
		tpl.Spec.Storage.RootDisk.StorageClass = defaultStorageClass
	}

	for idx, disk := range tpl.Spec.Storage.DataDisks {
		if len(disk.StorageClass) == 0 {
			tpl.Spec.Storage.DataDisks[idx].StorageClass = defaultStorageClass
		}
	}

	return nil
}

func (webhook *Webhook) mutateNetwork(_ context.Context, tpl *v1alpha1.Template) *field.Error {
	if tpl.Spec.Network == nil {
		tpl.Spec.Network = &v1alpha1.NetworkSpec{}
	}

	if len(tpl.Spec.Network.Interfaces) == 0 {
		tpl.Spec.Network.Interfaces = []v1alpha1.NetworkInterface{
			{
				Nad:    defaultNetworkNad,
				Subnet: defaultNetworkSubnet,
				Type:   defaultNetworkType,
			},
		}
	}

	return nil
}

func (webhook *Webhook) mutateAccess(_ context.Context, tpl *v1alpha1.Template) *field.Error {
	if tpl.Spec.Access == nil {
		tpl.Spec.Access = &v1alpha1.AccessSpec{}
	}
	if tpl.Spec.Access.Ssh == nil {
		tpl.Spec.Access.Ssh = &v1alpha1.SshAccessSpec{Enabled: true}
	}
	if tpl.Spec.Access.Console == nil {
		tpl.Spec.Access.Console = &v1alpha1.ConsoleAccessSpec{
			Serial: true,
			Vnc:    true,
		}
	}
	return nil
}
