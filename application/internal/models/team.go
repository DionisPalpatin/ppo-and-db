package models

import "time"

type Team struct {
	Id               int
	Name             string
	RegistrationDate time.Time
}
