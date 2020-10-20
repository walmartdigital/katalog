package repositories_test

import (
	"errors"
	"testing"

	"github.com/emirpasic/gods/lists/arraylist"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/walmartdigital/katalog/domain"
	"github.com/walmartdigital/katalog/server/repositories"
)

type fakePersistence struct {
	memory map[string]interface{}
	fail   bool
}

func (f *fakePersistence) Create(id string, obj interface{}) error {
	if id == "" {
		return errors.New("")
	}
	f.memory[id] = obj
	return nil
}

// Update ...
func (f *fakePersistence) Update(id string, obj interface{}) error {
	if id == "" {
		return errors.New("")
	}
	f.memory[id] = obj
	return nil
}

func (f *fakePersistence) Get(id string) (interface{}, error) {
	if id == "" {
		return nil, errors.New("")
	}
	res := f.memory[id]
	return res, nil
}

func (f *fakePersistence) Delete(id string) error {
	if id == "" {
		return errors.New("")
	}
	delete(f.memory, id)
	return nil
}

func (f *fakePersistence) GetAll() ([]interface{}, error) {
	list := arraylist.New()
	if f.fail {
		return nil, errors.New("")
	}
	for _, value := range f.memory {
		list.Add(value)
	}
	return list.Values(), nil
}

func TestAll(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "service repository")
}

var _ = Describe("create resource", func() {
	It("should persist a service resource", func() {
		id := "10174c96-a835-4e9e-b49e-9085f6e63368"
		resource := domain.Resource{K8sResource: &domain.Service{ID: id}}
		memory := make(map[string]interface{})
		fake := fakePersistence{memory: memory, fail: false}
		resourceRepository := repositories.CreateResourceRepository(&fake)

		error := resourceRepository.CreateResource(resource)

		Expect(error).To(BeNil())
		allServices, _ := fake.GetAll()
		Expect(allServices[0]).To(Equal(resource))
	})

	It("should fail if missing id in service resource", func() {
		resource := domain.Resource{K8sResource: &domain.Service{}}
		memory := make(map[string]interface{})
		fake := fakePersistence{memory: memory}
		resourceRepository := repositories.CreateResourceRepository(&fake)

		error := resourceRepository.CreateResource(resource)

		Expect(error).NotTo(BeNil())
	})

	It("should persist a deployment resource", func() {
		id := "10174c96-a835-4e9e-b49e-9085f6e63368"
		resource := domain.Resource{K8sResource: &domain.Deployment{ID: id}}
		memory := make(map[string]interface{})
		fake := fakePersistence{memory: memory, fail: false}
		resourceRepository := repositories.CreateResourceRepository(&fake)

		error := resourceRepository.CreateResource(resource)

		Expect(error).To(BeNil())
		allDeployments, _ := fake.GetAll()
		Expect(allDeployments[0]).To(Equal(resource))
	})

	It("should fail if missing id in service resource", func() {
		resource := domain.Resource{K8sResource: &domain.Deployment{}}
		memory := make(map[string]interface{})
		fake := fakePersistence{memory: memory}
		resourceRepository := repositories.CreateResourceRepository(&fake)

		error := resourceRepository.CreateResource(resource)

		Expect(error).NotTo(BeNil())
	})
})

var _ = Describe("get resource", func() {
	It("should retrieve a given resource", func() {
		id := "10174c96-a835-4e9e-b49e-9085f6e63368"
		resource := domain.Resource{K8sResource: &domain.Service{ID: id}}
		memory := make(map[string]interface{})
		fake := fakePersistence{memory: memory}
		resourceRepository := repositories.CreateResourceRepository(&fake)
		resourceRepository.CreateResource(resource)

		res, error := resourceRepository.GetResource(id)

		Expect(error).To(BeNil())
		Expect(res).To(Equal(resource))
	})

	It("should not retrieve a non-existing resource", func() {
		id := "10174c96-a835-4e9e-b49e-9085f6e63368"
		resource := domain.Resource{K8sResource: &domain.Service{ID: id}}
		memory := make(map[string]interface{})
		fake := fakePersistence{memory: memory}
		resourceRepository := repositories.CreateResourceRepository(&fake)
		resourceRepository.CreateResource(resource)

		res, _ := resourceRepository.GetResource("10174c96-a835-4e9e-b49e-9085f6e633ff")

		Expect(res).To(BeNil())
	})

	It("should return an error if trying to retrieve a resource with empty ID", func() {
		id := "10174c96-a835-4e9e-b49e-9085f6e63368"
		resource := domain.Resource{K8sResource: &domain.Service{ID: id}}
		memory := make(map[string]interface{})
		fake := fakePersistence{memory: memory}
		resourceRepository := repositories.CreateResourceRepository(&fake)
		resourceRepository.CreateResource(resource)

		res, error := resourceRepository.GetResource("")

		Expect(res).To(BeNil())
		Expect(error).NotTo(BeNil())
	})
})

