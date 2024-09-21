package da

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"notebook_app/internal/business_logic"
	"notebook_app/internal/models"
)

func (sr *SectionRepository) GetSectionByID(id int) (*models.Section, *bl.MyError) {
	sr.MyLogger.WriteLog("GetSectionByID is called (Repo)", slog.LevelInfo, nil)

	if id < 0 {
		myErr := bl.CreateError(bl.ErrGetSectionByID, bl.ErrGetSectionByIDError(), "GetSectionByID")
		sr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		sr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return nil, myErr
	}

	var section models.Section
	db := sr.DbConfigs.DB
	schemaName := sr.DbConfigs.SchemaName
	query := fmt.Sprintf("SELECT * FROM %s.sections WHERE id = $1", schemaName)
	ctx := context.Background()

	err := db.QueryRowContext(ctx, query, id).Scan(&section.Id, &section.CreationDate)

	if err != nil {
		var myErr *bl.MyError

		if errors.Is(err, sql.ErrNoRows) {
			myErr = bl.CreateError(bl.ErrGetSectionByID, bl.ErrGetSectionByIDError(), "GetSectionByID")
		} else {
			myErr = bl.CreateError(bl.ErrGetSectionByID, err, "GetSectionByID")
		}

		sr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return nil, myErr
	}

	myOk := bl.CreateError(bl.AllIsOk, nil, "")
	return &section, myOk
}

func (sr *SectionRepository) GetSectionByTeamName(teamName string) (*models.Section, *bl.MyError) {
	sr.MyLogger.WriteLog("GetSectionByTeamName is called (Repo)", slog.LevelInfo, nil)

	if teamName == "" {
		myErr := bl.CreateError(bl.ErrGetSectionByTeamName, bl.ErrGetSectionByTeamNameError(), "GetSectionByTeamName")
		sr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		sr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return nil, myErr
	}

	var section models.Section
	db := sr.DbConfigs.DB
	schemaName := sr.DbConfigs.SchemaName
	query := fmt.Sprintf("SELECT id, creation_date FROM %s.sections s JOIN %s.teams_sections ts ON s.id = ts.section_id WHERE t.name = $1", schemaName, schemaName)
	ctx := context.Background()

	err := db.QueryRowContext(ctx, query, teamName).Scan(&section.Id, &section.CreationDate)

	if err != nil {
		var myErr *bl.MyError

		if errors.Is(err, sql.ErrNoRows) {
			myErr = bl.CreateError(bl.ErrGetSectionByTeamName, bl.ErrGetSectionByTeamNameError(), "GetSectionByTeamName")
		} else {
			myErr = bl.CreateError(bl.ErrGetSectionByTeamName, err, "GetSectionByTeamName")
		}

		sr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return nil, myErr
	}

	myOk := bl.CreateError(bl.AllIsOk, nil, "")
	return &section, myOk
}

func (sr *SectionRepository) GetAllSections() ([]*models.Section, *bl.MyError) {
	sr.MyLogger.WriteLog("GetAllSections is called (Repo)", slog.LevelInfo, nil)

	db := sr.DbConfigs.DB
	schemaName := sr.DbConfigs.SchemaName
	query := fmt.Sprintf("SELECT * FROM %s.sections", schemaName)
	ctx := context.Background()

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		myErr := bl.CreateError(bl.ErrGetAllSections, err, "GetAllSections")
		sr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return nil, myErr
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
			myErr := bl.CreateError(bl.ErrGetAllSections, err, "GetAllSections")
			sr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
			return nil, myErr
		}

		sections = append(sections, &section)
	}

	if err := rows.Err(); err != nil {
		myErr := bl.CreateError(bl.ErrGetAllSections, err, "GetAllSections")
		sr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return nil, myErr
	}

	myOk := bl.CreateError(bl.AllIsOk, nil, "")
	return sections, myOk
}

