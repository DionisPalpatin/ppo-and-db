package transport_models

type SectionInfo struct {
	ID     int `json:"id" validate:"required"`
	TeamID int `json:"team_id" validate:"required"`
}

type SectionFullInfo struct {
	ID               int    `json:"id" validate:"required"`
	RegistrationDate string `json:"registration_date" validate:"required"`
	TeamID           int    `json:"team_id" validate:"required"`
}
