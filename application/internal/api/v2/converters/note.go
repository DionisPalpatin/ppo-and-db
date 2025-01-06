package converters

import (
	"time"

	"github.com/DionisPalpatin/ppo-and-db/application/internal/api/v2/transport_models"
	"github.com/DionisPalpatin/ppo-and-db/application/internal/models"
)

func ToNoteInfo(note *models.Note) transport_models.NoteInfo {
	return transport_models.NoteInfo{
		ID:       note.Id,
		Name:     note.Name,
		Likes:    note.Likes,
		Dislikes: note.Dislikes,
	}
}

func ToNoteFullInfo(note *models.Note) transport_models.NoteFullInfo {
	return transport_models.NoteFullInfo{
		ID:               note.Id,
		Access:           note.Access,
		Name:             note.Name,
		Likes:            note.Likes,
		Dislikes:         note.Dislikes,
		RegistrationDate: note.RegistrationDate.String(),
		OwnerID:          note.OwnerID,
		SectionID:        note.SectionID,
	}
}

func ToNoteFullData(note *models.Note) transport_models.NoteFullData {
	noteFullInfo := ToNoteFullInfo(note)

	noteFullData := transport_models.NoteFullData{
		Info:          noteFullInfo,
		TextFile:      note.Content.Text,
		Image:         note.Content.Img,
		OtherFile:     note.Content.Raw,
		TxtFileType:   note.Content.TextExt,
		ImgFileType:   note.Content.ImgExt,
		OtherFileType: note.Content.RawExt,
	}

	return noteFullData
}

func toContent(note *transport_models.NoteFullData) models.Content {
	return models.Content{
		Text: note.TextFile,
		Img:  note.Image,
		Raw:  note.OtherFile,

		TextExt: note.TxtFileType,
		ImgExt:  note.ImgFileType,
		RawExt:  note.OtherFileType,
	}
}

func FromNoteInfo(noteInfo *transport_models.NoteInfo) models.Note {
	return models.Note{
		Id:       noteInfo.ID,
		Name:     noteInfo.Name,
		Likes:    noteInfo.Likes,
		Dislikes: noteInfo.Dislikes,
	}
}

func FromNoteFullInfo(noteFullInfo *transport_models.NoteFullInfo) (models.Note, error) {
	parsedTime, err := time.Parse("2006-01-02 15:04:05-07", noteFullInfo.RegistrationDate)
	if err != nil {
		return models.Note{}, err
	}

	return models.Note{
		Id:               noteFullInfo.ID,
		Access:           noteFullInfo.Access,
		Name:             noteFullInfo.Name,
		Likes:            noteFullInfo.Likes,
		Dislikes:         noteFullInfo.Dislikes,
		RegistrationDate: parsedTime,
		OwnerID:          noteFullInfo.OwnerID,
		SectionID:        noteFullInfo.SectionID,
	}, nil
}

func FromNoteFullData(note *transport_models.NoteFullData) (models.Note, error) {
	noteInfo, err := FromNoteFullInfo(&note.Info)
	if err != nil {
		return models.Note{}, err
	}

	content := toContent(note)
	noteModel := models.Note{
		Id:               noteInfo.Id,
		Access:           noteInfo.Access,
		Name:             noteInfo.Name,
		Likes:            noteInfo.Likes,
		Dislikes:         noteInfo.Dislikes,
		RegistrationDate: noteInfo.RegistrationDate,
		OwnerID:          noteInfo.OwnerID,
		SectionID:        noteInfo.SectionID,
		Content:          content,
	}

	return noteModel, nil
}
