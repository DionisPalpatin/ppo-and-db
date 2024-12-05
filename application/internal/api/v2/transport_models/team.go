package transport_models

type TeamFullInfo struct {
	ID               int    `json:"id" validate:"required"`
	Name             string `json:"name" validate:"required"`
	RegistrationDate string `json:"registration_date" validate:"required"`
}

type TeamInfo struct {
	ID   int    `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}
