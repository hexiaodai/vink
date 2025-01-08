module github.com/kubevm.io/vink/apis

go 1.22.7

toolchain go1.22.10

replace github.com/kubevm.io/vink => ../

require (
	github.com/envoyproxy/protoc-gen-validate v1.0.4
	github.com/golang/protobuf v1.5.4
	github.com/spf13/cobra v1.8.1
	github.com/tkrajina/typescriptify-golang-structs v0.2.0
	google.golang.org/grpc v1.66.2
	google.golang.org/protobuf v1.34.2
	sigs.k8s.io/yaml v1.4.0
)

require (
	github.com/kubeovn/kube-ovn v1.12.23 // indirect
	github.com/ovn-org/libovsdb v0.7.0 // indirect
	kubevirt.io/api v1.4.0 // indirect
	sigs.k8s.io/controller-runtime v0.19.0 // indirect
)

require (
	github.com/fxamacker/cbor/v2 v2.7.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kubevm.io/vink v0.0.0-00010101000000-000000000000
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/openshift/api v0.0.0-20231207204216-5efc6fca4b2d // indirect
	github.com/openshift/custom-resource-status v1.1.2 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/tkrajina/go-reflector v0.5.5 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	golang.org/x/net v0.30.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
	golang.org/x/text v0.19.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240903143218-8af14fe29dc1 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	k8s.io/api v0.31.1 // indirect
	k8s.io/apiextensions-apiserver v0.31.0 // indirect
	k8s.io/apimachinery v0.31.1 // indirect
	k8s.io/klog/v2 v2.130.1 // indirect
	k8s.io/utils v0.0.0-20240902221715-702e33fdd3c3 // indirect
	kubevirt.io/containerized-data-importer-api v1.59.0 // indirect
	kubevirt.io/controller-lifecycle-operator-sdk/api v0.0.0-20220329064328-f3cc58c6ed90 // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.4.1 // indirect
)
