package sfdc_test

import (
	"testing"

	"github.com/joefitzgerald/sfdc"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
)

func testInstance(t *testing.T, when spec.G, it spec.S) {
	it.Before(func() {
		RegisterTestingT(t)
	})

	it("New() creates an instance", func() {
		instance, err := sfdc.New(sfdc.WithNoAuthentication())
		Expect(err).NotTo(HaveOccurred())
		Expect(instance).NotTo(BeNil())
	})
}
