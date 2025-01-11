package interg_da

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
)

type TeamRepoIntegrationTestSuite struct {
	suite.Suite
	db       *sql.DB
	teamRepo da.TeamRepository

	idTeamToDelete  int
	idTeamToUpdate  int
	idTeamToUserAdd int
	idTeamToUserDel int

	idUserToAdd    int
	idUserToDelete int
}

func (s *TeamRepoIntegrationTestSuite) BeforeEach(t provider.T) {
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

	s.teamRepo = da.NewTeamRepository(&dbconfig, &logger)

	s.idUserToAdd = 6
	s.idUserToDelete = 3

	s.idTeamToUpdate = 9
	s.idTeamToDelete = 10
	s.idTeamToUserAdd = 6
	s.idTeamToUserDel = 3
}

func (s *TeamRepoIntegrationTestSuite) AfterEach(t provider.T) {
	_ = s.db.Close()
}

func (s *TeamRepoIntegrationTestSuite) TestGetTeamByIDSuccess(t provider.T) {
	t.Title("GetTeamByID: Success")
	t.Tags("TeamRepoIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		team, err := s.teamRepo.GetTeamByID(1)

		sCtx.Assert().NotNil(team)
		sCtx.Assert().NotNil(err)
		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *TeamRepoIntegrationTestSuite) TestGetTeamByIDNotFound(t provider.T) {
	t.Title("GetTeamByID: Team Not Found")
	t.Tags("TeamRepoIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Team not found", func(sCtx provider.StepCtx) {
		team, err := s.teamRepo.GetTeamByID(999)

		sCtx.Assert().Nil(team)
		sCtx.Assert().NotNil(err)
		sCtx.Assert().Equal(bl.NoSuchTeam, err.ErrNum)
	})
}

func (s *TeamRepoIntegrationTestSuite) TestGetTeamByNameSuccess(t provider.T) {
	t.Title("GetTeamByID: Success")
	t.Tags("TeamRepoIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		team, err := s.teamRepo.GetTeamByName("Alpha")

		sCtx.Assert().NotNil(team)
		sCtx.Assert().NotNil(err)
		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *TeamRepoIntegrationTestSuite) TestGetTeamByNameNotFound(t provider.T) {
	t.Title("GetTeamByID: Team Not Found")
	t.Tags("TeamRepoIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Team not found", func(sCtx provider.StepCtx) {
		team, err := s.teamRepo.GetTeamByName("bukva")

		sCtx.Assert().Nil(team)
		sCtx.Assert().NotNil(err)
		sCtx.Assert().Equal(bl.NoSuchTeam, err.ErrNum)
	})
}

func (s *TeamRepoIntegrationTestSuite) TestGetAllTeamsSuccess(t provider.T) {
	t.Title("GetAllTeams: Success")
	t.Tags("TeamRepoIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		teams, err := s.teamRepo.GetAllTeams()

		sCtx.Assert().NotNil(teams)
		sCtx.Assert().NotNil(err)
		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *TeamRepoIntegrationTestSuite) TestUpdateTeamSuccess(t provider.T) {
	t.Title("UpdateTeam: Success")
	t.Tags("TeamRepoIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		teamToUpdate := data_builders.NewTeamBuilder().WithName("Updated Team").WithId(s.idTeamToUpdate).Build()

		err := s.teamRepo.UpdateTeam(teamToUpdate)

		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *TeamRepoIntegrationTestSuite) TestDeleteTeamSuccess(t provider.T) {
	t.Title("DeleteTeam: Success")
	t.Tags("TeamRepoIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		err := s.teamRepo.DeleteTeam(s.idTeamToDelete)

		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *TeamRepoIntegrationTestSuite) TestAddTeamSuccess(t provider.T) {
	t.Title("AddTeam: Success")
	t.Tags("TeamRepoIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		newTeam := data_builders.NewTeamBuilder().WithName("New Team").Build()

		_, err := s.teamRepo.AddTeam(newTeam)

		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *TeamRepoIntegrationTestSuite) TestAddTeamConflict(t provider.T) {
	t.Title("AddTeam: Conflict")
	t.Tags("TeamRepoIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Team already exists", func(sCtx provider.StepCtx) {
		existingTeam := data_builders.NewTeamBuilder().WithName("Alpha").WithId(1).Build()

		id, err := s.teamRepo.AddTeam(existingTeam)

		sCtx.Assert().Equal(0, id)
		sCtx.Assert().NotNil(err)
		sCtx.Assert().Equal(bl.DatabaseError, err.ErrNum)
	})
}

func (s *TeamRepoIntegrationTestSuite) TestAddUserToTeamSuccess(t provider.T) {
	t.Title("AddUserToTeam: Success")
	t.Tags("TeamRepoIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		err := s.teamRepo.AddUserToTeam(s.idUserToAdd, s.idTeamToUserAdd)

		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *TeamRepoIntegrationTestSuite) TestDeleteUserFromTeamSuccess(t provider.T) {
	t.Title("DeleteUserFromTeam: Success")
	t.Tags("TeamRepoIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		err := s.teamRepo.DeleteUserFromTeam(s.idUserToDelete, s.idTeamToUserDel)

		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *TeamRepoIntegrationTestSuite) TestGetTeamMembersSuccess(t provider.T) {
	t.Title("GetTeamMembers: Success")
	t.Tags("TeamRepoIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		members, err := s.teamRepo.GetTeamMembers(1)

		sCtx.Assert().NotNil(members)
		sCtx.Assert().NotNil(err)
		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *TeamRepoIntegrationTestSuite) TestGetTeamMembersNotFound(t provider.T) {
	t.Title("GetTeamMembers: Team Not Found")
	t.Tags("TeamRepoIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Team not found", func(sCtx provider.StepCtx) {
		members, err := s.teamRepo.GetTeamMembers(999)

		sCtx.Assert().Nil(members)
		sCtx.Assert().NotNil(err)
		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *TeamRepoIntegrationTestSuite) TestGetUserTeamSuccess(t provider.T) {
	t.Title("GetUserTeam: Success")
	t.Tags("TeamRepoIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		requester := data_builders.NewUserBuilder().WithRole(bl.Admin).WithUserID(1).Build()

		team, err := s.teamRepo.GetUserTeam(requester)

		sCtx.Assert().NotNil(team)
		sCtx.Assert().NotNil(err)
		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func (s *TeamRepoIntegrationTestSuite) TestGetUserTeamNotFound(t provider.T) {
	t.Title("GetUserTeam: User Not Found")
	t.Tags("TeamRepoIntegrationTest")

	// if integr_ut.IsUnitTestsFailed() {
	// 	t.Skip()
	// }

	t.WithNewStep("User not found", func(sCtx provider.StepCtx) {
		requester := data_builders.NewUserBuilder().WithRole(bl.Admin).WithUserID(999).Build()

		team, err := s.teamRepo.GetUserTeam(requester)

		sCtx.Assert().Nil(team)
		sCtx.Assert().NotNil(err)
		sCtx.Assert().Equal(bl.NoSuchTeam, err.ErrNum)
	})
}
