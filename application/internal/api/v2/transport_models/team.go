package transport_models

type TeamFullInfo struct {
	ID               int    `my_json:"id" validate:"required"`
	Name             string `my_json:"name" validate:"required"`
	RegistrationDate string `my_json:"registration_date" validate:"required"`
}

type TeamInfo struct {
	ID   int    `my_json:"id" validate:"required"`
	Name string `my_json:"name" validate:"required"`
}
