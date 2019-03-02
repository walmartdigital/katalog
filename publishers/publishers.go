package publishers

import (
	"github.com/hashicorp/consul/api"
	"github.com/seadiaz/katalog/domain"
)

// Publisher ...
type Publisher interface {
	Publish(obj interface{})
}

func createService(service domain.Service) []api.AgentServiceRegistration {
	output := make([]api.AgentServiceRegistration, len(service.Endpoints))
	for i, endpoint := range service.Endpoints {
		output[i] = api.AgentServiceRegistration{
			ID:      service.ID + "-" + string(i),
			Name:    service.Name,
			Address: endpoint.Address,
			Port:    service.Port,
		}
	}
	return output
}
