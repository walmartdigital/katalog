package repositories

import (
  "github.com/golang/glog"
	"github.com/seadiaz/katalog/src/domain"
	"github.com/seadiaz/katalog/src/server/persistence"
  "github.com/emirpasic/gods/lists/arraylist"
)

const kind = "service"

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

// GetAllServices ...
func (r *ServiceRepository) GetAllServices() []interface{} {
  glog.Info("get all services called")
  list := arraylist.New()
  services := r.persistence.GetAll(kind)
  for _, item := range services {
    service := item.(domain.Service)
    list.Add(service)
  }

  return list.Values()
}
