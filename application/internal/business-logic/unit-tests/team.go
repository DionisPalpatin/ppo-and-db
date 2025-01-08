package unit_tests

import (
	"testing"

	v1 "github.com/DionisPalpatin/ppo-and-db/application/internal/database/mocks/v1"
	"github.com/stretchr/testify/assert"
)

func TestGetTeam(t *testing.T) {
	t.Run("SuccessGetTeamByID", func(t *testing.T) {
		retOk := bl.CreateError(bl.Ok, "", nil)
		reqUser := &bl.User{Role: bl.Admin}
		returnTeam := &bl.Team{}

		mockTeamRepo := new(v1.MockITeamRepository)
		mockTeamRepo.On("GetTeamByID", 1).Return(returnTeam, retOk)

		tsSrv := bl.TeamService{}
		_, err := tsSrv.GetTeam(1, "", bl.SearchByID, reqUser, mockTeamRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.Ok, err.ErrNum)

		mockTeamRepo.AssertExpectations(t)
	})

	t.Run("SuccessGetTeamByName", func(t *testing.T) {
		retErr := bl.CreateError(bl.Ok, "", nil)
		reqUser := &bl.User{Role: bl.Admin}
		returnTeam := &bl.Team{}

		mockTeamRepo := new(v1.MockITeamRepository)
		mockTeamRepo.On("GetTeamByName", "Team 1").Return(returnTeam, retErr)

		tsSrv := bl.TeamService{}
		_, err := tsSrv.GetTeam(0, "Team 1", bl.SearchByString, reqUser, mockTeamRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.Ok, err.ErrNum)

		mockTeamRepo.AssertExpectations(t)
	})

	t.Run("ErrorRequesterNotAdmin", func(t *testing.T) {
		reqUser := &bl.User{Role: bl.Reader}

		mockRepo := new(v1.MockITeamRepository)

		tsSrv := bl.TeamService{}
		_, err := tsSrv.GetTeam(1, "", bl.SearchByID, reqUser, mockRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrAccessDenied, err.ErrNum)

		mockRepo.AssertExpectations(t)
	})

	t.Run("ErrorUnknownSearchParameter", func(t *testing.T) {
		reqUser := &bl.User{Role: bl.Admin}

		mockRepo := new(v1.MockITeamRepository)

		tsSrv := bl.TeamService{}
		_, err := tsSrv.GetTeam(1, "", 100, reqUser, mockRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrSearchParameter, err.ErrNum)

		mockRepo.AssertExpectations(t)
	})

	t.Run("ErrorGetTeamByID", func(t *testing.T) {
		retErr := bl.CreateError(bl.ErrGetTeamByID, "GetTeamByID", bl.ErrGetTeamByIDError())
		reqUser := &bl.User{Role: bl.Admin}
		returnTeam := &bl.Team{}

		mockTeamRepo := new(v1.MockITeamRepository)
		mockTeamRepo.On("GetTeamByID", 1).Return(returnTeam, retErr)

		tsSrv := bl.TeamService{}
		_, err := tsSrv.GetTeam(1, "", bl.SearchByID, reqUser, mockTeamRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrGetTeamByID, err.ErrNum)

		mockTeamRepo.AssertExpectations(t)
	})

	t.Run("ErrorGetTeamByName", func(t *testing.T) {
		retErr := bl.CreateError(bl.ErrGetTeamByName, "GetTeamByName", bl.ErrGetTeamByNameError())
		reqUser := &bl.User{Role: bl.Admin}
		returnTeam := &bl.Team{}

		mockTeamRepo := new(v1.MockITeamRepository)
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
		retErr := bl.CreateError(bl.Ok, "", nil)
		reqUser := &bl.User{Role: bl.Admin}
		returnTeams := []*bl.Team{&bl.Team{}, &bl.Team{}}

		mockTeamRepo := new(v1.MockITeamRepository)
		mockTeamRepo.On("GetAllTeams").Return(returnTeams, retErr)

		tsSrv := bl.TeamService{}
		_, err := tsSrv.GetAllTeams(reqUser, mockTeamRepo)

		assert.NotNil(t, err)
		assert.Equal(t, err.ErrNum, bl.Ok)

		mockTeamRepo.AssertExpectations(t)
	})

	t.Run("ErrorRequesterNotAdmin", func(t *testing.T) {
		reqUser := &bl.User{Role: bl.Reader}

		mockRepo := new(v1.MockITeamRepository)

		tsSrv := bl.TeamService{}
		_, err := tsSrv.GetAllTeams(reqUser, mockRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrAccessDenied, err.ErrNum)

		mockRepo.AssertExpectations(t)
	})

	t.Run("ErrorGetAllTeams", func(t *testing.T) {
		retErr := bl.CreateError(bl.ErrGetAllTeams, "GetAllTeams", bl.ErrGetAllTeamsError())
		returnTeams := []*bl.Team{&bl.Team{}, &bl.Team{}}
		reqUser := &bl.User{Role: bl.Admin}

		mockTeamRepo := new(v1.MockITeamRepository)
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
		retErr := bl.CreateError(bl.Ok, "", nil)
		reqUser := &bl.User{Role: bl.Admin}
		updateTeam := &bl.Team{Id: 1, Name: "Team 1"}

		mockTeamRepo := new(v1.MockITeamRepository)
		mockTeamRepo.On("UpdateTeam", updateTeam).Return(retErr)

		tsSrv := bl.TeamService{}
		err := tsSrv.UpdateTeam(reqUser, updateTeam, mockTeamRepo)

		assert.NotNil(t, err)
		assert.Equal(t, err.ErrNum, bl.Ok)

		mockTeamRepo.AssertExpectations(t)
	})

	t.Run("ErrorRequesterNotAdmin", func(t *testing.T) {
		reqUser := &bl.User{Role: bl.Reader}

		mockRepo := new(v1.MockITeamRepository)

		tsSrv := bl.TeamService{}
		err := tsSrv.UpdateTeam(reqUser, nil, mockRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrAccessDenied, err.ErrNum)

		mockRepo.AssertExpectations(t)
	})

	t.Run("ErrorUpdateTeam", func(t *testing.T) {
		retErr := bl.CreateError(bl.ErrUpdateTeam, "UpdateTeam", bl.ErrUpdateTeamError())
		reqUser := &bl.User{Role: bl.Admin}
		updateTeam := &bl.Team{Id: 1, Name: "Team 1"}

		mockTeamRepo := new(v1.MockITeamRepository)
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
		retErr := bl.CreateError(bl.Ok, "", nil)
		reqUser := &bl.User{Role: bl.Admin}

		mockTeamRepo := new(v1.MockITeamRepository)
		mockTeamRepo.On("DeleteTeam", 1).Return(retErr)

		tsSrv := bl.TeamService{}
		err := tsSrv.DeleteTeam(reqUser, 1, mockTeamRepo)

		assert.NotNil(t, err)
		assert.Equal(t, err.ErrNum, bl.Ok)

		mockTeamRepo.AssertExpectations(t)
	})

	t.Run("ErrorRequesterNotAdmin", func(t *testing.T) {
		reqUser := &bl.User{Role: bl.Reader}

		mockRepo := new(v1.MockITeamRepository)

		tsSrv := bl.TeamService{}
		err := tsSrv.DeleteTeam(reqUser, 1, mockRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrAccessDenied, err.ErrNum)

		mockRepo.AssertExpectations(t)
	})

	t.Run("ErrorDeleteTeam", func(t *testing.T) {
		retErr := bl.CreateError(bl.ErrDeleteTeam, "DeleteTeam", bl.ErrDeleteTeamError())
		reqUser := &bl.User{Role: bl.Admin}

		mockTeamRepo := new(v1.MockITeamRepository)
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
		retErr := bl.CreateError(bl.Ok, "", nil)
		reqUser := &bl.User{Role: bl.Admin}
		returnTeam := &bl.Team{Id: 1, Name: "Team 1"}

		mockTeamRepo := new(v1.MockITeamRepository)
		mockTeamRepo.On("AddTeam", returnTeam).Return(retErr)

		tsSrv := bl.TeamService{}
		_, err := tsSrv.AddTeam(reqUser, returnTeam, mockTeamRepo)

		assert.NotNil(t, err)
		assert.Equal(t, err.ErrNum, bl.Ok)

		mockTeamRepo.AssertExpectations(t)
	})

	t.Run("ErrorRequesterNotAdmin", func(t *testing.T) {
		reqUser := &bl.User{Role: bl.Reader}

		mockRepo := new(v1.MockITeamRepository)

		tsSrv := bl.TeamService{}
		_, err := tsSrv.AddTeam(reqUser, nil, mockRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrAccessDenied, err.ErrNum)

		mockRepo.AssertExpectations(t)
	})

	t.Run("ErrorAddTeam", func(t *testing.T) {
		retErr := bl.CreateError(bl.ErrAddTeam, "AddTeam", bl.ErrAddTeamError())
		reqUser := &bl.User{Role: bl.Admin}
		returnTeam := &bl.Team{Id: 1, Name: "Team 1"}

		mockTeamRepo := new(v1.MockITeamRepository)
		mockTeamRepo.On("AddTeam", returnTeam).Return(retErr)

		tsSrv := bl.TeamService{}
		_, err := tsSrv.AddTeam(reqUser, returnTeam, mockTeamRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrAddTeam, err.ErrNum)

		mockTeamRepo.AssertExpectations(t)
	})
}

