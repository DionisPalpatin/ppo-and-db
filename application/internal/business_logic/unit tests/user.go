package UnitTests

import (
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"

	bl "github.com/DionisPalpatin/ppo-and-db/application/internal/business_logic"
	myerr_builder "github.com/DionisPalpatin/ppo-and-db/application/internal/business_logic/myerror_builder"
	mocks "github.com/DionisPalpatin/ppo-and-db/application/internal/database/mocks/v2"
	mylogger "github.com/DionisPalpatin/ppo-and-db/application/internal/logger"
	"github.com/DionisPalpatin/ppo-and-db/application/internal/models"
	"github.com/DionisPalpatin/ppo-and-db/application/internal/models/builders"
)

type UserServiceTestSuite struct {
	suite.Suite
}

type extDependenses struct {
	userRepo *mocks.MockUserRepository
	logger     *mylogger.MyLogger
}


func initExtDependeses(t provider.T) *extDependenses {
	mockArtistRepo := mocks.NewMockUserRepository(t)

	f := &extDependenses{
		userRepo: mockArtistRepo,
		logger:     &mylogger.MyLogger{},
	}

	f.logger.InitLogger("./tests_user_logs.log", "debug")

	return f
}

func (s *UserServiceTestSuite) TestUserService_GetUser_OK(t provider.T) {
	t.Title("Get user: OK")
	t.Tags("UserService")
	t.Parallel()

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		df := initExtDependeses(t)
		in_user := builders.NewUserBuilder().Build()
		out_user := builders.NewUserBuilder().WithUserID(in_user.Id)
		req_user := builders.NewUserBuilder().WithRole(bl.Admin)
		myerr := bl.NewMyErrorBuilder().WithErrNum(bl.Ok)

		df.userRepo.EXPECT().GetUserByID(in_user.Id).Return(out_user.User, myerr.MyError).Once()
		userService := bl.NewUserService(df.userRepo, df.logger)

		user, err := userService.GetUser(in_user.Id, "", bl.SearchByID, req_user.User)

		sCtx.Assert().Equal(err.ErrNum, bl.Ok)
        sCtx.Assert().Equal(user.Id, out_user.User.Id)
	})
}

func (s *UserServiceTestSuite) TestUserService_GetUser_AccessDenied(t provider.T) {
	t.Title("Get user: Access Denied")
	t.Tags("UserService")
	t.Parallel()

	t.WithNewStep("Access Denied", func(sCtx provider.StepCtx) {
		df := initExtDependeses(t)
		in_user := builders.NewUserBuilder().Build()
		req_user := builders.NewUserBuilder().WithRole(bl.Reader)
		myerr := bl.NewMyErrorBuilder().WithErrNum(bl.ErrAccessDenied)

		df.userRepo.EXPECT().GetUserByID(in_user.Id).Return(nil, myerr.MyError).Once()
		userService := bl.NewUserService(df.userRepo, df.logger)

		user, err := userService.GetUser(in_user.Id, "", bl.SearchByID, req_user.User)

		sCtx.Assert().Nil(user)
		sCtx.Assert().Equal(err.ErrNum, bl.ErrAccessDenied)
	})
}

func (s *UserServiceTestSuite) TestUserService_GetAllUsers_OK(t provider.T) {
	t.Title("Get all users: OK")
	t.Tags("UserService")
	t.Parallel()

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		df := initExtDependeses(t)
		req_user := builders.NewUserBuilder().WithRole(bl.Admin)  // Админская роль
		users := []*models.User{
			builders.NewUserBuilder().Build(),
			builders.NewUserBuilder().Build(),
		}
		myerr := bl.NewMyErrorBuilder().WithErrNum(bl.Ok)

		df.userRepo.EXPECT().GetAllUsers().Return(users, myerr.MyError).Once()
		userService := bl.NewUserService(df.userRepo, df.logger)

		resultUsers, err := userService.GetAllUsers(req_user.User)

		sCtx.Assert().Equal(err.ErrNum, bl.Ok)
		sCtx.Assert().Equal(len(resultUsers), len(users))
	})
}

