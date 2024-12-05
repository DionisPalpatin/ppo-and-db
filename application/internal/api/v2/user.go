package handlersv2

import (
	"encoding/json"
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/api/v2/converters"
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/api/v2/transport_models"
	bl "github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/business_logic"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"strconv"
)

func (app *App) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	reqUser := getRequester(r, w, app.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	vars := mux.Vars(r)
	targetUserID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	user, myErr := app.IServices.IUsrSvc.GetUser(targetUserID, "", bl.SearchByID, reqUser)
	if myErr == nil {
		app.Configs.LogConfigs.Logger.WriteLog("myErr is nil", slog.LevelError, nil)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	} else if myErr.ErrNum == bl.ErrAccessDenied {
		app.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	} else if myErr.ErrNum == bl.NoSuchUser {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if myErr.ErrNum != bl.Ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	userConverted := converters.ToUserFullInfo(user)

	w.Header().Set("Content-Type", "application/my_json")
	err = json.NewEncoder(w).Encode(userConverted)

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (app *App) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	reqUser := getRequester(r, w, app.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	users, myErr := app.IServices.IUsrSvc.GetAllUsers(reqUser)
	if myErr == nil {
		app.Configs.LogConfigs.Logger.WriteLog("myErr is nil", slog.LevelError, nil)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	} else if myErr.ErrNum == bl.ErrAccessDenied {
		app.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	} else if myErr.ErrNum != bl.Ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	usersConverted := make([]transport_models.UserPublicInfo, 0, len(users))
	for _, user := range users {
		usersConverted = append(usersConverted, converters.ToUserPublicInfo(user))
	}

	w.Header().Set("Content-Type", "application/my_json")
	err := json.NewEncoder(w).Encode(usersConverted)

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (app *App) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	reqUser := getRequester(r, w, app.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	vars := mux.Vars(r)
	targetUserID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	myErr := app.IServices.IUsrSvc.DeleteUser(reqUser, targetUserID)
	if myErr.ErrNum == bl.ErrAccessDenied {
		app.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	} else if myErr.ErrNum == bl.OperationError {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if myErr.ErrNum != bl.Ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *App) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	reqUser := getRequester(r, w, app.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	var userInfo transport_models.UserPrivateInfo
	if err := json.NewDecoder(r.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var validate = validator.New()
	if err := validate.Struct(&userInfo); err != nil {
		app.Configs.LogConfigs.Logger.WriteLog("Invalid request payload", slog.LevelError, nil)
		http.Error(w, "Invalid request payload", http.StatusUnprocessableEntity)
		return
	}

	vars := mux.Vars(r)
	targetUserID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	user := converters.FromUserPrivateInfo(userInfo)
	userInfo.ID = targetUserID

	myErr := app.IServices.IUsrSvc.UpdateUser(reqUser, &user)

	if myErr.ErrNum == bl.ErrAccessDenied {
		app.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	} else if myErr.ErrNum == bl.OperationError {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if myErr.ErrNum != bl.Ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
