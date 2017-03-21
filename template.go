package main

import (
	"strings"
	"text/template"
)

func typefor(sfdctype string) string {
	if len(sfdctype) == 0 {
		return "interface{}"
	}
	switch strings.ToLower(sfdctype) {
	case "address":
		return "*Address"
	case "boolean":
		return "bool"
	case
		"calculated",
		"combobox",
		"date",
		"datetime",
		"email",
		"encryptedstring",
		"id",
		"masterrecord",
		"multipicklist",
		"phone",
		"picklist",
		"reference",
		"string",
		"textarea",
		"time",
		"url":
		return "string"
	case
		"currency",
		"double",
		"percent":
		return "float64"
	case "int":
		return "int"
	case "junctionidlist":
		return "[]string"
	case "location":
		return "*Geolocation"
	}

	return "interface{}"
}

func cleanname(name string) string {
	name = strings.Replace(name, "__c", "", -1)
	if strings.Contains(name, "_") {
		newName := ""
		splitName := strings.Split(name, "_")
		for _, item := range splitName {
			if len(item) == 0 {
				continue
			}
			if len(item) == 1 {
				newName += strings.ToUpper(item)
				continue
			}
			if len(item) == 2 {
				if strings.ToUpper(item) == "ID" {
					newName += strings.ToUpper(item)
					continue
				}
			}
			newName += strings.ToUpper(string(item[0])) + string(item[1:])
		}
		name = newName
	}

	if strings.HasSuffix(name, "Id") {
		name = strings.TrimRight(name, "Id") + "ID"
	}
	if strings.HasSuffix(name, "Url") {
		name = strings.TrimRight(name, "Url") + "URL"
	}
	if strings.HasSuffix(name, "Api") {
		name = strings.TrimRight(name, "Api") + "API"
	}
	if strings.HasPrefix(name, "Api") {
		name = "API" + strings.TrimLeft(name, "Api")
	}

	return strings.Title(name)
}

func jsontag(tagName string) string {
	return "`json:\"" + tagName + "\"`"
}

func cleannamelower(name string) string {
	return strings.ToLower(cleanname(name))
}

func valueforkey(key string, m map[string]string) string {
	return m[key]
}

func backtick() string {
	return "`"
}

var generatedTmpl = template.Must(template.New("generated").Funcs(template.FuncMap{
	"jsontag":        jsontag,
	"cleanname":      cleanname,
	"cleannamelower": cleannamelower,
	"typefor":        typefor,
	"tolower":        strings.ToLower,
	"backtick":       backtick,
}).Parse(`
// generated by sfdc; DO NOT EDIT

package {{.PackageName}}

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"
)

// {{cleanname .SObject.Name}} is the {{.SObject.Label}} SObject
type {{cleanname .SObject.Name}} struct {
  {{range .Fields}}{{cleanname .FieldName}} {{typefor .SObjectField.Type}} {{jsontag .SObjectField.Name}} // {{.SObjectField.Type}}
  {{end}}
}

type {{cleannamelower .SObject.Name}} struct {
	config      *Config
	InstanceURL string
	DescribeURL string
	SobjectURL  string
	QueryURL    string
	QueryAllURL string
	allFields   string
}

// {{cleanname .SObject.Name}}QueryResponse is the result of a SOQL Query for the
// {{cleanname .SObject.Name}} object.
type {{cleanname .SObject.Name}}QueryResponse struct {
	Done           bool                          {{jsontag "done"}}
	NextRecordsURL string                        {{jsontag "nextRecordsUrl"}}
	Records        []{{cleanname .SObject.Name}} {{jsontag "records"}}
	TotalSize      int                           {{jsontag "totalSize"}}
}

func (o *{{cleannamelower .SObject.Name}}) AllFields() string {
	if o.allFields != "" {
		return o.allFields
	}

	s := []string{
		{{range .Fields}}"{{.SObjectField.Name}}",
		{{end}}
	}
	o.allFields = strings.Join(s, ", ")
 	return o.allFields
}

func (o *{{cleannamelower .SObject.Name}}) Get(ctx context.Context, id string) (*{{cleanname .SObject.Name}}, error) {
	uri := fmt.Sprintf("%v/%v/%v", o.InstanceURL, o.SobjectURL, id)
	req, err := BuildRequest(ctx, uri)
	if err != nil {
		return nil, err
	}
	res, err := o.config.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var r {{cleanname .SObject.Name}}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}
	return &r, nil
}

func (o *{{cleannamelower .SObject.Name}}) Query(ctx context.Context, fields string, constraints string) ([]{{cleanname .SObject.Name}}, error) {
	query := fmt.Sprintf("SELECT %v FROM {{cleanname .SObject.Name}}", fields)
	if utf8.RuneCountInString(constraints) > 0 {
		query = fmt.Sprintf("%v WHERE %v", query, constraints)
	}
	uri, _ := url.Parse(fmt.Sprintf("%v%v", o.InstanceURL, o.QueryURL))
	q := uri.Query()
	q.Set("q", query)
	uri.RawQuery = q.Encode()
	reqURI := uri.String()

	var r {{cleanname .SObject.Name}}QueryResponse
	results := r.Records
	for !r.Done {
		if r.NextRecordsURL != "" {
			reqURI = fmt.Sprintf("%v%v", o.InstanceURL, r.NextRecordsURL)
		}
		req, err := BuildRequest(ctx, reqURI)
		if err != nil {
			return nil, err
		}

		res, err := o.config.Client.Do(req)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			return nil, err
		}
		results = append(results, r.Records...)
	}

	return results, nil
}

// QueryAsync returns a channel that {{cleannamelower .SObject.Name}} are written to.
// The channel is closed when all records have been written.
// Errors are written to the returned error channel.
// The query aborts when an error is encountered.
func (o *{{cleannamelower .SObject.Name}}) QueryAsync(ctx context.Context, fields string, constraints string) (<- chan []{{cleanname .SObject.Name}}, <- chan error) {
	query := fmt.Sprintf("SELECT %v FROM {{cleanname .SObject.Name}}", fields)
	if utf8.RuneCountInString(constraints) > 0 {
		query = fmt.Sprintf("%v WHERE %v", query, constraints)
	}
	uri, _ := url.Parse(fmt.Sprintf("%v%v", o.InstanceURL, o.QueryURL))
	q := uri.Query()
	q.Set("q", query)
	uri.RawQuery = q.Encode()
	reqURI := uri.String()

	result := make(chan []{{cleanname .SObject.Name}})
	errs := make(chan error, 1)
	go func() {
		var r {{cleanname .SObject.Name}}QueryResponse
		for !r.Done {
			if r.NextRecordsURL != "" {
				reqURI = fmt.Sprintf("%v%v", o.InstanceURL, r.NextRecordsURL)
			}
			req, err := BuildRequest(ctx, reqURI)
			if err != nil {
				errs <- err
				return
			}

			res, err := o.config.Client.Do(req)
			if err != nil {
				errs <- err
				return
			}
			defer res.Body.Close()

			if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
				errs <- err
				return
			}
			result <- r.Records
		}
		close(result)
	}()

	return result, errs
}
`))

