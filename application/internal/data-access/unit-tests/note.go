package da_unit

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"

	"github.com/DionisPalpatin/ppo-and-db/application/internal/config"
	mylogger "github.com/DionisPalpatin/ppo-and-db/application/internal/logger"
	"github.com/DionisPalpatin/ppo-and-db/application/internal/models"

	bl "github.com/DionisPalpatin/ppo-and-db/application/internal/business-logic"
	da "github.com/DionisPalpatin/ppo-and-db/application/internal/data-access"

	// "github.com/DionisPalpatin/ppo-and-db/application/internal/models"
	data_builders "github.com/DionisPalpatin/ppo-and-db/application/internal/models/builders"
)

type NoteRepositorySuite struct {
	suite.Suite

	db   *sql.DB
	mock sqlmock.Sqlmock
	repo da.NoteRepository
	ctx  context.Context
}

func (s *NoteRepositorySuite) BeforeEach(t provider.T) {
	var err error

	s.db, s.mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}

	confs := config.DBConfigs{
		DB:         s.db,
		SchemaName: "test",
	}
	logger := mylogger.MyLogger{}
	logger.InitLogger("./logs/unit-tests-note.log", "debug")

	s.repo = da.NewNoteRepository(&confs, &logger)
	s.ctx = context.Background()
}

func (s *NoteRepositorySuite) AfterEach(t provider.T) {
	s.db.Close()
}

func TestNoteRepositorySuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(NoteRepositorySuite))
}

func (s *NoteRepositorySuite) Test_GetNoteByID_Success(t provider.T) {
	t.Title("GetNoteByID: Success")
	t.Tags("NoteRepository")

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		id := 1
		expecContent := data_builders.NewContentBuilder().
			WithText([]byte("Note text content")).
			WithTextExt(".txt").
			WithImg([]byte("image_data")).
			WithImgExt(".png").
			WithRaw([]byte("raw_data")).
			WithRawExt(".bin").
			Build()
		expectedNote := data_builders.NewNoteBuilder().
			WithId(id).
			WithAccess(1).
			WithName("Test Note").
			WithLikes(10).
			WithDislikes(2).
			WithRegistrationDate(ZEROTIME).
			WithOwnerID(100).
			WithSectionID(5).
			WithContent(*expecContent).
			Build()

		// Mock query for note details
		var content_type int = 1
		query := fmt.Sprintf(getNoteByIDQuery, s.repo.DbConfigs.SchemaName)
		s.mock.ExpectQuery(query).
			WithArgs(id).
			WillReturnRows(sqlmock.NewRows([]string{"id", "access", "name", "content_type", "likes", "dislikes", "registration_date", "owner_id", "section_id"}).
				AddRow(expectedNote.Id, expectedNote.Access, expectedNote.Name, content_type, expectedNote.Likes, expectedNote.Dislikes, expectedNote.RegistrationDate, expectedNote.OwnerID, expectedNote.SectionID))

		// Mock query for text content
		query = fmt.Sprintf(getNoteTextContent, s.repo.DbConfigs.SchemaName)
		s.mock.ExpectQuery(query).
			WithArgs(id).
			WillReturnRows(sqlmock.NewRows([]string{"text", "text_ext"}).
				AddRow(expectedNote.Content.Text, expectedNote.Content.TextExt))

		// Mock query for image content
		query = fmt.Sprintf(getNoteImageContent, s.repo.DbConfigs.SchemaName)
		s.mock.ExpectQuery(query).
			WithArgs(id).
			WillReturnRows(sqlmock.NewRows([]string{"img", "img_ext"}).
				AddRow(expectedNote.Content.Img, expectedNote.Content.ImgExt))

		// Mock query for raw data content
		query = fmt.Sprintf(getNoteRawDataContent, s.repo.DbConfigs.SchemaName)
		s.mock.ExpectQuery(query).
			WithArgs(id).
			WillReturnRows(sqlmock.NewRows([]string{"raw", "raw_ext"}).
				AddRow(expectedNote.Content.Raw, expectedNote.Content.RawExt))

		// Call the method
		note, err := s.repo.GetNoteByID(id)

		// Assertions
		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
		sCtx.Assert().NotNil(note)
		sCtx.Assert().Equal(expectedNote, note)
	})
}

