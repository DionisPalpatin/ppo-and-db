package display

import (
	"fmt"

	"notebook_app/internal/models"
)

func DisplayCollectionInfo(collection *models.Collection) {
	if collection != nil {
		fmt.Println("## Подборка:")
		fmt.Printf("ID: %d\n", collection.Id)
		fmt.Printf("Название: %s\n", collection.Name)
		fmt.Printf("Дата создания: %s\n", collection.CreationDate)
		fmt.Printf("ID владельца: %d\n\n", collection.OwnerID)
	}
}

func DisplayCollections(colls []*models.Collection) {
	fmt.Println("## Подборки:")
	for i, col := range colls {
		fmt.Printf("%d) %s\n", i, col.Name)
	}
}

func DisplayNotesInCollection(notes []*models.Note) {
	fmt.Println("## Записке в подборке:")
	for i, note := range notes {
		fmt.Printf("%d) %s\n", i, note.Name)
	}
}
