# This is a wrapper to manage helm chart
#
# All make targets related to helm are defined in this file.

.PHONY: helm.package
helm.package:
	@$(LOG_TARGET)
	@helm package helm --destination temp --debug --version $(VERSION) --app-version $(VERSION)

.PHONY: helm.generate-template
helm.generate-template:
	@$(LOG_TARGET)
	@helm -n vink template \
		--set vink.image.repository=$(IMAGE_AGENT) \
		deploy/$(HELM_NAME)-$(VERSION).tgz > deploy/vink.yaml
	@cp -r charts/crds deploy/

.PHONY: helm.push
helm.push:
	@$(LOG_TARGET)
	@helm push temp/$(HELM_NAME)-$(VERSION).tgz $(OCI_REGISTRY)

# kk create cluster --with-kubernetes v1.24.17 --container-manager containerd

# k create ns vink
# k create ns kubevirt
# k create ns cdi
# k create ns local-path-storage

# kubectl label node -lbeta.kubernetes.io/os=linux kubernetes.io/os=linux --overwrite
# kubectl label node -lnode-role.kubernetes.io/control-plane kube-ovn/role=master --overwrite

# helm pull oci://registry-1.docker.io/hejianmin/chart-vink --version 0.0.1-4e596d93

# kubectl delete configmap -n vink --selector=owner=helm,name=vink
# kubectl delete secret -n vink --selector=owner=helm,name=vink

# k delete serviceaccounts -n kube-system kube-ovn-pre-delete-hook
# k delete clusterroles.rbac.authorization.k8s.io system:kube-ovn-pre-delete-hook
# k delete clusterrolebindings.rbac.authorization.k8s.io kube-ovn-pre-delete-hook

##@ Helm

.PHONY: helm.release
helm.release: ## Package virtnet helm chart for release.
helm.release: helm.package helm.generate-template helm.push
