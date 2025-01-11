package bl

import (
	"log/slog"
	"github.com/DionisPalpatin/ppo-and-db/application/internal/models"
	mylogger "github.com/DionisPalpatin/ppo-and-db/application/internal/logger"
)

type TeamService struct {
	itr ITeamRepository
	logger *mylogger.MyLogger
}

func NewTeamService(itr ITeamRepository, logger *mylogger.MyLogger) ITeamService {
	return &TeamService{itr: itr, logger: logger}
}

func (ts *TeamService) GetTeam(id int, name string, searchBy int, requester *models.User) (*models.Team, *MyError) {
	ts.logger.WriteLog("GetTeam is called (Service)", slog.LevelInfo, nil)

	var team *models.Team
	var err *MyError

	if requester.Role != Admin {
		return nil, CreateError(ErrAccessDenied, "GetTeam", "bl")
	}

	switch searchBy {
	case SearchByID:
		team, err = ts.itr.GetTeamByID(id)

	case SearchByString:
		team, err = ts.itr.GetTeamByName(name)

	default:
		team = nil
		err = CreateError(ErrSearchParameter, "GetTeam", "bl")
	}

	resState := CreateError(Ok, "GetTeam", "business_logic")
	ts.logger.WriteLog(resState.ConcatenateFields(), slog.LevelInfo, nil)
	return team, err
}

func (ts *TeamService) GetAllTeams(requester *models.User) ([]*models.Team, *MyError) {
	ts.logger.WriteLog("GetAllTeams is called (Service)", slog.LevelInfo, nil)

	if requester.Role != Admin {
		return nil, CreateError(ErrAccessDenied, "GetAllTeams", "bl")
	}

	return ts.itr.GetAllTeams()
}

func (ts *TeamService) UpdateTeam(requester *models.User, team *models.Team) *MyError {
	ts.logger.WriteLog("UpdateTeam is called (Service)", slog.LevelInfo, nil)

	if requester.Role != Admin {
		return CreateError(ErrAccessDenied, "UpdateTeam", "bl")
	}

	return ts.itr.UpdateTeam(team)
}

func (ts *TeamService) DeleteTeam(requester *models.User, id int) *MyError {
	ts.logger.WriteLog("DeleteTeam is called (Service)", slog.LevelInfo, nil)

	if requester.Role != Admin {
		return CreateError(ErrAccessDenied, "DeleteTeam", "bl")
	}

	return ts.itr.DeleteTeam(id)
}

func (ts *TeamService) AddTeam(requester *models.User, team *models.Team) (int, *MyError) {
	ts.logger.WriteLog("AddTeam is called (Service)", slog.LevelInfo, nil)

	if requester.Role != Admin {
		return 0, CreateError(ErrAccessDenied, "AddTeam", "bl")
	}

	return ts.itr.AddTeam(team)
}

func (ts *TeamService) AddUserToTeam(requester *models.User, userID int, teamID int) *MyError {
	ts.logger.WriteLog("AddUserToTeam is called (Service)", slog.LevelInfo, nil)

	if requester.Role != Admin {
		return CreateError(ErrAccessDenied, "AddUserToTeam", "bl")
	}

	resState := CreateError(Ok, "AddUserToTeam", "business_logic")
	ts.logger.WriteLog(resState.ConcatenateFields(), slog.LevelInfo, nil)
	return ts.itr.AddUserToTeam(userID, teamID)
}

func (ts *TeamService) DeleteUserFromTeam(requester *models.User, userID int, teamID int) *MyError {
	ts.logger.WriteLog("DeleteUserFromTeam is called (Service)", slog.LevelInfo, nil)

	if requester.Role != Admin {
		return CreateError(ErrAccessDenied, "AddUserToTeam", "bl")
	}

	resState := CreateError(Ok, "DeleteUserFromTeam", "business_logic")
	ts.logger.WriteLog(resState.ConcatenateFields(), slog.LevelInfo, nil)
	return ts.itr.DeleteUserFromTeam(userID, teamID)
}

func (ts *TeamService) GetTeamMembers(teamID int, requester *models.User) ([]*models.User, *MyError) {
	ts.logger.WriteLog("GetTeamMembers is called (service)", slog.LevelInfo, nil)

	if requester.Role != Admin {
		return nil, CreateError(ErrAccessDenied, "GetTeamMembers", "bl")
	}

	_, err := ts.itr.GetTeamByID(teamID)
	if err.ErrNum != Ok {
		return nil, err
	}

	return ts.itr.GetTeamMembers(teamID)
}

func (ts *TeamService) GetUserTeam(user *models.User) (*models.Team, *MyError) {
	ts.logger.WriteLog("GetUserTeam is called (Service)", slog.LevelInfo, nil)
	return ts.itr.GetUserTeam(user)
}

// func (ts *TeamService) GetSectionTeam(secID int, requester *models.User) (*models.Team, *MyError) {
// 	ts.logger.WriteLog("GetSectionTeam is called (Service)", slog.LevelInfo, nil)

// 	if requester.Role != Admin {
// 		return nil, CreateError(ErrAccessDenied, "GetSectionTeam", "bl")
// 	}

// 	return ts.itr.GetTeamBySectionID(secID)
// }
