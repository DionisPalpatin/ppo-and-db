package transport_models

type UserRegistrationInfo struct {
	FIO      string `my_json:"fio" validate:"required"`
	Login    string `my_json:"login" validate:"required"`
	Password string `my_json:"password" validate:"required"`
}

type UserLoginInfo struct {
	Login    string `my_json:"login" validate:"required"`
	Password string `my_json:"password" validate:"required"`
}

type UserPublicInfo struct {
	ID   int    `my_json:"id" validate:"required"`
	FIO  string `my_json:"fio" validate:"required"`
	Role int    `my_json:"role" validate:"required"`
}

type UserPrivateInfo struct {
	ID       int    `my_json:"id" validate:"required"`
	FIO      string `my_json:"fio" validate:"required"`
	Login    string `my_json:"login" validate:"required"`
	Password string `my_json:"password" validate:"required"`
	Role     int    `my_json:"role" validate:"required"`
}

type UserFullInfo struct {
	ID               int    `my_json:"id" validate:"required"`
	FIO              string `my_json:"fio" validate:"required"`
	RegistrationDate string `my_json:"registration_date" validate:"required"`
	Login            string `my_json:"login" validate:"required"`
	Password         string `my_json:"password" validate:"required"`
	Role             int    `my_json:"role" validate:"required"`
}
