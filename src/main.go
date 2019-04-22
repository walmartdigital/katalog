package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	k8sdriver "github.com/walmartdigital/katalog/src/collector/k8s-driver"
	"github.com/walmartdigital/katalog/src/collector/publishers"
	"github.com/walmartdigital/katalog/src/server"
	"github.com/walmartdigital/katalog/src/server/persistence"
	"github.com/walmartdigital/katalog/src/server/repositories"
)

const roleCollector = "collector"
const roleServer = "server"

const publisherHTTP = "http"

var role = flag.String("role", roleCollector, "collector or server")
var consulAddress = flag.String("consul-addr", "127.0.0.1:8500", "consul address")
var httpURL = flag.String("http-url", "http://127.0.0.1:10000", "http url")
var excludeSystemNamespace = flag.Bool("exclude-system-namespace", false, "exclude all services from kube-system namespace")
var publisher = flag.String("publisher", publisherHTTP, "select where to publis: http, consul")

func usage() {
	flag.PrintDefaults()
	os.Exit(2)
}

func init() {
	flag.Usage = usage
	flag.Set("logtostderr", "true")
	flag.Set("stderrthreshold", "WARNING")
	flag.Set("v", "2")
	flag.Parse()
}

func main() {
	kubeconfig := filepath.Join(
		os.Getenv("HOME"), ".kube", "config",
	)

	switch *role {
	case roleCollector:
		mainCollector(kubeconfig)
	case roleServer:
		mainServer()
	default:
		glog.Warning("role should be server or collector")
	}
}

func mainCollector(kubeconfig string) {
	glog.Info("collector starting...")
	serviceEvents := make(chan interface{})
	k8sDriver := k8sdriver.BuildDriver(kubeconfig, *excludeSystemNamespace)
	publisher := resolvePublisher()
	go k8sDriver.StartWatchingServices(serviceEvents)
	for {
		select {
		case event := <-serviceEvents:
			publisher.Publish(event)
		}
	}
}

func resolvePublisher() publishers.Publisher {
	switch *publisher {
	case publisherHTTP:
		return publishers.BuildHTTPPublisher(*httpURL)
	default:
		return nil
	}
}

func mainServer() {
	glog.Info("server starting...")
	memory := make(map[string]interface{})
	persistence := persistence.BuildMemoryPersistence(memory)
	serviceRepository := repositories.CreateServiceRepository(persistence)
	router := mux.NewRouter().StrictSlash(true)
	server := server.CreateServer(serviceRepository, router)
	server.Run()
}
