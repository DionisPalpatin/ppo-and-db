package handlersv2

import (
	"encoding/json"
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/api/v2/converters"
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/api/v2/transport_models"
	bl "github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/business_logic"
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

func (hs *HandlersStruct) RegisterHandler(c *gin.Context) {
	var userRegInfo transport_models.UserRegistrationInfo

	err := c.ShouldBindJSON(&userRegInfo)

	if err != nil {
		hs.Configs.LogConfigs.Logger.WriteLog("Invalid request payload", slog.LevelError, nil)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid request payload"})
		return
	}

	var validate = validator.New()

	if err := validate.Struct(&userRegInfo); err != nil {
		hs.Configs.LogConfigs.Logger.WriteLog("Invalid request payload", slog.LevelError, nil)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid request payload"})
		return
	}

	existingUser, myErr := hs.IServices.IUsrSvc.GetUser(-1, userRegInfo.Login, bl.SearchByString, &models.User{Role: bl.Admin})

	if myErr == nil {
		hs.Configs.LogConfigs.Logger.WriteLog("myErr is nil", slog.LevelError, nil)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	} else if myErr.ErrNum == bl.Ok && existingUser != nil {
		hs.Configs.LogConfigs.Logger.WriteLog(myErr.Error(), slog.LevelError, nil)
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	user := converters.FromRegistrationInfo(&userRegInfo)
	_, myErr = hs.IServices.IOAuthSvc.RegisterUser(user.Fio, user.Login, user.Password)

	if myErr.ErrNum != bl.Ok {
		hs.Configs.LogConfigs.Logger.WriteLog("Registration error", slog.LevelError, nil)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Status(http.StatusOK)
}

func (hs *HandlersStruct) LoginHandler(c *gin.Context) {
	logOnNull(hs, "LoginHandler")

	hs.Configs.LogConfigs.Logger.WriteLog("Start login", slog.LevelInfo, nil)

	var userLoginInfo transport_models.UserLoginInfo

	err := c.ShouldBindJSON(&userLoginInfo)

	if err != nil {
		hs.Configs.LogConfigs.Logger.WriteLog("Invalid request payload", slog.LevelError, nil)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid request payload"})
		return
	}

	var validate = validator.New()

	if err := validate.Struct(&userLoginInfo); err != nil {
		hs.Configs.LogConfigs.Logger.WriteLog("Invalid request payload", slog.LevelError, nil)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid request payload"})
		return
	}

	hs.Configs.LogConfigs.Logger.WriteLog("JSON in validated. Start to convert structs", slog.LevelInfo, nil)

	user := converters.FromLoginInfo(&userLoginInfo)

	hs.Configs.LogConfigs.Logger.WriteLog("Start get user to check existanse", slog.LevelInfo, nil)

	_, myErr := hs.IServices.IUsrSvc.GetUser(-1, user.Login, bl.SearchByString, &models.User{Role: bl.Admin})

	hs.Configs.LogConfigs.Logger.WriteLog("End get user to check existanse", slog.LevelInfo, nil)

	if myErr == nil {
		hs.Configs.LogConfigs.Logger.WriteLog("myErr is nil", slog.LevelError, nil)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	} else if myErr.ErrNum == bl.NoSuchUser {
		hs.Configs.LogConfigs.Logger.WriteLog("User not found", slog.LevelError, nil)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	} else if myErr.ErrNum != bl.Ok {
		hs.Configs.LogConfigs.Logger.WriteLog("Error checking user existence", slog.LevelError, nil)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	token, myErr := hs.IServices.IOAuthSvc.SignInUser(user.Login, user.Password)

	if myErr == nil {
		hs.Configs.LogConfigs.Logger.WriteLog("myErr is nil", slog.LevelError, nil)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	} else if myErr.ErrNum == bl.OperationError {
		hs.Configs.LogConfigs.Logger.WriteLog("Invalid login or password", slog.LevelError, nil)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid login or password"})
		return
	} else if myErr.ErrNum != bl.Ok {
		hs.Configs.LogConfigs.Logger.WriteLog("Error during login", slog.LevelError, nil)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	hs.Configs.LogConfigs.Logger.WriteLog("Start convert token", slog.LevelInfo, nil)

	convToken := converters.ToTransportToken(token)

	hs.Configs.LogConfigs.Logger.WriteLog("Start send token", slog.LevelInfo, nil)

	w.Header().Set("Content-Type", "application/my_json")
	err = json.NewEncoder(w).Encode(convToken)

	if err != nil {
		hs.Configs.LogConfigs.Logger.WriteLog("Error checking user existence", slog.LevelError, nil)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
}
