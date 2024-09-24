package UnitTests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/business_logic"
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/business_logic/UnitTests/mocks"
)

func TestGetSection(t *testing.T) {
	t.Run("SuccessGetSectionByID", func(t *testing.T) {
		retErr := bl.CreateError(bl.AllIsOk, nil, "")
		retSection := &bl.Section{Id: 1, CommandID: 1}
		reqUser := &bl.User{Role: bl.Admin}

		mockSectionRepo := new(mocks.MockISectionRepository)
		mockSectionRepo.On("GetSectionByID", 1).Return(retSection, retErr)

		ssSrv := bl.SectionService{}
		section, err := ssSrv.GetSection(1, "", reqUser, bl.SearchByID, mockSectionRepo)

		assert.NotNil(t, err)
		assert.Equal(t, err.ErrNum, bl.AllIsOk)
		assert.Equal(t, 1, section.Id)

		mockSectionRepo.AssertExpectations(t)
	})

	t.Run("SuccessGetSectionByTeamName", func(t *testing.T) {
		retErr := bl.CreateError(bl.AllIsOk, nil, "")
		retSection := &bl.Section{Id: 1, CommandID: 1}
		reqUser := &bl.User{Role: bl.Admin}

		mockSectionRepo := new(mocks.MockISectionRepository)
		mockSectionRepo.On("GetSectionByTeamName", "Team Name").Return(retSection, retErr)

		ssSrv := bl.SectionService{}
		section, err := ssSrv.GetSection(0, "Team Name", reqUser, bl.SearchByString, mockSectionRepo)

		assert.NotNil(t, err)
		assert.Equal(t, err.ErrNum, bl.AllIsOk)
		assert.Equal(t, 1, section.Id)

		mockSectionRepo.AssertExpectations(t)
	})

	t.Run("ErrorSearchParameter", func(t *testing.T) {
		reqUser := &bl.User{Role: bl.Admin}

		ssSrv := bl.SectionService{}
		_, err := ssSrv.GetSection(0, "", reqUser, -1, nil)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrSearchParameter, err.ErrNum)
	})

	t.Run("ErrorAccessDenied", func(t *testing.T) {
		reqUser := &bl.User{Role: bl.Reader}

		ssSrv := bl.SectionService{}
		_, err := ssSrv.GetSection(0, "", reqUser, bl.SearchByID, nil)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrAccessDenied, err.ErrNum)
	})
}

func TestGetAllSections(t *testing.T) {
	t.Run("SuccessGetAllSections", func(t *testing.T) {
		retErr := bl.CreateError(bl.AllIsOk, nil, "")
		retSections := []*bl.Section{
			{Id: 1, CommandID: 1},
			{Id: 2, CommandID: 2},
		}
		reqUser := &bl.User{Role: bl.Admin}

		mockSectionRepo := new(mocks.MockISectionRepository)
		mockSectionRepo.On("GetAllSections").Return(retSections, retErr)

		ssSrv := bl.SectionService{}
		sections, err := ssSrv.GetAllSections(reqUser, mockSectionRepo)

		assert.NotNil(t, err)
		assert.Equal(t, err.ErrNum, bl.AllIsOk)
		assert.Len(t, sections, 2)
		assert.Equal(t, 1, sections[0].Id)
		assert.Equal(t, 2, sections[1].Id)

		mockSectionRepo.AssertExpectations(t)
	})

	t.Run("ErrorAccessDenied", func(t *testing.T) {
		reqUser := &bl.User{Role: bl.Reader}

		ssSrv := bl.SectionService{}
		sections, err := ssSrv.GetAllSections(reqUser, nil)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrAccessDenied, err.ErrNum)

		assert.Nil(t, sections)
	})
}

func TestGetAllNotesInSection(t *testing.T) {
	t.Run("SuccessGetAllNotesInSection", func(t *testing.T) {
		retErr := bl.CreateError(bl.AllIsOk, nil, "")
		retSection := &bl.Section{Id: 1, CommandID: 1}
		retNotes := []*bl.Note{
			{Id: 1, OwnerID: 1},
			{Id: 2, OwnerID: 2},
		}
		reqUser := &bl.User{Role: bl.Admin, CommandID: 1}

		mockSectionRepo := new(mocks.MockISectionRepository)
		mockSectionRepo.On("GetSectionByID", 1).Return(retSection, retErr)
		mockSectionRepo.On("GetAllNotesInSection", retSection).Return(retNotes, retErr)

		ssSrv := bl.SectionService{}
		notes, err := ssSrv.GetAllNotesInSection(1, reqUser, mockSectionRepo)

		assert.NotNil(t, err)
		assert.Equal(t, err.ErrNum, bl.AllIsOk)
		assert.Len(t, notes, 2)
		assert.Equal(t, 1, notes[0].Id)
		assert.Equal(t, 2, notes[1].Id)

		mockSectionRepo.AssertExpectations(t)
	})

	t.Run("ErrorGetSectionByID", func(t *testing.T) {
		retErr := bl.CreateError(bl.ErrGetSectionByID, bl.ErrGetSectionByIDError(), "GetSectionByID")
		reqUser := &bl.User{Role: bl.Admin, CommandID: 1}
		reqSection := &bl.Section{}

		mockSectionRepo := new(mocks.MockISectionRepository)
		mockSectionRepo.On("GetSectionByID", 1).Return(reqSection, retErr)

		ssSrv := bl.SectionService{}
		notes, err := ssSrv.GetAllNotesInSection(1, reqUser, mockSectionRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrGetSectionByID, err.ErrNum)

		assert.Nil(t, notes)

		mockSectionRepo.AssertExpectations(t)
	})

	t.Run("ErrorAccessDenied", func(t *testing.T) {
		retErr := bl.CreateError(bl.AllIsOk, nil, "")
		retSection := &bl.Section{Id: 1, CommandID: 2}
		reqUser := &bl.User{Role: bl.Author, CommandID: 1}

		mockSectionRepo := new(mocks.MockISectionRepository)
		mockSectionRepo.On("GetSectionByID", 1).Return(retSection, retErr)

		ssSrv := bl.SectionService{}
		_, err := ssSrv.GetAllNotesInSection(1, reqUser, mockSectionRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrAccessDenied, err.ErrNum)

		mockSectionRepo.AssertExpectations(t)
	})
}

