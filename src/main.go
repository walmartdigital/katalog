package main

import (
	"flag"
	"net/http"
	"os"
	"path/filepath"

	"github.com/walmartdigital/katalog/src/domain"

	"github.com/avast/retry-go"
	"github.com/gorilla/mux"
	k8sdriver "github.com/walmartdigital/katalog/src/collector/k8s-driver"
	"github.com/walmartdigital/katalog/src/collector/publishers"
	"github.com/walmartdigital/katalog/src/server"
	"github.com/walmartdigital/katalog/src/server/persistence"
	"github.com/walmartdigital/katalog/src/server/repositories"
	"k8s.io/klog"
)

const roleCollector = "collector"
const roleServer = "server"

const publisherHTTP = "http"

var role = flag.String("role", roleCollector, "collector or server")
var httpURL = flag.String("http-url", "http://127.0.0.1:10000", "http url")
var excludeSystemNamespace = flag.Bool("exclude-system-namespace", false, "exclude all services from kube-system namespace")
var publisher = flag.String("publisher", publisherHTTP, "select where to publish: http")
var configfile = flag.Bool("kubeconfig", false, "true if a $HOME/.kube/config file exists")

func main() {
	klog.InitFlags(nil)
	flag.Set("logtostderr", "true")
	flag.Parse()
	var kubeconfig string

	if *configfile {
		kubeconfig = filepath.Join(
			os.Getenv("HOME"), ".kube", "config",
		)
	} else {
		kubeconfig = ""
	}

	switch *role {
	case roleCollector:
		mainCollector(kubeconfig)
	case roleServer:
		mainServer()
	default:
		klog.Warning("role should be server or collector")
	}
}

func mainCollector(kubeconfig string) {
	klog.Info("collector starting...")
	serviceEvents := make(chan interface{})
	deploymentEvents := make(chan interface{})
	statefulsetEvents := make(chan interface{})
	k8sDriver := k8sdriver.BuildDriver(kubeconfig, *excludeSystemNamespace)
	publisher := resolvePublisher()
	go k8sDriver.StartWatchingResources(serviceEvents, domain.Resource{K8sResource: &domain.Service{}})
	go k8sDriver.StartWatchingResources(deploymentEvents, domain.Resource{K8sResource: &domain.Deployment{}})
	go k8sDriver.StartWatchingResources(statefulsetEvents, domain.Resource{K8sResource: &domain.StatefulSet{}})
	for {
		select {
		case event := <-serviceEvents:
			publisher.Publish(event)
		case event := <-deploymentEvents:
			publisher.Publish(event)
		case event := <-statefulsetEvents:
			publisher.Publish(event)
		}
	}
}

func resolvePublisher() publishers.Publisher {
	switch *publisher {
	case publisherHTTP:
		return publishers.BuildHTTPPublisher(*httpURL, retry.Do)
	default:
		return nil
	}
}

func mainServer() {
	klog.Info("server starting...")
	memory := make(map[string]interface{})
	persistence := persistence.BuildMemoryPersistence(memory)
	resourceRepository := repositories.CreateResourceRepository(persistence)
	router := mux.NewRouter().StrictSlash(true)
	routerWrapper := &routerWrapper{router: router}
	httpServer := &http.Server{Addr: ":10000", Handler: router}
	server := server.CreateServer(httpServer, resourceRepository, routerWrapper)

	server.Run()
}

type routerWrapper struct {
	router *mux.Router
}

func (r *routerWrapper) HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) server.Route {
	return &routeWrapper{route: r.router.HandleFunc(path, f)}
}

type routeWrapper struct {
	route *mux.Route
}

func (r *routeWrapper) Methods(methods ...string) server.Route {
	r.route.Methods(methods[0])
	return r
}
