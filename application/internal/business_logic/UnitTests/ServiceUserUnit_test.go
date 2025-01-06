package UnitTests

import (
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/suite"

	"github.com/stretchr/testify/assert"

	v1 "github.com/DionisPalpatin/ppo-and-db/application/internal/database/mocks/v1"
	"github.com/DionisPalpatin/ppo-and-db/application/internal/models"
)

type userTestSuite struct {
	suite.Suite
}

func TestGetUser(t *testing.T) {
	t.Run("SuccessGetUserByID", func(t *testing.T) {
		retErr := bl.CreateError(bl.Ok, "", nil)
		retUser := &models.User{Role: bl.Admin}
		reqUser := &models.User{Role: bl.Admin}

		mockRepo := new(v1.MockIUserRepository)
		mockRepo.On("GetUserByID", 1).Return(retUser, retErr)

		usSrv := bl.UserService{}
		_, err := usSrv.GetUser(1, "", bl.SearchByID, reqUser, mockRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.Ok, err.ErrNum)

		mockRepo.AssertExpectations(t)
	})

	t.Run("SuccessGetUserByLogin", func(t *testing.T) {
		retErr := bl.CreateError(bl.Ok, "", nil)
		retUser := &models.User{Role: bl.Admin}
		reqUser := &models.User{Role: bl.Admin}

		mockRepo := new(v1.MockIUserRepository)
		mockRepo.On("GetUserByLogin", "testuser").Return(retUser, retErr)

		usSrv := bl.UserService{}
		_, err := usSrv.GetUser(0, "testuser", bl.SearchByString, reqUser, mockRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.Ok, err.ErrNum)

		mockRepo.AssertExpectations(t)
	})

	t.Run("ErrorRequesterNotAdmin", func(t *testing.T) {
		reqUser := &models.User{Role: bl.Reader}

		mockRepo := new(v1.MockIUserRepository)

		usSrv := bl.UserService{}
		_, err := usSrv.GetUser(1, "", bl.SearchByID, reqUser, mockRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrAccessDenied, err.ErrNum)

		mockRepo.AssertExpectations(t)
	})

	t.Run("ErrorUnknownSearchParameter", func(t *testing.T) {
		reqUser := &models.User{Role: bl.Admin}

		mockRepo := new(v1.MockIUserRepository)

		usSrv := bl.UserService{}
		_, err := usSrv.GetUser(-1, "", 100, reqUser, mockRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrSearchParameter, err.ErrNum)

		mockRepo.AssertExpectations(t)
	})

	t.Run("ErrorGetUserByID", func(t *testing.T) {
		retErr := bl.CreateError(bl.ErrGetUserByID, "GetUserByID", bl.ErrGetUserByIDError())
		retUser := &models.User{Role: bl.Admin}
		reqUser := &models.User{Role: bl.Admin}

		mockRepo := new(v1.MockIUserRepository)
		mockRepo.On("GetUserByID", -1).Return(retUser, retErr)

		usSrv := bl.UserService{}
		_, err := usSrv.GetUser(-1, "", bl.SearchByID, reqUser, mockRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrGetUserByID, err.ErrNum)

		mockRepo.AssertExpectations(t)
	})

	t.Run("ErrorGetUserByLogin", func(t *testing.T) {
		retErr := bl.CreateError(bl.ErrGetUserByLoginOrFio, "GetUserByLogin", bl.ErrGetUserByLoginError())
		retUser := &models.User{Role: bl.Admin}
		reqUser := &models.User{Role: bl.Admin}

		mockRepo := new(v1.MockIUserRepository)
		mockRepo.On("GetUserByLogin", "").Return(retUser, retErr)

		usSrv := bl.UserService{}
		_, err := usSrv.GetUser(0, "", bl.SearchByString, reqUser, mockRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrGetUserByLoginOrFio, err.ErrNum)

		mockRepo.AssertExpectations(t)
	})
}

