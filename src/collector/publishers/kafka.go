package publishers

import (
	"github.com/walmartdigital/katalog/src/domain"

	kafka "github.com/segmentio/kafka-go"
)

func getKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{kafkaURL},
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})
}

// KafkaPublisher ...
type KafkaPublisher struct {
	url          string
	topicPrefix  string //katalog.artifact.[created|deleted|updated]
	kafkaWriters map[string]*kafka.Writer
}

// BuildKafkaPublisher ...
func BuildKafkaPublisher(url string, topicPrefix string) Publisher {
	return &KafkaPublisher{url: url, topicPrefix: topicPrefix}
}

// CreateProducers ...
func (c *KafkaPublisher) CreateProducers() error {
	c.kafkaWriters = map[string]*kafka.Writer{
		"created": getKafkaWriter(c.url, c.topicPrefix+".created"),
		"deleted": getKafkaWriter(c.url, c.topicPrefix+".updated"),
		"updated": getKafkaWriter(c.url, c.topicPrefix+".updated"),
	}

	return nil
}

// Close ...
func (c *KafkaPublisher) Close() error {
	var err error
	errCreated := c.kafkaWriters["created"].Close()
	if errCreated != nil {
		log.Fatal(errCreated)
		err = errCreated
	}

	errDeleted := c.kafkaWriters["deleted"].Close()
	if errDeleted != nil {
		log.Fatal(errDeleted)
		err = errDeleted
	}

	errUpdated := c.kafkaWriters["updated"].Close()
	if errUpdated != nil {
		log.Fatal(errUpdated)
		err = errUpdated
	}

	return err
}

// Publish ...
func (c *KafkaPublisher) Publish(obj interface{}) error {
	operation := obj.(domain.Operation)

	_ = operation

	return nil
}
