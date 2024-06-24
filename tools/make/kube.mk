##@ Kubernetes

.PHONY: kube.generate
kube.generate: ## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
kube.generate:
	@$(LOG_TARGET)
