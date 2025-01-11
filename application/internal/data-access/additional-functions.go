package data_access

import (
	"context"
	"database/sql"

	"github.com/DionisPalpatin/ppo-and-db/application/internal/models"
)

// ----------------------------------------------------------------------------
// Вспомогательные функции
// ----------------------------------------------------------------------------

func addOrUpdateNoteInfo(note *models.Note, db *sql.DB, ctx context.Context, query string, add bool) error {
	tx, err := db.BeginTx(ctx, nil)
	defer deferTransaction(err, tx)

	var content_type int = 1

	if err == nil {
		if add {
			err = tx.QueryRowContext(ctx, query,
				note.Access,
				note.Name,
				content_type,
				note.Likes,
				note.Dislikes,
				note.RegistrationDate,
				note.OwnerID,
				note.SectionID,
			).Scan(&note.Id)
		} else {
			_, err = tx.ExecContext(ctx, query,
				note.Access,
				note.Name,
				content_type,
				note.Likes,
				note.Dislikes,
				note.RegistrationDate,
				note.OwnerID,
				note.SectionID,
				note.Id,
			)
		}
		
	}

	return err
}

func addOrUpdateNoteContent(noteId int, data []byte, fileExt string, db *sql.DB, ctx context.Context, query string) error {
	tx, err := db.BeginTx(ctx, nil)
	defer deferTransaction(err, tx)

	if err == nil {
		_, err = tx.ExecContext(ctx, query, data, fileExt, noteId)
	}

	return err
}

func deferTransaction(err error, tx *sql.Tx) {
	if err != nil {
		_ = tx.Rollback()
	} else {
		_ = tx.Commit()
	}
}
