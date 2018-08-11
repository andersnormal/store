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

type TestRich struct {
	String string
	Int    int
}

var _ = Describe("Store", func() {
	var (
		s   *store.Store
		err error
	)

	BeforeEach(func() {
		os.Remove(DefaultStore)

		s, err = store.Open(DefaultStore, nil)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
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
			By("Putting and getting 42")

			err = s.Put(store.Byte(DefaultKeyString), "42")
			Expect(err).NotTo(HaveOccurred())

			var val string
			err = s.Get(store.Byte(DefaultKeyString), &val)
			Expect(err).NotTo(HaveOccurred())
			Expect(val).To(Equal("42"))
		})

		It("Should allow to Delete() keys", func() {
			By("Deleting 42")

			err = s.Put(store.Byte(DefaultKeyString), "42")
			Expect(err).NotTo(HaveOccurred())

			err = s.Delete(store.Byte(DefaultKeyString))
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("Put() and Get() rich types", func() {
		It("Should allow to Put() rich types", func() {
			r := &TestRich{"42", 42}

			By("Putting rich type")

			err = s.Put(store.Byte(DefaultKeyString), r)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Should allow to Put() and Get() rich types", func() {
			r := &TestRich{"42", 42}

			By("Putting and getting rich type")

			err = s.Put(store.Byte(DefaultKeyString), r)
			Expect(err).NotTo(HaveOccurred())

			var val TestRich
			err = s.Get(store.Byte(DefaultKeyString), &val)
			Expect(err).NotTo(HaveOccurred())
			Expect(val.Int).To(Equal(42))
			Expect(val.String).To(Equal("42"))
		})
	})

	Measure("it should Put() very efficient", func(b Benchmarker) {
		runtime := b.Time("runtime", func() {
			err = s.Put(store.Byte(DefaultKeyString), "42")
			Expect(err).NotTo(HaveOccurred())
		})

		Ω(runtime.Seconds()).Should(BeNumerically("<", 0.2), "Put() shouldn't take too long.")
	}, 10)
})
