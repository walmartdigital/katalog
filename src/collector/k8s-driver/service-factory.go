package k8sdriver

import (
	"time"

	"github.com/walmartdigital/katalog/src/domain"
	corev1 "k8s.io/api/core/v1"
)

func buildServiceFromK8sService(sourceService *corev1.Service) domain.Service {
	port := 0
	if len(sourceService.Spec.Ports) > 0 {
		port = int(sourceService.Spec.Ports[0].Port)
	}

	destinationService := &domain.Service{
		ID:        string(sourceService.GetUID()),
		Name:      sourceService.GetName(),
		Address:   sourceService.Spec.ClusterIP,
		Port:      port,
		Namespace: sourceService.GetNamespace(),
		Labels:    sourceService.GetLabels(),
		Timestamp: time.Now(),
	}

	return *destinationService
}
