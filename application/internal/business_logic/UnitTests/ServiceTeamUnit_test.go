package UnitTests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/business_logic"
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/business_logic/UnitTests/mocks"
)

func TestGetTeam(t *testing.T) {
	t.Run("SuccessGetTeamByID", func(t *testing.T) {
		retOk := bl.CreateError(bl.AllIsOk, nil, "")
		reqUser := &bl.User{Role: bl.Admin}
		returnTeam := &bl.Team{}

		mockTeamRepo := new(mocks.MockITeamRepository)
		mockTeamRepo.On("GetTeamByID", 1).Return(returnTeam, retOk)

		tsSrv := bl.TeamService{}
		_, err := tsSrv.GetTeam(1, "", bl.SearchByID, reqUser, mockTeamRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.AllIsOk, err.ErrNum)

		mockTeamRepo.AssertExpectations(t)
	})

	t.Run("SuccessGetTeamByName", func(t *testing.T) {
		retErr := bl.CreateError(bl.AllIsOk, nil, "")
		reqUser := &bl.User{Role: bl.Admin}
		returnTeam := &bl.Team{}

		mockTeamRepo := new(mocks.MockITeamRepository)
		mockTeamRepo.On("GetTeamByName", "Team 1").Return(returnTeam, retErr)

		tsSrv := bl.TeamService{}
		_, err := tsSrv.GetTeam(0, "Team 1", bl.SearchByString, reqUser, mockTeamRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.AllIsOk, err.ErrNum)

		mockTeamRepo.AssertExpectations(t)
	})

	t.Run("ErrorRequesterNotAdmin", func(t *testing.T) {
		reqUser := &bl.User{Role: bl.Reader}

		mockRepo := new(mocks.MockITeamRepository)

		tsSrv := bl.TeamService{}
		_, err := tsSrv.GetTeam(1, "", bl.SearchByID, reqUser, mockRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrAccessDenied, err.ErrNum)

		mockRepo.AssertExpectations(t)
	})

	t.Run("ErrorUnknownSearchParameter", func(t *testing.T) {
		reqUser := &bl.User{Role: bl.Admin}

		mockRepo := new(mocks.MockITeamRepository)

		tsSrv := bl.TeamService{}
		_, err := tsSrv.GetTeam(1, "", 100, reqUser, mockRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrSearchParameter, err.ErrNum)

		mockRepo.AssertExpectations(t)
	})

	t.Run("ErrorGetTeamByID", func(t *testing.T) {
		retErr := bl.CreateError(bl.ErrGetTeamByID, bl.ErrGetTeamByIDError(), "GetTeamByID")
		reqUser := &bl.User{Role: bl.Admin}
		returnTeam := &bl.Team{}

		mockTeamRepo := new(mocks.MockITeamRepository)
		mockTeamRepo.On("GetTeamByID", 1).Return(returnTeam, retErr)

		tsSrv := bl.TeamService{}
		_, err := tsSrv.GetTeam(1, "", bl.SearchByID, reqUser, mockTeamRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrGetTeamByID, err.ErrNum)

		mockTeamRepo.AssertExpectations(t)
	})

	t.Run("ErrorGetTeamByName", func(t *testing.T) {
		retErr := bl.CreateError(bl.ErrGetTeamByName, bl.ErrGetTeamByNameError(), "GetTeamByName")
		reqUser := &bl.User{Role: bl.Admin}
		returnTeam := &bl.Team{}

		mockTeamRepo := new(mocks.MockITeamRepository)
		mockTeamRepo.On("GetTeamByName", "Team 1").Return(returnTeam, retErr)

		tsSrv := bl.TeamService{}
		_, err := tsSrv.GetTeam(0, "Team 1", bl.SearchByString, reqUser, mockTeamRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrGetTeamByName, err.ErrNum)

		mockTeamRepo.AssertExpectations(t)
	})
}

