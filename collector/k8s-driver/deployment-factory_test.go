package k8sdriver

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Deployment builder struct", func() {

	BeforeEach(func() {})

	It("should build a Deployment object when pass k8Deployment resource", func() {
		deployment := buildDeploymentFromK8sDeployment(buildDeployment())

		Expect(deployment.GetID()).To(Equal("UIDExample"))
		Expect(deployment.GetObservedGeneration()).To(Equal(int64(1)))
		Expect(deployment.GetGeneration()).To(Equal(int64(5)))
		Expect(deployment.GetName()).To(Equal("NameExample"))
		Expect(deployment.GetNamespace()).To(Equal("NameSpaceExample"))
		Expect(deployment.GetLabels()).To(Equal(map[string]string{"keyLabelExample": "valueLabelExample"}))
		Expect(deployment.GetAnnotations()).To(Equal(map[string]string{"keyAnnotationsExample": "valueAnnotationsExample"}))
		Expect(deployment.GetContainers()).To(Equal(map[string]string{"containerNameExample": "containerImageExample"}))
		Expect(deployment.GetTimestamp()).Should(BeTemporally(">", time.Time{}))
	})

})

func buildDeployment() *appsv1.Deployment {
	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:        "NameExample",
			Namespace:   "NameSpaceExample",
			UID:         "UIDExample",
			Generation:  5,
			Labels:      map[string]string{"keyLabelExample": "valueLabelExample"},
			Annotations: map[string]string{"keyAnnotationsExample": "valueAnnotationsExample"},
		},
		Spec: appsv1.DeploymentSpec{
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
		Status: appsv1.DeploymentStatus{
			ObservedGeneration: 1,
		},
	}
}
