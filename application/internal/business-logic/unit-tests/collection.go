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

type CollectionServiceTestSuite struct {
	suite.Suite
}

type extDependenciesCollection struct {
	collectionRepo *mocks.MockCollectionRepository
	logger         *mylogger.MyLogger
}

func initExtDependenciesCollection(t provider.T) *extDependenciesCollection {
	mockRepo := mocks.NewMockCollectionRepository(t)

	f := &extDependenciesCollection{
		collectionRepo: mockRepo,
		logger:         &mylogger.MyLogger{},
	}

	f.logger.InitLogger("./logs/unit-test-collection.log", "debug")

	return f
}

func (s *CollectionServiceTestSuite) TestCollectionService_GetCollection_OK(t provider.T) {
	t.Title("Get collection: OK")
	t.Tags("CollectionService")
	t.Parallel()

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		df := initExtDependenciesCollection(t)
		inCollection := builders.NewCollectionBuilder().Build()
		outCollection := builders.NewCollectionBuilder().WithId(inCollection.Id)
		myerr := myerror_builder.NewMyErrorBuilder().WithErrNum(bl.Ok).Build()

		df.collectionRepo.EXPECT().
			GetCollectionByID(inCollection.Id).
			Return(outCollection.Collection, myerr).Once()

		collectionService := bl.NewCollectionService(df.collectionRepo, df.logger)

		collection, err := collectionService.GetCollection(inCollection.Id, "", bl.SearchByID)

		sCtx.Assert().Equal(err.ErrNum, bl.Ok)
		sCtx.Assert().Equal(collection.Id, outCollection.Collection.Id)
	})
}

func (s *CollectionServiceTestSuite) TestCollectionService_GetCollection_ErrSearchParameter(t provider.T) {
	t.Title("Get collection: ErrSearchParameter")
	t.Tags("CollectionService")
	t.Parallel()

	t.WithNewStep("ErrSearchParameter", func(sCtx provider.StepCtx) {
		df := initExtDependenciesCollection(t)
		inCollection := builders.NewCollectionBuilder().Build()
		// myerr := myerror_builder.NewMyErrorBuilder().WithErrNum(bl.ErrAccessDenied).Build()

		// df.collectionRepo.EXPECT().
		// 	GetCollection(inCollection.Id, inCollection.Name, bl.SearchByID).
		// 	Return(nil, myerr.MyError).
		// 	Once()

		collectionService := bl.NewCollectionService(df.collectionRepo, df.logger)

		collection, err := collectionService.GetCollection(inCollection.Id, "", 100500)

		sCtx.Assert().Nil(collection)
		sCtx.Assert().Equal(err.ErrNum, bl.ErrSearchParameter)
	})
}

func (s *CollectionServiceTestSuite) TestCollectionService_GetAllCollections_OK(t provider.T) {
	t.Title("Get all collections: OK")
	t.Tags("CollectionService")
	t.Parallel()

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		df := initExtDependenciesCollection(t)
		reqUser := builders.NewUserBuilder().WithRole(bl.Admin)
		collections := []*models.Collection{
			builders.NewCollectionBuilder().Build(),
			builders.NewCollectionBuilder().Build(),
			builders.NewCollectionBuilder().Build(),
		}
		myerr := myerror_builder.NewMyErrorBuilder().WithErrNum(bl.Ok).Build()

		df.collectionRepo.EXPECT().
			GetAllCollections().
			Return(collections, myerr).Once()

		collectionService := bl.NewCollectionService(df.collectionRepo, df.logger)

		result, err := collectionService.GetAllCollections(reqUser.User)

		sCtx.Assert().Equal(err.ErrNum, bl.Ok)
		sCtx.Assert().Equal(len(result), 3)
	})
}

