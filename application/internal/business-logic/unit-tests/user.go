package unit_tests

import (
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"

	mocks "github.com/DionisPalpatin/ppo-and-db/application/internal/database/mocks/v2"
	mylogger "github.com/DionisPalpatin/ppo-and-db/application/internal/logger"

	bl "github.com/DionisPalpatin/ppo-and-db/application/internal/business-logic"
	myerror_builder "github.com/DionisPalpatin/ppo-and-db/application/internal/business-logic/builder"

	"github.com/DionisPalpatin/ppo-and-db/application/internal/models"
	"github.com/DionisPalpatin/ppo-and-db/application/internal/models/builders"
)

type UserServiceTestSuite struct {
	suite.Suite
}

type extDependensesUser struct {
	repo   *mocks.MockUserRepository
	logger *mylogger.MyLogger
}

func initExtDependenciesUser(t provider.T) *extDependensesUser {
	mockRepo := mocks.NewMockUserRepository(t)

	f := &extDependensesUser{
		repo:   mockRepo,
		logger: &mylogger.MyLogger{},
	}

	f.logger.InitLogger("./logs/unit-test-user.log", "debug")

	return f
}

func (s *UserServiceTestSuite) TestUserService_GetUser_OK(t provider.T) {
	t.Title("Get user: OK")
	t.Tags("UserService")
	t.Parallel()

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		df := initExtDependenciesUser(t)
		inUser := builders.NewUserBuilder().Build()
		outUser := builders.NewUserBuilder().WithUserID(inUser.Id)
		reqUser := builders.NewUserBuilder().WithRole(bl.Admin)
		myerr := myerror_builder.NewMyErrorBuilder().WithErrNum(bl.Ok).Build()

		df.repo.EXPECT().GetUserByID(inUser.Id).Return(outUser.User, myerr).Once()
		userService := bl.NewUserService(df.repo, df.logger)

		user, err := userService.GetUser(inUser.Id, "", bl.SearchByID, reqUser.User)

		sCtx.Assert().Equal(err.ErrNum, bl.Ok)
		sCtx.Assert().Equal(user.Id, outUser.User.Id)
	})
}

func (s *UserServiceTestSuite) TestUserService_GetUser_AccessDenied(t provider.T) {
	t.Title("Get user: Access Denied")
	t.Tags("UserService")
	t.Parallel()

	t.WithNewStep("Access Denied", func(sCtx provider.StepCtx) {
		df := initExtDependenciesUser(t)
		inUser := builders.NewUserBuilder().Build()
		reqUser := builders.NewUserBuilder().WithRole(bl.Reader)
		myerr := myerror_builder.NewMyErrorBuilder().WithErrNum(bl.ErrAccessDenied).Build()

		df.repo.EXPECT().GetUserByID(inUser.Id).Return(nil, myerr).Once()
		userService := bl.NewUserService(df.repo, df.logger)

		user, err := userService.GetUser(inUser.Id, "", bl.SearchByID, reqUser.User)

		sCtx.Assert().Nil(user)
		sCtx.Assert().Equal(err.ErrNum, bl.ErrAccessDenied)
	})
}

func (s *UserServiceTestSuite) TestUserService_GetAllUsers_OK(t provider.T) {
	t.Title("Get all users: OK")
	t.Tags("UserService")
	t.Parallel()

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		df := initExtDependenciesUser(t)
		reqUser := builders.NewUserBuilder().WithRole(bl.Admin) // Админская роль
		users := []*models.User{
			builders.NewUserBuilder().Build(),
			builders.NewUserBuilder().Build(),
		}
		myerr := myerror_builder.NewMyErrorBuilder().WithErrNum(bl.Ok).Build()

		df.repo.EXPECT().GetAllUsers().Return(users, myerr).Once()
		userService := bl.NewUserService(df.repo, df.logger)

		resultUsers, err := userService.GetAllUsers(reqUser.User)

		sCtx.Assert().Equal(err.ErrNum, bl.Ok)
		sCtx.Assert().Equal(len(resultUsers), len(users))
	})
}

