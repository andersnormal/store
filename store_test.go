package store_test

import (
	"math/rand"
	"os"
	"sync"
	"time"

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

	Describe("Put() and Get() on nil", func() {
		It("Should not allow to Put() nil value", func() {
			err = s.Put(store.Byte(DefaultKeyString), nil)
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(store.ErrBadValue))
		})

		It("Should allow to Get() nil value", func() {
			err = s.Put(store.Byte(DefaultKeyString), "42")
			Expect(err).NotTo(HaveOccurred())

			err = s.Get(store.Byte(DefaultKeyString), nil)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("Put() and Get() on goroutine", func() {
		It("Should be goroutine-safe", func(done Done) {
			var wg sync.WaitGroup

			rand.Seed(time.Now().UnixNano())

			for i := 0; i < 1000; i++ {
				wg.Add(1)
				go func() {
					defer GinkgoRecover()
					defer wg.Done()

					switch rand.Intn(3) {
					case 0:
						err = s.Put(store.Byte(DefaultKeyString), "42")
						Expect(err).NotTo(HaveOccurred())
					case 1:
						var val string
						err = s.Get(store.Byte(DefaultKeyString), &val)
						if err != nil && err != store.ErrKeyNotExist {
							Expect(err).To(HaveOccurred())
							return
						}
						Expect(val).To(Equal("42"))
					case 2:
						err = s.Delete(store.Byte(DefaultKeyString))
						Expect(err).NotTo(HaveOccurred())
					}
				}()
			}
			wg.Wait()
			close(done)
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
