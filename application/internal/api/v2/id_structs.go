package handlersv2

type userIDStruct struct {
	UserID int `json:"user_id" validate:"required"`
}

type teamIDStruct struct {
	TeamID int `json:"team_id" validate:"required"`
}

type noteIDStruct struct {
	NoteID int `json:"note_id" validate:"required"`
}

type collectionIDStruct struct {
	CollID int `json:"collection_id" validate:"required"`
}

type sectionIDStruct struct {
	SecID int `json:"section_id" validate:"required"`
}
