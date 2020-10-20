package domain

import (
	"reflect"
	"time"
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
	GetAnnotations() map[string]string
	GetTimestamp() time.Time
	GetObserveGeneration() string
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

// GetAnnotations ...
func (r *Resource) GetAnnotations() map[string]string {
	return r.K8sResource.GetAnnotations()
}

// GetTimestamp ...
func (r *Resource) GetTimestamp() time.Time {
	return r.K8sResource.GetTimestamp()
}

// GetObserveGeneration ...
func (r *Resource) GetObserveGeneration() string {
	return r.K8sResource.GetObserveGeneration()
}
