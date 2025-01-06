package bl

import "github.com/DionisPalpatin/ppo-and-db/application/internal/models"

type TeamService struct {
	itr ITeamRepository
}

func (ts *TeamService) GetTeam(id int, name string, searchBy int, requester *models.User) (*models.Team, *MyError) {
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

	return team, err
}

func (ts *TeamService) GetAllTeams(requester *models.User) ([]*models.Team, *MyError) {
	if requester.Role != Admin {
		return nil, CreateError(ErrAccessDenied, "GetAllTeams", "bl")
	}

	return ts.itr.GetAllTeams()
}

func (ts *TeamService) UpdateTeam(requester *models.User, team *models.Team) *MyError {
	if requester.Role != Admin {
		return CreateError(ErrAccessDenied, "UpdateTeam", "bl")
	}

	return ts.itr.UpdateTeam(team)
}

func (ts *TeamService) DeleteTeam(requester *models.User, id int) *MyError {
	if requester.Role != Admin {
		return CreateError(ErrAccessDenied, "DeleteTeam", "bl")
	}

	return ts.itr.DeleteTeam(id)
}

func (ts *TeamService) AddTeam(requester *models.User, team *models.Team) (int, *MyError) {
	if requester.Role != Admin {
		return 0, CreateError(ErrAccessDenied, "AddTeam", "bl")
	}

	return ts.itr.AddTeam(team)
}

func (ts *TeamService) AddUserToTeam(requester *models.User, userID int, teamID int) *MyError {
	if requester.Role != Admin {
		return CreateError(ErrAccessDenied, "AddUserToTeam", "bl")
	}

	return ts.itr.AddUserToTeam(userID, teamID)
}

func (ts *TeamService) DeleteUserFromTeam(requester *models.User, userID int, teamID int) *MyError {
	if requester.Role != Admin {
		return CreateError(ErrAccessDenied, "AddUserToTeam", "bl")
	}

	return ts.itr.DeleteUserFromTeam(userID, teamID)
}

func (ts *TeamService) GetTeamMembers(teamID int, requester *models.User) ([]*models.User, *MyError) {
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
	return ts.itr.GetUserTeam(user)
}

func (ts *TeamService) GetSectionTeam(secID int, requester *models.User) (*models.Team, *MyError) {
	if requester.Role != Admin {
		return nil, CreateError(ErrAccessDenied, "GetSectionTeam", "bl")
	}

	return ts.itr.GetTeamBySectionID(secID)
}
