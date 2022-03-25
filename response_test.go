package sfdc

import (
	"strings"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
)

func testResponse(t *testing.T, when spec.G, it spec.S) {
	it.Before(func() {
		RegisterTestingT(t)
	})

	when("something goes wrong", func() {
		it("returns an error", func() {
			response := strings.NewReader(`[
				{"message": "test error!", "errorCode": "ERR_TEST"}
			]`)

			err := errorForResponse(response)
			Expect(err.Error()).To(ContainSubstring("test error"))
			Expect(err.Error()).To(ContainSubstring("ERR_TEST"))
		})
	})
}
