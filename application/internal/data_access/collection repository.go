package data_access

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	bl "github.com/DionisPalpatin/ppo-and-db/application/internal/business_logic"
	mylogger "github.com/DionisPalpatin/ppo-and-db/application/internal/logger"
	"github.com/DionisPalpatin/ppo-and-db/application/internal/models"
)

func (cr *CollectionRepository) GetCollectionByID(id int) (*models.Collection, *bl.MyError) {
	cr.MyLogger.WriteLog("GetCollectionByID is called (Repo)", slog.LevelInfo, nil)

	if id < 0 {
		resState := bl.CreateError(bl.ErrInParameter, "GetCollectionByID", "data_access")
		cr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		return nil, resState
	}

	var collection models.Collection
	db := cr.DbConfigs.DB
	schemaName := cr.DbConfigs.SchemaName
	query := fmt.Sprintf(getCollectionByIDQuery, schemaName)
	ctx := context.Background()

	err := db.QueryRowContext(ctx, query, id).Scan(
		&collection.Id,
		&collection.Name,
		&collection.CreationDate,
	)

	if err != nil {
		var resState *bl.MyError

		if errors.Is(err, sql.ErrNoRows) {
			resState = bl.CreateError(bl.NoSuchColl, "GetCollectionByID", "data_access")
			cr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		} else {
			resState := bl.CreateError(bl.DatabaseError, "GetCollectionByID", "data_access")
			cr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		}

		return nil, resState
	}

	resState := bl.CreateError(bl.Ok, "GetCollectionByID", "data_access")
	cr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelInfo, nil)
	return &collection, resState
}

func (cr *CollectionRepository) GetCollectionByName(name string) (*models.Collection, *bl.MyError) {
	cr.MyLogger.WriteLog("GetCollectionByName is called (Repo)", slog.LevelInfo, nil)

	if name == "" {
		resState := bl.CreateError(bl.ErrInParameter, "GetCollectionByName", "data_access")
		cr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		return nil, resState
	}

	var collection models.Collection
	db := cr.DbConfigs.DB
	schemaName := cr.DbConfigs.SchemaName
	query := fmt.Sprintf(getCollectionByNameQuery, schemaName)
	ctx := context.Background()

	err := db.QueryRowContext(ctx, query, name).Scan(
		&collection.Id,
		&collection.Name,
		&collection.CreationDate,
	)

	if err != nil {
		var resState *bl.MyError

		if errors.Is(err, sql.ErrNoRows) {
			resState = bl.CreateError(bl.NoSuchColl, "GetCollectionByName", "data_access")
			cr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		} else {
			resState := bl.CreateError(bl.DatabaseError, "GetCollectionByName", "data_access")
			cr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		}

		return nil, resState
	}

	resState := bl.CreateError(bl.Ok, "GetCollectionByName", "data_access")
	cr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelInfo, nil)
	return &collection, resState
}

func (cr *CollectionRepository) GetAllCollections() ([]*models.Collection, *bl.MyError) {
	cr.MyLogger.WriteLog("GetAllCollections is called (Repo)", slog.LevelInfo, nil)

	db := cr.DbConfigs.DB
	schemaName := cr.DbConfigs.SchemaName
	query := fmt.Sprintf(getAllCollectionsQuery, schemaName)
	ctx := context.Background()

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		var resState *bl.MyError

		if errors.Is(err, sql.ErrNoRows) {
			resState = bl.CreateError(bl.NoSuchColl, "GetAllCollections", "data_access")
			cr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		} else {
			resState := bl.CreateError(bl.DatabaseError, "GetAllCollections", "data_access")
			cr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		}

		return nil, resState
	}
	defer rows.Close()

	var collections []*models.Collection

	for rows.Next() {
		var collection models.Collection
		err = rows.Scan(
			&collection.Id,
			&collection.Name,
			&collection.CreationDate,
		)

		if err != nil {
			resState := bl.CreateError(bl.DatabaseError, "GetAllCollections", "data_access")
			cr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
			return nil, resState
		}

		collections = append(collections, &collection)
	}

	if err = rows.Err(); err != nil {
		resState := bl.CreateError(bl.DatabaseError, "GetAllCollections", "data_access")
		cr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		return nil, resState
	}

	resState := bl.CreateError(bl.Ok, "GetAllCollections", "data_access")
	cr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelInfo, nil)
	return collections, resState
}

