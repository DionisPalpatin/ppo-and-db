package models

import "time"

type User struct {
	Id               int
	Fio              string
	RegistrationDate time.Time
	Login            string
	Password         string
	Role             int
}
