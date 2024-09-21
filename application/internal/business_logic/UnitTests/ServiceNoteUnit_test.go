package UnitTests

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"notebook_app/internal/business_logic"
	"notebook_app/internal/business_logic/UnitTests/mocks"
)

func TestGetNote(t *testing.T) {
	t.Run("SuccessGetNote", func(t *testing.T) {
		retErr := bl.CreateError(bl.AllIsOk, nil, "")
		retNote := &bl.Note{Id: 1, SectionID: 1, OwnerID: 1}
		retSection := &bl.Section{Id: 1, CommandID: 1}
		requester := &bl.User{Id: 1, Role: bl.Admin}

		mockNoteRepo := new(mocks.MockINoteRepository)
		mockNoteRepo.On("GetNoteByID", 1).Return(retNote, retErr)
		mockSectionRepo := new(mocks.MockISectionRepository)
		mockSectionRepo.On("GetSectionByID", 1).Return(retSection, retErr)

		nsSrv := bl.NoteService{}
		note, err := nsSrv.GetNote(1, requester, mockNoteRepo, mockSectionRepo)

		assert.NotNil(t, err)
		assert.NotNil(t, note)
		assert.Equal(t, err.ErrNum, bl.AllIsOk)

		mockNoteRepo.AssertExpectations(t)
		mockSectionRepo.AssertExpectations(t)
	})

	t.Run("ErrorGetNoteByID", func(t *testing.T) {
		retErr := bl.CreateError(bl.ErrGetNoteByID, bl.ErrGetNoteByIDError(), "GetNoteByID")
		requester := &bl.User{Id: 1, Role: bl.Admin}
		retNote := &bl.Note{}

		mockNoteRepo := new(mocks.MockINoteRepository)
		mockSectionRepo := new(mocks.MockISectionRepository)
		mockNoteRepo.On("GetNoteByID", 1).Return(retNote, retErr)

		nsSrv := bl.NoteService{}
		_, err := nsSrv.GetNote(1, requester, mockNoteRepo, mockSectionRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrGetNoteByID, err.ErrNum)

		mockNoteRepo.AssertExpectations(t)
	})

	t.Run("ErrorGetSectionByID", func(t *testing.T) {
		retErr := bl.CreateError(bl.ErrGetSectionByID, bl.ErrGetSectionByIDError(), "GetSectionByID")
		retNote := &bl.Note{Id: 1, SectionID: 1, OwnerID: 1}
		requester := &bl.User{Id: 1, CommandID: 1, Role: bl.Admin}
		returnSec := &bl.Section{}

		mockNoteRepo := new(mocks.MockINoteRepository)
		mockNoteRepo.On("GetNoteByID", 1).Return(retNote, bl.CreateError(bl.AllIsOk, nil, ""))
		mockSectionRepo := new(mocks.MockISectionRepository)
		mockSectionRepo.On("GetSectionByID", 1).Return(returnSec, retErr)

		nsSrv := bl.NoteService{}
		_, err := nsSrv.GetNote(1, requester, mockNoteRepo, mockSectionRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrGetSectionByID, err.ErrNum)

		mockNoteRepo.AssertExpectations(t)
		mockSectionRepo.AssertExpectations(t)
	})

	t.Run("ErrorAccessDenied", func(t *testing.T) {
		retOk := bl.CreateError(bl.AllIsOk, nil, "")
		retNote := &bl.Note{Id: 1, SectionID: 1, OwnerID: 1}
		retSection := &bl.Section{Id: 1, CommandID: 2}
		requester := &bl.User{Id: 1, CommandID: 1, Role: bl.Admin}

		mockNoteRepo := new(mocks.MockINoteRepository)
		mockNoteRepo.On("GetNoteByID", 1).Return(retNote, retOk)
		mockSectionRepo := new(mocks.MockISectionRepository)
		mockSectionRepo.On("GetSectionByID", 1).Return(retSection, retOk)

		nsSrv := bl.NoteService{}
		_, err := nsSrv.GetNote(1, requester, mockNoteRepo, mockSectionRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrAccessDenied, err.ErrNum)

		mockNoteRepo.AssertExpectations(t)
		mockSectionRepo.AssertExpectations(t)
	})
}

