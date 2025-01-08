package bl

import "github.com/DionisPalpatin/ppo-and-db/application/internal/models"

type CollectionService struct {
	icr ICollectionRepository
}

func (cs *CollectionService) GetCollection(colID int, name string, searchBy int) (*models.Collection, *MyError) {
	switch searchBy {
	case SearchByID:
		return cs.icr.GetCollectionByID(colID)

	case SearchByString:
		return cs.icr.GetCollectionByName(name)

	default:
		return nil, CreateError(ErrSearchParameter, "GetCollection", "bl")
	}
}

func (cs *CollectionService) GetAllCollections(user *models.User) ([]*models.Collection, *MyError) {
	if user.Role != Admin {
		return nil, CreateError(ErrAccessDenied, "GetAllCollections", "bl")
	}
	return cs.icr.GetAllCollections()
}

func (cs *CollectionService) GetAllUsersCollections(user *models.User) ([]*models.Collection, *MyError) {
	return cs.icr.GetAllUserCollections(user)
}

func (cs *CollectionService) AddCollection(coll *models.Collection) (int, *MyError) {
	return cs.icr.AddCollection(coll)
}

func (cs *CollectionService) DeleteCollection(id int, user *models.User) *MyError {
	col, myErr := cs.icr.GetCollectionByID(id)
	if myErr.ErrNum != Ok {
		return myErr
	}

	if col.OwnerID != user.Id && user.Role != Admin {
		myErr = CreateError(OperationError, "DeleteCollection", "bl")
		return myErr
	}

	return cs.icr.DeleteCollection(id)
}

func (cs *CollectionService) UpdateCollection(collection *models.Collection) *MyError {
	return cs.icr.UpdateCollection(collection)
}

func (cs *CollectionService) GetAllNotesInCollection(collection *models.Collection) ([]*models.Note, *MyError) {
	return cs.icr.GetAllNotesInCollection(collection)
}
