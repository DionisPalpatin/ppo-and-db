package data_access

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	mylogger "github.com/DionisPalpatin/ppo-and-db/application/internal/logger"
	"github.com/DionisPalpatin/ppo-and-db/application/internal/models"
)

func (sr *SectionRepository) GetSectionByID(id int) (*models.Section, *bl.MyError) {
	sr.MyLogger.WriteLog("GetSectionByID is called (Repo)", slog.LevelInfo, nil)

	if id < 0 {
		resState := bl.CreateError(bl.ErrInParameter, "GetSectionByID", "data_access")
		sr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		return nil, resState
	}

	var section models.Section
	db := sr.DbConfigs.DB
	schemaName := sr.DbConfigs.SchemaName
	query := fmt.Sprintf(getSectionByIDQuery, schemaName)
	ctx := context.Background()

	err := db.QueryRowContext(ctx, query, id).Scan(&section.Id, &section.CreationDate)
	if err != nil {
		var resState *bl.MyError

		if errors.Is(err, sql.ErrNoRows) {
			resState = bl.CreateError(bl.NoSuchSect, "GetSectionByID", "data_access")
			sr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		} else {
			resState := bl.CreateError(bl.DatabaseError, "GetSectionByID", "data_access")
			sr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		}

		return nil, resState
	}

	resState := bl.CreateError(bl.Ok, "GetSectionByID", "data_access")
	sr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelInfo, nil)
	return &section, resState
}

func (sr *SectionRepository) GetSectionByTeamName(teamName string) (*models.Section, *bl.MyError) {
	sr.MyLogger.WriteLog("GetSectionByTeamName is called (Repo)", slog.LevelInfo, nil)

	if teamName == "" {
		resState := bl.CreateError(bl.ErrInParameter, "GetSectionByTeamName", "data_access")
		sr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		return nil, resState
	}

	var section models.Section
	db := sr.DbConfigs.DB
	schemaName := sr.DbConfigs.SchemaName
	query := fmt.Sprintf(getSectionByTeamNameQuery, schemaName, schemaName)
	ctx := context.Background()

	err := db.QueryRowContext(ctx, query, teamName).Scan(&section.Id, &section.CreationDate)

	if err != nil {
		var resState *bl.MyError

		if errors.Is(err, sql.ErrNoRows) {
			resState = bl.CreateError(bl.NoSuchSect, "GetSectionByTeamName", "data_access")
			sr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		} else {
			resState := bl.CreateError(bl.DatabaseError, "GetSectionByTeamName", "data_access")
			sr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		}

		return nil, resState
	}

	resState := bl.CreateError(bl.Ok, "GetSectionByTeamName", "data_access")
	sr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelInfo, nil)
	return &section, resState
}

func (sr *SectionRepository) GetAllSections() ([]*models.Section, *bl.MyError) {
	sr.MyLogger.WriteLog("GetAllSections is called (Repo)", slog.LevelInfo, nil)

	db := sr.DbConfigs.DB
	schemaName := sr.DbConfigs.SchemaName
	query := fmt.Sprintf(getAllSectionsQuery, schemaName)
	ctx := context.Background()

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		var resState *bl.MyError

		if errors.Is(err, sql.ErrNoRows) {
			resState = bl.CreateError(bl.EmptyResult, "GetAllSections", "data_access")
			sr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		} else {
			resState := bl.CreateError(bl.DatabaseError, "GetAllSections", "data_access")
			sr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		}

		return nil, resState
	}
	defer rows.Close()

	var sections []*models.Section
	for rows.Next() {
		var section models.Section
		err := rows.Scan(
			&section.Id,
			&section.CreationDate,
		)

		if err != nil {
			resState := bl.CreateError(bl.DatabaseError, "GetAllSections", "data_access")
			sr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
			return nil, resState
		}

		sections = append(sections, &section)
	}

	if err := rows.Err(); err != nil {
		resState := bl.CreateError(bl.DatabaseError, "GetAllSections", "data_access")
		sr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		return nil, resState
	}

	resState := bl.CreateError(bl.Ok, "GetAllSections", "data_access")
	sr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelInfo, nil)
	return sections, resState
}

