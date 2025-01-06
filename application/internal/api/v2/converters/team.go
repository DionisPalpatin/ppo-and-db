package converters

import (
	"time"

	"github.com/DionisPalpatin/ppo-and-db/application/internal/api/v2/transport_models"
	"github.com/DionisPalpatin/ppo-and-db/application/internal/models"
)

func ToTeamInfo(team *models.Team) transport_models.TeamInfo {
	return transport_models.TeamInfo{
		ID:   team.Id,
		Name: team.Name,
	}
}

func ToTeamFullInfo(team *models.Team) transport_models.TeamFullInfo {
	return transport_models.TeamFullInfo{
		ID:               team.Id,
		Name:             team.Name,
		RegistrationDate: team.RegistrationDate.String(),
	}
}

func FromTeamInfo(teamInfo *transport_models.TeamInfo) models.Team {
	return models.Team{
		Id:   teamInfo.ID,
		Name: teamInfo.Name,
	}
}

func FromTeamFullInfo(teamFullInfo *transport_models.TeamFullInfo) (models.Team, error) {
	parsedTime, err := time.Parse("2006-01-02 15:04:05-07", teamFullInfo.RegistrationDate)
	if err != nil {
		return models.Team{}, err
	}

	return models.Team{
		Id:               teamFullInfo.ID,
		Name:             teamFullInfo.Name,
		RegistrationDate: parsedTime,
	}, nil
}
