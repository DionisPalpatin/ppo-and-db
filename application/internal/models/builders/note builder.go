package builders

import "github.com/DionisPalpatin/ppo-and-db/application/internal/models"
import "time"


type ContentBuilder struct {
	content *models.Content
}


type NoteBuilder struct {
	note *models.Note
}

func NewContentBuilder() *ContentBuilder {
	return &ContentBuilder{
		content: &models.Content{
			Text:    []byte("Default text"),
			Img:     []byte{},
			Raw:     []byte{},
			TextExt: ".txt",
			ImgExt:  ".jpg",
			RawExt:  ".bin",
		},
	}
}

func NewNoteBuilder() *NoteBuilder {
	return &NoteBuilder{
		note: &models.Note{
			Id:               1,
			Access:           0,
			Name:             "Default Note",
			Content:          *NewContentBuilder().Build(),
			Likes:            0,
			Dislikes:         0,
			RegistrationDate: time.Now(),
			OwnerID:          1,
			SectionID:        1,
		},
	}
}

func (b *ContentBuilder) WithText(text []byte) *ContentBuilder {
	b.content.Text = text
	return b
}

func (b *ContentBuilder) WithImg(img []byte) *ContentBuilder {
	b.content.Img = img
	return b
}

func (b *ContentBuilder) WithRaw(raw []byte) *ContentBuilder {
	b.content.Raw = raw
	return b
}

func (b *ContentBuilder) WithTextExt(ext string) *ContentBuilder {
	b.content.TextExt = ext
	return b
}

func (b *ContentBuilder) WithImgExt(ext string) *ContentBuilder {
	b.content.ImgExt = ext
	return b
}

func (b *ContentBuilder) WithRawExt(ext string) *ContentBuilder {
	b.content.RawExt = ext
	return b
}

func (b *ContentBuilder) Build() *models.Content {
	return b.content
}

func (b *NoteBuilder) WithId(id int) *NoteBuilder {
	b.note.Id = id
	return b
}

func (b *NoteBuilder) WithAccess(access int) *NoteBuilder {
	b.note.Access = access
	return b
}

func (b *NoteBuilder) WithName(name string) *NoteBuilder {
	b.note.Name = name
	return b
}

func (b *NoteBuilder) WithContent(content models.Content) *NoteBuilder {
	b.note.Content = content
	return b
}

func (b *NoteBuilder) WithLikes(likes int) *NoteBuilder {
	b.note.Likes = likes
	return b
}

func (b *NoteBuilder) WithDislikes(dislikes int) *NoteBuilder {
	b.note.Dislikes = dislikes
	return b
}

func (b *NoteBuilder) WithRegistrationDate(date time.Time) *NoteBuilder {
	b.note.RegistrationDate = date
	return b
}

func (b *NoteBuilder) WithOwnerID(ownerID int) *NoteBuilder {
	b.note.OwnerID = ownerID
	return b
}

func (b *NoteBuilder) WithSectionID(sectionID int) *NoteBuilder {
	b.note.SectionID = sectionID
	return b
}

func (b *NoteBuilder) Build() *models.Note {
	return b.note
}
