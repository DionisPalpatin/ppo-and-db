package converters

import (
	"time"

	"github.com/DionisPalpatin/ppo-and-db/application/internal/api/transport_models"
	"github.com/DionisPalpatin/ppo-and-db/application/internal/models"
)

func ToCollectionInfo(collection *models.Collection) transport_models.CollectionInfo {
	return transport_models.CollectionInfo{
		ID:   collection.Id,
		Name: collection.Name,
	}
}

func ToCollectionAddInfo(collection *models.Collection) transport_models.CollectionAddInfo {
	return transport_models.CollectionAddInfo{
		ID:      collection.Id,
		Name:    collection.Name,
		OwnerID: collection.OwnerID,
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

func FromCollectionInfo(collectionInfo *transport_models.CollectionInfo) models.Collection {
	return models.Collection{
		Id:   collectionInfo.ID,
		Name: collectionInfo.Name,
	}
}

func FromCollectionAddInfo(collectionInfo *transport_models.CollectionAddInfo) models.Collection {
	return models.Collection{
		Id:      collectionInfo.ID,
		Name:    collectionInfo.Name,
		OwnerID: collectionInfo.OwnerID,
	}
}

func FromCollectionFullInfo(collectionFullInfo *transport_models.CollectionFullInfo) (models.Collection, error) {
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