var commonTmpl = template.Must(template.New("common").Funcs(template.FuncMap{
	"jsontag":        jsontag,
	"cleanname":      cleanname,
	"cleannamelower": cleannamelower,
	"typefor":        typefor,
	"tolower":        strings.ToLower,
	"valueforkey":    valueforkey,
	"backtick":       backtick,
}).Parse(`
// generated by sfdc; DO NOT EDIT

package {{.PackageName}}

import (
	"context"
	"errors"
	"net/http"

	"golang.org/x/oauth2"
)

const (
	// SFDCDateLayout is used to convert SFDC DateTime strings correctly
	SFDCDateLayout = "2006-01-02T15:04:05.000Z"
)

// Config is SFDC API configuration
type Config struct {
	Token         *oauth2.Token {{backtick}}ignored:"true"{{backtick}}
	Oauth2Config  *oauth2.Config {{backtick}}ignored:"true"{{backtick}}
	Client        *http.Client {{backtick}}ignored:"true"{{backtick}}
	Version       string {{backtick}}default:"v37.0"{{backtick}}
	ClientID      string {{backtick}}required:"true" json:"sfdc_clientid,omitempty"{{backtick}}
	ClientSecret  string {{backtick}}required:"true" json:"sfdc_clientsecret,omitempty"{{backtick}}
	Username      string {{backtick}}required:"true" json:"sfdc_username,omitempty"{{backtick}}
	Password      string {{backtick}}required:"true" json:"sfdc_password,omitempty"{{backtick}}
	SecurityToken string {{backtick}}required:"true" json:"sfdc_securitytoken,omitempty"{{backtick}}
	Environment   string {{backtick}}default:"Production" json:"sfdc_environment,omitempty"{{backtick}}
}

// API is a SalesForce REST API client.
type API struct {
	Config       *Config
	InstanceURL  string
	ResourceURL  string
	QueryURL     string
	QueryAllURL  string
	{{range .Objects}}{{.ObjectName}} *{{tolower .ObjectName}}
	{{end}}
}

// New creates a SalesForce API that you can use to access the SalesForce REST
// API.
func New(c *Config) (*API, error) {
	if c.Oauth2Config == nil {
		c.Oauth2Config = &oauth2.Config{
			ClientID:     c.ClientID,
			ClientSecret: c.ClientSecret,
			Scopes:       nil,
			Endpoint:     oauth2.Endpoint{
				AuthURL:  "https://login.salesforce.com/services/oauth2/authorize",
				TokenURL: "https://login.salesforce.com/services/oauth2/token",
			},
		}
	}
	if c.Token == nil {
		token, err := c.Oauth2Config.PasswordCredentialsToken(context.Background(), c.Username, c.Password)
		if err != nil {
			return nil, err
		}
		c.Token = token
	}
	c.Client = c.Oauth2Config.Client(context.Background(), c.Token)
	instanceURL, ok := c.Token.Extra("instance_url").(string)
	if !ok {
		return nil, errors.New("instance_url not available in the token")
	}
	api := &API{
		Config: c,
		InstanceURL:  instanceURL,
		ResourceURL:  "{{.ResourcesURI}}",
		QueryURL:     "{{.ResourcesURI}}/query",
		QueryAllURL:  "{{.ResourcesURI}}/queryAll",
		{{range .Objects}}{{.ObjectName}}: &{{tolower .ObjectName}} {
			config: c,
			InstanceURL: instanceURL,
			DescribeURL: "{{valueforkey "describe" .SObject.URLs}}",
			SobjectURL:  "{{valueforkey "sobject" .SObject.URLs}}",
			QueryURL:    "{{.ResourcesURI}}/query",
			QueryAllURL: "{{.ResourcesURI}}/queryAll",
		},
		{{end}}
	}

	return api, nil
}

// BuildRequest creates an http.Request with defaults set for SFDC
func BuildRequest(ctx context.Context, uri string) (*http.Request, error) {
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// Address is a compound data type that contains address field data.
type Address struct {
  Accuracy string
  City string
  Country string
  CountryCode string
  Latitude string
  Longitude string
  PostalCode string
  State string
  StateCode string
  Street string
}

// Geolocation is a compound data type that contains latitude and logitude
// values for geolocation fields.
type Geolocation struct {
  Latitude string
  Logitude string
}
`))
