package publishers_test

import (
	"encoding/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	kafka "github.com/segmentio/kafka-go"
	"golang.org/x/net/context"

	"github.com/walmartdigital/katalog/src/collector/publishers"
	"github.com/walmartdigital/katalog/src/domain"
	"github.com/walmartdigital/katalog/src/mocks/mock_publishers"
)

var _ = Describe("Run Consumer on 'created' topic", func() {
	var (
		fakeWriterFactory *mock_publishers.MockWriterFactory
		fakeWriter        *mock_publishers.MockWriter
		publisher         publishers.Publisher
		ctx               context.Context
		cancel            context.CancelFunc
	)

	BeforeEach(func() {
		// Initialize the mocked Kafka related objects
		fakeWriterFactory = mock_publishers.NewMockWriterFactory(ctrl)
		fakeWriter = mock_publishers.NewMockWriter(ctrl)
		fakeWriterFactory.EXPECT().Create("", ".created").Return(
			fakeWriter,
		).Times(1)
		fakeWriterFactory.EXPECT().Create("", ".deleted").Return(
			fakeWriter,
		).Times(1)
		fakeWriterFactory.EXPECT().Create("", ".updated").Return(
			fakeWriter,
		).Times(1)
		fakeWriterFactory.EXPECT().Create("", ".health").Return(
			fakeWriter,
		).Times(1)
		ctx, cancel = context.WithCancel(context.Background())
		_ = cancel

		publisher = publishers.BuildKafkaPublisher(ctx, "", "", fakeWriterFactory)
	})

	It("should create a publisher", func() {
		Expect(publisher).NotTo(BeNil())
	})

	It("should publish a Deployment creation event", func() {
		deployment := domain.Deployment{
			ID:         "276797fa-b207-11e9-8527-000d3af9d6b6",
			Name:       "queue-node",
			Generation: 7,
			Namespace:  "amida",
			Labels: map[string]string{
				"HEAD":                   "569de2ecd9f9357b3380664f43c90d07ec6acaff",
				"app":                    "nats",
				"fluxcd.io/sync-gc-mark": "sha256.0fRlq9kqkh2eSDRqXANMzgN8_8jeguja3eDLoE5E0Xo",
			},
			Containers: map[string]string{
				"nats-exporter":  "synadia/prometheus-nats-exporter:0.4.0",
				"nats-streaming": "nats-streaming:0.15.1",
			},
		}

		operation := domain.Operation{
			Kind:     domain.OperationTypeAdd,
			Resource: domain.Resource{K8sResource: &deployment},
		}

		dbytes, _ := json.Marshal(deployment)

		message := kafka.Message{
			Key:   []byte("/deployments/276797fa-b207-11e9-8527-000d3af9d6b6"),
			Value: dbytes,
		}

		fakeWriter.EXPECT().WriteMessages(ctx, message).Return(
			nil,
		).Times(1)
		publisher.Publish(operation)
	})

	It("should publish a StatefulSet creation event", func() {
		ss := domain.StatefulSet{
			ID:         "276797fa-b207-11e9-8527-000d3af9d6b6",
			Name:       "queue-node",
			Generation: 7,
			Namespace:  "amida",
			Labels: map[string]string{
				"HEAD":                   "569de2ecd9f9357b3380664f43c90d07ec6acaff",
				"app":                    "nats",
				"fluxcd.io/sync-gc-mark": "sha256.0fRlq9kqkh2eSDRqXANMzgN8_8jeguja3eDLoE5E0Xo",
			},
			Containers: map[string]string{
				"nats-exporter":  "synadia/prometheus-nats-exporter:0.4.0",
				"nats-streaming": "nats-streaming:0.15.1",
			},
		}

		operation := domain.Operation{
			Kind:     domain.OperationTypeAdd,
			Resource: domain.Resource{K8sResource: &ss},
		}

		ssbytes, _ := json.Marshal(ss)

		message := kafka.Message{
			Key:   []byte("/statefulsets/276797fa-b207-11e9-8527-000d3af9d6b6"),
			Value: ssbytes,
		}

		fakeWriter.EXPECT().WriteMessages(ctx, message).Return(
			nil,
		).Times(1)
		publisher.Publish(operation)
	})

	It("should publish a Service creation event", func() {
		i := domain.Instance{Address: "hello"}
		ss := domain.Service{
			ID:         "276797fa-b207-11e9-8527-000d3af9d6b6",
			Name:       "queue-node",
			Port:       1212,
			Address:    "someservice",
			Generation: 7,
			Namespace:  "amida",
			Instances:  []domain.Instance{i},
		}

		operation := domain.Operation{
			Kind:     domain.OperationTypeAdd,
			Resource: domain.Resource{K8sResource: &ss},
		}

		ssbytes, _ := json.Marshal(ss)

		message := kafka.Message{
			Key:   []byte("/services/276797fa-b207-11e9-8527-000d3af9d6b6"),
			Value: ssbytes,
		}

		fakeWriter.EXPECT().WriteMessages(ctx, message).Return(
			nil,
		).Times(1)
		publisher.Publish(operation)
	})

	It("should publish a Service update event", func() {
		i := domain.Instance{Address: "hello"}
		ss := domain.Service{
			ID:         "276797fa-b207-11e9-8527-000d3af9d6b6",
			Name:       "queue-node",
			Port:       1212,
			Address:    "someservice",
			Generation: 7,
			Namespace:  "amida",
			Instances:  []domain.Instance{i},
		}

		operation := domain.Operation{
			Kind:     domain.OperationTypeUpdate,
			Resource: domain.Resource{K8sResource: &ss},
		}

		ssbytes, _ := json.Marshal(ss)

		message := kafka.Message{
			Key:   []byte("/services/276797fa-b207-11e9-8527-000d3af9d6b6"),
			Value: ssbytes,
		}

		fakeWriter.EXPECT().WriteMessages(ctx, message).Return(
			nil,
		).Times(1)
		publisher.Publish(operation)
	})

	It("should publish a Service delete event", func() {
		i := domain.Instance{Address: "hello"}
		ss := domain.Service{
			ID:         "276797fa-b207-11e9-8527-000d3af9d6b6",
			Name:       "queue-node",
			Port:       1212,
			Address:    "someservice",
			Generation: 7,
			Namespace:  "amida",
			Instances:  []domain.Instance{i},
		}

		operation := domain.Operation{
			Kind:     domain.OperationTypeDelete,
			Resource: domain.Resource{K8sResource: &ss},
		}

		ssbytes, _ := json.Marshal(ss)

		message := kafka.Message{
			Key:   []byte("/services/276797fa-b207-11e9-8527-000d3af9d6b6"),
			Value: ssbytes,
		}

		fakeWriter.EXPECT().WriteMessages(ctx, message).Return(
			nil,
		).Times(1)
		publisher.Publish(operation)
	})

	AfterEach(func() {
	})
})
