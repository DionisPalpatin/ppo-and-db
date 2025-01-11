package mocks

import (
	bl "github.com/DionisPalpatin/ppo-and-db/application/internal/business-logic"
	"github.com/DionisPalpatin/ppo-and-db/application/internal/models"
)

//go:generate mockery --name MockNoteRepository --with-expecter
type MockNoteRepository interface {
	GetNoteByID(id int) (*models.Note, *bl.MyError)
	GetNoteByName(name string) (*models.Note, *bl.MyError)
	GetAllNotes() ([]*models.Note, *bl.MyError)
	GetAllPublicNotes() ([]*models.Note, *bl.MyError)
	AddNote(note *models.Note) (int, *bl.MyError)
	DeleteNote(id int) *bl.MyError
	UpdateNoteContent(note *models.Note) *bl.MyError
	UpdateNoteInfo(note *models.Note) *bl.MyError
	AddNoteToCollection(collectionID int, noteID int) *bl.MyError
	DeleteNoteFromCollection(collectionID int, noteID int) *bl.MyError
}