func (s *CollectionServiceTestSuite) TestCollectionService_GetAllCollections_AccessDenied(t provider.T) {
	t.Title("Get all collections: Access Denied")
	t.Tags("CollectionService")
	t.Parallel()

	t.WithNewStep("Access Denied", func(sCtx provider.StepCtx) {
		df := initExtDependenciesCollection(t)
		reqUser := builders.NewUserBuilder().WithRole(bl.Reader)
		// collections := []*models.Collection{
		// 	builders.NewCollectionBuilder().Build(),
		// 	builders.NewCollectionBuilder().Build(),
		// 	builders.NewCollectionBuilder().Build(),
		// }
		// myerr := myerror_builder.NewMyErrorBuilder().WithErrNum(bl.Ok).Build()

		// df.collectionRepo.EXPECT().
		// 	GetAllCollections().
		// 	Return(collections, myerr.MyError).
		// 	Once()

		collectionService := bl.NewCollectionService(df.collectionRepo, df.logger)

		colls, err := collectionService.GetAllCollections(reqUser.User)

		sCtx.Assert().Nil(colls)
		sCtx.Assert().Equal(err.ErrNum, bl.ErrAccessDenied)
	})
}

func (s *CollectionServiceTestSuite) TestCollectionService_GetAllUsersCollections_OK(t provider.T) {
	t.Title("Get all users collections: OK")
	t.Tags("CollectionService")
	t.Parallel()

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		df := initExtDependenciesCollection(t)
		reqUser := builders.NewUserBuilder().WithRole(bl.Admin)
		collections := []*models.Collection{
			builders.NewCollectionBuilder().Build(),
			builders.NewCollectionBuilder().Build(),
			builders.NewCollectionBuilder().Build(),
		}
		myerr := myerror_builder.NewMyErrorBuilder().WithErrNum(bl.Ok).Build()

		df.collectionRepo.EXPECT().
			GetAllUserCollections(reqUser.User).
			Return(collections, myerr).Once()

		collectionService := bl.NewCollectionService(df.collectionRepo, df.logger)

		result, err := collectionService.GetAllUsersCollections(reqUser.User)

		sCtx.Assert().Equal(err.ErrNum, bl.Ok)
		sCtx.Assert().Equal(len(result), 3)
	})
}

func (s *CollectionServiceTestSuite) TestCollectionService_AddCollection_OK(t provider.T) {
	t.Title("Add collection: OK")
	t.Tags("CollectionService")
	t.Parallel()

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		df := initExtDependenciesCollection(t)
		newCollection := builders.NewCollectionBuilder().Build()
		myerr := myerror_builder.NewMyErrorBuilder().WithErrNum(bl.Ok).Build()

		df.collectionRepo.EXPECT().
			AddCollection(newCollection).
			Return(myerr).Once()

		collectionService := bl.NewCollectionService(df.collectionRepo, df.logger)

		id, err := collectionService.AddCollection(newCollection)

		sCtx.Assert().Equal(err.ErrNum, bl.Ok)
		sCtx.Assert().Equal(id, 0)
	})
}

func (s *CollectionServiceTestSuite) TestCollectionService_DeleteCollection_OK(t provider.T) {
	t.Title("Delete collection: OK")
	t.Tags("CollectionService")
	t.Parallel()

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		df := initExtDependenciesCollection(t)
		reqUser := builders.NewUserBuilder().WithRole(bl.Admin)
		newCollection := builders.NewCollectionBuilder().WithId(1).WithOwnerID(reqUser.User.Id)
		myerr := myerror_builder.NewMyErrorBuilder().WithErrNum(bl.Ok).Build()

		df.collectionRepo.EXPECT().DeleteCollection(newCollection.Collection.Id).Return(myerr).Once()
		df.collectionRepo.EXPECT().GetCollectionByID(newCollection.Collection.Id).Return(newCollection.Collection, myerr).Once()

		collectionService := bl.NewCollectionService(df.collectionRepo, df.logger)

		err := collectionService.DeleteCollection(1, reqUser.User)

		sCtx.Assert().Equal(err.ErrNum, bl.Ok)
	})
}

func (s *CollectionServiceTestSuite) TestCollectionService_DeleteCollection_AccessDenied(t provider.T) {
	t.Title("Delete collection: Access Denied")
	t.Tags("CollectionService")
	t.Parallel()

	t.WithNewStep("Access Denied", func(sCtx provider.StepCtx) {
		df := initExtDependenciesCollection(t)
		reqUser := builders.NewUserBuilder().WithRole(bl.Reader)
		newCollection := builders.NewCollectionBuilder().WithId(1).WithOwnerID(reqUser.User.Id + 1)
		// myerr := myerror_builder.NewMyErrorBuilder().WithErrNum(bl.ErrAccessDenied).Build()
		myok := myerror_builder.NewMyErrorBuilder().WithErrNum(bl.Ok).Build()

		// df.collectionRepo.EXPECT().DeleteCollection(newCollection.Collection.Id).Return(myerr).Once()
		df.collectionRepo.EXPECT().GetCollectionByID(newCollection.Collection.Id).Return(newCollection.Collection, myok).Once()

		collectionService := bl.NewCollectionService(df.collectionRepo, df.logger)

		err := collectionService.DeleteCollection(newCollection.Collection.Id, reqUser.User)

		sCtx.Assert().Equal(err.ErrNum, bl.ErrAccessDenied)
	})
}

