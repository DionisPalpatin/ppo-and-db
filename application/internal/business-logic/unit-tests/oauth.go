package unit_tests

import (
	"time"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"

	mocks "github.com/DionisPalpatin/ppo-and-db/application/internal/database/mocks/v2"
	mylogger "github.com/DionisPalpatin/ppo-and-db/application/internal/logger"

	bl "github.com/DionisPalpatin/ppo-and-db/application/internal/business-logic"
	myerror_builder "github.com/DionisPalpatin/ppo-and-db/application/internal/business-logic/builder"

	"github.com/DionisPalpatin/ppo-and-db/application/internal/models/builders"
)

type OAuthServiceTestSuite struct {
	suite.Suite
}

type extDependenciesOAuth struct {
	userRepo *mocks.MockUserRepository
	logger   *mylogger.MyLogger
}

func initExtDependenciesOAuth(t provider.T) *extDependenciesOAuth {
	mockRepo := mocks.NewMockUserRepository(t)

	f := &extDependenciesOAuth{
		userRepo: mockRepo,
		logger:   &mylogger.MyLogger{},
	}

	f.logger.InitLogger("./logs/unit-test-oauth.log", "debug")

	return f
}

func (s *OAuthServiceTestSuite) TestOAuthService_RegisterUser_OK(t provider.T) {
	t.Title("Register user: OK")
	t.Tags("OAuthService")
	t.Parallel()

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		df := initExtDependenciesOAuth(t)

		login := "newuser"
		fio := "New User"
		password := "securepassword"
		id := 0
		var regDate time.Time

		errNotFound := myerror_builder.NewMyErrorBuilder().WithErrNum(bl.NoSuchUser).Build()
		newUser := builders.NewUserBuilder().WithUserID(id).WithFio(fio).WithLogin(login).WithPassword(password).WithRegistrationDate(regDate).Build()
		errOk := myerror_builder.NewMyErrorBuilder().WithErrNum(bl.Ok).Build()

		df.userRepo.EXPECT().GetUserByLogin(login).Return(nil, errNotFound).Once()
		df.userRepo.EXPECT().AddUser(newUser).Return(errOk).Once()
		oAuthService := bl.NewOAuthService(df.userRepo, df.logger)

		user, err := oAuthService.RegisterUser(fio, login, password)

		sCtx.Assert().NotNil(user)
		sCtx.Assert().Equal(user.Login, login)
		sCtx.Assert().Equal(err.ErrNum, bl.Ok)
	})
}

func (s *OAuthServiceTestSuite) TestOAuthService_RegisterUser_UserExists(t provider.T) {
	t.Title("Register user: User already exists")
	t.Tags("OAuthService")
	t.Parallel()

	t.WithNewStep("User already exists", func(sCtx provider.StepCtx) {
		df := initExtDependenciesOAuth(t)

		login := "testuser"
		fio := "Test User"
		password := "password123"

		existingUser := builders.NewUserBuilder().WithLogin(login).Build()
		errExists := myerror_builder.NewMyErrorBuilder().WithErrNum(bl.Ok).Build()

		df.userRepo.EXPECT().GetUserByLogin(login).Return(existingUser, errExists).Once()
		oAuthService := bl.NewOAuthService(df.userRepo, df.logger)

		user, err := oAuthService.RegisterUser(fio, login, password)

		sCtx.Assert().Nil(user)
		sCtx.Assert().Equal(err.ErrNum, bl.UserExists)
	})
}

func (s *OAuthServiceTestSuite) TestOAuthService_SignInUser_OK(t provider.T) {
	t.Title("Sign in user: OK")
	t.Tags("OAuthService")
	t.Parallel()

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		df := initExtDependenciesOAuth(t)
		login := "testuser"
		password := "correctpassword"
		user := builders.NewUserBuilder().WithLogin(login).WithPassword(password).Build()
		errOk := myerror_builder.NewMyErrorBuilder().WithErrNum(bl.Ok).Build()

		df.userRepo.EXPECT().GetUserByLogin(login).Return(user, errOk).Once()

		oAuthService := bl.NewOAuthService(df.userRepo, df.logger)

		token, err := oAuthService.SignInUser(login, password)

		sCtx.Assert().NotNil(token)
		sCtx.Assert().Equal(err.ErrNum, bl.Ok)
	})
}

func (s *OAuthServiceTestSuite) TestOAuthService_SignInUser_InvalidPassword(t provider.T) {
	t.Title("Sign in user: Invalid password")
	t.Tags("OAuthService")
	t.Parallel()

	t.WithNewStep("Invalid password", func(sCtx provider.StepCtx) {
		df := initExtDependenciesOAuth(t)
		login := "testuser"
		password := "wrongpassword"
		user := builders.NewUserBuilder().WithLogin(login).WithPassword("correctpassword").Build()
		errOk := myerror_builder.NewMyErrorBuilder().WithErrNum(bl.Ok).Build()

		df.userRepo.EXPECT().GetUserByLogin(login).Return(user, errOk).Once()

		oAuthService := bl.NewOAuthService(df.userRepo, df.logger)

		token, err := oAuthService.SignInUser(login, password)

		sCtx.Assert().Nil(token)
		sCtx.Assert().Equal(err.ErrNum, bl.AuthenticationError)
	})
}

