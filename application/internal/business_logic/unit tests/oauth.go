package UnitTests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	v2 "github.com/DionisPalpatin/ppo-and-db/application/internal/database/mocks/v2"
	"github.com/DionisPalpatin/ppo-and-db/application/internal/models"
)

func TestRegisterUser(t *testing.T) {
	t.Run("SuccessRegisterUser", func(t *testing.T) {
		retOk := bl.CreateError(bl.Ok, nil, "")
		retErr := bl.CreateError(bl.ErrGetUserByLoginOrFio, bl.ErrGetUserByLoginError(), "GetUserByLogin")
		regUser := &models.User{Fio: "Test User", Login: "testuser", Password: "password", Role: bl.Reader}
		retUser := &models.User{}

		mockUserRepo := new(v2.MockIUserRepository)
		mockUserRepo.On("GetUserByLogin", "testuser").Return(retUser, retErr)
		mockUserRepo.On("AddUserToTeam", regUser).Return(retOk)

		osSrv := bl.OAuthService{}
		_, err := osSrv.RegisterUser("Test User", "testuser", "password")

		assert.NotNil(t, err)
		assert.Equal(t, err.ErrNum, bl.Ok)

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("ErrorUserExists", func(t *testing.T) {
		retOk := bl.CreateError(bl.Ok, nil, "")
		retUser := &models.User{Fio: "Test User", Login: "testuser", Password: "password", Role: bl.Reader}

		mockUserRepo := new(v2.MockIUserRepository)
		mockUserRepo.On("GetUserByLogin", "testuser").Return(retUser, retOk)

		osSrv := bl.OAuthService{}
		_, err := osSrv.RegisterUser("Test User", "testuser", "password")

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrRegisterUser, err.ErrNum)

		mockUserRepo.AssertExpectations(t)
	})
}

func TestSignInUser(t *testing.T) {
	t.Run("SuccessSignInUser", func(t *testing.T) {
		retOk := bl.CreateError(bl.Ok, nil, "")
		retUser := &models.User{Fio: "Test User", Login: "testuser", Password: "password", Role: bl.Reader}

		mockUserRepo := new(v2.MockIUserRepository)
		mockUserRepo.On("GetUserByLogin", "testuser").Return(retUser, retOk)

		osSrv := bl.OAuthService{}
		_, err := osSrv.SignInUser("testuser", "password", mockUserRepo)

		assert.NotNil(t, err)
		assert.Equal(t, err.ErrNum, bl.Ok)

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("ErrorUserNotFound", func(t *testing.T) {
		retErr := bl.CreateError(bl.ErrGetUserByLoginOrFio, bl.ErrGetUserByLoginError(), "GetUserByLogin")
		retUser := &models.User{}

		mockUserRepo := new(v2.MockIUserRepository)
		mockUserRepo.On("GetUserByLogin", "testuser").Return(retUser, retErr)

		osSrv := bl.OAuthService{}
		_, err := osSrv.SignInUser("testuser", "password", mockUserRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrGetUserByLoginOrFio, err.ErrNum)

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("ErrorIncorrectPassword", func(t *testing.T) {
		retOk := bl.CreateError(bl.Ok, nil, "")
		retUser := &models.User{Fio: "Test User", Login: "testuser", Password: "password", Role: bl.Reader}

		mockUserRepo := new(v2.MockIUserRepository)
		mockUserRepo.On("GetUserByLogin", "testuser").Return(retUser, retOk)

		osSrv := bl.OAuthService{}
		_, err := osSrv.SignInUser("testuser", "wrongpassword", mockUserRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrSignInUser, err.ErrNum)

		mockUserRepo.AssertExpectations(t)
	})
}
