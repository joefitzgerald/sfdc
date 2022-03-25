package sfdc_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/joefitzgerald/sfdc"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
)

func testEntity(t *testing.T, when spec.G, it spec.S) {
	type testEntity struct {
		ID string `json:"Id"`
	}

	var (
		instance *sfdc.Instance
		entity   *sfdc.Entity[testEntity]
	)
	it.Before(func() {
		RegisterTestingT(t)
	})

	when("building a query", func() {
		it("uses the entity name in the query", func() {
			type testEntity struct {
				ID string `json:"Id"`
			}
			e := sfdc.NewEntity[testEntity](&sfdc.Instance{})
			q := e.BuildQuery("Id", "")
			Expect(q).To(Equal("SELECT Id FROM testEntity"))
		})

		it("allows the entity name to be set", func() {
			type testEntity struct {
				ID string `json:"Id"`
			}
			e := sfdc.NewEntity[testEntity](&sfdc.Instance{})
			e.SetName("Opportunity")
			q := e.BuildQuery("Id", "")
			Expect(q).To(Equal("SELECT Id FROM Opportunity"))
		})
	})

	when("using an invalid instance url", func() {
		it.Before(func() {
			var err error
			instance, err = sfdc.New(sfdc.WithNoAuthentication(), sfdc.WithURL(":::::::---INVALID---"))
			Expect(err).NotTo(HaveOccurred())
			Expect(instance).NotTo(BeNil())
			entity = sfdc.NewEntity[testEntity](instance)
			Expect(entity).NotTo(BeNil())
		})

		it("Query fails", func() {
			_, err := entity.Query(context.Background(), "SELECT Id FROM Account")
			Expect(err).To(HaveOccurred())
		})

		it("QueryAsync returns an error", func() {
			records, errs := entity.QueryAsync(context.Background(), "SELECT Id FROM Account")
			for rec := range records {
				Expect(rec).NotTo(BeNil())
				Expect(rec).To(HaveLen(0))
			}
			var err error
			select {
			case err = <-errs:
			default:
			}
			Expect(err).To(HaveOccurred())
		})
	})

	when("using a valid instance", func() {
		var (
			server  *httptest.Server
			handler func(w http.ResponseWriter, r *http.Request)
		)

		it.Before(func() {
			handler = func(w http.ResponseWriter, r *http.Request) {
				// No-op handler
			}
			server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handler(w, r)
			}))
			var err error
			instance, err = sfdc.New(sfdc.WithNoAuthentication(), sfdc.WithURL(server.URL))
			Expect(err).NotTo(HaveOccurred())
			Expect(instance).NotTo(BeNil())
		})

		it.After(func() {
			server.Close()
		})

		when("using a test entity", func() {
			it.Before(func() {
				entity = sfdc.NewEntity[testEntity](instance)
				Expect(entity).NotTo(BeNil())
			})

			when("Query()", func() {
				it("returns records", func() {
					handler = func(w http.ResponseWriter, r *http.Request) {
						Expect(r.Method).To(Equal(http.MethodGet))
						Expect(r.URL.Path).To(Equal("/services/data/v54.0/queryAll"))
						Expect(r.URL.Query().Get("q")).To(Equal("SELECT Id FROM Account"))
						w.Write([]byte(`{"records":[], "done":true}`))
					}
					result, err := entity.Query(context.Background(), "SELECT Id FROM Account")
					Expect(err).NotTo(HaveOccurred())
					Expect(result).NotTo(BeNil())
					Expect(result).To(HaveLen(0))
				})

				it("fetches all records", func() {
					callCount := 0
					handler = func(w http.ResponseWriter, r *http.Request) {
						callCount += 1
						Expect(r.Method).To(Equal(http.MethodGet))
						if callCount == 1 {
							Expect(r.URL.Path).To(Equal("/services/data/v54.0/queryAll"))
							Expect(r.URL.Query().Get("q")).To(Equal("SELECT Id FROM Account"))
						} else {
							Expect(r.URL.Path).To(Equal("/next"))
						}

						done := "false"
						if callCount > 1 {
							done = "true"
						}
						w.Write([]byte(strings.ReplaceAll(`{"records":[{"Id": "test"}], "done":{done}, "nextRecordsUrl": "/next"}`, "{done}", done)))
					}
					result, err := entity.Query(context.Background(), "SELECT Id FROM Account")
					Expect(err).NotTo(HaveOccurred())
					Expect(callCount).To(Equal(2))
					Expect(result).NotTo(BeNil())
					Expect(result).To(HaveLen(2))
				})
			})

			when("QueryAsync", func() {
				it("returns records", func() {
					callCount := 0
					handler = func(w http.ResponseWriter, r *http.Request) {
						callCount += 1
						Expect(r.Method).To(Equal(http.MethodGet))
						if callCount == 1 {
							Expect(r.URL.Path).To(Equal("/services/data/v54.0/queryAll"))
							Expect(r.URL.Query().Get("q")).To(Equal("SELECT Id FROM Account"))
						} else {
							Expect(r.URL.Path).To(Equal("/next"))
						}

						done := "false"
						if callCount > 1 {
							done = "true"
						}
						w.Write([]byte(strings.ReplaceAll(`{"records":[{"Id": "test"}], "done":{done}, "nextRecordsUrl": "/next"}`, "{done}", done)))
					}
					var result []testEntity
					records, errs := entity.QueryAsync(context.Background(), "SELECT Id FROM Account")
					for rec := range records {
						result = append(result, rec...)
						Expect(rec).NotTo(BeNil())
						Expect(rec).To(HaveLen(1))
					}
					var err error
					select {
					case err = <-errs:
					default:
					}
					Expect(err).NotTo(HaveOccurred())
					Expect(callCount).To(Equal(2))
					Expect(result).NotTo(BeNil())
					Expect(result).To(HaveLen(2))
				})
			})
		})
	})
}
