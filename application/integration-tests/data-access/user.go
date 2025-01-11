package interg_da

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
)

type UserRepoIntegrationTestSuite struct {
	suite.Suite
	db       *sql.DB
	userRepo da.UserRepository

	idUserToDelete int
	idUserToUpdate int
}

func (s *UserRepoIntegrationTestSuite) BeforeEach(t provider.T) {
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

	s.userRepo = da.NewUserRepository(&dbconfig, &logger)

	s.idUserToUpdate = 9
	s.idUserToDelete = 10
}

func (s *UserRepoIntegrationTestSuite) AfterEach(t provider.T) {
	_ = s.db.Close()
}

func (s *UserRepoIntegrationTestSuite) TestGetUserByIDSuccess(t provider.T) {
	t.Title("GetUser: Success")
	t.Tags("UserRepoIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		user, err := s.userRepo.GetUserByID(1)

		sCtx.Assert().NotNil(user)
		sCtx.Assert().NotNil(err)
		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *UserRepoIntegrationTestSuite) TestGetUserByLoginSuccess(t provider.T) {
	t.Title("GetUser: Success")
	t.Tags("UserRepoIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		user, err := s.userRepo.GetUserByLogin("user1")

		sCtx.Assert().NotNil(user)
		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *UserRepoIntegrationTestSuite) TestGetAllUsersSuccess(t provider.T) {
	t.Title("GetAllUsers: Success")
	t.Tags("UserRepoIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		users, err := s.userRepo.GetAllUsers()

		sCtx.Assert().NotNil(users)
		sCtx.Assert().NotNil(err)
		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *UserRepoIntegrationTestSuite) TestUpdateUserSuccess(t provider.T) {
	t.Title("UpdateUser: Success")
	t.Tags("UserRepoIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		userToUpdate := data_builders.NewUserBuilder().WithUserID(s.idUserToUpdate).WithRole(bl.Reader).Build()

		err := s.userRepo.UpdateUser(userToUpdate)

		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *UserRepoIntegrationTestSuite) TestDeleteUserSuccess(t provider.T) {
	t.Title("DeleteUser: Success")
	t.Tags("UserRepoIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		err := s.userRepo.DeleteUser(s.idUserToDelete)

		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *UserRepoIntegrationTestSuite) TestAddUserSuccess(t provider.T) {
	t.Title("DeleteUser: Success")
	t.Tags("UserRepoIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		user := data_builders.NewUserBuilder().WithLogin("unique login").WithPassword("unique password").Build()

		err := s.userRepo.AddUser(user)

		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *UserRepoIntegrationTestSuite) TestAddUserDatabaseError(t provider.T) {
	t.Title("DeleteUser: user exists")
	t.Tags("UserRepoIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("User already exists", func(sCtx provider.StepCtx) {
		user := data_builders.NewUserBuilder().WithLogin("user1").Build()

		err := s.userRepo.AddUser(user)

		sCtx.Assert().Equal(bl.DatabaseError, err.ErrNum)
	})
}
