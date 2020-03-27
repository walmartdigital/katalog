package domain

import "reflect"

// StatefulSet ...
type StatefulSet struct {
	ID          string            `json:",omitempty"`
	Name        string            `json:",omitempty"`
	Generation  int64             `json:",omitempty"`
	Namespace   string            `json:",omitempty"`
	Labels      map[string]string `json:",omitempty"`
	Annotations map[string]string `json:",omitempty"`
	Containers  map[string]string `json:",omitempty"`
}

// GetID ...
func (s *StatefulSet) GetID() string {
	return s.ID
}

// GetType ...
func (s *StatefulSet) GetType() reflect.Type {
	return reflect.TypeOf(s)
}

// GetK8sResource ...
func (s *StatefulSet) GetK8sResource() interface{} {
	return s
}

// GetGeneration ...
func (s *StatefulSet) GetGeneration() int64 {
	return s.Generation
}

// GetNamespace ...
func (s *StatefulSet) GetNamespace() string {
	return s.Namespace
}

// GetName ...
func (s *StatefulSet) GetName() string {
	return s.Name
}

// GetLabels ...
func (s *StatefulSet) GetLabels() map[string]string {
	return s.Labels
}

// GetContainers ...
func (s *StatefulSet) GetContainers() map[string]string {
	return s.Containers
}

// GetAnnotations ...
func (s *StatefulSet) GetAnnotations() map[string]string {
	return s.Annotations
}
