package unit_tests

import (
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"

	bl "github.com/DionisPalpatin/ppo-and-db/application/internal/business-logic"
	mocks "github.com/DionisPalpatin/ppo-and-db/application/internal/database/mocks/v2"
	mylogger "github.com/DionisPalpatin/ppo-and-db/application/internal/logger"

	myerror_builder "github.com/DionisPalpatin/ppo-and-db/application/internal/business-logic/builder"

	"github.com/DionisPalpatin/ppo-and-db/application/internal/models"
	"github.com/DionisPalpatin/ppo-and-db/application/internal/models/builders"
)

type NoteServiceTestSuite struct {
	suite.Suite
}

type extDependensesNote struct {
	noteRepo *mocks.MockNoteRepository
	logger   *mylogger.MyLogger
}

func initExtDependenciesNote(t provider.T) *extDependensesNote {
	mockRepo := mocks.NewMockNoteRepository(t)

	f := &extDependensesNote{
		noteRepo: mockRepo,
		logger:   &mylogger.MyLogger{},
	}

	f.logger.InitLogger("./logs/unit-test-note.log", "debug")

	return f
}

func (s *NoteServiceTestSuite) TestNoteService_GetNote_OK(t provider.T) {
	t.Title("Get note: OK")
	t.Tags("NoteService")
	t.Parallel()

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		df := initExtDependenciesNote(t)
		inNote := builders.NewNoteBuilder().Build()
		outNote := builders.NewNoteBuilder().WithId(inNote.Id)
		reqUser := builders.NewUserBuilder().WithRole(bl.Admin)
		myerr := myerror_builder.NewMyErrorBuilder().WithErrNum(bl.Ok).Build()

		df.noteRepo.EXPECT().GetNoteByID(inNote.Id).Return(outNote.Note, myerr).Once()
		noteService := bl.NewNoteService(df.noteRepo, df.logger)

		note, err := noteService.GetNote(inNote.Id, "", bl.SearchByID, reqUser.User)

		sCtx.Assert().Equal(err.ErrNum, bl.Ok)
		sCtx.Assert().Equal(note.Id, outNote.Note.Id)
	})
}

func (s *NoteServiceTestSuite) TestNoteService_GetNote_AccessDenied(t provider.T) {
	t.Title("Get note: Access Denied")
	t.Tags("NoteService")
	t.Parallel()

	t.WithNewStep("Access Denied", func(sCtx provider.StepCtx) {
		df := initExtDependenciesNote(t)
		inNote := builders.NewNoteBuilder().Build()
		reqUser := builders.NewUserBuilder().WithRole(bl.Reader)
		myerr := myerror_builder.NewMyErrorBuilder().WithErrNum(bl.ErrAccessDenied).Build()

		df.noteRepo.EXPECT().GetNoteByID(inNote.Id).Return(nil, myerr).Once()
		noteService := bl.NewNoteService(df.noteRepo, df.logger)

		note, err := noteService.GetNote(inNote.Id, "", bl.SearchByID, reqUser.User)

		sCtx.Assert().Nil(note)
		sCtx.Assert().Equal(err.ErrNum, bl.ErrAccessDenied)
	})
}

func (s *NoteServiceTestSuite) TestNoteService_GetAllNotes_OK(t provider.T) {
	t.Title("Get all notes: OK")
	t.Tags("NoteService")
	t.Parallel()

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		df := initExtDependenciesNote(t)
		notes := []*models.Note{
			builders.NewNoteBuilder().Build(),
			builders.NewNoteBuilder().Build(),
		}
		reqUser := builders.NewUserBuilder().WithRole(bl.Admin)
		myerr := myerror_builder.NewMyErrorBuilder().WithErrNum(bl.Ok).Build()

		df.noteRepo.EXPECT().GetAllNotes().Return(notes, myerr).Once()
		noteService := bl.NewNoteService(df.noteRepo, df.logger)

		resNotes, err := noteService.GetAllNotes(false, reqUser.User)

		sCtx.Assert().Equal(err.ErrNum, bl.Ok)
		sCtx.Assert().Equal(len(resNotes), len(notes))
	})
}

