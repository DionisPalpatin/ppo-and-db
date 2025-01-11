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

type CollectionRepoIntegrationTest struct {
    suite.Suite
    db              *sql.DB
    collectionRepo da.CollectionRepository

	idCollToDel      int
	idCollToUpd      int
	idCollToGetNotes int
}

func (s *CollectionRepoIntegrationTest) BeforeEach(t provider.T) {
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
        Host:     pgInfo.Host,
        User:     pgInfo.User,
        Password: pgInfo.Password,
        Port:     5432,
        Name:     pgInfo.DBName,
        DB:       s.db,
		SchemaName: "interg_tests",
    }

    logger := mylogger.MyLogger{}
    logger.InitLogger("./logs/integration-test-collection.log", "debug")
 
    s.collectionRepo = da.NewCollectionRepository(&dbconfig, &logger)

	s.idCollToDel      = 10
	s.idCollToUpd      = 9
	s.idCollToGetNotes = 1
}

func (s *CollectionRepoIntegrationTest) AfterEach(t provider.T) {
    _ = s.db.Close()
}

func (s *CollectionRepoIntegrationTest) TestGetCollectionByIDSuccess(t provider.T) {
    t.Title("GetCollectionByID: Success")
    t.Tags("CollectionRepoIntegrationTest")

    t.WithNewStep("Fetch existing collection", func(sCtx provider.StepCtx) {
        _, err := s.collectionRepo.GetCollectionByID(1)

        sCtx.Assert().Equal(bl.Ok, err.ErrNum)
    })
}

func (s *CollectionRepoIntegrationTest) TestGetCollectionByIDNotFound(t provider.T) {
    t.Title("GetCollectionByID: Not Found")
    t.Tags("CollectionRepoIntegrationTest")

    t.WithNewStep("Collection does not exist", func(sCtx provider.StepCtx) {
        collection, err := s.collectionRepo.GetCollectionByID(999)

        sCtx.Assert().Nil(collection)
        sCtx.Assert().NotNil(err)
        sCtx.Assert().Equal(bl.NoSuchColl, err.ErrNum)
    })
}

func (s *CollectionRepoIntegrationTest) TestGetCollectionByNameSuccess(t provider.T) {
    t.Title("GetCollectionByID: Success")
    t.Tags("CollectionRepoIntegrationTest")

    t.WithNewStep("Fetch existing collection", func(sCtx provider.StepCtx) {
        _, err := s.collectionRepo.GetCollectionByName("Collection 1")

        sCtx.Assert().Equal(bl.Ok, err.ErrNum)
    })
}

func (s *CollectionRepoIntegrationTest) TestGetCollectionByNameNotFound(t provider.T) {
    t.Title("GetCollectionByName: Not Found")
    t.Tags("CollectionRepoIntegrationTest")

    t.WithNewStep("Collection does not exist", func(sCtx provider.StepCtx) {
        collection, err := s.collectionRepo.GetCollectionByName("dajkflajflajfkljaklfjk")

        sCtx.Assert().Nil(collection)
        sCtx.Assert().NotNil(err)
        sCtx.Assert().Equal(bl.NoSuchColl, err.ErrNum)
    })
}

func (s *CollectionRepoIntegrationTest) TestGetAllCollectionsSuccess(t provider.T) {
    t.Title("GetAllCollections: Success")
    t.Tags("CollectionRepoIntegrationTest")

    t.WithNewStep("Fetch collections for a valid user", func(sCtx provider.StepCtx) {
        _, err := s.collectionRepo.GetAllCollections()

        sCtx.Assert().Equal(bl.Ok, err.ErrNum)
    })
}

func (s *CollectionRepoIntegrationTest) TestGetAllUserCollectionsSuccess(t provider.T) {
    t.Title("GetAllUserCollections: Success")
    t.Tags("CollectionRepoIntegrationTest")

    t.WithNewStep("Get all user's collections", func(sCtx provider.StepCtx) {
        user := data_builders.NewUserBuilder().WithUserID(1).Build()

        _, err := s.collectionRepo.GetAllUserCollections(user)

        sCtx.Assert().Equal(bl.Ok, err.ErrNum)
    })
}

func (s *CollectionRepoIntegrationTest) TestAddCollectionSuccess(t provider.T) {
    t.Title("AddCollection: Success")
    t.Tags("CollectionRepoIntegrationTest")

    t.WithNewStep("Add a new collection", func(sCtx provider.StepCtx) {
        collection := data_builders.NewCollectionBuilder().WithOwnerID(1).Build()
		
        _, err := s.collectionRepo.AddCollection(collection)

        sCtx.Assert().Equal(bl.Ok, err.ErrNum)
    })
}

