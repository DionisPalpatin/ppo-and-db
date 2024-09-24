package methods

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/UI/TechUI/display"
	bl "github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/business_logic"
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/models"
)

func (MenuPoints) AddTeam(user *models.User, dateTimeFormat string, ireps *bl.IRepositories, isvcs *bl.IServices) {
	teamName := ""
	fmt.Print("Введите название команды: ")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		teamName = scanner.Text()
	}

	itr := ireps.ITeamRepo
	its := isvcs.ITeamSvc

	team := &models.Team{Name: teamName, RegistrationDate: time.Now().Format(dateTimeFormat)}

	myErr := its.AddTeam(user, team, itr)
	if myErr.ErrNum == bl.ErrAccessDenied {
		display.DisplayError("Ошибка: у Вас нет прав на выполнение данной операции")
	} else if myErr.ErrNum == bl.ErrDeleteTeam {
		display.DisplayError("Ошибка: попробуйте еще раз")
	} else {
		display.DisplaySuccess("Команда успешно добавлена")
	}
}

func (MenuPoints) DeleteTeam(user *models.User, ireps *bl.IRepositories, isvcs *bl.IServices) {
	id := 0
	fmt.Print("Введите id команды: ")
	fmt.Scanln(&id)

	itr := ireps.ITeamRepo
	its := isvcs.ITeamSvc

	myErr := its.DeleteTeam(user, id, itr)
	if myErr.ErrNum == bl.ErrAccessDenied {
		display.DisplayError("Ошибка: у Вас нет прав на выполнение данной операции")
	} else if myErr.ErrNum == bl.ErrDeleteTeam {
		display.DisplayError("Ошибка: попробуйте еще раз")
	} else {
		display.DisplaySuccess("Команда успешно удалена")
	}
}

func (MenuPoints) FindTeam(user *models.User, ireps *bl.IRepositories, isvcs *bl.IServices) {
	itr := ireps.ITeamRepo
	its := isvcs.ITeamSvc
	choice := 0

	for choice < 1 || choice > 3 {
		fmt.Println("\n1) Искать по названию")
		fmt.Println("2) Искать по ID")
		fmt.Println("3) Назад")
		fmt.Print("Введите номер пункта меню: ")
		fmt.Scanln(&choice)
	}

	if choice == 1 {
		name := ""
		fmt.Print("Введите название команды: ")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			name = scanner.Text()
		}

		team, myErr := its.GetTeam(0, name, bl.SearchByString, user, itr)
		if myErr.ErrNum == bl.ErrAccessDenied {
			display.DisplayError("Ошибка: у Вас нет прав на выполнение данной операции")
		} else if myErr.ErrNum == bl.ErrGetNoteByName {
			display.DisplayError("Ошибка: попробуйте еще раз")
		} else {
			display.DisplayTeam(team)
		}
	} else if choice == 2 {
		id := 0
		fmt.Print("Введите id команды: ")
		fmt.Scanln(&id)

		team, myErr := its.GetTeam(id, "", bl.SearchByID, user, itr)
		if myErr.ErrNum == bl.ErrAccessDenied {
			display.DisplayError("Ошибка: у Вас нет прав на выполнение данной операции")
		} else if myErr.ErrNum == bl.ErrGetNoteByName {
			display.DisplayError("Ошибка: попробуйте еще раз")
		} else {
			display.DisplayTeam(team)
		}
	}
}

func (MenuPoints) DisplayTeamMembers(user *models.User, ireps *bl.IRepositories, isvcs *bl.IServices) {
	itr := ireps.ITeamRepo
	its := isvcs.ITeamSvc

	id := 0
	fmt.Print("Введите id команды: ")
	fmt.Scanln(&id)

	members, myErr := its.GetTeamMembers(id, user, itr)
	if myErr.ErrNum == bl.ErrAccessDenied {
		display.DisplayError("Ошибка: у Вас нет прав на выполнение данной операции")
	} else if myErr.ErrNum == bl.ErrGetNoteByName {
		display.DisplayError("Ошибка: попробуйте еще раз")
	} else {
		display.DisplayTeamMembers(members)
	}
}

func (MenuPoints) DisplayAllTeams(user *models.User, ireps *bl.IRepositories, isvcs *bl.IServices) {
	itr := ireps.ITeamRepo
	its := isvcs.ITeamSvc

	teams, myErr := its.GetAllTeams(user, itr)
	if myErr.ErrNum == bl.ErrAccessDenied {
		display.DisplayError("Ошибка: у Вас нет прав на выполнение данной операции")
	} else if myErr.ErrNum == bl.ErrGetNoteByName {
		display.DisplayError("Ошибка: попробуйте еще раз")
	} else {
		display.DisplayAllTeams(teams)
	}
}

func (MenuPoints) AddUserToTeam(user *models.User, ireps *bl.IRepositories, isvcs *bl.IServices) {
	teamID := 0
	fmt.Print("Введите ID команды: ")
	fmt.Scanln(&teamID)

	userID := 0
	fmt.Print("Введите ID Пользователя: ")
	fmt.Scanln(&userID)

	itr := ireps.ITeamRepo
	its := isvcs.ITeamSvc

	myErr := its.AddUserToTeam(user, userID, teamID, itr)
	if myErr.ErrNum == bl.ErrAccessDenied {
		display.DisplayError("Ошибка: у Вас нет прав на выполнение данной операции")
	} else if myErr.ErrNum == bl.ErrDeleteTeam {
		display.DisplayError("Ошибка: попробуйте еще раз")
	} else {
		display.DisplaySuccess("Пользователя успешно добавлен в команду")
	}
}

func (MenuPoints) DeleteUserFromTeam(user *models.User, ireps *bl.IRepositories, isvcs *bl.IServices) {
	teamID := 0
	fmt.Print("Введите ID команды: ")
	fmt.Scanln(&teamID)

	userID := 0
	fmt.Print("Введите ID Пользователя: ")
	fmt.Scanln(&userID)

	itr := ireps.ITeamRepo
	its := isvcs.ITeamSvc

	myErr := its.DeleteUserFromTeam(user, userID, teamID, itr)
	if myErr.ErrNum == bl.ErrAccessDenied {
		display.DisplayError("Ошибка: у Вас нет прав на выполнение данной операции")
	} else if myErr.ErrNum == bl.ErrDeleteTeam {
		display.DisplayError("Ошибка: попробуйте еще раз")
	} else {
		display.DisplaySuccess("Пользователь успешно удален из команды")
	}
}
