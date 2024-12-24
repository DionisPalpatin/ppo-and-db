package handlersv2

import (
	bl "github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/business_logic"
	mylogger "github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/logger"
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/models"
	"log/slog"
	"net/http"
)

func getRequester(r *http.Request, w http.ResponseWriter, logger *mylogger.MyLogger) *models.User {
	logger.WriteLog("Start get requester", slog.LevelInfo, nil)
	
	claims, err := bl.ValidateAndParseToken(r.Header.Get("Authorization"))

	if err != nil {
		logger.WriteLog("JWT error: "+err.Error(), slog.LevelError, nil)
		http.Error(w, "Authorization error", http.StatusUnauthorized)
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
