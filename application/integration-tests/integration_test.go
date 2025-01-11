package integration_tests

import (
	integr_bl "github.com/DionisPalpatin/ppo-and-db/application/integration-tests/business-logic"
	integr_da "github.com/DionisPalpatin/ppo-and-db/application/integration-tests/data-access"
	integr_ut "github.com/DionisPalpatin/ppo-and-db/application/integration-tests/utils"

	"fmt"
	"os"
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func BeforeAll() {
	pgInfo := integr_ut.PostgresInfo{
		Host:     "localhost",
		User:     "postgres",
		Password: "password",
		Port:     15423,
		DBName:   "NotebookAppIntTests",
	}

	db, err := integr_ut.InitDB(pgInfo)
	if err != nil {
		fmt.Println(err)
		return
	}

	text, err := os.ReadFile("./init.sql")
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = db.Exec(string(text))
	if err != nil {
		fmt.Println(err)
		return
	}

	db.Close()
	fmt.Println("Соединение с БД в BeforeAll закрыто!")
}

func TestIntegrationRunner(t *testing.T) {
	BeforeAll()

	suite.RunSuite(t, new(integr_bl.UserServiceIntegrationTest))
	suite.RunSuite(t, new(integr_bl.TeamServiceIntegrationTest))
	suite.RunSuite(t, new(integr_bl.NoteServiceIntegrationTest))
	suite.RunSuite(t, new(integr_bl.CollectionServiceIntegrationTest))
	suite.RunSuite(t, new(integr_bl.AuthServiceIntegrationTestSuite))

	suite.RunSuite(t, new(integr_da.UserRepoIntegrationTestSuite))
	suite.RunSuite(t, new(integr_da.TeamRepoIntegrationTestSuite))
	suite.RunSuite(t, new(integr_da.NoteRepoIntegrationTestSuite))
	suite.RunSuite(t, new(integr_da.CollectionRepoIntegrationTest))

	if t.Failed() {
		_ = os.Setenv("INTEGRATION_SUCCESS", "0")
	} else if cur := os.Getenv("INTEGRATION_SUCCESS"); cur == "" || cur == "1" {
		_ = os.Setenv("INTEGRATION_SUCCESS", "1")
	}
}
