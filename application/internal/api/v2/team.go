package handlersv2

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/DionisPalpatin/ppo-and-db/application/internal/api/v2/converters"
	"github.com/DionisPalpatin/ppo-and-db/application/internal/api/v2/transport_models"
	bl "github.com/DionisPalpatin/ppo-and-db/application/internal/business_logic"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func (hs *HandlersStruct) GetAllTeamsHandler(c *gin.Context) {
	reqUser := getRequester(c, hs.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	data, myErr := hs.IServices.ITeamSvc.GetAllTeams(reqUser)
	if myErr == nil {
		hs.Configs.LogConfigs.Logger.WriteLog("myErr is nil", slog.LevelError, nil)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	} else if myErr.ErrNum == bl.ErrAccessDenied {
		hs.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	} else if myErr.ErrNum != bl.Ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	dataConverted := make([]transport_models.TeamInfo, 0, len(data))
	for _, team := range data {
		dataConverted = append(dataConverted, converters.ToTeamInfo(team))
	}

	c.Header("Content-Type", "application/my_json")
	c.JSON(http.StatusOK, dataConverted)
}

func (hs *HandlersStruct) GetTeamHandler(c *gin.Context) {
	reqUser := getRequester(c, hs.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	targetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	srcData, myErr := hs.IServices.ITeamSvc.GetTeam(targetID, "", bl.SearchByID, reqUser)
	if myErr == nil {
		hs.Configs.LogConfigs.Logger.WriteLog("myErr is nil", slog.LevelError, nil)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	} else if myErr.ErrNum == bl.ErrAccessDenied {
		hs.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	} else if myErr.ErrNum == bl.NoSuchTeam {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	} else if myErr.ErrNum != bl.Ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	convertedData := converters.ToTeamFullInfo(srcData)

	c.Header("Content-Type", "application/my_json")
	c.JSON(http.StatusOK, convertedData)
}

func (hs *HandlersStruct) AddTeamHandler(c *gin.Context) {
	reqUser := getRequester(c, hs.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	var teamInfo transport_models.TeamInfo
	err := c.ShouldBindJSON(&teamInfo)

	if err != nil {
		hs.Configs.LogConfigs.Logger.WriteLog("Invalid request payload", slog.LevelError, nil)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid request payload"})
		return
	}

	var validate = validator.New()

	if err := validate.Struct(&teamInfo); err != nil {
		hs.Configs.LogConfigs.Logger.WriteLog("Invalid request payload", slog.LevelError, nil)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid request payload"})
		return
	}

	convertedData := converters.FromTeamInfo(&teamInfo)
	var idStruct teamIDStruct
	var myErr *bl.MyError
	idStruct.TeamID, myErr = hs.IServices.ITeamSvc.AddTeam(reqUser, &convertedData)

	if myErr == nil {
		hs.Configs.LogConfigs.Logger.WriteLog("myErr is nil", slog.LevelError, nil)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	} else if myErr.ErrNum == bl.NoSuchTeam {
		hs.Configs.LogConfigs.Logger.WriteLog("Team not found", slog.LevelError, nil)
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	} else if myErr.ErrNum != bl.Ok {
		hs.Configs.LogConfigs.Logger.WriteLog("Error checking user existence", slog.LevelError, nil)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Header("Content-Type", "application/my_json")
	c.JSON(http.StatusOK, idStruct)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
}

func (hs *HandlersStruct) DeleteTeamHandler(c *gin.Context) {
	reqUser := getRequester(c, hs.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	targetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID format"})
		return
	}

	myErr := hs.IServices.ITeamSvc.DeleteTeam(reqUser, targetID)
	if myErr.ErrNum == bl.ErrAccessDenied {
		hs.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	} else if myErr.ErrNum == bl.OperationError {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	} else if myErr.ErrNum != bl.Ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Status(http.StatusOK)
}

func (hs *HandlersStruct) UpdateTeamHandler(c *gin.Context) {
	reqUser := getRequester(c, hs.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	var data transport_models.TeamInfo
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	var validate = validator.New()
	if err := validate.Struct(&data); err != nil {
		hs.Configs.LogConfigs.Logger.WriteLog("Invalid request payload", slog.LevelError, nil)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid request payload"})
		return
	}

	targetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	convData := converters.FromTeamInfo(&data)
	convData.Id = targetID
	myErr := hs.IServices.ITeamSvc.UpdateTeam(reqUser, &convData)

	if myErr.ErrNum == bl.ErrAccessDenied {
		hs.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	} else if myErr.ErrNum == bl.OperationError {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	} else if myErr.ErrNum != bl.Ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Status(http.StatusOK)
}

func (hs *HandlersStruct) AddUserToTeamHandler(c *gin.Context) {
	reqUser := getRequester(c, hs.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	var data userIDStruct
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	var validate = validator.New()
	if err := validate.Struct(&data); err != nil {
		hs.Configs.LogConfigs.Logger.WriteLog("Invalid request payload", slog.LevelError, nil)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid request payload"})
		return
	}

	targetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID format"})
		return
	}

	myErr := hs.IServices.ITeamSvc.AddUserToTeam(reqUser, data.UserID, targetID)

	if myErr.ErrNum == bl.ErrAccessDenied {
		hs.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	} else if myErr.ErrNum == bl.OperationError {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	} else if myErr.ErrNum != bl.Ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Status(http.StatusOK)
}

func (hs *HandlersStruct) DeleteUserFromTeamHandler(c *gin.Context) {
	reqUser := getRequester(c, hs.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	teamID, err := strconv.Atoi(c.Param("teamID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID format"})
		return
	}
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	myErr := hs.IServices.ITeamSvc.DeleteUserFromTeam(reqUser, userID, teamID)

	if myErr.ErrNum == bl.ErrAccessDenied {
		hs.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	} else if myErr.ErrNum == bl.OperationError {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	} else if myErr.ErrNum != bl.Ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Status(http.StatusOK)
}

func (hs *HandlersStruct) GetTeamMembersHandler(c *gin.Context) {
	reqUser := getRequester(c, hs.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	teamID, err := strconv.Atoi(c.Param("teamID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID format"})
		return
	}

	data, myErr := hs.IServices.ITeamSvc.GetTeamMembers(teamID, reqUser)
	if myErr.ErrNum == bl.ErrAccessDenied {
		hs.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	} else if myErr.ErrNum == bl.OperationError {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	} else if myErr.ErrNum != bl.Ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	convertedData := make([]transport_models.UserPublicInfo, 0, len(data))
	for _, user := range data {
		convertedData = append(convertedData, converters.ToUserPublicInfo(user))
	}

	c.Header("Content-Type", "application/my_json")
	c.JSON(http.StatusOK, convertedData)
}
