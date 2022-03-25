package sfdc

import (
	"net/http"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
)

func testInstanceOption(t *testing.T, when spec.G, it spec.S) {
	it.Before(func() {
		RegisterTestingT(t)
	})

	when("withURL", func() {
		it("applies the URL to the instance", func() {
			instance, err := New(WithNoAuthentication(), WithURL("https://example.com"))
			Expect(err).NotTo(HaveOccurred())
			Expect(instance).NotTo(BeNil())
			Expect(instance.url).To(Equal("https://example.com"))
		})
	})

	when("withHTTPClient", func() {
		it("is not supplied, http.DefaultClient is used", func() {
			instance, err := New(WithNoAuthentication())
			Expect(err).NotTo(HaveOccurred())
			Expect(instance).NotTo(BeNil())
			Expect(instance.client).To(Equal(http.DefaultClient))
		})

		it("the provided client is used", func() {
			client := &http.Client{
				Timeout: time.Second * 10,
			}
			instance, err := New(WithNoAuthentication(), WithHTTPClient(client))
			Expect(err).NotTo(HaveOccurred())
			Expect(instance).NotTo(BeNil())
			Expect(instance.client).NotTo(Equal(http.DefaultClient))
			Expect(instance.client).To(Equal(client))
		})
	})

	when("withAPIVersion", func() {
		it("is not supplied, a default is set", func() {
			instance, err := New(WithNoAuthentication())
			Expect(err).NotTo(HaveOccurred())
			Expect(instance).NotTo(BeNil())
			Expect(instance.apiVersion).NotTo(BeZero())
		})

		it("the provided API version is used", func() {
			instance, err := New(WithNoAuthentication(), WithAPIVersion("v100.0"))
			Expect(err).NotTo(HaveOccurred())
			Expect(instance).NotTo(BeNil())
			Expect(instance.apiVersion).To(Equal("v100.0"))
		})
	})
}
