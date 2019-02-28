package k8sdriver

import (
	"github.com/seadiaz/katalog/domain"
	"k8s.io/api/core/v1"
)

func buildOperationFromK8sService(kind domain.OperationType, sourceService *v1.Service, endpoints v1.Endpoints) domain.Operation {
	destinationService := buildServiceFromK8sService(sourceService)
	for _, endpoint := range buildEndpointFromK8sEndpoints(endpoints) {
		destinationService.AddEndpoint(endpoint)
	}
	operation := &domain.Operation{
		Kind:    kind,
		Service: destinationService,
	}

	return *operation
}
