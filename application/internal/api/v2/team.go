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

type requestUserIDBody struct {
	UserID int `json:"userID" validate:"required"`
}

func (app *App) GetAllTeamsHandler(w http.ResponseWriter, r *http.Request) {
	reqUser := getRequester(r, w, app.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	data, myErr := app.IServices.ITeamSvc.GetAllTeams(reqUser)
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

	dataConverted := make([]transport_models.TeamInfo, 0, len(data))
	for _, team := range data {
		dataConverted = append(dataConverted, converters.ToTeamInfo(team))
	}

	w.Header().Set("Content-Type", "application/my_json")
	err := json.NewEncoder(w).Encode(dataConverted)

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (app *App) GetTeamHandler(w http.ResponseWriter, r *http.Request) {
	reqUser := getRequester(r, w, app.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	vars := mux.Vars(r)
	targetID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	srcData, myErr := app.IServices.ITeamSvc.GetTeam(targetID, "", bl.SearchByID, reqUser)
	if myErr == nil {
		app.Configs.LogConfigs.Logger.WriteLog("myErr is nil", slog.LevelError, nil)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	} else if myErr.ErrNum == bl.ErrAccessDenied {
		app.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	} else if myErr.ErrNum == bl.NoSuchTeam {
		http.Error(w, "Team not found", http.StatusNotFound)
		return
	} else if myErr.ErrNum != bl.Ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	convertedData := converters.ToTeamFullInfo(srcData)

	w.Header().Set("Content-Type", "application/my_json")
	err = json.NewEncoder(w).Encode(convertedData)

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (app *App) AddTeamHandler(w http.ResponseWriter, r *http.Request) {
	reqUser := getRequester(r, w, app.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	var teamInfo transport_models.TeamInfo
	err := json.NewDecoder(r.Body).Decode(&teamInfo)

	if err != nil {
		app.Configs.LogConfigs.Logger.WriteLog("Invalid request payload", slog.LevelError, nil)
		http.Error(w, "Invalid request payload", http.StatusUnprocessableEntity)
		return
	}

	var validate = validator.New()

	if err := validate.Struct(&teamInfo); err != nil {
		app.Configs.LogConfigs.Logger.WriteLog("Invalid request payload", slog.LevelError, nil)
		http.Error(w, "Invalid request payload", http.StatusUnprocessableEntity)
		return
	}

	convertedData := converters.FromTeamInfo(teamInfo)
	myErr := app.IServices.ITeamSvc.AddTeam(reqUser, &convertedData)

	if myErr == nil {
		app.Configs.LogConfigs.Logger.WriteLog("myErr is nil", slog.LevelError, nil)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	} else if myErr.ErrNum == bl.NoSuchTeam {
		app.Configs.LogConfigs.Logger.WriteLog("Team not found", slog.LevelError, nil)
		http.Error(w, "Team not found", http.StatusNotFound)
		return
	} else if myErr.ErrNum != bl.Ok {
		app.Configs.LogConfigs.Logger.WriteLog("Error checking user existence", slog.LevelError, nil)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *App) DeleteTeamHandler(w http.ResponseWriter, r *http.Request) {
	reqUser := getRequester(r, w, app.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	vars := mux.Vars(r)
	targetID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid team ID format", http.StatusBadRequest)
		return
	}

	myErr := app.IServices.ITeamSvc.DeleteTeam(reqUser, targetID)
	if myErr.ErrNum == bl.ErrAccessDenied {
		app.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	} else if myErr.ErrNum == bl.OperationError {
		http.Error(w, "Team not found", http.StatusNotFound)
		return
	} else if myErr.ErrNum != bl.Ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *App) UpdateTeamHandler(w http.ResponseWriter, r *http.Request) {
	reqUser := getRequester(r, w, app.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	var data transport_models.TeamInfo
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var validate = validator.New()
	if err := validate.Struct(&data); err != nil {
		app.Configs.LogConfigs.Logger.WriteLog("Invalid request payload", slog.LevelError, nil)
		http.Error(w, "Invalid request payload", http.StatusUnprocessableEntity)
		return
	}

	vars := mux.Vars(r)
	targetID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	convData := converters.FromTeamInfo(data)
	convData.Id = targetID
	myErr := app.IServices.ITeamSvc.UpdateTeam(reqUser, &convData)

	if myErr.ErrNum == bl.ErrAccessDenied {
		app.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	} else if myErr.ErrNum == bl.OperationError {
		http.Error(w, "Team not found", http.StatusNotFound)
		return
	} else if myErr.ErrNum != bl.Ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *App) AddUserToTeamHandler(w http.ResponseWriter, r *http.Request) {
	reqUser := getRequester(r, w, app.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	var data requestUserIDBody
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var validate = validator.New()
	if err := validate.Struct(&data); err != nil {
		app.Configs.LogConfigs.Logger.WriteLog("Invalid request payload", slog.LevelError, nil)
		http.Error(w, "Invalid request payload", http.StatusUnprocessableEntity)
		return
	}

	vars := mux.Vars(r)
	targetID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid team ID format", http.StatusBadRequest)
		return
	}

	myErr := app.IServices.ITeamSvc.AddUserToTeam(reqUser, data.UserID, targetID)

	if myErr.ErrNum == bl.ErrAccessDenied {
		app.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	} else if myErr.ErrNum == bl.OperationError {
		http.Error(w, "Team not found", http.StatusNotFound)
		return
	} else if myErr.ErrNum != bl.Ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *App) DeleteUserFromTeamHandler(w http.ResponseWriter, r *http.Request) {
	reqUser := getRequester(r, w, app.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	vars := mux.Vars(r)
	teamID, err := strconv.Atoi(vars["teamID"])
	if err != nil {
		http.Error(w, "Invalid team ID format", http.StatusBadRequest)
		return
	}
	userID, err := strconv.Atoi(vars["userID"])
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	myErr := app.IServices.ITeamSvc.DeleteUserFromTeam(reqUser, userID, teamID)

	if myErr.ErrNum == bl.ErrAccessDenied {
		app.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	} else if myErr.ErrNum == bl.OperationError {
		http.Error(w, "Team not found", http.StatusNotFound)
		return
	} else if myErr.ErrNum != bl.Ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *App) GetTeamMembersHandler(w http.ResponseWriter, r *http.Request) {
	reqUser := getRequester(r, w, app.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	vars := mux.Vars(r)
	teamID, err := strconv.Atoi(vars["teamID"])
	if err != nil {
		http.Error(w, "Invalid team ID format", http.StatusBadRequest)
		return
	}

	data, myErr := app.IServices.ITeamSvc.GetTeamMembers(teamID, reqUser)
	if myErr.ErrNum == bl.ErrAccessDenied {
		app.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	} else if myErr.ErrNum == bl.OperationError {
		http.Error(w, "Team not found", http.StatusNotFound)
		return
	} else if myErr.ErrNum != bl.Ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	convertedData := make([]transport_models.UserPublicInfo, 0, len(data))
	for _, user := range data {
		convertedData = append(convertedData, converters.ToUserPublicInfo(user))
	}

	w.Header().Set("Content-Type", "application/my_json")
	err = json.NewEncoder(w).Encode(convertedData)

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
