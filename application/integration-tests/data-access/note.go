package interg_da

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
)

type NoteRepoIntegrationTestSuite struct {
	suite.Suite
	db       *sql.DB
	noteRepo da.NoteRepository

	idNoteToGet int
	idNoteToUpd int
	idNoteToDel int

	idCollNoteAdd int
	idNoteCollAdd int
	idCollNoteDel int
	idNoteCollDel int
}

func (s *NoteRepoIntegrationTestSuite) BeforeEach(t provider.T) {
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

	s.noteRepo = da.NewNoteRepository(&dbconfig, &logger)

	s.idNoteToGet = 1
	s.idNoteToUpd = 5
	s.idNoteToDel = 10

	s.idCollNoteAdd = 6
	s.idNoteCollAdd = 6
	s.idCollNoteDel = 7
	s.idNoteCollDel = 7
}

func (s *NoteRepoIntegrationTestSuite) AfterEach(t provider.T) {
	_ = s.db.Close()
}

func (s *NoteRepoIntegrationTestSuite) TestGetNoteByIDSuccess(t provider.T) {
	t.Title("GetNoteByID: Success")
	t.Tags("NoteRepoIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		note, err := s.noteRepo.GetNoteByID(s.idNoteToGet)

		sCtx.Assert().NotNil(note)
		sCtx.Assert().NotNil(err)
		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *NoteRepoIntegrationTestSuite) TestGetNoteByIDNotFound(t provider.T) {
	t.Title("GetNoteByID: Note Not Found")
	t.Tags("NoteRepoIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Note not found", func(sCtx provider.StepCtx) {
		note, err := s.noteRepo.GetNoteByID(999)

		sCtx.Assert().Nil(note)
		sCtx.Assert().NotNil(err)
		sCtx.Assert().Equal(bl.NoSuchNote, err.ErrNum)
	})
}

func (s *NoteRepoIntegrationTestSuite) TestGetNoteByNameSuccess(t provider.T) {
	t.Title("GetNoteByName: Success")
	t.Tags("NoteRepoIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		note, err := s.noteRepo.GetNoteByName("Note1")

		sCtx.Assert().NotNil(note)
		sCtx.Assert().NotNil(err)
		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *NoteRepoIntegrationTestSuite) TestGetNoteByNameDNotFound(t provider.T) {
	t.Title("GetNoteByName: Note Not Found")
	t.Tags("NoteRepoIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Note not found", func(sCtx provider.StepCtx) {
		note, err := s.noteRepo.GetNoteByName("akldfjkladfjsakljdf")

		sCtx.Assert().Nil(note)
		sCtx.Assert().NotNil(err)
		sCtx.Assert().Equal(bl.NoSuchNote, err.ErrNum)
	})
}

func (s *NoteRepoIntegrationTestSuite) TestGetAllNotesSuccess(t provider.T) {
	t.Title("TestGetAllNotes: Success")
	t.Tags("NoteServiceIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		note, err := s.noteRepo.GetAllNotes()

		sCtx.Assert().NotNil(note)
		sCtx.Assert().NotNil(err)
		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *NoteRepoIntegrationTestSuite) TestGetAllPublicNotesSuccess(t provider.T) {
	t.Title("TestGetAllPublicNotes: Note Not Found")
	t.Tags("NoteServiceIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Access denied", func(sCtx provider.StepCtx) {
		note, err := s.noteRepo.GetAllPublicNotes()

		sCtx.Assert().NotNil(note)
		sCtx.Assert().NotNil(err)
		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *NoteRepoIntegrationTestSuite) TestAddNoteSuccess(t provider.T) {
	t.Title("AddNote: Success")
	t.Tags("NoteRepoIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		note := data_builders.NewNoteBuilder().WithName("Not unique note name but for da-test-hahaha").Build()

		_, err := s.noteRepo.AddNote(note)

		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *NoteRepoIntegrationTestSuite) TestDeleteNoteSuccess(t provider.T) {
	t.Title("DeleteNote: Success")
	t.Tags("NoteRepoIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		err := s.noteRepo.DeleteNote(s.idNoteToDel)

		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *NoteRepoIntegrationTestSuite) TestDeleteNoteNotFound(t provider.T) {
	t.Title("DeleteNote: Note Not Found")
	t.Tags("NoteRepoIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Note not found", func(sCtx provider.StepCtx) {
		err := s.noteRepo.DeleteNote(999)

		sCtx.Assert().NotNil(err)
		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *NoteRepoIntegrationTestSuite) TestUpdateNoteInfoSuccess(t provider.T) {
	t.Title("UpdateNoteInfo: Success")
	t.Tags("NoteRepoIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		note := data_builders.NewNoteBuilder().WithId(s.idNoteToUpd).Build()

		err := s.noteRepo.UpdateNoteInfo(note)

		sCtx.Assert().NotNil(err)
		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *NoteRepoIntegrationTestSuite) TestUpdateNoteContentSuccess(t provider.T) {
	t.Title("UpdateNote: Unauthorized")
	t.Tags("NoteRepoIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Unauthorized", func(sCtx provider.StepCtx) {
		note := data_builders.NewNoteBuilder().WithId(s.idNoteToUpd).Build()

		err := s.noteRepo.UpdateNoteContent(note)

		sCtx.Assert().NotNil(err)
		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *NoteRepoIntegrationTestSuite) TestAddNoteToCollectionSuccess(t provider.T) {
	t.Title("AddNoteToCollection: Success")
	t.Tags("NoteRepoIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		err := s.noteRepo.AddNoteToCollection(s.idNoteCollAdd, s.idCollNoteAdd)

		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *NoteRepoIntegrationTestSuite) TestDeleteNoteFromCollectionSuccess(t provider.T) {
	t.Title("DeleteNoteFromCollection: Success")
	t.Tags("NoteRepoIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		err := s.noteRepo.DeleteNoteFromCollection(s.idNoteCollDel, s.idCollNoteDel)

		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}
