package store_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	DefaultStore = "store"
)

func TestStore(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Store Suite")
}
