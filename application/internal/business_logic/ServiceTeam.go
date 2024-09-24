package bl

import "github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/models"

type TeamService struct{}

func (ts TeamService) GetTeam(id int, name string, searchBy int, requester *models.User, itr ITeamRepository) (*models.Team, *MyError) {
	// if requester.Role != Admin {
	// 	return nil, CreateError(ErrAccessDenied, ErrAccessDeniedError(), "GetTeam")
	// }

	var team *models.Team
	var err *MyError

	switch searchBy {
	case SearchByID:
		team, err = itr.GetTeamByID(id)

	case SearchByString:
		team, err = itr.GetTeamByName(name)

	default:
		team = nil
		err = CreateError(ErrSearchParameter, ErrSearchParameterError(), "GetTeam")
	}

	return team, err
}

func (ts TeamService) GetAllTeams(requester *models.User, itr ITeamRepository) ([]*models.Team, *MyError) {
	if requester.Role != Admin {
		return nil, CreateError(ErrAccessDenied, ErrAccessDeniedError(), "GetAllTeams")
	}

	return itr.GetAllTeams()
}

func (ts TeamService) UpdateTeam(requester *models.User, team *models.Team, itr ITeamRepository) *MyError {
	if requester.Role != Admin {
		return CreateError(ErrAccessDenied, ErrAccessDeniedError(), "UpdateTeam")
	}

	return itr.UpdateTeam(team)
}

func (ts TeamService) DeleteTeam(requester *models.User, id int, itr ITeamRepository) *MyError {
	if requester.Role != Admin {
		return CreateError(ErrAccessDenied, ErrAccessDeniedError(), "DeleteTeam")
	}

	return itr.DeleteTeam(id)
}

func (ts TeamService) AddTeam(requester *models.User, team *models.Team, itr ITeamRepository) *MyError {
	if requester.Role != Admin {
		return CreateError(ErrAccessDenied, ErrAccessDeniedError(), "AddTeam")
	}

	return itr.AddTeam(team)
}

func (ts TeamService) AddUserToTeam(requester *models.User, userID int, teamID int, itr ITeamRepository) *MyError {
	if requester.Role != Admin {
		return CreateError(ErrAccessDenied, ErrAccessDeniedError(), "AddUserToTeam")
	}

	return itr.AddUserToTeam(userID, teamID)
}

func (ts TeamService) DeleteUserFromTeam(requester *models.User, userID int, teamID int, itr ITeamRepository) *MyError {
	if requester.Role != Admin {
		return CreateError(ErrAccessDenied, ErrAccessDeniedError(), "AddUserToTeam")
	}

	return itr.DeleteUserFromTeam(userID, teamID)
}

func (ts TeamService) GetTeamMembers(teamID int, requester *models.User, itr ITeamRepository) ([]*models.User, *MyError) {
	if requester.Role != Admin {
		return nil, CreateError(ErrAccessDenied, ErrAccessDeniedError(), "GetTeamMembers")
	}

	_, err := itr.GetTeamByID(teamID)
	if err.ErrNum != AllIsOk {
		return nil, err
	}

	return itr.GetTeamMembers(teamID)
}

func (ts TeamService) GetUserTeam(user *models.User, itr ITeamRepository) (*models.Team, *MyError) {
	return itr.GetUserTeam(user)
}
