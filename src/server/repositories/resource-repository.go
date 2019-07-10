package repositories

import (
	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/golang/glog"
	"github.com/mitchellh/mapstructure"
	"github.com/walmartdigital/katalog/src/domain"
	"github.com/walmartdigital/katalog/src/server/persistence"
)

// ResourceRepository ...
type ResourceRepository struct {
	persistence persistence.Persistence
}

// CreateResourceRepository ...
func CreateResourceRepository(persistence persistence.Persistence) *ResourceRepository {
	return &ResourceRepository{
		persistence: persistence,
	}
}

// CreateResource ...
func (r *ResourceRepository) CreateResource(resource interface{}) error {
	res := resource.(domain.Resource)
	if err := r.persistence.Create(res.GetID(), res); err != nil {
		return err
	}

	return nil
}

// UpdateResource ...
func (r *ResourceRepository) UpdateResource(resource interface{}) error {
	res := resource.(domain.Resource)
	if err := r.persistence.Update(res.GetID(), res); err != nil {
		return err
	}

	return nil
}

// DeleteResource ...
func (r *ResourceRepository) DeleteResource(obj interface{}) error {
	id := obj.(string)
	if err := r.persistence.Delete(id); err != nil {
		return err
	}
	return nil
}

// GetAllResources ...
func (r *ResourceRepository) GetAllResources() []interface{} {
	glog.Info("get all resourcess called")
	list := arraylist.New()
	resources := r.persistence.GetAll()
	for _, item := range resources {
		var resources domain.Resource
		mapstructure.Decode(item, &resources)
		list.Add(resources)
	}

	return list.Values()
}