func (s *NoteServiceTestSuite) TestNoteService_GetAllNotes_AccessDenied(t provider.T) {
	t.Title("Get all notes: Access Denied")
	t.Tags("NoteService")
	t.Parallel()

	t.WithNewStep("Access Denied", func(sCtx provider.StepCtx) {
		df := initExtDependenciesNote(t)
		reqUser := builders.NewUserBuilder().WithRole(bl.Reader)

		noteService := bl.NewNoteService(df.noteRepo, df.logger)

		resNotes, err := noteService.GetAllNotes(false, reqUser.User)

		sCtx.Assert().Nil(resNotes)
		sCtx.Assert().Equal(err.ErrNum, bl.ErrAccessDenied)
	})
}

func (s *NoteServiceTestSuite) TestNoteService_AddNote_OK(t provider.T) {
	t.Title("Add note: OK")
	t.Tags("NoteService")
	t.Parallel()

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		df := initExtDependenciesNote(t)
		inNote := builders.NewNoteBuilder().WithId(0)
		reqUser := builders.NewUserBuilder().WithRole(bl.Admin)
		myerr := myerror_builder.NewMyErrorBuilder().WithErrNum(bl.Ok).Build()

		df.noteRepo.EXPECT().AddNote(inNote.Note).Return(myerr).Once()
		noteService := bl.NewNoteService(df.noteRepo, df.logger)

		id, err := noteService.AddNote(inNote.Note, reqUser.User)

		sCtx.Assert().Equal(err.ErrNum, bl.Ok)
		sCtx.Assert().Equal(id, inNote.Note.Id)
	})
}

func (s *NoteServiceTestSuite) TestNoteService_AddNote_AccessDenied(t provider.T) {
	t.Title("Add note: Access Denied")
	t.Tags("NoteService")
	t.Parallel()

	t.WithNewStep("Access Denied", func(sCtx provider.StepCtx) {
		df := initExtDependenciesNote(t)
		inNote := builders.NewNoteBuilder().Build()
		reqUser := builders.NewUserBuilder().WithRole(bl.Reader)

		noteService := bl.NewNoteService(df.noteRepo, df.logger)

		id, err := noteService.AddNote(inNote, reqUser.User)

		sCtx.Assert().Equal(id, 0)
		sCtx.Assert().Equal(err.ErrNum, bl.ErrAccessDenied)
	})
}

func (s *NoteServiceTestSuite) TestNoteService_DeleteNote_OK(t provider.T) {
	t.Title("Delete note: OK")
	t.Tags("NoteService")
	t.Parallel()

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		df := initExtDependenciesNote(t)
		reqUser := builders.NewUserBuilder().WithRole(bl.Admin)
		myerr := myerror_builder.NewMyErrorBuilder().WithErrNum(bl.Ok).Build()
		inNote := builders.NewNoteBuilder().Build()

		df.noteRepo.EXPECT().GetNoteByID(inNote.Id).Return(inNote, myerr).Once()

		df.noteRepo.EXPECT().DeleteNote(1).Return(myerr).Once()
		noteService := bl.NewNoteService(df.noteRepo, df.logger)

		err := noteService.DeleteNote(1, reqUser.User)

		sCtx.Assert().Equal(err.ErrNum, bl.Ok)
	})
}

func (s *NoteServiceTestSuite) TestNoteService_DeleteNote_AccessDenied(t provider.T) {
	t.Title("Delete note: Access Denied")
	t.Tags("NoteService")
	t.Parallel()

	t.WithNewStep("Access Denied", func(sCtx provider.StepCtx) {
		df := initExtDependenciesNote(t)
		reqUser := builders.NewUserBuilder().WithRole(bl.Reader)
		myok := myerror_builder.NewMyErrorBuilder().WithErrNum(bl.Ok).Build()
		inNote := builders.NewNoteBuilder().Build()

		df.noteRepo.EXPECT().GetNoteByID(inNote.Id).Return(inNote, myok).Once()
		noteService := bl.NewNoteService(df.noteRepo, df.logger)

		err := noteService.DeleteNote(1, reqUser.User)

		sCtx.Assert().Equal(err.ErrNum, bl.ErrAccessDenied)
	})
}

func (s *NoteServiceTestSuite) TestNoteService_UpdateNote_OK(t provider.T) {
	t.Title("Update note: OK")
	t.Tags("NoteService")
	t.Parallel()

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		df := initExtDependenciesNote(t)
		reqUser := builders.NewUserBuilder().WithRole(bl.Admin)
		myerr := myerror_builder.NewMyErrorBuilder().WithErrNum(bl.Ok).Build()
		inNote := builders.NewNoteBuilder().Build()

		df.noteRepo.EXPECT().UpdateNoteContent(inNote).Return(myerr).Once()
		df.noteRepo.EXPECT().UpdateNoteInfo(inNote).Return(myerr).Once()
		noteService := bl.NewNoteService(df.noteRepo, df.logger)

		err := noteService.UpdateNote(inNote, reqUser.User, "./tests-files/text.txt", "./tests-files/image.jpg", "./tests-files/raw.file")

		sCtx.Assert().Equal(err.ErrNum, bl.Ok)
	})
}

