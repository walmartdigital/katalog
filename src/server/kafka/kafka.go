package kafka

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/walmartdigital/katalog/src/regex"
	"github.com/walmartdigital/katalog/src/server"
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

// Check ...
func (c *Consumer) Check() bool {
	return true
}

// Reader ...
type Reader interface {
	Close() error
	ReadMessage(context.Context) (kafka.Message, error)
}

// ReaderFactory ...
type ReaderFactory interface {
	Create(string, string) Reader
}

// Consumer ...
type Consumer struct {
	url         string
	event       string
	topicPrefix string //katalog.artifact.[created|deleted|updated]
	reader      Reader
	context     context.Context
	wg          *sync.WaitGroup
	service     *server.Service
}

// CreateConsumer ...
func CreateConsumer(ctx context.Context, wg *sync.WaitGroup, kafkaURL string, topicPrefix string, event string, readerFactory ReaderFactory, service *server.Service) *Consumer {
	return &Consumer{
		url:         kafkaURL,
		topicPrefix: topicPrefix,
		reader:      readerFactory.Create(kafkaURL, topicPrefix+"."+event),
		context:     ctx,
		wg:          wg,
		event:       event,
		service:     service,
	}
}

// Run ...
func (c *Consumer) Run() {
	defer c.wg.Done()
	defer c.reader.Close()

	for {
		select {
		case <-c.context.Done():
			log.Info("Received cancel signal from parent context")
			return
		default:
			m, err := c.reader.ReadMessage(c.context)
			if err != nil {
				break
			}

			key := string(m.Key)
			value := string(m.Value)

			log.WithFields(logrus.Fields{
				"key":    key,
				"offset": m.Offset,
			}).Debug("Message Received")

			matchedNamedGroups := regex.GetParams(
				"/(?P<artifact>.+)/(?P<id>.+)",
				key,
			)

			artifact := matchedNamedGroups["artifact"]
			id := matchedNamedGroups["id"]

			log.WithFields(logrus.Fields{
				"event":    c.event,
				"artifact": artifact,
				"id":       id,
			}).Debug("Event processing")

			switch c.event {
			case "created":
				switch artifact {
				case "services":
					go c.CreateService(value)
				case "deployments":
					go c.CreateDeployment(value)
				case "statefulsets":
					go c.CreateStatefulSet(value)
				default:
					log.WithFields(logrus.Fields{
						"event":    c.event,
						"artifact": artifact,
					}).Warn("Artifact not recognized")
				}
			case "updated":
				switch artifact {
				case "services":
					go c.UpdateService(value)
				case "deployments":
					go c.UpdateDeployment(value)
				case "statefulsets":
					go c.UpdateStatefulSet(value)
				default:
					log.WithFields(logrus.Fields{
						"event":    c.event,
						"artifact": artifact,
					}).Warn("Artifact not recognized")
				}
			case "deleted":
				switch artifact {
				case "services":
					go c.DeleteService(id)
				case "deployments":
					go c.DeleteDeployment(id)
				case "statefulsets":
					go c.DeleteStatefulSet(id)
				default:
					log.WithFields(logrus.Fields{
						"event":    c.event,
						"artifact": artifact,
					}).Warn("Artifact not recognized")
				}
			default:
				log.WithFields(logrus.Fields{
					"event":    c.event,
					"artifact": artifact,
				}).Warn("Event not recognized")
			}

			log.WithFields(logrus.Fields{
				"event":    c.event,
				"artifact": artifact,
				"id":       id,
			}).Debug("Event process task launched")
		}
	}
}
