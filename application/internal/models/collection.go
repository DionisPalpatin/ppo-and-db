package models

import "time"

type Collection struct {
	Id           int
	Name         string
	CreationDate time.Time
	OwnerID      int
}
