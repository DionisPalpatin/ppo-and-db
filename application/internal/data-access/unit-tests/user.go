package da_unit

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"

	"github.com/DionisPalpatin/ppo-and-db/application/internal/config"
	mylogger "github.com/DionisPalpatin/ppo-and-db/application/internal/logger"

	bl "github.com/DionisPalpatin/ppo-and-db/application/internal/business-logic"
	da "github.com/DionisPalpatin/ppo-and-db/application/internal/data-access"

	"github.com/DionisPalpatin/ppo-and-db/application/internal/models"
	data_builders "github.com/DionisPalpatin/ppo-and-db/application/internal/models/builders"
)

var ZEROTIME time.Time

type UserRepositorySuite struct {
	suite.Suite

	db   *sql.DB
	mock sqlmock.Sqlmock
	repo da.UserRepository
	ctx  context.Context
}

func (s *UserRepositorySuite) BeforeEach(t provider.T) {
	var err error

	s.db, s.mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}

	confs := config.DBConfigs {
		DB: s.db,
		SchemaName: "test",
	}
	logger := mylogger.MyLogger{}
	logger.InitLogger("./logs/unit-tests-user.log", "debug")

	s.repo = da.NewUserRepository(&confs, &logger)
	s.ctx = context.Background()
}

func (s *UserRepositorySuite) AfterEach(t provider.T) {
	s.db.Close()
}

func TestUserRepositorySuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(UserRepositorySuite))
}

func (s *UserRepositorySuite) Test_UserRepository_GetUserByID_Success(t provider.T) {
	t.Title("GetUserByID: Success")
	t.Tags("UserRepository")

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		query := fmt.Sprintf(getUserByIDQuery, s.repo.DbConfigs.SchemaName)
		s.mock.ExpectQuery(query).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "fio", "registration_date", "login", "password", "role"}).
				AddRow(1, "John Doe", ZEROTIME, "johndoe", "password", bl.Reader))

		expectedUser := data_builders.NewUserBuilder().WithUserID(1).Build()

		user, err := s.repo.GetUserByID(1)

		sCtx.Assert().Equal(err.ErrNum, bl.Ok)
		sCtx.Assert().Equal(expectedUser.Id, user.Id)
	})
}

func (s *UserRepositorySuite) Test_UserRepository_GetUserByID_Failure(t provider.T) {
	t.Title("GetUserByID: Failure")
	t.Tags("UserRepository")

	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		query := fmt.Sprintf(getUserByIDQuery, s.repo.DbConfigs.SchemaName)
		s.mock.ExpectQuery(query).
			WithArgs(1).
			WillReturnError(sql.ErrConnDone)

		user, err := s.repo.GetUserByID(1)

		sCtx.Assert().Nil(user)
		sCtx.Assert().Equal(err.ErrNum, bl.DatabaseError)
	})
}

func (s *UserRepositorySuite) Test_UserRepository_GetUserByLogin_Success(t provider.T) {
	t.Title("GetUserByLogin: Success")
	t.Tags("UserRepository")

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		query := fmt.Sprintf(getUserByStringQuery, s.repo.DbConfigs.SchemaName)
		s.mock.ExpectQuery(query).
			WithArgs("johndoe").
			WillReturnRows(sqlmock.NewRows([]string{"id", "fio", "registration_date", "login", "password", "role"}).
				AddRow(1, "John Doe", ZEROTIME, "johndoe", "password", bl.Reader))

		expectedUser := data_builders.NewUserBuilder().WithUserID(1).WithLogin("johndoe").Build()

		user, err := s.repo.GetUserByLogin("johndoe")

		sCtx.Assert().Equal(err.ErrNum, bl.Ok)
		sCtx.Assert().Equal(expectedUser.Login, user.Login)
	})
}

func (s *UserRepositorySuite) Test_UserRepository_GetUserByLogin_Failure(t provider.T) {
	t.Title("GetUserByLogin: Failure")
	t.Tags("UserRepository")

	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		query := fmt.Sprintf(getUserByStringQuery, s.repo.DbConfigs.SchemaName)
		s.mock.ExpectQuery(query).
			WithArgs("johndoe").
			WillReturnError(sql.ErrConnDone)

		user, err := s.repo.GetUserByLogin("johndoe")

		sCtx.Assert().Nil(user)
		sCtx.Assert().Equal(bl.DatabaseError, err.ErrNum)
	})
}

func (s *UserRepositorySuite) Test_UserRepository_GetAllUsers_Success(t provider.T) {
	t.Title("GetAllUsers: Success")
	t.Tags("UserRepository")

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		query := fmt.Sprintf(getAllUsersQuery, s.repo.DbConfigs.SchemaName)
		s.mock.ExpectQuery(query).
			WillReturnRows(sqlmock.NewRows([]string{"id", "fio", "registration_date", "login", "password", "role"}).
				AddRow(1, "John Doe", ZEROTIME, "johndoe", "password", bl.Reader).
				AddRow(2, "Jane Smith", ZEROTIME, "janesmith", "password", bl.Admin))

		expectedUsers := []*models.User{
			data_builders.NewUserBuilder().WithUserID(1).WithFio("John Doe").WithLogin("johndoe").WithRole(bl.Reader).WithRegistrationDate(ZEROTIME).Build(),
			data_builders.NewUserBuilder().WithUserID(2).WithFio("Jane Smith").WithLogin("janesmith").WithRole(bl.Admin).WithRegistrationDate(ZEROTIME).Build(),
		}

		users, err := s.repo.GetAllUsers()

		sCtx.Assert().Equal(err.ErrNum, bl.Ok)
		sCtx.Assert().Equal(len(expectedUsers), len(users))
	})
}

