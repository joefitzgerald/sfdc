package main

import (
	"flag"
	"log"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Version       string `default:"v37.0"`
	ClientID      string `required:"true"`
	ClientSecret  string `required:"true"`
	Username      string `required:"true"`
	Password      string `required:"true"`
	SecurityToken string `required:"true"`
	Environment   string `default:"Production"`
}

var (
	objectNames  = flag.String("object", "", "comma-separated list of SFDC object names; must be set")
	outputPrefix = flag.String("prefix", "", "prefix to be added to the output file")
	outputSuffix = flag.String("suffix", "_sfdc", "suffix to be added to the output file")
)

func main() {
	flag.Parse()
	var c config
	err := envconfig.Process("sfdc", &c)
	if err != nil {
		log.Fatal(err.Error())
	}
	if len(*objectNames) == 0 {
		log.Fatalf("the flag -object must be set")
	}

	sfdc := New(c, *objectNames)

	err = sfdc.getResources()
	if err != nil {
		log.Fatal(err)
	}
	err = sfdc.getSObjects()
	if err != nil {
		log.Fatal(err)
	}

	// Only one directory at a time can be processed, and the default is ".".
	dir := "."
	if args := flag.Args(); len(args) == 1 {
		dir = args[0]
	} else if len(args) > 1 {
		log.Fatalf("only one directory at a time")
	}

	pkg := getPackageName(dir)
	sfdc.writeCommonFile(dir, pkg)
	sfdc.writeModelFiles(dir, pkg)
}
