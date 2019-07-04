package k8sdriver

import (
	"github.com/walmartdigital/katalog/src/domain"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func buildOperationFromK8sService(kind domain.OperationType, sourceService *corev1.Service, endpoints corev1.Endpoints) domain.Operation {
	destinationService := buildServiceFromK8sService(sourceService)
	for _, endpoint := range buildEndpointFromK8sEndpoints(endpoints) {
		destinationService.AddInstance(endpoint)
	}
	resource := &domain.Resource{
		Type:   "Service",
		Object: destinationService,
	}
	operation := &domain.Operation{
		Kind:     kind,
		Resource: *resource,
	}
	return *operation
}

func buildOperationFromK8sDeployment(kind domain.OperationType, sourceDeployment *appsv1.Deployment) domain.Operation {
	destinationDeployment := buildDeploymentFromK8sDeployment(sourceDeployment)
	resource := &domain.Resource{
		Type:   "Deployment",
		Object: destinationDeployment,
	}
	operation := &domain.Operation{
		Kind:     kind,
		Resource: *resource,
	}
	return *operation
}
