package models

type Note struct {
	Id               int
	Access           int
	Name             string
	ContentType      int
	Likes            int
	Dislikes         int
	RegistrationDate string
	OwnerID          int
	SectionID        int
}
