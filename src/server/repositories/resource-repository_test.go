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

var _ = Describe("create resource", func() {
	It("should persist a service resource", func() {
		id := "10174c96-a835-4e9e-b49e-9085f6e63368"
		service := domain.Service{ID: id}
		resource := domain.Resource{Type: "Service", Object: service}
		memory := make(map[string]interface{})
		fake := fakePersistence{memory: memory}
		resourceRepository := repositories.CreateResourceRepository(&fake)

		error := resourceRepository.CreateResource(resource)

		Expect(error).To(BeNil())
		Expect(fake.GetAll()[0]).To(Equal(resource))
	})

	It("should fail if missing id in service resource", func() {
		service := domain.Service{}
		resource := domain.Resource{Type: "Service", Object: service}
		memory := make(map[string]interface{})
		fake := fakePersistence{memory: memory}
		resourceRepository := repositories.CreateResourceRepository(&fake)

		error := resourceRepository.CreateResource(resource)

		Expect(error).NotTo(BeNil())
	})

	It("should persist a deployment resource", func() {
		id := "10174c96-a835-4e9e-b49e-9085f6e63368"
		deployment := domain.Deployment{ID: id}
		resource := domain.Resource{Type: "Deployment", Object: deployment}
		memory := make(map[string]interface{})
		fake := fakePersistence{memory: memory}
		resourceRepository := repositories.CreateResourceRepository(&fake)

		error := resourceRepository.CreateResource(resource)

		Expect(error).To(BeNil())
		Expect(fake.GetAll()[0]).To(Equal(resource))
	})

	It("should fail if missing id in service resource", func() {
		service := domain.Deployment{}
		resource := domain.Resource{Type: "Deployment", Object: service}
		memory := make(map[string]interface{})
		fake := fakePersistence{memory: memory}
		resourceRepository := repositories.CreateResourceRepository(&fake)

		error := resourceRepository.CreateResource(resource)

		Expect(error).NotTo(BeNil())
	})
})

var _ = Describe("delete service resource", func() {
	It("should delete a given service resource", func() {
		id := "10174c96-a835-4e9e-b49e-9085f6e63368"
		service := domain.Service{ID: id}
		resource := domain.Resource{Type: "Service", Object: service}
		memory := make(map[string]interface{})
		fake := fakePersistence{memory: memory}
		resourceRepository := repositories.CreateResourceRepository(&fake)
		resourceRepository.CreateResource(resource)

		error := resourceRepository.DeleteResource(id)

		Expect(error).To(BeNil())
		Expect(len(memory)).To(Equal(0))
	})

	It("should fail if missing id for service resouce", func() {
		id := "10174c96-a835-4e9e-b49e-9085f6e63368"
		service := domain.Service{ID: id}
		resource := domain.Resource{Type: "Service", Object: service}
		memory := make(map[string]interface{})
		fake := fakePersistence{memory: memory}
		resourceRepository := repositories.CreateResourceRepository(&fake)
		resourceRepository.CreateResource(resource)

		error := resourceRepository.DeleteResource("")

		Expect(error).NotTo(BeNil())
	})
})

var _ = Describe("get all resources", func() {
	It("should return all values", func() {
		id := "10174c96-a835-4e9e-b49e-9085f6e63368"
		service := domain.Service{ID: id}
		resource := domain.Resource{Type: "Service", Object: service}
		memory := make(map[string]interface{})
		fake := fakePersistence{memory: memory}
		resourceRepository := repositories.CreateResourceRepository(&fake)
		error := resourceRepository.CreateResource(resource)

		results := resourceRepository.GetAllResources()

		Expect(error).To(BeNil())
		Expect(fake.GetAll()).To(Equal(results))
	})
})
