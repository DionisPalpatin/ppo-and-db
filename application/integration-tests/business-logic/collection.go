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

type CollectionServiceIntegrationTest struct {
	suite.Suite
	db                *sql.DB
	collectionService bl.ICollectionService

	idCollToDel      int
	idCollToUpd      int
	idCollToGetNotes int
}

func (s *CollectionServiceIntegrationTest) BeforeEach(t provider.T) {
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
	logger.InitLogger("./logs/integration-test-collection.log", "debug")

	collectionRepo := da.NewCollectionRepository(&dbconfig, &logger)
	s.collectionService = bl.NewCollectionService(&collectionRepo, &logger)

	s.idCollToDel = 11
	s.idCollToUpd = 8
	s.idCollToGetNotes = 1
}

func (s *CollectionServiceIntegrationTest) AfterEach(t provider.T) {
	_ = s.db.Close()
}

func (s *CollectionServiceIntegrationTest) TestGetCollectionSuccess(t provider.T) {
	t.Title("GetCollection: Success")
	t.Tags("CollectionServiceIntegrationTest")

	t.WithNewStep("Fetch existing collection", func(sCtx provider.StepCtx) {
		collection, err := s.collectionService.GetCollection(1, "", bl.SearchByID)

		sCtx.Assert().NotNil(collection)
		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *CollectionServiceIntegrationTest) TestGetCollectionNotFound(t provider.T) {
	t.Title("GetCollection: Not Found")
	t.Tags("CollectionServiceIntegrationTest")

	t.WithNewStep("Collection does not exist", func(sCtx provider.StepCtx) {
		collection, err := s.collectionService.GetCollection(999, "", bl.SearchByID)

		sCtx.Assert().Nil(collection)
		sCtx.Assert().NotNil(err)
		sCtx.Assert().Equal(bl.NoSuchColl, err.ErrNum)
	})
}

func (s *CollectionServiceIntegrationTest) TestGetAllCollectionsSuccess(t provider.T) {
	t.Title("GetAllCollections: Success")
	t.Tags("CollectionServiceIntegrationTest")

	t.WithNewStep("Fetch collections for a valid user", func(sCtx provider.StepCtx) {
		user := data_builders.NewUserBuilder().WithRole(bl.Admin).Build()

		_, err := s.collectionService.GetAllCollections(user)

		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *CollectionServiceIntegrationTest) TestGetAllCollectionsAccessDenied(t provider.T) {
	t.Title("GetAllCollections: AccessDenied")
	t.Tags("CollectionServiceIntegrationTest")

	t.WithNewStep("Access denied for reader", func(sCtx provider.StepCtx) {
		user := data_builders.NewUserBuilder().WithRole(bl.Reader).Build()

		collections, err := s.collectionService.GetAllCollections(user)

		sCtx.Assert().Empty(collections)
		sCtx.Assert().Equal(bl.ErrAccessDenied, err.ErrNum)
	})
}

func (s *CollectionServiceIntegrationTest) TestAddCollectionSuccess(t provider.T) {
	t.Title("AddCollection: Success")
	t.Tags("CollectionServiceIntegrationTest")

	t.WithNewStep("Add a new collection", func(sCtx provider.StepCtx) {
		collection := data_builders.NewCollectionBuilder().WithOwnerID(1).Build()

		_, err := s.collectionService.AddCollection(collection)

		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

// func (s *CollectionServiceIntegrationTest) TestAddCollectionDuplicate(t provider.T) {
//     t.Title("AddCollection: Duplicate Collection")
//     t.Tags("CollectionServiceIntegrationTest")

//     t.WithNewStep("Try to add a duplicate collection", func(sCtx provider.StepCtx) {
//         collection := data_builders.NewCollectionBuilder().WithName("Collection 1").WithOwnerID(1).Build()

//         _, err := s.collectionService.AddCollection(collection)

//         sCtx.Assert().NotNil(err)
//         sCtx.Assert().Equal(bl.DatabaseError, err.ErrNum)
//     })
// }

func (s *CollectionServiceIntegrationTest) TestDeleteCollectionSuccess(t provider.T) {
	t.Title("DeleteCollection: Success")
	t.Tags("CollectionServiceIntegrationTest")

	t.WithNewStep("Delete an existing collection", func(sCtx provider.StepCtx) {
		user := data_builders.NewUserBuilder().WithRole(bl.Admin).Build()

		err := s.collectionService.DeleteCollection(s.idCollToDel, user)

		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *CollectionServiceIntegrationTest) TestDeleteCollectionNotFound(t provider.T) {
	t.Title("DeleteCollection: Collection Not Found")
	t.Tags("CollectionServiceIntegrationTest")

	t.WithNewStep("Try to delete a non-existent collection", func(sCtx provider.StepCtx) {
		user := data_builders.NewUserBuilder().Build()

		err := s.collectionService.DeleteCollection(999, user)

		sCtx.Assert().NotNil(err)
		sCtx.Assert().Equal(bl.NoSuchColl, err.ErrNum)
	})
}

func (s *CollectionServiceIntegrationTest) TestUpdateCollectionSuccess(t provider.T) {
	t.Title("UpdateCollection: Success")
	t.Tags("CollectionServiceIntegrationTest")

	t.WithNewStep("Update an existing collection", func(sCtx provider.StepCtx) {
		collection := data_builders.NewCollectionBuilder().WithId(s.idCollToUpd).WithName("UpdatedName").Build()

		err := s.collectionService.UpdateCollection(collection)

		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *CollectionServiceIntegrationTest) TestUpdateCollectionNotFound(t provider.T) {
	t.Title("UpdateCollection: Collection Not Found")
	t.Tags("CollectionServiceIntegrationTest")

	t.WithNewStep("Try to update a non-existent collection", func(sCtx provider.StepCtx) {
		collection := data_builders.NewCollectionBuilder().WithId(999).WithName("UpdatedName").WithOwnerID(1).Build()

		err := s.collectionService.UpdateCollection(collection)

		sCtx.Assert().NotNil(err)
		sCtx.Assert().Equal(bl.NoSuchColl, err.ErrNum)
	})
}

func (s *CollectionServiceIntegrationTest) TestGetAllNotesInCollectionSuccess(t provider.T) {
	t.Title("GetAllNotesInCollection: Success")
	t.Tags("CollectionServiceIntegrationTest")

	t.WithNewStep("Fetch notes for an existing collection", func(sCtx provider.StepCtx) {
		collection := data_builders.NewCollectionBuilder().WithId(s.idCollToGetNotes).Build()

		_, err := s.collectionService.GetAllNotesInCollection(collection)

		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

// func (s *CollectionServiceIntegrationTest) TestGetAllNotesInCollectionNotFound(t provider.T) {
//     t.Title("GetAllNotesInCollection: Collection Not Found")
//     t.Tags("CollectionServiceIntegrationTest")

//     t.WithNewStep("Try to fetch notes for a non-existent collection", func(sCtx provider.StepCtx) {
//         collection := data_builders.NewCollectionBuilder().WithId(999).Build()

//         notes, err := s.collectionService.GetAllNotesInCollection(collection)

//         sCtx.Assert().Empty(notes)
//         sCtx.Assert().NotNil(err)
//         sCtx.Assert().Equal(bl.NoSuchColl, err.ErrNum)
//     })
// }
