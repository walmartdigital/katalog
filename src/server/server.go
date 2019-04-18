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
	router            Router
}

// Router ...
type Router interface {
	HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) *mux.Route
	ServeHTTP(http.ResponseWriter, *http.Request)
}

// CreateServer ...
func CreateServer(repository repositories.Repository, router Router) Server {
	return Server{
		serviceRepository: repository,
		router:            router,
	}
}

// Run ...
func (s *Server) Run() {
	s.handleRequests()
}

func (s *Server) handleRequests() {
	s.router.HandleFunc("/services", s.getAllServices).Methods("GET")
	s.router.HandleFunc("/services/_count", s.countServices).Methods("GET")
	s.router.HandleFunc("/services/{id}", s.createService).Methods("PUT")
	s.router.HandleFunc("/services/{id}", s.deleteService).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":10000", s.router))
}
