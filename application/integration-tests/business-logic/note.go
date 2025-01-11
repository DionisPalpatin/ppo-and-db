package interg_bl

import (
	integr_ut "github.com/DionisPalpatin/ppo-and-db/application/integration-tests/utils"
	bl "github.com/DionisPalpatin/ppo-and-db/application/internal/business-logic"
	"github.com/DionisPalpatin/ppo-and-db/application/internal/config"
	da "github.com/DionisPalpatin/ppo-and-db/application/internal/data-access"
	mylogger "github.com/DionisPalpatin/ppo-and-db/application/internal/logger"
	data_builders "github.com/DionisPalpatin/ppo-and-db/application/internal/models/builders"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"

	"database/sql"
	"fmt"
	"os"
)

type NoteServiceIntegrationTest struct {
	suite.Suite
	db          *sql.DB
	noteService bl.INoteService

	idNoteToGet int
	idNoteToUpd int
	idNoteToDel int

	idCollNoteAdd int
	idNoteCollAdd int
	idCollNoteDel int
	idNoteCollDel int
}

func (s *NoteServiceIntegrationTest) BeforeEach(t provider.T) {
	var err error
	pgInfo := integr_ut.PostgresInfo{
		Host:     "localhost",
		User:     "postgres",
		Password: "password",
		Port:     15423,
		DBName:   "NotebookAppIntTests",
	}

	s.db, err = integr_ut.InitDB(pgInfo)
	if err != nil {
		fmt.Println(err)
		return
	}

	if s.db == nil {
		return
	}

	text, err := os.ReadFile("./init.sql")
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = s.db.Exec(string(text))
	if err != nil {
		fmt.Println(err)
		return
	}

	dbconfig := config.DBConfigs{
		Host:       pgInfo.Host,
		User:       pgInfo.User,
		Password:   pgInfo.Password,
		Port:       5432,
		Name:       pgInfo.DBName,
		DB:         s.db,
		SchemaName: "interg_tests",
	}

	logger := mylogger.MyLogger{}
	logger.InitLogger("./logs/integration-test-note.log", "debug")

	noteRepo := da.NewNoteRepository(&dbconfig, &logger)
	s.noteService = bl.NewNoteService(&noteRepo, &logger)

	s.idNoteToGet = 1
	s.idNoteToUpd = 5
	s.idNoteToDel = 11

	s.idCollNoteAdd = 4
	s.idNoteCollAdd = 4
	s.idCollNoteDel = 5
	s.idNoteCollDel = 5
}

func (s *NoteServiceIntegrationTest) AfterEach(t provider.T) {
	_ = s.db.Close()
}

