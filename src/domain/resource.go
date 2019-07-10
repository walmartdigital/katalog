package domain

import (
	"reflect"
)

// K8sResource ...
type K8sResource interface {
	GetType() reflect.Type
	GetID() string
	GetK8sResource() interface{}
	GetGeneration() int64
}

// Resource ...
type Resource struct {
	K8sResource K8sResource `json:"k8s-resource"`
}

// GetType ...
func (r *Resource) GetType() reflect.Type {
	return r.K8sResource.GetType()
}

// GetK8sResource ...
func (r *Resource) GetK8sResource() interface{} {
	return r.K8sResource.GetK8sResource()
}

// GetID ...
func (r *Resource) GetID() string {
	return r.K8sResource.GetID()
}

// GetGeneration ...
func (r *Resource) GetGeneration() int64 {
	return r.K8sResource.GetGeneration()
}
