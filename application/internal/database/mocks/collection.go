package mocks

import (
	bl "github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/business_logic"
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/models"
)

//go:generate mockery --name MockCollectionRepository --with-expecter
type MockCollectionRepository interface {
	GetCollectionByID(id int) (*models.Collection, *bl.MyError)
	GetCollectionByName(name string) (*models.Collection, *bl.MyError)
	GetAllCollections() ([]*models.Collection, *bl.MyError)
	GetAllUserCollections(user *models.User) ([]*models.Collection, *bl.MyError)
	AddCollection(collection *models.Collection) (int, *bl.MyError)
	DeleteCollection(id int) *bl.MyError
	UpdateCollection(collection *models.Collection) *bl.MyError
	GetAllNotesInCollection(collection *models.Collection) ([]*models.Note, *bl.MyError)
}
