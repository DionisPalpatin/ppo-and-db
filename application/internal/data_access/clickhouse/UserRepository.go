package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	bl "github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/business_logic"
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/models"
)

func (ur *UserRepository) GetUserByID(id int) (*models.User, *bl.MyError) {
	ur.MyLogger.WriteLog("GetUserByID is called (Repo)", slog.LevelInfo, nil)

	if id == 0 {
		myErr := bl.CreateError(bl.ErrGetUserByID, bl.ErrGetUserByIDError(), "GetUserByID")
		ur.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return nil, myErr
	}

	var user models.User
	db := ur.DbConfigs.DB
	schemaName := ur.DbConfigs.Name
	query := fmt.Sprintf("SELECT * FROM %s.users WHERE id = ?", schemaName)
	ctx := context.Background()

	err := db.QueryRowContext(ctx, query, id).Scan(
		&user.Id,
		&user.Fio,
		&user.RegistrationDate,
		&user.Login,
		&user.Password,
		&user.Role,
	)

	if err != nil {
		var myErr *bl.MyError

		if errors.Is(err, sql.ErrNoRows) {
			myErr = bl.CreateError(bl.ErrGetUserByID, bl.ErrGetUserByIDError(), "GetUserByID")
		} else {
			myErr = bl.CreateError(bl.ErrGetSectionByID, err, "GetSectionByID")
		}

		ur.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return nil, myErr
	}

	myOk := bl.CreateError(bl.AllIsOk, nil, "")
	return &user, myOk
}

func (ur *UserRepository) GetUserByLogin(login string) (*models.User, *bl.MyError) {
	ur.MyLogger.WriteLog("GetUserByLogin is called (Repo)", slog.LevelInfo, nil)

	if login == "" {
		myErr := bl.CreateError(bl.ErrGetUserByLogin, bl.ErrGetUserByLoginError(), "GetUserByLogin")
		ur.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return nil, myErr
	}

	var user models.User
	db := ur.DbConfigs.DB
	schemaName := ur.DbConfigs.Name
	query := fmt.Sprintf("SELECT * FROM %s.users WHERE login = ? OR fio = ?", schemaName)
	ctx := context.Background()

	err := db.QueryRowContext(ctx, query, login).Scan(
		&user.Id,
		&user.Fio,
		&user.RegistrationDate,
		&user.Login,
		&user.Password,
		&user.Role,
	)

	if err != nil {
		var myErr *bl.MyError

		if errors.Is(err, sql.ErrNoRows) {
			myErr = bl.CreateError(bl.ErrGetUserByLogin, bl.ErrGetUserByLoginError(), "GetUserByLogin")
		} else {
			myErr = bl.CreateError(bl.ErrGetUserByLogin, err, "GetUserByLogin")
		}

		ur.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return nil, myErr
	}

	myOk := bl.CreateError(bl.AllIsOk, nil, "")
	return &user, myOk
}

func (ur *UserRepository) GetAllUsers() ([]*models.User, *bl.MyError) {
	ur.MyLogger.WriteLog("GetAllUsers is called (Repo)", slog.LevelInfo, nil)

	db := ur.DbConfigs.DB
	schemaName := ur.DbConfigs.Name
	query := fmt.Sprintf("SELECT * FROM %s.users", schemaName)
	ctx := context.Background()

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		myErr := bl.CreateError(bl.ErrGetAllUsers, err, "GetAllUsers")
		ur.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return nil, myErr
	}

	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.Id,
			&user.Fio,
			&user.RegistrationDate,
			&user.Login,
			&user.Password,
			&user.Role,
		)

		if err != nil {
			myErr := bl.CreateError(bl.ErrGetAllUsers, err, "GetAllUsers")
			ur.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
			return nil, myErr
		}

		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		myErr := bl.CreateError(bl.ErrGetAllUsers, err, "GetAllUsers")
		ur.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return nil, myErr
	}

	myOk := bl.CreateError(bl.AllIsOk, nil, "")
	return users, myOk
}

