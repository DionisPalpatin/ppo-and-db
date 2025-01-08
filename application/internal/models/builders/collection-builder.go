package builders

import (
	"time"

	"github.com/DionisPalpatin/ppo-and-db/application/internal/models"
)

type CollectionBuilder struct {
	collection *models.Collection
}

func NewCollectionBuilder() *CollectionBuilder {
	return &CollectionBuilder{
		collection: &models.Collection{
			Id:           1,
			Name:         "Default Collection",
			CreationDate: time.Now(),
			OwnerID:      1,
		},
	}
}

func (b *CollectionBuilder) WithId(id int) *CollectionBuilder {
	b.collection.Id = id
	return b
}

func (b *CollectionBuilder) WithName(name string) *CollectionBuilder {
	b.collection.Name = name
	return b
}

func (b *CollectionBuilder) WithCreationDate(date time.Time) *CollectionBuilder {
	b.collection.CreationDate = date
	return b
}

func (b *CollectionBuilder) WithOwnerID(ownerID int) *CollectionBuilder {
	b.collection.OwnerID = ownerID
	return b
}

func (b *CollectionBuilder) Build() *models.Collection {
	return b.collection
}
