package models

import "time"

type Content struct {
	Text []byte
	Img  []byte
	Raw  []byte

	TextExt string
	ImgExt  string
	RawExt  string
}

type Note struct {
	Id               int
	Access           int
	Name             string
	Content          Content
	Likes            int
	Dislikes         int
	RegistrationDate time.Time
	OwnerID          int
	SectionID        int
}
