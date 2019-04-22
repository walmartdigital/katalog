package server

import (
	"encoding/json"
	"net/http"

	"github.com/walmartdigital/katalog/src/domain"
)

func (s *Server) createService(w http.ResponseWriter, r *http.Request) {
	var service domain.Service
	json.NewDecoder(r.Body).Decode(&service)
	s.serviceRepository.CreateService(service)
	json.NewEncoder(w).Encode(service)
}

func (s *Server) deleteService(w http.ResponseWriter, r *http.Request) {
	var service domain.Service
	json.NewDecoder(r.Body).Decode(&service)
	s.serviceRepository.DeleteService(service)
	json.NewEncoder(w).Encode(service)
}

func (s *Server) getAllServices(w http.ResponseWriter, r *http.Request) {
	services := s.serviceRepository.GetAllServices()
	json.NewEncoder(w).Encode(services)
}

func (s *Server) countServices(w http.ResponseWriter, r *http.Request) {
	services := s.serviceRepository.GetAllServices()
	json.NewEncoder(w).Encode(struct{ Count int }{len(services)})
}