func (s *UserServiceTestSuite) TestUserService_GetAllUsers_AccessDenied(t provider.T) {
	t.Title("Get all users: Access Denied")
	t.Tags("UserService")
	t.Parallel()

	t.WithNewStep("Access Denied", func(sCtx provider.StepCtx) {
		df := initExtDependeses(t)
		req_user := builders.NewUserBuilder().WithRole(bl.Reader)  // Несоответствующие права доступа
		myerr := bl.NewMyErrorBuilder().WithErrNum(bl.ErrAccessDenied)

		df.userRepo.EXPECT().GetAllUsers().Return(nil, myerr.MyError).Once()
		userService := bl.NewUserService(df.userRepo, df.logger)

		users, err := userService.GetAllUsers(req_user.User)

		sCtx.Assert().Nil(users)
		sCtx.Assert().Equal(err.ErrNum, bl.ErrAccessDenied)
	})
}

func (s *UserServiceTestSuite) TestUserService_UpdateUser_OK(t provider.T) {
	t.Title("Update user: OK")
	t.Tags("UserService")
	t.Parallel()

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		df := initExtDependeses(t)
		in_user := builders.NewUserBuilder().Build()
		req_user := builders.NewUserBuilder().WithRole(bl.Admin)  // Админская роль
		myerr := bl.NewMyErrorBuilder().WithErrNum(bl.Ok)

		df.userRepo.EXPECT().UpdateUser(in_user).Return(myerr.MyError).Once()
		userService := bl.NewUserService(df.userRepo, df.logger)

		err := userService.UpdateUser(req_user.User, in_user)

		sCtx.Assert().Equal(err.ErrNum, bl.Ok)
	})
}

func (s *UserServiceTestSuite) TestUserService_UpdateUser_AccessDenied(t provider.T) {
	t.Title("Update user: Access Denied")
	t.Tags("UserService")
	t.Parallel()

	t.WithNewStep("Access Denied", func(sCtx provider.StepCtx) {
		df := initExtDependeses(t)
		in_user := builders.NewUserBuilder().Build()
		req_user := builders.NewUserBuilder().WithRole(bl.Reader)
		myerr := bl.NewMyErrorBuilder().WithErrNum(bl.ErrAccessDenied)

		df.userRepo.EXPECT().UpdateUser(in_user).Return(myerr.MyError).Once()
		userService := bl.NewUserService(df.userRepo, df.logger)

		err := userService.UpdateUser(req_user.User, in_user)

		sCtx.Assert().Equal(err.ErrNum, bl.ErrAccessDenied)
	})
}

func (s *UserServiceTestSuite) TestUserService_DeleteUser_OK(t provider.T) {
	t.Title("Delete user: OK")
	t.Tags("UserService")
	t.Parallel()

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		df := initExtDependeses(t)
		in_user := builders.NewUserBuilder().Build()
		req_user := builders.NewUserBuilder().WithRole(bl.Admin)  // Админская роль
		myerr := bl.NewMyErrorBuilder().WithErrNum(bl.Ok)

		df.userRepo.EXPECT().DeleteUser(in_user.Id).Return(myerr.MyError).Once()
		userService := bl.NewUserService(df.userRepo, df.logger)

		err := userService.DeleteUser(req_user.User, in_user.Id)

		sCtx.Assert().Equal(err.ErrNum, bl.Ok)
	})
}

func (s *UserServiceTestSuite) TestUserService_DeleteUser_AccessDenied(t provider.T) {
	t.Title("Delete user: Access Denied")
	t.Tags("UserService")
	t.Parallel()

	t.WithNewStep("Access Denied", func(sCtx provider.StepCtx) {
		df := initExtDependeses(t)
		in_user := builders.NewUserBuilder().Build()
		req_user := builders.NewUserBuilder().WithRole(bl.Reader)  // Несоответствующие права доступа
		myerr := bl.NewMyErrorBuilder().WithErrNum(bl.ErrAccessDenied)

		df.userRepo.EXPECT().DeleteUser(in_user.Id).Return(myerr.MyError).Once()
		userService := bl.NewUserService(df.userRepo, df.logger)

		err := userService.DeleteUser(req_user.User, in_user.Id)

		sCtx.Assert().Equal(err.ErrNum, bl.ErrAccessDenied)  // Ожидаем ошибку доступа
	})
}
