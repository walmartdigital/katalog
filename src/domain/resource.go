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
	GetNamespace() string
	GetName() string
	GetLabels() map[string]string
}

// Resource ...
type Resource struct {
	K8sResource K8sResource `json:"K8sResource"`
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

// GetNamespace ...
func (r *Resource) GetNamespace() string {
	return r.K8sResource.GetNamespace()
}

// GetName ...
func (r *Resource) GetName() string {
	return r.K8sResource.GetName()
}

// GetLabels ...
func (r *Resource) GetLabels() map[string]string {
	return r.K8sResource.GetLabels()
}