func (s *UserRepositorySuite) Test_UserRepository_GetAllUsers_Failure(t provider.T) {
	t.Title("GetAllUsers: Failure")
	t.Tags("UserRepository")

	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		query := fmt.Sprintf(getAllUsersQuery, s.repo.DbConfigs.SchemaName)
		s.mock.ExpectQuery(query).
			WillReturnError(sql.ErrConnDone)

		users, err := s.repo.GetAllUsers()

		sCtx.Assert().Nil(users)
		sCtx.Assert().Equal(err.ErrNum, bl.DatabaseError)
	})
}

func (s *UserRepositorySuite) Test_UserRepository_AddUser_Success(t provider.T) {
	t.Title("AddUser: Success")
	t.Tags("UserRepository")

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		query := fmt.Sprintf(addUserQuery, s.repo.DbConfigs.SchemaName)

		s.mock.ExpectBegin()
		s.mock.ExpectExec(query).
			WithArgs("John Doe", ZEROTIME, "johndoe", "password", bl.Reader).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectCommit()	

		user := data_builders.NewUserBuilder().WithFio("John Doe").WithLogin("johndoe").WithPassword("password").WithRole(bl.Reader).WithRegistrationDate(ZEROTIME).Build()

		myErr := s.repo.AddUser(user)

		sCtx.Assert().Equal(myErr.ErrNum, bl.Ok)
	})
}

func (s *UserRepositorySuite) Test_UserRepository_AddUser_Failure(t provider.T) {
	t.Title("AddUser: Failure")
	t.Tags("UserRepository")

	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		query := fmt.Sprintf(addUserQuery, s.repo.DbConfigs.SchemaName)
		s.mock.ExpectQuery(query).
			WithArgs("John Doe", ZEROTIME, "johndoe", "password", bl.Reader).
			WillReturnError(sql.ErrConnDone)

		user := data_builders.NewUserBuilder().WithFio("John Doe").WithLogin("johndoe").WithRole(bl.Reader).WithRegistrationDate(ZEROTIME).Build()

		myErr := s.repo.AddUser(user)

		sCtx.Assert().Equal(myErr.ErrNum, bl.DatabaseError)
	})
}

func (s *UserRepositorySuite) Test_UserRepository_DeleteUser_Success(t provider.T) {
	t.Title("DeleteUser: Success")
	t.Tags("UserRepository")

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		query := fmt.Sprintf(deleteUserQuery, s.repo.DbConfigs.SchemaName)

		s.mock.ExpectBegin()
		s.mock.ExpectExec(query).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectCommit()	

		err := s.repo.DeleteUser(1)

		sCtx.Assert().Equal(err.ErrNum, bl.Ok)
	})
}

func (s *UserRepositorySuite) Test_UserRepository_DeleteUser_Failure(t provider.T) {
	t.Title("DeleteUser: Failure")
	t.Tags("UserRepository")

	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		query := fmt.Sprintf(deleteUserQuery, s.repo.DbConfigs.SchemaName)
		s.mock.ExpectExec(query).
			WithArgs(1).
			WillReturnError(sql.ErrConnDone)

		err := s.repo.DeleteUser(1)

		sCtx.Assert().Equal(err.ErrNum, bl.DatabaseError)
	})
}

func (s *UserRepositorySuite) Test_UserRepository_UpdateUser_Success(t provider.T) {
	t.Title("UpdateUser: Success")
	t.Tags("UserRepository")

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		query := fmt.Sprintf(updateUserQuery, s.repo.DbConfigs.SchemaName)

		s.mock.ExpectBegin()
		s.mock.ExpectExec(query).
			WithArgs("John Doe", ZEROTIME, "johndoe", "password", bl.Reader, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectCommit()
		
		user := data_builders.NewUserBuilder().WithUserID(1).WithFio("John Doe").WithLogin("johndoe").WithRole(bl.Reader).WithRegistrationDate(ZEROTIME).Build()

		err := s.repo.UpdateUser(user)

		sCtx.Assert().Equal(err.ErrNum, bl.Ok)
	})
}

func (s *UserRepositorySuite) Test_UserRepository_UpdateUser_Failure(t provider.T) {
	t.Title("UpdateUser: Failure")
	t.Tags("UserRepository")

	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		query := fmt.Sprintf(updateUserQuery, s.repo.DbConfigs.SchemaName)
		s.mock.ExpectExec(query).
			WithArgs("John Doe", ZEROTIME, "johndoe", "password", bl.Reader, 1).
			WillReturnError(sql.ErrConnDone)

		user := data_builders.NewUserBuilder().WithUserID(1).WithFio("John Doe").WithLogin("johndoe").WithRole(bl.Reader).WithRegistrationDate(ZEROTIME).Build()

		err := s.repo.UpdateUser(user)

		sCtx.Assert().Equal(err.ErrNum, bl.DatabaseError)
	})
}


