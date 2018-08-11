package store_test

import (
	"os"

	"github.com/andersnormal/store"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	DefaultKeyString = "42"
)

var _ = Describe("Store", func() {
	var (
		s   *store.Store
		err error
	)

	BeforeSuite(func() {
		os.Remove(DefaultStore)

		s, err = store.Open(DefaultStore, nil)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterSuite(func() {
		err = s.Close()
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("Put(), Delete() and Get() values from the store", func() {
		It("Should allow to Put() values", func() {
			By("Entering 42")

			err = s.Put(store.Byte(DefaultKeyString), "42")
			Expect(err).NotTo(HaveOccurred())
		})

		It("Should allow to Get() values", func() {
			By("Getting 42")

			var val string
			err = s.Get(store.Byte(DefaultKeyString), &val)
			Expect(err).NotTo(HaveOccurred())
			Expect(val).To(Equal("42"))
		})

		It("Should allow to Delete() keys", func() {
			By("Deleting 42")

			err = s.Delete(store.Byte(DefaultKeyString))
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
