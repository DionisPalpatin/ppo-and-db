package mocks

import (
	bl "github.com/DionisPalpatin/ppo-and-db/application/internal/business-logic"
	"github.com/DionisPalpatin/ppo-and-db/application/internal/models"
)

//go:generate mockery --name MockSectionRepository --with-expecter
type MockSectionRepository interface {
	GetSectionByID(id int) (*models.Section, *bl.MyError)
	GetSectionByTeamName(teamName string) (*models.Section, *bl.MyError)
	GetAllSections() ([]*models.Section, *bl.MyError)
	AddSection(section *models.Section, team *models.Team) (int, *bl.MyError)
	DeleteSection(id int) *bl.MyError
	UpdateSection(section *models.Section) *bl.MyError
	GetAllNotesInSection(section *models.Section) ([]*models.Note, *bl.MyError)
	AddNoteToSection(note *models.Note, section *models.Section) *bl.MyError
	DeleteNoteFromSection(note *models.Note, section *models.Section) *bl.MyError
}