func (sr *SectionRepository) AddSection(section *models.Section, team *models.Team) (int, *bl.MyError) {
	sr.MyLogger.WriteLog("AddSection is called (Repo)", slog.LevelInfo, nil)

	if section == nil {
		resState := bl.CreateError(bl.ErrInParameter, "AddSection", "data_access")
		sr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		return 0, resState
	}

	db := sr.DbConfigs.DB
	schemaName := sr.DbConfigs.SchemaName
	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		resState := bl.CreateError(bl.DatabaseError, "AddSection", "data_access")
		sr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		return 0, resState
	}
	defer tx.Rollback()

	query := fmt.Sprintf(addSectionQuery, schemaName)
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		resState := bl.CreateError(bl.DatabaseError, "AddSection", "data_access")
		sr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		return 0, resState
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, section.CreationDate).Scan(section.Id)

	if err != nil {
		var resState *bl.MyError

		if errors.Is(err, sql.ErrNoRows) {
			resState = bl.CreateError(bl.NoSuchSect, "AddSection", "data_access")
			sr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		} else {
			resState := bl.CreateError(bl.DatabaseError, "AddSection", "data_access")
			sr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		}

		return 0, resState
	}

	query = fmt.Sprintf(addSectionToTeamQuery, schemaName)
	_, err = tx.ExecContext(ctx, query, team.Id, section.Id)

	if err != nil {
		resState := bl.CreateError(bl.DatabaseError, "AddSection", "data_access")
		sr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		return 0, resState
	}

	err = tx.Commit()
	if err != nil {
		resState := bl.CreateError(bl.DatabaseError, "AddSection", "data_access")
		sr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		return 0, resState
	}

	resState := bl.CreateError(bl.Ok, "AddSection", "data_access")
	sr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelInfo, nil)
	return section.Id, resState
}

func (sr *SectionRepository) DeleteSection(id int) *bl.MyError {
	sr.MyLogger.WriteLog("DeleteSection is called (Repo)", slog.LevelInfo, nil)

	if id <= 0 {
		resState := bl.CreateError(bl.ErrInParameter, "DeleteSection", "data_access")
		sr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		return resState
	}

	db := sr.DbConfigs.DB
	schemaName := sr.DbConfigs.SchemaName
	query := fmt.Sprintf(deleteSectionQuery, schemaName, schemaName)
	ctx := context.Background()

	// Delete section from sections
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		resState := bl.CreateError(bl.DatabaseError, "DeleteSection", "data_access")
		sr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		return resState
	}
	defer deferTransaction(err, tx)

	_, err = tx.ExecContext(ctx, query, id)

	if err != nil {
		var resState *bl.MyError

		if errors.Is(err, sql.ErrNoRows) {
			resState = bl.CreateError(bl.NoSuchSect, "DeleteSection", "data_access")
			sr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		} else {
			resState := bl.CreateError(bl.DatabaseError, "DeleteSection", "data_access")
			sr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		}

		return resState
	}

	resState := bl.CreateError(bl.Ok, "DeleteSection", "data_access")
	sr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelInfo, nil)
	return resState
}

func (sr *SectionRepository) UpdateSection(section *models.Section) *bl.MyError {
	sr.MyLogger.WriteLog("UpdateSection is called (Repo)", slog.LevelInfo, nil)

	if section == nil {
		resState := bl.CreateError(bl.ErrInParameter, "UpdateSection", "data_access")
		sr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		return resState
	}

	db := sr.DbConfigs.DB
	schemaName := sr.DbConfigs.SchemaName
	query := fmt.Sprintf(updateSectionQuery, schemaName)
	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		resState := bl.CreateError(bl.DatabaseError, "UpdateSection", "data_access")
		sr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		return resState
	}
	defer deferTransaction(err, tx)

	_, err = tx.ExecContext(ctx, query, section.CreationDate, section.Id)

	if err != nil {
		var resState *bl.MyError

		if errors.Is(err, sql.ErrNoRows) {
			resState = bl.CreateError(bl.NoSuchSect, "UpdateSection", "data_access")
			sr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		} else {
			resState := bl.CreateError(bl.DatabaseError, "UpdateSection", "data_access")
			sr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		}

		return resState
	}

	resState := bl.CreateError(bl.Ok, "UpdateSection", "data_access")
	sr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelInfo, nil)
	return resState
}

