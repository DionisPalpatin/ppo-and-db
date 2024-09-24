package clickhouse

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log/slog"

	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/business_logic"
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/models"
)

func (nr *NoteRepository) GetNoteByID(id int) (*models.Note, []byte, *bl.MyError) {
	nr.MyLogger.WriteLog("GetNoteByID is called (Repo)", slog.LevelInfo, nil)

	if id == 0 {
		myErr := bl.CreateError(bl.ErrGetNoteByID, bl.ErrGetNoteByIDError(), "GetNoteByID")
		return nil, nil, myErr
	}

	var note models.Note
	db := nr.DbConfigs.DB
	schemaName := nr.DbConfigs.SchemaName
	query := fmt.Sprintf("SELECT * FROM %s.notes WHERE id = ?", schemaName)
	ctx := context.Background()

	err := db.QueryRowContext(ctx, query, id).Scan(
		&note.Id,
		&note.Access,
		&note.Name,
		&note.ContentType,
		&note.Likes,
		&note.Dislikes,
		&note.RegistrationDate,
		&note.OwnerID,
		&note.SectionID,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			myErr := bl.CreateError(bl.ErrGetNoteByID, bl.ErrGetNoteByIDError(), "GetNoteByID")
			return nil, nil, myErr
		}
		return nil, nil, bl.CreateError(bl.ErrGetNoteByID, err, "GetNoteByID")
	}

	if note.ContentType == bl.TextCont {
		query = fmt.Sprintf("SELECT * FROM %s.texts WHERE note_id = $1", schemaName)
	} else if note.ContentType == bl.ImgCont {
		query = fmt.Sprintf("SELECT * FROM %s.images WHERE note_id = $1", schemaName)
	} else {
		query = fmt.Sprintf("SELECT * FROM %s.raw_datas WHERE note_id = $1", schemaName)
	}

	var dataId, noteId int
	var data []byte

	err = db.QueryRowContext(ctx, query, note.Id).Scan(&dataId, &data, &noteId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			myErr := bl.CreateError(bl.ErrGetNoteByName, bl.ErrGetNoteByNameError(), "GetNoteByName")
			return nil, nil, myErr
		}

		myErr := bl.CreateError(bl.ErrGetNoteByName, err, "GetNoteByName")
		return nil, nil, myErr
	}

	myOk := bl.CreateError(bl.AllIsOk, nil, "")
	return &note, data, myOk
}

func (nr *NoteRepository) GetNoteByName(name string) (*models.Note, []byte, *bl.MyError) {
	nr.MyLogger.WriteLog("GetNoteByName is called (Repo)", slog.LevelInfo, nil)

	if name == "" {
		myErr := bl.CreateError(bl.ErrGetNoteByName, bl.ErrGetNoteByNameError(), "GetNoteByName")
		return nil, nil, myErr
	}

	var note models.Note
	db := nr.DbConfigs.DB
	schemaName := nr.DbConfigs.SchemaName
	query := fmt.Sprintf("SELECT * FROM %s.notes WHERE name = ?", schemaName)
	ctx := context.Background()

	err := db.QueryRowContext(ctx, query, name).Scan(
		&note.Id,
		&note.Access,
		&note.Name,
		&note.ContentType,
		&note.Likes,
		&note.Dislikes,
		&note.RegistrationDate,
		&note.OwnerID,
		&note.SectionID,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			myErr := bl.CreateError(bl.ErrGetNoteByName, bl.ErrGetNoteByNameError(), "GetNoteByName")
			return nil, nil, myErr
		}

		myErr := bl.CreateError(bl.ErrGetNoteByName, err, "GetNoteByName")
		return nil, nil, myErr
	}

	if note.ContentType == bl.TextCont {
		query = fmt.Sprintf("SELECT * FROM %s.texts WHERE note_id = $1", schemaName)
	} else if note.ContentType == bl.ImgCont {
		query = fmt.Sprintf("SELECT * FROM %s.images WHERE note_id = $1", schemaName)
	} else {
		query = fmt.Sprintf("SELECT * FROM %s.raw_datas WHERE note_id = $1", schemaName)
	}

	var id, noteId int
	var data []byte

	err = db.QueryRowContext(ctx, query, note.Id).Scan(&id, &data, &noteId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			myErr := bl.CreateError(bl.ErrGetNoteByName, bl.ErrGetNoteByNameError(), "GetNoteByName")
			return nil, nil, myErr
		}

		myErr := bl.CreateError(bl.ErrGetNoteByName, err, "GetNoteByName")
		return nil, nil, myErr
	}

	myOk := bl.CreateError(bl.AllIsOk, nil, "")
	return &note, data, myOk
}