func (s *NoteServiceIntegrationTest) TestGetNoteSuccess(t provider.T) {
	t.Title("GetNote: Success")
	t.Tags("NoteServiceIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		requester := data_builders.NewUserBuilder().WithRole(bl.Admin).Build()
		note, err := s.noteService.GetNote(s.idNoteToGet, "", bl.SearchByID, requester)

		sCtx.Assert().NotNil(note)
		sCtx.Assert().NotNil(err)
		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *NoteServiceIntegrationTest) TestGetNoteNotFound(t provider.T) {
	t.Title("GetNote: Note Not Found")
	t.Tags("NoteServiceIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Note not found", func(sCtx provider.StepCtx) {
		requester := data_builders.NewUserBuilder().WithRole(bl.Admin).Build()
		note, err := s.noteService.GetNote(999, "", bl.SearchByID, requester)

		sCtx.Assert().Nil(note)
		sCtx.Assert().NotNil(err)
		sCtx.Assert().Equal(bl.NoSuchNote, err.ErrNum)
	})
}

func (s *NoteServiceIntegrationTest) TestGetAllNotesSuccess(t provider.T) {
	t.Title("TestGetAllNotes: Success")
	t.Tags("NoteServiceIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		requester := data_builders.NewUserBuilder().WithRole(bl.Admin).Build()

		notes, err := s.noteService.GetAllNotes(false, requester)

		sCtx.Assert().NotNil(notes)
		sCtx.Assert().NotNil(err)
		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *NoteServiceIntegrationTest) TestGetAllNotesAccessDenied(t provider.T) {
	t.Title("TestGetAllNotes: Note Not Found")
	t.Tags("NoteServiceIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Access denied", func(sCtx provider.StepCtx) {
		requester := data_builders.NewUserBuilder().WithRole(bl.Reader).Build()
		note, err := s.noteService.GetAllNotes(false, requester)

		sCtx.Assert().Nil(note)
		sCtx.Assert().NotNil(err)
		sCtx.Assert().Equal(bl.ErrAccessDenied, err.ErrNum)
	})
}

func (s *NoteServiceIntegrationTest) TestAddNoteSuccess(t provider.T) {
	t.Title("AddNote: Success")
	t.Tags("NoteServiceIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		note := data_builders.NewNoteBuilder().WithName("Not unique note name").Build()
		requester := data_builders.NewUserBuilder().WithRole(bl.Admin).Build()

		_, err := s.noteService.AddNote(note, requester)

		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *NoteServiceIntegrationTest) TestAddNoteUnauthorized(t provider.T) {
	t.Title("AddNote: Unauthorized")
	t.Tags("NoteServiceIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Unauthorized", func(sCtx provider.StepCtx) {
		note := data_builders.NewNoteBuilder().Build()
		requester := data_builders.NewUserBuilder().WithRole(bl.Reader).Build()

		id, err := s.noteService.AddNote(note, requester)

		sCtx.Assert().Zero(id)
		sCtx.Assert().NotNil(err)
		sCtx.Assert().Equal(bl.ErrAccessDenied, err.ErrNum)
	})
}

func (s *NoteServiceIntegrationTest) TestDeleteNoteSuccess(t provider.T) {
	t.Title("DeleteNote: Success")
	t.Tags("NoteServiceIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		requester := data_builders.NewUserBuilder().WithRole(bl.Admin).Build()

		err := s.noteService.DeleteNote(s.idNoteToDel, requester)

		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *NoteServiceIntegrationTest) TestDeleteNoteNotFound(t provider.T) {
	t.Title("DeleteNote: Note Not Found")
	t.Tags("NoteServiceIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Note not found", func(sCtx provider.StepCtx) {
		requester := data_builders.NewUserBuilder().WithRole(bl.Admin).Build()

		err := s.noteService.DeleteNote(999, requester)

		sCtx.Assert().NotNil(err)
		sCtx.Assert().Equal(bl.NoSuchNote, err.ErrNum)
	})
}

func (s *NoteServiceIntegrationTest) TestUpdateNoteSuccess(t provider.T) {
	t.Title("UpdateNote: Success")
	t.Tags("NoteServiceIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		requester := data_builders.NewUserBuilder().WithRole(bl.Admin).Build()
		note := data_builders.NewNoteBuilder().WithId(s.idNoteToUpd).Build()

		err := s.noteService.UpdateNote(note, requester, "./tests-files/text.txt", "./tests-files/image.jpg", "./tests-files/raw.file")

		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *NoteServiceIntegrationTest) TestUpdateNoteUnauthorized(t provider.T) {
	t.Title("UpdateNote: Unauthorized")
	t.Tags("NoteServiceIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Unauthorized", func(sCtx provider.StepCtx) {
		note := data_builders.NewNoteBuilder().Build()
		requester := data_builders.NewUserBuilder().WithRole(bl.Reader).Build()

		err := s.noteService.UpdateNote(note, requester, "./tests-files/text.txt", "./tests-files/image.png", "./tests-files/raw.file")

		sCtx.Assert().NotNil(err)
		sCtx.Assert().Equal(bl.ErrAccessDenied, err.ErrNum)
	})
}

func (s *NoteServiceIntegrationTest) TestAddNoteToCollectionSuccess(t provider.T) {
	t.Title("AddNoteToCollection: Success")
	t.Tags("NoteServiceIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		err := s.noteService.AddNoteToCollection(s.idNoteCollAdd, s.idCollNoteAdd)

		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *NoteServiceIntegrationTest) TestDeleteNoteFromCollectionSuccess(t provider.T) {
	t.Title("DeleteNoteFromCollection: Success")
	t.Tags("NoteServiceIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		err := s.noteService.DeleteNoteFromCollection(s.idNoteCollDel, s.idCollNoteDel)

		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}