func (s *NoteRepositorySuite) Test_GetNoteByID_InvalidID(t provider.T) {
	t.Title("GetNoteByID: Invalid ID")
	t.Tags("NoteRepository")

	t.WithNewStep("Invalid ID", func(sCtx provider.StepCtx) {
		id := -1

		note, err := s.repo.GetNoteByID(id)

		sCtx.Assert().Nil(note)
		sCtx.Assert().NotNil(err)
		sCtx.Assert().Equal(bl.ErrInParameter, err.ErrNum)
	})
}

func (s *NoteRepositorySuite) Test_GetNoteByName_Success(t provider.T) {
	t.Title("GetNoteByID: Success")
	t.Tags("NoteRepository")

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		id := 1
		name := "Test Note"
		expecContent := data_builders.NewContentBuilder().
			WithText([]byte("Note text content")).
			WithTextExt(".txt").
			WithImg([]byte("image_data")).
			WithImgExt(".png").
			WithRaw([]byte("raw_data")).
			WithRawExt(".bin").
			Build()
		expectedNote := data_builders.NewNoteBuilder().
			WithId(id).
			WithAccess(1).
			WithName(name).
			WithLikes(10).
			WithDislikes(2).
			WithRegistrationDate(ZEROTIME).
			WithOwnerID(100).
			WithSectionID(5).
			WithContent(*expecContent).
			Build()

		// Mock query for note details
		var contentType int = 1
		query := fmt.Sprintf(getNoteByNameQuery, s.repo.DbConfigs.SchemaName)
		s.mock.ExpectQuery(query).
			WithArgs(name).
			WillReturnRows(sqlmock.NewRows([]string{"id", "access", "name", "content_type", "likes", "dislikes", "registration_date", "owner_id", "section_id"}).
				AddRow(expectedNote.Id, expectedNote.Access, expectedNote.Name, contentType, expectedNote.Likes, expectedNote.Dislikes, expectedNote.RegistrationDate, expectedNote.OwnerID, expectedNote.SectionID))

		// Mock query for text content
		query = fmt.Sprintf(getNoteTextContent, s.repo.DbConfigs.SchemaName)
		s.mock.ExpectQuery(query).
			WithArgs(id).
			WillReturnRows(sqlmock.NewRows([]string{"text", "text_ext"}).
				AddRow(expectedNote.Content.Text, expectedNote.Content.TextExt))

		// Mock query for image content
		query = fmt.Sprintf(getNoteImageContent, s.repo.DbConfigs.SchemaName)
		s.mock.ExpectQuery(query).
			WithArgs(id).
			WillReturnRows(sqlmock.NewRows([]string{"img", "img_ext"}).
				AddRow(expectedNote.Content.Img, expectedNote.Content.ImgExt))

		// Mock query for raw data content
		query = fmt.Sprintf(getNoteRawDataContent, s.repo.DbConfigs.SchemaName)
		s.mock.ExpectQuery(query).
			WithArgs(id).
			WillReturnRows(sqlmock.NewRows([]string{"raw", "raw_ext"}).
				AddRow(expectedNote.Content.Raw, expectedNote.Content.RawExt))

		// Call the method
		note, err := s.repo.GetNoteByName(name)

		// Assertions
		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
		sCtx.Assert().NotNil(note)
		sCtx.Assert().Equal(expectedNote, note)
	})
}

func (s *NoteRepositorySuite) Test_GetNoteByName_InvalidName(t provider.T) {
	t.Title("GetNoteByID: Invalid ID")
	t.Tags("NoteRepository")

	t.WithNewStep("Invalid Name", func(sCtx provider.StepCtx) {
		name := ""

		note, err := s.repo.GetNoteByName(name)

		sCtx.Assert().Nil(note)
		sCtx.Assert().NotNil(err)
		sCtx.Assert().Equal(bl.ErrInParameter, err.ErrNum)
	})
}

