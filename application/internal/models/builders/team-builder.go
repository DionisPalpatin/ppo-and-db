package builders

import "time"
import "github.com/DionisPalpatin/ppo-and-db/application/internal/models"

type TeamBuilder struct {
	Team *models.Team
}

func NewTeamBuilder() *TeamBuilder {
	return &TeamBuilder{
		Team: &models.Team{
			Id:               1,
			Name:             "Default Team",
			RegistrationDate: time.Now(),
		},
	}
}

func (b *TeamBuilder) WithId(id int) *TeamBuilder {
	b.Team.Id = id
	return b
}

func (b *TeamBuilder) WithName(name string) *TeamBuilder {
	b.Team.Name = name
	return b
}

func (b *TeamBuilder) WithRegistrationDate(date time.Time) *TeamBuilder {
	b.Team.RegistrationDate = date
	return b
}

func (b *TeamBuilder) Build() *models.Team {
	return b.Team
}
