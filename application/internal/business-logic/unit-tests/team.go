package unit_tests

import (
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"

	mocks "github.com/DionisPalpatin/ppo-and-db/application/internal/database/mocks/v2"
	mylogger "github.com/DionisPalpatin/ppo-and-db/application/internal/logger"

	bl "github.com/DionisPalpatin/ppo-and-db/application/internal/business-logic"
	myerror_builder "github.com/DionisPalpatin/ppo-and-db/application/internal/business-logic/builder"

	"github.com/DionisPalpatin/ppo-and-db/application/internal/models"
	"github.com/DionisPalpatin/ppo-and-db/application/internal/models/builders"
)

type TeamServiceTestSuite struct {
	suite.Suite
}

type extDependensesTeam struct {
	teamRepo *mocks.MockTeamRepository
	logger   *mylogger.MyLogger
}

func initExtDependenciesTeam(t provider.T) *extDependensesTeam {
	teamRepo := mocks.NewMockTeamRepository(t)

	f := &extDependensesTeam{
		teamRepo: teamRepo,
		logger:   &mylogger.MyLogger{},
	}

	f.logger.InitLogger("./logs/unit-test-team.log", "debug")

	return f
}

func (s *TeamServiceTestSuite) TestTeamService_GetTeam_OK(t provider.T) {
	t.Title("Get team: OK")
	t.Tags("TeamService")
	t.Parallel()

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		df := initExtDependenciesTeam(t)
		inTeam := builders.NewTeamBuilder().Build()
		outTeam := builders.NewTeamBuilder().WithId(inTeam.Id)
		reqUser := builders.NewUserBuilder().WithRole(bl.Admin)
		myerr := myerror_builder.NewMyErrorBuilder().WithErrNum(bl.Ok).Build()

		df.teamRepo.EXPECT().GetTeamByID(inTeam.Id).Return(outTeam.Team, myerr).Once()
		TeamService := bl.NewTeamService(df.teamRepo, df.logger)

		team, err := TeamService.GetTeam(inTeam.Id, "", bl.SearchByID, reqUser.User)

		sCtx.Assert().Equal(err.ErrNum, bl.Ok)
		sCtx.Assert().Equal(team.Id, outTeam.Team.Id)
	})
}

func (s *TeamServiceTestSuite) TestTeamService_GetTeam_AccessDenied(t provider.T) {
	t.Title("Get user: Access Denied")
	t.Tags("TeamService")
	t.Parallel()

	t.WithNewStep("Access Denied", func(sCtx provider.StepCtx) {
		df := initExtDependenciesTeam(t)
		inTeam := builders.NewTeamBuilder().Build()
		reqUser := builders.NewUserBuilder().WithRole(bl.Reader)
		// myerr := myerror_builder.NewMyErrorBuilder().WithErrNum(bl.ErrAccessDenied).Build()

		// df.teamRepo.EXPECT().GetTeamByID(inTeam.Id).Return(nil, myerr).Once()
		teamService := bl.NewTeamService(df.teamRepo, df.logger)

		team, err := teamService.GetTeam(inTeam.Id, "", bl.SearchByID, reqUser.User)

		sCtx.Assert().Nil(team)
		sCtx.Assert().Equal(err.ErrNum, bl.ErrAccessDenied)
	})
}

func (s *TeamServiceTestSuite) TestTeamService_GetAllTeams_OK(t provider.T) {
	t.Title("Get all users: OK")
	t.Tags("TeamService")
	t.Parallel()

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		df := initExtDependenciesTeam(t)
		reqUser := builders.NewUserBuilder().WithRole(bl.Admin)
		users := []*models.Team{
			builders.NewTeamBuilder().Build(),
			builders.NewTeamBuilder().Build(),
		}
		myerr := myerror_builder.NewMyErrorBuilder().WithErrNum(bl.Ok).Build()

		df.teamRepo.EXPECT().GetAllTeams().Return(users, myerr).Once()
		teamSvc := bl.NewTeamService(df.teamRepo, df.logger)

		resultTeams, err := teamSvc.GetAllTeams(reqUser.User)

		sCtx.Assert().Equal(err.ErrNum, bl.Ok)
		sCtx.Assert().Equal(len(resultTeams), len(users))
	})
}

func (s *TeamServiceTestSuite) TestTeamService_GetAllTeams_AccessDenied(t provider.T) {
	t.Title("Get all users: Access Denied")
	t.Tags("TeamService")
	t.Parallel()

	t.WithNewStep("Access Denied", func(sCtx provider.StepCtx) {
		df := initExtDependenciesTeam(t)
		reqUser := builders.NewUserBuilder().WithRole(bl.Reader)
		// myerr := myerror_builder.NewMyErrorBuilder().WithErrNum(bl.ErrAccessDenied).Build()

		// df.teamRepo.EXPECT().GetAllTeams().Return(nil, myerr).Once()
		teamService := bl.NewTeamService(df.teamRepo, df.logger)

		teams, err := teamService.GetAllTeams(reqUser.User)

		sCtx.Assert().Nil(teams)
		sCtx.Assert().Equal(err.ErrNum, bl.ErrAccessDenied)
	})
}

