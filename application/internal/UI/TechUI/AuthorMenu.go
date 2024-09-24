package TechUI

import (
	"fmt"

	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/config"
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/UI/TechUI/methods"
	bl "github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/business_logic"
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/models"
)

func AuthorMenu(user *models.User, configs *config.Configs, ireps *bl.IRepositories, isvcs *bl.IServices) {
	for {
		points := methods.MenuPoints{}

		fmt.Println("\n# Меню для автора:")

		fmt.Println("\n## Меню Записок:")
		fmt.Println("1) Добавить Записку")
		fmt.Println("2) Удалить Записку")
		fmt.Println("3) Найти Записку")
		fmt.Println("4) Вывести все общедоступные Записки")
		fmt.Println("5) Вывести все Записки в коллекции")
		fmt.Println("6) Вывести все Записки в разделе своей команды")

		fmt.Println("\n## Меню Подборок:")
		fmt.Println("7) Добавить Подборку")
		fmt.Println("8) Удалить Подборку")
		fmt.Println("9) Вывести все свои Подборки")
		fmt.Println("10) Добавить Записку в Подборку")
		fmt.Println("11) Удалить Записку из Подборки")

		fmt.Println("\n## Меню Разделов:")
		fmt.Println("12) Добавить Записку в раздел")
		fmt.Println("13) Удалить Записку из раздела")

		fmt.Println("\n14) Выйти из аккаунта")

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
			points.DisplayNotesInCollection(user, ireps, isvcs)

		case "6":
			points.DisplayNotesInUserSection(user, ireps, isvcs)

		case "7":
			points.AddCollection(user, configs.DateTimeFormat, ireps, isvcs)

		case "8":
			points.DeleteCollection(user, ireps, isvcs)

		case "9":
			points.DisplayAllUserCollections(user, ireps, isvcs)

		case "10":
			points.AddNoteToCollection(ireps, isvcs)

		case "11":
			points.DeleteNoteFromCollection(user, ireps, isvcs)

		case "12":
			points.AddNoteToSection(user, ireps, isvcs)

		case "13":
			points.DeleteNoteFromSection(user, ireps, isvcs)

		case "14":
			fmt.Println("Выход из аккаунта...")
			return

		default:
			fmt.Println("Неверный выбор. Попробуйте снова.")
		}
	}
}
