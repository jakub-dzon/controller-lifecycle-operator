package v1beta1

import (
	sdkapi "github.com/jakub-dzon/controller-lifecycle-operator-sdk/pkg/sdk/api"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ConfigCrManager provides test CR management functionality
type ConfigCrManager struct {
}

// IsCreating checks whether creation of the managed resources will be executed
func (m *ConfigCrManager) IsCreating(cr controllerutil.Object) (bool, error) {
	return false, nil
}

// Create creates empty CR
func (m *ConfigCrManager) Create() controllerutil.Object {
	return new(Config)
}

// Status extracts status from the cr
func (m *ConfigCrManager) Status(cr runtime.Object) *sdkapi.Status {
	return &cr.(*Config).Status.Status
}

// GetAllResources provides all resources managed by the cr
func (m *ConfigCrManager) GetAllResources(cr runtime.Object) ([]runtime.Object, error) {
	return nil, nil
}

// GetDependantResourcesListObjects returns resource list objects of dependant resources
func (m *ConfigCrManager) GetDependantResourcesListObjects() []runtime.Object {
	return nil
}