func (s *NoteRepositorySuite) Test_GetAllNotes_Success(t provider.T) {
	t.Title("GetAllNotes: Success")
	t.Tags("NoteRepository")

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		query := fmt.Sprintf(getAllNotesQuery, s.repo.DbConfigs.SchemaName)
		s.mock.ExpectQuery(query).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "access", "name", "content_type", "likes", "dislikes", "registration_date", "owner_id", "section_id",
			}).
				AddRow(1, 1, "Note1", 1, 10, 1, ZEROTIME, 1, 1).
				AddRow(2, 1, "Note 2", 1, 20, 0, ZEROTIME, 2, 2))

		notes, err := s.repo.GetAllNotes()

		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
		sCtx.Assert().Len(notes, 2)
		sCtx.Assert().Equal("Note1", notes[0].Name)
		sCtx.Assert().Equal("Note 2", notes[1].Name)
	})
}

func (s *NoteRepositorySuite) Test_GetAllNotes_Failure(t provider.T) {
	t.Title("GetAllNotes: Failure")
	t.Tags("NoteRepository")

	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		query := fmt.Sprintf(getAllNotesQuery, s.repo.DbConfigs.SchemaName)
		s.mock.ExpectQuery(query).WillReturnError(sql.ErrConnDone)

		notes, err := s.repo.GetAllNotes()

		sCtx.Assert().Nil(notes)
		sCtx.Assert().NotNil(err)
		sCtx.Assert().Equal(bl.DatabaseError, err.ErrNum)
	})
}

func (s *NoteRepositorySuite) Test_DeleteNote_Success(t provider.T) {
	t.Title("DeleteNote: Success")
	t.Tags("NoteRepository")

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		query := fmt.Sprintf(deleteNoteQuery, s.repo.DbConfigs.SchemaName)
		s.mock.ExpectBegin()
		s.mock.ExpectExec(query).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 1))
		s.mock.ExpectCommit()

		err := s.repo.DeleteNote(1)

		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *NoteRepositorySuite) Test_DeleteNote_Failure(t provider.T) {
	t.Title("DeleteNote: Failure")
	t.Tags("NoteRepository")

	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		query := fmt.Sprintf(deleteNoteQuery, s.repo.DbConfigs.SchemaName)
		s.mock.ExpectBegin()
		s.mock.ExpectExec(query).
			WithArgs(1).
			WillReturnError(sql.ErrConnDone)
		s.mock.ExpectRollback()

		err := s.repo.DeleteNote(1)

		sCtx.Assert().NotNil(err)
		sCtx.Assert().Equal(bl.DatabaseError, err.ErrNum)
	})
}

func (s *NoteRepositorySuite) Test_GetAllPublicNotes_Success(t provider.T) {
	t.Title("GetAllPublicNotes: Success")
	t.Tags("NoteRepository")

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		query := fmt.Sprintf(getAllPublicNotesQuery, s.repo.DbConfigs.SchemaName)
		s.mock.ExpectQuery(query).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "access", "name", "content_type", "likes", "dislikes", "registration_date", "owner_id", "section_id",
			}).
				AddRow(1, 1, "Note1", 1, 10, 1, ZEROTIME, 1, 1).
				AddRow(2, 1, "Note 2", 1, 20, 0, ZEROTIME, 2, 2))

		notes, err := s.repo.GetAllPublicNotes()

		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
		sCtx.Assert().Len(notes, 2)
		sCtx.Assert().Equal("Note1", notes[0].Name)
	})
}

func (s *NoteRepositorySuite) Test_GetAllPublicNotes_Failure(t provider.T) {
	t.Title("GetAllPublicNotes: Failure")
	t.Tags("NoteRepository")

	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		query := fmt.Sprintf(getAllPublicNotesQuery, s.repo.DbConfigs.SchemaName)
		s.mock.ExpectQuery(query).
			WillReturnError(sql.ErrConnDone)

		notes, err := s.repo.GetAllPublicNotes()

		sCtx.Assert().Nil(notes)
		sCtx.Assert().NotNil(err)
		sCtx.Assert().Equal(bl.DatabaseError, err.ErrNum)
	})
}

