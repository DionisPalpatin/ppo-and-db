package interg_bl

import (
	integr_ut "github.com/DionisPalpatin/ppo-and-db/application/integration-tests/utils"
	bl "github.com/DionisPalpatin/ppo-and-db/application/internal/business-logic"
	"github.com/DionisPalpatin/ppo-and-db/application/internal/config"
	da "github.com/DionisPalpatin/ppo-and-db/application/internal/data-access"
	mylogger "github.com/DionisPalpatin/ppo-and-db/application/internal/logger"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"

	"database/sql"
	"fmt"
	"os"
)

type AuthServiceIntegrationTestSuite struct {
	suite.Suite
	db           *sql.DB
	oauthService bl.IOAuthService
}

func (s *AuthServiceIntegrationTestSuite) BeforeEach(t provider.T) {
	var err error
	pgInfo := integr_ut.PostgresInfo{
		Host:     "localhost",
		User:     "postgres",
		Password: "password",
		Port:     15423,
		DBName:   "NotebookAppIntTests",
	}

	s.db, err = integr_ut.InitDB(pgInfo)
	if err != nil {
		fmt.Println(err)
		return
	}

	if s.db == nil {
		return
	}

	text, err := os.ReadFile("./init.sql")
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = s.db.Exec(string(text))
	if err != nil {
		fmt.Println(err)
		return
	}

	dbconfig := config.DBConfigs{
		Host:       pgInfo.Host,
		User:       pgInfo.User,
		Password:   pgInfo.Password,
		Port:       pgInfo.Port,
		Name:       pgInfo.DBName,
		DB:         s.db,
		SchemaName: "interg_tests",
	}
	logger := mylogger.MyLogger{}
	logger.InitLogger("./logs/integration-test-auth.log", "debug")

	userRepo := da.NewUserRepository(&dbconfig, &logger)
	s.oauthService = bl.NewOAuthService(&userRepo, &logger)
}

func (s *AuthServiceIntegrationTestSuite) AfterEach(t provider.T) {
	_ = s.db.Close()
}

func (s *AuthServiceIntegrationTestSuite) TestRegisterUserSuccess(t provider.T) {
	t.Title("RegisterUser: Success")
	t.Tags("AuthServiceIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Registering user successfully", func(sCtx provider.StepCtx) {
		fio := "John Doe"
		login := "john.doe"
		password := "securepassword"

		user, err := s.oauthService.RegisterUser(fio, login, password)

		sCtx.Assert().NotNil(user)
		sCtx.Assert().NotNil(err)
		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
		sCtx.Assert().Equal(login, user.Login)
	})
}

func (s *AuthServiceIntegrationTestSuite) TestRegisterUserDuplicateLogin(t provider.T) {
	t.Title("RegisterUser: Duplicate Login")
	t.Tags("AuthServiceIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Attempting to register user with duplicate login", func(sCtx provider.StepCtx) {
		fio := "Jane Doe"
		login := "user1"
		password := "securepassword"

		user, err := s.oauthService.RegisterUser(fio, login, password)

		sCtx.Assert().Nil(user)
		sCtx.Assert().NotNil(err)
		sCtx.Assert().Equal(bl.UserExists, err.ErrNum)
	})
}

func (s *AuthServiceIntegrationTestSuite) TestSignInUserSuccess(t provider.T) {
	t.Title("SignInUser: Success")
	t.Tags("AuthServiceIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Signing in user successfully", func(sCtx provider.StepCtx) {
		login := "user1"
		password := "password1"

		token, err := s.oauthService.SignInUser(login, password)

		sCtx.Assert().NotNil(token)
		sCtx.Assert().Nil(err)
	})
}

func (s *AuthServiceIntegrationTestSuite) TestSignInUserInvalidCredentials(t provider.T) {
	t.Title("SignInUser: Invalid Credentials")
	t.Tags("AuthServiceIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Attempting to sign in with invalid credentials", func(sCtx provider.StepCtx) {
		login := "user1"
		password := "wrongpassword"

		token, err := s.oauthService.SignInUser(login, password)

		sCtx.Assert().Nil(token)
		sCtx.Assert().NotNil(err)
		sCtx.Assert().Equal(bl.AuthenticationError, err.ErrNum)
	})
}
