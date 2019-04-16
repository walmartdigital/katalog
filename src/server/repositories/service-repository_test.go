package repositories_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/walmartdigital/katalog/src/domain"
	"github.com/walmartdigital/katalog/src/server/repositories"
)

type fakePersistence struct {
	memory map[string]interface{}
}

func (f *fakePersistence) Create(kind string, id string, obj interface{}) {
	f.memory[kind] = obj
}

func (f *fakePersistence) Delete(kind string, id string) {
	delete(f.memory, kind)
}

func (f *fakePersistence) GetAll(kind string) []interface{} {
	s := make([]interface{}, 1)
	s[0] = f.memory[kind]
	return s
}

func TestAll(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "service repository")
}

var _ = Describe("create service", func() {
	It("should persist a service", func() {
		service := domain.Service{}
		memory := make(map[string]interface{})
		fake := fakePersistence{memory: memory}
		serviceRepository := repositories.CreateServiceRepository(&fake)

		serviceRepository.CreateService(service)

		Expect(fake.GetAll("services")[0]).To(Equal(service))
	})
})

var _ = Describe("delete service", func() {
	It("should delete a given service", func() {
		service := domain.Service{}
		memory := make(map[string]interface{})
		fake := fakePersistence{memory: memory}
		serviceRepository := repositories.CreateServiceRepository(&fake)
		serviceRepository.CreateService(service)

		serviceRepository.DeleteService(service)

		Expect(len(memory)).To(Equal(0))
		Expect(fake.GetAll("services")[0]).To(BeNil())
	})
})

var _ = Describe("get all services", func() {
	It("should return all values", func() {
		service := domain.Service{}
		memory := make(map[string]interface{})
		fake := fakePersistence{memory: memory}
		serviceRepository := repositories.CreateServiceRepository(&fake)
		serviceRepository.CreateService(service)

		results := serviceRepository.GetAllServices()

		Expect(fake.GetAll("services")).To(Equal(results))
	})
})
