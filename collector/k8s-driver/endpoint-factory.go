package k8sdriver

import (
	"github.com/walmartdigital/katalog/domain"
	v1 "k8s.io/api/core/v1"
)

func buildEndpointFromK8sEndpoints(endpoints v1.Endpoints) []domain.Instance {
	if len(endpoints.Subsets) <= 0 {
		return make([]domain.Instance, 0)
	}

	output := make([]domain.Instance, len(endpoints.Subsets[0].Addresses))
	for i, address := range endpoints.Subsets[0].Addresses {
		output[i] = domain.Instance{
			Address: address.IP,
		}
	}

	return output
}
