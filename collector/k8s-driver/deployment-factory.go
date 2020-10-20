package k8sdriver

import (
	"time"

	"github.com/walmartdigital/katalog/domain"
	appsv1 "k8s.io/api/apps/v1"
)

func buildDeploymentFromK8sDeployment(sourceDeployment *appsv1.Deployment) domain.Deployment {
	m := make(map[string]string)

	for _, c := range sourceDeployment.Spec.Template.Spec.Containers {
		m[c.Name] = c.Image
	}

	destinationDeployment := &domain.Deployment{
<<<<<<< HEAD:src/collector/k8s-driver/deployment-factory.go
		ID:          string(sourceDeployment.GetUID()),
		Name:        sourceDeployment.GetName(),
		Generation:  sourceDeployment.GetGeneration(),
		Namespace:   sourceDeployment.GetNamespace(),
		Labels:      sourceDeployment.GetLabels(),
		Annotations: sourceDeployment.GetAnnotations(),
		Containers:  m,
		Timestamp:   time.Now().UTC(),
=======
		ID:                 string(sourceDeployment.GetUID()),
		Name:               sourceDeployment.GetName(),
		Generation:         sourceDeployment.GetGeneration(),
		Namespace:          sourceDeployment.GetNamespace(),
		Labels:             sourceDeployment.GetLabels(),
		Annotations:        sourceDeployment.GetAnnotations(),
		Containers:         m,
		Timestamp:          time.Now().UTC(),
		ObservedGeneration: sourceDeployment.Status.ObservedGeneration,
>>>>>>> 13a3de9... refactor(gomod): convert from dep to gomod [CI SKIP]:collector/k8s-driver/deployment-factory.go
	}

	return *destinationDeployment
}
