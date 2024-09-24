package clickhouse

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/business_logic"
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/models"
)

func (cr *CollectionRepository) GetCollectionByID(id int) (*models.Collection, *bl.MyError) {
	cr.MyLogger.WriteLog("GetCollectionByID is called (Repo)", slog.LevelInfo, nil)

	if id == 0 {
		myErr := bl.CreateError(bl.ErrGetCollectionByID, bl.ErrGetCollectionByIDError(), "GetCollectionByID")
		cr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return nil, myErr
	}

	var collection models.Collection
	db := cr.DbConfigs.DB
	schemaName := cr.DbConfigs.Name
	query := fmt.Sprintf("SELECT * FROM %s.collections WHERE id = ?", schemaName)
	ctx := context.Background()

	err := db.QueryRowContext(ctx, query, id).Scan(
		&collection.Id,
		&collection.Name,
		&collection.CreationDate,
	)

	if err != nil {
		var myErr *bl.MyError

		if errors.Is(err, sql.ErrNoRows) {
			myErr = bl.CreateError(bl.ErrGetCollectionByID, bl.ErrGetCollectionByIDError(), "GetCollectionByID")
		} else {
			myErr = bl.CreateError(bl.ErrGetCollectionByID, err, "GetCollectionByID")
		}

		cr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return nil, myErr
	}

	myOk := bl.CreateError(bl.AllIsOk, nil, "")
	return &collection, myOk
}

func (cr *CollectionRepository) GetCollectionByName(name string) (*models.Collection, *bl.MyError) {
	cr.MyLogger.WriteLog("GetCollectionByName is called (Repo)", slog.LevelInfo, nil)

	if name == "" {
		myErr := bl.CreateError(bl.ErrGetCollectionByName, bl.ErrGetCollectionByNameError(), "GetCollectionByName")
		cr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return nil, myErr
	}

	var collection models.Collection
	db := cr.DbConfigs.DB
	schemaName := cr.DbConfigs.Name
	query := fmt.Sprintf("SELECT * FROM %s.collections WHERE name = ?", schemaName)
	ctx := context.Background()

	err := db.QueryRowContext(ctx, query, name).Scan(
		&collection.Id,
		&collection.Name,
		&collection.CreationDate,
	)

	if err != nil {
		var myErr *bl.MyError

		if errors.Is(err, sql.ErrNoRows) {
			myErr = bl.CreateError(bl.ErrGetCollectionByName, bl.ErrGetCollectionByNameError(), "GetCollectionByName")
		} else {
			myErr = bl.CreateError(bl.ErrGetCollectionByName, err, "GetCollectionByName")
		}

		cr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return nil, myErr
	}

	myOk := bl.CreateError(bl.AllIsOk, nil, "")
	return &collection, myOk
}

func (cr *CollectionRepository) GetAllCollections() ([]*models.Collection, *bl.MyError) {
	cr.MyLogger.WriteLog("GetAllCollections is called (Repo)", slog.LevelInfo, nil)

	db := cr.DbConfigs.DB
	schemaName := cr.DbConfigs.Name
	query := fmt.Sprintf("SELECT * FROM %s.collections", schemaName)
	ctx := context.Background()

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		myErr := bl.CreateError(bl.ErrGetAllCollections, err, "GetAllCollections")
		cr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return nil, myErr
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
			myErr := bl.CreateError(bl.ErrGetAllCollections, err, "GetAllCollections")
			cr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
			return nil, myErr
		}

		collections = append(collections, &collection)
	}

	if err := rows.Err(); err != nil {
		myErr := bl.CreateError(bl.ErrGetAllCollections, err, "GetAllCollections")
		cr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return nil, myErr
	}

	myOk := bl.CreateError(bl.AllIsOk, nil, "")
	return collections, myOk
}

func (cr *CollectionRepository) GetAllUserCollections(user *models.User) ([]*models.Collection, *bl.MyError) {
	cr.MyLogger.WriteLog("GetAllUserCollections is called (Repo)", slog.LevelInfo, nil)

	if user == nil {
		myErr := bl.CreateError(bl.ErrGetAllUserCollections, bl.ErrGetAllUserCollectionsError(), "GetAllUserCollections")
		cr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return nil, myErr
	}

	db := cr.DbConfigs.DB
	schemaName := cr.DbConfigs.Name
	query := fmt.Sprintf("SELECT c.* FROM %s.collections c JOIN %s.notes n ON c.id = n.collection_id WHERE n.owner_id = ?", schemaName, schemaName)
	ctx := context.Background()

	rows, err := db.QueryContext(ctx, query, user.Id)

	if err != nil {
		myErr := bl.CreateError(bl.ErrGetAllUserCollections, err, "GetAllUserCollections")
		cr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return nil, myErr
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
			myErr := bl.CreateError(bl.ErrGetAllUserCollections, err, "GetAllUserCollections")
			cr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
			return nil, myErr
		}

		collections = append(collections, &collection)
	}

	if err := rows.Err(); err != nil {
		myErr := bl.CreateError(bl.ErrGetAllUserCollections, err, "GetAllUserCollections")
		cr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return nil, myErr
	}

	myOk := bl.CreateError(bl.AllIsOk, nil, "")
	return collections, myOk
}

