package TechUI

import (
	"fmt"
	"os"

	"notebook_app/config"
	bl "notebook_app/internal/business_logic"
	"notebook_app/internal/models"
)

func AuthorizationMenu(configs *config.Configs, ireps *bl.IRepositories, isvcs *bl.IServices) *models.User {
	for {
		fmt.Print("\n\n\n")
		fmt.Println("1) Зарегистрироваться")
		fmt.Println("2) Войти в аккаунт")
		fmt.Println("3) Выйти из программы")

		var choice string
		fmt.Print("Введите номер пункта меню: ")
		fmt.Scanln(&choice)

		var user *models.User

		switch choice {
		case "1":
			user = Registration(configs, ireps, isvcs)

		case "2":
			user = LogIn(configs, ireps, isvcs)

		case "3":
			user = nil
		}

		return user
	}
}

func LogIn(configs *config.Configs, ireps *bl.IRepositories, isvcs *bl.IServices) *models.User {
	var login string
	fmt.Print("Введите логин: ")
	fmt.Fscan(os.Stdin, &login)

	var password string
	fmt.Print("Введите пароль: ")
	fmt.Fscan(os.Stdin, &password)

	ioas := isvcs.IOAuthSvc
	iur := ireps.IUsrRepo
	user, err := ioas.SignInUser(login, password, iur)

	if err.ErrNum == bl.AllIsOk {
		if user.Password != password {
			fmt.Println("Неправильный логин или пароль")
			return nil
		}

		return user
	} else {
		fmt.Println("Неправильный логин или пароль")
	}

	return nil
}

func Registration(configs *config.Configs, ireps *bl.IRepositories, isvcs *bl.IServices) *models.User {
	user := new(models.User)

	var fio string
	fmt.Print("Введите Ваше ФИО: ")
	fmt.Fscan(os.Stdin, &fio)

	var login string
	fmt.Print("Введите Ваш логин: ")
	fmt.Fscan(os.Stdin, &login)

	var password string
	fmt.Print("Введите Ваш пароль: ")
	fmt.Fscan(os.Stdin, &password)

	ioas := isvcs.IOAuthSvc
	iur := ireps.IUsrRepo
	user, err := ioas.RegisterUser(fio, login, password, iur)

	if err.ErrNum == bl.AllIsOk {
		fmt.Println("Ошибка: попробуйте еще раз")
		return nil
	}

	return user
}
