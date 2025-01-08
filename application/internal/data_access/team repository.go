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

func (tr *TeamRepository) GetTeamByID(id int) (*models.Team, *bl.MyError) {
	tr.MyLogger.WriteLog("GetTeamByID is called (Repo)", slog.LevelInfo, nil)

	if id < 0 {
		resState := bl.CreateError(bl.ErrInParameter, "GetTeamByID", "data_access")
		tr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		return nil, resState
	}

	var team models.Team
	db := tr.DbConfigs.DB
	schemaName := tr.DbConfigs.SchemaName
	query := fmt.Sprintf(getTeamByIDQuery, schemaName)
	ctx := context.Background()

	err := db.QueryRowContext(ctx, query, id).Scan(
		&team.Id,
		&team.Name,
		&team.RegistrationDate,
	)

	if err != nil {
		var resState *bl.MyError

		if errors.Is(err, sql.ErrNoRows) {
			resState = bl.CreateError(bl.NoSuchTeam, "GetTeamByID", "data_access")
			tr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		} else {
			resState := bl.CreateError(bl.DatabaseError, "GetTeamByID", "data_access")
			tr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		}

		return nil, resState
	}

	resState := bl.CreateError(bl.Ok, "GetAllSections", "data_access")
	return &team, resState
}

func (tr *TeamRepository) GetTeamByName(name string) (*models.Team, *bl.MyError) {
	tr.MyLogger.WriteLog("GetTeamByName is called (Repo)", slog.LevelInfo, nil)

	if name == "" {
		resState := bl.CreateError(bl.ErrInParameter, "GetTeamByName", "data_access")
		tr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		return nil, resState
	}

	var team models.Team
	db := tr.DbConfigs.DB
	schemaName := tr.DbConfigs.SchemaName
	query := fmt.Sprintf(getTeamByNameQuery, schemaName)
	ctx := context.Background()

	err := db.QueryRowContext(ctx, query, name).Scan(
		&team.Id,
		&team.Name,
		&team.RegistrationDate,
	)

	if err != nil {
		var resState *bl.MyError

		if errors.Is(err, sql.ErrNoRows) {
			resState = bl.CreateError(bl.NoSuchTeam, "GetTeamByName", "data_access")
			tr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		} else {
			resState := bl.CreateError(bl.DatabaseError, "GetTeamByName", "data_access")
			tr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		}

		return nil, resState
	}

	resState := bl.CreateError(bl.Ok, "GetTeamByName", "data_access")
	tr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelInfo, nil)
	return &team, resState
}

func (tr *TeamRepository) GetTeamBySectionID(id int) (*models.Team, *bl.MyError) {
	tr.MyLogger.WriteLog("GetTeamBySectionID is called (Repo)", slog.LevelInfo, nil)

	if id < 0 {
		resState := bl.CreateError(bl.NoSuchTeam, "GetTeamBySectionID", "data_access")
		tr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		return nil, resState
	}

	var team models.Team
	db := tr.DbConfigs.DB
	schemaName := tr.DbConfigs.SchemaName
	query := fmt.Sprintf(getTeamBySectionIDQuery, schemaName, schemaName)
	ctx := context.Background()

	err := db.QueryRowContext(ctx, query, id).Scan(
		&team.Id,
		&team.Name,
		&team.RegistrationDate,
	)

	if err != nil {
		var resState *bl.MyError

		if errors.Is(err, sql.ErrNoRows) {
			resState = bl.CreateError(bl.NoSuchTeam, "GetTeamBySectionID", "data_access")
			tr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		} else {
			resState := bl.CreateError(bl.DatabaseError, "GetTeamBySectionID", "data_access")
			tr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		}

		return nil, resState
	}

	resState := bl.CreateError(bl.Ok, "GetTeamBySectionID", "data_access")
	tr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelInfo, nil)
	return &team, resState
}

func (tr *TeamRepository) GetAllTeams() ([]*models.Team, *bl.MyError) {
	tr.MyLogger.WriteLog("GetAllTeams is called (Repo)", slog.LevelInfo, nil)

	db := tr.DbConfigs.DB
	schemaName := tr.DbConfigs.SchemaName
	query := fmt.Sprintf(getAllTeamsQuery, schemaName)
	ctx := context.Background()

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		var resState *bl.MyError

		if errors.Is(err, sql.ErrNoRows) {
			resState = bl.CreateError(bl.EmptyResult, "GetAllTeams", "data_access")
			tr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		} else {
			resState := bl.CreateError(bl.DatabaseError, "GetAllTeams", "data_access")
			tr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		}

		return nil, resState
	}
	defer rows.Close()

	var teams []*models.Team
	for rows.Next() {
		var team models.Team
		err := rows.Scan(
			&team.Id,
			&team.Name,
			&team.RegistrationDate,
		)

		if err != nil {
			resState := bl.CreateError(bl.DatabaseError, "GetAllTeams", "data_access")
			tr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
			return nil, resState
		}

		teams = append(teams, &team)
	}

	if err := rows.Err(); err != nil {
		resState := bl.CreateError(bl.DatabaseError, "GetAllTeams", "data_access")
		tr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		return nil, resState
	}

	resState := bl.CreateError(bl.Ok, "GetAllTeams", "data_access")
	tr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelInfo, nil)
	return teams, resState
}