func TestGetAllTeams(t *testing.T) {
	t.Run("SuccessGetAllTeams", func(t *testing.T) {
		retErr := bl.CreateError(bl.AllIsOk, nil, "")
		reqUser := &bl.User{Role: bl.Admin}
		returnTeams := []*bl.Team{&bl.Team{}, &bl.Team{}}

		mockTeamRepo := new(mocks.MockITeamRepository)
		mockTeamRepo.On("GetAllTeams").Return(returnTeams, retErr)

		tsSrv := bl.TeamService{}
		_, err := tsSrv.GetAllTeams(reqUser, mockTeamRepo)

		assert.NotNil(t, err)
		assert.Equal(t, err.ErrNum, bl.AllIsOk)

		mockTeamRepo.AssertExpectations(t)
	})

	t.Run("ErrorRequesterNotAdmin", func(t *testing.T) {
		reqUser := &bl.User{Role: bl.Reader}

		mockRepo := new(mocks.MockITeamRepository)

		tsSrv := bl.TeamService{}
		_, err := tsSrv.GetAllTeams(reqUser, mockRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrAccessDenied, err.ErrNum)

		mockRepo.AssertExpectations(t)
	})

	t.Run("ErrorGetAllTeams", func(t *testing.T) {
		retErr := bl.CreateError(bl.ErrGetAllTeams, bl.ErrGetAllTeamsError(), "GetAllTeams")
		returnTeams := []*bl.Team{&bl.Team{}, &bl.Team{}}
		reqUser := &bl.User{Role: bl.Admin}

		mockTeamRepo := new(mocks.MockITeamRepository)
		mockTeamRepo.On("GetAllTeams").Return(returnTeams, retErr)

		tsSrv := bl.TeamService{}
		_, err := tsSrv.GetAllTeams(reqUser, mockTeamRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrGetAllTeams, err.ErrNum)

		mockTeamRepo.AssertExpectations(t)
	})
}

func TestUpdateTeam(t *testing.T) {
	t.Run("SuccessUpdateTeam", func(t *testing.T) {
		retErr := bl.CreateError(bl.AllIsOk, nil, "")
		reqUser := &bl.User{Role: bl.Admin}
		updateTeam := &bl.Team{Id: 1, Name: "Team 1"}

		mockTeamRepo := new(mocks.MockITeamRepository)
		mockTeamRepo.On("UpdateTeam", updateTeam).Return(retErr)

		tsSrv := bl.TeamService{}
		err := tsSrv.UpdateTeam(reqUser, updateTeam, mockTeamRepo)

		assert.NotNil(t, err)
		assert.Equal(t, err.ErrNum, bl.AllIsOk)

		mockTeamRepo.AssertExpectations(t)
	})

	t.Run("ErrorRequesterNotAdmin", func(t *testing.T) {
		reqUser := &bl.User{Role: bl.Reader}

		mockRepo := new(mocks.MockITeamRepository)

		tsSrv := bl.TeamService{}
		err := tsSrv.UpdateTeam(reqUser, nil, mockRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrAccessDenied, err.ErrNum)

		mockRepo.AssertExpectations(t)
	})

	t.Run("ErrorUpdateTeam", func(t *testing.T) {
		retErr := bl.CreateError(bl.ErrUpdateTeam, bl.ErrUpdateTeamError(), "UpdateTeam")
		reqUser := &bl.User{Role: bl.Admin}
		updateTeam := &bl.Team{Id: 1, Name: "Team 1"}

		mockTeamRepo := new(mocks.MockITeamRepository)
		mockTeamRepo.On("UpdateTeam", updateTeam).Return(retErr)

		tsSrv := bl.TeamService{}
		err := tsSrv.UpdateTeam(reqUser, updateTeam, mockTeamRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrUpdateTeam, err.ErrNum)

		mockTeamRepo.AssertExpectations(t)
	})
}

