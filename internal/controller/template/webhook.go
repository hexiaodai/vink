package template

import (
	"context"

	"github.com/kubevm.io/vink/pkg/k8s/apis/vink/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

type Webhook struct {
	Client client.Client
}

func (webhook *Webhook) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(&v1alpha1.Template{}).
		WithDefaulter(webhook).
		WithValidator(webhook).
		Complete()
}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the type.
func (webhook *Webhook) Default(ctx context.Context, obj runtime.Object) error {
	return webhook.mutate(ctx, obj)
}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (webhook *Webhook) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	return nil, webhook.validate(ctx, obj)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (webhook *Webhook) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	return nil, webhook.validate(ctx, newObj)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (webhook *Webhook) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	return nil, nil
}
