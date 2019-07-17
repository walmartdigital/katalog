package k8sdriver

import (
	"github.com/walmartdigital/katalog/src/domain"
	appsv1 "k8s.io/api/apps/v1"
)

func buildStatefulSetFromK8sStatefulSet(sourceStatefulSet *appsv1.StatefulSet) domain.StatefulSet {
	destinationStatefulSet := &domain.StatefulSet{
		ID:         string(sourceStatefulSet.GetUID()),
		Name:       sourceStatefulSet.GetName(),
		Generation: sourceStatefulSet.GetGeneration(),
		Namespace:  sourceStatefulSet.GetNamespace(),
	}

	return *destinationStatefulSet
}
