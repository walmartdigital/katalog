package domain

import "reflect"

// Deployment ...
type Deployment struct {
	ID         string            `json:",omitempty"`
	Name       string            `json:",omitempty"`
	Generation int64             `json:",omitempty"`
	Namespace  string            `json:",omitempty"`
	Labels     map[string]string `json:",omitempty"`
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
