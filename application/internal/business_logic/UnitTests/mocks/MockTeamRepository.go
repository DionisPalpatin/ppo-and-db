package mocks

import (
	"github.com/stretchr/testify/mock"

	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/business_logic"
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/models"
)

// MockITeamRepository Мок-реализация ITeamRepository
type MockITeamRepository struct {
	mock.Mock
}

func (m *MockITeamRepository) GetTeamByID(id int) (*models.Team, *bl.MyError) {
	args := m.Called(id)
	return args.Get(0).(*models.Team), args.Get(1).(*bl.MyError)
}

func (m *MockITeamRepository) GetTeamByName(name string) (*models.Team, *bl.MyError) {
	args := m.Called(name)
	return args.Get(0).(*models.Team), args.Get(1).(*bl.MyError)
}

func (m *MockITeamRepository) GetAllTeams() ([]*models.Team, *bl.MyError) {
	args := m.Called()
	return args.Get(0).([]*models.Team), args.Get(1).(*bl.MyError)
}

func (m *MockITeamRepository) AddTeam(team *models.Team) *bl.MyError {
	args := m.Called(team)
	return args.Get(0).(*bl.MyError)
}

func (m *MockITeamRepository) DeleteTeam(id int) *bl.MyError {
	args := m.Called(id)
	return args.Get(0).(*bl.MyError)
}

func (m *MockITeamRepository) AddUser(userID int, teamID int) *bl.MyError {
	args := m.Called(userID, teamID)
	return args.Get(0).(*bl.MyError)
}

func (m *MockITeamRepository) UpdateTeam(team *models.Team) *bl.MyError {
	args := m.Called(team)
	return args.Get(0).(*bl.MyError)
}

func (m *MockITeamRepository) GetTeamMembers(teamID int) ([]*models.User, *bl.MyError) {
	args := m.Called(teamID)
	return args.Get(0).([]*models.User), args.Get(1).(*bl.MyError)
}
