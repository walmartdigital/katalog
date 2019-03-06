package k8sdriver

import (
	"github.com/seadiaz/katalog/src/domain"
	"k8s.io/api/core/v1"
)

func buildServiceFromK8sService(sourceService *v1.Service) domain.Service {
	destinationService := &domain.Service{
		ID:      string(sourceService.GetUID()),
		Name:    sourceService.GetName(),
		Address: sourceService.Spec.ClusterIP,
		Port:    int(sourceService.Spec.Ports[0].Port),
	}

	return *destinationService
}
