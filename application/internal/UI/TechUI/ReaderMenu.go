package TechUI

import (
	"fmt"

	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/config"
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/UI/TechUI/methods"
	bl "github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/business_logic"
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/models"
)

func ReaderMenu(user *models.User, configs *config.Configs, ireps *bl.IRepositories, isvcs *bl.IServices) {
	for {
		points := methods.MenuPoints{}

		fmt.Println("\n# Меню для читателя:")

		fmt.Println("\n## Меню Записок:")
		fmt.Println("1) Найти Записку")
		fmt.Println("2) Вывести все общедоступные Записки")
		fmt.Println("3) Вывести все Записки в коллекции")
		fmt.Println("4) Вывести все Записки в разделе своей команды")

		fmt.Println("\n## Меню Подборок:")
		fmt.Println("5) Добавить Подборку")
		fmt.Println("6) Удалить Подборку")
		fmt.Println("7) Вывести все свои Подборки")
		fmt.Println("8) Добавить Записку в Подборку")
		fmt.Println("9) Удалить Записку из Подборки")

		fmt.Println("\n10) Выйти из аккаунта")

		var choice string
		fmt.Print("Введите номер пункта меню: ")
		fmt.Scanln(choice)

		switch choice {
		case "1":
			points.FindNote(user, ireps, isvcs)

		case "2":
			points.DisplayAllOpenNotes(ireps, isvcs)

		case "3":
			points.DisplayNotesInCollection(user, ireps, isvcs)

		case "4":
			points.DisplayNotesInUserSection(user, ireps, isvcs)

		case "5":
			points.AddCollection(user, configs.DateTimeFormat, ireps, isvcs)

		case "6":
			points.DeleteCollection(nil, ireps, isvcs)

		case "7":
			points.DisplayAllUserCollections(user, ireps, isvcs)

		case "8":
			points.AddNoteToCollection(ireps, isvcs)

		case "9":
			points.DeleteNoteFromCollection(user, ireps, isvcs)

		case "10":
			fmt.Println("Выход из аккаунта...")
			return

		default:
			fmt.Println("Неверный выбор. Попробуйте снова.")
		}
	}
}