func (cr *CollectionRepository) GetAllUserCollections(user *models.User) ([]*models.Collection, *bl.MyError) {
	cr.MyLogger.WriteLog("GetAllUserCollections is called (Repo)", slog.LevelInfo, nil)

	if user == nil {
		resState := bl.CreateError(bl.ErrInParameter, "GetAllUserCollections", "data_access")
		cr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		return nil, resState
	}

	db := cr.DbConfigs.DB
	schemaName := cr.DbConfigs.SchemaName
	query := fmt.Sprintf(getAllUserCollectionsQuery, schemaName, schemaName)
	ctx := context.Background()

	rows, err := db.QueryContext(ctx, query, user.Id)
	if err != nil {
		var resState *bl.MyError

		if errors.Is(err, sql.ErrNoRows) {
			resState = bl.CreateError(bl.NoSuchColl, "GetAllUserCollections", "data_access")
			cr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		} else {
			resState := bl.CreateError(bl.DatabaseError, "GetAllUserCollections", "data_access")
			cr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		}

		return nil, resState
	}
	defer rows.Close()

	var collections []*models.Collection
	for rows.Next() {
		var collection models.Collection
		err := rows.Scan(
			&collection.Id,
			&collection.Name,
			&collection.CreationDate,
		)
		if err != nil {
			resState := bl.CreateError(bl.DatabaseError, "GetAllUserCollections", "data_access")
			cr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
			return nil, resState
		}

		collections = append(collections, &collection)
	}

	if err := rows.Err(); err != nil {
		resState := bl.CreateError(bl.DatabaseError, "GetAllUserCollections", "data_access")
		cr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		return nil, resState
	}

	resState := bl.CreateError(bl.Ok, "GetAllUserCollections", "data_access")
	cr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelInfo, nil)
	return collections, resState
}

func (cr *CollectionRepository) AddCollection(collection *models.Collection) (int, *bl.MyError) {
	cr.MyLogger.WriteLog("AddCollection is called (Repo)", slog.LevelInfo, nil)

	if collection == nil {
		resState := bl.CreateError(bl.ErrInParameter, "AddCollection", "data_access")
		cr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		return 0, resState
	}

	db := cr.DbConfigs.DB
	schemaName := cr.DbConfigs.SchemaName
	query := fmt.Sprintf(addCollectionQuery, schemaName)
	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)

	if err != nil {
		resState := bl.CreateError(bl.DatabaseError, "AddCollection", "data_access")
		cr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		return 0, resState
	}
	defer deferTransaction(err, tx)

	err = tx.QueryRowContext(ctx, query, collection.Name, collection.CreationDate).Scan(collection.Id)

	if err != nil {
		var resState *bl.MyError

		if errors.Is(err, sql.ErrNoRows) {
			resState = bl.CreateError(bl.NoSuchColl, "AddCollection", "data_access")
			cr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		} else {
			resState := bl.CreateError(bl.DatabaseError, "AddCollection", "data_access")
			cr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		}

		return 0, resState
	}

	resState := bl.CreateError(bl.Ok, "AddCollection", "data_access")
	cr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelInfo, nil)
	return collection.Id, resState
}

