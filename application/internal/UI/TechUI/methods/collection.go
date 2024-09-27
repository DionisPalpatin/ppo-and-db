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

func (MenuPoints) AddCollection(user *models.User, dateTimeFormat string, ireps *bl.IRepositories, isvcs *bl.IServices) {
	collName := ""
	fmt.Print("Введите название Подборки: ")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		collName = scanner.Text()
	}

	icr := ireps.IColRepo
	ics := isvcs.IColSvc

	collection := &models.Collection{Name: collName, CreationDate: time.Now(), OwnerID: user.Id}

	myErr := ics.AddCollection(collection, icr)
	if myErr.ErrNum == bl.ErrAccessDenied {
		display.DisplayError("Ошибка: у Вас нет прав на выполнение данной операции")
	} else if myErr.ErrNum == bl.ErrDeleteTeam {
		display.DisplayError("Ошибка: попробуйте еще раз")
	} else {
		display.DisplaySuccess("Подборка успешно добавлена")
	}
}

func (MenuPoints) DeleteCollection(user *models.User, ireps *bl.IRepositories, isvcs *bl.IServices) {
	id := 0
	fmt.Print("Введите id Подборки: ")
	fmt.Scanln(&id)

	icr := ireps.IColRepo
	ics := isvcs.IColSvc

	myErr := ics.DeleteCollection(id, user, icr)
	if myErr.ErrNum == bl.ErrAccessDenied {
		display.DisplayError("Ошибка: у Вас нет прав на выполнение данной операции")
	} else if myErr.ErrNum == bl.ErrDeleteTeam {
		display.DisplayError("Ошибка: попробуйте еще раз")
	} else {
		display.DisplaySuccess("Подборка успешно удалена")
	}
}

func (MenuPoints) DisplayAllUserCollections(user *models.User, ireps *bl.IRepositories, isvcs *bl.IServices) {
	icr := ireps.IColRepo
	ics := isvcs.IColSvc

	collections, myErr := ics.GetAllUsersCollections(user, icr)
	if myErr.ErrNum == bl.ErrAccessDenied {
		display.DisplayError("Ошибка: у Вас нет прав на выполнение данной операции")
	} else if myErr.ErrNum == bl.ErrGetNoteByName {
		display.DisplayError("Ошибка: попробуйте еще раз")
	} else {
		display.DisplayCollections(collections)
	}
}

func (MenuPoints) DisplayAllCollections(user *models.User, ireps *bl.IRepositories, isvcs *bl.IServices) {
	icr := ireps.IColRepo
	ics := isvcs.IColSvc

	collections, myErr := ics.GetAllCollections(user, icr)
	if myErr.ErrNum == bl.ErrAccessDenied {
		display.DisplayError("Ошибка: у Вас нет прав на выполнение данной операции")
	} else if myErr.ErrNum == bl.ErrGetNoteByName {
		display.DisplayError("Ошибка: попробуйте еще раз")
	} else {
		display.DisplayCollections(collections)
	}
}

func (MenuPoints) AddNoteToCollection(ireps *bl.IRepositories, isvcs *bl.IServices) {
	noteID := 0
	fmt.Print("Введите ID записки: ")
	fmt.Scanln(&noteID)

	collID := 0
	fmt.Print("Введите ID подборки: ")
	fmt.Scanln(&collID)

	inr := ireps.INoteRepo
	ins := isvcs.INoteSvc

	myErr := ins.AddNoteToCollection(noteID, collID, inr)
	if myErr.ErrNum == bl.ErrAccessDenied {
		display.DisplayError("Ошибка: у Вас нет прав на выполнение данной операции")
	} else if myErr.ErrNum == bl.ErrDeleteTeam {
		display.DisplayError("Ошибка: попробуйте еще раз")
	} else {
		display.DisplaySuccess("Записка успешно добавлен в подборку")
	}
}

func (MenuPoints) DeleteNoteFromCollection(user *models.User, ireps *bl.IRepositories, isvcs *bl.IServices) {
	noteID := 0
	fmt.Print("Введите ID команды: ")
	fmt.Scanln(&noteID)

	collID := 0
	fmt.Print("Введите ID Пользователя: ")
	fmt.Scanln(&collID)

	inr := ireps.INoteRepo
	ins := isvcs.INoteSvc

	myErr := ins.DeleteNoteFromCollection(noteID, collID, inr)
	if myErr.ErrNum == bl.ErrAccessDenied {
		display.DisplayError("Ошибка: у Вас нет прав на выполнение данной операции")
	} else if myErr.ErrNum == bl.ErrDeleteTeam {
		display.DisplayError("Ошибка: попробуйте еще раз")
	} else {
		display.DisplaySuccess("Записка успешно удален из подборки")
	}
}
