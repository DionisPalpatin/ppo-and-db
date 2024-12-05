package transport_models

type CollectionInfo struct {
	ID   int    `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}

type CollectionAddInfo struct {
	ID      int    `json:"id" validate:"required"`
	Name    string `json:"name" validate:"required"`
	OwnerID int    `json:"owner_id" validate:"required"`
}

type CollectionFullInfo struct {
	ID               int    `json:"id" validate:"required"`
	Name             string `json:"name" validate:"required"`
	RegistrationDate string `json:"registration_date" validate:"required"`
	OwnerID          int    `json:"owner_id" validate:"required"`
}