func (s *CollectionServiceTestSuite) TestCollectionService_UpdateCollection_OK(t provider.T) {
	t.Title("Update collection: OK")
	t.Tags("CollectionService")
	t.Parallel()

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		df := initExtDependenciesCollection(t)
		updatedCollection := builders.NewCollectionBuilder().Build()
		myerr := myerror_builder.NewMyErrorBuilder().WithErrNum(bl.Ok).Build()

		df.collectionRepo.EXPECT().
			UpdateCollection(updatedCollection).
			Return(myerr).Once()

		collectionService := bl.NewCollectionService(df.collectionRepo, df.logger)

		err := collectionService.UpdateCollection(updatedCollection)

		sCtx.Assert().Equal(err.ErrNum, bl.Ok)
	})
}

func (s *CollectionServiceTestSuite) TestCollectionService_UpdateCollection_AccessDenied(t provider.T) {
	t.Title("Update collection: Access Denied")
	t.Tags("CollectionService")
	t.Parallel()

	t.WithNewStep("Access Denied", func(sCtx provider.StepCtx) {
		df := initExtDependenciesCollection(t)
		updatedCollection := builders.NewCollectionBuilder().Build()
		myerr := myerror_builder.NewMyErrorBuilder().WithErrNum(bl.ErrAccessDenied).Build()

		df.collectionRepo.EXPECT().
			UpdateCollection(updatedCollection).
			Return(myerr).Once()

		collectionService := bl.NewCollectionService(df.collectionRepo, df.logger)

		err := collectionService.UpdateCollection(updatedCollection)

		sCtx.Assert().Equal(err.ErrNum, bl.ErrAccessDenied)
	})
}

func (s *CollectionServiceTestSuite) TestCollectionService_GetAllNotesInCollection_OK(t provider.T) {
	t.Title("Get all notes in collection: OK")
	t.Tags("CollectionService")
	t.Parallel()

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		df := initExtDependenciesCollection(t)
		collection := builders.NewCollectionBuilder().Build()
		notes := []*models.Note{
			builders.NewNoteBuilder().Build(),
			builders.NewNoteBuilder().Build(),
			builders.NewNoteBuilder().Build(),
		}
		myerr := myerror_builder.NewMyErrorBuilder().WithErrNum(bl.Ok).Build()

		df.collectionRepo.EXPECT().
			GetAllNotesInCollection(collection).
			Return(notes, myerr).Once()

		collectionService := bl.NewCollectionService(df.collectionRepo, df.logger)

		result, err := collectionService.GetAllNotesInCollection(collection)

		sCtx.Assert().Equal(err.ErrNum, bl.Ok)
		sCtx.Assert().Equal(len(result), 3)
	})
}

func (s *CollectionServiceTestSuite) TestCollectionService_GetAllNotesInCollection_AccessDenied(t provider.T) {
	t.Title("Get all notes in collection: Access Denied")
	t.Tags("CollectionService")
	t.Parallel()

	t.WithNewStep("Access Denied", func(sCtx provider.StepCtx) {
		df := initExtDependenciesCollection(t)
		collection := builders.NewCollectionBuilder().Build()
		myerr := myerror_builder.NewMyErrorBuilder().WithErrNum(bl.ErrAccessDenied).Build()

		df.collectionRepo.EXPECT().
			GetAllNotesInCollection(collection).
			Return(nil, myerr).Once()

		collectionService := bl.NewCollectionService(df.collectionRepo, df.logger)

		result, err := collectionService.GetAllNotesInCollection(collection)

		sCtx.Assert().Nil(result)
		sCtx.Assert().Equal(err.ErrNum, bl.ErrAccessDenied)
	})
}
