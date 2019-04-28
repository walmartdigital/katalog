package persistence_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/walmartdigital/katalog/src/server/persistence"
)

func TestAll(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "memory persistence")
}

var _ = Describe("create", func() {
	It("should add a single object", func() {
		memory := make(map[string]interface{})
		persistence := persistence.BuildMemoryPersistence(memory)
		id := "4128cbf6-b279-46b3-ae19-9f90ea190978"
		value := struct{ id string }{id}

		error := persistence.Create(id, value)

		Expect(error).To(BeNil())
		Expect(memory["4128cbf6-b279-46b3-ae19-9f90ea190978"]).To(Equal(value))
	})

	It("should fail when id is empty", func() {
		memory := make(map[string]interface{})
		persistence := persistence.BuildMemoryPersistence(memory)
		id := ""
		value := struct{ id string }{id}

		error := persistence.Create(id, value)

		Expect(error).NotTo(BeNil())
	})
})

var _ = Describe("delete", func() {
	It("should remove one item", func() {
		memory := make(map[string]interface{})
		persistence := persistence.BuildMemoryPersistence(memory)
		id := "7879d950-e511-4798-a074-a951d9eddbb8"
		value := struct{ id string }{id: id}
		memory[id] = value
		expected := make([]interface{}, 1)
		expected[0] = value

		error := persistence.Delete(id)

		Expect(error).To(BeNil())
		Expect(len(memory)).To(Equal(0))
		Expect(memory[id]).To(BeNil())
	})

	It("should fail when id is empty", func() {
		memory := make(map[string]interface{})
		persistence := persistence.BuildMemoryPersistence(memory)
		id := ""
		value := struct{ id string }{id: id}
		memory[id] = value
		expected := make([]interface{}, 1)
		expected[0] = value

		error := persistence.Delete(id)

		Expect(error).NotTo(BeNil())
		Expect(memory[id]).NotTo(BeNil())
	})
})

var _ = Describe("get all", func() {
	It("should return all values", func() {
		memory := make(map[string]interface{})
		persistence := persistence.BuildMemoryPersistence(memory)
		value := struct{ id string }{""}
		memory[""] = value
		expected := make([]interface{}, 1)
		expected[0] = value

		results := persistence.GetAll()

		Expect(results).To(Equal(expected))
	})
})
