package k8sdriver

import (
	"time"

	"github.com/walmartdigital/katalog/domain"
	appsv1 "k8s.io/api/apps/v1"
)

// BuildStatefulSetFromK8sStatefulSet ...
func buildStatefulSetFromK8sStatefulSet(sourceStatefulSet *appsv1.StatefulSet) domain.StatefulSet {
	m := make(map[string]string)

	for _, c := range sourceStatefulSet.Spec.Template.Spec.Containers {
		m[c.Name] = c.Image
	}

	destinationStatefulSet := &domain.StatefulSet{
		ID:                 string(sourceStatefulSet.GetUID()),
		Name:               sourceStatefulSet.GetName(),
		Generation:         sourceStatefulSet.GetGeneration(),
		Namespace:          sourceStatefulSet.GetNamespace(),
		Labels:             sourceStatefulSet.GetLabels(),
		Containers:         m,
		Timestamp:          time.Now().UTC().Format("2006-01-02 15:04:05"),
		ObservedGeneration: sourceStatefulSet.Status.ObservedGeneration,
	}

	return *destinationStatefulSet
}
