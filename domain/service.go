package domain

import (
	"reflect"
	"time"
)

// Service ...
type Service struct {
	ID          string            `json:"ID"`
	Name        string            `json:"Name"`
	Port        int               `json:"Port"`
	Address     string            `json:"Address"`
	Generation  int64             `json:"Generation"`
	Namespace   string            `json:"Namespace"`
	Instances   []Instance        `json:"Instances"`
	Labels      map[string]string `json:",omitempty"`
	Annotations map[string]string `json:",omitempty"`
	Timestamp   time.Time         `json:"Timestamp"`
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

// GetAnnotations ...
func (s *Service) GetAnnotations() map[string]string {
	return s.Annotations
}

// GetTimestamp ...
func (s *Service) GetTimestamp() time.Time {
	return s.Timestamp
}
