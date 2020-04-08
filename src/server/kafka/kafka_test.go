package kafka_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/walmartdigital/katalog/src/mocks/mock_kafka"
	"github.com/walmartdigital/katalog/src/mocks/mock_repositories"
	"github.com/walmartdigital/katalog/src/mocks/mock_server"
	"github.com/walmartdigital/katalog/src/server/kafka"
)

var ctrl *gomock.Controller

func TestAll(t *testing.T) {
	ctrl = gomock.NewController(t)
	defer ctrl.Finish()

	RegisterFailHandler(Fail)
	RunSpecs(t, "Server")
}

var _ = Describe("run server", func() {
	var ()

	BeforeEach(func() {
		// Initialize the mocked Kafka related objects
		fakeReaderFactory := mock_kafka.NewMockReaderFactory(ctrl)
		fakeReader := mock_kafka.NewMockReader(ctrl)
		fakeReaderFactory.EXPECT().Create(gomock.Any(), gomock.Any()).Return(
			fakeReader,
		).Times(3)

		// Initialize the mocked Repository related objects
		fakeRepoFactory := mock_repositories.NewMockRepositoryFactory(ctrl)
		fakeRepo := mock_repositories.NewMockRepository(ctrl)
		fakeRepoFactory.EXPECT().Create().Return(
			fakeRepo,
		).Times(1)

		// Initialize the mocked Metrics related objects
		fakeMetricsFactory := mock_server.NewMockMetricsFactory(ctrl)
		fakeMetrics := mock_server.NewMockMetrics(ctrl)
		fakeMetrics.EXPECT().IncrementCounter(gomock.Any(), gomock.Any()).AnyTimes()
		fakeMetricsFactory.EXPECT().Create().Return(
			fakeMetrics,
		).Times(1)

		consumer := kafka.CreateConsumer("", "", fakeReaderFactory, fakeRepoFactory, fakeMetricsFactory)
		_ = consumer
	})

	It("should create a service", func() {
		Expect("hello").To(Equal("hello"))
	})

	AfterEach(func() {
	})
})
