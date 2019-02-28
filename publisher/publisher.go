package publisher

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/hashicorp/consul/api"
	"github.com/seadiaz/katalog/domain"
)

// Publish ...
func Publish(obj interface{}) {
	operation := obj.(domain.Operation)
	switch operation.Kind {
	case (domain.OperationTypeAdd):
		consulCreate(operation.Service)
	case (domain.OperationTypeDelete):
		consulDelete(operation.Service)
	}
}

func consulCreate(service domain.Service) {
	client, _ := api.NewClient(api.DefaultConfig())
	agent := client.Agent()
	for _, destinationService := range createService(service) {
		err := agent.ServiceRegister(&destinationService)
		if err != nil {
			glog.Errorln(err)
		}
	}
	fmt.Printf("service %s registered with %d endpoints\n", service.Name, len(service.Endpoints))
}

func consulDelete(service domain.Service) {
	client, _ := api.NewClient(api.DefaultConfig())
	agent := client.Agent()
	err := agent.ServiceDeregister(service.ID)
	if err != nil {
		glog.Errorln(err)
	}
	fmt.Printf("service %s deregistered\n", service.Name)
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