// func (s *CollectionRepoIntegrationTest) TestAddCollectionDuplicate(t provider.T) {
//     t.Title("AddCollection: Duplicate Collection")
//     t.Tags("CollectionRepoIntegrationTest")

//     t.WithNewStep("Try to add a duplicate collection", func(sCtx provider.StepCtx) {
//         collection := data_builders.NewCollectionBuilder().WithName("Collection 1").WithOwnerID(1).Build()

//         _, err := s.collectionRepo.AddCollection(collection)

//         sCtx.Assert().NotNil(err)
//         sCtx.Assert().Equal(bl.DatabaseError, err.ErrNum)
//     })
// }

func (s *CollectionRepoIntegrationTest) TestDeleteCollectionSuccess(t provider.T) {
    t.Title("DeleteCollection: Success")
    t.Tags("CollectionRepoIntegrationTest")

    t.WithNewStep("Delete an existing collection", func(sCtx provider.StepCtx) {
        err := s.collectionRepo.DeleteCollection(s.idCollToDel)

        sCtx.Assert().Equal(bl.Ok, err.ErrNum)
    })
}

// func (s *CollectionRepoIntegrationTest) TestDeleteCollectionNotFound(t provider.T) {
//     t.Title("DeleteCollection: Collection Not Found")
//     t.Tags("CollectionRepoIntegrationTest")

//     t.WithNewStep("Try to delete a non-existent collection", func(sCtx provider.StepCtx) {
//         err := s.collectionRepo.DeleteCollection(999)

//         sCtx.Assert().NotNil(err)
//         sCtx.Assert().Equal(bl.NoSuchColl, err.ErrNum)
//     })
// }

func (s *CollectionRepoIntegrationTest) TestUpdateCollectionSuccess(t provider.T) {
    t.Title("UpdateCollection: Success")
    t.Tags("CollectionRepoIntegrationTest")

    t.WithNewStep("Update an existing collection", func(sCtx provider.StepCtx) {
        collection := data_builders.NewCollectionBuilder().WithId(s.idCollToUpd).WithName("UpdatedName").Build()

        err := s.collectionRepo.UpdateCollection(collection)

        sCtx.Assert().Equal(bl.Ok, err.ErrNum)
    })
}

func (s *CollectionRepoIntegrationTest) TestUpdateCollectionNotFound(t provider.T) {
    t.Title("UpdateCollection: Collection Not Found")
    t.Tags("CollectionRepoIntegrationTest")

    t.WithNewStep("Try to update a non-existent collection", func(sCtx provider.StepCtx) {
        collection := data_builders.NewCollectionBuilder().WithId(999).WithName("UpdatedName").WithOwnerID(1).Build()

        err := s.collectionRepo.UpdateCollection(collection)

        sCtx.Assert().NotNil(err)
        sCtx.Assert().Equal(bl.NoSuchColl, err.ErrNum)
    })
}

func (s *CollectionRepoIntegrationTest) TestGetAllNotesInCollectionSuccess(t provider.T) {
    t.Title("GetAllNotesInCollection: Success")
    t.Tags("CollectionRepoIntegrationTest")

    t.WithNewStep("Fetch notes for an existing collection", func(sCtx provider.StepCtx) {
        collection := data_builders.NewCollectionBuilder().WithId(s.idCollToGetNotes).Build()

        _, err := s.collectionRepo.GetAllNotesInCollection(collection)

        sCtx.Assert().Equal(bl.Ok, err.ErrNum)
    })
}

// func (s *CollectionRepoIntegrationTest) TestGetAllNotesInCollectionNotFound(t provider.T) {
//     t.Title("GetAllNotesInCollection: Collection Not Found")
//     t.Tags("CollectionRepoIntegrationTest")

//     t.WithNewStep("Try to fetch notes for a non-existent collection", func(sCtx provider.StepCtx) {
//         collection := data_builders.NewCollectionBuilder().WithId(999).Build()

//         notes, err := s.collectionRepo.GetAllNotesInCollection(collection)

//         sCtx.Assert().Empty(notes)
//         sCtx.Assert().NotNil(err)
//         sCtx.Assert().Equal(bl.NoSuchColl, err.ErrNum)
//     })
// }
