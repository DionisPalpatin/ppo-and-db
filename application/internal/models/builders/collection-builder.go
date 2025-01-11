package builders

import (
	"time"

	"github.com/DionisPalpatin/ppo-and-db/application/internal/models"
)

type CollectionBuilder struct {
	Collection *models.Collection
}

func NewCollectionBuilder() *CollectionBuilder {
	return &CollectionBuilder{
		Collection: &models.Collection{
			Id:           1,
			Name:         "Default Collection",
			CreationDate: time.Now(),
			OwnerID:      1,
		},
	}
}

func (b *CollectionBuilder) WithId(id int) *CollectionBuilder {
	b.Collection.Id = id
	return b
}

func (b *CollectionBuilder) WithName(name string) *CollectionBuilder {
	b.Collection.Name = name
	return b
}

func (b *CollectionBuilder) WithCreationDate(date time.Time) *CollectionBuilder {
	b.Collection.CreationDate = date
	return b
}

func (b *CollectionBuilder) WithOwnerID(ownerID int) *CollectionBuilder {
	b.Collection.OwnerID = ownerID
	return b
}

func (b *CollectionBuilder) Build() *models.Collection {
	return b.Collection
}
