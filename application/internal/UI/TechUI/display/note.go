package display

import (
	"fmt"

	bl "github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/business_logic"
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/models"
)

func DisplayNoteInfo(note *models.Note) {
	if note != nil {
		fmt.Println("## Записка:")
		fmt.Printf("ID: %d\n", note.Id)
		fmt.Printf("Доступ: %d\n", note.Access)
		fmt.Printf("Название: %s\n", note.Name)
		fmt.Printf("Тип контента: %d\n", note.ContentType)
		fmt.Printf("Лайки: %d\n", note.Likes)
		fmt.Printf("Дизлайки: %d\n", note.Dislikes)
		fmt.Printf("Дата регистрации: %s\n", note.RegistrationDate)
		fmt.Printf("ID владельца: %d\n", note.OwnerID)
		fmt.Printf("ID раздела: %d\n\n", note.SectionID)
	}
}

func DisplayNote(note *models.Note, data []byte) {
	if note != nil && len(data) > 0 {
		DisplayNoteInfo(note)
		if note.ContentType == bl.TextCont {
			fmt.Printf("%s\n", string(data))
		} else if note.ContentType == bl.ImgCont {
			fmt.Println("Вывести содержимое этой Записки в консоль на возможно.")
		} else {
			fmt.Println("Вывести содержимое этой Записки в консоль на возможно.")
		}
	}
}

func DisplayAllNotes(notes []*models.Note, openOnly int) {
	fmt.Println("## Записки:")
	i := 1
	for _, note := range notes {
		if openOnly == 1 {
			if note.Access == bl.OpenCont {
				fmt.Printf("%d) %s\n", i, note.Name)
			}
		}
	}
}
