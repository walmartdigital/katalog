package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/walmartdigital/katalog/src/server/repositories"
)

// HTTPServer ...
type HTTPServer interface {
	ListenAndServe() error
}

// Server ...
type Server struct {
	httpServer        HTTPServer
	serviceRepository repositories.Repository
	router            Router
}

// Router ...
type Router interface {
	HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) *mux.Route
	ServeHTTP(http.ResponseWriter, *http.Request)
}

// CreateServer ...
func CreateServer(server HTTPServer, repository repositories.Repository, router Router) Server {
	return Server{
		httpServer:        server,
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
	s.router.HandleFunc("/services/{service}", s.createService).Methods("PUT")
	s.router.HandleFunc("/services/{id}", s.deleteService).Methods("DELETE")
	s.httpServer.ListenAndServe()
}
