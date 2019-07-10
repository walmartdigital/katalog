package repositories

import (
	"github.com/walmartdigital/katalog/src/domain"
)

// Repository ...
type Repository interface {
	CreateResource(obj interface{}) error
	UpdateResource(obj interface{}) (*domain.Resource, error)
	DeleteResource(obj interface{}) error
	GetAllResources() []interface{}
}
