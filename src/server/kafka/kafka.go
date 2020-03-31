package kafka

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/walmartdigital/katalog/src/regex"
	"github.com/walmartdigital/katalog/src/server"
	"github.com/walmartdigital/katalog/src/server/repositories"
	"github.com/walmartdigital/katalog/src/utils"

	kafka "github.com/segmentio/kafka-go"
)

var log = logrus.New()

func init() {
	err := utils.LogInit(log)
	if err != nil {
		log.Fatal(err)
	}
}

func getKafkaReader(kafkaURL string, topic string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{kafkaURL},
		Topic:     topic,
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
	})
}

// Consumer ...
type Consumer struct {
	url                 string
	topicPrefix         string //katalog.artifact.[created|deleted|updated]
	KafkaReaders        map[string]*kafka.Reader
	resourcesRepository repositories.Repository
	metrics             *map[string]interface{}
	service             server.Service
}

// CreateConsumer ...
func CreateConsumer(kafkaURL string, topicPrefix string, repository repositories.Repository) Consumer {
	current := Consumer{
		url:                 kafkaURL,
		topicPrefix:         topicPrefix,
		resourcesRepository: repository,
		KafkaReaders: map[string]*kafka.Reader{
			"created": getKafkaReader(kafkaURL, topicPrefix+".created"),
			"deleted": getKafkaReader(kafkaURL, topicPrefix+".updated"),
			"updated": getKafkaReader(kafkaURL, topicPrefix+".updated"),
		},
		metrics: server.InitMetrics(),
	}

	current.service = server.MakeService(current.resourcesRepository, current.metrics)

	return current
}

// Run ...
func (c *Consumer) Run() {
	go c.ConsumeEvent("created")
	go c.ConsumeEvent("deleted")
	go c.ConsumeEvent("updated")
}

// ConsumeEvent ...
func (c *Consumer) ConsumeEvent(event string) {
	consumer := c.KafkaReaders[event]

	defer consumer.Close()

	for {
		m, err := consumer.ReadMessage(context.Background())
		if err != nil {
			break
		}

		key := string(m.Key)
		value := string(m.Value)

		log.Debug("message at offset %d: %s = %s\n", m.Offset, key, value)

		matchedNamedGroups := regex.GetParams(
			"/(?P<artifact>.+)/(?P<id>.+)",
			key,
		)

		artifact := matchedNamedGroups["artifact"]
		id := matchedNamedGroups["id"]

		switch event {
		case "create":
			switch artifact {
			case "service":
				c.CreateService(value)
			case "deployment":
				c.CreateDeployment(value)
			case "statefulset":
				c.CreateStatefulSet(value)
			}
		case "update":
			switch "artifact" {
			case "service":
				c.UpdateService(value)
			case "deployment":
				c.UpdateDeployment(value)
			case "statefulset":
				c.UpdateStatefulSet(value)
			}
		case "delete":
			switch "artifact" {
			case "service":
				c.DeleteService(id)
			case "deployment":
				c.DeleteDeployment(id)
			case "statefulset":
				c.DeleteStatefulSet(id)
			}
		}
	}
}

// DestroyMetrics ...
func (c *Consumer) DestroyMetrics() {
	for _, v := range *(c.metrics) {
		prometheus.Unregister(v.(prometheus.Collector))
	}
}
