package reconciler_test

import (
	"github.com/jakub-dzon/controller-lifecycle-operator-sdk/pkg/sdk"
	sdkapi "github.com/jakub-dzon/controller-lifecycle-operator-sdk/pkg/sdk/api"
	"github.com/jakub-dzon/controller-lifecycle-operator-sdk/pkg/sdk/reconciler"
	testcr "github.com/jakub-dzon/controller-lifecycle-operator-sdk/tests/cr"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	realClient "sigs.k8s.io/controller-runtime/pkg/client"
	fakeClient "sigs.k8s.io/controller-runtime/pkg/client/fake"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

const (
	finalizerName = "my-finalizer"
	namespace     = "my-namespace"
)

var _ = Describe("Reconciler", func() {

	It("should init the CR", func() {
		cr := createConfig("test", "unique-id")
		operatorVersion := "v0.0.1"

		s := scheme.Scheme
		testcr.AddToScheme(s)
		client := createClient(s, cr)

		r := createReconciler(client, s)

		err := r.CrInit(cr, operatorVersion)

		Expect(err).ToNot(HaveOccurred())

		Expect(cr.Status.Phase).To(BeEquivalentTo(sdkapi.PhaseDeploying))
		Expect(cr.Status.TargetVersion).To(BeEquivalentTo(operatorVersion))
		Expect(cr.Status.OperatorVersion).To(BeEquivalentTo(operatorVersion))

		Expect(cr.GetFinalizers()).To(HaveLen(1))
		Expect(cr.GetFinalizers()[0]).To(BeEquivalentTo(finalizerName))
	})
})

func createClient(scheme *runtime.Scheme, objs ...runtime.Object) realClient.Client {
	return fakeClient.NewFakeClientWithScheme(scheme, objs...)
}

func createReconciler(client realClient.Client, s *runtime.Scheme) *reconciler.Reconciler {
	var configCrManager = testcr.ConfigCrManager{}
	var log = logf.Log.WithName("tests")
	var callbackDispatcher = sdk.NewCallbackDispatcher(log, client, client, s, namespace)

	return reconciler.NewReconciler(&configCrManager, log, client, callbackDispatcher, s, "create-version", "update-version", "last-applied-config", 0, finalizerName)
}

func createConfig(name, uid string) *testcr.Config {
	return &testcr.Config{ObjectMeta: metav1.ObjectMeta{Name: name, UID: types.UID(uid)}, Status: testcr.ConfigStatus{}}
}
