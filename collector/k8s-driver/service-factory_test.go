package k8sdriver_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/walmartdigital/katalog/collector/k8s-driver"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Service builder struct", func() {

	BeforeEach(func() {})

	It("should build a Service object when pass k8sService resource", func() {
		service := BuildServiceFromK8sService(buildService())

		Expect(service.GetID()).To(Equal("UIDExample"))
		Expect(service.GetObservedGeneration()).To(Equal(int64(0)))
		Expect(service.GetName()).To(Equal("ServiceNameExample"))
		Expect(service.GetPort()).To(Equal(3200))
		Expect(service.GetAddress()).To(Equal("127.0.0.1"))
		Expect(service.GetNamespace()).To(Equal("ServiceNameSpaceExample"))
		Expect(service.GetLabels()).To(Equal(map[string]string{"keyLabelExample": "valueLabelExample"}))
		Expect(service.GetTimestamp()).Should(BeTemporally(">", time.Time{}))
	})

})

func buildService() *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:       "ServiceNameExample",
			Namespace:  "ServiceNameSpaceExample",
			UID:        "UIDExample",
			Generation: 5,
			Labels:     map[string]string{"keyLabelExample": "valueLabelExample"},
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: "127.0.0.1",
			Ports:     []corev1.ServicePort{{Port: 3200}, {Port: 8900}},
		},
	}
}
