package mocks

import (
	bl "github.com/DionisPalpatin/ppo-and-db/application/internal/business-logic"
	"github.com/DionisPalpatin/ppo-and-db/application/internal/models"
)

//go:generate mockery --name MockUserRepository --with-expecter
type MockUserRepository interface {
	GetUserByID(id int) (*models.User, *bl.MyError)
	GetUserByLogin(login string) (*models.User, *bl.MyError)
	GetAllUsers() ([]*models.User, *bl.MyError)
	AddUser(user *models.User) *bl.MyError
	DeleteUser(id int) *bl.MyError
	UpdateUser(user *models.User) *bl.MyError
}
