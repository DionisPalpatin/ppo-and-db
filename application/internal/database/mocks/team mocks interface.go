package mocks

import (
	bl "github.com/DionisPalpatin/ppo-and-db/application/internal/business_logic"
	"github.com/DionisPalpatin/ppo-and-db/application/internal/models"
)

//go:generate mockery --name MockTeamRepository --with-expecter
type MockTeamRepository interface {
	GetTeamByID(id int) (*models.Team, *bl.MyError)
	GetTeamByName(name string) (*models.Team, *bl.MyError)
	GetTeamBySectionID(id int) (*models.Team, *bl.MyError)
	GetAllTeams() ([]*models.Team, *bl.MyError)
	AddTeam(team *models.Team) (int, *bl.MyError)
	DeleteTeam(id int) *bl.MyError
	AddUserToTeam(userID int, teamID int) *bl.MyError
	DeleteUserFromTeam(uid int, tid int) *bl.MyError
	UpdateTeam(team *models.Team) *bl.MyError
	GetTeamMembers(teamID int) ([]*models.User, *bl.MyError)
	GetUserTeam(user *models.User) (*models.Team, *bl.MyError)
}