func (nr *NoteRepository) GetAllNotes() ([]*models.Note, *bl.MyError) {
	nr.MyLogger.WriteLog("GetAllNotes is called (Repo)", slog.LevelInfo, nil)

	db := nr.DbConfigs.DB
	schemaName := nr.DbConfigs.SchemaName
	query := fmt.Sprintf("SELECT * FROM %s.notes", schemaName)
	ctx := context.Background()

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		myErr := bl.CreateError(bl.ErrGetAllNotes, err, "GetAllNotes")
		nr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return nil, myErr
	}
	defer rows.Close()

	var notes []*models.Note
	for rows.Next() {
		var note models.Note
		err := rows.Scan(
			&note.Id,
			&note.Access,
			&note.Name,
			&note.ContentType,
			&note.Likes,
			&note.Dislikes,
			&note.RegistrationDate,
			&note.OwnerID,
			&note.SectionID,
		)

		if err != nil {
			myErr := bl.CreateError(bl.ErrGetAllNotes, err, "GetAllNotes")
			nr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
			return nil, myErr
		}

		notes = append(notes, &note)
	}

	if err := rows.Err(); err != nil {
		myErr := bl.CreateError(bl.ErrGetAllNotes, err, "GetAllNotes")
		nr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return nil, myErr
	}

	myOk := bl.CreateError(bl.AllIsOk, nil, "")
	return notes, myOk
}

func (nr *NoteRepository) GetAllPublicNotes() ([]*models.Note, *bl.MyError) {
	nr.MyLogger.WriteLog("GetAllNotes is called (Repo)", slog.LevelInfo, nil)

	db := nr.DbConfigs.DB
	schemaName := nr.DbConfigs.SchemaName
	query := fmt.Sprintf("SELECT * FROM %s.notes WHERE access = 1", schemaName)
	ctx := context.Background()

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		myErr := bl.CreateError(bl.ErrGetAllNotes, err, "GetAllNotes")
		nr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return nil, myErr
	}
	defer rows.Close()

	var notes []*models.Note
	for rows.Next() {
		var note models.Note
		err := rows.Scan(
			&note.Id,
			&note.Access,
			&note.Name,
			&note.ContentType,
			&note.Likes,
			&note.Dislikes,
			&note.RegistrationDate,
			&note.OwnerID,
			&note.SectionID,
		)

		if err != nil {
			myErr := bl.CreateError(bl.ErrGetAllNotes, err, "GetAllNotes")
			nr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
			return nil, myErr
		}

		notes = append(notes, &note)
	}

	if err := rows.Err(); err != nil {
		myErr := bl.CreateError(bl.ErrGetAllNotes, err, "GetAllNotes")
		nr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return nil, myErr
	}

	myOk := bl.CreateError(bl.AllIsOk, nil, "")
	return notes, myOk
}

func (nr *NoteRepository) AddNote(note *models.Note) *bl.MyError {
	nr.MyLogger.WriteLog("AddNote is called (Repo)", slog.LevelInfo, nil)

	if note == nil {
		myErr := bl.CreateError(bl.ErrAddNote, bl.ErrAddNoteError(), "AddNote")
		nr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}

	db := nr.DbConfigs.DB
	schemaName := nr.DbConfigs.SchemaName
	query := fmt.Sprintf("INSERT INTO %s.notes (access, name, content_type, likes, dislikes, registration_date, owner_id, section_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", schemaName)
	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		myErr := bl.CreateError(bl.ErrAddNote, err, "AddNote")
		nr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	_, err = tx.ExecContext(ctx, query,
		note.Access,
		note.Name,
		note.ContentType,
		note.Likes,
		note.Dislikes,
		note.RegistrationDate,
		note.OwnerID,
		note.SectionID,
	)

	if err != nil {
		myErr := bl.CreateError(bl.ErrAddNote, err, "AddNote")
		nr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}

	myOk := bl.CreateError(bl.AllIsOk, nil, "")
	return myOk
}

