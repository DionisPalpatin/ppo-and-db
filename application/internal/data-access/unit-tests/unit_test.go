package da_unit

import (
	"testing"
	"os"

	"github.com/ozontech/allure-go/pkg/framework/suite"

)

func TestSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(UserRepositorySuite))
	suite.RunSuite(t, new(TeamRepositorySuite))
	suite.RunSuite(t, new(NoteRepositorySuite))
	suite.RunSuite(t, new(CollectionRepositorySuite))

	if t.Failed() {
		_ = os.Setenv("UNIT_SUCCESS", "0")
	} else if cur := os.Getenv("UNIT_SUCCESS"); cur == "" || cur == "1" {
		_ = os.Setenv("UNIT_SUCCESS", "1")
	}
}