func (s *TeamServiceTestSuite) TestTeamService_UpdateTeam_OK(t provider.T) {
	t.Title("Update user: OK")
	t.Tags("TeamService")
	t.Parallel()

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		df := initExtDependenciesTeam(t)
		inTeam := builders.NewTeamBuilder().Build()
		reqUser := builders.NewUserBuilder().WithRole(bl.Admin) // Админская роль
		myerr := myerror_builder.NewMyErrorBuilder().WithErrNum(bl.Ok).Build()

		df.teamRepo.EXPECT().UpdateTeam(inTeam).Return(myerr).Once()
		teamService := bl.NewTeamService(df.teamRepo, df.logger)

		err := teamService.UpdateTeam(reqUser.User, inTeam)

		sCtx.Assert().Equal(err.ErrNum, bl.Ok)
	})
}

func (s *TeamServiceTestSuite) TestTeamService_UpdateTeam_AccessDenied(t provider.T) {
	t.Title("Update user: Access Denied")
	t.Tags("TeamService")
	t.Parallel()

	t.WithNewStep("Access Denied", func(sCtx provider.StepCtx) {
		df := initExtDependenciesTeam(t)
		inTeam := builders.NewTeamBuilder().Build()
		reqUser := builders.NewUserBuilder().WithRole(bl.Reader)
		// myerr := myerror_builder.NewMyErrorBuilder().WithErrNum(bl.ErrAccessDenied).Build()

		// df.teamRepo.EXPECT().UpdateTeam(inTeam).Return(myerr).Once()
		teamService := bl.NewTeamService(df.teamRepo, df.logger)

		err := teamService.UpdateTeam(reqUser.User, inTeam)

		sCtx.Assert().Equal(err.ErrNum, bl.ErrAccessDenied)
	})
}

func (s *TeamServiceTestSuite) TestTeamService_DeleteTeam_OK(t provider.T) {
	t.Title("Delete user: OK")
	t.Tags("TeamService")
	t.Parallel()

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		df := initExtDependenciesTeam(t)
		inTeam := builders.NewTeamBuilder().Build()
		reqUser := builders.NewUserBuilder().WithRole(bl.Admin) // Админская роль
		myerr := myerror_builder.NewMyErrorBuilder().WithErrNum(bl.Ok).Build()

		df.teamRepo.EXPECT().DeleteTeam(inTeam.Id).Return(myerr).Once()
		teamService := bl.NewTeamService(df.teamRepo, df.logger)

		err := teamService.DeleteTeam(reqUser.User, inTeam.Id)

		sCtx.Assert().Equal(err.ErrNum, bl.Ok)
	})
}

func (s *TeamServiceTestSuite) TestTeamService_DeleteTeam_AccessDenied(t provider.T) {
	t.Title("Delete user: Access Denied")
	t.Tags("TeamService")
	t.Parallel()

	t.WithNewStep("Access Denied", func(sCtx provider.StepCtx) {
		df := initExtDependenciesTeam(t)
		inTeam := builders.NewTeamBuilder().Build()
		reqUser := builders.NewUserBuilder().WithRole(bl.Reader)
		// myerr := myerror_builder.NewMyErrorBuilder().WithErrNum(bl.ErrAccessDenied).Build()

		// df.teamRepo.EXPECT().DeleteTeam(inTeam.Id).Return(myerr).Once()
		teamService := bl.NewTeamService(df.teamRepo, df.logger)

		err := teamService.DeleteTeam(reqUser.User, inTeam.Id)

		sCtx.Assert().Equal(err.ErrNum, bl.ErrAccessDenied)
	})
}

func (s *TeamServiceTestSuite) TestTeamService_AddTeam_OK(t provider.T) {
	t.Title("Add team: OK")
	t.Tags("TeamService")
	t.Parallel()

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		df := initExtDependenciesTeam(t)
		requester := builders.NewUserBuilder().WithRole(bl.Admin)
		inTeam := builders.NewTeamBuilder().Build()
		myerr := myerror_builder.NewMyErrorBuilder().WithErrNum(bl.Ok).Build()

		df.teamRepo.EXPECT().AddTeam(inTeam).Return(myerr).Once()
		teamService := bl.NewTeamService(df.teamRepo, df.logger)

		teamID, err := teamService.AddTeam(requester.User, inTeam)

		sCtx.Assert().Equal(err.ErrNum, bl.Ok)
		sCtx.Assert().Equal(teamID, 0)
	})
}

func (s *TeamServiceTestSuite) TestTeamService_AddTeam_AccessDenied(t provider.T) {
	t.Title("Add team: Access Denied")
	t.Tags("TeamService")
	t.Parallel()

	t.WithNewStep("Access Denied", func(sCtx provider.StepCtx) {
		df := initExtDependenciesTeam(t)
		requester := builders.NewUserBuilder().WithRole(bl.Reader)
		inTeam := builders.NewTeamBuilder().Build()

		teamService := bl.NewTeamService(df.teamRepo, df.logger)
		teamID, err := teamService.AddTeam(requester.User, inTeam)

		sCtx.Assert().Equal(err.ErrNum, bl.ErrAccessDenied)
		sCtx.Assert().Equal(teamID, 0)
	})
}

