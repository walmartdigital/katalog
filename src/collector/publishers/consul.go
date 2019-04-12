package publishers

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/hashicorp/consul/api"
	"github.com/walmartdigital/katalog/src/domain"
)

// ConsulPublisher ...
type ConsulPublisher struct {
	client api.Client
}

// CreateConsulPublisher ...
func CreateConsulPublisher(addr string) Publisher {
	client, _ := api.NewClient(&api.Config{Address: addr})
	return &ConsulPublisher{client: *client}
}

// Publish ...
func (c *ConsulPublisher) Publish(obj interface{}) {
	operation := obj.(domain.Operation)
	switch operation.Kind {
	case (domain.OperationTypeAdd):
		c.consulCreate(operation.Service)
	case (domain.OperationTypeDelete):
		c.consulDelete(operation.Service)
	}
}

func (c *ConsulPublisher) consulCreate(service domain.Service) {
	agent := c.client.Agent()
	for _, destinationService := range createService(service) {
		err := agent.ServiceRegister(&destinationService)
		if err != nil {
			glog.Errorln(err)
		}
	}
	fmt.Printf("service %s registered with %d endpoints\n", service.Name, len(service.Instances))
}

func (c *ConsulPublisher) consulDelete(service domain.Service) {
	agent := c.client.Agent()
	err := agent.ServiceDeregister(service.ID)
	if err != nil {
		glog.Errorln(err)
	}
	fmt.Printf("service %s deregistered\n", service.Name)
}
