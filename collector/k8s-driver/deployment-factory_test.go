package k8sdriver_test

import (
	"testing"

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

var _ = Describe("Deployment builder struct", func() {

	BeforeEach(func() {})

	It("should build a deployment when pass k8Deployment", func() {

		sourceDeployment := &appsv1.Deployment{
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
				Replicas: nil,
				Selector: nil,
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{{
							Name:  "containerNameExample",
							Image: "containerImageExample",
						}},
					},
				},
				Strategy:                appsv1.DeploymentStrategy{},
				MinReadySeconds:         0,
				RevisionHistoryLimit:    nil,
				Paused:                  false,
				ProgressDeadlineSeconds: nil,
			},
			Status: appsv1.DeploymentStatus{
				ObservedGeneration: 1,
			},
		}
		deployment := BuildDeploymentFromK8sDeployment(sourceDeployment)

		Expect(deployment.GetID()).To(Equal("UIDExample"))
		Expect(deployment.GetObservedGeneration()).To(Equal(int64(1)))
		Expect(deployment.GetGeneration()).To(Equal(int64(5)))
		Expect(deployment.GetName()).To(Equal("NameExample"))
		Expect(deployment.GetNamespace()).To(Equal("NameSpaceExample"))
		Expect(deployment.GetLabels()).NotTo(BeNil())
		Expect(deployment.GetAnnotations()).NotTo(BeNil())
		Expect(deployment.GetContainers()).To(Equal(map[string]string{"containerNameExample": "containerImageExample"}))
		Expect(deployment.GetTimestamp()).NotTo(BeNil())
	})

})
