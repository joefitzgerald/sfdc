package sfdc

import (
	"fmt"
	"net/http"
	"net/url"
)

type Instance struct {
	url        string
	client     *http.Client
	apiVersion string
}

func New(auth AuthOption, options ...InstanceOption) (*Instance, error) {
	result := &Instance{
		apiVersion: "v54.0",
		client:     http.DefaultClient,
	}

	for i := range options {
		options[i].applyToInstance(result)
	}

	err := auth.applyAuth(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (i *Instance) InstanceURL() string {
	return i.url
}

func (i *Instance) QueryURL() url.URL {
	uri, _ := url.Parse(fmt.Sprintf("%s/services/data/%s/query", i.url, i.apiVersion))
	return *uri
}

func (i *Instance) QueryAllURL() (*url.URL, error) {
	return url.Parse(fmt.Sprintf("%s/services/data/%s/queryAll", i.url, i.apiVersion))
}
