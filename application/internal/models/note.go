package models

import "time"

type Note struct {
	Id               int
	Access           int
	Name             string
	ContentType      int
	Likes            int
	Dislikes         int
	RegistrationDate time.Time
	OwnerID          int
	SectionID        int
}
