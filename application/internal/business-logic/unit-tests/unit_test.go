package unit_tests

import (
	"testing"
	"os"

	"github.com/ozontech/allure-go/pkg/framework/suite"

)

func TestSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(UserServiceTestSuite))
	suite.RunSuite(t, new(TeamServiceTestSuite))
	suite.RunSuite(t, new(NoteServiceTestSuite))
	suite.RunSuite(t, new(CollectionServiceTestSuite))
	suite.RunSuite(t, new(OAuthServiceTestSuite))

	if t.Failed() {
		_ = os.Setenv("UNIT_SUCCESS", "0")
	} else if cur := os.Getenv("UNIT_SUCCESS"); cur == "" || cur == "1" {
		_ = os.Setenv("UNIT_SUCCESS", "1")
	}
}