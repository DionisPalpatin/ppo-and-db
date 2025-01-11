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

	bl "github.com/DionisPalpatin/ppo-and-db/application/internal/business-logic"
	da "github.com/DionisPalpatin/ppo-and-db/application/internal/data-access"

	// "github.com/DionisPalpatin/ppo-and-db/application/internal/models"
	data_builders "github.com/DionisPalpatin/ppo-and-db/application/internal/models/builders"
)

type CollectionRepositorySuite struct {
	suite.Suite

	db   *sql.DB
	mock sqlmock.Sqlmock
	repo da.CollectionRepository
	ctx  context.Context
}

func (s *CollectionRepositorySuite) BeforeEach(t provider.T) {
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
	logger.InitLogger("./logs/unit-tests-collection.log", "debug")

	s.repo = da.NewCollectionRepository(&confs, &logger)
	s.ctx = context.Background()
}

func (s *CollectionRepositorySuite) AfterEach(t provider.T) {
	s.db.Close()
}

func TestCollectionRepositorySuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(CollectionRepositorySuite))
}

func (s *CollectionRepositorySuite) Test_CollectionRepository_GetCollectionByID_Success(t provider.T) {
	t.Title("GetCollectionByID: Success")
	t.Tags("CollectionRepository")

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		query := fmt.Sprintf(getCollectionByIDQuery, s.repo.DbConfigs.SchemaName)
		s.mock.ExpectQuery(query).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "creation_date", "owner_id"}).
				AddRow(1, "Test Collection", ZEROTIME, 1))

		expectedCollection := data_builders.NewCollectionBuilder().WithId(1).WithName("Test Collection").Build()

		collection, err := s.repo.GetCollectionByID(1)

		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
		sCtx.Assert().Equal(expectedCollection.Id, collection.Id)
	})
}

func (s *CollectionRepositorySuite) Test_CollectionRepository_GetCollectionByID_Failure(t provider.T) {
	t.Title("GetCollectionByID: Failure")
	t.Tags("CollectionRepository")

	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		query := fmt.Sprintf(getCollectionByIDQuery, s.repo.DbConfigs.SchemaName)
		s.mock.ExpectQuery(query).
			WithArgs(1).
			WillReturnError(sql.ErrConnDone)

		collection, err := s.repo.GetCollectionByID(1)

		sCtx.Assert().Nil(collection)
		sCtx.Assert().Equal(bl.DatabaseError, err.ErrNum)
	})
}

func (s *CollectionRepositorySuite) Test_CollectionRepository_GetCollectionByName_Success(t provider.T) {
	t.Title("GetCollectionByName: Success")
	t.Tags("CollectionRepository")

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		query := fmt.Sprintf(getCollectionByNameQuery, s.repo.DbConfigs.SchemaName)
		s.mock.ExpectQuery(query).
			WithArgs("Test Collection").
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "creation_date", "owner_id"}).
				AddRow(1, "Test Collection", ZEROTIME, 1))

		collection, err := s.repo.GetCollectionByName("Test Collection")

		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
		sCtx.Assert().Equal("Test Collection", collection.Name)
	})
}

func (s *CollectionRepositorySuite) Test_CollectionRepository_GetCollectionByName_Failure(t provider.T) {
	t.Title("GetCollectionByName: Failure")
	t.Tags("CollectionRepository")

	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		query := fmt.Sprintf(getCollectionByNameQuery, s.repo.DbConfigs.SchemaName)
		s.mock.ExpectQuery(query).
			WithArgs("Test Collection").
			WillReturnError(sql.ErrConnDone)

		collection, err := s.repo.GetCollectionByName("Test Collection")

		sCtx.Assert().Nil(collection)
		sCtx.Assert().Equal(bl.DatabaseError, err.ErrNum)
	})
}

func (s *CollectionRepositorySuite) Test_CollectionRepository_GetAllCollections_Success(t provider.T) {
	t.Title("GetAllCollections: Success")
	t.Tags("CollectionRepository")

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		query := fmt.Sprintf(getAllCollectionsQuery, s.repo.DbConfigs.SchemaName)
		s.mock.ExpectQuery(query).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "creation_date", "owner_id"}).
				AddRow(1, "Collection 1", ZEROTIME, 1).
				AddRow(2, "Collection 2", ZEROTIME, 1))

		collections, err := s.repo.GetAllCollections()

		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
		sCtx.Assert().Len(collections, 2)
	})
}

