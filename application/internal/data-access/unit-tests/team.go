package da_unit

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"

	"github.com/DionisPalpatin/ppo-and-db/application/internal/config"
	mylogger "github.com/DionisPalpatin/ppo-and-db/application/internal/logger"

	bl "github.com/DionisPalpatin/ppo-and-db/application/internal/business-logic"
	da "github.com/DionisPalpatin/ppo-and-db/application/internal/data-access"

	"github.com/DionisPalpatin/ppo-and-db/application/internal/models"
	data_builders "github.com/DionisPalpatin/ppo-and-db/application/internal/models/builders"
)

type TeamRepositorySuite struct {
	suite.Suite

	db   *sql.DB
	mock sqlmock.Sqlmock
	repo da.TeamRepository
	ctx  context.Context
}

func (s *TeamRepositorySuite) BeforeEach(t provider.T) {
	var err error

	s.db, s.mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}

	confs := config.DBConfigs {
		DB: s.db,
		SchemaName: "test",
	}
	logger := mylogger.MyLogger{}
	logger.InitLogger("./logs/unit-tests-team.log", "debug")

	s.repo = da.NewTeamRepository(&confs, &logger)
	s.ctx = context.Background()
}

func (s *TeamRepositorySuite) AfterEach(t provider.T) {
	s.db.Close()
}

func TestTeamRepositorySuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(TeamRepositorySuite))
}


func (s *TeamRepositorySuite) Test_TeamRepository_GetTeamByID_Success(t provider.T) {
	t.Title("GetTeamByID: Success")
	t.Tags("TeamRepository")

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		query := fmt.Sprintf(getTeamByIDQuery, s.repo.DbConfigs.SchemaName)
		s.mock.ExpectQuery(query).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "registration_date"}).
				AddRow(1, "Team A", ZEROTIME))

		expectedTeam := data_builders.NewTeamBuilder().WithId(1).WithName("Team A").Build()

		team, err := s.repo.GetTeamByID(1)

		sCtx.Assert().Equal(err.ErrNum, bl.Ok)
		sCtx.Assert().Equal(expectedTeam.Id, team.Id)
	})
}

func (s *TeamRepositorySuite) Test_TeamRepository_GetTeamByID_Failure(t provider.T) {
	t.Title("GetTeamByID: Failure")
	t.Tags("TeamRepository")

	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		query := fmt.Sprintf(getTeamByIDQuery, s.repo.DbConfigs.SchemaName)
		s.mock.ExpectQuery(query).
			WithArgs(1).
			WillReturnError(sql.ErrConnDone)

		team, err := s.repo.GetTeamByID(1)

		sCtx.Assert().Nil(team)
		sCtx.Assert().Equal(err.ErrNum, bl.DatabaseError)
	})
}

func (s *TeamRepositorySuite) Test_TeamRepository_GetTeamByName_Success(t provider.T) {
	t.Title("GetTeamByName: Success")
	t.Tags("TeamRepository")

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		query := fmt.Sprintf(getTeamByNameQuery, s.repo.DbConfigs.SchemaName)
		s.mock.ExpectQuery(query).
			WithArgs("Team A").
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "registration_date"}).
				AddRow(1, "Team A", ZEROTIME))

		expectedTeam := data_builders.NewTeamBuilder().WithId(1).WithName("Team A").Build()

		team, err := s.repo.GetTeamByName("Team A")

		sCtx.Assert().Equal(err.ErrNum, bl.Ok)
		sCtx.Assert().Equal(expectedTeam.Id, team.Id)
	})
}

func (s *TeamRepositorySuite) Test_TeamRepository_GetTeamByName_Failure(t provider.T) {
	t.Title("GetTeamByName: Failure")
	t.Tags("TeamRepository")

	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		query := fmt.Sprintf(getTeamByNameQuery, s.repo.DbConfigs.SchemaName)
		s.mock.ExpectQuery(query).
			WithArgs("Team A").
			WillReturnError(sql.ErrConnDone)

		team, err := s.repo.GetTeamByName("Team A")

		sCtx.Assert().Nil(team)
		sCtx.Assert().Equal(err.ErrNum, bl.DatabaseError)
	})
}

func (s *TeamRepositorySuite) Test_TeamRepository_GetAllTeams_Success(t provider.T) {
	t.Title("GetAllTeams: Success")
	t.Tags("TeamRepository")

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		query := fmt.Sprintf(getAllTeamsQuery, s.repo.DbConfigs.SchemaName)
		s.mock.ExpectQuery(query).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "registration_date"}).
				AddRow(1, "Team A", ZEROTIME).
				AddRow(2, "Team B", ZEROTIME))

		expectedTeams := []*models.Team{
			data_builders.NewTeamBuilder().WithId(1).WithName("Team A").Build(),
			data_builders.NewTeamBuilder().WithId(2).WithName("Team B").Build(),
		}

		teams, err := s.repo.GetAllTeams()

		sCtx.Assert().Equal(err.ErrNum, bl.Ok)
		sCtx.Assert().Len(teams, 2)
		sCtx.Assert().Equal(expectedTeams[0].Id, teams[0].Id)
	})
}

