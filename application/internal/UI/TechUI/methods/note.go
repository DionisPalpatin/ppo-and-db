package methods

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/config"
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/UI/TechUI/display"
	bl "github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/business_logic"
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/models"
)

func (MenuPoints) FindNote(user *models.User, ireps *bl.IRepositories, isvcs *bl.IServices) {
	var noteName string
	fmt.Print("Введите название заметки: ")
	fmt.Fscan(os.Stdin, &noteName)

	ins := isvcs.INoteSvc
	inr := ireps.INoteRepo
	isr := ireps.ISecRepo
	itr := ireps.ITeamRepo

	note, data, _, err := ins.GetNote(0, noteName, bl.SearchByString, user, inr, isr, itr)
	if err.ErrNum == bl.ErrAccessDenied {
		display.DisplayError("Ошибка: у Вас нет доступа к данной Записке")
	} else if err.ErrNum == bl.ErrGetNoteByName {
		display.DisplayError("Ошибка: Записка с таким названием не найдена")
	} else {
		display.DisplayNote(note, data)
	}
}

func (MenuPoints) DisplayAllOpenNotes(ireps *bl.IRepositories, isvcs *bl.IServices) {
	ins := isvcs.INoteSvc
	inr := ireps.INoteRepo

	notes, err := ins.GetAllNotes(false, &models.User{Role: bl.Admin}, inr)
	if err.ErrNum == bl.ErrAccessDenied {
		display.DisplayError("Ошибка: у Вас нет доступа к данной Записке")
	} else if err.ErrNum == bl.ErrGetNoteByName {
		display.DisplayError("Ошибка: попробуйте еще раз")
	} else {
		display.DisplayAllNotes(notes, 1)
	}
}

func (MenuPoints) DisplayNotesInSection(user *models.User, ireps *bl.IRepositories, isvcs *bl.IServices) {
	isr := ireps.ISecRepo
	iss := isvcs.ISecSvc
	itr := ireps.ITeamRepo

	var teamName string
	fmt.Println("Введите имя команды: ")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		teamName = scanner.Text()
	}

	sec, myErr := iss.GetSection(0, teamName, user, bl.SearchByString, isr)
	if myErr.ErrNum == bl.ErrGetSectionByTeamName {
		display.DisplayError("Ошибка: у этой команды нет своего Раздела либо такой Команды нет")
		return
	} else if myErr.ErrNum == bl.ErrAccessDenied {
		display.DisplayError("Ошибка: у Вас нет доступа к Разделу данной команды")
	}

	notes, myErr := iss.GetAllNotesInSection(sec.Id, user, isr, itr)
	if myErr.ErrNum == bl.ErrAccessDenied {
		display.DisplayError("Ошибка: у Вас нет доступа к какому-либо Разделу")
	} else if myErr.ErrNum == bl.ErrGetNoteByName {
		display.DisplayError("Ошибка: попробуйте еще раз")
	} else {
		display.DisplayNotesInSection(notes)
	}

}

func (MenuPoints) DisplayNotesInUserSection(user *models.User, ireps *bl.IRepositories, isvcs *bl.IServices) {
	isr := ireps.ISecRepo
	iss := isvcs.ISecSvc
	itr := ireps.ITeamRepo
	its := isvcs.ITeamSvc

	team, myErr := its.GetUserTeam(user, itr)
	if myErr.ErrNum == bl.ErrGetUserTeam {
		display.DisplayError("Ошибка: Вы не состоите ни в какой команде")
		return
	}

	sec, myErr := iss.GetSection(0, team.Name, user, bl.SearchByString, isr)
	if myErr.ErrNum == bl.ErrGetSectionByTeamName {
		display.DisplayError("Ошибка: у Вашей команды нет своего Раздела")
		return
	}

	notes, myErr := iss.GetAllNotesInSection(sec.Id, user, isr, itr)
	if myErr.ErrNum == bl.ErrAccessDenied {
		display.DisplayError("Ошибка: у Вас нет доступа к какому-либо Разделу")
	} else if myErr.ErrNum == bl.ErrGetNoteByName {
		display.DisplayError("Ошибка: попробуйте еще раз")
	} else {
		display.DisplayNotesInSection(notes)
	}
}