func TestAddNoteToSection(t *testing.T) {
	t.Run("SuccessAddNoteToSection", func(t *testing.T) {
		retErr := bl.CreateError(bl.AllIsOk, nil, "")
		targetSection := &bl.Section{Id: 1, CommandID: 1}
		srcNote := &bl.Note{Id: 1, OwnerID: 1}
		reqUser := &bl.User{Role: bl.Admin, CommandID: 1}

		mockSectionRepo := new(mocks.MockISectionRepository)
		mockSectionRepo.On("AddNoteToSection", srcNote, targetSection).Return(retErr)

		ssSrv := bl.SectionService{}
		err := ssSrv.AddNoteToSection(targetSection, srcNote, reqUser, mockSectionRepo)

		assert.NotNil(t, err)
		assert.Equal(t, err.ErrNum, bl.AllIsOk)

		mockSectionRepo.AssertExpectations(t)
	})

	t.Run("ErrorAccessDenied", func(t *testing.T) {
		targetSection := &bl.Section{Id: 1, CommandID: 2}
		srcNote := &bl.Note{Id: 1, OwnerID: 1}
		reqUser := &bl.User{Role: bl.Author, CommandID: 1}

		mockSectionRepo := new(mocks.MockISectionRepository)

		ssSrv := bl.SectionService{}
		err := ssSrv.AddNoteToSection(targetSection, srcNote, reqUser, mockSectionRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrAccessDenied, err.ErrNum)

		mockSectionRepo.AssertExpectations(t)
	})
}

func TestAddSection(t *testing.T) {
	t.Run("SuccessAddSection", func(t *testing.T) {
		retErr := bl.CreateError(bl.AllIsOk, nil, "")
		targetSection := &bl.Section{Id: 1, CommandID: 1}
		reqUser := &bl.User{Role: bl.Admin}

		mockSectionRepo := new(mocks.MockISectionRepository)
		mockSectionRepo.On("AddSection", targetSection).Return(retErr)

		ssSrv := bl.SectionService{}
		err := ssSrv.AddSection(targetSection, reqUser, mockSectionRepo)

		assert.NotNil(t, err)
		assert.Equal(t, err.ErrNum, bl.AllIsOk)

		mockSectionRepo.AssertExpectations(t)
	})

	t.Run("ErrorAccessDenied", func(t *testing.T) {
		targetSection := &bl.Section{Id: 1, CommandID: 1}
		reqUser := &bl.User{Role: bl.Reader}

		mockSectionRepo := new(mocks.MockISectionRepository)

		ssSrv := bl.SectionService{}
		err := ssSrv.AddSection(targetSection, reqUser, mockSectionRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrAccessDenied, err.ErrNum)

		mockSectionRepo.AssertExpectations(t)
	})
}

func TestDeleteSection(t *testing.T) {
	t.Run("SuccessDeleteSection", func(t *testing.T) {
		retErr := bl.CreateError(bl.AllIsOk, nil, "")
		reqUser := &bl.User{Role: bl.Admin}

		mockSectionRepo := new(mocks.MockISectionRepository)
		mockSectionRepo.On("DeleteSection", 1).Return(retErr)

		ssSrv := bl.SectionService{}
		err := ssSrv.DeleteSection(1, reqUser, mockSectionRepo)

		assert.NotNil(t, err)
		assert.Equal(t, err.ErrNum, bl.AllIsOk)

		mockSectionRepo.AssertExpectations(t)
	})

	t.Run("ErrorAccessDenied", func(t *testing.T) {
		reqUser := &bl.User{Role: bl.Reader}

		mockSectionRepo := new(mocks.MockISectionRepository)

		ssSrv := bl.SectionService{}
		err := ssSrv.DeleteSection(1, reqUser, mockSectionRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrAccessDenied, err.ErrNum)

		mockSectionRepo.AssertExpectations(t)
	})
}

func TestUpdateSection(t *testing.T) {
	t.Run("SuccessUpdateSection", func(t *testing.T) {
		retErr := bl.CreateError(bl.AllIsOk, nil, "")
		targetSection := &bl.Section{Id: 1, CommandID: 1}
		reqUser := &bl.User{Role: bl.Admin}

		mockSectionRepo := new(mocks.MockISectionRepository)
		mockSectionRepo.On("UpdateSection", targetSection).Return(retErr)

		ssSrv := bl.SectionService{}
		err := ssSrv.UpdateSection(targetSection, reqUser, mockSectionRepo)

		assert.NotNil(t, err)
		assert.Equal(t, err.ErrNum, bl.AllIsOk)

		mockSectionRepo.AssertExpectations(t)
	})

	t.Run("ErrorAccessDenied", func(t *testing.T) {
		targetSection := &bl.Section{Id: 1, CommandID: 1}
		reqUser := &bl.User{Role: bl.Reader}

		mockSectionRepo := new(mocks.MockISectionRepository)

		ssSrv := bl.SectionService{}
		err := ssSrv.UpdateSection(targetSection, reqUser, mockSectionRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrAccessDenied, err.ErrNum)

		mockSectionRepo.AssertExpectations(t)
	})
}