func TestDeleteTeam(t *testing.T) {
	t.Run("SuccessDeleteTeam", func(t *testing.T) {
		retErr := bl.CreateError(bl.AllIsOk, nil, "")
		reqUser := &bl.User{Role: bl.Admin}

		mockTeamRepo := new(mocks.MockITeamRepository)
		mockTeamRepo.On("DeleteTeam", 1).Return(retErr)

		tsSrv := bl.TeamService{}
		err := tsSrv.DeleteTeam(reqUser, 1, mockTeamRepo)

		assert.NotNil(t, err)
		assert.Equal(t, err.ErrNum, bl.AllIsOk)

		mockTeamRepo.AssertExpectations(t)
	})

	t.Run("ErrorRequesterNotAdmin", func(t *testing.T) {
		reqUser := &bl.User{Role: bl.Reader}

		mockRepo := new(mocks.MockITeamRepository)

		tsSrv := bl.TeamService{}
		err := tsSrv.DeleteTeam(reqUser, 1, mockRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrAccessDenied, err.ErrNum)

		mockRepo.AssertExpectations(t)
	})

	t.Run("ErrorDeleteTeam", func(t *testing.T) {
		retErr := bl.CreateError(bl.ErrDeleteTeam, bl.ErrDeleteTeamError(), "DeleteTeam")
		reqUser := &bl.User{Role: bl.Admin}

		mockTeamRepo := new(mocks.MockITeamRepository)
		mockTeamRepo.On("DeleteTeam", 1).Return(retErr)

		tsSrv := bl.TeamService{}
		err := tsSrv.DeleteTeam(reqUser, 1, mockTeamRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrDeleteTeam, err.ErrNum)

		mockTeamRepo.AssertExpectations(t)
	})
}

func TestAddTeam(t *testing.T) {
	t.Run("SuccessAddTeam", func(t *testing.T) {
		retErr := bl.CreateError(bl.AllIsOk, nil, "")
		reqUser := &bl.User{Role: bl.Admin}
		returnTeam := &bl.Team{Id: 1, Name: "Team 1"}

		mockTeamRepo := new(mocks.MockITeamRepository)
		mockTeamRepo.On("AddTeam", returnTeam).Return(retErr)

		tsSrv := bl.TeamService{}
		err := tsSrv.AddTeam(reqUser, returnTeam, mockTeamRepo)

		assert.NotNil(t, err)
		assert.Equal(t, err.ErrNum, bl.AllIsOk)

		mockTeamRepo.AssertExpectations(t)
	})

	t.Run("ErrorRequesterNotAdmin", func(t *testing.T) {
		reqUser := &bl.User{Role: bl.Reader}

		mockRepo := new(mocks.MockITeamRepository)

		tsSrv := bl.TeamService{}
		err := tsSrv.AddTeam(reqUser, nil, mockRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrAccessDenied, err.ErrNum)

		mockRepo.AssertExpectations(t)
	})

	t.Run("ErrorAddTeam", func(t *testing.T) {
		retErr := bl.CreateError(bl.ErrAddTeam, bl.ErrAddTeamError(), "AddTeam")
		reqUser := &bl.User{Role: bl.Admin}
		returnTeam := &bl.Team{Id: 1, Name: "Team 1"}

		mockTeamRepo := new(mocks.MockITeamRepository)
		mockTeamRepo.On("AddTeam", returnTeam).Return(retErr)

		tsSrv := bl.TeamService{}
		err := tsSrv.AddTeam(reqUser, returnTeam, mockTeamRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrAddTeam, err.ErrNum)

		mockTeamRepo.AssertExpectations(t)
	})
}

