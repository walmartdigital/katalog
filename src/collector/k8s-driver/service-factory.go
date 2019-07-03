package k8sdriver

import (
	"github.com/walmartdigital/katalog/src/domain"
	appsv1 "k8s.io/api/core/v1"
)

func buildDeploymentFromK8sService(sourceService *appsv1.Deployment) domain.Deployment {
	destinationService := &domain.Service{
		ID:      string(sourceService.GetUID()),
		Name:    sourceService.GetName(),
		Address: sourceService.Spec.ClusterIP,
		Port:    int(sourceService.Spec.Ports[0].Port),
	}

	return *destinationService
}
