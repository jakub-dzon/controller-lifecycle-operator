package resources

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// ResourceBuilder helps in building k8s resources
type ResourceBuilder struct {
	commonLabels   map[string]string
	operatorLabels map[string]string
}

// NewResourceBuilder creates new ResourceBuilder
func NewResourceBuilder(commonLabels map[string]string, operatorLabels map[string]string) ResourceBuilder {
	return ResourceBuilder{
		commonLabels:   commonLabels,
		operatorLabels: operatorLabels,
	}
}

// WithCommonLabels aggregates common lables
func (b *ResourceBuilder) WithCommonLabels(labels map[string]string) map[string]string {
	if labels == nil {
		labels = make(map[string]string)
	}

	for k, v := range b.commonLabels {
		_, ok := labels[k]
		if !ok {
			labels[k] = v
		}
	}

	return labels
}

// WithOperatorLabels aggregates common lables
func (b *ResourceBuilder) WithOperatorLabels(labels map[string]string) map[string]string {
	if labels == nil {
		labels = make(map[string]string)
	}

	for k, v := range b.operatorLabels {
		_, ok := labels[k]
		if !ok {
			labels[k] = v
		}
	}

	return labels
}

// CreateServiceAccount creates service account
func (b *ResourceBuilder) CreateServiceAccount(name string) *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ServiceAccount",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: b.WithCommonLabels(nil),
		},
	}
}

// CreateOperatorServiceAccount creates service account
func (b *ResourceBuilder) CreateOperatorServiceAccount(name, namespace string) *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ServiceAccount",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    b.WithOperatorLabels(nil),
		},
	}
}

// CreateOperatorDeploymentSpec creates deployment
func (b *ResourceBuilder) CreateOperatorDeploymentSpec(name, matchKey, matchValue, serviceAccount string, numReplicas int32) *appsv1.DeploymentSpec {
	matchMap := map[string]string{matchKey: matchValue}
	spec := &appsv1.DeploymentSpec{
		Replicas: &numReplicas,
		Selector: &metav1.LabelSelector{
			MatchLabels: b.WithOperatorLabels(matchMap),
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: b.WithOperatorLabels(matchMap),
			},
			Spec: corev1.PodSpec{
				SecurityContext: &corev1.PodSecurityContext{
					RunAsNonRoot: &[]bool{true}[0],
				},
			},
		},
	}

	if serviceAccount != "" {
		spec.Template.Spec.ServiceAccountName = serviceAccount
	}

	return spec
}

// CreateOperatorDeployment creates deployment
func (b *ResourceBuilder) CreateOperatorDeployment(name, namespace, matchKey, matchValue, serviceAccount string, numReplicas int32) *appsv1.Deployment {
	deployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: *b.CreateOperatorDeploymentSpec(name, matchKey, matchValue, serviceAccount, numReplicas),
	}
	if serviceAccount != "" {
		deployment.Spec.Template.Spec.ServiceAccountName = serviceAccount
	}
	return deployment
}

// CreateDeployment creates deployment
func (b *ResourceBuilder) CreateDeployment(name, matchKey, matchValue, serviceAccount string, numReplicas int32) *appsv1.Deployment {
	matchMap := map[string]string{matchKey: matchValue}
	deployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: b.WithCommonLabels(matchMap),
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &numReplicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					matchKey: matchValue,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: b.WithCommonLabels(matchMap),
				},
				Spec: corev1.PodSpec{
					SecurityContext: &corev1.PodSecurityContext{
						RunAsNonRoot: &[]bool{true}[0],
					},
				},
			},
		},
	}
	if serviceAccount != "" {
		deployment.Spec.Template.Spec.ServiceAccountName = serviceAccount
	}
	return deployment
}

// CreatePortsContainer creates container
func (b *ResourceBuilder) CreatePortsContainer(name, image, verbosity string, pullPolicy corev1.PullPolicy, ports *[]corev1.ContainerPort) corev1.Container {
	return corev1.Container{
		Name:            name,
		Image:           image,
		Ports:           *ports,
		Args:            []string{"-v=" + verbosity},
		ImagePullPolicy: pullPolicy,
	}
}

// CreateContainer creates container
func (b *ResourceBuilder) CreateContainer(name, image, verbosity string, pullPolicy corev1.PullPolicy) corev1.Container {
	return corev1.Container{
		Name:                     name,
		Image:                    image,
		ImagePullPolicy:          pullPolicy,
		Args:                     []string{"-v=" + verbosity},
		TerminationMessagePolicy: corev1.TerminationMessageReadFile,
		TerminationMessagePath:   corev1.TerminationMessagePathDefault,
	}
}

// CreateService creates service
func (b *ResourceBuilder) CreateService(name, matchKey, matchValue string) *corev1.Service {
	matchMap := map[string]string{matchKey: matchValue}
	labelMap := map[string]string{matchKey: matchValue}
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: b.WithCommonLabels(labelMap),
		},
		Spec: corev1.ServiceSpec{
			Selector: matchMap,
		},
	}
}

// ValidateGVKs makes sure all resources have initialized GVKs
func ValidateGVKs(objects []runtime.Object) {
	for _, obj := range objects {
		gvk := obj.GetObjectKind().GroupVersionKind()
		if gvk.Version == "" || gvk.Kind == "" {
			panic(fmt.Sprintf("Uninitialized GVK for %+v", obj))
		}
	}
}
