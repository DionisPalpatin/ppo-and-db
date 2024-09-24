package bl

import (
	"os"

	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/models"
)

type NoteService struct{}

func (NoteService) GetNote(id int, name string, searchBy int, requester *models.User, inr INoteRepository, isr ISectionRepository, itr ITeamRepository) (*models.Note, []byte, *MyError) {
	var note *models.Note
	var data []byte
	var myErr *MyError

	switch searchBy {
	case SearchByID:
		note, data, myErr = inr.GetNoteByID(id)

	case SearchByString:
		note, data, myErr = inr.GetNoteByName(name)

	default:
		myErr = CreateError(ErrSearchParameter, ErrSearchParameterError(), "GetUser")
		return nil, nil, myErr
	}

	if note.SectionID >= 0 {
		var section *models.Section
		section, myErr = isr.GetSectionByID(note.SectionID)
		if myErr.ErrNum != AllIsOk {
			return nil, nil, myErr
		}

		var team *models.Team
		team, myErr = itr.GetUserTeam(requester)
		if myErr.ErrNum != AllIsOk {
			return nil, nil, myErr
		}

		var sectionTeam *models.Section
		sectionTeam, myErr = isr.GetSectionByTeamName(team.Name)
		if myErr.ErrNum != AllIsOk {
			return nil, nil, myErr
		}

		if section.Id != sectionTeam.Id {
			myErr = CreateError(ErrAccessDenied, ErrAccessDeniedError(), "GetNote")
			return nil, nil, myErr
		}
	}

	return note, data, myErr
}

func (NoteService) GetAllNotes(open bool, requester *models.User, inr INoteRepository) ([]*models.Note, *MyError) {
	if requester.Role != Admin {
		return nil, CreateError(ErrAccessDenied, ErrAccessDeniedError(), "GetAllTeams")
	}

	if open {
		return inr.GetAllPublicNotes()
	} else {
		return inr.GetAllNotes()
	}

}

func (NoteService) AddNote(note *models.Note, requester *models.User, inr INoteRepository) *MyError {
	if requester.Role == Reader {
		return CreateError(ErrAccessDenied, ErrAccessDeniedError(), "GetAllTeams")
	}

	return inr.AddNote(note)
}

func (NoteService) DeleteNote(id int, requester *models.User, inr INoteRepository) *MyError {
	var note *models.Note
	note, _, err := inr.GetNoteByID(id)
	if err.ErrNum != AllIsOk {
		return err
	}

	if requester.Role == Reader || requester.Role == Author && requester.Id != note.OwnerID {
		return CreateError(ErrAccessDenied, ErrAccessDeniedError(), "GetAllTeams")
	}

	return inr.DeleteNote(id)
}

func (NoteService) UpdateNoteContent(noteID int, requester *models.User, filePath string, inr INoteRepository) *MyError {
	var note *models.Note
	note, _, err := inr.GetNoteByID(noteID)
	if err.ErrNum != AllIsOk {
		return err
	}

	if requester.Role == Reader || requester.Role == Author && requester.Id != note.OwnerID {
		err = CreateError(ErrAccessDenied, ErrAccessDeniedError(), "GetAllTeams")
		return err
	}

	file, err1 := os.Open(filePath)
	if err1 != nil {
		err = CreateError(ErrNoFile, ErrNoFileError(), "UpdateNoteContent")
		return err
	}

	if note.ContentType == TextCont {
		err = inr.UpdateNoteContentText(file, note)
	} else if note.ContentType == ImgCont {
		err = inr.UpdateNoteContentImg(file, note)
	} else {
		err = inr.UpdateNoteContentRawData(file, note)
	}

	return err
}

func (NoteService) UpdateNoteInfo(requester *models.User, note *models.Note, inr INoteRepository) *MyError {
	if requester.Role == Reader || requester.Role == Author && requester.Id != note.OwnerID {
		return CreateError(ErrAccessDenied, ErrAccessDeniedError(), "GetAllTeams")
	}

	return inr.UpdateNoteInfo(note)
}

func (NoteService) AddNoteToCollection(noteID int, collID int, inr INoteRepository) *MyError {
	return inr.AddNoteToCollection(collID, noteID)
}

func (NoteService) DeleteNoteFromCollection(noteID int, collID int, inr INoteRepository) *MyError {
	return inr.DeleteNoteFromCollection(collID, noteID)
}
