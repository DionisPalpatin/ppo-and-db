package interg_bl

import (
	integr_ut "github.com/DionisPalpatin/ppo-and-db/application/integration-tests/utils"
	bl "github.com/DionisPalpatin/ppo-and-db/application/internal/business-logic"
	"github.com/DionisPalpatin/ppo-and-db/application/internal/config"
	da "github.com/DionisPalpatin/ppo-and-db/application/internal/data-access"
	mylogger "github.com/DionisPalpatin/ppo-and-db/application/internal/logger"
	data_builders "github.com/DionisPalpatin/ppo-and-db/application/internal/models/builders"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"

	"database/sql"
	"fmt"
	"os"
)

type UserServiceIntegrationTest struct {
	suite.Suite
	db          *sql.DB
	userService bl.IUserService

	idUserToDelete int
	idUserToUpdate int
}

func (s *UserServiceIntegrationTest) BeforeEach(t provider.T) {
	var err error
	pgInfo := integr_ut.PostgresInfo{
		// Host:     "postgres-int-tests",
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
	logger.InitLogger("./logs/integration-test-user.log", "debug")

	userRepo := da.NewUserRepository(&dbconfig, &logger)
	s.userService = bl.NewUserService(&userRepo, &logger)

	s.idUserToUpdate = 8
	s.idUserToDelete = 12
}

func (s *UserServiceIntegrationTest) AfterEach(t provider.T) {
	_ = s.db.Close()
}

func (s *UserServiceIntegrationTest) TestGetUserSuccess(t provider.T) {
	t.Title("GetUser: Success")
	t.Tags("UserServiceIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		requester := data_builders.NewUserBuilder().WithRole(bl.Admin).Build()

		user, err := s.userService.GetUser(1, "", bl.SearchByID, requester)

		sCtx.Assert().NotNil(user)
		sCtx.Assert().NotNil(err)
		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *UserServiceIntegrationTest) TestGetUserNotFound(t provider.T) {
	t.Title("GetUser: User Not Found")
	t.Tags("UserServiceIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("User not found", func(sCtx provider.StepCtx) {
		requester := data_builders.NewUserBuilder().WithRole(bl.Admin).Build()

		user, err := s.userService.GetUser(999, "", bl.SearchByID, requester)

		sCtx.Assert().Nil(user)
		sCtx.Assert().NotEqual(bl.Ok, err.ErrNum)
		sCtx.Assert().Equal(bl.NoSuchUser, err.ErrNum)
	})
}

func (s *UserServiceIntegrationTest) TestGetAllUsersSuccess(t provider.T) {
	t.Title("GetAllUsers: Success")
	t.Tags("UserServiceIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		requester := data_builders.NewUserBuilder().WithRole(bl.Admin).Build()

		users, err := s.userService.GetAllUsers(requester)

		sCtx.Assert().NotNil(users)
		sCtx.Assert().NotNil(err)
		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *UserServiceIntegrationTest) TestGetAllUsersForbidden(t provider.T) {
	t.Title("GetAllUsers: Forbidden")
	t.Tags("UserServiceIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Forbidden access", func(sCtx provider.StepCtx) {
		requester := data_builders.NewUserBuilder().WithRole(bl.Reader).Build()

		users, err := s.userService.GetAllUsers(requester)

		sCtx.Assert().Nil(users)
		sCtx.Assert().NotEqual(bl.Ok, err.ErrNum)
		sCtx.Assert().Equal(bl.ErrAccessDenied, err.ErrNum)
	})
}

func (s *UserServiceIntegrationTest) TestUpdateUserSuccess(t provider.T) {
	t.Title("UpdateUser: Success")
	t.Tags("UserServiceIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		requester := data_builders.NewUserBuilder().WithRole(bl.Admin).Build()
		userToUpdate := data_builders.NewUserBuilder().WithUserID(s.idUserToUpdate).WithRole(bl.Reader).Build()

		err := s.userService.UpdateUser(requester, userToUpdate)

		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *UserServiceIntegrationTest) TestDeleteUserSuccess(t provider.T) {
	t.Title("DeleteUser: Success")
	t.Tags("UserServiceIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		requester := data_builders.NewUserBuilder().WithRole(bl.Admin).Build()

		err := s.userService.DeleteUser(requester, s.idUserToDelete)

		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}
