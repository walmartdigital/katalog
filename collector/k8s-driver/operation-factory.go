package k8sdriver

import (
	"github.com/walmartdigital/katalog/domain"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func buildOperationFromK8sService(kind domain.OperationType, sourceService *corev1.Service, endpoints corev1.Endpoints) domain.Operation {
	destinationService := buildServiceFromK8sService(sourceService)
	for _, endpoint := range buildEndpointFromK8sEndpoints(endpoints) {
		destinationService.AddInstance(endpoint)
	}
	resource := &domain.Resource{
		K8sResource: &destinationService,
	}
	operation := &domain.Operation{
		Kind:     kind,
		Resource: *resource,
	}
	return *operation
}

func buildOperationFromK8sDeployment(kind domain.OperationType, sourceDeployment *appsv1.Deployment) domain.Operation {
	destinationDeployment := BuildDeploymentFromK8sDeployment(sourceDeployment)
	resource := &domain.Resource{
		K8sResource: &destinationDeployment,
	}
	operation := &domain.Operation{
		Kind:     kind,
		Resource: *resource,
	}
	return *operation
}

func buildOperationFromK8sStatefulSet(kind domain.OperationType, sourceStatefulSet *appsv1.StatefulSet) domain.Operation {
	destinationStatefulSet := BuildStatefulSetFromK8sStatefulSet(sourceStatefulSet)
	resource := &domain.Resource{
		K8sResource: &destinationStatefulSet,
	}
	operation := &domain.Operation{
		Kind:     kind,
		Resource: *resource,
	}
	return *operation
}
