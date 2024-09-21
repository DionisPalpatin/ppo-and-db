package mocks

import (
	"github.com/stretchr/testify/mock"

	"notebook_app/internal/business_logic"
	"notebook_app/internal/models"
)

// MockISectionRepository Мок-реализация ISectionRepository
type MockISectionRepository struct {
	mock.Mock
}

func (m *MockISectionRepository) GetSectionByID(id int) (*models.Section, *bl.MyError) {
	args := m.Called(id)
	return args.Get(0).(*models.Section), args.Get(1).(*bl.MyError)
}

func (m *MockISectionRepository) GetSectionByTeamName(teamName string) (*models.Section, *bl.MyError) {
	args := m.Called(teamName)
	return args.Get(0).(*models.Section), args.Get(1).(*bl.MyError)
}

func (m *MockISectionRepository) GetAllSections() ([]*models.Section, *bl.MyError) {
	args := m.Called()
	return args.Get(0).([]*models.Section), args.Get(1).(*bl.MyError)
}

func (m *MockISectionRepository) AddSection(section *models.Section, team *models.Team) *bl.MyError {
	args := m.Called(section)
	return args.Get(0).(*bl.MyError)
}

func (m *MockISectionRepository) DeleteSection(id int) *bl.MyError {
	args := m.Called(id)
	return args.Get(0).(*bl.MyError)
}

func (m *MockISectionRepository) UpdateSection(section *models.Section) *bl.MyError {
	args := m.Called(section)
	return args.Get(0).(*bl.MyError)
}

func (m *MockISectionRepository) GetAllNotesInSection(section *models.Section) ([]*models.Note, *bl.MyError) {
	args := m.Called(section)
	return args.Get(0).([]*models.Note), args.Get(1).(*bl.MyError)
}

func (m *MockISectionRepository) AddNoteToSection(note *models.Note, section *models.Section) *bl.MyError {
	args := m.Called(note, section)
	return args.Get(0).(*bl.MyError)
}