func TestAddUser(t *testing.T) {
	t.Run("SuccessAddUser", func(t *testing.T) {
		retErr := bl.CreateError(bl.AllIsOk, nil, "")
		reqUser := &bl.User{Role: bl.Admin}

		mockTeamRepo := new(mocks.MockITeamRepository)
		mockTeamRepo.On("AddUserToTeam", 1, 1).Return(retErr)

		tsSrv := bl.TeamService{}
		err := tsSrv.AddUserToTeam(reqUser, 1, 1, mockTeamRepo)

		assert.NotNil(t, err)
		assert.Equal(t, err.ErrNum, bl.AllIsOk)

		mockTeamRepo.AssertExpectations(t)
	})

	t.Run("ErrorRequesterNotAdmin", func(t *testing.T) {
		reqUser := &bl.User{Role: bl.Reader}

		mockRepo := new(mocks.MockITeamRepository)

		tsSrv := bl.TeamService{}
		err := tsSrv.AddUserToTeam(reqUser, 1, 1, mockRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrAccessDenied, err.ErrNum)

		mockRepo.AssertExpectations(t)
	})

	t.Run("ErrorAddUser", func(t *testing.T) {
		retErr := bl.CreateError(bl.ErrAddUser, bl.ErrAddUserError(), "AddUserToTeam")
		reqUser := &bl.User{Role: bl.Admin}

		mockTeamRepo := new(mocks.MockITeamRepository)
		mockTeamRepo.On("AddUserToTeam", 1, 1).Return(retErr)

		tsSrv := bl.TeamService{}
		err := tsSrv.AddUserToTeam(reqUser, 1, 1, mockTeamRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrAddUser, err.ErrNum)

		mockTeamRepo.AssertExpectations(t)
	})
}

func TestGetTeamMembers(t *testing.T) {
	t.Run("SuccessGetTeamMembers", func(t *testing.T) {
		retErr := bl.CreateError(bl.AllIsOk, nil, "")
		returnTeam := &bl.Team{Id: 1, Name: "Team 1"}
		returnMembers := []*bl.User{{Id: 1}, {Id: 2}}
		reqUser := &bl.User{Role: bl.Admin}

		mockTeamRepo := new(mocks.MockITeamRepository)
		mockTeamRepo.On("GetTeamByID", 1).Return(returnTeam, retErr)
		mockTeamRepo.On("GetTeamMembers", 1).Return(returnMembers, retErr)

		tsSrv := bl.TeamService{}
		_, err := tsSrv.GetTeamMembers(1, reqUser, mockTeamRepo)

		assert.NotNil(t, err)
		assert.Equal(t, err.ErrNum, bl.AllIsOk)

		mockTeamRepo.AssertExpectations(t)
	})

	t.Run("ErrorRequesterNotAdmin", func(t *testing.T) {
		reqUser := &bl.User{Role: bl.Reader}

		mockRepo := new(mocks.MockITeamRepository)

		tsSrv := bl.TeamService{}
		_, err := tsSrv.GetTeamMembers(1, reqUser, mockRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrAccessDenied, err.ErrNum)

		mockRepo.AssertExpectations(t)
	})

	t.Run("ErrorGetTeamByID", func(t *testing.T) {
		retErr := bl.CreateError(bl.ErrGetTeamByID, bl.ErrGetTeamByIDError(), "GetTeamByID")
		team := &bl.Team{}
		reqUser := &bl.User{Role: bl.Admin}

		mockTeamRepo := new(mocks.MockITeamRepository)
		mockTeamRepo.On("GetTeamByID", 1).Return(team, retErr)

		tsSrv := bl.TeamService{}
		_, err := tsSrv.GetTeamMembers(1, reqUser, mockTeamRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrGetTeamByID, err.ErrNum)

		mockTeamRepo.AssertExpectations(t)
	})

	t.Run("ErrorGetTeamMembers", func(t *testing.T) {
		retOk := bl.CreateError(bl.AllIsOk, nil, "")
		retErr := bl.CreateError(bl.ErrGetTeamMembers, bl.ErrGetTeamMembersError(), "GetTeamMembers")
		returnTeam := &bl.Team{Id: 1, Name: "Team 1"}
		returnTeamMems := []*bl.User{&bl.User{Role: -1}}
		reqUser := &bl.User{Role: bl.Admin}

		mockTeamRepo := new(mocks.MockITeamRepository)
		mockTeamRepo.On("GetTeamByID", 1).Return(returnTeam, retOk)
		mockTeamRepo.On("GetTeamMembers", 1).Return(returnTeamMems, retErr)

		tsSrv := bl.TeamService{}
		_, err := tsSrv.GetTeamMembers(1, reqUser, mockTeamRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrGetTeamMembers, err.ErrNum)

		mockTeamRepo.AssertExpectations(t)
	})
}
