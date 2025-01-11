package interg_bl

import (
	integr_ut "github.com/DionisPalpatin/ppo-and-db/application/integration-tests/utils"
	bl "github.com/DionisPalpatin/ppo-and-db/application/internal/business-logic"
	"github.com/DionisPalpatin/ppo-and-db/application/internal/config"
	da "github.com/DionisPalpatin/ppo-and-db/application/internal/data-access"
	mylogger "github.com/DionisPalpatin/ppo-and-db/application/internal/logger"
	data_builders "github.com/DionisPalpatin/ppo-and-db/application/internal/models/builders"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"

	"database/sql"
	"fmt"
	"os"
)

type TeamServiceIntegrationTest struct {
	suite.Suite
	db          *sql.DB
	teamService bl.ITeamService
	userService bl.IUserService

	idTeamToDelete  int
	idTeamToUpdate  int
	idTeamToUserAdd int
	idTeamToUserDel int

	idUserToAdd    int
	idUserToDelete int
}

func (s *TeamServiceIntegrationTest) BeforeEach(t provider.T) {
	var err error
	pgInfo := integr_ut.PostgresInfo{
		// Host:     "postgres-int-tests",
		Host:     "localhost",
		User:     "postgres",
		Password: "password",
		Port:     15423,
		DBName:   "NotebookAppIntTests",
	}

	s.db, err = integr_ut.InitDB(pgInfo)
	if err != nil {
		fmt.Println(err)
		return
	}

	if s.db == nil {
		return
	}

	text, err := os.ReadFile("./init.sql")
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = s.db.Exec(string(text))
	if err != nil {
		fmt.Println(err)
		return
	}

	dbconfig := config.DBConfigs{
		Host:       pgInfo.Host,
		User:       pgInfo.User,
		Password:   pgInfo.Password,
		Port:       pgInfo.Port,
		Name:       pgInfo.DBName,
		DB:         s.db,
		SchemaName: "interg_tests",
	}

	logger := mylogger.MyLogger{}
	logger.InitLogger("./logs/integration-test-team.log", "debug")

	userRepo := da.NewUserRepository(&dbconfig, &logger)
	teamRepo := da.NewTeamRepository(&dbconfig, &logger)
	s.userService = bl.NewUserService(&userRepo, &logger)
	s.teamService = bl.NewTeamService(&teamRepo, &logger)

	s.idUserToAdd = 4
	s.idUserToDelete = 5

	s.idTeamToUpdate = 8
	s.idTeamToDelete = 11
	s.idTeamToUserAdd = 4
	s.idTeamToUserDel = 5
}

func (s *TeamServiceIntegrationTest) AfterEach(t provider.T) {
	_ = s.db.Close()
}

