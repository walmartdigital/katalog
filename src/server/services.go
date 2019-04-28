package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/walmartdigital/katalog/src/domain"
)

func (s *Server) createService(w http.ResponseWriter, r *http.Request) {
	var service domain.Service
	json.NewDecoder(r.Body).Decode(&service)
	s.serviceRepository.CreateService(service)
	json.NewEncoder(w).Encode(service)
}

func (s *Server) deleteService(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	s.serviceRepository.DeleteService(id)
	fmt.Fprintf(w, "deleted service id: %s", id)
}

func (s *Server) getAllServices(w http.ResponseWriter, r *http.Request) {
	services := s.serviceRepository.GetAllServices()
	json.NewEncoder(w).Encode(services)
}

func (s *Server) countServices(w http.ResponseWriter, r *http.Request) {
	services := s.serviceRepository.GetAllServices()
	json.NewEncoder(w).Encode(struct{ Count int }{len(services)})
}
