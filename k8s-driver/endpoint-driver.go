package k8sdriver

import (
	"github.com/seadiaz/katalog/domain"
	"k8s.io/api/core/v1"
)

func buildEndpointFromK8sEndpoints(endpoints v1.Endpoints) []domain.Endpoint {
	output := make([]domain.Endpoint, len(endpoints.Subsets[0].Addresses))
	for i, address := range endpoints.Subsets[0].Addresses {
		output[i] = domain.Endpoint{
			Address: address.IP,
		}
	}

	return output
}
