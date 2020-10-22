package k8sdriver_test

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/walmartdigital/katalog/collector/k8s-driver"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestAll(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Factory")
}

var _ = Describe("StatefulSet builder struct", func() {

	BeforeEach(func() {})

	It("should build a StatefulSet when pass k8StatefulSet", func() {
		statefulSet := BuildStatefulSetFromK8sStatefulSet(buildStatefulSet())

		Expect(statefulSet.GetID()).To(Equal("UIDExample"))
		Expect(statefulSet.GetObservedGeneration()).To(Equal(int64(1)))
		Expect(statefulSet.GetGeneration()).To(Equal(int64(5)))
		Expect(statefulSet.GetName()).To(Equal("NameExample"))
		Expect(statefulSet.GetNamespace()).To(Equal("NameSpaceExample"))
		Expect(statefulSet.GetLabels()).To(Equal(map[string]string{"keyLabelExample": "valueLabelExample"}))
		Expect(statefulSet.GetContainers()).To(Equal(map[string]string{"containerNameExample": "containerImageExample"}))
		Expect(statefulSet.GetTimestamp()).Should(BeTemporally(">", time.Time{}))
	})

})

func buildStatefulSet() *appsv1.StatefulSet {
	return &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:        "NameExample",
			Namespace:   "NameSpaceExample",
			UID:         "UIDExample",
			Generation:  5,
			Labels:      map[string]string{"keyLabelExample": "valueLabelExample"},
			Annotations: map[string]string{"keyAnnotationsExample": "valueAnnotationsExample"},
		},
		Spec: appsv1.StatefulSetSpec{
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  "containerNameExample",
						Image: "containerImageExample",
					}},
				},
			},
		},
		Status: appsv1.StatefulSetStatus{ObservedGeneration: 1},
	}
}
