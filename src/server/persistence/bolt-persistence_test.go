package persistence_test

import (
	"errors"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/walmartdigital/katalog/src/server/persistence"
)

type dummyDriver struct {
}

func (d *dummyDriver) Update(fn func(persistence.BoltTxInterface) error) error {
	return nil
}

func (d *dummyDriver) View(fn func(persistence.BoltTxInterface) error) error {
	return errors.New("Dummy Error")
}

func (d *dummyDriver) Close() error {
	return nil
}

func TestAll(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Utils")
}

var _ = Describe("Persistence | Bolt Persistence", func() {
	Describe("GetAll", func() {
		It("should got empty list on error", func() {
			kind := "services"
			driver := &dummyDriver{}
			persistence := persistence.CreateBoltDriver(driver)

			output := persistence.GetAll(kind)

			Expect(output).To(BeEmpty())
		})
	})
})
