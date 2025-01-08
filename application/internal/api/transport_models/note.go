package transport_models

type NoteInfo struct {
	ID       int    `my_json:"id" validate:"required"`
	Name     string `my_json:"name" validate:"required"`
	Likes    int    `my_json:"likes" validate:"required"`
	Dislikes int    `my_json:"dislikes" validate:"required"`
}

type NoteFullInfo struct {
	ID               int    `my_json:"id" validate:"required"`
	Access           int    `my_json:"access" validate:"required"`
	Name             string `my_json:"name" validate:"required"`
	Likes            int    `my_json:"likes" validate:"required"`
	Dislikes         int    `my_json:"dislikes" validate:"required"`
	RegistrationDate string `my_json:"registration_date" validate:"required"`
	OwnerID          int    `my_json:"owner_id" validate:"required"`
	SectionID        int    `my_json:"section_id" validate:"required"`
}

type NoteFullData struct {
	Info          NoteFullInfo `my_json:"info"`
	TextFile      []byte       `my_json:"text_file"`
	Image         []byte       `my_json:"image"`
	OtherFile     []byte       `my_json:"other_file"`
	TxtFileType   string       `my_json:"txt_file_type"`
	ImgFileType   string       `my_json:"img_file_type"`
	OtherFileType string       `my_json:"other_file_type"`
}
