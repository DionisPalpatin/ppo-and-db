package builders

import (
	"time"

	bl "github.com/DionisPalpatin/ppo-and-db/application/internal/business-logic"
	"github.com/DionisPalpatin/ppo-and-db/application/internal/models"
)

type UserBuilder struct {
	User *models.User
}

func NewUserBuilder() *UserBuilder {
	return &UserBuilder{
		User: &models.User{
			Id:               7,
			Fio:              "uzi",
			Login:            "uzi@gmail.com",
			Password:         "password",
			RegistrationDate: time.Now(),
			Role:             bl.Reader,
		},
	}
}

func (b *UserBuilder) WithUserID(id int) *UserBuilder {
	b.User.Id = id
	return b
}

func (b *UserBuilder) WithFio(fio string) *UserBuilder {
	b.User.Fio = fio
	return b
}

func (b *UserBuilder) WithLogin(login string) *UserBuilder {
	b.User.Login = login
	return b
}

func (b *UserBuilder) WithRegistrationDate(date time.Time) *UserBuilder {
	b.User.RegistrationDate = date
	return b
}

func (b *UserBuilder) WithPassword(password string) *UserBuilder {
	b.User.Password = password
	return b
}

func (b *UserBuilder) WithRole(role int) *UserBuilder {
	b.User.Role = role
	return b
}

func (b *UserBuilder) Build() *models.User {
	return b.User
}
