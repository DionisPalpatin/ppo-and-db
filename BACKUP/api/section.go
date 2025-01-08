package transport_models

type SectionInfo struct {
	ID     int `my_json:"id" validate:"required"`
	TeamID int `my_json:"team_id" validate:"required"`
}

type SectionFullInfo struct {
	ID               int    `my_json:"id" validate:"required"`
	RegistrationDate string `my_json:"registration_date" validate:"required"`
	TeamID           int    `my_json:"team_id" validate:"required"`
}
