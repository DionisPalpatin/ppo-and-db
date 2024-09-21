package display

import (
	"fmt"

	"notebook_app/internal/models"
)

func DisplayTeam(team *models.Team) {
	if team != nil {
		fmt.Println("## Команда:")
		fmt.Printf("ID: %d\n", team.Id)
		fmt.Printf("Дата регистрации: %s\n", team.RegistrationDate)
	}
}

func DisplayAllTeams(teams []*models.Team) {
	for i, team := range teams {
		fmt.Printf("%d) %s\n", i, team.Name)
	}
}

func DisplayTeamMembers(members []*models.User) {
	for i, us := range members {
		fmt.Printf("%d) %s\n", i, us.Fio)
	}
}
