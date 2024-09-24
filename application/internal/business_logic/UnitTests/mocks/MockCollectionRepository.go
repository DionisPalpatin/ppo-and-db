package mocks

import (
	"github.com/stretchr/testify/mock"

	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/business_logic"
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/models"
)

type MockICollectionRepository struct {
	mock.Mock
}

func (m *MockICollectionRepository) GetCollectionByID(id int) (*models.Collection, *bl.MyError) {
	args := m.Called(id)
	return args.Get(0).(*models.Collection), args.Get(1).(*bl.MyError)
}

func (m *MockICollectionRepository) GetCollectionByName(name string) (*models.Collection, *bl.MyError) {
	args := m.Called(name)
	return args.Get(0).(*models.Collection), args.Get(1).(*bl.MyError)
}

func (m *MockICollectionRepository) GetAllCollections() ([]*models.Collection, *bl.MyError) {
	args := m.Called()
	return args.Get(0).([]*models.Collection), args.Get(1).(*bl.MyError)
}

func (m *MockICollectionRepository) GetAllUserCollections(user *models.User) ([]*models.Collection, *bl.MyError) {
	args := m.Called(user)
	return args.Get(0).([]*models.Collection), args.Get(1).(*bl.MyError)
}

func (m *MockICollectionRepository) AddCollection(collection *models.Collection) *bl.MyError {
	args := m.Called(collection)
	return args.Get(0).(*bl.MyError)
}

func (m *MockICollectionRepository) DeleteCollection(id int) *bl.MyError {
	args := m.Called(id)
	return args.Get(0).(*bl.MyError)
}

func (m *MockICollectionRepository) UpdateCollection(collection *models.Collection) *bl.MyError {
	args := m.Called(collection)
	return args.Get(0).(*bl.MyError)
}

func (m *MockICollectionRepository) GetAllNotesInCollection(collection *models.Collection) ([]*models.Note, *bl.MyError) {
	args := m.Called(collection)
	return args.Get(0).([]*models.Note), args.Get(1).(*bl.MyError)
}
