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
func (r *ServiceRepository) CreateService(obj interface{}) error {
	service := obj.(domain.Service)
	if err := r.persistence.Create(service.ID, service); err != nil {
		return err
	}
	return nil
}

// DeleteService ...
func (r *ServiceRepository) DeleteService(obj interface{}) error {
	id := obj.(string)
	if err := r.persistence.Delete(id); err != nil {
		return err
	}
	return nil
}

// GetAllServices ...
func (r *ServiceRepository) GetAllServices() []interface{} {
	glog.Info("get all services called")
	list := arraylist.New()
	services := r.persistence.GetAll()
	for _, item := range services {
		var service domain.Service
		mapstructure.Decode(item, &service)
		list.Add(service)
	}

	return list.Values()
}
