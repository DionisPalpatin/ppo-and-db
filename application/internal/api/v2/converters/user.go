package converters

import (
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/api/v2/transport_models"
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/models"
	"time"
)

func FromRegistrationInfo(user *transport_models.UserRegistrationInfo) models.User {
	return models.User{
		Fio:      user.FIO,
		Login:    user.Login,
		Password: user.Password,
	}
}

func FromLoginInfo(user *transport_models.UserLoginInfo) models.User {
	return models.User{
		Login:    user.Login,
		Password: user.Password,
	}
}

func ToUserPublicInfo(user *models.User) transport_models.UserPublicInfo {
	scheme := transport_models.UserPublicInfo{
		ID:   user.Id,
		FIO:  user.Fio,
		Role: user.Role,
	}

	return scheme
}

func ToUserPrivateInfo(user *models.User) transport_models.UserPrivateInfo {
	scheme := transport_models.UserPrivateInfo{
		ID:       user.Id,
		FIO:      user.Fio,
		Login:    user.Login,
		Password: user.Password,
		Role:     user.Role,
	}

	return scheme
}

func ToUserFullInfo(user *models.User) transport_models.UserFullInfo {
	scheme := transport_models.UserFullInfo{
		ID:               user.Id,
		FIO:              user.Fio,
		RegistrationDate: user.RegistrationDate.String(),
		Login:            user.Login,
		Password:         user.Password,
		Role:             user.Role,
	}

	return scheme
}

func FromUserPublicInfo(user *transport_models.UserPublicInfo) models.User {
	scheme := models.User{
		Id:   user.ID,
		Fio:  user.FIO,
		Role: user.Role,
	}

	return scheme
}

func FromUserPrivateInfo(user *transport_models.UserPrivateInfo) models.User {
	scheme := models.User{
		Id:       user.ID,
		Fio:      user.FIO,
		Login:    user.Login,
		Password: user.Password,
		Role:     user.Role,
	}

	return scheme
}

func FromUserFullInfo(user *transport_models.UserFullInfo) (models.User, error) {
	t, err := time.Parse(user.RegistrationDate, "2006-01-02 15:04:05-07")
	if err != nil {
		return models.User{}, err
	}

	scheme := models.User{
		Id:               user.ID,
		Fio:              user.FIO,
		RegistrationDate: t,
		Login:            user.Login,
		Password:         user.Password,
		Role:             user.Role,
	}

	return scheme, nil
}