func (s *CollectionRepositorySuite) Test_CollectionRepository_GetAllCollections_Failure(t provider.T) {
	t.Title("GetAllCollections: Failure")
	t.Tags("CollectionRepository")

	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		query := fmt.Sprintf(getAllCollectionsQuery, s.repo.DbConfigs.SchemaName)
		s.mock.ExpectQuery(query).
			WillReturnError(sql.ErrConnDone)

		collections, err := s.repo.GetAllCollections()

		sCtx.Assert().Nil(collections)
		sCtx.Assert().Equal(bl.DatabaseError, err.ErrNum)
	})
}

func (s *CollectionRepositorySuite) Test_CollectionRepository_GetAllUserCollections_Success(t provider.T) {
	t.Title("GetAllUserCollections: Success")
	t.Tags("CollectionRepository")

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		query := fmt.Sprintf(getAllUserCollectionsQuery, s.repo.DbConfigs.SchemaName, s.repo.DbConfigs.SchemaName, s.repo.DbConfigs.SchemaName)
		user := data_builders.NewUserBuilder().WithUserID(1).Build()
		s.mock.ExpectQuery(query).
			WithArgs(user.Id).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "creation_date", "owner_id"}).
				AddRow(1, "User Collection 1", ZEROTIME, 1))

		collections, err := s.repo.GetAllUserCollections(user)

		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
		sCtx.Assert().Len(collections, 1)
	})
}

func (s *CollectionRepositorySuite) Test_CollectionRepository_GetAllUserCollections_Failure(t provider.T) {
	t.Title("GetAllUserCollections: Failure")
	t.Tags("CollectionRepository")

	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		query := fmt.Sprintf(getAllUserCollectionsQuery, s.repo.DbConfigs.SchemaName, s.repo.DbConfigs.SchemaName, s.repo.DbConfigs.SchemaName)
		user := data_builders.NewUserBuilder().WithUserID(1).Build()
		s.mock.ExpectQuery(query).
			WithArgs(user.Id).
			WillReturnError(sql.ErrConnDone)

		collections, err := s.repo.GetAllUserCollections(user)

		sCtx.Assert().Nil(collections)
		sCtx.Assert().Equal(bl.DatabaseError, err.ErrNum)
	})
}

func (s *CollectionRepositorySuite) Test_CollectionRepository_AddCollection_Success(t provider.T) {
	t.Title("AddCollection: Success")
	t.Tags("CollectionRepository")

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		query := fmt.Sprintf(addCollectionQuery, s.repo.DbConfigs.SchemaName)
		collection := data_builders.NewCollectionBuilder().WithName("New Collection").Build()
		s.mock.ExpectBegin()
		s.mock.ExpectQuery(query).
			WithArgs(collection.Name, collection.CreationDate, collection.OwnerID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		s.mock.ExpectCommit()

		id, err := s.repo.AddCollection(collection)

		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
		sCtx.Assert().Equal(1, id)
	})
}

func (s *CollectionRepositorySuite) Test_CollectionRepository_AddCollection_Failure(t provider.T) {
	t.Title("AddCollection: Failure")
	t.Tags("CollectionRepository")

	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		query := fmt.Sprintf(addCollectionQuery, s.repo.DbConfigs.SchemaName)
		collection := data_builders.NewCollectionBuilder().WithName("New Collection").Build()
		s.mock.ExpectBegin()
		s.mock.ExpectQuery(query).
			WithArgs(collection.Name, collection.CreationDate, collection.OwnerID).
			WillReturnError(sql.ErrConnDone)
		s.mock.ExpectRollback()

		id, err := s.repo.AddCollection(collection)

		sCtx.Assert().Equal(0, id)
		sCtx.Assert().Equal(bl.DatabaseError, err.ErrNum)
	})
}

