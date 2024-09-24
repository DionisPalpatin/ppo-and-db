package bl

import "github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/models"

type UserService struct{}

func (usSrv UserService) GetUser(id int, login string, searchBy int, requester *models.User, iur IUserRepository) (*models.User, *MyError) {
	// Проверка, что пользователь имеет право запрашивать других пользователей
	if requester.Role != Admin {
		err := CreateError(ErrAccessDenied, ErrAccessDeniedError(), "GetUser")
		return nil, err
	}

	var user *models.User
	var err *MyError

	// Получаем пользователя
	switch searchBy {
	case SearchByID:
		user, err = iur.GetUserByID(id)

	case SearchByString:
		user, err = iur.GetUserByLogin(login)

	default:
		user = nil
		err = CreateError(ErrSearchParameter, ErrSearchParameterError(), "GetUser")
	}

	return user, err
}

func (usSrv UserService) GetAllUsers(requester *models.User, iur IUserRepository) ([]*models.User, *MyError) {
	// Проверка, что пользователь имеет право запрашивать других пользователей
	if requester.Role != Admin {
		err := CreateError(ErrAccessDenied, ErrAccessDeniedError(), "GetAllUsers")
		return nil, err
	}

	return iur.GetAllUsers()
}

func (usSrv UserService) UpdateUser(requester *models.User, user *models.User, iur IUserRepository) *MyError {
	// Проверка, что пользователь имеет право изменять данные других пользователей
	if requester.Role != Admin {
		err := CreateError(ErrAccessDenied, ErrAccessDeniedError(), "UpdateUser")
		return err
	}

	return iur.UpdateUser(user)
}

func (usSrv UserService) DeleteUser(requester *models.User, id int, iur IUserRepository) *MyError {
	// Проверка, что пользователь имеет право удалять других пользователей
	if requester.Role != Admin {
		err := CreateError(ErrAccessDenied, ErrAccessDeniedError(), "DeleteUser")
		return err
	}

	return iur.DeleteUser(id)
}
