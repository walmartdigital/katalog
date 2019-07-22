package repositories

import (
	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/mitchellh/mapstructure"
	"github.com/walmartdigital/katalog/src/domain"
	"github.com/walmartdigital/katalog/src/server/persistence"
	"k8s.io/klog"
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

// GetResource ...
func (r *ResourceRepository) GetResource(id string) (interface{}, error) {
	resource, err := r.persistence.Get(id)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

// UpdateResource ...
func (r *ResourceRepository) UpdateResource(resource interface{}) (*domain.Resource, error) {
	res := resource.(domain.Resource)
	savedResource, err := r.persistence.Get(res.GetID())
	if err != nil {
		return nil, err
	}
	sr := savedResource.(domain.Resource)
	if &sr != nil {
		if sr.GetGeneration() < res.GetGeneration() {
			r.persistence.Update(res.GetID(), res)
			return &res, nil
		}
	}
	return nil, nil
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
func (r *ResourceRepository) GetAllResources() ([]interface{}, error) {
	klog.Info("get all resourcess called")
	list := arraylist.New()
	resources, err := r.persistence.GetAll()
	if err != nil {
		return nil, err
	}
	for _, item := range resources {
		var resource domain.Resource
		mapstructure.Decode(item, &resource)
		list.Add(resource)
	}

	return list.Values(), nil
}
