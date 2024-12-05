package transport_models

type NoteInfo struct {
	ID       int    `json:"id" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Likes    int    `json:"likes" validate:"required"`
	Dislikes int    `json:"dislikes" validate:"required"`
}

type NoteFullInfo struct {
	ID               int    `json:"id" validate:"required"`
	Access           int    `json:"access" validate:"required"`
	Name             string `json:"name" validate:"required"`
	Likes            int    `json:"likes" validate:"required"`
	Dislikes         int    `json:"dislikes" validate:"required"`
	RegistrationDate string `json:"registration_date" validate:"required"`
	OwnerID          int    `json:"owner_id" validate:"required"`
	SectionID        int    `json:"section_id" validate:"required"`
}

type NoteFullData struct {
	Info      NoteFullInfo `json:"info"`
	TextFile  []byte       `json:"text_file"`
	Image     []byte       `json:"image"`
	OtherFile []byte       `json:"other_file"`
	FileType  string       `json:"file_type"`
}
