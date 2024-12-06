package handlersv2

import (
	"encoding/json"
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/api/v2/converters"
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/api/v2/transport_models"
	bl "github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/business_logic"
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/models"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

func (app *App) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var userRegInfo transport_models.UserRegistrationInfo

	err := json.NewDecoder(r.Body).Decode(&userRegInfo)

	if err != nil {
		app.Configs.LogConfigs.Logger.WriteLog("Invalid request payload", slog.LevelError, nil)
		http.Error(w, "Invalid request payload", http.StatusUnprocessableEntity)
		return
	}

	var validate = validator.New()

	if err := validate.Struct(&userRegInfo); err != nil {
		app.Configs.LogConfigs.Logger.WriteLog("Invalid request payload", slog.LevelError, nil)
		http.Error(w, "Invalid request payload", http.StatusUnprocessableEntity)
		return
	}

	existingUser, myErr := app.IServices.IUsrSvc.GetUser(-1, userRegInfo.Login, bl.SearchByString, &models.User{Role: bl.Admin})

	if myErr == nil {
		app.Configs.LogConfigs.Logger.WriteLog("myErr is nil", slog.LevelError, nil)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	} else if myErr.ErrNum == bl.Ok && existingUser != nil {
		app.Configs.LogConfigs.Logger.WriteLog(myErr.Error(), slog.LevelError, nil)
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	user := converters.FromRegistrationInfo(&userRegInfo)
	_, myErr = app.IServices.IOAuthSvc.RegisterUser(user.Fio, user.Login, user.Password)

	if myErr.ErrNum != bl.Ok {
		app.Configs.LogConfigs.Logger.WriteLog("Registration error", slog.LevelError, nil)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *App) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var userLoginInfo transport_models.UserLoginInfo

	err := json.NewDecoder(r.Body).Decode(&userLoginInfo)

	if err != nil {
		app.Configs.LogConfigs.Logger.WriteLog("Invalid request payload", slog.LevelError, nil)
		http.Error(w, "Invalid request payload", http.StatusUnprocessableEntity)
		return
	}

	var validate = validator.New()

	if err := validate.Struct(&userLoginInfo); err != nil {
		app.Configs.LogConfigs.Logger.WriteLog("Invalid request payload", slog.LevelError, nil)
		http.Error(w, "Invalid request payload", http.StatusUnprocessableEntity)
		return
	}

	user := converters.FromLoginInfo(&userLoginInfo)
	_, myErr := app.IServices.IUsrSvc.GetUser(-1, user.Login, bl.SearchByString, &models.User{Role: bl.Admin})

	if myErr == nil {
		app.Configs.LogConfigs.Logger.WriteLog("myErr is nil", slog.LevelError, nil)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	} else if myErr.ErrNum == bl.NoSuchUser {
		app.Configs.LogConfigs.Logger.WriteLog("User not found", slog.LevelError, nil)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if myErr.ErrNum != bl.Ok {
		app.Configs.LogConfigs.Logger.WriteLog("Error checking user existence", slog.LevelError, nil)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	token, myErr := app.IServices.IOAuthSvc.SignInUser(user.Login, user.Password)

	if myErr == nil {
		app.Configs.LogConfigs.Logger.WriteLog("myErr is nil", slog.LevelError, nil)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	} else if myErr.ErrNum == bl.OperationError {
		app.Configs.LogConfigs.Logger.WriteLog("Invalid login or password", slog.LevelError, nil)
		http.Error(w, "Invalid login or password", http.StatusUnauthorized)
		return
	} else if myErr.ErrNum != bl.Ok {
		app.Configs.LogConfigs.Logger.WriteLog("Error during login", slog.LevelError, nil)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/my_json")
	err = json.NewEncoder(w).Encode(models.Token{AccessToken: token})

	if err != nil {
		app.Configs.LogConfigs.Logger.WriteLog("Error checking user existence", slog.LevelError, nil)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
