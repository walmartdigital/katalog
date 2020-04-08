package kafka

import (
	"context"
	"sync"

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
	url                 string
	topicPrefix         string //katalog.artifact.[created|deleted|updated]
	KafkaReaders        map[string]*Reader
	resourcesRepository repositories.Repository
	service             server.Service
}

// CreateConsumer ...
func CreateConsumer(kafkaURL string, topicPrefix string, rfactory ReaderFactory, repository repositories.Repository, mfactory server.MetricsFactory) *Consumer {
	created := rfactory.Create(kafkaURL, topicPrefix+".created")
	deleted := rfactory.Create(kafkaURL, topicPrefix+".deleted")
	updated := rfactory.Create(kafkaURL, topicPrefix+".updated")

	current := &Consumer{
		url:                 kafkaURL,
		topicPrefix:         topicPrefix,
		resourcesRepository: repository,
		KafkaReaders: map[string]*Reader{
			"created": &created,
			"deleted": &deleted,
			"updated": &updated,
		},
	}

	current.service = server.MakeService(current.resourcesRepository, mfactory)

	return current
}

var wg sync.WaitGroup

// Run ...
func (c *Consumer) Run() {
	wg.Add(3)
	go c.ConsumeEvent("created")
	go c.ConsumeEvent("deleted")
	go c.ConsumeEvent("updated")
	wg.Wait()
}

// ConsumeEvent ...
func (c *Consumer) ConsumeEvent(event string) {
	defer wg.Done()

	consumer := c.KafkaReaders[event]

	defer (*consumer).Close()

	for {
		m, err := (*consumer).ReadMessage(context.Background())
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
			"event":    event,
			"artifact": artifact,
			"id":       id,
		}).Debug("Event processing")

		switch event {
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
					"event":    event,
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
					"event":    event,
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
					"event":    event,
					"artifact": artifact,
				}).Warn("Artifact not recognized")
			}
		default:
			log.WithFields(logrus.Fields{
				"event":    event,
				"artifact": artifact,
			}).Warn("Event not recognized")
		}

		log.WithFields(logrus.Fields{
			"event":    event,
			"artifact": artifact,
			"id":       id,
		}).Debug("Event process task launched")
	}
}