func (sr *SectionRepository) AddSection(section *models.Section, team *models.Team) *bl.MyError {
	sr.MyLogger.WriteLog("AddSection is called (Repo)", slog.LevelInfo, nil)

	if section == nil {
		myErr := bl.CreateError(bl.ErrAddSection, bl.ErrAddSectionError(), "AddSection")
		sr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}

	db := sr.DbConfigs.DB
	schemaName := sr.DbConfigs.SchemaName
	query := fmt.Sprintf("INSERT INTO %s.sections (creation_date) VALUES ($1)", schemaName)
	query1 := fmt.Sprintf("INSERT INTO %s.teams_sections (team_id, section_id) VALUES ($2, $3)", schemaName)
	result_query := fmt.Sprintf("%s; %s;", query, query1)
	ctx := context.Background()

	// Add section to sections
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return bl.CreateError(bl.ErrAddSection, err, "AddSection")
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	_, err = tx.ExecContext(ctx, result_query, section.CreationDate, team.Id, section.Id)

	if err != nil {
		myErr := bl.CreateError(bl.ErrAddSection, err, "AddSection")
		sr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}

	// Connect section and team
	tx, err = db.BeginTx(ctx, nil)
	if err != nil {
		return bl.CreateError(bl.ErrAddSection, err, "AddSection")
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	_, err = tx.ExecContext(ctx, query1, team.Id, section.Id)

	if err != nil {
		myErr := bl.CreateError(bl.ErrAddSection, err, "AddSection")
		sr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}

	myOk := bl.CreateError(bl.AllIsOk, nil, "")
	return myOk
}

func (sr *SectionRepository) DeleteSection(id int) *bl.MyError {
	sr.MyLogger.WriteLog("DeleteSection is called (Repo)", slog.LevelInfo, nil)

	if id < 0 {
		myErr := bl.CreateError(bl.ErrDeleteSection, bl.ErrDeleteSectionError(), "DeleteSection")
		sr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}

	db := sr.DbConfigs.DB
	schemaName := sr.DbConfigs.SchemaName
	query := fmt.Sprintf("DELETE FROM %s.sections WHERE id = $1", schemaName)
	query1 := fmt.Sprintf("DELETE FROM %s.teams_sections WHERE section_id = $1", schemaName)
	ctx := context.Background()

	// Delete section from sections
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return bl.CreateError(bl.ErrDeleteSection, err, "DeleteSection")
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	_, err = tx.ExecContext(ctx, query, id)

	if err != nil {
		myErr := bl.CreateError(bl.ErrDeleteSection, err, "DeleteSection")
		sr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}

	// Disconnect section and team
	tx, err = db.BeginTx(ctx, nil)
	if err != nil {
		return bl.CreateError(bl.ErrDeleteSection, err, "DeleteSection")
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	_, err = tx.ExecContext(ctx, query1, id)

	if err != nil {
		myErr := bl.CreateError(bl.ErrDeleteSection, err, "DeleteSection")
		sr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}

	myOk := bl.CreateError(bl.AllIsOk, nil, "")
	return myOk
}

func (sr *SectionRepository) UpdateSection(section *models.Section) *bl.MyError {
	sr.MyLogger.WriteLog("UpdateSection is called (Repo)", slog.LevelInfo, nil)

	if section == nil {
		myErr := bl.CreateError(bl.ErrUpdateSection, bl.ErrUpdateSectionError(), "UpdateSection")
		sr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}

	db := sr.DbConfigs.DB
	schemaName := sr.DbConfigs.SchemaName
	query := fmt.Sprintf("UPDATE %s.sections SET creation_date = $1 WHERE id = $2", schemaName)
	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return bl.CreateError(bl.ErrUpdateSection, err, "UpdateSection")
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	_, err = tx.ExecContext(ctx, query, section.CreationDate, section.Id)

	if err != nil {
		myErr := bl.CreateError(bl.ErrUpdateSection, err, "UpdateSection")
		sr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}

	myOk := bl.CreateError(bl.AllIsOk, nil, "")
	return myOk
}

func (sr *SectionRepository) GetAllNotesInSection(section *models.Section) ([]*models.Note, *bl.MyError) {
	sr.MyLogger.WriteLog("GetAllNotesInSection is called (Repo)", slog.LevelInfo, nil)

	if section == nil {
		myErr := bl.CreateError(bl.ErrGetAllNotesInSection, bl.ErrGetAllNotesInSectionError(), "GetAllNotesInSection")
		sr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return nil, myErr
	}

	db := sr.DbConfigs.DB
	schemaName := sr.DbConfigs.SchemaName
	query := fmt.Sprintf("SELECT * FROM %s.notes WHERE section_id = $1", schemaName)
	ctx := context.Background()

	rows, err := db.QueryContext(ctx, query, section.Id)
	if err != nil {
		myErr := bl.CreateError(bl.ErrGetAllNotesInSection, err, "GetAllNotesInSection")
		sr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
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
			myErr := bl.CreateError(bl.ErrGetAllNotesInSection, err, "GetAllNotesInSection")
			sr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
			return nil, myErr
		}

		notes = append(notes, &note)
	}

	if err := rows.Err(); err != nil {
		myErr := bl.CreateError(bl.ErrGetAllNotesInSection, err, "GetAllNotesInSection")
		sr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return nil, myErr
	}

	myOk := bl.CreateError(bl.AllIsOk, nil, "")
	return notes, myOk
}

func (sr *SectionRepository) AddNoteToSection(note *models.Note, section *models.Section) *bl.MyError {
	sr.MyLogger.WriteLog("AddNoteToSection is called (Repo)", slog.LevelInfo, nil)

	if note == nil || section == nil {
		myErr := bl.CreateError(bl.ErrAddNoteToSection, bl.ErrAddNoteToSectionError(), "AddNoteToSection")
		sr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}

	db := sr.DbConfigs.DB
	schemaName := sr.DbConfigs.SchemaName
	query := fmt.Sprintf("UPDATE %s.notes SET section_id = $1 WHERE id = $2", schemaName)
	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return bl.CreateError(bl.ErrAddNoteToSection, err, "AddNoteToSection")
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	_, err = tx.ExecContext(ctx, query, section.Id, note.Id)
	if err != nil {
		myErr := bl.CreateError(bl.ErrAddNoteToSection, err, "AddNoteToSection")
		sr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}

	myOk := bl.CreateError(bl.AllIsOk, nil, "")
	return myOk
}

func (sr *SectionRepository) DeleteNoteFromSection(note *models.Note, section *models.Section) *bl.MyError {
	sr.MyLogger.WriteLog("AddNoteToSection is called (Repo)", slog.LevelInfo, nil)

	if note == nil || section == nil {
		myErr := bl.CreateError(bl.ErrAddNoteToSection, bl.ErrAddNoteToSectionError(), "AddNoteToSection")
		sr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}

	db := sr.DbConfigs.DB
	schemaName := sr.DbConfigs.SchemaName
	query := fmt.Sprintf("UPDATE %s.notes SET section_id = -1 WHERE id = $1 AND section_id = $2", schemaName)
	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return bl.CreateError(bl.ErrAddNoteToSection, err, "AddNoteToSection")
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	_, err = tx.ExecContext(ctx, query, note.Id, section.Id)
	if err != nil {
		myErr := bl.CreateError(bl.ErrAddNoteToSection, err, "AddNoteToSection")
		sr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}

	myOk := bl.CreateError(bl.AllIsOk, nil, "")
	return myOk
}
