package bl

import (
	"time"

	"notebook_app/config"
	"notebook_app/internal/models"
)

type OAuthService struct{}

func (OAuthService) RegisterUser(fio string, login string, password string, iur IUserRepository) (*models.User, *MyError) {
	user, err := iur.GetUserByLogin(login)
	if err.ErrNum == AllIsOk {
		return nil, CreateError(ErrRegisterUser, ErrRegistrationError(), "RegisterUse")
	}

	user = &models.User{
		Fio:              fio,
		Login:            login,
		Password:         password,
		Role:             Reader,
		RegistrationDate: time.Now().Format(config.Configs{}.DateTimeFormat),
	}
	return user, iur.AddUser(user)
}

func (OAuthService) SignInUser(login string, password string, iur IUserRepository) (*models.User, *MyError) {
	user, err := iur.GetUserByLogin(login)

	if err.ErrNum != AllIsOk {
		return nil, err
	}
	if user.Password != password {
		return nil, CreateError(ErrSignInUser, ErrSignInUserError(), "SignInUser")
	}

	return user, &MyError{ErrNum: AllIsOk, FuncName: "", Err: nil}
}
