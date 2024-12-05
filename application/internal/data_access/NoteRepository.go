package data_access

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/business_logic"
	mylogger "github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/logger"
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/models"
)

// ----------------------------------------------------------------------------
// Основные методы
// ----------------------------------------------------------------------------

func (nr *NoteRepository) GetNoteByID(id int) (*models.Note, *bl.MyError) {
	nr.MyLogger.WriteLog("GetNoteByID is called (Repo)", slog.LevelInfo, nil)

	if id < 0 {
		resState := bl.CreateError(bl.ErrInParameter, "GetNoteByID", "data_access")
		nr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		return nil, resState
	}

	var note models.Note
	db := nr.DbConfigs.DB
	schemaName := nr.DbConfigs.SchemaName
	ctx := context.Background()

	query := fmt.Sprintf(getNoteByIDQuery, schemaName)
	err := db.QueryRowContext(ctx, query, id).Scan(
		&note.Id,
		&note.Access,
		&note.Name,
		&note.Likes,
		&note.Dislikes,
		&note.RegistrationDate,
		&note.OwnerID,
		&note.SectionID,
	)

	if err != nil {
		var resState *bl.MyError

		if errors.Is(err, sql.ErrNoRows) {
			resState = bl.CreateError(bl.NoSuchNote, "GetNoteByID", "data_access")
			nr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		} else {
			resState = bl.CreateError(bl.DatabaseError, "GetNoteByID", "data_access")
			nr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		}

		return nil, resState
	}

	query = fmt.Sprintf(getNoteTextContent, schemaName)
	err = db.QueryRowContext(ctx, query, note.Id).Scan(&note.Content.Text, &note.Content.TextExt)

	if err == nil {
		query = fmt.Sprintf(getNoteImageContent, schemaName)
		err = db.QueryRowContext(ctx, query, note.Id).Scan(&note.Content.Img, &note.Content.ImgExt)
	}

	if err == nil {
		query = fmt.Sprintf(getNoteRawDataContent, schemaName)
		err = db.QueryRowContext(ctx, query, note.Id).Scan(&note.Content.Raw, &note.Content.RawExt)
	}

	if err != nil {
		var resState *bl.MyError

		if errors.Is(err, sql.ErrNoRows) {
			resState = bl.CreateError(bl.NoSuchNote, "GetNoteByID", "data_access")
			nr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		} else {
			resState = bl.CreateError(bl.DatabaseError, "GetNoteByID", "data_access")
			nr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		}

		return nil, resState
	}

	resState := bl.CreateError(bl.Ok, "GetNoteByID", "data_access")
	nr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelInfo, nil)
	return &note, resState
}

func (nr *NoteRepository) GetNoteByName(name string) (*models.Note, *bl.MyError) {
	nr.MyLogger.WriteLog("GetNoteByName is called (Repo)", slog.LevelInfo, nil)

	if name == "" {
		resState := bl.CreateError(bl.ErrInParameter, "GetNoteByName", "data_access")
		nr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		return nil, resState
	}

	var note models.Note
	db := nr.DbConfigs.DB
	schemaName := nr.DbConfigs.SchemaName
	query := fmt.Sprintf(getNoteByNameQuery, schemaName)
	ctx := context.Background()

	err := db.QueryRowContext(ctx, query, name).Scan(
		&note.Id,
		&note.Access,
		&note.Name,
		&note.Likes,
		&note.Dislikes,
		&note.RegistrationDate,
		&note.OwnerID,
		&note.SectionID,
	)

	if err != nil {
		var resState *bl.MyError

		if errors.Is(err, sql.ErrNoRows) {
			resState = bl.CreateError(bl.NoSuchNote, "GetNoteByName", "data_access")
			nr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		} else {
			resState = bl.CreateError(bl.DatabaseError, "GetNoteByName", "data_access")
			nr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		}

		return nil, resState
	}

	query = fmt.Sprintf(getNoteTextContent, schemaName)
	err = db.QueryRowContext(ctx, query, note.Id).Scan(&note.Content.Text, &note.Content.TextExt)

	if err == nil {
		query = fmt.Sprintf(getNoteImageContent, schemaName)
		err = db.QueryRowContext(ctx, query, note.Id).Scan(&note.Content.Img, &note.Content.ImgExt)
	}

	if err == nil {
		query = fmt.Sprintf(getNoteRawDataContent, schemaName)
		err = db.QueryRowContext(ctx, query, note.Id).Scan(&note.Content.Raw, &note.Content.RawExt)
	}

	if err != nil {
		var resState *bl.MyError

		if errors.Is(err, sql.ErrNoRows) {
			resState = bl.CreateError(bl.NoSuchNote, "GetNoteByName", "data_access")
			nr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		} else {
			resState = bl.CreateError(bl.DatabaseError, "GetNoteByName", "data_access")
			nr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		}

		return nil, resState
	}

	resState := bl.CreateError(bl.Ok, "GetNoteByName", "data_access")
	nr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelInfo, nil)
	return &note, resState
}

