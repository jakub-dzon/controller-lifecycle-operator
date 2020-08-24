package certrotation

import (
	"context"
	"time"

	"github.com/jakub-dzon/controller-lifecycle-operator-sdk/pkg/sdk/resources"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	certapi "github.com/jakub-dzon/controller-lifecycle-operator-sdk/pkg/sdk/certrotation/api"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

const day = 24 * time.Hour

func newCertManagerForTest(client kubernetes.Interface, namespace string) CertManager {
	return NewCertManagerForClient(client, namespace)
}

func checkSecret(client kubernetes.Interface, namespace, name string, exists bool) {
	s, err := client.CoreV1().Secrets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if !exists {
		Expect(errors.IsNotFound(err)).To(BeTrue())
		return
	}
	Expect(err).ToNot(HaveOccurred())
	Expect(s.Data["tls.crt"]).ShouldNot(BeEmpty())
	Expect(s.Data["tls.crt"]).ShouldNot(BeEmpty())
}

func checkConfigMap(client kubernetes.Interface, namespace, name string, exists bool) {
	cm, err := client.CoreV1().ConfigMaps(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if !exists {
		Expect(errors.IsNotFound(err)).To(BeTrue())
		return
	}
	Expect(cm.Data["ca-bundle.crt"]).ShouldNot(BeEmpty())
}

func checkCerts(client kubernetes.Interface, namespace string, exists bool) {
	checkSecret(client, namespace, "signer", exists)
	checkConfigMap(client, namespace, "signer-bundle", exists)

	checkSecret(client, namespace, "client-signer", exists)
	checkConfigMap(client, namespace, "client-signer-bundle", exists)
	checkSecret(client, namespace, "client-cert", exists)
}

var _ = Describe("Cert rotation tests", func() {
	const namespace = "certificates"
	Context("with clean slate", func() {
		client := fake.NewSimpleClientset()
		cm := newCertManagerForTest(client, namespace)

		It("should create everything", func() {
			checkCerts(client, namespace, false)

			certs := createCertificateDefinitions(namespace)
			err := cm.Sync(certs)
			Expect(err).ToNot(HaveOccurred())

			checkCerts(client, namespace, true)
		})

		It("should not do anything", func() {
			checkCerts(client, namespace, true)

			certs := createCertificateDefinitions(namespace)
			err := cm.Sync(certs)
			Expect(err).ToNot(HaveOccurred())

			checkCerts(client, namespace, true)
		})
	})
})

func createCertificateDefinitions(namespace string) []certapi.CertificateDefinition {
	return []certapi.CertificateDefinition{
		{
			SignerSecret:        createSecret("signer", namespace),
			SignerValidity:      10 * 365 * day,
			SignerRefresh:       8 * 365 * day,
			CertBundleConfigmap: createConfigMap("signer-bundle", namespace),
		},
		{
			SignerSecret:        createSecret("client-signer", namespace),
			SignerValidity:      10 * 365 * day,
			SignerRefresh:       8 * 365 * day,
			CertBundleConfigmap: createConfigMap("client-signer-bundle", namespace),
			TargetSecret:        createSecret("client-cert", namespace),
			TargetValidity:      48 * time.Hour,
			TargetRefresh:       24 * time.Hour,
			TargetUser:          &[]string{"client.server.kubevirt.io"}[0],
		},
	}
}

func createConfigMap(name string, namespace string) *corev1.ConfigMap {
	configMap := resources.CreateConfigMap(name, nil)
	configMap.SetNamespace(namespace)
	return configMap
}

func createSecret(name string, namespace string) *corev1.Secret {
	secret := resources.CreateSecret(name, nil)
	secret.SetNamespace(namespace)
	return secret
}
