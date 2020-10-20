package domain

import (
	"reflect"
	"time"
)

// Deployment ...
type Deployment struct {
	ID                 string            `json:",omitempty"`
	Name               string            `json:",omitempty"`
	Generation         int64             `json:",omitempty"`
	Namespace          string            `json:",omitempty"`
	Labels             map[string]string `json:",omitempty"`
	Annotations        map[string]string `json:",omitempty"`
	Containers         map[string]string `json:",omitempty"`
	Timestamp          time.Time         `json:"Timestamp"`
	ObservedGeneration int64             `json:",omitempty"`
}

// GetID ...
func (s *Deployment) GetID() string {
	return s.ID
}

// GetType ...
func (s *Deployment) GetType() reflect.Type {
	return reflect.TypeOf(s)
}

// GetK8sResource ...
func (s *Deployment) GetK8sResource() interface{} {
	return s
}

// GetGeneration ...
func (s *Deployment) GetGeneration() int64 {
	return s.Generation
}

// GetNamespace ...
func (s *Deployment) GetNamespace() string {
	return s.Namespace
}

// GetName ...
func (s *Deployment) GetName() string {
	return s.Name
}

// GetLabels ...
func (s *Deployment) GetLabels() map[string]string {
	return s.Labels
}

// GetAnnotations ...
func (s *Deployment) GetAnnotations() map[string]string {
	return s.Annotations
}

// GetContainers ...
func (s *Deployment) GetContainers() map[string]string {
	return s.Containers
}

// GetTimestamp ...
func (s *Deployment) GetTimestamp() time.Time {
	return s.Timestamp
}

// GetObservedGeneration ...
func (s *Deployment) GetObservedGeneration() int64 {
	return s.ObservedGeneration
}