package bl

import (
	"time"

	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/models"
)

type OAuthService struct {
	iur IUserRepository
}

func (oas *OAuthService) RegisterUser(fio string, login string, password string) (*models.User, *MyError) {
	user, err := oas.iur.GetUserByLogin(login)
	if err.ErrNum == Ok {
		return nil, CreateError(UserExists, "RegisterUse", "bl")
	}

	user = &models.User{
		Fio:              fio,
		Login:            login,
		Password:         password,
		Role:             Reader,
		RegistrationDate: time.Now(),
	}

	return user, oas.iur.AddUser(user)
}

func (oas *OAuthService) SignInUser(login string, password string) (*models.Token, *MyError) {
	user, myErr := oas.iur.GetUserByLogin(login)

	if myErr.ErrNum != Ok {
		return nil, myErr
	}
	if user.Password != password {
		myErr := CreateError(AuthenticationError, "SignInUser", "bl")
		return nil, myErr
	}

	token, _ := generateToken(user.Id, string(rune(user.Role)))

	myOk := CreateError(Ok, "", "bl")
	return token, myOk
}
