package domain

import "reflect"

// Service ...
type Service struct {
	ID         string `json:",omitempty"`
	Name       string `json:",omitempty"`
	Port       int    `json:",omitempty"`
	Address    string `json:",omitempty"`
	Generation int64  `json:",omitempty"`
	Namespace  string `json:",omitempty"`
	Instances  []Instance
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
