package methods

import (
	"bufio"
	"fmt"
	"os"

	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/UI/TechUI/display"
	bl "github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/business_logic"
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/models"
)

func (MenuPoints) DeleteUser(user *models.User, ireps *bl.IRepositories, isvcs *bl.IServices) {
	id := 0
	fmt.Print("Введите id пользователя: ")
	fmt.Scanln(&id)

	ur := ireps.IUsrRepo
	us := isvcs.IUsrSvc

	myErr := us.DeleteUser(user, id, ur)
	if myErr.ErrNum == bl.ErrAccessDenied {
		display.DisplayError("Ошибка: у Вас нет доступа к данному пользователю")
	} else if myErr.ErrNum == bl.ErrDeleteUser {
		display.DisplayError("Ошибка: попробуйте еще раз")
	} else {
		display.DisplaySuccess("Пользователь успешно удален")
	}
}

func (MenuPoints) ChangeUserFIO(user *models.User, ireps *bl.IRepositories, isvcs *bl.IServices) {
	id := 0
	fmt.Print("Введите id пользователя: ")
	fmt.Scanln(&id)

	fio := ""
	fmt.Print("Введите ФИО пользователя: ")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		fio = scanner.Text()
	}

	ur := ireps.IUsrRepo
	us := isvcs.IUsrSvc

	changedUser, myErr := us.GetUser(id, "", bl.SearchByID, user, ur)
	if myErr.ErrNum == bl.ErrAccessDenied {
		display.DisplayError("Ошибка: у Вас нет доступа к данному пользователю")
		return
	} else if myErr.ErrNum == bl.ErrGetUserByID {
		display.DisplayError("Ошибка: попробуйте еще раз")
		return
	}

	changedUser.Fio = fio
	myErr = us.UpdateUser(user, changedUser, ur)
	if myErr.ErrNum == bl.ErrAccessDenied {
		display.DisplayError("Ошибка: у Вас нет доступа к данному пользователю")
	} else if myErr.ErrNum == bl.ErrUpdateUser {
		display.DisplayError("Ошибка: попробуйте еще раз")
	} else {
		display.DisplaySuccess("Пользователь успешно удален")
	}
}

func (MenuPoints) ChangeUserRole(user *models.User, ireps *bl.IRepositories, isvcs *bl.IServices) {
	id := 0
	fmt.Print("Введите id пользователя: ")
	fmt.Scanln(&id)

	role := 2
	fmt.Print("Введите новую роль пользователя (0 -- Читатель, 1 -- Писатель, 2 -- Администратор): ")
	fmt.Scanln(&role)

	if role < 0 || role > 2 {
		display.DisplayError("Ошибка: попробуйте еще раз")
		return
	}

	ur := ireps.IUsrRepo
	us := isvcs.IUsrSvc

	changedUser, myErr := us.GetUser(id, "", bl.SearchByID, user, ur)
	if myErr.ErrNum == bl.ErrAccessDenied {
		display.DisplayError("Ошибка: у Вас нет доступа к данному пользователю")
		return
	} else if myErr.ErrNum == bl.ErrGetUserByID {
		display.DisplayError("Ошибка: попробуйте еще раз")
		return
	}

	changedUser.Role = role
	myErr = us.UpdateUser(user, changedUser, ur)
	if myErr.ErrNum == bl.ErrAccessDenied {
		display.DisplayError("Ошибка: у Вас нет доступа к данному пользователю")
	} else if myErr.ErrNum == bl.ErrUpdateUser {
		display.DisplayError("Ошибка: попробуйте еще раз")
	} else {
		display.DisplaySuccess("Пользователь успешно удален")
	}
}

func (MenuPoints) FindUser(user *models.User, ireps *bl.IRepositories, isvcs *bl.IServices) {
	choice := 0
	ius := isvcs.IUsrSvc
	iur := ireps.IUsrRepo

	for choice < 1 || choice > 3 {
		fmt.Println("\n1) Искать по имени")
		fmt.Println("2) Искать по ID")
		fmt.Println("3) Назад")
		fmt.Print("Введите номер пункта меню: ")
		fmt.Scanln(&choice)
	}

	if choice == 1 {
		name := ""
		fmt.Print("Введите имя пользователя: ")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			name = scanner.Text()
		}

		gotUser, myErr := ius.GetUser(0, name, bl.SearchByString, user, iur)
		if myErr.ErrNum == bl.ErrAccessDenied {
			display.DisplayError("Ошибка: у Вас нет прав на выполнение данной операции")
		} else if myErr.ErrNum == bl.ErrGetNoteByName {
			display.DisplayError("Ошибка: попробуйте еще раз")
		} else {
			display.DisplayUser(gotUser)
		}
	} else if choice == 2 {
		id := 0
		fmt.Print("Введите id пользователя: ")
		fmt.Scanln(&id)

		gotUser, myErr := ius.GetUser(id, "", bl.SearchByID, user, iur)
		if myErr.ErrNum == bl.ErrAccessDenied {
			display.DisplayError("Ошибка: у Вас нет прав на выполнение данной операции")
		} else if myErr.ErrNum == bl.ErrGetNoteByName {
			display.DisplayError("Ошибка: попробуйте еще раз")
		} else {
			display.DisplayUser(gotUser)
		}
	}
}

func (MenuPoints) DisplayAllUsers(user *models.User, ireps *bl.IRepositories, isvcs *bl.IServices) {
	ius := isvcs.IUsrSvc
	iur := ireps.IUsrRepo

	users, err := ius.GetAllUsers(user, iur)
	if err.ErrNum == bl.ErrAccessDenied {
		display.DisplayError("Ошибка: у Вас нет доступа к данной Записке")
	} else if err.ErrNum == bl.ErrGetNoteByName {
		display.DisplayError("Ошибка: попробуйте еще раз")
	} else {
		display.DisplayAllUsers(users)
	}
}
