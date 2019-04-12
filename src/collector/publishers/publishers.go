package publishers

import (
	"github.com/hashicorp/consul/api"
	"github.com/walmartdigital/katalog/src/domain"
)

// Publisher ...
type Publisher interface {
	Publish(obj interface{})
}

func createService(service domain.Service) []api.AgentServiceRegistration {
	output := make([]api.AgentServiceRegistration, len(service.Instances))
	for i, endpoint := range service.Instances {
		output[i] = api.AgentServiceRegistration{
			ID:      service.ID + "-" + string(i),
			Name:    service.Name,
			Address: endpoint.Address,
			Port:    service.Port,
		}
	}
	return output
}