func (s *UserServiceTestSuite) TestUserService_GetAllUsers_AccessDenied(t provider.T) {
	t.Title("Get all users: Access Denied")
	t.Tags("UserService")
	t.Parallel()

	t.WithNewStep("Access Denied", func(sCtx provider.StepCtx) {
		df := initExtDependenciesUser(t)
		reqUser := builders.NewUserBuilder().WithRole(bl.Reader)

		userService := bl.NewUserService(df.repo, df.logger)

		users, err := userService.GetAllUsers(reqUser.User)

		sCtx.Assert().Nil(users)
		sCtx.Assert().Equal(err.ErrNum, bl.ErrAccessDenied)
	})
}

func (s *UserServiceTestSuite) TestUserService_UpdateUser_OK(t provider.T) {
	t.Title("Update user: OK")
	t.Tags("UserService")
	t.Parallel()

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		df := initExtDependenciesUser(t)
		inUser := builders.NewUserBuilder().Build()
		reqUser := builders.NewUserBuilder().WithRole(bl.Admin)
		myerr := myerror_builder.NewMyErrorBuilder().WithErrNum(bl.Ok).Build()

		df.repo.EXPECT().UpdateUser(inUser).Return(myerr).Once()
		userService := bl.NewUserService(df.repo, df.logger)

		err := userService.UpdateUser(reqUser.User, inUser)

		sCtx.Assert().Equal(err.ErrNum, bl.Ok)
	})
}

func (s *UserServiceTestSuite) TestUserService_UpdateUser_AccessDenied(t provider.T) {
	t.Title("Update user: Access Denied")
	t.Tags("UserService")
	t.Parallel()

	t.WithNewStep("Access Denied", func(sCtx provider.StepCtx) {
		df := initExtDependenciesUser(t)
		inUser := builders.NewUserBuilder().Build()
		reqUser := builders.NewUserBuilder().WithRole(bl.Reader)
		
		userService := bl.NewUserService(df.repo, df.logger)

		err := userService.UpdateUser(reqUser.User, inUser)

		sCtx.Assert().Equal(err.ErrNum, bl.ErrAccessDenied)
	})
}

func (s *UserServiceTestSuite) TestUserService_DeleteUser_OK(t provider.T) {
	t.Title("Delete user: OK")
	t.Tags("UserService")
	t.Parallel()

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		df := initExtDependenciesUser(t)
		inUser := builders.NewUserBuilder().Build()
		reqUser := builders.NewUserBuilder().WithRole(bl.Admin)
		myerr := myerror_builder.NewMyErrorBuilder().WithErrNum(bl.Ok).Build()

		df.repo.EXPECT().DeleteUser(inUser.Id).Return(myerr).Once()
		userService := bl.NewUserService(df.repo, df.logger)

		err := userService.DeleteUser(reqUser.User, inUser.Id)

		sCtx.Assert().Equal(err.ErrNum, bl.Ok)
	})
}

func (s *UserServiceTestSuite) TestUserService_DeleteUser_AccessDenied(t provider.T) {
	t.Title("Delete user: Access Denied")
	t.Tags("UserService")
	t.Parallel()

	t.WithNewStep("Access Denied", func(sCtx provider.StepCtx) {
		df := initExtDependenciesUser(t)
		inUser := builders.NewUserBuilder().Build()
		reqUser := builders.NewUserBuilder().WithRole(bl.Reader)
		// myerr := myerror_builder.NewMyErrorBuilder().WithErrNum(bl.ErrAccessDenied).Build()

		// df.repo.EXPECT().DeleteUser(inUser.Id).Return(myerr).Once()
		userService := bl.NewUserService(df.repo, df.logger)

		err := userService.DeleteUser(reqUser.User, inUser.Id)

		sCtx.Assert().Equal(err.ErrNum, bl.ErrAccessDenied)
	})
}


