package handlersv2

import (
	"encoding/json"
	bl "github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/business_logic"
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/models"
	"log/slog"
	"net/http"
)

type UserRegScheme struct {
	Fio      string `json:"fio"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserLoginScheme struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (app *App) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var userRegInfo UserRegScheme
	err := json.NewDecoder(r.Body).Decode(&userRegInfo)

	if err != nil {
		app.Configs.LogConfigs.Logger.WriteLog("Invalid request payload", slog.LevelError, nil)
		http.Error(w, "Invalid request payload", http.StatusUnprocessableEntity)
		return
	}

	existingUser, myErr := app.IServices.IUsrSvc.GetUser(-1, userRegInfo.Login, bl.SearchByString, &models.User{Role: bl.Admin})
	if myErr.ErrNum == bl.Ok && existingUser != nil {
		app.Configs.LogConfigs.Logger.WriteLog(myErr.Error(), slog.LevelError, nil)
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	_, myErr = app.IServices.IOAuthSvc.RegisterUser(userRegInfo.Fio, userRegInfo.Login, userRegInfo.Password)
	if myErr.ErrNum != bl.Ok {
		app.Configs.LogConfigs.Logger.WriteLog("Registration error", slog.LevelError, nil)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *App) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var userLoginInfo UserLoginScheme
	err := json.NewDecoder(r.Body).Decode(&userLoginInfo)

	if err != nil {
		app.Configs.LogConfigs.Logger.WriteLog("Invalid request payload", slog.LevelError, nil)
		http.Error(w, "Invalid request payload", http.StatusUnprocessableEntity)
		return
	}

	// Проверяем существование пользователя
	_, myErr := app.IServices.IUsrSvc.GetUser(-1, userLoginInfo.Login, bl.SearchByString, &models.User{Role: bl.Admin})
	if myErr.ErrNum == bl.NoSuchUser {
		app.Configs.LogConfigs.Logger.WriteLog("User not found", slog.LevelError, nil)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if myErr.ErrNum != bl.Ok {
		app.Configs.LogConfigs.Logger.WriteLog("Error checking user existence", slog.LevelError, nil)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Аутентификация пользователя
	token, myErr := app.IServices.IOAuthSvc.SignInUser(userLoginInfo.Login, userLoginInfo.Password)
	if myErr.ErrNum == bl.ErrSignInUser {
		app.Configs.LogConfigs.Logger.WriteLog("Invalid login or password", slog.LevelError, nil)
		http.Error(w, "Invalid login or password", http.StatusUnauthorized)
		return
	} else if myErr.ErrNum != bl.Ok {
		app.Configs.LogConfigs.Logger.WriteLog("Error during login", slog.LevelError, nil)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Успешная аутентификация
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.Token{AccessToken: token})
}
