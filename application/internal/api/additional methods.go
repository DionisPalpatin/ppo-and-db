package handlersv2

import (
	"log/slog"
	"net/http"
	"reflect"

	bl "github.com/DionisPalpatin/ppo-and-db/application/internal/business_logic"
	mylogger "github.com/DionisPalpatin/ppo-and-db/application/internal/logger"
	"github.com/DionisPalpatin/ppo-and-db/application/internal/models"
	"github.com/gin-gonic/gin"
)

func getRequester(c *gin.Context, logger *mylogger.MyLogger) *models.User {
	logger.WriteLog("Start get requester", slog.LevelInfo, nil)

	claims, err := bl.ValidateAndParseToken(c.GetHeader("Authorization"))

	if err != nil {
		logger.WriteLog("JWT error: "+err.Error(), slog.LevelError, nil)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization error"})
		return nil
	}

	reqUser := &models.User{Id: claims.UserID}

	if claims.UserRole == "admin" {
		reqUser.Role = bl.Admin
	} else if claims.UserRole == "reader" {
		reqUser.Role = bl.Reader
	} else if claims.UserRole == "author" {
		reqUser.Role = bl.Author
	}

	return reqUser
}

func logOnNull(hs *HandlersStruct, methodName string) {
	val := reflect.ValueOf(hs)
	if val.Kind() == reflect.Ptr && val.IsNil() {
		hs.Configs.LogConfigs.Logger.WriteLog("HandlersStruct is nil when method "+methodName+" is called\n", slog.LevelError, nil)
		return
	}

	if hs.IServices == nil {
		hs.Configs.LogConfigs.Logger.WriteLog("IServices is nil in HandlersStruct during method: "+methodName+"\n", slog.LevelError, nil)
		return
	}

	if hs.IRepos == nil {
		hs.Configs.LogConfigs.Logger.WriteLog("IRepositories is nil in HandlersStruct during method: "+methodName+"\n", slog.LevelError, nil)
		return
	}

	servicesVal := reflect.ValueOf(hs.IServices).Elem()
	for i := 0; i < servicesVal.NumField(); i++ {
		field := servicesVal.Field(i)
		if field.IsNil() {
			fieldName := servicesVal.Type().Field(i).Name
			hs.Configs.LogConfigs.Logger.WriteLog("Service "+fieldName+" is nil in HandlersStruct during method: "+methodName+"\n", slog.LevelError, nil)
		}
	}

	reposVal := reflect.ValueOf(hs.IRepos).Elem()
	for i := 0; i < reposVal.NumField(); i++ {
		field := reposVal.Field(i)
		if field.IsNil() {
			fieldName := reposVal.Type().Field(i).Name
			hs.Configs.LogConfigs.Logger.WriteLog("Repository "+fieldName+" is nil in HandlersStruct during method: "+methodName+"\n", slog.LevelError, nil)
		}
	}
}
