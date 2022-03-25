package sfdc_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joefitzgerald/sfdc"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"golang.org/x/oauth2"
)

type invalidTokenJSON struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int32  `json:"expires_in"`
}

type tokenJSON struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int32  `json:"expires_in"`
	InstanceURL  string `json:"instance_url"`
}

func testAuthOptions(t *testing.T, when spec.G, it spec.S) {
	it.Before(func() {
		RegisterTestingT(t)
	})

	when("using WithToken()", func() {
		var (
			server *httptest.Server
			calls  int
			config *oauth2.Config
		)

		when("the token has no refresh token", func() {
			it.Before(func() {
				config = &oauth2.Config{}
			})

			it("should return an error", func() {
				_, err := sfdc.New(sfdc.WithToken(context.Background(), config, &oauth2.Token{}))
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("refresh token is not set"))
			})
		})

		when("the token has a refresh token", func() {
			var handlerFunc func(w http.ResponseWriter, r *http.Request)
			it.Before(func() {
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					handlerFunc(w, r)
				}))
				config = &oauth2.Config{
					Endpoint: oauth2.Endpoint{
						AuthURL:  server.URL + "/auth",
						TokenURL: server.URL + "/token",
					},
				}
			})

			it.After(func() {
				server.Close()
			})

			when("the server returns a token without an instanceURL", func() {
				it.Before(func() {
					handlerFunc = func(w http.ResponseWriter, r *http.Request) {
						w.Header().Set("Content-Type", "application/json")
						p := &invalidTokenJSON{
							AccessToken:  "test-access-token",
							RefreshToken: "test-refresh-token",
							TokenType:    "Bearer",
							ExpiresIn:    1000,
						}
						json.NewEncoder(w).Encode(p)
						calls = calls + 1
					}
				})

				it("returns an error", func() {
					token := &oauth2.Token{
						RefreshToken: "test-refresh-token",
					}
					_, err := sfdc.New(sfdc.WithToken(context.Background(), config, token), sfdc.WithURL(server.URL))
					Expect(calls).To(Equal(1))
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("instance_url"))
				})
			})

			when("the server returns a token with an instanceURL", func() {
				it.Before(func() {
					handlerFunc = func(w http.ResponseWriter, r *http.Request) {
						w.Header().Set("Content-Type", "application/json")
						p := &tokenJSON{
							AccessToken:  "test-access-token",
							RefreshToken: "test-refresh-token",
							TokenType:    "Bearer",
							InstanceURL:  "test-instance",
							ExpiresIn:    1000,
						}
						json.NewEncoder(w).Encode(p)
						calls = calls + 1
					}
				})

				it("successfully fetches a new token", func() {
					token := &oauth2.Token{
						RefreshToken: "test-refresh-token",
					}
					instance, err := sfdc.New(sfdc.WithToken(context.Background(), config, token), sfdc.WithURL(server.URL))
					Expect(calls).To(Equal(1))
					Expect(err).NotTo(HaveOccurred())
					Expect(instance).NotTo(BeNil())
				})
			})
		})
	})

	when("using a valid token", func() {
		it("return", func() {

		})
	})

}
