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

func (tr *TeamRepository) GetTeamByID(id int) (*models.Team, *bl.MyError) {
	tr.MyLogger.WriteLog("GetTeamByID is called (Repo)", slog.LevelInfo, nil)

	if id < 0 {
		myErr := bl.CreateError(bl.ErrGetTeamByID, bl.ErrGetTeamByIDError(), "GetTeamByID")
		tr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return nil, myErr
	}

	var team models.Team
	db := tr.DbConfigs.DB
	schemaName := tr.DbConfigs.SchemaName
	query := fmt.Sprintf("SELECT * FROM %s.teams WHERE id = $1", schemaName)
	ctx := context.Background()

	err := db.QueryRowContext(ctx, query, id).Scan(
		&team.Id,
		&team.Name,
		&team.RegistrationDate,
	)

	if err != nil {
		var myErr *bl.MyError

		if errors.Is(err, sql.ErrNoRows) {
			myErr = bl.CreateError(bl.ErrGetTeamByID, bl.ErrGetTeamByIDError(), "GetTeamByID")
		} else {
			myErr = bl.CreateError(bl.ErrGetTeamByID, err, "GetTeamByID")
		}

		tr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return nil, myErr
	}

	myOk := bl.CreateError(bl.AllIsOk, nil, "")
	return &team, myOk
}

func (tr *TeamRepository) GetTeamByName(name string) (*models.Team, *bl.MyError) {
	tr.MyLogger.WriteLog("GetTeamByName is called (Repo)", slog.LevelInfo, nil)

	if name == "" {
		myErr := bl.CreateError(bl.ErrGetTeamByName, bl.ErrGetTeamByNameError(), "GetTeamByName")
		tr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return nil, myErr
	}

	var team models.Team
	db := tr.DbConfigs.DB
	schemaName := tr.DbConfigs.SchemaName
	query := fmt.Sprintf("SELECT * FROM %s.teams WHERE name = $1", schemaName)
	ctx := context.Background()

	err := db.QueryRowContext(ctx, query, name).Scan(
		&team.Id,
		&team.Name,
		&team.RegistrationDate,
	)

	if err != nil {
		var myErr *bl.MyError

		if errors.Is(err, sql.ErrNoRows) {
			myErr = bl.CreateError(bl.ErrGetTeamByName, bl.ErrGetTeamByNameError(), "GetTeamByName")
		} else {
			myErr = bl.CreateError(bl.ErrGetTeamByName, err, "GetTeamByName")
		}

		tr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return nil, myErr
	}

	myOk := bl.CreateError(bl.AllIsOk, nil, "")
	return &team, myOk
}

func (tr *TeamRepository) GetTeamBySectionID(id int) (*models.Team, *bl.MyError) {
	tr.MyLogger.WriteLog("GetTeamBySectionID is called (Repo)", slog.LevelInfo, nil)

	if id < 0 {
		myErr := bl.CreateError(bl.ErrGetTeamBySectionID, bl.ErrGetTeamBySectionIDError(), "GetTeamBySectionID")
		tr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return nil, myErr
	}

	var team models.Team
	db := tr.DbConfigs.DB
	schemaName := tr.DbConfigs.SchemaName
	query := fmt.Sprintf("SELECT * FROM %s.teams t JOIN %s.teams_sections ts ON t.id = ts.team_id WHERE ts.section_id = $1", schemaName, schemaName)
	ctx := context.Background()

	err := db.QueryRowContext(ctx, query, id).Scan(
		&team.Id,
		&team.Name,
		&team.RegistrationDate,
	)

	if err != nil {
		var myErr *bl.MyError

		if errors.Is(err, sql.ErrNoRows) {
			myErr = bl.CreateError(bl.ErrGetTeamBySectionID, bl.ErrGetTeamBySectionIDError(), "GetTeamBySectionID")
		} else {
			myErr = bl.CreateError(bl.ErrGetTeamBySectionID, err, "GetTeamBySectionID")
		}

		tr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return nil, myErr
	}

	myOk := bl.CreateError(bl.AllIsOk, nil, "")
	return &team, myOk
}

func (tr *TeamRepository) GetAllTeams() ([]*models.Team, *bl.MyError) {
	tr.MyLogger.WriteLog("GetAllTeams is called (Repo)", slog.LevelInfo, nil)

	db := tr.DbConfigs.DB
	schemaName := tr.DbConfigs.SchemaName
	query := fmt.Sprintf("SELECT * FROM %s.teams", schemaName)
	ctx := context.Background()

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		myErr := bl.CreateError(bl.ErrGetAllTeams, err, "GetAllTeams")
		tr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return nil, myErr
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
			myErr := bl.CreateError(bl.ErrGetAllTeams, err, "GetAllTeams")
			tr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
			return nil, myErr
		}

		teams = append(teams, &team)
	}

	if err := rows.Err(); err != nil {
		myErr := bl.CreateError(bl.ErrGetAllTeams, err, "GetAllTeams")
		tr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return nil, myErr
	}

	myOk := bl.CreateError(bl.AllIsOk, nil, "")
	return teams, myOk
}