func TestGetAllNotes(t *testing.T) {
	t.Run("SuccessGetAllNotes", func(t *testing.T) {
		retErr := bl.CreateError(bl.AllIsOk, nil, "")
		retNotes := []*bl.Note{
			{Id: 1, OwnerID: 1},
			{Id: 2, OwnerID: 2},
		}
		requester := &bl.User{Id: 1, CommandID: 1, Role: bl.Admin}

		mockNoteRepo := new(mocks.MockINoteRepository)
		mockNoteRepo.On("GetAllNotes").Return(retNotes, retErr)

		nsSrv := bl.NoteService{}
		notes, err := nsSrv.GetAllNotes(requester, mockNoteRepo)

		assert.NotNil(t, err)
		assert.NotNil(t, notes)
		assert.Equal(t, err.ErrNum, bl.AllIsOk)

		mockNoteRepo.AssertExpectations(t)
	})

	t.Run("ErrorAccessDenied", func(t *testing.T) {
		requester := &bl.User{Id: 1, CommandID: 1, Role: bl.Reader}

		mockNoteRepo := new(mocks.MockINoteRepository)

		nsSrv := bl.NoteService{}
		_, err := nsSrv.GetAllNotes(requester, mockNoteRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrAccessDenied, err.ErrNum)
	})
}

func TestAddNote(t *testing.T) {
	t.Run("SuccessAddNote", func(t *testing.T) {
		retErr := bl.CreateError(bl.AllIsOk, nil, "")
		retNote := &bl.Note{Id: 1, OwnerID: 1}
		requester := &bl.User{Id: 1, CommandID: 1, Role: bl.Admin}

		mockNoteRepo := new(mocks.MockINoteRepository)
		mockNoteRepo.On("AddNote", retNote).Return(retErr)

		nsSrv := bl.NoteService{}
		err := nsSrv.AddNote(retNote, requester, mockNoteRepo)

		assert.NotNil(t, err)
		assert.Equal(t, err.ErrNum, bl.AllIsOk)

		mockNoteRepo.AssertExpectations(t)
	})

	t.Run("ErrorAccessDenied", func(t *testing.T) {
		retNote := &bl.Note{Id: 1, OwnerID: 1}
		requester := &bl.User{Id: 1, CommandID: 1, Role: bl.Reader}

		mockNoteRepo := new(mocks.MockINoteRepository)

		nsSrv := bl.NoteService{}
		err := nsSrv.AddNote(retNote, requester, mockNoteRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrAccessDenied, err.ErrNum)
	})
}

func TestDeleteNote(t *testing.T) {
	t.Run("SuccessDeleteNote", func(t *testing.T) {
		retOk := bl.CreateError(bl.AllIsOk, nil, "")
		retNote := &bl.Note{Id: 1, OwnerID: 1}
		requester := &bl.User{Id: 1, CommandID: 1, Role: bl.Admin}

		mockNoteRepo := new(mocks.MockINoteRepository)
		mockNoteRepo.On("GetNoteByID", 1).Return(retNote, retOk)
		mockNoteRepo.On("DeleteNote", 1).Return(retOk)

		nsSrv := bl.NoteService{}
		err := nsSrv.DeleteNote(1, requester, mockNoteRepo)

		assert.NotNil(t, err)
		assert.Equal(t, err.ErrNum, bl.AllIsOk)

		mockNoteRepo.AssertExpectations(t)
	})

	t.Run("ErrorGetNoteByID", func(t *testing.T) {
		retErr := bl.CreateError(bl.ErrGetNoteByID, bl.ErrGetNoteByIDError(), "GetNoteByID")
		requester := &bl.User{Id: 1, CommandID: 1, Role: bl.Admin}
		retNote := &bl.Note{}

		mockNoteRepo := new(mocks.MockINoteRepository)
		mockNoteRepo.On("GetNoteByID", 1).Return(retNote, retErr)

		nsSrv := bl.NoteService{}
		err := nsSrv.DeleteNote(1, requester, mockNoteRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrGetNoteByID, err.ErrNum)

		mockNoteRepo.AssertExpectations(t)
	})

	t.Run("ErrorAccessDenied", func(t *testing.T) {
		retOk := bl.CreateError(bl.AllIsOk, nil, "")
		retNote := &bl.Note{Id: 1, OwnerID: 2}
		requester := &bl.User{Id: 1, CommandID: 1, Role: bl.Author}

		mockNoteRepo := new(mocks.MockINoteRepository)
		mockNoteRepo.On("GetNoteByID", 1).Return(retNote, retOk)

		nsSrv := bl.NoteService{}
		err := nsSrv.DeleteNote(1, requester, mockNoteRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrAccessDenied, err.ErrNum)

		mockNoteRepo.AssertExpectations(t)
	})
}