func (tr *TeamRepository) AddTeam(team *models.Team) (int, *bl.MyError) {
	tr.MyLogger.WriteLog("AddTeam is called (Repo)", slog.LevelInfo, nil)

	if team == nil {
		resState := bl.CreateError(bl.ErrInParameter, "AddTeam", "data_access")
		tr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		return 0, resState
	}

	db := tr.DbConfigs.DB
	schemaName := tr.DbConfigs.SchemaName
	query := fmt.Sprintf(addTeamQuery, schemaName)
	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		resState := bl.CreateError(bl.DatabaseError, "AddTeam", "data_access")
		tr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		return 0, resState
	}
	defer deferTransaction(err, tx)

	err = tx.QueryRowContext(ctx, query,
		team.Name,
		team.RegistrationDate,
	).Scan(team.Id)

	if err != nil {
		var resState *bl.MyError

		if errors.Is(err, sql.ErrNoRows) {
			resState = bl.CreateError(bl.NoSuchTeam, "AddTeam", "data_access")
			tr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		} else {
			resState := bl.CreateError(bl.OperationError, "AddTeam", "data_access")
			tr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		}

		return 0, resState
	}

	resState := bl.CreateError(bl.Ok, "AddTeam", "data_access")
	tr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelInfo, nil)
	return team.Id, resState
}

func (tr *TeamRepository) DeleteTeam(teamID int) *bl.MyError {
	tr.MyLogger.WriteLog("DeleteTeam is called (Repo)", slog.LevelInfo, nil)

	if teamID < 0 {
		resState := bl.CreateError(bl.ErrInParameter, "DeleteTeam", "data_access")
		tr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		return resState
	}

	db := tr.DbConfigs.DB
	schemaName := tr.DbConfigs.SchemaName
	ctx := context.Background()
	query := fmt.Sprintf(deleteTeamQuery, schemaName)

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		resState := bl.CreateError(bl.DatabaseError, "DeleteTeam", "data_access")
		tr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		return resState
	}
	defer deferTransaction(err, tx)

	_, err = tx.ExecContext(ctx, query, teamID)
	if err != nil {
		var resState *bl.MyError

		if errors.Is(err, sql.ErrNoRows) {
			resState = bl.CreateError(bl.NoSuchTeam, "DeleteTeam", "data_access")
			tr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		} else {
			resState := bl.CreateError(bl.OperationError, "DeleteTeam", "data_access")
			tr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		}

		return resState
	}

	resState := bl.CreateError(bl.Ok, "DeleteTeam", "data_access")
	tr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelInfo, nil)
	return resState
}

func (tr *TeamRepository) AddUserToTeam(uid int, tid int) *bl.MyError {
	tr.MyLogger.WriteLog("AddUserToTeam is called (Repo)", slog.LevelInfo, nil)

	if uid < 0 || tid < 0 {
		resState := bl.CreateError(bl.ErrInParameter, "AddUserToTeam", "data_access")
		tr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		return resState
	}

	db := tr.DbConfigs.DB
	schemaName := tr.DbConfigs.SchemaName
	query := fmt.Sprintf(addUserToTeamQuery, schemaName)
	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		resState := bl.CreateError(bl.DatabaseError, "AddUserToTeam", "data_access")
		tr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		return resState
	}
	defer deferTransaction(err, tx)

	_, err = tx.ExecContext(ctx, query, tid, uid)
	if err != nil {
		resState := bl.CreateError(bl.DatabaseError, "AddUserToTeam", "data_access")
		tr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		return resState
	}

	resState := bl.CreateError(bl.Ok, "AddUserToTeam", "data_access")
	tr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelInfo, nil)
	return resState
}

