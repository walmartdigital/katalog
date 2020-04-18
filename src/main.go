package main

import (
	"context"
	"errors"
	"flag"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"github.com/walmartdigital/katalog/src/domain"
	"github.com/walmartdigital/katalog/src/server"
	"github.com/walmartdigital/katalog/src/utils"

	"github.com/avast/retry-go"
	"github.com/gorilla/mux"
	k8sdriver "github.com/walmartdigital/katalog/src/collector/k8s-driver"
	"github.com/walmartdigital/katalog/src/collector/publishers"
	webhookServer "github.com/walmartdigital/katalog/src/server/http"
	kafkaServer "github.com/walmartdigital/katalog/src/server/kafka"
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
var publisher = flag.String("publisher", publisherHTTP, "select where to publish: kafka | http")
var configfile = flag.Bool("kubeconfig", false, "true if a $HOME/.kube/config file exists")

func main() {
	err := utils.LogInit(log)
	if err != nil {
		log.Fatal(err)
	}
	flag.Parse()
	var kubeconfig string

	if value, ok := os.LookupEnv("ROLE"); ok {
		switch value {
		case "SERVER":
			*role = roleServer
		case "COLLECTOR":
			*role = roleCollector
		default:
			panic("Role not supported")
		}
	}

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
		var wg sync.WaitGroup
		switch *publisher {
		case publisherHTTP:
			wg.Add(1)
			go mainServer(&wg, true)
		case publisherKafka:
			wg.Add(2)
			go mainServer(&wg, false)
			go mainConsumer(&wg, true)
		default:
			wg.Add(1)
			go mainServer(&wg, true)
		}
		wg.Wait()
	default:
		panic(errors.New("role should be server or collector"))
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
				log.Error(err)
			}
		case event := <-deploymentEvents:
			err := publisher.Publish(event)
			if err != nil {
				log.Error(err)
			}
		case event := <-statefulsetEvents:
			err := publisher.Publish(event)
			if err != nil {
				log.Error(err)
			}
		}
	}
}

var ticker *time.Ticker
var done chan bool

func check(checkable server.Checkable) {
	// Liveness probe
	ticker = time.NewTicker(30 * time.Second)
	done = make(chan bool)
	go func() {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			if checkable.Check() {
				log.Debug("(LIVE) Health check at " + t.Local().String())
				_, errOpen := os.OpenFile("/tmp/imalive", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
				if errOpen != nil {
					log.Error("Error opening health check file", errOpen)
				}

			} else {
				log.Debug("(DEAD) Health check at " + t.Local().String())
				errRemove := os.Remove("/tmp/imalive")
				if errRemove != nil {
					log.Error("Error removing health check file", errRemove)
				}
			}
		}
	}()
}

// KafkaWriterFactory ...
type KafkaWriterFactory struct{}

// Create ...
func (f KafkaWriterFactory) Create(kafkaURL string, topic string) publishers.Writer {
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{kafkaURL},
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})
}

func resolvePublisher() publishers.Publisher {
	var current publishers.Publisher
	switch *publisher {
	case publisherKafka:
		current = publishers.BuildKafkaPublisher(*kafkaURL, *kafkaTopicPrefix, KafkaWriterFactory{})
	case publisherHTTP:
		current = publishers.BuildHTTPPublisher(*httpURL, retry.Do)
	default:
		panic(errors.New("A publusher must be selected"))
	}

	check(current)

	return current
}

func closeProbes() {
	log.Info("Closing health checks")
	ticker.Stop()
	done <- true
}

func mainServer(wg *sync.WaitGroup, doCheck bool) {
	defer wg.Done()
	log.Info("http (webhook) server starting...")
	memory := new(sync.Map)
	persistence := persistence.BuildMemoryPersistence(memory)
	resourceRepository := repositories.CreateResourceRepository(persistence)
	router := mux.NewRouter().StrictSlash(true)
	routerWrapper := &routerWrapper{router: router}
	httpServer := &http.Server{Addr: ":10000", Handler: router}
	webhookServer := webhookServer.CreateServer(httpServer, resourceRepository, routerWrapper, PrometheusMetricsFactory{})
	if doCheck {
		check(webhookServer)
	}
	webhookServer.Run()
	log.Info("http (webhook) server started...")
}

// KafkaReaderFactory ...
type KafkaReaderFactory struct{}

// Create ...
func (f KafkaReaderFactory) Create(kafkaURL string, topic string) kafkaServer.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{kafkaURL},
		Topic:     topic,
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
	})
}

// ResourceRepositoryFactory ...
type ResourceRepositoryFactory struct {
	persistenceFactory MemoryPersistenceFactory
}

// Create ...
func (f ResourceRepositoryFactory) Create() repositories.Repository {
	persistence := f.persistenceFactory.Create()
	return repositories.CreateResourceRepository(persistence)
}

// MemoryPersistenceFactory ...
type MemoryPersistenceFactory struct{}

// Create ...
func (f MemoryPersistenceFactory) Create() persistence.Persistence {
	memory := new(sync.Map)
	return persistence.BuildMemoryPersistence(memory)
}

func mainConsumer(wg *sync.WaitGroup, doCheck bool) {
	consumerWg := new(sync.WaitGroup)

	defer wg.Done()

	log.Info("kafka consumer starting...")
	memFactory := MemoryPersistenceFactory{}
	service := server.MakeService(ResourceRepositoryFactory{persistenceFactory: memFactory}.Create(), PrometheusMetricsFactory{})

	created := kafkaServer.CreateConsumer(context.Background(), consumerWg, *kafkaURL, *kafkaTopicPrefix, "created", KafkaReaderFactory{}, &service)
	updated := kafkaServer.CreateConsumer(context.Background(), consumerWg, *kafkaURL, *kafkaTopicPrefix, "updated", KafkaReaderFactory{}, &service)
	deleted := kafkaServer.CreateConsumer(context.Background(), consumerWg, *kafkaURL, *kafkaTopicPrefix, "deleted", KafkaReaderFactory{}, &service)

	if doCheck {
		check(created)
		check(updated)
		check(deleted)
	}

	consumerWg.Add(3)
	go created.Run()
	go updated.Run()
	go deleted.Run()
	log.Info("kafka consumers started...")
	consumerWg.Wait()
}

type routerWrapper struct {
	router *mux.Router
}

func (r *routerWrapper) HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) webhookServer.Route {
	return &routeWrapper{route: r.router.HandleFunc(path, f)}
}

type routeWrapper struct {
	route *mux.Route
}

func (r *routeWrapper) Methods(methods ...string) webhookServer.Route {
	r.route.Methods(methods[0])
	return r
}

// PrometheusMetricsFactory ...
type PrometheusMetricsFactory struct {
}

// Create ...
func (p PrometheusMetricsFactory) Create() server.Metrics {
	var metrics server.Metrics
	metrics = server.PrometheusMetrics{}
	metrics.InitMetrics()
	return metrics
}
