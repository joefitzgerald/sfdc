package devicecode

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

func New(clientID string, clientSecret string, scopes []string) *Config {
	return newWithBase("login", clientID, clientSecret, scopes)
}

func NewWithDomain(domain string, clientID string, clientSecret string, scopes []string) *Config {
	return newWithBase(fmt.Sprintf("%s.my", domain), clientID, clientSecret, scopes)
}

func newWithBase(base string, clientID string, clientSecret string, scopes []string) *Config {
	return &Config{
		Config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Endpoint: oauth2.Endpoint{
				TokenURL: fmt.Sprintf("https://%s.salesforce.com/services/oauth2/token", base),
				AuthURL:  fmt.Sprintf("https://%s.salesforce.com/services/oauth2/authorize", base),
			},
			Scopes: scopes,
		},
		DeviceCodeURL: fmt.Sprintf("https://%s.salesforce.com/services/oauth2/token", base),
	}
}

// A deviceCode represents the user-visible code, verification URL and
// device-visible code used to allow for user authorisation of this app. The
// app should show UserCode and VerificationURL to the user.
type deviceCode struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURL string `json:"verification_uri"`
	ExpiresIn       int64  `json:"expires_in"`
	Interval        int64  `json:"interval"`
}

// A version of oauth2.Config augmented with device endpoints
type Config struct {
	*oauth2.Config
	DeviceCodeURL string
}

// A tokenOrError is either an OAuth2 Token response or an error indicating why
// such a response failed.
type tokenOrError struct {
	*oauth2.Token
	Error string `json:"error,omitempty"`
}

var (
	// ErrAccessDenied is an error returned when the user has denied this
	// app access to their account.
	ErrAccessDenied = errors.New("access denied by user")
)

const (
	deviceGrantType = "device"
)

// RequestDeviceCode will initiate the OAuth2 device authorization flow. It
// requests a device code and information on the code and URL to show to the
// user. Pass the returned DeviceCode to WaitForDeviceAuthorization.
func (c *Config) requestDeviceCode(client *http.Client) (*deviceCode, error) {
	resp, err := client.PostForm(c.DeviceCodeURL,
		url.Values{
			"client_id":     {c.ClientID},
			"scope":         {strings.Join(c.Scopes, " ")},
			"response_type": {"device_code"},
		})

	if err != nil {
		return nil, fmt.Errorf("error making request for device code: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"request for device code authorization returned status %v (%v)",
			resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	// Unmarshal response
	var dcr deviceCode
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&dcr); err != nil {
		return nil, fmt.Errorf("error decoding device code: %w", err)
	}

	return &dcr, nil
}

// WaitForDeviceAuthorization polls the token URL waiting for the user to
// authorize the app. Upon authorization, it returns the new token. If
// authorization fails then an error is returned. If that failure was due to a
// user explicitly denying access, the error is ErrAccessDenied.
func (c *Config) Token(ctx context.Context, client *http.Client) (*oauth2.Token, error) {
	code, err := c.requestDeviceCode(client)
	if err != nil {
		return nil, fmt.Errorf("error requesting device code: %w", err)
	}
	fmt.Printf("Visit: %v and enter: %v\n", code.VerificationURL, code.UserCode)
	for {
		select {
		case <-time.After(time.Duration(code.Interval) * time.Second):
			resp, err := client.PostForm(c.Endpoint.TokenURL,
				url.Values{
					// "client_secret": {config.ClientSecret},
					"client_id":  {c.ClientID},
					"code":       {code.DeviceCode},
					"grant_type": {deviceGrantType}})
			if err != nil {
				return nil, fmt.Errorf("error polling for token: %w", err)
			}
			// if resp.StatusCode != http.StatusOK {
			// 	b, err := ioutil.ReadAll(resp.Body)
			// 	if err != nil {
			// 		return nil, err
			// 	}
			// 	log.Println(string(b))
			// 	return nil, fmt.Errorf("HTTP error %v (%v) when polling for OAuth token",
			// 		resp.StatusCode, http.StatusText(resp.StatusCode))
			// }

			// Unmarshal response, checking for errors
			var token tokenOrError
			body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1<<20))
			resp.Body.Close()
			if err != nil {
				return nil, fmt.Errorf("oauth2: cannot fetch token: %v", err)
			}
			err = json.Unmarshal(body, &token)
			if err != nil {
				return nil, fmt.Errorf("cannot unmarshal token: %w", err)
			}

			switch token.Error {
			case "":
				raw := make(map[string]interface{})
				err = json.Unmarshal(body, &raw)
				if err != nil {
					return nil, err
				}
				return token.Token.WithExtra(raw), nil
			case "authorization_pending":

			case "slow_down":
				code.Interval *= 2
			case "access_denied":

				return nil, ErrAccessDenied
			default:
				return nil, fmt.Errorf("authorization failed: %v", token.Error)
			}
		case <-ctx.Done():
			return nil, errors.New("timed out waiting for authorization")
		}
	}
}