func (s *TeamRepositorySuite) Test_TeamRepository_GetAllTeams_Failure(t provider.T) {
	t.Title("GetAllTeams: Failure")
	t.Tags("TeamRepository")

	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		query := fmt.Sprintf(getAllTeamsQuery, s.repo.DbConfigs.SchemaName)
		s.mock.ExpectQuery(query).
			WillReturnError(sql.ErrConnDone)

		teams, err := s.repo.GetAllTeams()

		sCtx.Assert().Nil(teams)
		sCtx.Assert().Equal(err.ErrNum, bl.DatabaseError)
	})
}

func (s *TeamRepositorySuite) Test_TeamRepository_AddTeam_Success(t provider.T) {
	t.Title("AddTeam: Success")
	t.Tags("TeamRepository")

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		query := fmt.Sprintf(addTeamQuery, s.repo.DbConfigs.SchemaName)

		s.mock.ExpectBegin()
		s.mock.ExpectQuery(query).
			WithArgs("Team A", ZEROTIME).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		s.mock.ExpectCommit()	
			
		team := data_builders.NewTeamBuilder().WithName("Team A").WithRegistrationDate(ZEROTIME).Build()

		id, err := s.repo.AddTeam(team)

		sCtx.Assert().Equal(err.ErrNum, bl.Ok)
		sCtx.Assert().Equal(1, id)
	})
}

func (s *TeamRepositorySuite) Test_TeamRepository_AddTeam_Failure(t provider.T) {
	t.Title("AddTeam: Failure")
	t.Tags("TeamRepository")

	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		query := fmt.Sprintf(addTeamQuery, s.repo.DbConfigs.SchemaName)
		
		s.mock.ExpectBegin()
		s.mock.ExpectQuery(query).
			WithArgs("Team A", ZEROTIME).
			WillReturnError(sql.ErrConnDone)
		s.mock.ExpectRollback()

		team := data_builders.NewTeamBuilder().WithName("Team A").WithRegistrationDate(ZEROTIME).Build()

		_, err := s.repo.AddTeam(team)

		sCtx.Assert().Equal(err.ErrNum, bl.DatabaseError)
	})
}

func (s *TeamRepositorySuite) Test_TeamRepository_DeleteTeam_Success(t provider.T) {
	t.Title("DeleteTeam: Success")
	t.Tags("TeamRepository")

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		query := fmt.Sprintf(deleteTeamQuery, s.repo.DbConfigs.SchemaName)

		s.mock.ExpectBegin()
		s.mock.ExpectExec(query).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 1))
		s.mock.ExpectCommit()	

		err := s.repo.DeleteTeam(1)

		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *TeamRepositorySuite) Test_TeamRepository_DeleteTeam_Failure(t provider.T) {
	t.Title("DeleteTeam: Failure")
	t.Tags("TeamRepository")

	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		query := fmt.Sprintf(deleteTeamQuery, s.repo.DbConfigs.SchemaName)

		s.mock.ExpectBegin()
		s.mock.ExpectExec(query).
			WithArgs(1).
			WillReturnError(sql.ErrConnDone)
		s.mock.ExpectRollback()	

		err := s.repo.DeleteTeam(1)

		sCtx.Assert().Equal(bl.DatabaseError, err.ErrNum)
	})
}

func (s *TeamRepositorySuite) Test_TeamRepository_AddUserToTeam_Success(t provider.T) {
	t.Title("AddUserToTeam: Success")
	t.Tags("TeamRepository")

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		query := fmt.Sprintf(addUserToTeamQuery, s.repo.DbConfigs.SchemaName)

		s.mock.ExpectBegin()
		s.mock.ExpectExec(query).
			WithArgs(2, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))
		s.mock.ExpectCommit()

		err := s.repo.AddUserToTeam(1, 2)

		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *TeamRepositorySuite) Test_TeamRepository_AddUserToTeam_Failure(t provider.T) {
	t.Title("AddUserToTeam: Failure")
	t.Tags("TeamRepository")

	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		query := fmt.Sprintf(addUserToTeamQuery, s.repo.DbConfigs.SchemaName)

		s.mock.ExpectBegin()
		s.mock.ExpectExec(query).
			WithArgs(1, 2).
			WillReturnError(sql.ErrConnDone)
		s.mock.ExpectRollback()

		err := s.repo.AddUserToTeam(1, 2)

		sCtx.Assert().Equal(bl.DatabaseError, err.ErrNum)
	})
}

