package builders

import "github.com/DionisPalpatin/ppo-and-db/application/internal/models"
import "time"


type ContentBuilder struct {
	Content *models.Content
}


type NoteBuilder struct {
	Note *models.Note
}

func NewContentBuilder() *ContentBuilder {
	return &ContentBuilder{
		Content: &models.Content{
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
		Note: &models.Note{
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
	b.Content.Text = text
	return b
}

func (b *ContentBuilder) WithImg(img []byte) *ContentBuilder {
	b.Content.Img = img
	return b
}

func (b *ContentBuilder) WithRaw(raw []byte) *ContentBuilder {
	b.Content.Raw = raw
	return b
}

func (b *ContentBuilder) WithTextExt(ext string) *ContentBuilder {
	b.Content.TextExt = ext
	return b
}

func (b *ContentBuilder) WithImgExt(ext string) *ContentBuilder {
	b.Content.ImgExt = ext
	return b
}

func (b *ContentBuilder) WithRawExt(ext string) *ContentBuilder {
	b.Content.RawExt = ext
	return b
}

func (b *ContentBuilder) Build() *models.Content {
	return b.Content
}

func (b *NoteBuilder) WithId(id int) *NoteBuilder {
	b.Note.Id = id
	return b
}

func (b *NoteBuilder) WithAccess(access int) *NoteBuilder {
	b.Note.Access = access
	return b
}

func (b *NoteBuilder) WithName(name string) *NoteBuilder {
	b.Note.Name = name
	return b
}

func (b *NoteBuilder) WithContent(content models.Content) *NoteBuilder {
	b.Note.Content = content
	return b
}

func (b *NoteBuilder) WithLikes(likes int) *NoteBuilder {
	b.Note.Likes = likes
	return b
}

func (b *NoteBuilder) WithDislikes(dislikes int) *NoteBuilder {
	b.Note.Dislikes = dislikes
	return b
}

func (b *NoteBuilder) WithRegistrationDate(date time.Time) *NoteBuilder {
	b.Note.RegistrationDate = date
	return b
}

func (b *NoteBuilder) WithOwnerID(ownerID int) *NoteBuilder {
	b.Note.OwnerID = ownerID
	return b
}

func (b *NoteBuilder) WithSectionID(sectionID int) *NoteBuilder {
	b.Note.SectionID = sectionID
	return b
}

func (b *NoteBuilder) Build() *models.Note {
	return b.Note
}
