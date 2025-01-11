package bl

import (
	"time"

	"github.com/DionisPalpatin/ppo-and-db/application/internal/models"
	mylogger "github.com/DionisPalpatin/ppo-and-db/application/internal/logger"
)

var TEST bool = true
var ZEROTIME time.Time

type OAuthService struct {
	iur IUserRepository
	logger *mylogger.MyLogger
}

func NewOAuthService(iur IUserRepository, logger *mylogger.MyLogger) IOAuthService {
	return &OAuthService{iur: iur, logger: logger}
}

func (oas *OAuthService) RegisterUser(fio string, login string, password string) (*models.User, *MyError) {
	_, err := oas.iur.GetUserByLogin(login)
	if err.ErrNum == Ok {
		return nil, CreateError(UserExists, "RegisterUse", "bl")
	}

	user := &models.User{
		Fio:              fio,
		Login:            login,
		Password:         password,
		Role:             Reader,
		RegistrationDate: time.Now(),
	}

	if TEST {
		user.RegistrationDate = ZEROTIME
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
