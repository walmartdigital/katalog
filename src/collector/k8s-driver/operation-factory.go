package k8sdriver

import (
	"github.com/walmartdigital/katalog/src/domain"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func buildOperationFromK8sService(kind domain.OperationType, sourceService *corev1.Service, endpoints v1.Endpoints) domain.Operation {
	destinationService := buildServiceFromK8sService(sourceService)
	for _, endpoint := range buildEndpointFromK8sEndpoints(endpoints) {
		destinationService.AddInstance(endpoint)
	}
	operation := &domain.Operation{
		Kind:    kind,
		Service: destinationService,
	}

	return *operation
}

func buildOperationFromK8sDeployment(kind domain.OperationType, sourceDeployment *appsv1.Deployment) domain.Operation {
	destinationDeployment := buildDeploymentFromK8sDeployment(sourceDeployment)

	operation := &domain.Operation{
		Kind: kind,
		Resource: &domain.Resource{
			Type:   "Deployment",
			Object: destinationDeployment.(*interface{}),
		},
	}

	return *operation
}
