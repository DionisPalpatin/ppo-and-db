package converters

import (
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/api/v2/transport_models"
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/models"
	"time"
)

func ToCollectionInfo(collection *models.Collection) transport_models.CollectionInfo {
	return transport_models.CollectionInfo{
		ID:   collection.Id,
		Name: collection.Name,
	}
}

func ToCollectionFullInfo(collection *models.Collection) transport_models.CollectionFullInfo {
	return transport_models.CollectionFullInfo{
		ID:               collection.Id,
		Name:             collection.Name,
		RegistrationDate: collection.CreationDate.String(),
		OwnerID:          collection.OwnerID,
	}
}

func FromCollectionInfo(collectionInfo transport_models.CollectionInfo) models.Collection {
	return models.Collection{
		Id:   collectionInfo.ID,
		Name: collectionInfo.Name,
	}
}

func FromCollectionFullInfo(collectionFullInfo transport_models.CollectionFullInfo) (models.Collection, error) {
	parsedTime, err := time.Parse(collectionFullInfo.RegistrationDate, "2006-01-02 15:04:05-07")
	if err != nil {
		return models.Collection{}, err
	}

	return models.Collection{
		Id:           collectionFullInfo.ID,
		Name:         collectionFullInfo.Name,
		CreationDate: parsedTime,
		OwnerID:      collectionFullInfo.OwnerID,
	}, nil
}
