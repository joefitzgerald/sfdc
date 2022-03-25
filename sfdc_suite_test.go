package sfdc_test

import (
	"testing"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

var suite spec.Suite

func init() {
	suite = spec.New("", spec.Report(report.Terminal{}))
	suite("instance", testInstance)
	suite("entity", testEntity)
	suite("auth options", testAuthOptions)
	suite("fields", testFields)
}

func Test(t *testing.T) {
	suite.Run(t)
}
