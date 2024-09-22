package display

import (
	"fmt"

	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/models"
)

func DisplayNotesInSection(notes []*models.Note) {
	for i, note := range notes {
		fmt.Printf("%d) %s\n", i, note.Name)
	}
}

func DisplaySection(section *models.Section) {
	if section != nil {
		fmt.Printf("%s\n", section.CreationDate)
	}
}

func DisplayAllSections(section []*models.Section) {
	for i, sec := range section {
		fmt.Printf("%d) %s\n", i, sec.CreationDate)
	}
}