func (nr *NoteRepository) GetAllNotes() ([]*models.Note, *bl.MyError) {
	nr.MyLogger.WriteLog("GetAllNotes is called (Repo)", slog.LevelInfo, nil)

	db := nr.DbConfigs.DB
	schemaName := nr.DbConfigs.SchemaName
	query := fmt.Sprintf(getAllNotesQuery, schemaName)
	ctx := context.Background()

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		var resState *bl.MyError

		if errors.Is(err, sql.ErrNoRows) {
			resState = bl.CreateError(bl.EmptyResult, "GetAllNotes", "data_access")
			nr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		} else {
			resState = bl.CreateError(bl.DatabaseError, "GetAllNotes", "data_access")
			nr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		}

		return nil, resState
	}
	defer rows.Close()

	var notes []*models.Note
	for rows.Next() {
		var note models.Note
		err := rows.Scan(
			&note.Id,
			&note.Access,
			&note.Name,
			&note.Likes,
			&note.Dislikes,
			&note.RegistrationDate,
			&note.OwnerID,
			&note.SectionID,
		)

		if err != nil {
			resState := bl.CreateError(bl.DatabaseError, "GetAllNotes", "data_access")
			nr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
			return nil, resState
		}

		notes = append(notes, &note)
	}

	if err := rows.Err(); err != nil {
		resState := bl.CreateError(bl.DatabaseError, "GetAllNotes", "data_access")
		nr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		return nil, resState
	}

	resState := bl.CreateError(bl.Ok, "GetAllNotes", "data_access")
	nr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelInfo, nil)
	return notes, resState
}

func (nr *NoteRepository) GetAllPublicNotes() ([]*models.Note, *bl.MyError) {
	nr.MyLogger.WriteLog("GetAllNotes is called (Repo)", slog.LevelInfo, nil)

	db := nr.DbConfigs.DB
	schemaName := nr.DbConfigs.SchemaName
	query := fmt.Sprintf(getAllPublicNotesQuery, schemaName)
	ctx := context.Background()

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		var resState *bl.MyError

		if errors.Is(err, sql.ErrNoRows) {
			resState = bl.CreateError(bl.EmptyResult, "GetAllPublicNotes", "data_access")
			nr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		} else {
			resState = bl.CreateError(bl.DatabaseError, "GetAllPublicNotes", "data_access")
			nr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		}

		return nil, resState
	}
	defer rows.Close()

	var notes []*models.Note
	for rows.Next() {
		var note models.Note
		err := rows.Scan(
			&note.Id,
			&note.Access,
			&note.Name,
			&note.Likes,
			&note.Dislikes,
			&note.RegistrationDate,
			&note.OwnerID,
			&note.SectionID,
		)

		if err != nil {
			resState := bl.CreateError(bl.DatabaseError, "GetAllPublicNotes", "data_access")
			nr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
			return nil, resState
		}

		notes = append(notes, &note)
	}

	if err := rows.Err(); err != nil {
		resState := bl.CreateError(bl.DatabaseError, "GetAllPublicNotes", "data_access")
		nr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		return nil, resState
	}

	resState := bl.CreateError(bl.Ok, "GetAllPublicNotes", "data_access")
	nr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelInfo, nil)
	return notes, resState
}

