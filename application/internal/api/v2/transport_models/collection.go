package transport_models

type CollectionInfo struct {
	ID   int    `my_json:"id" validate:"required"`
	Name string `my_json:"name" validate:"required"`
}

type CollectionAddInfo struct {
	ID      int    `my_json:"id" validate:"required"`
	Name    string `my_json:"name" validate:"required"`
	OwnerID int    `my_json:"owner_id" validate:"required"`
}

type CollectionFullInfo struct {
	ID               int    `my_json:"id" validate:"required"`
	Name             string `my_json:"name" validate:"required"`
	RegistrationDate string `my_json:"registration_date" validate:"required"`
	OwnerID          int    `my_json:"owner_id" validate:"required"`
}
