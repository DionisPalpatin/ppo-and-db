package mocks

import (
	"github.com/stretchr/testify/mock"

	"notebook_app/internal/business_logic"
	"notebook_app/internal/models"
)

// MockIUserRepository Мок-реализация IUserRepository
type MockIUserRepository struct {
	mock.Mock
}

func (m *MockIUserRepository) GetUserByID(id int) (*models.User, *bl.MyError) {
	args := m.Called(id)
	return args.Get(0).(*models.User), args.Get(1).(*bl.MyError)
}

func (m *MockIUserRepository) GetUserByLogin(login string) (*models.User, *bl.MyError) {
	args := m.Called(login)
	return args.Get(0).(*models.User), args.Get(1).(*bl.MyError)
}

func (m *MockIUserRepository) GetAllUsers() ([]*models.User, *bl.MyError) {
	args := m.Called()
	return args.Get(0).([]*models.User), args.Get(1).(*bl.MyError)
}

func (m *MockIUserRepository) AddUser(user *models.User) *bl.MyError {
	args := m.Called(user)
	return args.Get(0).(*bl.MyError)
}

func (m *MockIUserRepository) DeleteUser(id int) *bl.MyError {
	args := m.Called(id)
	return args.Get(0).(*bl.MyError)
}

func (m *MockIUserRepository) UpdateUser(user *models.User) *bl.MyError {
	args := m.Called(user)
	return args.Get(0).(*bl.MyError)
}