func (nr *NoteRepository) DeleteNote(id int) *bl.MyError {
	nr.MyLogger.WriteLog("DeleteNote is called (Repo)", slog.LevelInfo, nil)

	if id == 0 {
		myErr := bl.CreateError(bl.ErrDeleteNote, bl.ErrDeleteNoteError(), "DeleteNote")
		nr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}

	db := nr.DbConfigs.DB
	schemaName := nr.DbConfigs.Name
	ctx := context.Background()

	query1 := fmt.Sprintf("ALTER TABLE %s.texts DELETE WHERE note_id = ?", schemaName)
	query2 := fmt.Sprintf("ALTER TABLE %s.images DELETE WHERE note_id = ?", schemaName)
	query3 := fmt.Sprintf("ALTER TABLE %s.raw_datas DELETE WHERE note_id = ?", schemaName)
	query4 := fmt.Sprintf("ALTER TABLE %s.notes_collections DELETE WHERE note_id = ?", schemaName)
	query5 := fmt.Sprintf("ALTER TABLE %s.notes DELETE WHERE id = ?", schemaName)
	result_query := fmt.Sprintf("%s; %s; %s; %s; %s;", query1, query2, query3, query4, query5)

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		myErr := bl.CreateError(bl.ErrDeleteNote, err, "DeleteNote")
		nr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	_, err = tx.ExecContext(ctx, result_query, id)

	if err != nil {
		myErr := bl.CreateError(bl.ErrDeleteNote, err, "DeleteNote")
		nr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}

	myOk := bl.CreateError(bl.AllIsOk, nil, "")
	return myOk
}

func (nr *NoteRepository) UpdateNoteContentText(reader io.Reader, note *models.Note) *bl.MyError {
	nr.MyLogger.WriteLog("UpdateNoteContentText is called (Repo)", slog.LevelInfo, nil)

	if note == nil {
		myErr := bl.CreateError(bl.ErrUpdateNoteContent, bl.ErrUpdateNoteContentError(), "UpdateNoteContentText")
		nr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}

	db := nr.DbConfigs.DB
	schemaName := nr.DbConfigs.Name
	query := fmt.Sprintf("INSERT INTO %s.texts (data, note_id) VALUES (?, ?)", schemaName)
	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		myErr := bl.CreateError(bl.ErrUpdateNoteContent, err, "UpdateNoteContentText")
		nr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	data, err := io.ReadAll(reader)
	if err != nil {
		myErr := bl.CreateError(bl.ErrUpdateNoteContent, err, "UpdateNoteContentText")
		nr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}

	_, err = tx.ExecContext(ctx, query, data, note.Id)
	if err != nil {
		myErr := bl.CreateError(bl.ErrUpdateNoteContent, err, "UpdateNoteContentText")
		nr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}

	myOk := bl.CreateError(bl.AllIsOk, nil, "")
	return myOk
}

func (nr *NoteRepository) UpdateNoteContentImg(reader io.Reader, note *models.Note) *bl.MyError {
	nr.MyLogger.WriteLog("UpdateNoteContentImg is called (Repo)", slog.LevelInfo, nil)

	if note == nil {
		myErr := bl.CreateError(bl.ErrUpdateNoteContent, bl.ErrUpdateNoteContentError(), "UpdateNoteContentImg")
		nr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}

	db := nr.DbConfigs.DB
	schemaName := nr.DbConfigs.Name
	query := fmt.Sprintf("INSERT INTO %s.images (data, note_id) VALUES (?, ?)", schemaName)
	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		myErr := bl.CreateError(bl.ErrUpdateNoteContent, err, "UpdateNoteContentImg")
		nr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	data, err := io.ReadAll(reader)
	if err != nil {
		myErr := bl.CreateError(bl.ErrUpdateNoteContent, err, "UpdateNoteContentImg")
		nr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}

	_, err = tx.ExecContext(ctx, query, data, note.Id)
	if err != nil {
		myErr := bl.CreateError(bl.ErrUpdateNoteContent, err, "UpdateNoteContentImg")
		nr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}

	myOk := bl.CreateError(bl.AllIsOk, nil, "")
	return myOk
}