func TestAddUser(t *testing.T) {
	t.Run("SuccessAddUser", func(t *testing.T) {
		retErr := bl.CreateError(bl.Ok, "", nil)
		reqUser := &bl.User{Role: bl.Admin}

		mockTeamRepo := new(v1.MockITeamRepository)
		mockTeamRepo.On("AddUserToTeam", 1, 1).Return(retErr)

		tsSrv := bl.TeamService{}
		err := tsSrv.AddUserToTeam(reqUser, 1, 1, mockTeamRepo)

		assert.NotNil(t, err)
		assert.Equal(t, err.ErrNum, bl.Ok)

		mockTeamRepo.AssertExpectations(t)
	})

	t.Run("ErrorRequesterNotAdmin", func(t *testing.T) {
		reqUser := &bl.User{Role: bl.Reader}

		mockRepo := new(v1.MockITeamRepository)

		tsSrv := bl.TeamService{}
		err := tsSrv.AddUserToTeam(reqUser, 1, 1, mockRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrAccessDenied, err.ErrNum)

		mockRepo.AssertExpectations(t)
	})

	t.Run("ErrorAddUser", func(t *testing.T) {
		retErr := bl.CreateError(bl.ErrAddUser, "AddUserToTeam", bl.ErrAddUserError())
		reqUser := &bl.User{Role: bl.Admin}

		mockTeamRepo := new(v1.MockITeamRepository)
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
		retErr := bl.CreateError(bl.Ok, nil, "")
		returnTeam := &bl.Team{Id: 1, Name: "Team 1"}
		returnMembers := []*bl.User{{Id: 1}, {Id: 2}}
		reqUser := &bl.User{Role: bl.Admin}

		mockTeamRepo := new(v1.MockITeamRepository)
		mockTeamRepo.On("GetTeamByID", 1).Return(returnTeam, retErr)
		mockTeamRepo.On("GetTeamMembers", 1).Return(returnMembers, retErr)

		tsSrv := bl.TeamService{}
		_, err := tsSrv.GetTeamMembers(1, reqUser, mockTeamRepo)

		assert.NotNil(t, err)
		assert.Equal(t, err.ErrNum, bl.Ok)

		mockTeamRepo.AssertExpectations(t)
	})

	t.Run("ErrorRequesterNotAdmin", func(t *testing.T) {
		reqUser := &bl.User{Role: bl.Reader}

		mockRepo := new(v1.MockITeamRepository)

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

		mockTeamRepo := new(v1.MockITeamRepository)
		mockTeamRepo.On("GetTeamByID", 1).Return(team, retErr)

		tsSrv := bl.TeamService{}
		_, err := tsSrv.GetTeamMembers(1, reqUser, mockTeamRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrGetTeamByID, err.ErrNum)

		mockTeamRepo.AssertExpectations(t)
	})

	t.Run("ErrorGetTeamMembers", func(t *testing.T) {
		retOk := bl.CreateError(bl.Ok, nil, "")
		retErr := bl.CreateError(bl.ErrGetTeamMembers, bl.ErrGetTeamMembersError(), "GetTeamMembers")
		returnTeam := &bl.Team{Id: 1, Name: "Team 1"}
		returnTeamMems := []*bl.User{&bl.User{Role: -1}}
		reqUser := &bl.User{Role: bl.Admin}

		mockTeamRepo := new(v1.MockITeamRepository)
		mockTeamRepo.On("GetTeamByID", 1).Return(returnTeam, retOk)
		mockTeamRepo.On("GetTeamMembers", 1).Return(returnTeamMems, retErr)

		tsSrv := bl.TeamService{}
		_, err := tsSrv.GetTeamMembers(1, reqUser, mockTeamRepo)

		assert.NotNil(t, err)
		assert.Equal(t, bl.ErrGetTeamMembers, err.ErrNum)

		mockTeamRepo.AssertExpectations(t)
	})
}
