package bl

import "notebook_app/internal/models"

type CollectionService struct{}

func (CollectionService) GetCollection(colID int, name string, searchBy int, icr ICollectionRepository) (*models.Collection, *MyError) {
	switch searchBy {
	case SearchByID:
		return icr.GetCollectionByID(colID)

	case SearchByString:
		return icr.GetCollectionByName(name)

	default:
		return nil, CreateError(ErrSearchParameter, ErrSearchParameterError(), "GetCollection")
	}
}

func (CollectionService) GetAllCollections(user *models.User, icr ICollectionRepository) ([]*models.Collection, *MyError) {
	if user.Role != Admin {
		return nil, CreateError(ErrAccessDenied, ErrAccessDeniedError(), "GetAllCollections")
	}
	return icr.GetAllCollections()
}

func (CollectionService) GetAllUsersCollections(user *models.User, icr ICollectionRepository) ([]*models.Collection, *MyError) {
	return icr.GetAllUserCollections(user)
}

func (CollectionService) AddCollection(coll *models.Collection, icr ICollectionRepository) *MyError {
	return icr.AddCollection(coll)
}

func (CollectionService) DeleteCollection(id int, user *models.User, icr ICollectionRepository) *MyError {
	col, myErr := icr.GetCollectionByID(id)
	if myErr.ErrNum != AllIsOk {
		return myErr
	}

	if col.OwnerID != user.Id && user.Role != Admin {
		myErr = CreateError(ErrDeleteCollection, ErrDeleteCollectionError(), "DeleteColleciton")
		return myErr
	}

	return icr.DeleteCollection(id)
}

func (CollectionService) UpdateCollection(collection *models.Collection, icr ICollectionRepository) *MyError {
	return icr.UpdateCollection(collection)
}

func (CollectionService) GetAllNotesInCollection(collection *models.Collection, icr ICollectionRepository) ([]*models.Note, *MyError) {
	return icr.GetAllNotesInCollection(collection)
}
