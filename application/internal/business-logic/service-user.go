package bl

import (
	"log/slog"

	mylogger "github.com/DionisPalpatin/ppo-and-db/application/internal/logger"
	"github.com/DionisPalpatin/ppo-and-db/application/internal/models"
)

type UserService struct {
	iur IUserRepository
	logger *mylogger.MyLogger
}

func NewUserService(iur IUserRepository, logger *mylogger.MyLogger) IUserService {
	return &UserService{iur: iur, logger: logger}
}

func (us *UserService) GetUser(id int, login string, searchBy int, requester *models.User) (*models.User, *MyError) {
	us.logger.WriteLog("GetUser is called (Repo)", slog.LevelInfo, nil)

	if requester.Role != Admin && requester.Id != id {
		err := CreateError(ErrAccessDenied, "GetUser", "")
		return nil, err
	}

	var user *models.User
	var err *MyError

	switch searchBy {
	case SearchByID:
		user, err = us.iur.GetUserByID(id)

	case SearchByString:
		user, err = us.iur.GetUserByLogin(login)

	default:
		user = nil
		err = CreateError(ErrSearchParameter, "GetUser", "")
	}

	resState := CreateError(Ok, "GetUser", "business_logic")
	us.logger.WriteLog(resState.ConcatenateFields(), slog.LevelInfo, nil)
	return user, err
}

func (us *UserService) GetAllUsers(requester *models.User) ([]*models.User, *MyError) {
	us.logger.WriteLog("GetAllUsers is called (Repo)", slog.LevelInfo, nil)

	if requester.Role != Admin {
		err := CreateError(ErrAccessDenied, "GetAllUsers", "")
		return nil, err
	}

	resState := CreateError(Ok, "GetAllUsers", "business_logic")
	us.logger.WriteLog(resState.ConcatenateFields(), slog.LevelInfo, nil)
	return us.iur.GetAllUsers()
}

func (us *UserService) UpdateUser(requester *models.User, user *models.User) *MyError {
	us.logger.WriteLog("UpdateUser is called (Repo)", slog.LevelInfo, nil)
	
	if requester.Role != Admin {
		err := CreateError(ErrAccessDenied, "UpdateUser", "")
		return err
	}

	resState := CreateError(Ok, "UpdateUser", "business_logic")
	us.logger.WriteLog(resState.ConcatenateFields(), slog.LevelInfo, nil)
	return us.iur.UpdateUser(user)
}

func (us *UserService) DeleteUser(requester *models.User, id int) *MyError {
	us.logger.WriteLog("DeleteUser is called (Repo)", slog.LevelInfo, nil)
	
	if requester.Role != Admin {
		err := CreateError(ErrAccessDenied, "DeleteUser", "")
		return err
	}

	resState := CreateError(Ok, "DeleteUser", "business_logic")
	us.logger.WriteLog(resState.ConcatenateFields(), slog.LevelInfo, nil)
	return us.iur.DeleteUser(id)
}
