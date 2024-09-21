package methods

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"notebook_app/config"
	"notebook_app/internal/UI/TechUI/display"
	bl "notebook_app/internal/business_logic"
	"notebook_app/internal/models"
)

func (MenuPoints) DisplayAllSections(user *models.User, ireps *bl.IRepositories, isvcs *bl.IServices) {
	isr := ireps.ISecRepo
	iss := isvcs.ISecSvc

	sections, myErr := iss.GetAllSections(user, isr)
	if myErr.ErrNum == bl.ErrAccessDenied {
		display.DisplayError("Ошибка: у Вас нет прав на выполнение данной операции")
	} else if myErr.ErrNum == bl.ErrGetNoteByName {
		display.DisplayError("Ошибка: попробуйте еще раз")
	} else {
		display.DisplayAllSections(sections)
	}
}

func (MenuPoints) DeleteSection(user *models.User, ireps *bl.IRepositories, isvcs *bl.IServices) {
	id := 0
	fmt.Print("Введите id раздела: ")
	fmt.Scanln(&id)

	isr := ireps.ISecRepo
	iss := isvcs.ISecSvc

	myErr := iss.DeleteSection(id, user, isr)
	if myErr.ErrNum == bl.ErrAccessDenied {
		display.DisplayError("Ошибка: у Вас нет доступа к данному разделу")
	} else if myErr.ErrNum == bl.ErrDeleteSection {
		display.DisplayError("Ошибка: попробуйте еще раз")
	} else {
		display.DisplaySuccess("Раздел успешно удален")
	}
}

func (MenuPoints) AddSection(user *models.User, configs *config.Configs, ireps *bl.IRepositories, isvcs *bl.IServices) {
	sec := &models.Section{
		Id:           0,
		CreationDate: time.Now().Format(configs.DateTimeFormat),
	}
	isr := ireps.ISecRepo
	iss := isvcs.ISecSvc
	itr := ireps.ITeamRepo
	its := isvcs.ITeamSvc

	var teamName string
	fmt.Print("Введите название команды, которой добавляется секция: ")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		teamName = scanner.Text()
	}

	team, myErr := its.GetTeam(0, teamName, bl.SearchByString, user, itr)
	if myErr.ErrNum == bl.ErrAccessDenied {
		display.DisplayError("Ошибка: у Вас нет прав на выполнение данной операции")
	} else if myErr.ErrNum == bl.ErrAddSection {
		display.DisplayError("Ошибка: попробуйте еще раз")
	} else {
		myErr = iss.AddSection(sec, team, user, isr)
		if myErr.ErrNum == bl.ErrAccessDenied {
			display.DisplayError("Ошибка: у Вас нет доступа к данной команде")
		} else if myErr.ErrNum == bl.ErrAddSection {
			display.DisplayError("Ошибка: попробуйте еще раз")
		} else {
			display.DisplaySuccess("Раздел успешно создан")
		}
	}

}

func (MenuPoints) AddNoteToSection(user *models.User, ireps *bl.IRepositories, isvcs *bl.IServices) {
	noteID := 0
	fmt.Print("Введите ID записки: ")
	fmt.Scanln(&noteID)

	secID := 0
	fmt.Print("Введите ID раздела: ")
	fmt.Scanln(&secID)

	isr := ireps.ISecRepo
	iss := isvcs.ISecSvc
	inr := ireps.INoteRepo
	ins := isvcs.INoteSvc
	itr := ireps.ITeamRepo

	note, _, myErr := ins.GetNote(noteID, "", bl.SearchByID, user, inr, isr, itr)
	if myErr.ErrNum == bl.ErrAccessDenied {
		display.DisplayError("Ошибка: у Вас нет прав на выполнение данной операции")
		return
	} else if myErr.ErrNum == bl.ErrDeleteTeam {
		display.DisplayError("Ошибка: попробуйте еще раз")
		return
	}

	sec, myErr := iss.GetSection(secID, "", user, bl.SearchByID, isr)
	if myErr.ErrNum == bl.ErrAccessDenied {
		display.DisplayError("Ошибка: у Вас нет прав на выполнение данной операции")
		return
	} else if myErr.ErrNum == bl.ErrDeleteTeam {
		display.DisplayError("Ошибка: попробуйте еще раз")
		return
	}

	myErr = iss.AddNoteToSection(sec, note, user, isr, itr)
	if myErr.ErrNum == bl.ErrAccessDenied {
		display.DisplayError("Ошибка: у Вас нет прав на выполнение данной операции")
	} else if myErr.ErrNum == bl.ErrDeleteTeam {
		display.DisplayError("Ошибка: попробуйте еще раз")
	} else {
		display.DisplaySuccess("Записка успешно добавлен в раздел")
	}
}

func (MenuPoints) DeleteNoteFromSection(user *models.User, ireps *bl.IRepositories, isvcs *bl.IServices) {
	noteID := 0
	fmt.Print("Введите ID записки: ")
	fmt.Scanln(&noteID)

	secID := 0
	fmt.Print("Введите ID раздела: ")
	fmt.Scanln(&secID)

	isr := ireps.ISecRepo
	iss := isvcs.ISecSvc
	inr := ireps.INoteRepo
	ins := isvcs.INoteSvc
	itr := ireps.ITeamRepo

	note, _, myErr := ins.GetNote(noteID, "", bl.SearchByID, user, inr, isr, itr)
	if myErr.ErrNum == bl.ErrAccessDenied {
		display.DisplayError("Ошибка: у Вас нет прав на выполнение данной операции")
		return
	} else if myErr.ErrNum == bl.ErrDeleteTeam {
		display.DisplayError("Ошибка: попробуйте еще раз")
		return
	}

	sec, myErr := iss.GetSection(secID, "", user, bl.SearchByID, isr)
	if myErr.ErrNum == bl.ErrAccessDenied {
		display.DisplayError("Ошибка: у Вас нет прав на выполнение данной операции")
		return
	} else if myErr.ErrNum == bl.ErrDeleteTeam {
		display.DisplayError("Ошибка: попробуйте еще раз")
		return
	}

	myErr = iss.DeleteNoteFromSection(sec, note, user, isr, itr)
	if myErr.ErrNum == bl.ErrAccessDenied {
		display.DisplayError("Ошибка: у Вас нет прав на выполнение данной операции")
	} else if myErr.ErrNum == bl.ErrDeleteTeam {
		display.DisplayError("Ошибка: попробуйте еще раз")
	} else {
		display.DisplaySuccess("Записка успешно удален из раздела")
	}
}
