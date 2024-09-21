package UnitTests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"notebook_app/internal/business_logic"
	"notebook_app/internal/business_logic/UnitTests/mocks"
	"notebook_app/internal/models"
)

func TestGetCollection(t *testing.T) {
	t.Run("SuccessGetCollectionByID", func(t *testing.T) {
		retErr := bl.CreateError(bl.AllIsOk, nil, "")
		retCollection := &models.Collection{}

		mockCollectionRepo := new(mocks.MockICollectionRepository)
		mockCollectionRepo.On("GetCollectionByID", 1).Return(retCollection, retErr)

		csSrv := bl.CollectionService{}
		_, err := csSrv.GetCollection(1, "", bl.SearchByID, mockCollectionRepo)

		assert.NotNil(t, err)
		assert.Equal(t, err.ErrNum, bl.AllIsOk)

		mockCollectionRepo.AssertExpectations(t)
	})

	t.Run("SuccessGetCollectionByName", func(t *testing.T) {
		retErr := bl.CreateError(bl.AllIsOk, nil, "")
		retCollection := &models.Collection{}

		mockCollectionRepo := new(mocks.MockICollectionRepository)
		mockCollectionRepo.On("GetCollectionByName", "Collection 1").Return(retCollection, retErr)

		csSrv := bl.CollectionService{}
		_, err := csSrv.GetCollection(0, "Collection 1", bl.SearchByString, mockCollectionRepo)

		assert.NotNil(t, err)
		assert.Equal(t, err.ErrNum, bl.AllIsOk)

		mockCollectionRepo.AssertExpectations(t)
	})

	t.Run("ErrorSearchParameter", func(t *testing.T) {
		csSrv := bl.CollectionService{}
		_, err := csSrv.GetCollection(0, "", -1, nil)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrSearchParameter, err.ErrNum)
	})
}

func TestGetAllCollections(t *testing.T) {
	t.Run("SuccessGetAllCollections", func(t *testing.T) {
		retErr := bl.CreateError(bl.AllIsOk, nil, "")
		retCollections := []*models.Collection{
			{Id: 1, Name: "Collection 1"},
			{Id: 2, Name: "Collection 2"},
		}
		retUser := &models.User{Role: bl.Admin}

		mockCollectionRepo := new(mocks.MockICollectionRepository)
		mockCollectionRepo.On("GetAllCollections").Return(retCollections, retErr)

		csSrv := bl.CollectionService{}
		collections, err := csSrv.GetAllCollections(retUser, mockCollectionRepo)

		assert.NotNil(t, err)
		assert.Equal(t, err.ErrNum, bl.AllIsOk)
		assert.Len(t, collections, 2)
		assert.Equal(t, 1, collections[0].Id)
		assert.Equal(t, "Collection 1", collections[0].Name)
		assert.Equal(t, 2, collections[1].Id)
		assert.Equal(t, "Collection 2", collections[1].Name)

		mockCollectionRepo.AssertExpectations(t)
	})

	t.Run("ErrorAccessDenied", func(t *testing.T) {
		retUser := &models.User{Role: bl.Reader}

		csSrv := bl.CollectionService{}
		collections, err := csSrv.GetAllCollections(retUser, nil)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrAccessDenied, err.ErrNum)

		assert.Nil(t, collections)
	})
}

func TestGetAllUsersCollections(t *testing.T) {
	t.Run("SuccessGetAllUsersCollections", func(t *testing.T) {
		retErr := bl.CreateError(bl.AllIsOk, nil, "")
		retCollections := []*models.Collection{
			{Id: 1, Name: "Collection 1"},
			{Id: 2, Name: "Collection 2"},
		}
		retUser := &models.User{Id: 1}

		mockCollectionRepo := new(mocks.MockICollectionRepository)
		mockCollectionRepo.On("GetAllUserCollections", retUser).Return(retCollections, retErr)

		csSrv := bl.CollectionService{}
		collections, err := csSrv.GetAllUsersCollections(retUser, mockCollectionRepo)

		assert.NotNil(t, err)
		assert.Equal(t, err.ErrNum, bl.AllIsOk)
		assert.Len(t, collections, 2)
		assert.Equal(t, 1, collections[0].Id)
		assert.Equal(t, "Collection 1", collections[0].Name)
		assert.Equal(t, 2, collections[1].Id)
		assert.Equal(t, "Collection 2", collections[1].Name)

		mockCollectionRepo.AssertExpectations(t)
	})
}

func TestAddCollection(t *testing.T) {
	t.Run("SuccessAddCollection", func(t *testing.T) {
		retErr := bl.CreateError(bl.AllIsOk, nil, "")
		retCollection := &models.Collection{Id: 1, Name: "Collection 1"}

		mockCollectionRepo := new(mocks.MockICollectionRepository)
		mockCollectionRepo.On("AddCollection", retCollection).Return(retErr)

		csSrv := bl.CollectionService{}
		err := csSrv.AddCollection(retCollection, mockCollectionRepo)

		assert.NotNil(t, err)
		assert.Equal(t, err.ErrNum, bl.AllIsOk)

		mockCollectionRepo.AssertExpectations(t)
	})
}

func TestDeleteCollection(t *testing.T) {
	t.Run("SuccessDeleteCollection", func(t *testing.T) {
		retErr := bl.CreateError(bl.AllIsOk, nil, "")

		mockCollectionRepo := new(mocks.MockICollectionRepository)
		mockCollectionRepo.On("DeleteCollection", 1).Return(retErr)

		csSrv := bl.CollectionService{}
		err := csSrv.DeleteCollection(1, nil, mockCollectionRepo)

		assert.NotNil(t, err)
		assert.Equal(t, err.ErrNum, bl.AllIsOk)

		mockCollectionRepo.AssertExpectations(t)
	})
}

func TestUpdateCollection(t *testing.T) {
	t.Run("SuccessUpdateCollection", func(t *testing.T) {
		retErr := bl.CreateError(bl.AllIsOk, nil, "")
		retCollection := &models.Collection{Id: 1, Name: "Collection 1"}

		mockCollectionRepo := new(mocks.MockICollectionRepository)
		mockCollectionRepo.On("UpdateCollection", retCollection).Return(retErr)

		csSrv := bl.CollectionService{}
		err := csSrv.UpdateCollection(retCollection, mockCollectionRepo)

		assert.NotNil(t, err)
		assert.Equal(t, err.ErrNum, bl.AllIsOk)

		mockCollectionRepo.AssertExpectations(t)
	})
}

func TestGetAllNotesInCollection(t *testing.T) {
	t.Run("SuccessGetAllNotesInCollection", func(t *testing.T) {
		retErr := bl.CreateError(bl.AllIsOk, nil, "")
		retNotes := []*models.Note{
			{Id: 1, OwnerID: 1},
			{Id: 2, OwnerID: 2},
		}
		retCollection := &models.Collection{Id: 1}

		mockCollectionRepo := new(mocks.MockICollectionRepository)
		mockCollectionRepo.On("GetAllNotesInCollection", retCollection).Return(retNotes, retErr)

		csSrv := bl.CollectionService{}
		notes, err := csSrv.GetAllNotesInCollection(retCollection, mockCollectionRepo)

		assert.NotNil(t, err)
		assert.Equal(t, err.ErrNum, bl.AllIsOk)
		assert.Len(t, notes, 2)
		assert.Equal(t, 1, notes[0].Id)
		assert.Equal(t, 2, notes[1].Id)

		mockCollectionRepo.AssertExpectations(t)
	})
}
