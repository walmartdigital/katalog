package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/golang/glog"
	"github.com/seadiaz/katalog/k8s-driver"
	"github.com/seadiaz/katalog/publishers"
)

const roleCollector = "collector"

var role = flag.String("role", "collector", "collector or server")
var consulAddress = flag.String("consul-addr", "127.0.0.1:8500", "consul address")
var excludeSysmteNamespace = flag.Bool("exclude-sysmte-namespace", false, "exclude all services from kube-system namespace")

func main() {
	flag.Parse()

	kubeconfig := filepath.Join(
		os.Getenv("HOME"), ".kube", "config",
	)

	switch *role {
	case roleCollector:
		mainCollector(kubeconfig)
	default:
		glog.Warning("role not found")
	}
}

func mainCollector(kubeconfig string) {
	serviceEvents := make(chan interface{})
	k8sDriver := k8sdriver.BuildDriver(kubeconfig, *excludeSysmteNamespace)
	publisher := publishers.Create(*consulAddress)
	go k8sDriver.StartWatchingServices(serviceEvents)
	for {
		select {
		case event := <-serviceEvents:
			publisher.Publish(event)
		}
	}
}
