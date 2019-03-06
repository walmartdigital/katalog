package persistence_test

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/seadiaz/katalog/src/server/persistence"
)

type dummyDriver struct {
}

func (d *dummyDriver) Open(path string, mode os.FileMode, options interface{}) (interface{}, error) {
	return nil, nil
}

func TestAll(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Utils")
}

var _ = Describe("Persistence | Bolt Persistence", func() {
	Describe("GetAll", func() {
		XIt("should encode a string", func() {
			kind := "services"
			driver := &dummyDriver{}
			persistence := persistence.CreateBoltDriver(driver)
			//
			// output := persistence.GetAll(kind)
			//
			// Expect(output).To(Equal(`"dummy text"`))
		})
	})
})
