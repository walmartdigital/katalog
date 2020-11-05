package k8sdriver

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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

	It("should build a StatefulSet object when pass k8StatefulSet resource", func() {
		statefulSet := buildStatefulSetFromK8sStatefulSet(buildStatefulSet())

		Expect(statefulSet.GetID()).To(Equal("UIDExample"))
		Expect(statefulSet.GetObservedGeneration()).To(Equal(int64(1)))
		Expect(statefulSet.GetGeneration()).To(Equal(int64(5)))
		Expect(statefulSet.GetName()).To(Equal("NameExample"))
		Expect(statefulSet.GetNamespace()).To(Equal("NameSpaceExample"))
		Expect(statefulSet.GetLabels()).To(Equal(map[string]string{"keyLabelExample": "valueLabelExample"}))
		Expect(statefulSet.GetContainers()).To(Equal(map[string]string{"containerNameExample": "containerImageExample"}))
		Expect(statefulSet.GetTimestamp()).Should(MatchRegexp(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}`))
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
