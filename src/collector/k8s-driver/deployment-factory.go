package k8sdriver

import (
	"github.com/walmartdigital/katalog/src/domain"
	appsv1 "k8s.io/api/apps/v1"
)

func buildDeploymentFromK8sDeployment(sourceDeployment *appsv1.Deployment) domain.Deployment {
	destinationDeployment := &domain.Deployment{
		ID:         string(sourceDeployment.GetUID()),
		Name:       sourceDeployment.GetName(),
		Generation: sourceDeployment.GetGeneration(),
		Namespace:  sourceDeployment.GetNamespace(),
	}

	return *destinationDeployment
}