func (nr *NoteRepository) AddNote(note *models.Note) *bl.MyError {
	nr.MyLogger.WriteLog("AddNote is called (Repo)", slog.LevelInfo, nil)

	if note == nil {
		resState := bl.CreateError(bl.ErrInParameter, "AddNote", "data_access")
		nr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		return resState
	}

	db := nr.DbConfigs.DB
	schemaName := nr.DbConfigs.SchemaName
	ctx := context.Background()

	query := fmt.Sprintf(addNoteInfoQuery, schemaName)
	err := addOrUpdateNoteInfo(note, db, ctx, query)
	if err != nil {
		var resState *bl.MyError

		if errors.Is(err, sql.ErrNoRows) {
			resState = bl.CreateError(bl.NoSuchNote, "AddNote", "data_access")
			nr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		} else {
			resState = bl.CreateError(bl.DatabaseError, "AddNote", "data_access")
			nr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		}

		return resState
	}

	if len(note.Content.Text) != 0 {
		query = fmt.Sprintf(addNoteTextQuery, schemaName)
		err = addOrUpdateNoteContent(note.Id, note.Content.Text, note.Content.TextExt, db, ctx, query)
		if err != nil {
			var resState *bl.MyError

			if errors.Is(err, sql.ErrNoRows) {
				resState = bl.CreateError(bl.NoSuchNote, "AddNote", "data_access")
				nr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
			} else {
				resState = bl.CreateError(bl.DatabaseError, "AddNote", "data_access")
				nr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
			}

			return resState
		}
	}

	if len(note.Content.Img) != 0 {
		query += fmt.Sprintf(addNoteImageQuery, schemaName)
		err = addOrUpdateNoteContent(note.Id, note.Content.Img, note.Content.ImgExt, db, ctx, query)
		if err != nil {
			var resState *bl.MyError

			if errors.Is(err, sql.ErrNoRows) {
				resState = bl.CreateError(bl.NoSuchNote, "AddNote", "data_access")
				nr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
			} else {
				resState = bl.CreateError(bl.DatabaseError, "AddNote", "data_access")
				nr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
			}

			return resState
		}
	}

	if len(note.Content.Raw) != 0 {
		query += fmt.Sprintf(addNoteRawDataQuery, schemaName)
		err = addOrUpdateNoteContent(note.Id, note.Content.Raw, note.Content.RawExt, db, ctx, query)
		if err != nil {
			var resState *bl.MyError

			if errors.Is(err, sql.ErrNoRows) {
				resState = bl.CreateError(bl.NoSuchNote, "AddNote", "data_access")
				nr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
			} else {
				resState = bl.CreateError(bl.DatabaseError, "AddNote", "data_access")
				nr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
			}

			return resState
		}
	}

	resState := bl.CreateError(bl.Ok, "AddNote", "data_access")
	nr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelInfo, nil)
	return resState
}

func (nr *NoteRepository) DeleteNote(id int) *bl.MyError {
	nr.MyLogger.WriteLog("DeleteNote is called (Repo)", slog.LevelInfo, nil)

	if id < 0 {
		resState := bl.CreateError(bl.ErrInParameter, "DeleteNote", "data_access")
		nr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		return resState
	}

	db := nr.DbConfigs.DB
	schemaName := nr.DbConfigs.SchemaName
	ctx := context.Background()

	query := fmt.Sprintf(deleteNoteQuery, schemaName)

	tx, err := db.BeginTx(ctx, nil)
	defer deferTransaction(err, tx)

	if err != nil {
		resState := bl.CreateError(bl.DatabaseError, "DeleteNote", "data_access")
		nr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		return resState
	}

	_, err = tx.ExecContext(ctx, query, id)

	if err != nil {
		var resState *bl.MyError

		if errors.Is(err, sql.ErrNoRows) {
			resState = bl.CreateError(bl.NoSuchNote, "DeleteNote", "data_access")
			nr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		} else {
			resState = bl.CreateError(bl.DatabaseError, "DeleteNote", "data_access")
			nr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		}

		return resState
	}

	resState := bl.CreateError(bl.Ok, "DeleteNote", "data_access")
	nr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelInfo, nil)
	return resState
}

func (nr *NoteRepository) UpdateNoteContent(note *models.Note) *bl.MyError {
	nr.MyLogger.WriteLog("UpdateNoteContent is called (Repo)", slog.LevelInfo, nil)

	if note == nil {
		resState := bl.CreateError(bl.ErrInParameter, "UpdateNoteContent", "data_access")
		nr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		return resState
	}

	db := nr.DbConfigs.DB
	schemaName := nr.DbConfigs.SchemaName
	ctx := context.Background()

	query := fmt.Sprintf(updateNoteTextContentQuery, schemaName)
	err := addOrUpdateNoteContent(note.Id, note.Content.Text, note.Content.TextExt, db, ctx, query)
	if err != nil {
		var resState *bl.MyError

		if errors.Is(err, sql.ErrNoRows) {
			resState = bl.CreateError(bl.NoSuchNote, "UpdateNoteContent", "data_access")
			nr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		} else {
			resState = bl.CreateError(bl.DatabaseError, "UpdateNoteContent", "data_access")
			nr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		}

		return resState
	}

	query = fmt.Sprintf(updateNoteImageContentQuery, schemaName)
	err = addOrUpdateNoteContent(note.Id, note.Content.Img, note.Content.ImgExt, db, ctx, query)
	if err != nil {
		var resState *bl.MyError

		if errors.Is(err, sql.ErrNoRows) {
			resState = bl.CreateError(bl.NoSuchNote, "UpdateNoteContent", "data_access")
			nr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		} else {
			resState = bl.CreateError(bl.DatabaseError, "UpdateNoteContent", "data_access")
			nr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		}

		return resState
	}

	query = fmt.Sprintf(updateNoteRawContentQuery, schemaName)
	err = addOrUpdateNoteContent(note.Id, note.Content.Raw, note.Content.RawExt, db, ctx, query)
	if err != nil {
		var resState *bl.MyError

		if errors.Is(err, sql.ErrNoRows) {
			resState = bl.CreateError(bl.NoSuchNote, "UpdateNoteContent", "data_access")
			nr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		} else {
			resState = bl.CreateError(bl.DatabaseError, "UpdateNoteContent", "data_access")
			nr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		}

		return resState
	}

	resState := bl.CreateError(bl.Ok, "UpdateNoteContent", "data_access")
	nr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelInfo, nil)
	return resState
}

