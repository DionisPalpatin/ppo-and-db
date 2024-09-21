package models

type User struct {
	Id               int
	Fio              string
	RegistrationDate string
	Login            string
	Password         string
	Role             int
}