func TestUpdateNoteContent(t *testing.T) {
	t.Run("SuccessUpdateNoteContentText", func(t *testing.T) {
		retErr := bl.CreateError(bl.AllIsOk, nil, "")
		retNote := &bl.Note{Id: 1, OwnerID: 1, ContentType: bl.TextCont}
		requester := &bl.User{Id: 1, CommandID: 1, Role: bl.Admin}

		filePath := "./Tests files/test.txt"
		file, err := os.Open(filePath)
		if err != nil {
			t.Fatal(err)
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				t.Fatal(err)
			}
		}(file)

		mockNoteRepo := new(mocks.MockINoteRepository)
		mockNoteRepo.On("GetNoteByID", 1).Return(retNote, bl.CreateError(bl.AllIsOk, nil, ""))
		mockNoteRepo.On("UpdateNoteContentText", mock.Anything, retNote).Return(retErr)

		nsSrv := bl.NoteService{}
		myErr := nsSrv.UpdateNoteContent(1, requester, filePath, mockNoteRepo)

		assert.Nil(t, err)
		assert.NotNil(t, myErr)
		assert.Equal(t, myErr.ErrNum, bl.AllIsOk)

		mockNoteRepo.AssertExpectations(t)
	})

	t.Run("SuccessUpdateNoteContentImg", func(t *testing.T) {
		retErr := bl.CreateError(bl.AllIsOk, nil, "")
		retNote := &bl.Note{Id: 1, OwnerID: 1, ContentType: bl.ImgCont}
		requester := &bl.User{Id: 1, CommandID: 1, Role: bl.Admin}

		filePath := "./Tests files/test.jpg"
		file, err := os.Open(filePath)
		if err != nil {
			t.Fatal(err)
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				t.Fatal(err)
			}
		}(file)

		mockNoteRepo := new(mocks.MockINoteRepository)
		mockNoteRepo.On("GetNoteByID", 1).Return(retNote, bl.CreateError(bl.AllIsOk, nil, ""))
		mockNoteRepo.On("UpdateNoteContentImg", mock.Anything, retNote).Return(retErr)

		nsSrv := bl.NoteService{}
		myErr := nsSrv.UpdateNoteContent(1, requester, filePath, mockNoteRepo)

		assert.Nil(t, err)
		assert.NotNil(t, myErr)
		assert.Equal(t, myErr.ErrNum, bl.AllIsOk)

		mockNoteRepo.AssertExpectations(t)
	})

	t.Run("SuccessUpdateNoteContentRawData", func(t *testing.T) {
		retErr := bl.CreateError(bl.AllIsOk, nil, "")
		retNote := &bl.Note{Id: 1, OwnerID: 1, ContentType: bl.RawData}
		requester := &bl.User{Id: 1, CommandID: 1, Role: bl.Admin}

		filePath := "./Tests files/test.bin"
		file, err := os.Open(filePath)
		if err != nil {
			t.Fatal(err)
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				t.Fatal(err)
			}
		}(file)

		mockNoteRepo := new(mocks.MockINoteRepository)
		mockNoteRepo.On("GetNoteByID", 1).Return(retNote, bl.CreateError(bl.AllIsOk, nil, ""))
		mockNoteRepo.On("UpdateNoteContentRawData", mock.Anything, retNote).Return(retErr)

		nsSrv := bl.NoteService{}
		myErr := nsSrv.UpdateNoteContent(1, requester, filePath, mockNoteRepo)

		assert.Nil(t, err)
		assert.NotNil(t, myErr)
		assert.Equal(t, myErr.ErrNum, bl.AllIsOk)

		mockNoteRepo.AssertExpectations(t)
	})

	t.Run("ErrorGetNoteByID", func(t *testing.T) {
		retErr := bl.CreateError(bl.ErrGetNoteByID, bl.ErrGetNoteByIDError(), "GetNoteByID")
		requester := &bl.User{Id: 1, CommandID: 1, Role: bl.Admin}
		retNote := &bl.Note{}
		filePath := "./Tests files/test.txt"

		mockNoteRepo := new(mocks.MockINoteRepository)
		mockNoteRepo.On("GetNoteByID", 1).Return(retNote, retErr)

		nsSrv := bl.NoteService{}
		err := nsSrv.UpdateNoteContent(1, requester, filePath, mockNoteRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrGetNoteByID, err.ErrNum)

		mockNoteRepo.AssertExpectations(t)
	})

	t.Run("ErrorAccessDenied", func(t *testing.T) {
		retOk := bl.CreateError(bl.AllIsOk, nil, "")
		retNote := &bl.Note{Id: 1, OwnerID: 2, ContentType: bl.TextCont}
		requester := &bl.User{Id: 1, CommandID: 1, Role: bl.Author}

		filePath := "./Tests files/test.txt"
		file, err := os.Open(filePath)
		if err != nil {
			t.Fatal(err)
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				t.Fatal(err)
			}
		}(file)

		mockNoteRepo := new(mocks.MockINoteRepository)
		mockNoteRepo.On("GetNoteByID", 1).Return(retNote, retOk)

		nsSrv := bl.NoteService{}
		myErr := nsSrv.UpdateNoteContent(1, requester, filePath, mockNoteRepo)

		assert.NotNil(t, myErr)
		assert.Equal(t, bl.ErrAccessDenied, myErr.ErrNum)

		mockNoteRepo.AssertExpectations(t)
	})

	t.Run("ErrorOpenFile", func(t *testing.T) {
		retOk := bl.CreateError(bl.AllIsOk, nil, "")
		retNote := &bl.Note{Id: 1, OwnerID: 1, ContentType: bl.TextCont}
		requester := &bl.User{Id: 1, CommandID: 1, Role: bl.Admin}

		filePath := "test.txt"

		mockNoteRepo := new(mocks.MockINoteRepository)
		mockNoteRepo.On("GetNoteByID", 1).Return(retNote, retOk)

		nsSrv := bl.NoteService{}
		err := nsSrv.UpdateNoteContent(1, requester, filePath, mockNoteRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrNoFile, err.ErrNum)

		mockNoteRepo.AssertExpectations(t)
	})
}

