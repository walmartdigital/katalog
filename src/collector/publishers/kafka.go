package publishers

import (
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/walmartdigital/katalog/src/domain"
	"github.com/walmartdigital/katalog/src/utils"

	kafka "github.com/segmentio/kafka-go"
)

var log = logrus.New()

func init() {
	err := utils.LogInit(log)
	if err != nil {
		logrus.Fatal(err)
	}
}

func getKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{kafkaURL},
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})
}

// KafkaPublisher ...
type KafkaPublisher struct {
	url   string
	topicPrefix string //katalog.artifact.[created|deleted|updated]
	kafkaWriters map[string]*kafka.writer
}

// BuildKafkaPublisher ...
func BuildKafkaPublisher(url string, topic string) Publisher {
	return &KafkaPublisher{url: url, topic: topic}
}

// CreateProducers ...
func (c *KafkaPublisher) CreateProducers() error {
	c.kafkaWriters := map[string]*kafka.writer{
		"created": getKafkaWriter(c.url, c.topicPrefix + "created"),
		"deleted": getKafkaWriter(c.url, c.topicPrefix + "updated"),
		"updated": getKafkaWriter(c.url, c.topicPrefix + "updated"),	
	}
}

// Close ...
func (c *KafkaPublisher) Close() error {
	c.kafkaWriter.Close()
}

// Publish ...
func (c *KafkaPublisher) Publish(obj interface{}) error {
	operation := obj.(domain.Operation)
	

}
