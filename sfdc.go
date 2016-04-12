package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/joefitzgerald/sfdc/parser"
	"github.com/nimajalali/go-force/force"
)

// New creates a SFDC instance
func New(c config, objectNames string) *SFDC {
	sfdc := &SFDC{
		Resources:           make(map[string]string),
		SObjects:            make(map[string]*force.SObjectMetaData),
		SObjectDescriptions: make(map[string]*force.SObjectDescription),
	}
	api, err := force.Create(c.Version, c.ClientID, c.ClientSecret, c.Username, c.Password, c.SecurityToken, c.Environment)
	if err != nil {
		log.Fatal(err)
	}
	sfdc.Version = c.Version
	sfdc.API = api
	sfdc.ObjectNames = strings.Split(objectNames, ",")
	return sfdc
}

func getPackageName(dir string) string {
	pkg, err := parser.GetPackageName(dir, *outputPrefix, *outputSuffix+".go")
	if err != nil {
		log.Fatalf("parsing package: %v", err)
	}
	return pkg
}

func (sfdc *SFDC) writeModelFiles(dir string, pkg string) {
	// Run generate for each object.
	for _, sobject := range sfdc.SObjects {
		if sfdc.isFiltered(sobject.Name) {
			continue
		}

		name := cleanname(sobject.Name)
		uri := sobject.URLs["describe"]
		sobjectDescription := &force.SObjectDescription{}
		err := sfdc.API.Get(uri, nil, sobjectDescription)
		if err != nil {
			log.Fatal(err)
		}

		var context = struct {
			PackageName        string
			TypeName           string
			SObject            *force.SObjectMetaData
			SObjectDescription *force.SObjectDescription
		}{
			PackageName:        pkg,
			TypeName:           name,
			SObject:            sobject,
			SObjectDescription: sobjectDescription,
		}

		var buf bytes.Buffer
		if err := generatedTmpl.Execute(&buf, context); err != nil {
			log.Fatalf("generating code: %v", err)
		}

		src, err := format.Source(buf.Bytes())
		if err != nil {
			// Should never happen, but can arise when developing this code.
			// The user can compile the output to see the error.
			log.Printf("warning: internal error: invalid Go generated: %s", err)
			log.Printf("warning: compile the package to analyze the error")
			src = buf.Bytes()
		}

		output := strings.ToLower(*outputPrefix + context.TypeName + *outputSuffix + ".go")
		outputPath := filepath.Join(dir, output)
		if err := ioutil.WriteFile(outputPath, src, 0644); err != nil {
			log.Fatalf("writing output: %s", err)
		}
	}
}

func writeCommonFile(dir string, pkg string) {
	var buf bytes.Buffer
	var context = struct {
		PackageName string
	}{
		PackageName: pkg,
	}
	if err := commonTmpl.Execute(&buf, context); err != nil {
		log.Fatalf("generating code: %v", err)
	}

	src, err := format.Source(buf.Bytes())
	if err != nil {
		// Should never happen, but can arise when developing this code.
		// The user can compile the output to see the error.
		log.Printf("warning: internal error: invalid Go generated: %s", err)
		log.Printf("warning: compile the package to analyze the error")
		src = buf.Bytes()
	}

	output := strings.ToLower(*outputPrefix + "common" + *outputSuffix + ".go")
	outputPath := filepath.Join(dir, output)
	if err := ioutil.WriteFile(outputPath, src, 0644); err != nil {
		log.Fatalf("writing output: %s", err)
	}
}

// SFDC describes a SFDC instance
type SFDC struct {
	Version             string
	API                 *force.ForceApi
	ResourcesURI        string
	Resources           map[string]string
	SObjects            map[string]*force.SObjectMetaData
	SObjectDescriptions map[string]*force.SObjectDescription
	MaxBatchSize        int64
	ObjectNames         []string
}

func (sfdc *SFDC) getResources() error {
	uri := fmt.Sprintf("/services/data/%v", sfdc.Version)
	return sfdc.API.Get(uri, nil, &sfdc.Resources)
}

func (sfdc *SFDC) getSObjects() error {
	uri := sfdc.Resources["sobjects"]
	list := &force.SObjectApiResponse{}
	err := sfdc.API.Get(uri, nil, list)
	if err != nil {
		return err
	}

	sfdc.MaxBatchSize = list.MaxBatchSize
	// The API doesn't return the list of sobjects in a map. Convert it.
	for _, object := range list.SObjects {
		sfdc.SObjects[object.Name] = object
	}

	return nil
}

func (sfdc *SFDC) isFiltered(object string) bool {
	for _, a := range sfdc.ObjectNames {
		if a == object {
			return false
		}
	}
	return true
}
