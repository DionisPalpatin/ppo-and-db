package bl

import "github.com/DionisPalpatin/ppo-and-db/application/internal/models"

type UserService struct {
	iur IUserRepository
}

func (us *UserService) GetUser(id int, login string, searchBy int, requester *models.User) (*models.User, *MyError) {
	if requester.Role != Admin && requester.Id != id {
		err := CreateError(ErrAccessDenied, "GetUser", "")
		return nil, err
	}

	var user *models.User
	var err *MyError

	// Получаем пользователя
	switch searchBy {
	case SearchByID:
		user, err = us.iur.GetUserByID(id)

	case SearchByString:
		user, err = us.iur.GetUserByLogin(login)

	default:
		user = nil
		err = CreateError(ErrSearchParameter, "GetUser", "")
	}

	return user, err
}

func (us *UserService) GetAllUsers(requester *models.User) ([]*models.User, *MyError) {
	if requester.Role != Admin {
		err := CreateError(ErrAccessDenied, "GetAllUsers", "")
		return nil, err
	}

	return us.iur.GetAllUsers()
}

func (us *UserService) UpdateUser(requester *models.User, user *models.User) *MyError {
	if requester.Role != Admin {
		err := CreateError(ErrAccessDenied, "UpdateUser", "")
		return err
	}

	return us.iur.UpdateUser(user)
}

func (us *UserService) DeleteUser(requester *models.User, id int) *MyError {
	if requester.Role != Admin {
		err := CreateError(ErrAccessDenied, "DeleteUser", "")
		return err
	}

	return us.iur.DeleteUser(id)
}
