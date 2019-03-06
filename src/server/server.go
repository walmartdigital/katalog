package server

import (
	"log"
	"net/http"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/seadiaz/katalog/src/server/persistence"
)

// Server ...
type Server struct {
	persistence persistence.Persistence
}

// CreateServer ...
func CreateServer(persistence persistence.Persistence) Server {
	return Server{
		persistence: persistence,
	}
}

// Run ...
func (s *Server) Run() {
	glog.Info("server starging...")
	s.handleRequests()
}

func (s *Server) handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/services", s.getAllServices).Methods("GET")
	myRouter.HandleFunc("/services/_count", s.countServices).Methods("GET")
	myRouter.HandleFunc("/services/{id}", s.createService).Methods("PUT")
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}