func (s *TeamServiceTestSuite) TestTeamService_AddUserToTeam_OK(t provider.T) {
	t.Title("Add user to team: OK")
	t.Tags("TeamService")
	t.Parallel()

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		df := initExtDependenciesTeam(t)
		requester := builders.NewUserBuilder().WithRole(bl.Admin)
		userID, teamID := 1, 1
		myerr := myerror_builder.NewMyErrorBuilder().WithErrNum(bl.Ok).Build()

		df.teamRepo.EXPECT().AddUserToTeam(userID, teamID).Return(myerr).Once()
		teamService := bl.NewTeamService(df.teamRepo, df.logger)

		err := teamService.AddUserToTeam(requester.User, userID, teamID)

		sCtx.Assert().Equal(err.ErrNum, bl.Ok)
	})
}

func (s *TeamServiceTestSuite) TestTeamService_AddUserToTeam_AccessDenied(t provider.T) {
	t.Title("Add user to team: Access Denied")
	t.Tags("TeamService")
	t.Parallel()

	t.WithNewStep("Access Denied", func(sCtx provider.StepCtx) {
		df := initExtDependenciesTeam(t)
		requester := builders.NewUserBuilder().WithRole(bl.Reader)
		userID, teamID := 1, 1

		teamService := bl.NewTeamService(df.teamRepo, df.logger)
		err := teamService.AddUserToTeam(requester.User, userID, teamID)

		sCtx.Assert().Equal(err.ErrNum, bl.ErrAccessDenied)
	})
}

func (s *TeamServiceTestSuite) TestTeamService_DeleteUserFromTeam_OK(t provider.T) {
	t.Title("Delete user from team: OK")
	t.Tags("TeamService")
	t.Parallel()

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		df := initExtDependenciesTeam(t)
		requester := builders.NewUserBuilder().WithRole(bl.Admin)
		userID, teamID := 1, 1
		myerr := myerror_builder.NewMyErrorBuilder().WithErrNum(bl.Ok).Build()

		df.teamRepo.EXPECT().DeleteUserFromTeam(userID, teamID).Return(myerr).Once()
		teamService := bl.NewTeamService(df.teamRepo, df.logger)

		err := teamService.DeleteUserFromTeam(requester.User, userID, teamID)

		sCtx.Assert().Equal(err.ErrNum, bl.Ok)
	})
}

func (s *TeamServiceTestSuite) TestTeamService_DeleteUserFromTeam_AccessDenied(t provider.T) {
	t.Title("Delete user from team: Access Denied")
	t.Tags("TeamService")
	t.Parallel()

	t.WithNewStep("Access Denied", func(sCtx provider.StepCtx) {
		df := initExtDependenciesTeam(t)
		requester := builders.NewUserBuilder().WithRole(bl.Reader)
		userID, teamID := 1, 1

		teamService := bl.NewTeamService(df.teamRepo, df.logger)
		err := teamService.DeleteUserFromTeam(requester.User, userID, teamID)

		sCtx.Assert().Equal(err.ErrNum, bl.ErrAccessDenied)
	})
}

func (s *TeamServiceTestSuite) TestTeamService_GetTeamMembers_OK(t provider.T) {
	t.Title("Get team members: OK")
	t.Tags("TeamService")
	t.Parallel()

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		df := initExtDependenciesTeam(t)
		requester := builders.NewUserBuilder().WithRole(bl.Admin)
		teamID := 1
		members := []*models.User{
			builders.NewUserBuilder().Build(),
			builders.NewUserBuilder().Build(),
		}
		myerr := myerror_builder.NewMyErrorBuilder().WithErrNum(bl.Ok).Build()

		df.teamRepo.EXPECT().GetTeamByID(teamID).Return(nil, myerr).Once()
		df.teamRepo.EXPECT().GetTeamMembers(teamID).Return(members, myerr).Once()
		teamService := bl.NewTeamService(df.teamRepo, df.logger)

		users, err := teamService.GetTeamMembers(teamID, requester.User)

		sCtx.Assert().Equal(err.ErrNum, bl.Ok)
		sCtx.Assert().NotNil(users)
		sCtx.Assert().Equal(len(users), len(members))
	})
}

func (s *TeamServiceTestSuite) TestTeamService_GetTeamMembers_AccessDenied(t provider.T) {
	t.Title("Get team members: Access Denied")
	t.Tags("TeamService")
	t.Parallel()

	t.WithNewStep("Access Denied", func(sCtx provider.StepCtx) {
		df := initExtDependenciesTeam(t)
		requester := builders.NewUserBuilder().WithRole(bl.Reader)
		teamID := 1

		teamService := bl.NewTeamService(df.teamRepo, df.logger)
		users, err := teamService.GetTeamMembers(teamID, requester.User)

		sCtx.Assert().Nil(users)
		sCtx.Assert().Equal(err.ErrNum, bl.ErrAccessDenied)
	})
}