func TestUpdateNoteInfo(t *testing.T) {
	t.Run("SuccessUpdateNoteInfo", func(t *testing.T) {
		retErr := bl.CreateError(bl.AllIsOk, nil, "")
		retNote := &bl.Note{Id: 1, OwnerID: 1}
		requester := &bl.User{Id: 1, CommandID: 1, Role: bl.Admin}

		mockNoteRepo := new(mocks.MockINoteRepository)
		mockNoteRepo.On("UpdateNoteInfo", retNote).Return(retErr)

		nsSrv := bl.NoteService{}
		err := nsSrv.UpdateNoteInfo(requester, retNote, mockNoteRepo)

		assert.NotNil(t, err)
		assert.Equal(t, err.ErrNum, bl.AllIsOk)

		mockNoteRepo.AssertExpectations(t)
	})

	t.Run("ErrorAccessDenied", func(t *testing.T) {
		retNote := &bl.Note{Id: 1, OwnerID: 2}
		requester := &bl.User{Id: 1, CommandID: 1, Role: bl.Author}

		mockNoteRepo := new(mocks.MockINoteRepository)

		nsSrv := bl.NoteService{}
		err := nsSrv.UpdateNoteInfo(requester, retNote, mockNoteRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrAccessDenied, err.ErrNum)
	})
}

func TestAddNoteToCollection(t *testing.T) {
	t.Run("SuccessAddNoteToCollection", func(t *testing.T) {
		retErr := bl.CreateError(bl.AllIsOk, nil, "")

		mockNoteRepo := new(mocks.MockINoteRepository)
		mockNoteRepo.On("AddNoteToCollection", 1, 1).Return(retErr)

		nsSrv := bl.NoteService{}
		err := nsSrv.AddNoteToCollection(1, 1, mockNoteRepo)

		assert.NotNil(t, err)
		assert.Equal(t, err.ErrNum, bl.AllIsOk)

		mockNoteRepo.AssertExpectations(t)
	})
}
