package repositories_test

import (
	"errors"
	"testing"

	"github.com/emirpasic/gods/lists/arraylist"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/walmartdigital/katalog/src/domain"
	"github.com/walmartdigital/katalog/src/server/repositories"
)

type fakePersistence struct {
	memory map[string]interface{}
}

func (f *fakePersistence) Create(id string, obj interface{}) error {
	if id == "" {
		return errors.New("")
	}
	f.memory[id] = obj
	return nil
}

func (f *fakePersistence) Delete(id string) error {
	if id == "" {
		return errors.New("")
	}
	delete(f.memory, id)
	return nil
}

func (f *fakePersistence) GetAll() []interface{} {
	list := arraylist.New()
	for _, value := range f.memory {
		list.Add(value)
	}
	return list.Values()
}

func TestAll(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "service repository")
}

var _ = Describe("create service", func() {
	It("should persist a service", func() {
		id := "10174c96-a835-4e9e-b49e-9085f6e63368"
		service := domain.Service{ID: id}
		memory := make(map[string]interface{})
		fake := fakePersistence{memory: memory}
		serviceRepository := repositories.CreateServiceRepository(&fake)

		error := serviceRepository.CreateService(service)

		Expect(error).To(BeNil())
		Expect(fake.GetAll()[0]).To(Equal(service))
	})

	It("should fail if missing id", func() {
		service := domain.Service{}
		memory := make(map[string]interface{})
		fake := fakePersistence{memory: memory}
		serviceRepository := repositories.CreateServiceRepository(&fake)

		error := serviceRepository.CreateService(service)

		Expect(error).NotTo(BeNil())
	})
})

var _ = Describe("delete service", func() {
	It("should delete a given service", func() {
		id := "10174c96-a835-4e9e-b49e-9085f6e63368"
		service := domain.Service{ID: id}
		memory := make(map[string]interface{})
		fake := fakePersistence{memory: memory}
		serviceRepository := repositories.CreateServiceRepository(&fake)
		serviceRepository.CreateService(service)

		error := serviceRepository.DeleteService(id)

		Expect(error).To(BeNil())
		Expect(len(memory)).To(Equal(0))
	})

	It("should fail if missing id", func() {
		service := domain.Service{}
		memory := make(map[string]interface{})
		fake := fakePersistence{memory: memory}
		serviceRepository := repositories.CreateServiceRepository(&fake)
		serviceRepository.CreateService(service)

		error := serviceRepository.DeleteService("")

		Expect(error).NotTo(BeNil())
	})
})

var _ = Describe("get all services", func() {
	It("should return all values", func() {
		id := "10174c96-a835-4e9e-b49e-9085f6e63368"
		service := domain.Service{ID: id}
		memory := make(map[string]interface{})
		fake := fakePersistence{memory: memory}
		serviceRepository := repositories.CreateServiceRepository(&fake)
		error := serviceRepository.CreateService(service)

		results := serviceRepository.GetAllServices()

		Expect(error).To(BeNil())
		Expect(fake.GetAll()).To(Equal(results))
	})
})