func (MenuPoints) AddNote(user *models.User, configs *config.Configs, ireps *bl.IRepositories, isvcs *bl.IServices) {
	newNote := models.Note{
		Id:          0,
		Access:      0,
		Name:        "",
		ContentType: 0,
		Likes:       0,
		Dislikes:    0,
		OwnerID:     0,
		SectionID:   0,
	}

	fmt.Print("Название Записки: ")
	fmt.Fscan(os.Stdin, &newNote.Name)

	var filePath string
	fmt.Print("Введите путь до файла с содержимым: ")
	fmt.Fscan(os.Stdin, &filePath)

	ind := strings.LastIndexByte(filePath, '.')
	if ind > 0 {
		fileType := filePath[ind+1:]

		i := 0
		for i = 0; i < len(configs.TextTypes) && fileType != configs.TextTypes[i]; i++ {
		}

		if i < len(configs.TextTypes) {
			newNote.ContentType = bl.TextCont
		} else {
			for i = 0; i < len(configs.ImageTypes) && fileType != configs.ImageTypes[i]; i++ {
			}

			if i < len(configs.ImageTypes) {
				newNote.ContentType = bl.ImgCont
			} else {
				newNote.ContentType = bl.RawData
			}
		}
	} else {
		newNote.ContentType = bl.RawData
	}

	newNote.RegistrationDate = time.Now()
	newNote.OwnerID = user.Id
	newNote.Access = bl.OpenCont
	newNote.SectionID = -1

	ins := isvcs.INoteSvc
	inr := ireps.INoteRepo
	myErr := ins.AddNote(&newNote, user, inr)
	if myErr.ErrNum == bl.ErrAccessDenied {
		display.DisplayError("Ошибка: у Вас нет прав на выполнение данной операции")
	} else if myErr.ErrNum == bl.ErrGetNoteByName {
		display.DisplayError("Ошибка: попробуйте еще раз")
	}
}

func (MenuPoints) DeleteNote(user *models.User, ireps *bl.IRepositories, isvcs *bl.IServices) {
	var id int
	fmt.Print("Введите id Записки: ")
	fmt.Fscan(os.Stdin, &id)

	inr := ireps.INoteRepo
	ins := isvcs.INoteSvc
	myErr := ins.DeleteNote(id, user, inr)
	if myErr.ErrNum == bl.ErrAccessDenied {
		display.DisplayError("Ошибка: у Вас нет прав на выполнение данной операции")
	} else if myErr.ErrNum == bl.ErrGetNoteByName {
		display.DisplayError("Ошибка: попробуйте еще раз")
	}
}

func (MenuPoints) DisplayAllNotes(user *models.User, ireps *bl.IRepositories, isvcs *bl.IServices) {
	inr := ireps.INoteRepo
	ins := isvcs.INoteSvc

	notes, myErr := ins.GetAllNotes(false, user, inr)
	if myErr.ErrNum == bl.ErrAccessDenied {
		display.DisplayError("Ошибка: у Вас нет прав на выполнение данной операции")
	} else if myErr.ErrNum == bl.ErrGetNoteByName {
		display.DisplayError("Ошибка: попробуйте еще раз")
	} else {
		display.DisplayAllNotes(notes, 0)
	}
}

func (MenuPoints) DisplayNotesInCollection(user *models.User, ireps *bl.IRepositories, isvcs *bl.IServices) {
	icr := ireps.IColRepo
	ics := isvcs.IColSvc

	colls, _ := ics.GetAllUsersCollections(user, icr)

	var choice int
	fmt.Print("Введите id коллекции: ")
	_, err := fmt.Fscan(os.Stdin, &choice)
	if err != nil {
		display.DisplayError("Ошибка: попробуйте еще раз")
	}
	choice--

	notes, myErr := ics.GetAllNotesInCollection(colls[choice], icr)
	if myErr.ErrNum == bl.ErrAccessDenied {
		display.DisplayError("Ошибка: у Вас нет прав на выполнение данной операции")
	} else if myErr.ErrNum == bl.ErrGetNoteByName {
		display.DisplayError("Ошибка: попробуйте еще раз")
	} else {
		display.DisplayAllNotes(notes, 0)
	}
}
