module github.com/operator-framework/operator-lifecycle-manager

go 1.12

require (
	github.com/blang/semver v3.5.1+incompatible
	github.com/coreos/go-semver v0.3.0
	github.com/davecgh/go-spew v1.1.1
	github.com/fsnotify/fsnotify v1.4.7
	github.com/ghodss/yaml v1.0.0
	github.com/go-openapi/spec v0.19.4
	github.com/golang/mock v1.3.1
	github.com/maxbrunsfeld/counterfeiter/v6 v6.2.2
	github.com/mikefarah/yq/v2 v2.4.1
	github.com/mitchellh/hashstructure v1.0.0
	github.com/onsi/ginkgo v1.10.1
	github.com/onsi/gomega v1.7.0
	github.com/openshift/api v0.0.0-20200205133042-34f0ec8dab87
	github.com/openshift/client-go v0.0.0-20190923180330-3b6373338c9b
	github.com/operator-framework/operator-registry v1.6.1
	github.com/otiai10/copy v1.0.2
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.2.1
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v0.0.5
	github.com/stretchr/testify v1.4.0
	golang.org/x/sys v0.0.0-20191210023423-ac6580df4449 // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0
	google.golang.org/grpc v1.27.0
	gopkg.in/yaml.v2 v2.2.8
	helm.sh/helm/v3 v3.1.2
	k8s.io/api v0.17.3
	k8s.io/apiextensions-apiserver v0.17.3
	k8s.io/apimachinery v0.17.3
	k8s.io/apiserver v0.17.3
	k8s.io/client-go v0.17.3
	k8s.io/code-generator v0.17.3
	k8s.io/component-base v0.17.3
	k8s.io/klog v1.0.0
	k8s.io/kube-aggregator v0.17.3
	k8s.io/kube-openapi v0.0.0-20191107075043-30be4d16710a
	rsc.io/letsencrypt v0.0.3 // indirect
	sigs.k8s.io/controller-tools v0.2.4
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.2+incompatible
	github.com/docker/docker => github.com/moby/moby v0.7.3-0.20190826074503-38ab9da00309 // Required by Helm
	github.com/openshift/api => github.com/openshift/api v3.9.1-0.20190924102528-32369d4db2ad+incompatible
	github.com/openshift/client-go => github.com/openshift/client-go v0.0.0-20190923180330-3b6373338c9b
	github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.3-0.20190127221311-3c4408c8b829
	google.golang.org/grpc => google.golang.org/grpc v1.26.0 // https://github.com/etcd-io/etcd/issues/11563
	sigs.k8s.io/structured-merge-diff => sigs.k8s.io/structured-merge-diff v1.0.1-0.20191108220359-b1b620dd3f06
)
