module github.com/jakub-dzon/controller-lifecycle-operator-sdk

go 1.14

require (
	github.com/appscode/jsonpatch v0.0.0-20190108182946-7c0e3b262f30
	github.com/blang/semver v3.5.1+incompatible
	github.com/evanphx/json-patch v4.5.0+incompatible
	github.com/go-logr/logr v0.1.0
	github.com/go-logr/zapr v0.1.1 // indirect
	github.com/onsi/ginkgo v1.12.1
	github.com/onsi/gomega v1.10.1
	github.com/openshift/custom-resource-status v0.0.0-20200602122900-c002fd1547ca
	github.com/openshift/library-go v0.0.0-20200821154433-215f00df72cc
	go.uber.org/multierr v1.3.0 // indirect
	k8s.io/api v0.19.0-rc.2
	k8s.io/apiextensions-apiserver v0.19.0-rc.2
	k8s.io/apimachinery v0.19.0-rc.2
	k8s.io/apiserver v0.19.0-rc.2
	k8s.io/client-go v8.0.0+incompatible
	k8s.io/kubernetes v1.14.0
	sigs.k8s.io/controller-runtime v0.6.2
)

replace (
	github.com/openshift/api => github.com/openshift/api v0.0.0-20200526144822-34f54f12813a
	github.com/openshift/library-go => github.com/mhenriks/library-go v0.0.0-20200804184258-4fc3a5379c7a

	k8s.io/api => k8s.io/api v0.18.6
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.18.6
	k8s.io/apimachinery => k8s.io/apimachinery v0.18.6
	k8s.io/apiserver => k8s.io/apiserver v0.18.6
	k8s.io/client-go => k8s.io/client-go v0.18.6
	k8s.io/component-base => k8s.io/component-base v0.18.6
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.18.6

	sigs.k8s.io/structured-merge-diff => sigs.k8s.io/structured-merge-diff v1.0.0
	vbom.ml/util => github.com/fvbommel/util v0.0.0-20180919145318-efcd4e0f9787
)
