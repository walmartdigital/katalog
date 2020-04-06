package http

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/walmartdigital/katalog/src/server"
	"github.com/walmartdigital/katalog/src/server/repositories"
	"github.com/walmartdigital/katalog/src/utils"
)

var log = logrus.New()

func init() {
	err := utils.LogInit(log)
	if err != nil {
		log.Fatal(err)
	}
}

// WebhookServer ...
type WebhookServer interface {
	ListenAndServe() error
}

// Server ...
type Server struct {
	httpServer          WebhookServer
	resourcesRepository repositories.Repository
	router              Router
	metrics             *map[string]interface{}
	service             server.Service
}

// Check ...
func (s *Server) Check() bool {
	return true
}

// Router ...
type Router interface {
	HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) Route
}

// Route ...
type Route interface {
	Methods(methods ...string) Route
}

// CreateServer ...
func CreateServer(webhook WebhookServer, repository repositories.Repository, router Router) *Server {
	current := &Server{
		httpServer:          webhook,
		resourcesRepository: repository,
		router:              router,
	}

	current.service = server.MakeService(current.resourcesRepository)

	return current
}

// Run ...
func (s *Server) Run() {
	s.handleRequests()
}

func (s *Server) handleRequests() {
	s.router.HandleFunc("/metrics", promhttp.Handler().ServeHTTP).Methods("GET")
	s.router.HandleFunc("/services", s.getAllServices).Methods("GET")
	s.router.HandleFunc("/services/_count", s.countServices).Methods("GET")
	s.router.HandleFunc("/services/{id}", s.CreateService).Methods("POST")
	s.router.HandleFunc("/services/{id}", s.UpdateService).Methods("PUT")
	s.router.HandleFunc("/services/{id}", s.DeleteService).Methods("DELETE")
	s.router.HandleFunc("/deployments", s.getAllDeployments).Methods("GET")
	s.router.HandleFunc("/deployments/_count", s.countDeployments).Methods("GET")
	s.router.HandleFunc("/deployments/{id}", s.CreateDeployment).Methods("POST")
	s.router.HandleFunc("/deployments/{id}", s.UpdateDeployment).Methods("PUT")
	s.router.HandleFunc("/deployments/{id}", s.DeleteDeployment).Methods("DELETE")
	s.router.HandleFunc("/statefulsets", s.getAllStatefulSets).Methods("GET")
	s.router.HandleFunc("/statefulsets/_count", s.countStatefulSets).Methods("GET")
	s.router.HandleFunc("/statefulsets/{id}", s.CreateStatefulSet).Methods("POST")
	s.router.HandleFunc("/statefulsets/{id}", s.UpdateStatefulSet).Methods("PUT")
	s.router.HandleFunc("/statefulsets/{id}", s.DeleteStatefulSet).Methods("DELETE")

	err := s.httpServer.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
