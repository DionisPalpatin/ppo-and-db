package data_access

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	bl "github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/business_logic"
	mylogger "github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/logger"
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/models"
)

func (ur *UserRepository) GetUserByID(id int) (*models.User, *bl.MyError) {
	ur.MyLogger.WriteLog("GetUserByID is called (Repo)", slog.LevelInfo, nil)

	if id < 0 {
		resState := bl.CreateError(bl.NoSuchUser, "GetUserByID", "data_access")
		ur.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		return nil, resState
	}

	var user models.User
	db := ur.DbConfigs.DB
	schemaName := ur.DbConfigs.SchemaName
	ctx := context.Background()

	query := fmt.Sprintf(getUserByIDQuery, schemaName)
	err := db.QueryRowContext(ctx, query, id).Scan(
		&user.Id,
		&user.Fio,
		&user.RegistrationDate,
		&user.Login,
		&user.Password,
		&user.Role,
	)

	if err != nil {
		var resState *bl.MyError

		if errors.Is(err, sql.ErrNoRows) {
			resState = bl.CreateError(bl.NoSuchUser, "GetUserByID", "data_access")
		} else {
			resState = bl.CreateError(bl.DatabaseError, "GetUserByID", "data_access")
		}

		ur.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		return nil, resState
	}

	resState := bl.CreateError(bl.Ok, "GetUserByID", "data_access")
	ur.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelInfo, nil)
	return &user, resState
}

func (ur *UserRepository) GetUserByLogin(loginOrFio string) (*models.User, *bl.MyError) {
	ur.MyLogger.WriteLog("GetUserByLogin is called (Repo)", slog.LevelInfo, nil)

	if loginOrFio == "" {
		resState := bl.CreateError(bl.NoSuchUser, "GetUserByLogin", "data_access")
		ur.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		return nil, resState
	}

	var user models.User
	db := ur.DbConfigs.DB
	schemaName := ur.DbConfigs.SchemaName
	ctx := context.Background()

	query := fmt.Sprintf(getUserByStringQuery, schemaName)
	err := db.QueryRowContext(ctx, query, loginOrFio).Scan(
		&user.Id,
		&user.Fio,
		&user.RegistrationDate,
		&user.Login,
		&user.Password,
		&user.Role,
	)

	if err != nil {
		var resState *bl.MyError

		if errors.Is(err, sql.ErrNoRows) {
			resState = bl.CreateError(bl.NoSuchUser, "GetUserByLogin", "data_access")
		} else {
			resState = bl.CreateError(bl.NoSuchUser, "GetUserByLogin", "data_access")
		}

		ur.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		return nil, resState
	}

	resState := bl.CreateError(bl.Ok, "GetUserByLogin", "data_access")
	ur.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelInfo, nil)
	return &user, resState
}

func (ur *UserRepository) GetAllUsers() ([]*models.User, *bl.MyError) {
	ur.MyLogger.WriteLog("GetAllUsers is called (Repo)", slog.LevelInfo, nil)

	db := ur.DbConfigs.DB
	schemaName := ur.DbConfigs.SchemaName
	ctx := context.Background()

	query := fmt.Sprintf(getAllUsersQuery, schemaName)
	rows, err := db.QueryContext(ctx, query)
	defer rows.Close()
	if err != nil {
		resState := bl.CreateError(bl.DatabaseError, "GetAllUsers", "data_access")
		ur.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		return nil, resState
	}

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
			resState := bl.CreateError(bl.DatabaseError, "GetAllUsers", "data_access")
			ur.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
			return nil, resState
		}

		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		resState := bl.CreateError(bl.DatabaseError, "GetAllUsers", "data_access")
		ur.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		return nil, resState
	}

	resState := bl.CreateError(bl.Ok, "GetAllUsers", "data_access")
	ur.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelInfo, nil)
	return users, resState
}

func (ur *UserRepository) AddUser(user *models.User) *bl.MyError {
	ur.MyLogger.WriteLog("AddUser is called (Repo)", slog.LevelInfo, nil)

	if user == nil {
		resState := bl.CreateError(bl.OperationError, "AddUser", "data_access")
		ur.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		return resState
	}

	db := ur.DbConfigs.DB
	schemaName := ur.DbConfigs.SchemaName
	ctx := context.Background()

	query := fmt.Sprintf(addUserQuery, schemaName)
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return bl.CreateError(bl.DatabaseError, "AddUser", "data_access")
	}
	defer deferTransaction(err, tx)

	_, err = tx.ExecContext(ctx, query,
		user.Fio,
		user.RegistrationDate,
		user.Login,
		user.Password,
		user.Role,
	)

	if err != nil {
		resState := bl.CreateError(bl.DatabaseError, "AddUser", "data_access")
		ur.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		return resState
	}

	resState := bl.CreateError(bl.Ok, "AddUser", "data_access")
	ur.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelInfo, nil)
	return resState
}

func (ur *UserRepository) DeleteUser(id int) *bl.MyError {
	ur.MyLogger.WriteLog("DeleteUser is called (Repo)", slog.LevelInfo, nil)

	if id < 0 {
		resState := bl.CreateError(bl.OperationError, "DeleteUser", "data_access")
		ur.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		return resState
	}

	db := ur.DbConfigs.DB
	schemaName := ur.DbConfigs.SchemaName
	ctx := context.Background()

	query := fmt.Sprintf(deleteUserQuery, schemaName)
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return bl.CreateError(bl.DatabaseError, "DeleteUser", "data_access")
	}
	defer deferTransaction(err, tx)

	_, err = tx.ExecContext(ctx, query, id)
	if err != nil {
		resState := bl.CreateError(bl.DatabaseError, "DeleteUser", "data_access")
		ur.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		return resState
	}

	resState := bl.CreateError(bl.Ok, "DeleteUser", "data_access")
	ur.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelInfo, nil)
	return resState
}

func (ur *UserRepository) UpdateUser(user *models.User) *bl.MyError {
	ur.MyLogger.WriteLog("UpdateUser is called (Repo)", slog.LevelInfo, nil)

	if user == nil {
		resState := bl.CreateError(bl.OperationError, "UpdateUser", "data_access")
		ur.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		return resState
	}

	db := ur.DbConfigs.DB
	schemaName := ur.DbConfigs.SchemaName
	query := fmt.Sprintf(updateUserQuery, schemaName)
	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return bl.CreateError(bl.DatabaseError, "UpdateUser", "data_access")
	}
	defer deferTransaction(err, tx)

	_, err = tx.ExecContext(ctx, query,
		user.Fio,
		user.RegistrationDate,
		user.Login,
		user.Password,
		user.Role,
		user.Id,
	)

	if err != nil {
		resState := bl.CreateError(bl.DatabaseError, "UpdateUser", "data_access")
		ur.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		return resState
	}

	resState := bl.CreateError(bl.Ok, "UpdateUser", "data_access")
	ur.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelInfo, nil)
	return resState
}
