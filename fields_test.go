package sfdc_test

import (
	"testing"

	"github.com/joefitzgerald/sfdc"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
)

func testFields(t *testing.T, when spec.G, it spec.S) {
	it.Before(func() {
		RegisterTestingT(t)
	})

	it("returns fields with json tags", func() {
		type t struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		}
		e := sfdc.NewEntity[t](&sfdc.Instance{})
		Expect(e.TaggedFields()).To(Equal("id,name"))
	})

	it("skips fields with json - tags", func() {
		type t struct {
			ID   string `json:"id"`
			Name string `json:"-"`
		}
		e := sfdc.NewEntity[t](&sfdc.Instance{})
		Expect(e.TaggedFields()).To(Equal("id"))
	})

	it("skips fields with sfdc - tags", func() {
		type t struct {
			ID   string `json:"id"`
			Name string `json:"name" sfdc:"-"`
		}
		e := sfdc.NewEntity[t](&sfdc.Instance{})
		Expect(e.TaggedFields()).To(Equal("id"))
	})

	it("returns fields with sfdc tags", func() {
		type t struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description" sfdc:"func(description)"`
		}
		e := sfdc.NewEntity[t](&sfdc.Instance{})
		Expect(e.TaggedFields()).To(Equal("id,name,func(description)"))
	})

	it("handles nested queries", func() {
		type t2 struct {
			ID           string `json:"id"`
			AttachmentID string `json:"attachmentid"`
		}
		type t struct {
			ID          string                 `json:"id"`
			Name        string                 `json:"name"`
			Attachments sfdc.QueryResponse[t2] `json:"attachments" sfdc:"(SELECT Id,AttachmentId FROM Attachment)"`
		}
		e := sfdc.NewEntity[t](&sfdc.Instance{})
		Expect(e.TaggedFields()).To(Equal("id,name,(SELECT Id,AttachmentId FROM Attachment)"))
	})
}
