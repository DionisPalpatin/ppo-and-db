package UnitTests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/business_logic"
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/business_logic/UnitTests/mocks"
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/models"
)

func TestRegisterUser(t *testing.T) {
	t.Run("SuccessRegisterUser", func(t *testing.T) {
		retOk := bl.CreateError(bl.AllIsOk, nil, "")
		retErr := bl.CreateError(bl.ErrGetUserByLogin, bl.ErrGetUserByLoginError(), "GetUserByLogin")
		regUser := &models.User{Fio: "Test User", Login: "testuser", Password: "password", Role: bl.Reader}
		retUser := &models.User{}

		mockUserRepo := new(mocks.MockIUserRepository)
		mockUserRepo.On("GetUserByLogin", "testuser").Return(retUser, retErr)
		mockUserRepo.On("AddUserToTeam", regUser).Return(retOk)

		osSrv := bl.OAuthService{}
		_, err := osSrv.RegisterUser("Test User", "testuser", "password", mockUserRepo)

		assert.NotNil(t, err)
		assert.Equal(t, err.ErrNum, bl.AllIsOk)

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("ErrorUserExists", func(t *testing.T) {
		retOk := bl.CreateError(bl.AllIsOk, nil, "")
		retUser := &models.User{Fio: "Test User", Login: "testuser", Password: "password", Role: bl.Reader}

		mockUserRepo := new(mocks.MockIUserRepository)
		mockUserRepo.On("GetUserByLogin", "testuser").Return(retUser, retOk)

		osSrv := bl.OAuthService{}
		_, err := osSrv.RegisterUser("Test User", "testuser", "password", mockUserRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrRegisterUser, err.ErrNum)

		mockUserRepo.AssertExpectations(t)
	})
}

func TestSignInUser(t *testing.T) {
	t.Run("SuccessSignInUser", func(t *testing.T) {
		retOk := bl.CreateError(bl.AllIsOk, nil, "")
		retUser := &models.User{Fio: "Test User", Login: "testuser", Password: "password", Role: bl.Reader}

		mockUserRepo := new(mocks.MockIUserRepository)
		mockUserRepo.On("GetUserByLogin", "testuser").Return(retUser, retOk)

		osSrv := bl.OAuthService{}
		_, err := osSrv.SignInUser("testuser", "password", mockUserRepo)

		assert.NotNil(t, err)
		assert.Equal(t, err.ErrNum, bl.AllIsOk)

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("ErrorUserNotFound", func(t *testing.T) {
		retErr := bl.CreateError(bl.ErrGetUserByLogin, bl.ErrGetUserByLoginError(), "GetUserByLogin")
		retUser := &models.User{}

		mockUserRepo := new(mocks.MockIUserRepository)
		mockUserRepo.On("GetUserByLogin", "testuser").Return(retUser, retErr)

		osSrv := bl.OAuthService{}
		_, err := osSrv.SignInUser("testuser", "password", mockUserRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrGetUserByLogin, err.ErrNum)

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("ErrorIncorrectPassword", func(t *testing.T) {
		retOk := bl.CreateError(bl.AllIsOk, nil, "")
		retUser := &models.User{Fio: "Test User", Login: "testuser", Password: "password", Role: bl.Reader}

		mockUserRepo := new(mocks.MockIUserRepository)
		mockUserRepo.On("GetUserByLogin", "testuser").Return(retUser, retOk)

		osSrv := bl.OAuthService{}
		_, err := osSrv.SignInUser("testuser", "wrongpassword", mockUserRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrSignInUser, err.ErrNum)

		mockUserRepo.AssertExpectations(t)
	})
}