func (sr *SectionRepository) GetAllNotesInSection(section *models.Section) ([]*models.Note, *bl.MyError) {
	sr.MyLogger.WriteLog("GetAllNotesInSection is called (Repo)", slog.LevelInfo, nil)

	if section == nil {
		resState := bl.CreateError(bl.ErrInParameter, "GetAllNotesInSection", "data_access")
		sr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		return nil, resState
	}

	db := sr.DbConfigs.DB
	schemaName := sr.DbConfigs.SchemaName
	query := fmt.Sprintf(getAllNotesInSectionQuery, schemaName)
	ctx := context.Background()

	rows, err := db.QueryContext(ctx, query, section.Id)
	if err != nil {
		var resState *bl.MyError

		if errors.Is(err, sql.ErrNoRows) {
			resState = bl.CreateError(bl.NoSuchSect, "GetAllNotesInSection", "data_access")
			sr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		} else {
			resState := bl.CreateError(bl.DatabaseError, "GetAllNotesInSection", "data_access")
			sr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
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
			resState := bl.CreateError(bl.DatabaseError, "GetAllNotesInSection", "data_access")
			sr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
			return nil, resState
		}

		notes = append(notes, &note)
	}

	if err := rows.Err(); err != nil {
		resState := bl.CreateError(bl.DatabaseError, "GetAllNotesInSection", "data_access")
		sr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		return nil, resState
	}

	resState := bl.CreateError(bl.Ok, "GetAllNotesInSection", "data_access")
	sr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelInfo, nil)
	return notes, resState
}

func (sr *SectionRepository) AddNoteToSection(note *models.Note, section *models.Section) *bl.MyError {
	sr.MyLogger.WriteLog("AddNoteToSection is called (Repo)", slog.LevelInfo, nil)

	if note == nil || section == nil {
		resState := bl.CreateError(bl.ErrInParameter, "AddNoteToSection", "data_access")
		sr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		return resState
	}

	db := sr.DbConfigs.DB
	schemaName := sr.DbConfigs.SchemaName
	query := fmt.Sprintf(addNoteToSectionQuery, schemaName)
	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		resState := bl.CreateError(bl.DatabaseError, "AddNoteToSection", "data_access")
		sr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		return resState
	}
	defer deferTransaction(err, tx)

	_, err = tx.ExecContext(ctx, query, section.Id, note.Id)
	if err != nil {
		var resState *bl.MyError

		if errors.Is(err, sql.ErrNoRows) {
			resState = bl.CreateError(bl.NoSuchSect, "AddNoteToSection", "data_access")
			sr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		} else {
			resState := bl.CreateError(bl.DatabaseError, "AddNoteToSection", "data_access")
			sr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		}

		return resState
	}

	resState := bl.CreateError(bl.Ok, "AddNoteToSection", "data_access")
	sr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelInfo, nil)
	return resState
}

func (sr *SectionRepository) DeleteNoteFromSection(note *models.Note, section *models.Section) *bl.MyError {
	sr.MyLogger.WriteLog("AddNoteToSection is called (Repo)", slog.LevelInfo, nil)

	if note == nil || section == nil {
		resState := bl.CreateError(bl.ErrInParameter, "DeleteNoteFromSection", "data_access")
		sr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		return resState
	}

	db := sr.DbConfigs.DB
	schemaName := sr.DbConfigs.SchemaName
	query := fmt.Sprintf(deleteNoteFromSectionQuery, schemaName)
	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		resState := bl.CreateError(bl.DatabaseError, "DeleteNoteFromSection", "data_access")
		sr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		return resState
	}
	defer deferTransaction(err, tx)

	_, err = tx.ExecContext(ctx, query, note.Id, section.Id)
	if err != nil {
		var resState *bl.MyError

		if errors.Is(err, sql.ErrNoRows) {
			resState = bl.CreateError(bl.NoSuchSect, "DeleteNoteFromSection", "data_access")
			sr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		} else {
			resState := bl.CreateError(bl.DatabaseError, "DeleteNoteFromSection", "data_access")
			sr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		}

		return resState
	}

	resState := bl.CreateError(bl.Ok, "DeleteNoteFromSection", "data_access")
	sr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelInfo, nil)
	return resState
}