func (tr *TeamRepository) AddTeam(team *models.Team) *bl.MyError {
	tr.MyLogger.WriteLog("AddTeam is called (Repo)", slog.LevelInfo, nil)

	if team == nil {
		myErr := bl.CreateError(bl.ErrAddTeam, bl.ErrAddTeamError(), "AddTeam")
		tr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}

	db := tr.DbConfigs.DB
	schemaName := tr.DbConfigs.SchemaName
	query := fmt.Sprintf("INSERT INTO %s.teams (name, registration_date) VALUES ($1, $2)", schemaName)
	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return bl.CreateError(bl.ErrAddTeam, err, "AddTeam")
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	_, err = tx.ExecContext(ctx, query,
		team.Name,
		team.RegistrationDate,
	)

	if err != nil {
		myErr := bl.CreateError(bl.ErrAddTeam, err, "AddTeam")
		tr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}

	myOk := bl.CreateError(bl.AllIsOk, nil, "")
	return myOk
}

func (tr *TeamRepository) DeleteTeam(team_id int) *bl.MyError {
	tr.MyLogger.WriteLog("DeleteTeam is called (Repo)", slog.LevelInfo, nil)

	if team_id < 0 {
		myErr := bl.CreateError(bl.ErrDeleteTeam, bl.ErrDeleteTeamError(), "DeleteTeam")
		tr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}

	db := tr.DbConfigs.DB
	schemaName := tr.DbConfigs.SchemaName
	ctx := context.Background()

	// Запросы для удаления и изменения всех затрагиваемых таблиц
	query0 := fmt.Sprintf("SELECT section_id FROM %s.teams_sections WHERE team_id = $1", schemaName)
	query1 := fmt.Sprintf("UPDATE %s.notes SET section_id = 0 WHERE section_id = $2", schemaName, schemaName)
	query2 := fmt.Sprintf("DELETE FROM %s.teams_sections WHERE team_id = $1", schemaName)
	query3 := fmt.Sprintf("DELETE FROM %s.team_members WHERE team_id = $1", schemaName)
	query4 := fmt.Sprintf("DELETE FROM %s.teams WHERE id = $1", schemaName)
	query5 := fmt.Sprintf("DELETE FROM %s.sections WHERE id = $1", schemaName)
	result_query := fmt.Sprintf("%s; %s; %s; %s; %s;", query1, query2, query3, query4, query5)

	sec_id := 0
	err := db.QueryRowContext(ctx, query0, team_id).Scan(&sec_id)

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return bl.CreateError(bl.ErrDeleteTeam, err, "DeleteTeam")
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	_, err = tx.ExecContext(ctx, result_query, team_id, sec_id)
	if err != nil {
		myErr := bl.CreateError(bl.ErrDeleteTeam, err, "DeleteTeam")
		tr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}

	myOk := bl.CreateError(bl.AllIsOk, nil, "")
	return myOk
}

func (tr *TeamRepository) AddUserToTeam(uid int, tid int) *bl.MyError {
	tr.MyLogger.WriteLog("AddUserToTeam is called (Repo)", slog.LevelInfo, nil)

	if uid < 0 || tid < 0 {
		myErr := bl.CreateError(bl.ErrAddUserToTeam, bl.ErrAddUserError(), "AddUserToTeam")
		tr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}

	db := tr.DbConfigs.DB
	schemaName := tr.DbConfigs.SchemaName
	query := fmt.Sprintf("INSERT INTO %s.team_members (team_id, user_id) values ($1, $2)", schemaName)
	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return bl.CreateError(bl.ErrAddUserToTeam, err, "AddUserToTeam")
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	_, err = tx.ExecContext(ctx, query, tid, uid)
	if err != nil {
		myErr := bl.CreateError(bl.ErrAddUserToTeam, err, "AddUserToTeam")
		tr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}

	myOk := bl.CreateError(bl.AllIsOk, nil, "")
	return myOk
}

