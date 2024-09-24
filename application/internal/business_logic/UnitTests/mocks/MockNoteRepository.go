package mocks

import (
	"io"

	"github.com/stretchr/testify/mock"

	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/business_logic"
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/models"
)

// MockINoteRepository Мок-реализация INoteRepository
type MockINoteRepository struct {
	mock.Mock
}

func (m *MockINoteRepository) GetNoteByID(id int) (*models.Note, *bl.MyError) {
	args := m.Called(id)
	return args.Get(0).(*models.Note), args.Get(1).(*bl.MyError)
}

func (m *MockINoteRepository) GetAllNotes() ([]*models.Note, *bl.MyError) {
	args := m.Called()
	return args.Get(0).([]*models.Note), args.Get(1).(*bl.MyError)
}

func (m *MockINoteRepository) AddNote(note *models.Note) *bl.MyError {
	args := m.Called(note)
	return args.Get(0).(*bl.MyError)
}

func (m *MockINoteRepository) DeleteNote(id int) *bl.MyError {
	args := m.Called(id)
	return args.Get(0).(*bl.MyError)
}

func (m *MockINoteRepository) UpdateNoteContentText(reader io.Reader, note *models.Note) *bl.MyError {
	args := m.Called(reader, note)
	return args.Get(0).(*bl.MyError)
}

func (m *MockINoteRepository) UpdateNoteContentImg(reader io.Reader, note *models.Note) *bl.MyError {
	args := m.Called(reader, note)
	return args.Get(0).(*bl.MyError)
}

func (m *MockINoteRepository) UpdateNoteContentRawData(reader io.Reader, note *models.Note) *bl.MyError {
	args := m.Called(reader, note)
	return args.Get(0).(*bl.MyError)
}

func (m *MockINoteRepository) UpdateNoteInfo(note *models.Note) *bl.MyError {
	args := m.Called(note)
	return args.Get(0).(*bl.MyError)
}

func (m *MockINoteRepository) AddNoteToCollection(collectionID int, noteID int) *bl.MyError {
	args := m.Called(collectionID, noteID)
	return args.Get(0).(*bl.MyError)
}