func (cr *CollectionRepository) DeleteCollection(id int) *bl.MyError {
	cr.MyLogger.WriteLog("DeleteCollection is called (Repo)", slog.LevelInfo, nil)

	if id < 0 {
		resState := bl.CreateError(bl.ErrInParameter, "DeleteCollection", "data_access")
		cr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		return resState
	}

	db := cr.DbConfigs.DB
	schemaName := cr.DbConfigs.SchemaName
	ctx := context.Background()
	query := fmt.Sprintf(deleteCollectionQuery, schemaName, schemaName)

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		resState := bl.CreateError(bl.DatabaseError, "DeleteCollection", "data_access")
		cr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		return resState
	}
	defer deferTransaction(err, tx)

	_, err = tx.ExecContext(ctx, query, id)

	if err != nil {
		var resState *bl.MyError

		if errors.Is(err, sql.ErrNoRows) {
			resState = bl.CreateError(bl.NoSuchColl, "DeleteCollection", "data_access")
			cr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		} else {
			resState := bl.CreateError(bl.DatabaseError, "DeleteCollection", "data_access")
			cr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		}

		return resState
	}

	resState := bl.CreateError(bl.Ok, "DeleteCollection", "data_access")
	cr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelInfo, nil)
	return resState
}

func (cr *CollectionRepository) UpdateCollection(collection *models.Collection) *bl.MyError {
	cr.MyLogger.WriteLog("UpdateCollection is called (Repo)", slog.LevelInfo, nil)

	if collection == nil {
		resState := bl.CreateError(bl.ErrInParameter, "UpdateCollection", "data_access")
		cr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		return resState
	}

	db := cr.DbConfigs.DB
	schemaName := cr.DbConfigs.SchemaName
	query := fmt.Sprintf(updateCollectionQuery, schemaName)
	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		resState := bl.CreateError(bl.DatabaseError, "UpdateCollection", "data_access")
		cr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		return resState
	}
	defer deferTransaction(err, tx)

	_, err = tx.ExecContext(ctx, query, collection.Name, collection.CreationDate, collection.Id)

	if err != nil {
		var resState *bl.MyError

		if errors.Is(err, sql.ErrNoRows) {
			resState = bl.CreateError(bl.NoSuchColl, "UpdateCollection", "data_access")
			cr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		} else {
			resState := bl.CreateError(bl.DatabaseError, "UpdateCollection", "data_access")
			cr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		}

		return resState
	}

	resState := bl.CreateError(bl.Ok, "UpdateCollection", "data_access")
	cr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelInfo, nil)
	return resState
}

func (cr *CollectionRepository) GetAllNotesInCollection(collection *models.Collection) ([]*models.Note, *bl.MyError) {
	cr.MyLogger.WriteLog("GetAllNotesInCollection is called (Repo)", slog.LevelInfo, nil)

	if collection == nil {
		resState := bl.CreateError(bl.ErrInParameter, "GetAllNotesInCollection", "data_access")
		cr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		return nil, resState
	}

	db := cr.DbConfigs.DB
	schemaName := cr.DbConfigs.SchemaName
	query := fmt.Sprintf(getAllNotesInCollectionQuery, schemaName, schemaName)
	ctx := context.Background()

	rows, err := db.QueryContext(ctx, query, collection.Id)

	if err != nil {
		var resState *bl.MyError

		if errors.Is(err, sql.ErrNoRows) {
			resState = bl.CreateError(bl.NoSuchColl, "GetAllNotesInCollection", "data_access")
			cr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		} else {
			resState := bl.CreateError(bl.DatabaseError, "GetAllNotesInCollection", "data_access")
			cr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
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
			resState := bl.CreateError(bl.DatabaseError, "GetAllNotesInCollection", "data_access")
			cr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
			return nil, resState
		}

		notes = append(notes, &note)
	}

	if err := rows.Err(); err != nil {
		resState := bl.CreateError(bl.DatabaseError, "GetAllNotesInCollection", "data_access")
		cr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		return nil, resState
	}

	resState := bl.CreateError(bl.Ok, "GetAllNotesInCollection", "data_access")
	cr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelInfo, nil)
	return notes, resState
}
