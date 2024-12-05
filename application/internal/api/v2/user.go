package handlersv2

import (
	"encoding/json"
	"fmt"
	bl "github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/business_logic"
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/models"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"strconv"
)

func (app *App) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	claims, err := bl.ValidateAndParseToken(r.Header.Get("Authorization"))
	if err != nil {
		app.Configs.LogConfigs.Logger.WriteLog("JWT error: "+err.Error(), slog.LevelError, nil)
		http.Error(w, "Authorization error", http.StatusUnauthorized)
		return
	}

	app.Configs.LogConfigs.Logger.WriteLog(fmt.Sprintf("User %d with role %s accessing GetUserHandler", claims.UserID, claims.UserRole), slog.LevelInfo, nil)

	vars := mux.Vars(r)
	targetUserID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	if claims.UserRole != "admin" && claims.UserID != targetUserID {
		app.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	}

	reqUser := &models.User{Id: claims.UserID}
	if claims.UserRole == "admin" {
		reqUser.Role = bl.Admin
	} else if claims.UserRole == "reader" {
		reqUser.Role = bl.Reader
	} else if claims.UserRole == "author" {
		reqUser.Role = bl.Author
	}

	user, myErr := app.IServices.IUsrSvc.GetUser(targetUserID, "", bl.SearchByID, reqUser)
	if myErr.ErrNum == bl.ErrGetUserByID {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if myErr.ErrNum != bl.Ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (app *App) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	// Извлечение роли и ID из токена
	claims, err := bl.ValidateAndParseToken(r.Header.Get("Authorization"))
	if err != nil {
		app.Configs.LogConfigs.Logger.WriteLog("JWT error: "+err.Error(), slog.LevelError, nil)
		http.Error(w, "Authorization error", http.StatusUnauthorized)
		return
	}

	// Логирование доступа
	app.Configs.LogConfigs.Logger.WriteLog(fmt.Sprintf("User %d with role %s accessing GetAllUsersHandler", claims.UserID, claims.UserRole), slog.LevelInfo, nil)

	reqUser := &models.User{Id: claims.UserID}
	if claims.UserRole == "admin" {
		reqUser.Role = bl.Admin
	} else if claims.UserRole == "reader" {
		reqUser.Role = bl.Reader
	} else if claims.UserRole == "author" {
		reqUser.Role = bl.Author
	}

	// Получение всех пользователей
	users, myErr := app.IServices.IUsrSvc.GetAllUsers(reqUser)
	if myErr.ErrNum == bl.ErrAccessDenied {
		app.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	} else if myErr.ErrNum != bl.Ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Ответ с пользователями
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (app *App) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	claims, err := bl.ValidateAndParseToken(r.Header.Get("Authorization"))
	if err != nil {
		app.Configs.LogConfigs.Logger.WriteLog("JWT error: "+err.Error(), slog.LevelError, nil)
		http.Error(w, "Authorization error", http.StatusUnauthorized)
		return
	}

	app.Configs.LogConfigs.Logger.WriteLog(fmt.Sprintf("User %d with role %s accessing DeleteUserHandler", claims.UserID, claims.UserRole), slog.LevelInfo, nil)

	vars := mux.Vars(r)
	targetUserID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	reqUser := &models.User{Id: claims.UserID}
	if claims.UserRole == "admin" {
		reqUser.Role = bl.Admin
	} else if claims.UserRole == "reader" {
		reqUser.Role = bl.Reader
	} else if claims.UserRole == "author" {
		reqUser.Role = bl.Author
	}

	myErr := app.IServices.IUsrSvc.DeleteUser(reqUser, targetUserID)
	if myErr.ErrNum == bl.ErrAccessDenied {
		app.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	} else if myErr.ErrNum == bl.ErrDeleteUser {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if myErr.ErrNum != bl.Ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Ответ об успешном удалении
	w.WriteHeader(http.StatusOK)
}

func (app *App) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	claims, err := bl.ValidateAndParseToken(r.Header.Get("Authorization"))
	if err != nil {
		app.Configs.LogConfigs.Logger.WriteLog("JWT error: "+err.Error(), slog.LevelError, nil)
		http.Error(w, "Authorization error", http.StatusUnauthorized)
		return
	}

	app.Configs.LogConfigs.Logger.WriteLog(fmt.Sprintf("User %d with role %s accessing UpdateUserHandler", claims.UserID, claims.UserRole), slog.LevelInfo, nil)

	vars := mux.Vars(r)
	targetUserID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	reqUser := &models.User{Id: claims.UserID}
	if claims.UserRole == "admin" {
		reqUser.Role = bl.Admin
	} else if claims.UserRole == "reader" {
		reqUser.Role = bl.Reader
	} else if claims.UserRole == "author" {
		reqUser.Role = bl.Author
	}

	var updatedUser models.User
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	updatedUser.Id = targetUserID

	myErr := app.IServices.IUsrSvc.UpdateUser(reqUser, &updatedUser)
	if myErr.ErrNum == bl.ErrAccessDenied {
		app.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	} else if myErr.ErrNum == bl.ErrDeleteUser {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if myErr.ErrNum != bl.Ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
