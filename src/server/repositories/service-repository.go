package repositories

import (
	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/golang/glog"
	"github.com/mitchellh/mapstructure"
	"github.com/walmartdigital/katalog/src/domain"
	"github.com/walmartdigital/katalog/src/server/persistence"
)

const kind = "services"

// ServiceRepository ...
type ServiceRepository struct {
	persistence persistence.Persistence
}

// CreateServiceRepository ...
func CreateServiceRepository(persistence persistence.Persistence) *ServiceRepository {
	return &ServiceRepository{
		persistence: persistence,
	}
}

// CreateService ...
func (r *ServiceRepository) CreateService(obj interface{}) {
	service := obj.(domain.Service)
	r.persistence.Create(kind, service.ID, service)
}

// DeleteService ...
func (r *ServiceRepository) DeleteService(obj interface{}) {
	service := obj.(domain.Service)
	r.persistence.Delete(kind, service.ID)
}

// GetAllServices ...
func (r *ServiceRepository) GetAllServices() []interface{} {
	glog.Info("get all services called")
	list := arraylist.New()
	services := r.persistence.GetAll(kind)
	for _, item := range services {
		var service domain.Service
		mapstructure.Decode(item, &service)
		list.Add(service)
	}

	return list.Values()
}
