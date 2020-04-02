package publishers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/sirupsen/logrus"
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
	url           string
	topicPrefix   string //katalog.artifact.[created|deleted|updated]
	kafkaWriters  map[string]*kafka.Writer
	healthCounter int
}

// BuildKafkaPublisher ...
func BuildKafkaPublisher(url string, topicPrefix string) Publisher {
	publisher := &KafkaPublisher{url: url, topicPrefix: topicPrefix}
	err := publisher.CreateProducers()
	if err != nil {
		logrus.Fatal(err)
	}
	return publisher
}

// CreateProducers ...
func (c *KafkaPublisher) CreateProducers() error {
	c.kafkaWriters = map[string]*kafka.Writer{
		"created": getKafkaWriter(c.url, c.topicPrefix+".created"),
		"deleted": getKafkaWriter(c.url, c.topicPrefix+".updated"),
		"updated": getKafkaWriter(c.url, c.topicPrefix+".updated"),
		"health":  getKafkaWriter(c.url, c.topicPrefix+".health"),
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

// getWriter ...
func (c *KafkaPublisher) getWriter(operation domain.Operation) *kafka.Writer {
	if c.kafkaWriters == nil {
		panic(errors.New("Writers not created, call GetProducers first"))
	}

	switch operation.Kind {
	case (domain.OperationTypeAdd):
		return c.kafkaWriters["created"]
	case (domain.OperationTypeUpdate):
		return c.kafkaWriters["updated"]
	case (domain.OperationTypeDelete):
		return c.kafkaWriters["deleted"]
	default:
		panic(errors.New("operation unknown"))
	}
}

// getPayload ...
func (c *KafkaPublisher) getPayload(resource domain.Resource) (string, error) {
	payloadBytes := new(bytes.Buffer)

	switch v := resource.GetType(); v {
	case reflect.TypeOf(new(domain.Service)):
		service := resource.GetK8sResource().(*domain.Service)
		err := json.NewEncoder(payloadBytes).Encode(*service)
		if err != nil {
			log.Fatal(err)
			return "", err
		}

	case reflect.TypeOf(new(domain.Deployment)):
		deployment := resource.GetK8sResource().(*domain.Deployment)
		err := json.NewEncoder(payloadBytes).Encode(*deployment)
		if err != nil {
			log.Fatal(err)
			return "", err
		}

	case reflect.TypeOf(new(domain.StatefulSet)):
		statefulset := resource.GetK8sResource().(*domain.StatefulSet)
		err := json.NewEncoder(payloadBytes).Encode(*statefulset)
		if err != nil {
			log.Fatal(err)
			return "", err
		}

	default:
		log.Errorf("Type %s not found", v)
		panic(errors.New("Type %s not found"))
	}

	return payloadBytes.String(), nil
}

// getKey ...
func (c *KafkaPublisher) getKey(resource domain.Resource) string {
	switch v := resource.GetType(); v {
	case reflect.TypeOf(new(domain.Service)):
		service := resource.GetK8sResource().(*domain.Service)
		return "/services/" + service.ID

	case reflect.TypeOf(new(domain.Deployment)):
		deployment := resource.GetK8sResource().(*domain.Deployment)
		return "/deployments/" + deployment.ID

	case reflect.TypeOf(new(domain.StatefulSet)):
		statefulset := resource.GetK8sResource().(*domain.StatefulSet)
		return "/statefulsets/" + statefulset.ID

	default:
		log.Errorf("Type %s not found", v)
		panic(errors.New("Type %s not found"))
	}
}

// Check ...
func (c *KafkaPublisher) Check() bool {
	if c.kafkaWriters == nil {
		return false
	}

	writer := c.kafkaWriters["health"]

	c.healthCounter++
	err := writer.WriteMessages(
		context.Background(),
		kafka.Message{
			Key:   []byte("check"),
			Value: []byte(fmt.Sprintf("{\"count\": %d}", c.healthCounter)),
		},
	)
	if err != nil {
		log.Fatal(err)
		return false
	}

	return true
}

// Publish ...
func (c *KafkaPublisher) Publish(obj interface{}) error {
	operation := obj.(domain.Operation)

	writer := c.getWriter(operation)

	key := c.getKey(operation.Resource)

	value, errGettingValue := c.getPayload(operation.Resource)
	if errGettingValue != nil {
		log.Fatal(errGettingValue)
		return errGettingValue
	}

	log.WithFields(logrus.Fields{
		"key": key,
	}).Debug("Sending message")

	errWritingMessage := writer.WriteMessages(
		context.Background(),
		kafka.Message{
			Key:   []byte(key),
			Value: []byte(value),
		},
	)
	if errWritingMessage != nil {
		log.Fatal(errWritingMessage)
		return errWritingMessage
	}

	return nil
}