func (s *CollectionRepositorySuite) Test_CollectionRepository_DeleteCollection_Success(t provider.T) {
	t.Title("DeleteCollection: Success")
	t.Tags("CollectionRepository")

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		query1 := fmt.Sprintf(deleteCollectionQuery1, s.repo.DbConfigs.SchemaName)
		query2 := fmt.Sprintf(deleteCollectionQuery2, s.repo.DbConfigs.SchemaName)

		s.mock.ExpectBegin()
		s.mock.ExpectExec(query1).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectExec(query2).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(1, 1))	
		s.mock.ExpectCommit()

		err := s.repo.DeleteCollection(1)

		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *CollectionRepositorySuite) Test_CollectionRepository_DeleteCollection_Failure(t provider.T) {
	t.Title("DeleteCollection: Failure")
	t.Tags("CollectionRepository")

	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		query1 := fmt.Sprintf(deleteCollectionQuery1, s.repo.DbConfigs.SchemaName)
		query2 := fmt.Sprintf(deleteCollectionQuery2, s.repo.DbConfigs.SchemaName)

		s.mock.ExpectBegin()
		s.mock.ExpectExec(query1).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectExec(query2).
			WithArgs(1).
			WillReturnError(sql.ErrConnDone)
		s.mock.ExpectRollback()

		err := s.repo.DeleteCollection(1)

		sCtx.Assert().Equal(bl.DatabaseError, err.ErrNum)
	})
}

func (s *CollectionRepositorySuite) Test_CollectionRepository_UpdateCollection_Success(t provider.T) {
	t.Title("UpdateCollection: Success")
	t.Tags("CollectionRepository")

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		query := fmt.Sprintf(updateCollectionQuery, s.repo.DbConfigs.SchemaName)
		collection := data_builders.NewCollectionBuilder().WithId(1).WithName("Updated Collection").Build()
		s.mock.ExpectBegin()
		s.mock.ExpectExec(query).
			WithArgs(collection.Name, collection.CreationDate, collection.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectCommit()

		err := s.repo.UpdateCollection(collection)

		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *CollectionRepositorySuite) Test_CollectionRepository_UpdateCollection_Failure(t provider.T) {
	t.Title("UpdateCollection: Failure")
	t.Tags("CollectionRepository")

	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		query := fmt.Sprintf(updateCollectionQuery, s.repo.DbConfigs.SchemaName)
		collection := data_builders.NewCollectionBuilder().WithId(1).WithName("Updated Collection").Build()
		s.mock.ExpectBegin()
		s.mock.ExpectExec(query).
			WithArgs(collection.Name, collection.CreationDate, collection.Id).
			WillReturnError(sql.ErrConnDone)
		s.mock.ExpectRollback()

		err := s.repo.UpdateCollection(collection)

		sCtx.Assert().Equal(bl.DatabaseError, err.ErrNum)
	})
}

func (s *CollectionRepositorySuite) Test_CollectionRepository_GetAllNotesInCollection_Success(t provider.T) {
	t.Title("GetAllNotesInCollection: Success")
	t.Tags("CollectionRepository")

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		// Подготовка данных
		query := fmt.Sprintf(getAllNotesInCollectionQuery, s.repo.DbConfigs.SchemaName, s.repo.DbConfigs.SchemaName)
		collection := data_builders.NewCollectionBuilder().WithId(1).Build()

		// Настройка mock-результатов
		s.mock.ExpectQuery(query).
			WithArgs(collection.Id).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "access", "name", "content_type", "likes", "dislikes", 
				"registration_date", "owner_id", "section_id"}).
				AddRow(1, 0, "Test Note", 1, 10, 2, ZEROTIME, 1, 1))

		// Вызов тестируемого метода
		notes, err := s.repo.GetAllNotesInCollection(collection)

		// Проверка результатов
		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
		sCtx.Assert().Len(notes, 1)
		sCtx.Assert().Equal(1, notes[0].Id)
		sCtx.Assert().Equal("Test Note", notes[0].Name)
		sCtx.Assert().Equal(10, notes[0].Likes)
	})
}

func (s *CollectionRepositorySuite) Test_CollectionRepository_GetAllNotesInCollection_Failure(t provider.T) {
	t.Title("GetAllNotesInCollection: Failure")
	t.Tags("CollectionRepository")

	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		query := fmt.Sprintf(getAllNotesInCollectionQuery, s.repo.DbConfigs.SchemaName, s.repo.DbConfigs.SchemaName)
		collection := data_builders.NewCollectionBuilder().WithId(1).Build()
		s.mock.ExpectQuery(query).
			WithArgs(collection.Id).
			WillReturnError(sql.ErrConnDone)

		notes, err := s.repo.GetAllNotesInCollection(collection)

		sCtx.Assert().Nil(notes)
		sCtx.Assert().Equal(bl.DatabaseError, err.ErrNum)
	})
}