func TestGetAllUsers(t *testing.T) {
	t.Run("SuccessGetAllUsers", func(t *testing.T) {
		retOk := bl.CreateError(bl.Ok, "", nil)
		retUsers := []*models.User{&models.User{}}
		reqUser := &models.User{Role: bl.Admin}

		mockRepo := new(v1.MockIUserRepository)
		mockRepo.On("GetAllUsers").Return(retUsers, retOk)

		usSrv := bl.UserService{}
		_, err := usSrv.GetAllUsers(reqUser, mockRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.Ok, err.ErrNum)

		mockRepo.AssertExpectations(t)
	})

	t.Run("ErrorRequesterNotAdmin", func(t *testing.T) {
		reqUser := &models.User{Role: bl.Reader}

		mockRepo := new(v1.MockIUserRepository)

		usSrv := bl.UserService{}
		_, err := usSrv.GetAllUsers(reqUser, mockRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrAccessDenied, err.ErrNum)

		mockRepo.AssertExpectations(t)
	})

	t.Run("ErrorGetAllUsers", func(t *testing.T) {
		retErr := bl.CreateError(bl.ErrGetAllUsers, "GetAllUsers", bl.ErrGetAllUsersError())
		reqUser := &models.User{Role: bl.Admin}
		retUsers := []*models.User{&models.User{}}

		mockRepo := new(v1.MockIUserRepository)
		mockRepo.On("GetAllUsers").Return(retUsers, retErr)

		usSrv := bl.UserService{}
		_, err := usSrv.GetAllUsers(reqUser, mockRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrGetAllUsers, err.ErrNum)

		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateUser(t *testing.T) {
	t.Run("SuccessUpdateUser", func(t *testing.T) {
		retOk := bl.CreateError(bl.Ok, "", nil)
		updateUser := &models.User{Id: 2, Role: bl.Reader}
		reqUser := &models.User{Role: bl.Admin}

		mockRepo := new(v1.MockIUserRepository)
		mockRepo.On("UpdateUser", updateUser).Return(retOk)

		usSrv := bl.UserService{}
		err := usSrv.UpdateUser(reqUser, updateUser, mockRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.Ok, err.ErrNum)

		mockRepo.AssertExpectations(t)
	})

	t.Run("ErrorRequesterNotAdmin", func(t *testing.T) {
		updateUser := &models.User{Id: 2, Role: bl.Reader}
		reqUser := &models.User{Role: bl.Reader}

		mockRepo := new(v1.MockIUserRepository)

		usSrv := bl.UserService{}
		err := usSrv.UpdateUser(reqUser, updateUser, mockRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrAccessDenied, err.ErrNum)

		mockRepo.AssertExpectations(t)
	})

	t.Run("ErrorUpdateUser", func(t *testing.T) {
		retErr := bl.CreateError(bl.ErrUpdateUser, "", bl.ErrUpdateUserError())
		updateUser := &models.User{Id: 2, Role: bl.Reader}
		reqUser := &models.User{Role: bl.Admin}

		mockRepo := new(v1.MockIUserRepository)
		mockRepo.On("UpdateUser", updateUser).Return(retErr)

		usSrv := bl.UserService{}
		err := usSrv.UpdateUser(reqUser, updateUser, mockRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrUpdateUser, err.ErrNum)

		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteUser(t *testing.T) {
	t.Run("SuccessDeleteUser", func(t *testing.T) {
		retOk := bl.CreateError(bl.Ok, "", nil)
		reqUser := &models.User{Role: bl.Admin}

		mockRepo := new(v1.MockIUserRepository)
		mockRepo.On("DeleteUser", 2).Return(retOk)

		usSrv := bl.UserService{}
		err := usSrv.DeleteUser(reqUser, 2, mockRepo)

		assert.NotNil(t, err)
		assert.Equal(t, err.ErrNum, bl.Ok)

		mockRepo.AssertExpectations(t)
	})

	t.Run("ErrorRequesterNotAdmin", func(t *testing.T) {
		reqUser := &models.User{Role: bl.Reader}

		mockRepo := new(v1.MockIUserRepository)

		usSrv := bl.UserService{}
		err := usSrv.DeleteUser(reqUser, 2, mockRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrAccessDenied, err.ErrNum)

		mockRepo.AssertExpectations(t)
	})

	t.Run("ErrorDeleteUser", func(t *testing.T) {
		retErr := bl.CreateError(bl.ErrDeleteUser, "", bl.ErrDeleteUserError())
		reqUser := &models.User{Role: bl.Admin}

		mockRepo := new(v1.MockIUserRepository)
		mockRepo.On("DeleteUser", 2).Return(retErr)

		usSrv := bl.UserService{}
		err := usSrv.DeleteUser(reqUser, 2, mockRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrDeleteUser, err.ErrNum)

		mockRepo.AssertExpectations(t)
	})
}
