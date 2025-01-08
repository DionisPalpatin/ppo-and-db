package handlersv2

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/DionisPalpatin/ppo-and-db/application/internal/api/converters"
	"github.com/DionisPalpatin/ppo-and-db/application/internal/api/transport_models"
	bl "github.com/DionisPalpatin/ppo-and-db/application/internal/business_logic"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func (hs *HandlersStruct) GetUserHandler(c *gin.Context) {
	reqUser := getRequester(c, hs.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	targetUserID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	user, myErr := hs.IServices.IUsrSvc.GetUser(targetUserID, "", bl.SearchByID, reqUser)
	if myErr == nil {
		hs.Configs.LogConfigs.Logger.WriteLog("myErr is nil", slog.LevelError, nil)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	} else if myErr.ErrNum == bl.ErrAccessDenied {
		hs.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	} else if myErr.ErrNum == bl.NoSuchUser {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	} else if myErr.ErrNum != bl.Ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	userConverted := converters.ToUserFullInfo(user)

	c.Header("Content-Type", "application/my_json")
	c.JSON(http.StatusOK, userConverted)
}

func (hs *HandlersStruct) GetAllUsersHandler(c *gin.Context) {
	reqUser := getRequester(c, hs.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	users, myErr := hs.IServices.IUsrSvc.GetAllUsers(reqUser)
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

	usersConverted := make([]transport_models.UserPublicInfo, 0, len(users))
	for _, user := range users {
		usersConverted = append(usersConverted, converters.ToUserPublicInfo(user))
	}

	c.Header("Content-Type", "application/my_json")
	c.JSON(http.StatusOK, usersConverted)
}

func (hs *HandlersStruct) DeleteUserHandler(c *gin.Context) {
	reqUser := getRequester(c, hs.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	targetUserID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	myErr := hs.IServices.IUsrSvc.DeleteUser(reqUser, targetUserID)
	if myErr.ErrNum == bl.ErrAccessDenied {
		hs.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	} else if myErr.ErrNum == bl.OperationError {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	} else if myErr.ErrNum != bl.Ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Status(http.StatusOK)
}

func (hs *HandlersStruct) UpdateUserHandler(c *gin.Context) {
	reqUser := getRequester(c, hs.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	var userInfo transport_models.UserPrivateInfo
	if err := c.ShouldBindJSON(&userInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	var validate = validator.New()
	if err := validate.Struct(&userInfo); err != nil {
		hs.Configs.LogConfigs.Logger.WriteLog("Invalid request payload", slog.LevelError, nil)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid request payload"})
		return
	}

	targetUserID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	user := converters.FromUserPrivateInfo(&userInfo)
	userInfo.ID = targetUserID

	myErr := hs.IServices.IUsrSvc.UpdateUser(reqUser, &user)

	if myErr.ErrNum == bl.ErrAccessDenied {
		hs.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	} else if myErr.ErrNum == bl.OperationError {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	} else if myErr.ErrNum != bl.Ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Status(http.StatusOK)
}