func (s *TeamServiceIntegrationTest) TestGetTeamSuccess(t provider.T) {
	t.Title("GetTeam: Success")
	t.Tags("TeamServiceIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		requester := data_builders.NewUserBuilder().WithRole(bl.Admin).Build()

		team, err := s.teamService.GetTeam(1, "", bl.SearchByID, requester)

		sCtx.Assert().NotNil(team)
		sCtx.Assert().NotNil(err)
		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *TeamServiceIntegrationTest) TestGetTeamNotFound(t provider.T) {
	t.Title("GetTeam: Team Not Found")
	t.Tags("TeamServiceIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Team not found", func(sCtx provider.StepCtx) {
		requester := data_builders.NewUserBuilder().WithRole(bl.Admin).Build()

		team, err := s.teamService.GetTeam(999, "", bl.SearchByID, requester)

		sCtx.Assert().Nil(team)
		sCtx.Assert().NotNil(err)
		sCtx.Assert().Equal(bl.NoSuchTeam, err.ErrNum)
	})
}

func (s *TeamServiceIntegrationTest) TestGetAllTeamsSuccess(t provider.T) {
	t.Title("GetAllTeams: Success")
	t.Tags("TeamServiceIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		requester := data_builders.NewUserBuilder().WithRole(bl.Admin).Build()

		// Запрос ко всем командам
		teams, err := s.teamService.GetAllTeams(requester)

		sCtx.Assert().NotNil(teams)
		sCtx.Assert().NotNil(err)
		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *TeamServiceIntegrationTest) TestGetAllTeamsForbidden(t provider.T) {
	t.Title("GetAllTeams: Forbidden")
	t.Tags("TeamServiceIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Forbidden access", func(sCtx provider.StepCtx) {
		requester := data_builders.NewUserBuilder().WithRole(bl.Reader).Build()

		teams, err := s.teamService.GetAllTeams(requester)

		sCtx.Assert().Nil(teams)
		sCtx.Assert().NotNil(err)
		sCtx.Assert().Equal(bl.ErrAccessDenied, err.ErrNum)
	})
}

func (s *TeamServiceIntegrationTest) TestUpdateTeamSuccess(t provider.T) {
	t.Title("UpdateTeam: Success")
	t.Tags("TeamServiceIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		requester := data_builders.NewUserBuilder().WithRole(bl.Admin).Build()
		teamToUpdate := data_builders.NewTeamBuilder().WithName("Updated Team").WithId(s.idTeamToUpdate).Build()

		err := s.teamService.UpdateTeam(requester, teamToUpdate)

		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *TeamServiceIntegrationTest) TestDeleteTeamSuccess(t provider.T) {
	t.Title("DeleteTeam: Success")
	t.Tags("TeamServiceIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		requester := data_builders.NewUserBuilder().WithRole(bl.Admin).Build()

		err := s.teamService.DeleteTeam(requester, s.idTeamToDelete)

		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *TeamServiceIntegrationTest) TestAddTeamSuccess(t provider.T) {
	t.Title("AddTeam: Success")
	t.Tags("TeamServiceIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		requester := data_builders.NewUserBuilder().WithRole(bl.Admin).Build()
		newTeam := data_builders.NewTeamBuilder().WithName("New Team").Build()

		_, err := s.teamService.AddTeam(requester, newTeam)

		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *TeamServiceIntegrationTest) TestAddTeamConflict(t provider.T) {
	t.Title("AddTeam: Conflict")
	t.Tags("TeamServiceIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Team already exists", func(sCtx provider.StepCtx) {
		requester := data_builders.NewUserBuilder().WithRole(bl.Admin).Build()
		existingTeam := data_builders.NewTeamBuilder().WithName("Alpha").WithId(1).Build()

		id, err := s.teamService.AddTeam(requester, existingTeam)

		sCtx.Assert().Equal(0, id)
		sCtx.Assert().NotNil(err)
		sCtx.Assert().Equal(bl.DatabaseError, err.ErrNum)
	})
}

func (s *TeamServiceIntegrationTest) TestAddUserToTeamSuccess(t provider.T) {
	t.Title("AddUserToTeam: Success")
	t.Tags("TeamServiceIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		requester := data_builders.NewUserBuilder().WithRole(bl.Admin).Build()

		err := s.teamService.AddUserToTeam(requester, s.idUserToAdd, s.idTeamToUserAdd)

		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *TeamServiceIntegrationTest) TestDeleteUserFromTeamSuccess(t provider.T) {
	t.Title("DeleteUserFromTeam: Success")
	t.Tags("TeamServiceIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		requester := data_builders.NewUserBuilder().WithRole(bl.Admin).Build()

		err := s.teamService.DeleteUserFromTeam(requester, s.idUserToDelete, s.idTeamToUserDel)

		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *TeamServiceIntegrationTest) TestGetTeamMembersSuccess(t provider.T) {
	t.Title("GetTeamMembers: Success")
	t.Tags("TeamServiceIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		requester := data_builders.NewUserBuilder().WithRole(bl.Admin).Build()

		members, err := s.teamService.GetTeamMembers(1, requester)

		sCtx.Assert().NotNil(members)
		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *TeamServiceIntegrationTest) TestGetTeamMembersNotFound(t provider.T) {
	t.Title("GetTeamMembers: Team Not Found")
	t.Tags("TeamServiceIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Team not found", func(sCtx provider.StepCtx) {
		requester := data_builders.NewUserBuilder().WithRole(bl.Admin).Build()

		members, err := s.teamService.GetTeamMembers(999, requester)

		sCtx.Assert().Nil(members)
		sCtx.Assert().NotNil(err)
		sCtx.Assert().Equal(bl.NoSuchTeam, err.ErrNum)
	})
}

func (s *TeamServiceIntegrationTest) TestGetUserTeamSuccess(t provider.T) {
	t.Title("GetUserTeam: Success")
	t.Tags("TeamServiceIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		requester := data_builders.NewUserBuilder().WithRole(bl.Admin).WithUserID(1).Build()

		team, err := s.teamService.GetUserTeam(requester)

		sCtx.Assert().NotNil(team)
		sCtx.Assert().NotNil(err)
		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *TeamServiceIntegrationTest) TestGetUserTeamNotFound(t provider.T) {
	t.Title("GetUserTeam: User Not Found")
	t.Tags("TeamServiceIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("User not found", func(sCtx provider.StepCtx) {
		requester := data_builders.NewUserBuilder().WithRole(bl.Admin).WithUserID(999).Build()

		team, err := s.teamService.GetUserTeam(requester)

		sCtx.Assert().Nil(team)
		sCtx.Assert().NotNil(err)
		sCtx.Assert().Equal(bl.NoSuchTeam, err.ErrNum)
	})
}