func (s *NoteRepositorySuite) Test_AddNoteToCollection_Success(t provider.T) {
	t.Title("AddNoteToCollection: Success")
	t.Tags("NoteRepository")

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		query := fmt.Sprintf(addNoteToCollectionQuery, s.repo.DbConfigs.SchemaName)
		s.mock.ExpectBegin()
		s.mock.ExpectExec(query).
			WithArgs(1, 2).
			WillReturnResult(sqlmock.NewResult(0, 1))
		s.mock.ExpectCommit()

		err := s.repo.AddNoteToCollection(2, 1)

		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *NoteRepositorySuite) Test_AddNoteToCollection_Failure(t provider.T) {
	t.Title("AddNoteToCollection: Failure")
	t.Tags("NoteRepository")

	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		query := fmt.Sprintf(addNoteToCollectionQuery, s.repo.DbConfigs.SchemaName)
		s.mock.ExpectBegin()
		s.mock.ExpectExec(query).
			WithArgs(1, 2).
			WillReturnError(sql.ErrConnDone)
		s.mock.ExpectRollback()

		err := s.repo.AddNoteToCollection(2, 1)

		sCtx.Assert().NotNil(err)
		sCtx.Assert().Equal(bl.DatabaseError, err.ErrNum)
	})
}

func (s *NoteRepositorySuite) Test_DeleteNoteFromCollection_Success(t provider.T) {
	t.Title("DeleteNoteFromCollection: Success")
	t.Tags("NoteRepository")

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		query := fmt.Sprintf(deleteNoteFromCollectionQuery, s.repo.DbConfigs.SchemaName)
		s.mock.ExpectBegin()
		s.mock.ExpectExec(query).
			WithArgs(1, 2).
			WillReturnResult(sqlmock.NewResult(0, 1))
		s.mock.ExpectCommit()

		err := s.repo.DeleteNoteFromCollection(2, 1)

		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *NoteRepositorySuite) Test_DeleteNoteFromCollection_Failure(t provider.T) {
	t.Title("DeleteNoteFromCollection: Failure")
	t.Tags("NoteRepository")

	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		query := fmt.Sprintf(deleteNoteFromCollectionQuery, s.repo.DbConfigs.SchemaName)
		s.mock.ExpectBegin()
		s.mock.ExpectExec(query).
			WithArgs(1, 2).
			WillReturnError(sql.ErrConnDone)
		s.mock.ExpectRollback()

		err := s.repo.DeleteNoteFromCollection(2, 1)

		sCtx.Assert().NotNil(err)
		sCtx.Assert().Equal(bl.DatabaseError, err.ErrNum)
	})
}

