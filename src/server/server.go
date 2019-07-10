package server

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/walmartdigital/katalog/src/server/repositories"
)

// HTTPServer ...
type HTTPServer interface {
	ListenAndServe() error
}

// Server ...
type Server struct {
	httpServer          HTTPServer
	resourcesRepository repositories.Repository
	router              Router
	metrics             *map[string]interface{}
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
func CreateServer(server HTTPServer, repository repositories.Repository, router Router) Server {
	metricsmap := make(map[string]interface{})

	metricsmap["createDeployment"] = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "katalog",
			Subsystem: "deployment",
			Name:      "katalog_deployment_create",
			Help:      "Total number of deployment creations",
		},
		[]string{"id"},
	)
	prometheus.MustRegister(metricsmap["createDeployment"].(*prometheus.CounterVec))

	metricsmap["updateDeployment"] = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "katalog",
			Subsystem: "deployment",
			Name:      "katalog_deployment_update",
			Help:      "Total number of deployment updates",
		},
		[]string{"id"},
	)
	prometheus.MustRegister(metricsmap["updateDeployment"].(*prometheus.CounterVec))

	return Server{
		httpServer:          server,
		resourcesRepository: repository,
		router:              router,
		metrics:             &metricsmap,
	}
}

// Run ...
func (s *Server) Run() {
	s.handleRequests()
}

func (s *Server) handleRequests() {
	s.router.HandleFunc("/metrics", promhttp.Handler().ServeHTTP).Methods("GET")
	s.router.HandleFunc("/services", s.getAllServices).Methods("GET")
	s.router.HandleFunc("/services/_count", s.countServices).Methods("GET")
	s.router.HandleFunc("/services/{id}", s.createService).Methods("POST")
	s.router.HandleFunc("/services/{id}", s.updateService).Methods("PUT")
	s.router.HandleFunc("/services/{id}", s.deleteService).Methods("DELETE")
	s.router.HandleFunc("/deployments", s.getAllDeployments).Methods("GET")
	s.router.HandleFunc("/deployments/_count", s.countDeployments).Methods("GET")
	s.router.HandleFunc("/deployments/{id}", s.createDeployment).Methods("POST")
	s.router.HandleFunc("/deployments/{id}", s.updateDeployment).Methods("PUT")
	s.router.HandleFunc("/deployments/{id}", s.deleteDeployment).Methods("DELETE")
	s.httpServer.ListenAndServe()
}
