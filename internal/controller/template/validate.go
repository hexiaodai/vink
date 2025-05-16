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

var (
	networkField *field.Path = field.NewPath("spec").Child("network").Child("interfaces").Child("type")
	userField    *field.Path = field.NewPath("spec").Child("general").Child("user")
	osField      *field.Path = field.NewPath("spec").Child("general").Child("os")
	sourceField  *field.Path = field.NewPath("spec").Child("general").Child("source")
	// storageField *field.Path = field.NewPath("spec").Child("storage")
)

var allowedInterTypes = map[string]struct{}{"bridge": {}, "sriov": {}, "masquerade": {}}

func (webhook *Webhook) validate(ctx context.Context, obj runtime.Object) error {
	tpl, ok := obj.(*v1alpha1.Template)
	if !ok {
		return fmt.Errorf("object is not a template")
	}

	var allErrs field.ErrorList

	if errs := webhook.validateSourceExclusive(ctx, tpl); errs != nil {
		allErrs = append(allErrs, errs...)
	}
	if errs := webhook.validateOs(ctx, tpl); errs != nil {
		allErrs = append(allErrs, errs...)
	}
	if errs := webhook.validateNetwork(ctx, tpl); errs != nil {
		allErrs = append(allErrs, errs...)
	}
	if errs := webhook.validateUser(ctx, tpl); errs != nil {
		allErrs = append(allErrs, errs...)
	}

	if len(allErrs) == 0 {
		return nil
	}
	return errors.NewInvalid(schema.GroupKind{Group: v1alpha1.GroupVersion.Group, Kind: "Template"}, tpl.Name, allErrs)
}

func (webhook *Webhook) validateNetwork(_ context.Context, tpl *v1alpha1.Template) (errs field.ErrorList) {
	var (
		network       = tpl.Spec.Network
		duplicatedNad = map[string]int{}
	)

	if network == nil {
		return
	}

	for _, iface := range network.Interfaces {
		duplicatedNad[iface.Nad]++
		if _, ok := allowedInterTypes[iface.Type]; !ok {
			errs = append(errs, field.Invalid(networkField, iface.Type, "invalid interface type"))
		}
	}
	for name, count := range duplicatedNad {
		if count > 1 {
			errs = append(errs, field.Invalid(networkField, name, "duplicated nad"))
		}
	}
	for _, iface := range network.Interfaces {
		if len(iface.Subnet) == 0 {
			errs = append(errs, field.Invalid(networkField, iface.Subnet, "subnet cannot be empty"))
		}
	}
	return
}

func (webhook *Webhook) validateUser(_ context.Context, tpl *v1alpha1.Template) (errs field.ErrorList) {
	var (
		user          = tpl.Spec.General.User
		passwordCount = 0
		sshKeyCount   = 0
	)
	if len(user.Password) > 0 {
		passwordCount++
	}
	if len(user.PasswordBase64) > 0 {
		passwordCount++
	}
	if len(user.PasswordSecretRef) > 0 {
		passwordCount++
	}
	if len(user.SshKey) > 0 {
		sshKeyCount++
	}
	if len(user.SshKeyBase64) > 0 {
		sshKeyCount++
	}
	if passwordCount > 1 {
		errs = append(errs, field.Invalid(userField, user, "only one of password, passwordBase64, passwordSecretRef can be set"))
	}
	if sshKeyCount > 1 {
		errs = append(errs, field.Invalid(userField, user, "only one of sshKey, sshKeyBase64, sshKeySecretRef can be set"))
	}
	return
}

func (webhook *Webhook) validateOs(_ context.Context, tpl *v1alpha1.Template) (errs field.ErrorList) {
	var (
		os     = tpl.Spec.General.Os
		source = tpl.Spec.General.Source
	)

	if source.Builtin == nil && (os == nil || len(os.Name) == 0 || len(os.Version) == 0) {
		errs = append(errs, field.Invalid(osField, os, "os is required"))
	}

	if source.Builtin != nil && os != nil {
		if os.Name != source.Builtin.Distribution && os.Version != source.Builtin.Version {
			errs = append(errs, field.Invalid(osField, os, "os name and version must be the same as the builtin image"))
		}
	}
	return
}

func (webhook *Webhook) validateSourceExclusive(_ context.Context, tpl *v1alpha1.Template) (errs field.ErrorList) {
	var (
		count  = 0
		source = tpl.Spec.General.Source
	)
	if source.Builtin != nil {
		count++
	}
	if source.Registry != nil {
		count++
	}
	if source.Http != nil {
		count++
	}
	if source.S3 != nil {
		count++
	}
	if source.Pvc != nil {
		count++
	}
	if source.DataVolume != nil {
		count++
	}
	if count != 1 {
		errs = append(errs, field.Invalid(sourceField, source, fmt.Sprintf("exactly one of [builtin, registry, http, s3, pvc, dataVolume] must be set, but got %d", count)))
	}
	return
}