func (s *NoteRepositorySuite) Test_AddNote_Success(t provider.T) {
	t.Title("AddNote: Success")
	t.Tags("NoteRepository")

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		id := 1
		content := data_builders.NewContentBuilder().
			WithText([]byte("Note text content")).
			WithTextExt(".txt").
			WithImg([]byte("image_data")).
			WithImgExt(".png").
			WithRaw([]byte("raw_data")).
			WithRawExt(".bin").
			Build()
		note := data_builders.NewNoteBuilder().
			WithId(id).
			WithAccess(1).
			WithName("Test Note").
			WithLikes(10).
			WithDislikes(2).
			WithRegistrationDate(ZEROTIME).
			WithOwnerID(100).
			WithSectionID(5).
			WithContent(*content).
			Build()

		mainQuery := fmt.Sprintf(addNoteInfoQuery, s.repo.DbConfigs.SchemaName)
		txtQuery := fmt.Sprintf(addNoteTextQuery, s.repo.DbConfigs.SchemaName)
		imgQuery := fmt.Sprintf(addNoteImageQuery, s.repo.DbConfigs.SchemaName)
		rawQuery := fmt.Sprintf(addNoteRawDataQuery, s.repo.DbConfigs.SchemaName)

		s.mock.ExpectBegin()
		s.mock.ExpectQuery(mainQuery).
			WithArgs(note.Access, note.Name, 1, note.Likes, note.Dislikes, note.RegistrationDate, note.OwnerID, note.SectionID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(note.Id))
		s.mock.ExpectCommit()

		s.mock.ExpectBegin()
		s.mock.ExpectExec(txtQuery).
			WithArgs(content.Text, content.TextExt, note.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectCommit()

		s.mock.ExpectBegin()
		s.mock.ExpectExec(imgQuery).
			WithArgs(content.Img, content.ImgExt, note.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectCommit()

		s.mock.ExpectBegin()
		s.mock.ExpectExec(rawQuery).
			WithArgs(content.Raw, content.RawExt, note.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectCommit()

		realId, err := s.repo.AddNote(note)

		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
		sCtx.Assert().NotNil(note)
		sCtx.Assert().Equal(note.Id, realId)
	})
}

func (s *NoteRepositorySuite) Test_AddNote_NoteIsNil(t provider.T) {
	t.Title("AddNote: Note Is Nil")
	t.Tags("NoteRepository")

	t.WithNewStep("Note Is Nil", func(sCtx provider.StepCtx) {
		var note *models.Note

		_, err := s.repo.AddNote(note)

		sCtx.Assert().Nil(note)
		sCtx.Assert().NotNil(err)
		sCtx.Assert().Equal(bl.ErrInParameter, err.ErrNum)
	})
}

func (s *NoteRepositorySuite) Test_UpdateNoteContent_Success(t provider.T) {
	t.Title("UpdateNoteContent: Success")
	t.Tags("NoteRepository")

	t.WithNewStep("Update content successfully", func(sCtx provider.StepCtx) {
		content := data_builders.NewContentBuilder().
			WithText([]byte("Updated text content")).
			WithTextExt(".txt").
			WithImg([]byte("updated_image_data")).
			WithImgExt(".png").
			WithRaw([]byte("updated_raw_data")).
			WithRawExt(".bin").
			Build()
		note := data_builders.NewNoteBuilder().
			WithId(1).
			WithContent(*content).
			Build()

		textQuery := fmt.Sprintf(updateNoteTextContentQuery, s.repo.DbConfigs.SchemaName)
		imageQuery := fmt.Sprintf(updateNoteImageContentQuery, s.repo.DbConfigs.SchemaName)
		rawQuery := fmt.Sprintf(updateNoteRawContentQuery, s.repo.DbConfigs.SchemaName)

		s.mock.ExpectBegin()
		s.mock.ExpectExec(textQuery).
			WithArgs(note.Content.Text, note.Content.TextExt, note.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectCommit()

		s.mock.ExpectBegin()
		s.mock.ExpectExec(imageQuery).
			WithArgs(note.Content.Img, note.Content.ImgExt, note.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectCommit()

		s.mock.ExpectBegin()
		s.mock.ExpectExec(rawQuery).
			WithArgs(note.Content.Raw, note.Content.RawExt, note.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectCommit()

		err := s.repo.UpdateNoteContent(note)

		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *NoteRepositorySuite) Test_UpdateNoteContent_NoteIsNil(t provider.T) {
	t.Title("UpdateNoteContent: Note is nil")
	t.Tags("NoteRepository")

	t.WithNewStep("Fail if note is nil", func(sCtx provider.StepCtx) {
		var note *models.Note

		err := s.repo.UpdateNoteContent(note)

		sCtx.Assert().NotNil(err)
		sCtx.Assert().Equal(bl.ErrInParameter, err.ErrNum)
	})
}

func (s *NoteRepositorySuite) Test_UpdateNoteInfo_Success(t provider.T) {
	t.Title("UpdateNoteInfo: Success")
	t.Tags("NoteRepository")

	t.WithNewStep("Update note info successfully", func(sCtx provider.StepCtx) {
		note := data_builders.NewNoteBuilder().
			WithId(1).
			WithAccess(1).
			WithName("Updated Note").
			WithLikes(100).
			WithDislikes(5).
			WithRegistrationDate(ZEROTIME).
			WithOwnerID(200).
			WithSectionID(10).
			Build()

		query := fmt.Sprintf(updateNoteInfoQuery, s.repo.DbConfigs.SchemaName)

		s.mock.ExpectBegin()
		s.mock.ExpectExec(query).
			WithArgs(note.Access, note.Name, 1, note.Likes, note.Dislikes, note.RegistrationDate, note.OwnerID, note.SectionID, note.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectCommit()

		err := s.repo.UpdateNoteInfo(note)

		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *NoteRepositorySuite) Test_UpdateNoteInfo_NoteIsNil(t provider.T) {
	t.Title("UpdateNoteInfo: Note is nil")
	t.Tags("NoteRepository")

	t.WithNewStep("Fail if note is nil", func(sCtx provider.StepCtx) {
		var note *models.Note

		err := s.repo.UpdateNoteInfo(note)

		sCtx.Assert().NotNil(err)
		sCtx.Assert().Equal(bl.ErrInParameter, err.ErrNum)
	})
}
