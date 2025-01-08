package bl

import (
	"io"
	"os"
	"path"

	"github.com/DionisPalpatin/ppo-and-db/application/internal/models"
)

type NoteService struct {
	inr INoteRepository
	// isr ISectionRepository
	itr ITeamRepository
}

func (ns *NoteService) GetNote(id int, name string, searchBy int, requester *models.User) (*models.Note, *MyError) {
	var note *models.Note
	var myErr *MyError

	switch searchBy {
	case SearchByID:
		note, myErr = ns.inr.GetNoteByID(id)

	case SearchByString:
		note, myErr = ns.inr.GetNoteByName(name)

	default:
		myErr = CreateError(ErrSearchParameter, "GetNote", "bl")
		return nil, myErr
	}

	// if note.SectionID >= 0 {
	// 	var section *models.Section
	// 	section, myErr = ns.isr.GetSectionByID(note.SectionID)
	// 	if myErr.ErrNum != Ok {
	// 		return nil, myErr
	// 	}

	// 	var team *models.Team
	// 	team, myErr = ns.itr.GetUserTeam(requester)
	// 	if myErr.ErrNum != Ok {
	// 		return nil, myErr
	// 	}

	// 	var sectionTeam *models.Section
	// 	sectionTeam, myErr = ns.isr.GetSectionByTeamName(team.Name)
	// 	if myErr.ErrNum != Ok {
	// 		return nil, myErr
	// 	}

	// 	if section.Id != sectionTeam.Id {
	// 		myErr = CreateError(ErrAccessDenied, "GetNote", "bl")
	// 		return nil, myErr
	// 	}
	// }

	return note, myErr
}

func (ns *NoteService) GetAllNotes(open bool, requester *models.User) ([]*models.Note, *MyError) {
	if requester.Role != Admin {
		return nil, CreateError(ErrAccessDenied, "GetAllNotes", "bl")
	}

	if open {
		return ns.inr.GetAllPublicNotes()
	} else {
		return ns.inr.GetAllNotes()
	}

}

func (ns *NoteService) AddNote(note *models.Note, requester *models.User) (int, *MyError) {
	if requester.Role == Reader {
		return 0, CreateError(ErrAccessDenied, "AddNote", "bl")
	}

	return ns.inr.AddNote(note)
}

func (ns *NoteService) DeleteNote(id int, requester *models.User) *MyError {
	var note *models.Note
	note, err := ns.inr.GetNoteByID(id)
	if err.ErrNum != Ok {
		return err
	}

	if requester.Role == Reader || requester.Role == Author && requester.Id != note.OwnerID {
		return CreateError(ErrAccessDenied, "DeleteNote", "bl")
	}

	return ns.inr.DeleteNote(id)
}

func (ns *NoteService) UpdateNote(note *models.Note, requester *models.User, textFilePath string, imgFilePath string, rawFilePath string) *MyError {
	var myErr *MyError
	var err error
	var dataFile *os.File
	var content models.Content

	if requester.Role == Reader || requester.Role == Author && requester.Id != note.OwnerID {
		myErr = CreateError(ErrAccessDenied, "UpdateNote", "bl")
		return myErr
	}

	dataFile, err = os.Open(textFilePath)
	if err != nil {
		myErr = CreateError(ErrNoFile, "UpdateNote", "bl")
		return myErr
	} else {
		content.Text, err = io.ReadAll(dataFile)
		if err != nil {
			myErr = CreateError(ErrReadFile, "UpdateNote", "bl")
			return myErr
		}
		content.TextExt = path.Ext(textFilePath)
	}

	dataFile, err = os.Open(imgFilePath)
	if err != nil {
		myErr = CreateError(ErrNoFile, "UpdateNote", "bl")
		return myErr
	} else {
		content.Img, err = io.ReadAll(dataFile)
		if err != nil {
			myErr = CreateError(ErrReadFile, "UpdateNote", "bl")
			return myErr
		}
		content.ImgExt = path.Ext(textFilePath)
	}

	dataFile, err = os.Open(rawFilePath)
	if err != nil {
		myErr = CreateError(ErrNoFile, "UpdateNote", "bl")
		return myErr
	} else {
		content.Raw, err = io.ReadAll(dataFile)
		if err != nil {
			myErr = CreateError(ErrReadFile, "UpdateNote", "bl")
			return myErr
		}
		content.RawExt = path.Ext(textFilePath)
	}

	note.Content = content

	myErr = ns.inr.UpdateNoteContent(note)
	if myErr.ErrNum != Ok {
		return myErr
	}

	return ns.inr.UpdateNoteInfo(note)
}

func (ns *NoteService) AddNoteToCollection(noteID int, collID int) *MyError {
	return ns.inr.AddNoteToCollection(collID, noteID)
}

func (ns *NoteService) DeleteNoteFromCollection(noteID int, collID int) *MyError {
	return ns.inr.DeleteNoteFromCollection(collID, noteID)
}
