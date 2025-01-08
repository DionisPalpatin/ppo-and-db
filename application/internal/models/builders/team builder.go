package builders

import "time"
import "github.com/DionisPalpatin/ppo-and-db/application/internal/models"

type TeamBuilder struct {
	team *models.Team
}

func NewTeamBuilder() *TeamBuilder {
	return &TeamBuilder{
		team: &models.Team{
			Id:               1,
			Name:             "Default Team",
			RegistrationDate: time.Now(),
		},
	}
}

func (b *TeamBuilder) WithId(id int) *TeamBuilder {
	b.team.Id = id
	return b
}

func (b *TeamBuilder) WithName(name string) *TeamBuilder {
	b.team.Name = name
	return b
}

func (b *TeamBuilder) WithRegistrationDate(date time.Time) *TeamBuilder {
	b.team.RegistrationDate = date
	return b
}

func (b *TeamBuilder) Build() *models.Team {
	return b.team
}