func (s *TeamRepositorySuite) Test_TeamRepository_DeleteUserFromTeam_Success(t provider.T) {
	t.Title("DeleteUserFromTeam: Success")
	t.Tags("TeamRepository")

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		query := fmt.Sprintf(deleteUserFromTeamQuery, s.repo.DbConfigs.SchemaName)

		s.mock.ExpectBegin()
		s.mock.ExpectExec(query).
			WithArgs(1, 2).
			WillReturnResult(sqlmock.NewResult(0, 1))
		s.mock.ExpectCommit()	

		err := s.repo.DeleteUserFromTeam(1, 2)

		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *TeamRepositorySuite) Test_TeamRepository_DeleteUserFromTeam_Failure(t provider.T) {
	t.Title("DeleteUserFromTeam: Failure")
	t.Tags("TeamRepository")

	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		query := fmt.Sprintf(deleteUserFromTeamQuery, s.repo.DbConfigs.SchemaName)

		s.mock.ExpectBegin()
		s.mock.ExpectExec(query).
			WithArgs(1, 2).
			WillReturnError(sql.ErrConnDone)
		s.mock.ExpectRollback()	

		err := s.repo.DeleteUserFromTeam(1, 2)

		sCtx.Assert().Equal(bl.DatabaseError, err.ErrNum)
	})
}

func (s *TeamRepositorySuite) Test_TeamRepository_UpdateTeam_Success(t provider.T) {
	t.Title("UpdateTeam: Success")
	t.Tags("TeamRepository")

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		query := fmt.Sprintf(updateTeamQuery, s.repo.DbConfigs.SchemaName)
		team := data_builders.NewTeamBuilder().WithId(1).WithName("New Name").WithRegistrationDate(ZEROTIME).Build()

		s.mock.ExpectBegin()
		s.mock.ExpectExec(query).
			WithArgs(team.Name, team.RegistrationDate, team.Id).
			WillReturnResult(sqlmock.NewResult(0, 1))
		s.mock.ExpectCommit()

		err := s.repo.UpdateTeam(team)

		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *TeamRepositorySuite) Test_TeamRepository_UpdateTeam_Failure(t provider.T) {
	t.Title("UpdateTeam: Failure")
	t.Tags("TeamRepository")

	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		query := fmt.Sprintf(updateTeamQuery, s.repo.DbConfigs.SchemaName)
		team := data_builders.NewTeamBuilder().WithId(1).Build()

		s.mock.ExpectBegin()
		s.mock.ExpectExec(query).
			WithArgs(team.Name, team.RegistrationDate, team.Id).
			WillReturnError(sql.ErrConnDone)
		s.mock.ExpectRollback()	

		err := s.repo.UpdateTeam(team)

		sCtx.Assert().Equal(bl.DatabaseError, err.ErrNum)
	})
}

func (s *TeamRepositorySuite) Test_TeamRepository_GetTeamMembers_Success(t provider.T) {
	t.Title("GetTeamMembers: Success")
	t.Tags("TeamRepository")

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		query := fmt.Sprintf(getTeamMembersQuery, s.repo.DbConfigs.SchemaName, s.repo.DbConfigs.SchemaName)
		s.mock.ExpectQuery(query).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "fio", "registration_date", "login", "password", "role"}).
				AddRow(1, "John Doe", ZEROTIME, "johndoe", "password", bl.Reader))

		members, err := s.repo.GetTeamMembers(1)

		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
		sCtx.Assert().Len(members, 1)
	})
}

func (s *TeamRepositorySuite) Test_TeamRepository_GetTeamMembers_Failure(t provider.T) {
	t.Title("GetTeamMembers: Failure")
	t.Tags("TeamRepository")

	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		query := fmt.Sprintf(getTeamMembersQuery, s.repo.DbConfigs.SchemaName, s.repo.DbConfigs.SchemaName)
		s.mock.ExpectQuery(query).
			WithArgs(1).
			WillReturnError(sql.ErrConnDone)

		members, err := s.repo.GetTeamMembers(1)

		sCtx.Assert().Nil(members)
		sCtx.Assert().Equal(bl.DatabaseError, err.ErrNum)
	})
}

func (s *TeamRepositorySuite) Test_TeamRepository_GetUserTeam_Success(t provider.T) {
	t.Title("GetUserTeam: Success")
	t.Tags("TeamRepository")

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		query := fmt.Sprintf(getUserTeamQuery, s.repo.DbConfigs.SchemaName, s.repo.DbConfigs.SchemaName)
		s.mock.ExpectQuery(query).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "registration_date"}).
				AddRow(1, "Team A", ZEROTIME))

		team, err := s.repo.GetUserTeam(&models.User{Id: 1})

		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
		sCtx.Assert().NotNil(team)
	})
}

func (s *TeamRepositorySuite) Test_TeamRepository_GetUserTeam_Failure(t provider.T) {
	t.Title("GetUserTeam: Failure")
	t.Tags("TeamRepository")

	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		query := fmt.Sprintf(getUserTeamQuery, s.repo.DbConfigs.SchemaName, s.repo.DbConfigs.SchemaName)
		s.mock.ExpectQuery(query).
			WithArgs(1).
			WillReturnError(sql.ErrConnDone)

		team, err := s.repo.GetUserTeam(&models.User{Id: 1})

		sCtx.Assert().Nil(team)
		sCtx.Assert().Equal(bl.DatabaseError, err.ErrNum)
	})
}