func (s *NoteServiceTestSuite) TestNoteService_UpdateNote_AccessDenied(t provider.T) {
	t.Title("Update note: Access Denied")
	t.Tags("NoteService")
	t.Parallel()

	t.WithNewStep("Access Denied", func(sCtx provider.StepCtx) {
		df := initExtDependenciesNote(t)
		reqUser := builders.NewUserBuilder().WithRole(bl.Reader)
		inNote := builders.NewNoteBuilder().Build()

		noteService := bl.NewNoteService(df.noteRepo, df.logger)

		err := noteService.UpdateNote(inNote, reqUser.User, "./tests-files/text.txt", "./tests-files/image.png", "./tests-files/raw.file")

		sCtx.Assert().Equal(err.ErrNum, bl.ErrAccessDenied)
	})
}

func (s *NoteServiceTestSuite) TestNoteService_AddNoteToCollection_OK(t provider.T) {
	t.Title("Add note to collection: OK")
	t.Tags("NoteService")
	t.Parallel()

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		df := initExtDependenciesNote(t)
		noteID := 1
		collID := 42
		myerr := myerror_builder.NewMyErrorBuilder().WithErrNum(bl.Ok).Build()

		df.noteRepo.EXPECT().AddNoteToCollection(collID, noteID).Return(myerr).Once()
		noteService := bl.NewNoteService(df.noteRepo, df.logger)

		err := noteService.AddNoteToCollection(noteID, collID)

		sCtx.Assert().Equal(err.ErrNum, bl.Ok)
	})
}

func (s *NoteServiceTestSuite) TestNoteService_AddNoteToCollection_AccessDenied(t provider.T) {
	t.Title("Add note to collection: Access Denied")
	t.Tags("NoteService")
	t.Parallel()

	t.WithNewStep("Access Denied", func(sCtx provider.StepCtx) {
		df := initExtDependenciesNote(t)
		noteID := -1
		collID := 42
		myerr := myerror_builder.NewMyErrorBuilder().WithErrNum(bl.ErrAccessDenied).Build()

		df.noteRepo.EXPECT().AddNoteToCollection(collID, noteID).Return(myerr).Once()
		noteService := bl.NewNoteService(df.noteRepo, df.logger)

		err := noteService.AddNoteToCollection(noteID, collID)

		sCtx.Assert().Equal(err.ErrNum, bl.ErrAccessDenied)
	})
}

func (s *NoteServiceTestSuite) TestNoteService_DeleteNoteFromCollection_OK(t provider.T) {
	t.Title("Delete note from collection: OK")
	t.Tags("NoteService")
	t.Parallel()

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		df := initExtDependenciesNote(t)
		noteID := 1
		collID := 42
		myerr := myerror_builder.NewMyErrorBuilder().WithErrNum(bl.Ok).Build()

		df.noteRepo.EXPECT().DeleteNoteFromCollection(collID, noteID).Return(myerr).Once()
		noteService := bl.NewNoteService(df.noteRepo, df.logger)

		err := noteService.DeleteNoteFromCollection(noteID, collID)

		sCtx.Assert().Equal(err.ErrNum, bl.Ok)
	})
}

func (s *NoteServiceTestSuite) TestNoteService_DeleteNoteFromCollection_ErrInParameter(t provider.T) {
	t.Title("Delete note from collection: ErrInParameter")
	t.Tags("NoteService")
	t.Parallel()

	t.WithNewStep("ErrInParameter", func(sCtx provider.StepCtx) {
		df := initExtDependenciesNote(t)
		noteID := -1
		collID := 42
		myerr := myerror_builder.NewMyErrorBuilder().WithErrNum(bl.ErrInParameter).Build()

		df.noteRepo.EXPECT().DeleteNoteFromCollection(collID, noteID).Return(myerr).Once()
		noteService := bl.NewNoteService(df.noteRepo, df.logger)

		err := noteService.DeleteNoteFromCollection(noteID, collID)

		sCtx.Assert().Equal(err.ErrNum, bl.ErrInParameter)
	})
}
