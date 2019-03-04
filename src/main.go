package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/golang/glog"
	"github.com/seadiaz/katalog/src/k8s-driver"
	"github.com/seadiaz/katalog/src/persistence"
	"github.com/seadiaz/katalog/src/publishers"
	"github.com/seadiaz/katalog/src/server"
)

const roleCollector = "collector"
const roleServer = "server"

const publisherHTTP = "http"
const publisherConsul = "consul"

var role = flag.String("role", roleCollector, "collector or server")
var consulAddress = flag.String("consul-addr", "127.0.0.1:8500", "consul address")
var httpUlr = flag.String("http-url", "http://127.0.0.1:10000", "http url")
var excludeSysmteNamespace = flag.Bool("exclude-sysmte-namespace", false, "exclude all services from kube-system namespace")
var publisher = flag.String("publisher", publisherHTTP, "select where to publis: http, consul")

func main() {
	flag.Parse()

	kubeconfig := filepath.Join(
		os.Getenv("HOME"), ".kube", "config",
	)

	switch *role {
	case roleCollector:
		mainCollector(kubeconfig)
	case roleServer:
		mainServer()
	default:
		glog.Warning("role not found")
	}
}

func mainCollector(kubeconfig string) {
	glog.Info("collector starting...")
	serviceEvents := make(chan interface{})
	k8sDriver := k8sdriver.BuildDriver(kubeconfig, *excludeSysmteNamespace)
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
		return publishers.CreateHTTPPublisher(*httpUlr)
	case publisherConsul:
		return publishers.CreateConsulPublisher(*consulAddress)
	default:
		return nil
	}
}

func mainServer() {
	persistence := persistence.CreateBoltDriver()
	server := server.CreateServer(persistence)
	server.Run()
}