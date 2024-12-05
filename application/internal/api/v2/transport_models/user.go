package transport_models

type UserRegistrationInfo struct {
	FIO      string `json:"fio" validate:"required"`
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserLoginInfo struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserPublicInfo struct {
	ID   int    `json:"id" validate:"required"`
	FIO  string `json:"fio" validate:"required"`
	Role int    `json:"role" validate:"required"`
}

type UserPrivateInfo struct {
	ID       int    `json:"id" validate:"required"`
	FIO      string `json:"fio" validate:"required"`
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
	Role     int    `json:"role" validate:"required"`
}

type UserFullInfo struct {
	ID               int    `json:"id" validate:"required"`
	FIO              string `json:"fio" validate:"required"`
	RegistrationDate string `json:"registration_date" validate:"required"`
	Login            string `json:"login" validate:"required"`
	Password         string `json:"password" validate:"required"`
	Role             int    `json:"role" validate:"required"`
}