func (ur *UserRepository) AddUser(user *models.User) *bl.MyError {
	ur.MyLogger.WriteLog("AddUserToTeam is called (Repo)", slog.LevelInfo, nil)

	if user == nil {
		myErr := bl.CreateError(bl.ErrAddUser, bl.ErrAddUserError(), "AddUserToTeam")
		ur.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}

	db := ur.DbConfigs.DB
	schemaName := ur.DbConfigs.Name
	query := fmt.Sprintf("INSERT INTO %s.users (fio, registration_date, login, password, role) VALUES (?, ?, ?, ?, ?)", schemaName)
	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return bl.CreateError(bl.ErrAddUser, err, "AddUserToTeam")
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	_, err = tx.ExecContext(ctx, query,
		user.Fio,
		user.RegistrationDate,
		user.Login,
		user.Password,
		user.Role,
	)

	if err != nil {
		myErr := bl.CreateError(bl.ErrAddUser, err, "AddUserToTeam")
		ur.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}

	myOk := bl.CreateError(bl.AllIsOk, nil, "")
	return myOk
}

func (ur *UserRepository) DeleteUser(id int) *bl.MyError {
	ur.MyLogger.WriteLog("DeleteUser is called (Repo)", slog.LevelInfo, nil)

	if id == 0 {
		myErr := bl.CreateError(bl.ErrDeleteUser, bl.ErrDeleteUserError(), "DeleteUser")
		ur.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}

	db := ur.DbConfigs.DB
	schemaName := ur.DbConfigs.Name
	ctx := context.Background()

	query1 := fmt.Sprintf("ALTER TABLE %s.notes_collections DELETE WHERE collection_id IN (SELECT id FROM %s.collections WHERE owner_id = ?)", schemaName, schemaName)
	query2 := fmt.Sprintf("ALTER TABLE %s.collections DELETE WHERE owner_id = ?", schemaName)
	query3 := fmt.Sprintf("ALTER TABLE %s.team_members DELETE WHERE user_id = ?", schemaName)
	query4 := fmt.Sprintf("ALTER TABLE %s.texts DELETE WHERE note_id IN (SELECT id FROM %s.notes WHERE owner_id = ?)", schemaName, schemaName)
	query5 := fmt.Sprintf("ALTER TABLE %s.images DELETE WHERE note_id IN (SELECT id FROM %s.notes WHERE owner_id = ?)", schemaName, schemaName)
	query6 := fmt.Sprintf("ALTER TABLE %s.raw_datas DELETE WHERE note_id IN (SELECT id FROM %s.notes WHERE owner_id = ?)", schemaName, schemaName)
	query7 := fmt.Sprintf("ALTER TABLE %s.notes DELETE WHERE owner_id = ?", schemaName)
	query8 := fmt.Sprintf("ALTER TABLE %s.users DELETE WHERE id = ?", schemaName)
	result_query := fmt.Sprintf("%s; %s; %s; %s; %s; %s; %s; %s;", query1, query2, query3, query4, query5, query6, query7, query8)

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return bl.CreateError(bl.ErrDeleteUser, err, "DeleteUser")
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
		myErr := bl.CreateError(bl.ErrDeleteUser, err, "DeleteUser")
		ur.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}

	myOk := bl.CreateError(bl.AllIsOk, nil, "")
	return myOk
}

func (ur *UserRepository) UpdateUser(user *models.User) *bl.MyError {
	ur.MyLogger.WriteLog("UpdateUser is called (Repo)", slog.LevelInfo, nil)

	if user == nil {
		myErr := bl.CreateError(bl.ErrUpdateUser, bl.ErrUpdateUserError(), "UpdateUser")
		ur.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}

	db := ur.DbConfigs.DB
	schemaName := ur.DbConfigs.Name
	query := fmt.Sprintf("ALTER TABLE %s.users UPDATE fio = ?, registration_date = ?, login = ?, password = ?, role = ? WHERE id = ?", schemaName)
	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return bl.CreateError(bl.ErrUpdateUser, err, "UpdateUser")
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	_, err = tx.ExecContext(ctx, query,
		user.Fio,
		user.RegistrationDate,
		user.Login,
		user.Password,
		user.Role,
		user.Id,
	)

	if err != nil {
		myErr := bl.CreateError(bl.ErrUpdateUser, err, "UpdateUser")
		ur.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}

	myOk := bl.CreateError(bl.AllIsOk, nil, "")
	return myOk
}