func (nr *NoteRepository) UpdateNoteContentRawData(reader io.Reader, note *models.Note) *bl.MyError {
	nr.MyLogger.WriteLog("UpdateNoteContentRawData is called (Repo)", slog.LevelInfo, nil)

	if note == nil {
		myErr := bl.CreateError(bl.ErrUpdateNoteContent, bl.ErrUpdateNoteContentError(), "UpdateNoteContentRawData")
		nr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}

	db := nr.DbConfigs.DB
	schemaName := nr.DbConfigs.Name
	query := fmt.Sprintf("INSERT INTO %s.raw_datas (data, note_id) VALUES (?, ?)", schemaName)
	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		myErr := bl.CreateError(bl.ErrUpdateNoteContent, err, "UpdateNoteContentRawData")
		nr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	data, err := io.ReadAll(reader)
	if err != nil {
		myErr := bl.CreateError(bl.ErrUpdateNoteContent, err, "UpdateNoteContentRawData")
		nr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}

	_, err = tx.ExecContext(ctx, query, data, note.Id)

	if err != nil {
		myErr := bl.CreateError(bl.ErrUpdateNoteContent, err, "UpdateNoteContentRawData")
		nr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}

	myOk := bl.CreateError(bl.AllIsOk, nil, "")
	return myOk
}

func (nr *NoteRepository) UpdateNoteInfo(note *models.Note) *bl.MyError {
	nr.MyLogger.WriteLog("UpdateNoteInfo is called (Repo)", slog.LevelInfo, nil)

	if note == nil {
		myErr := bl.CreateError(bl.ErrUpdateNoteInfo, bl.ErrUpdateNoteInfoError(), "UpdateNoteInfo")
		nr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}

	db := nr.DbConfigs.DB
	schemaName := nr.DbConfigs.Name
	query := fmt.Sprintf("ALTER TABLE %s.notes UPDATE access = ?, name = ?, content_type = ?, likes = ?, dislikes = ?, registration_date = ?, owner_id = ?, section_id = ? WHERE id = ?", schemaName)
	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		myErr := bl.CreateError(bl.ErrUpdateNoteInfo, err, "UpdateNoteInfo")
		nr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	_, err = tx.ExecContext(ctx, query,
		note.Access,
		note.Name,
		note.ContentType,
		note.Likes,
		note.Dislikes,
		note.RegistrationDate,
		note.OwnerID,
		note.SectionID,
		note.Id,
	)

	if err != nil {
		myErr := bl.CreateError(bl.ErrUpdateNoteInfo, err, "UpdateNoteInfo")
		nr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}

	myOk := bl.CreateError(bl.AllIsOk, nil, "")
	return myOk
}

func (nr *NoteRepository) AddNoteToCollection(collectionID int, noteID int) *bl.MyError {
	nr.MyLogger.WriteLog("AddNoteToCollection is called (Repo)", slog.LevelInfo, nil)

	if collectionID < 0 || noteID < 0 {
		myErr := bl.CreateError(bl.ErrAddNoteToCollection, bl.ErrAddNoteToCollectionError(), "AddNoteToCollection")
		nr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}

	db := nr.DbConfigs.DB
	schemaName := nr.DbConfigs.Name
	query := fmt.Sprintf("INSERT INTO %s.note_collections (note_id, collection_id) VALUES (?, ?)", schemaName)
	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		myErr := bl.CreateError(bl.ErrAddNoteToCollection, err, "AddNoteToCollection")
		nr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	_, err = tx.ExecContext(ctx, query, noteID, collectionID)
	if err != nil {
		myErr := bl.CreateError(bl.ErrAddNoteToCollection, err, "AddNoteToCollection")
		nr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}

	myOk := bl.CreateError(bl.AllIsOk, nil, "")
	return myOk
}

func (nr *NoteRepository) DeleteNoteFromCollection(collectionID int, noteID int) *bl.MyError {
	nr.MyLogger.WriteLog("DeleteNoteFromCollection is called (Repo)", slog.LevelInfo, nil)

	if collectionID < 0 || noteID < 0 {
		myErr := bl.CreateError(bl.ErrDeleteNoteFromCollection, bl.ErrDeleteNoteFromCollectionError(), "DeleteNoteFromCollection")
		nr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}

	db := nr.DbConfigs.DB
	schemaName := nr.DbConfigs.Name
	query := fmt.Sprintf("ALTER TABLE %s.note_collections DELETE WHERE note_id = ? AND collection_id = ?", schemaName)
	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		myErr := bl.CreateError(bl.ErrDeleteNoteFromCollection, err, "DeleteNoteFromCollection")
		nr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	_, err = tx.ExecContext(ctx, query, noteID, collectionID)
	if err != nil {
		myErr := bl.CreateError(bl.ErrDeleteNoteFromCollection, err, "DeleteNoteFromCollection")
		nr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}

	myOk := bl.CreateError(bl.AllIsOk, nil, "")
	return myOk
}
