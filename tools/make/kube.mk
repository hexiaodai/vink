##@ Kubernetes

.PHONY: kube.generate
kube.generate: ## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
kube.generate:
	@$(LOG_TARGET)
	@tools/bin/controller-gen object:headerFile="$(ROOT_DIR)/tools/boilerplate/boilerplate.go.txt" paths="$(ROOT_DIR)/pkg/k8s/apis/vink/v1alpha1/..."
	@tools/bin/controller-gen crd paths="$(ROOT_DIR)/pkg/k8s/apis/vink/v1alpha1/..." output:crd:dir="$(ROOT_DIR)/manifests/crds"

	@cd $(ROOT_DIR)/pkg/k8s/apis/vink/v1alpha1
	$(GOPATH)/pkg/mod/k8s.io/code-generator@v0.24.2/generate-groups.sh all \
	github.com/kubevm.io/vink/pkg/k8s/client \
	github.com/kubevm.io/vink/pkg/k8s/apis \
	vink:v1alpha1 \
	-h $(ROOT_DIR)/tools/boilerplate/boilerplate.go.txt \
	-v 10
