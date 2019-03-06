package server

import (
	"encoding/json"
	"net/http"

	"github.com/seadiaz/katalog/src/domain"
)

// Service ...
type Service struct {
	Title   string `json:"Title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

// Services ...
type Services []Service

type genericResponse struct {
	Count int `json:",omitempty"`
}

func (s *Server) getAllServices(w http.ResponseWriter, r *http.Request) {
	services := s.serviceRepository.GetAllServices()
	json.NewEncoder(w).Encode(services)
}

func (s *Server) createService(w http.ResponseWriter, r *http.Request) {
	var service domain.Service
	json.NewDecoder(r.Body).Decode(&service)
	s.serviceRepository.CreateService(service)
	json.NewEncoder(w).Encode(service)
}

func (s *Server) countServices(w http.ResponseWriter, r *http.Request) {
	services := s.serviceRepository.GetAllServices()
	json.NewEncoder(w).Encode(genericResponse{Count: len(services)})
}
