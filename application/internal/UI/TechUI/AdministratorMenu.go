package TechUI

import (
	"fmt"

	"notebook_app/config"
	"notebook_app/internal/UI/TechUI/methods"
	bl "notebook_app/internal/business_logic"
	"notebook_app/internal/models"
)

func AdminMenu(user *models.User, configs *config.Configs, ireps *bl.IRepositories, isvcs *bl.IServices) {
	for {
		points := methods.MenuPoints{}

		fmt.Println("\n# Меню Администратора:")

		fmt.Println("\n## Меню Записок:")
		fmt.Println("1) Добавить Записку")
		fmt.Println("2) Удалить Записку")
		fmt.Println("3) Найти Записку")
		fmt.Println("4) Вывести все общедоступные Записки")
		fmt.Println("5) Вывести все Записки")
		fmt.Println("6) Вывести все Записки в коллекции")
		fmt.Println("7) Вывести все Записки в разделе своей команды")
		fmt.Println("8) Вывести все Записки в разделе команды")

		fmt.Println("\n## Меню Подборок:")
		fmt.Println("9) Добавить Подборку")
		fmt.Println("10) Удалить Подборку")
		fmt.Println("11) Вывести все свои Подборки")
		fmt.Println("12) Вывести Подборки всех пользователей")
		fmt.Println("13) Добавить Записку в Подборку")
		fmt.Println("14) Удалить Записку из Подборки")

		fmt.Println("\n## Меню Пользователей:")
		fmt.Println("15) Удалить пользователя")
		fmt.Println("16) Изменить ФИО пользователя")
		fmt.Println("17) Изменить роль пользователя")
		fmt.Println("18) Найти пользователя")
		fmt.Println("19) Вывести всех пользователей")

		fmt.Println("\n## Меню Команд:")
		fmt.Println("20) Добавить команду")
		fmt.Println("21) Удалить команду")
		fmt.Println("22) Найти команду")
		fmt.Println("23) Вывести всех членов команды")
		fmt.Println("24) Вывести все команды")
		fmt.Println("25) Добавить пользователя в команду")
		fmt.Println("26) Удалить пользователя из команды")

		fmt.Println("\n## Меню Разделов:")
		fmt.Println("27) Добавить раздел")
		fmt.Println("28) Удалить раздел")
		fmt.Println("29) Вывести все разделы")
		fmt.Println("30) Добавить Записку в раздел")
		fmt.Println("31) Удалить Записку из раздела")

		fmt.Println("\n32) Выйти из аккаунта")

		var choice string
		fmt.Print("\n\nВведите номер пункта меню: ")
		fmt.Scanln(&choice)

		switch choice {
		case "1":
			points.AddNote(user, configs, ireps, isvcs)

		case "2":
			points.DeleteNote(user, ireps, isvcs)

		case "3":
			points.FindNote(user, ireps, isvcs)

		case "4":
			points.DisplayAllOpenNotes(ireps, isvcs)

		case "5":
			points.DisplayAllNotes(user, ireps, isvcs)

		case "6":
			points.DisplayNotesInCollection(user, ireps, isvcs)

		case "7":
			points.DisplayNotesInUserSection(user, ireps, isvcs)

		case "8":
			points.DisplayNotesInSection(user, ireps, isvcs)

		case "9":
			points.AddCollection(user, configs.DateTimeFormat, ireps, isvcs)

		case "10":
			points.DeleteCollection(user, ireps, isvcs)

		case "11":
			points.DisplayAllUserCollections(user, ireps, isvcs)

		case "12":
			points.DisplayAllCollections(user, ireps, isvcs)

		case "13":
			points.AddNoteToCollection(ireps, isvcs)

		case "14":
			points.DeleteNoteFromCollection(user, ireps, isvcs)

		case "15":
			points.DeleteUser(user, ireps, isvcs)

		case "16":
			points.ChangeUserFIO(user, ireps, isvcs)

		case "17":
			points.ChangeUserRole(user, ireps, isvcs)

		case "18":
			points.FindUser(user, ireps, isvcs)

		case "19":
			points.DisplayAllUsers(user, ireps, isvcs)

		case "20":
			points.AddTeam(user, configs.DateTimeFormat, ireps, isvcs)

		case "21":
			points.DeleteTeam(user, ireps, isvcs)

		case "22":
			points.FindTeam(user, ireps, isvcs)

		case "23":
			points.DisplayTeamMembers(user, ireps, isvcs)

		case "24":
			points.DisplayAllTeams(user, ireps, isvcs)

		case "25":
			points.AddUserToTeam(user, ireps, isvcs)

		case "26":
			points.DeleteUserFromTeam(user, ireps, isvcs)

		case "27":
			points.AddSection(user, configs, ireps, isvcs)

		case "28":
			points.DeleteSection(user, ireps, isvcs)

		case "29":
			points.DisplayAllSections(user, ireps, isvcs)

		case "30":
			points.AddNoteToSection(user, ireps, isvcs)

		case "31":
			points.DeleteNoteFromSection(user, ireps, isvcs)

		case "32":
			print("Выход из аккаунта")
			break

		default:
			print("Неверный выбор. Попробуйте снова.")
		}
	}
}