func (tr *TeamRepository) DeleteUserFromTeam(uid int, tid int) *bl.MyError {
	tr.MyLogger.WriteLog("DeleteUserFromTeam is called (Repo)", slog.LevelInfo, nil)

	if uid < 0 || tid < 0 {
		myErr := bl.CreateError(bl.ErrDeleteUserFromTeam, bl.ErrAddUserError(), "DeleteUserFromTeam")
		tr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}

	db := tr.DbConfigs.DB
	schemaName := tr.DbConfigs.SchemaName
	query := fmt.Sprintf("DELETE FROM %s.teams_sections WHERE user_id = $1 AND team_id = $2", schemaName)
	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return bl.CreateError(bl.ErrDeleteUserFromTeam, err, "DeleteUserFromTeam")
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	_, err = tx.ExecContext(ctx, query, uid, tid)
	if err != nil {
		myErr := bl.CreateError(bl.ErrDeleteUserFromTeam, err, "DeleteUserFromTeam")
		tr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}

	myOk := bl.CreateError(bl.AllIsOk, nil, "")
	return myOk
}

func (tr *TeamRepository) UpdateTeam(team *models.Team) *bl.MyError {
	tr.MyLogger.WriteLog("UpdateTeam is called (Repo)", slog.LevelInfo, nil)

	if team == nil {
		myErr := bl.CreateError(bl.ErrUpdateTeam, bl.ErrUpdateTeamError(), "UpdateTeam")
		tr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}

	db := tr.DbConfigs.DB
	schemaName := tr.DbConfigs.SchemaName
	query := fmt.Sprintf("UPDATE %s.teams SET name = $1, registration_date = $2 WHERE id = $3", schemaName)
	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return bl.CreateError(bl.ErrUpdateTeam, err, "UpdateTeam")
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	_, err = tx.ExecContext(ctx, query,
		team.Name,
		team.RegistrationDate,
		team.Id,
	)

	if err != nil {
		myErr := bl.CreateError(bl.ErrUpdateTeam, err, "UpdateTeam")
		tr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return myErr
	}

	myOk := bl.CreateError(bl.AllIsOk, nil, "")
	return myOk
}

func (tr *TeamRepository) GetTeamMembers(teamID int) ([]*models.User, *bl.MyError) {
	tr.MyLogger.WriteLog("GetTeamMembers is called (Repo)", slog.LevelInfo, nil)

	if teamID < 0 {
		myErr := bl.CreateError(bl.ErrGetTeamMembers, bl.ErrGetTeamMembersError(), "GetTeamMembers")
		tr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return nil, myErr
	}

	db := tr.DbConfigs.DB
	schemaName := tr.DbConfigs.SchemaName
	query := fmt.Sprintf("SELECT * FROM %s.users u JOIN %s.team_members tm ON u.id = tm.user_id WHERE tm.team_id = $1", schemaName, schemaName)
	ctx := context.Background()

	rows, err := db.QueryContext(ctx, query, teamID)
	if err != nil {
		return nil, &bl.MyError{ErrNum: bl.ErrGetTeamMembers, FuncName: "GetTeamMembers", Err: err}
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
			myErr := bl.CreateError(bl.ErrGetTeamMembers, err, "GetTeamMembers")
			tr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
			return nil, myErr
		}

		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		myErr := bl.CreateError(bl.ErrGetTeamMembers, err, "GetTeamMembers")
		tr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return nil, myErr
	}

	myOk := bl.CreateError(bl.AllIsOk, nil, "")
	return users, myOk
}

func (tr *TeamRepository) GetUserTeam(user *models.User) (*models.Team, *bl.MyError) {
	tr.MyLogger.WriteLog("GetUserTeam is called (Repo)", slog.LevelInfo, nil)

	db := tr.DbConfigs.DB
	schemaName := tr.DbConfigs.SchemaName
	query := fmt.Sprintf("SELECT * FROM %s.teams t JOIN %s.team_members tm ON t.id = tm.team_id WHERE tm.user_id = $1", schemaName, schemaName)
	ctx := context.Background()
	team := models.Team{}

	err := db.QueryRowContext(ctx, query, user.Id).Scan(
		&team.Id,
		&team.Name,
		&team.RegistrationDate,
	)

	if err != nil {
		var myErr *bl.MyError

		if errors.Is(err, sql.ErrNoRows) {
			myErr = bl.CreateError(bl.ErrGetUserTeam, bl.ErrGetUserTeamError(), "GetUserTeam")
		} else {
			myErr = bl.CreateError(bl.ErrGetUserTeam, err, "GetUserTeam")
		}

		tr.MyLogger.WriteLog(myErr.Err.Error(), slog.LevelError, nil)
		return nil, myErr
	}

	myOk := bl.CreateError(bl.AllIsOk, nil, "")
	return &team, myOk
}