var _ = Describe("delete service resource", func() {
	It("should delete a given service resource", func() {
		id := "10174c96-a835-4e9e-b49e-9085f6e63368"
		resource := domain.Resource{K8sResource: &domain.Service{ID: id}}
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
		resource := domain.Resource{K8sResource: &domain.Service{ID: id}}
		memory := make(map[string]interface{})
		fake := fakePersistence{memory: memory}
		resourceRepository := repositories.CreateResourceRepository(&fake)
		resourceRepository.CreateResource(resource)

		error := resourceRepository.DeleteResource("")

		Expect(error).NotTo(BeNil())
	})
})

var _ = Describe("update service resource", func() {
	It("should update a given service resource", func() {
		r1 := domain.Resource{K8sResource: &domain.Service{ID: "10174c96-a835-4e9e-b49e-9085f6e63368", Generation: 1}}
		r2 := domain.Resource{K8sResource: &domain.Service{ID: "10174c96-a835-4e9e-b49e-9085f6e63368", Generation: 2}}

		memory := make(map[string]interface{})
		fake := fakePersistence{memory: memory}
		resourceRepository := repositories.CreateResourceRepository(&fake)
		resourceRepository.CreateResource(r1)

		r, error := resourceRepository.UpdateResource(r2)

		Expect(error).To(BeNil())
		Expect(r2).To(Equal(*r))
	})

	It("should not update a given service resource because Generation is not greater than stored object", func() {
		r1 := domain.Resource{K8sResource: &domain.Service{ID: "10174c96-a835-4e9e-b49e-9085f6e63368", Generation: 1}}
		r2 := domain.Resource{K8sResource: &domain.Service{ID: "10174c96-a835-4e9e-b49e-9085f6e63368", Generation: 1}}

		memory := make(map[string]interface{})
		fake := fakePersistence{memory: memory}
		resourceRepository := repositories.CreateResourceRepository(&fake)
		resourceRepository.CreateResource(r1)

		r, error := resourceRepository.UpdateResource(r2)

		Expect(error).To(BeNil())
		Expect(r).To(BeNil())
	})

	It("should fail if missing id for service resource", func() {
		id := "10174c96-a835-4e9e-b49e-9085f6e63368"
		resource := domain.Resource{K8sResource: &domain.Service{ID: id}}
		memory := make(map[string]interface{})
		fake := fakePersistence{memory: memory}
		resourceRepository := repositories.CreateResourceRepository(&fake)
		resourceRepository.CreateResource(resource)

		badres := domain.Resource{K8sResource: &domain.Service{ID: ""}}
		_, error := resourceRepository.UpdateResource(badres)

		Expect(error).NotTo(BeNil())
	})
})

var _ = Describe("get all resources", func() {
	It("should return all values", func() {
		id := "10174c96-a835-4e9e-b49e-9085f6e63368"
		resource := domain.Resource{K8sResource: &domain.Service{ID: id}}
		memory := make(map[string]interface{})
		fake := fakePersistence{memory: memory, fail: false}
		resourceRepository := repositories.CreateResourceRepository(&fake)
		error := resourceRepository.CreateResource(resource)

		results, _ := resourceRepository.GetAllResources()

		Expect(error).To(BeNil())
		Expect(fake.GetAll()).To(Equal(results))
	})

	It("should return an error when retrieving all elements from persistence fails", func() {
		id := "10174c96-a835-4e9e-b49e-9085f6e63368"
		resource := domain.Resource{K8sResource: &domain.Service{ID: id}}
		memory := make(map[string]interface{})
		fake := fakePersistence{memory: memory, fail: true}
		resourceRepository := repositories.CreateResourceRepository(&fake)
		resourceRepository.CreateResource(resource)

		results, error := resourceRepository.GetAllResources()

		Expect(error).NotTo(BeNil())
		Expect(results).To(BeNil())
	})
})
