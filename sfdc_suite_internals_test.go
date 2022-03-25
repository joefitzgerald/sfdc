package sfdc

import (
	"testing"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

var suite spec.Suite

func init() {
	suite = spec.New("sfdc-internals", spec.Report(report.Terminal{}))
	suite("instance option", testInstanceOption)
	suite("response", testResponse)
}

func Test(t *testing.T) {
	suite.Run(t)
}
