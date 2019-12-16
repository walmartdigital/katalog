package domain

import "reflect"

// Service ...
type Service struct {
	ID         string            `json:"ID"`
	Name       string            `json:"Name"`
	Port       int               `json:"Port"`
	Address    string            `json:"Address"`
	Generation int64             `json:"Generation"`
	Namespace  string            `json:"Namespace"`
	Instances  []Instance        `json:"Instances"`
	Labels     map[string]string `json:",omitempty"`
}

// AddInstance ...
func (s *Service) AddInstance(endpoint Instance) {
	s.Instances = append(s.Instances, endpoint)
}

// GetID ...
func (s *Service) GetID() string {
	return s.ID
}

// GetType ...
func (s *Service) GetType() reflect.Type {
	return reflect.TypeOf(s)
}

// GetK8sResource ...
func (s *Service) GetK8sResource() interface{} {
	return s
}

// GetGeneration ...
func (s *Service) GetGeneration() int64 {
	return s.Generation
}

// GetNamespace ...
func (s *Service) GetNamespace() string {
	return s.Namespace
}

// GetName ...
func (s *Service) GetName() string {
	return s.Name
}

// GetLabels ...
func (s *Service) GetLabels() map[string]string {
	return s.Labels
}
