package k8sdriver

import (
	"github.com/walmartdigital/katalog/src/domain"
	corev1 "k8s.io/api/core/v1"
)

func buildServiceFromK8sService(sourceService *corev1.Service) domain.Service {
	destinationService := &domain.Service{
		ID:        string(sourceService.GetUID()),
		Name:      sourceService.GetName(),
		Address:   sourceService.Spec.ClusterIP,
		Port:      int(sourceService.Spec.Ports[0].Port),
		Namespace: sourceService.GetNamespace(),
	}

	return *destinationService
}