func (cr *CollectionRepository) AddCollection(collection *models.Collection) *bl.MyError {
	cr.MyLogger.WriteLog("AddCollection is called (Repo)", slog.LevelInfo, nil)

	if collection == nil {
		myErr := bl.CreateError(bl.ErrAddCollection, bl.ErrAddCollectionError(), "AddCollection")
		cr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}

	db := cr.DbConfigs.DB
	schemaName := cr.DbConfigs.Name
	query := fmt.Sprintf("INSERT INTO %s.collections (name, creation_date) VALUES (?, ?)", schemaName)
	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return bl.CreateError(bl.ErrAddCollection, err, "AddCollection")
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	_, err = tx.ExecContext(ctx, query, collection.Name, collection.CreationDate)

	if err != nil {
		myErr := bl.CreateError(bl.ErrAddCollection, err, "AddCollection")
		cr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}

	myOk := bl.CreateError(bl.AllIsOk, nil, "")
	return myOk
}

func (cr *CollectionRepository) DeleteCollection(id int) *bl.MyError {
	cr.MyLogger.WriteLog("DeleteCollection is called (Repo)", slog.LevelInfo, nil)

	if id == 0 {
		myErr := bl.CreateError(bl.ErrDeleteCollection, bl.ErrDeleteCollectionError(), "DeleteCollection")
		cr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}

	db := cr.DbConfigs.DB
	schemaName := cr.DbConfigs.Name
	ctx := context.Background()
	query1 := fmt.Sprintf("ALTER TABLE %s.notes_collections DELETE WHERE collection_id = ?", schemaName)
	query2 := fmt.Sprintf("ALTER TABLE %s.collections DELETE WHERE id = ?", schemaName)
	result_query := fmt.Sprintf("%s; %s;", query1, query2)

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return bl.CreateError(bl.ErrDeleteCollection, err, "DeleteCollection")
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
		myErr := bl.CreateError(bl.ErrDeleteCollection, err, "DeleteCollection")
		cr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}

	myOk := bl.CreateError(bl.AllIsOk, nil, "")
	return myOk
}

func (cr *CollectionRepository) UpdateCollection(collection *models.Collection) *bl.MyError {
	cr.MyLogger.WriteLog("UpdateCollection is called (Repo)", slog.LevelInfo, nil)

	if collection == nil {
		myErr := bl.CreateError(bl.ErrUpdateCollection, bl.ErrUpdateCollectionError(), "UpdateCollection")
		cr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}

	db := cr.DbConfigs.DB
	schemaName := cr.DbConfigs.Name
	query := fmt.Sprintf("ALTER TABLE %s.collections UPDATE name = ?, creation_date = ? WHERE id = ?", schemaName)
	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return bl.CreateError(bl.ErrUpdateCollection, err, "UpdateCollection")
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	_, err = tx.ExecContext(ctx, query, collection.Name, collection.CreationDate, collection.Id)

	if err != nil {
		myErr := bl.CreateError(bl.ErrUpdateCollection, err, "UpdateCollection")
		cr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}

	myOk := bl.CreateError(bl.AllIsOk, nil, "")
	return myOk
}

func (cr *CollectionRepository) GetAllNotesInCollection(collection *models.Collection) ([]*models.Note, *bl.MyError) {
	cr.MyLogger.WriteLog("GetAllNotesInCollection is called (Repo)", slog.LevelInfo, nil)

	if collection == nil {
		myErr := bl.CreateError(bl.ErrGetAllNotesInCollection, bl.ErrGetAllNotesInCollectionError(), "GetAllNotesInCollection")
		cr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return nil, myErr
	}

	db := cr.DbConfigs.DB
	schemaName := cr.DbConfigs.Name
	query := fmt.Sprintf("SELECT * FROM %s.notes n JOIN %s.note_collections nc ON n.id = nc.note_id WHERE nc.collection_id = ?", schemaName, schemaName)
	ctx := context.Background()

	rows, err := db.QueryContext(ctx, query, collection.Id)

	if err != nil {
		myErr := bl.CreateError(bl.ErrGetAllNotesInCollection, err, "GetAllNotesInCollection")
		cr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
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
			myErr := bl.CreateError(bl.ErrGetAllNotesInCollection, err, "GetAllNotesInCollection")
			cr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
			return nil, myErr
		}

		notes = append(notes, &note)
	}

	if err := rows.Err(); err != nil {
		myErr := bl.CreateError(bl.ErrGetAllNotesInCollection, err, "GetAllNotesInCollection")
		cr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return nil, myErr
	}

	myOk := bl.CreateError(bl.AllIsOk, nil, "")
	return notes, myOk
}