func (nr *NoteRepository) UpdateNoteInfo(note *models.Note) *bl.MyError {
	nr.MyLogger.WriteLog("UpdateNoteInfo is called (Repo)", slog.LevelInfo, nil)

	if note == nil {
		resState := bl.CreateError(bl.ErrInParameter, "UpdateNoteInfo", "data_access")
		nr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		return resState
	}

	db := nr.DbConfigs.DB
	schemaName := nr.DbConfigs.SchemaName
	ctx := context.Background()

	query := fmt.Sprintf(updateNoteInfoQuery, schemaName)

	err := addOrUpdateNoteInfo(note, db, ctx, query)
	if err != nil {
		var resState *bl.MyError

		if errors.Is(err, sql.ErrNoRows) {
			resState = bl.CreateError(bl.NoSuchNote, "UpdateNoteInfo", "data_access")
			nr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		} else {
			resState = bl.CreateError(bl.DatabaseError, "UpdateNoteInfo", "data_access")
			nr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		}

		return resState
	}

	resState := bl.CreateError(bl.Ok, "UpdateNoteInfo", "data_access")
	nr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelInfo, nil)
	return resState
}

func (nr *NoteRepository) AddNoteToCollection(collectionID int, noteID int) *bl.MyError {
	nr.MyLogger.WriteLog("AddNoteToCollection is called (Repo)", slog.LevelInfo, nil)

	if collectionID < 0 || noteID < 0 {
		resState := bl.CreateError(bl.ErrInParameter, "AddNoteToCollection", "data_access")
		nr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		return resState
	}

	db := nr.DbConfigs.DB
	schemaName := nr.DbConfigs.SchemaName
	query := fmt.Sprintf(addNoteToCollectionQuery, schemaName)
	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		resState := bl.CreateError(bl.DatabaseError, "AddNoteToCollection", "data_access")
		nr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		return resState
	}
	defer deferTransaction(err, tx)

	_, err = tx.ExecContext(ctx, query, noteID, collectionID)
	if err != nil {
		var resState *bl.MyError

		if errors.Is(err, sql.ErrNoRows) {
			resState = bl.CreateError(bl.NoSuchNote, "AddNoteToCollection", "data_access")
			nr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		} else {
			resState = bl.CreateError(bl.DatabaseError, "AddNoteToCollection", "data_access")
			nr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		}

		return resState
	}

	resState := bl.CreateError(bl.Ok, "AddNoteToCollection", "data_access")
	nr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelInfo, nil)
	return resState
}

func (nr *NoteRepository) DeleteNoteFromCollection(collectionID int, noteID int) *bl.MyError {
	nr.MyLogger.WriteLog("DeleteNoteFromCollection is called (Repo)", slog.LevelInfo, nil)

	if collectionID < 0 || noteID < 0 {
		resState := bl.CreateError(bl.ErrInParameter, "DeleteNoteFromCollection", "data_access")
		nr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		return resState
	}

	db := nr.DbConfigs.DB
	schemaName := nr.DbConfigs.SchemaName
	query := fmt.Sprintf(deleteNoteFromCollectionQuery, schemaName)
	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)
	defer deferTransaction(err, tx)
	if err != nil {
		resState := bl.CreateError(bl.DatabaseError, "DeleteNoteFromCollection", "data_access")
		nr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		return resState
	}

	_, err = tx.ExecContext(ctx, query, noteID, collectionID)
	if err != nil {
		var resState *bl.MyError

		if errors.Is(err, sql.ErrNoRows) {
			resState = bl.CreateError(bl.NoSuchNote, "DeleteNoteFromCollection", "data_access")
			nr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		} else {
			resState = bl.CreateError(bl.DatabaseError, "DeleteNoteFromCollection", "data_access")
			nr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		}

		return resState
	}

	resState := bl.CreateError(bl.Ok, "DeleteNoteFromCollection", "data_access")
	nr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelInfo, nil)
	return resState
}
