##@ grpc

API_VERSIONS = v1alpha1 v1alpha2

.PHONY: grpc.generate
grpc.generate: ## Generated client and server code.
	@$(LOG_TARGET)
