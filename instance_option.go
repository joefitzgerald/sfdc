package sfdc

import "net/http"

type InstanceOption interface {
	applyToInstance(i *Instance)
}

type withAPIVersion struct {
	apiVersion string
}

func (w *withAPIVersion) applyToInstance(i *Instance) {
	i.apiVersion = w.apiVersion
}

func WithAPIVersion(version string) InstanceOption {
	return &withAPIVersion{apiVersion: version}
}

type withURL struct {
	url string
}

func (w *withURL) applyToInstance(i *Instance) {
	i.url = w.url
}

func WithURL(url string) InstanceOption {
	return &withURL{url: url}
}

type withHTTPClient struct {
	client *http.Client
}

func (w *withHTTPClient) applyToInstance(i *Instance) {
	i.client = w.client
}

func WithHTTPClient(client *http.Client) InstanceOption {
	return &withHTTPClient{client: client}
}