func (tr *TeamRepository) DeleteUserFromTeam(uid int, tid int) *bl.MyError {
	tr.MyLogger.WriteLog("DeleteUserFromTeam is called (Repo)", slog.LevelInfo, nil)

	if uid < 0 || tid < 0 {
		resState := bl.CreateError(bl.ErrInParameter, "DeleteUserFromTeam", "data_access")
		tr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		return resState
	}

	db := tr.DbConfigs.DB
	schemaName := tr.DbConfigs.SchemaName
	query := fmt.Sprintf(deleteUserFromTeamQuery, schemaName)
	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		resState := bl.CreateError(bl.DatabaseError, "DeleteUserFromTeam", "data_access")
		tr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		return resState
	}
	defer deferTransaction(err, tx)

	_, err = tx.ExecContext(ctx, query, uid, tid)
	if err != nil {
		resState := bl.CreateError(bl.DatabaseError, "DeleteUserFromTeam", "data_access")
		tr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		return resState
	}

	resState := bl.CreateError(bl.Ok, "DeleteUserFromTeam", "data_access")
	tr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelInfo, nil)
	return resState
}

func (tr *TeamRepository) UpdateTeam(team *models.Team) *bl.MyError {
	tr.MyLogger.WriteLog("UpdateTeam is called (Repo)", slog.LevelInfo, nil)

	if team == nil {
		resState := bl.CreateError(bl.ErrInParameter, "UpdateTeam", "data_access")
		tr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		return resState
	}

	db := tr.DbConfigs.DB
	schemaName := tr.DbConfigs.SchemaName
	query := fmt.Sprintf(updateTeamQuery, schemaName)
	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		resState := bl.CreateError(bl.DatabaseError, "UpdateTeam", "data_access")
		tr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		return resState
	}
	defer deferTransaction(err, tx)

	_, err = tx.ExecContext(ctx, query,
		team.Name,
		team.RegistrationDate,
		team.Id,
	)

	if err != nil {
		resState := bl.CreateError(bl.DatabaseError, "UpdateTeam", "data_access")
		tr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		return resState
	}

	resState := bl.CreateError(bl.Ok, "UpdateTeam", "data_access")
	tr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelInfo, nil)
	return resState
}

func (tr *TeamRepository) GetTeamMembers(teamID int) ([]*models.User, *bl.MyError) {
	tr.MyLogger.WriteLog("GetTeamMembers is called (Repo)", slog.LevelInfo, nil)

	if teamID < 0 {
		resState := bl.CreateError(bl.ErrInParameter, "GetTeamMembers", "data_access")
		tr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		return nil, resState
	}

	db := tr.DbConfigs.DB
	schemaName := tr.DbConfigs.SchemaName
	query := fmt.Sprintf(getTeamMembersQuery, schemaName, schemaName)
	ctx := context.Background()

	rows, err := db.QueryContext(ctx, query, teamID)
	if err != nil {
		var resState *bl.MyError

		if errors.Is(err, sql.ErrNoRows) {
			resState = bl.CreateError(bl.EmptyResult, "GetTeamMembers", "data_access")
			tr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		} else {
			resState = bl.CreateError(bl.DatabaseError, "GetTeamMembers", "data_access")
			tr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		}

		return nil, resState
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
			resState := bl.CreateError(bl.DatabaseError, "GetTeamMembers", "data_access")
			tr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
			return nil, resState
		}

		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		resState := bl.CreateError(bl.DatabaseError, "GetTeamMembers", "data_access")
		tr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		return nil, resState
	}

	resState := bl.CreateError(bl.Ok, "GetTeamMembers", "data_access")
	tr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelInfo, nil)
	return users, resState
}

func (tr *TeamRepository) GetUserTeam(user *models.User) (*models.Team, *bl.MyError) {
	tr.MyLogger.WriteLog("GetUserTeam is called (Repo)", slog.LevelInfo, nil)

	db := tr.DbConfigs.DB
	schemaName := tr.DbConfigs.SchemaName
	query := fmt.Sprintf(getUserTeamQuery, schemaName, schemaName)
	ctx := context.Background()
	team := models.Team{}

	err := db.QueryRowContext(ctx, query, user.Id).Scan(
		&team.Id,
		&team.Name,
		&team.RegistrationDate,
	)

	if err != nil {
		var resState *bl.MyError

		if errors.Is(err, sql.ErrNoRows) {
			resState = bl.CreateError(bl.NoSuchTeam, "GetUserTeam", "data_access")
			tr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelError, mylogger.LogCallerInfo())
		} else {
			resState = bl.CreateError(bl.DatabaseError, "GetUserTeam", "data_access")
			tr.MyLogger.WriteLog(resState.ConcatenateWithExternalErr(err), slog.LevelError, mylogger.LogCallerInfo())
		}

		return nil, resState
	}

	resState := bl.CreateError(bl.Ok, "GetUserTeam", "data_access")
	tr.MyLogger.WriteLog(resState.ConcatenateFields(), slog.LevelInfo, nil)
	return &team, resState
}
