package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/walmartdigital/katalog/src/server/repositories"
)

// Server ...
type Server struct {
	serviceRepository repositories.Repository
}

// CreateServer ...
func CreateServer(repository repositories.Repository) Server {
	return Server{
		serviceRepository: repository,
	}
}

// Run ...
func (s *Server) Run() {
	s.handleRequests()
}

func (s *Server) handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/services", s.getAllServices).Methods("GET")
	myRouter.HandleFunc("/services/_count", s.countServices).Methods("GET")
	myRouter.HandleFunc("/services/{id}", s.createService).Methods("PUT")
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}
