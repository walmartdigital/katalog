package main

import (
	"errors"
	"flag"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/walmartdigital/katalog/src/domain"
	"github.com/walmartdigital/katalog/src/utils"

	"github.com/avast/retry-go"
	"github.com/gorilla/mux"
	k8sdriver "github.com/walmartdigital/katalog/src/collector/k8s-driver"
	"github.com/walmartdigital/katalog/src/collector/publishers"
	"github.com/walmartdigital/katalog/src/server"
	"github.com/walmartdigital/katalog/src/server/persistence"
	"github.com/walmartdigital/katalog/src/server/repositories"
)

var log = logrus.New()

const roleCollector = "collector"
const roleServer = "server"
const publisherHTTP = "http"
const publisherKafka = "kafka"

var role = flag.String("role", roleCollector, "collector or server")
var httpURL = flag.String("http-url", "http://127.0.0.1:10000", "http url")
var kafkaURL = flag.String("kafka-url", "localhost:9092", "kafka url")
var kafkaTopicPrefix = flag.String("kafka-topic-prefix", "_katalog.artifact", "kafka topic prefix")
var excludeSystemNamespace = flag.Bool("exclude-system-namespace", false, "exclude all services from kube-system namespace")
var publisher = flag.String("publisher", publisherKafka, "select where to publish: kafka | http")
var configfile = flag.Bool("kubeconfig", false, "true if a $HOME/.kube/config file exists")

func main() {
	err := utils.LogInit(log)
	if err != nil {
		log.Fatal(err)
	}
	flag.Parse()
	var kubeconfig string

	if value, ok := os.LookupEnv("PUBLISHER"); ok {
		publisher = &value
	}

	if value, ok := os.LookupEnv("HTTP_URL"); ok {
		httpURL = &value
	}

	if value, ok := os.LookupEnv("KAFKA_URL"); ok {
		kafkaURL = &value
	}

	if value, ok := os.LookupEnv("KAFKA_TOPIC_PREFIX"); ok {
		kafkaTopicPrefix = &value
	}

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
		log.Warning("role should be server or collector")
	}
}

func mainCollector(kubeconfig string) {
	log.Info("collector starting...")
	serviceEvents := make(chan interface{})
	deploymentEvents := make(chan interface{})
	statefulsetEvents := make(chan interface{})
	k8sDriver := k8sdriver.BuildDriver(kubeconfig, *excludeSystemNamespace)
	publisher := resolvePublisher()
	defer closeProbes()
	go k8sDriver.StartWatchingResources(serviceEvents, domain.Resource{K8sResource: &domain.Service{}})
	go k8sDriver.StartWatchingResources(deploymentEvents, domain.Resource{K8sResource: &domain.Deployment{}})
	go k8sDriver.StartWatchingResources(statefulsetEvents, domain.Resource{K8sResource: &domain.StatefulSet{}})
	for {
		select {
		case event := <-serviceEvents:
			err := publisher.Publish(event)
			if err != nil {
				log.Fatal(err)
			}
		case event := <-deploymentEvents:
			err := publisher.Publish(event)
			if err != nil {
				log.Fatal(err)
			}
		case event := <-statefulsetEvents:
			err := publisher.Publish(event)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

var ticker *time.Ticker
var done chan bool

func resolvePublisher() publishers.Publisher {
	var current publishers.Publisher
	switch *publisher {
	case publisherKafka:
		current = publishers.BuildKafkaPublisher(*kafkaURL, *kafkaTopicPrefix)
	case publisherHTTP:
		current = publishers.BuildHTTPPublisher(*httpURL, retry.Do)
	default:
		panic(errors.New("A publusher must be selected"))
	}

	// Liveness probe
	ticker = time.NewTicker(30 * time.Second)
	done = make(chan bool)
	go func() {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			if current.Check() {
				log.Debug("(LIVE) Health check at " + t.Local().String())
				_, errOpen := os.OpenFile("/tmp/imlive", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
				if errOpen != nil {
					log.Fatal(errOpen)
				}

			} else {
				log.Debug("(DEAD) Health check at " + t.Local().String())
				errRemove := os.Remove("/tmp/imlive")
				if errRemove != nil {
					log.Fatal(errRemove)
				}
			}
		}
	}()

	return current
}

func closeProbes() {
	log.Info("Closing health checks")
	ticker.Stop()
	done <- true
}

func mainServer() {
	log.Info("server starting...")
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
