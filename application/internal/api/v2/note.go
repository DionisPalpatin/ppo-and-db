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

func (app *App) GetNoteHandler(w http.ResponseWriter, r *http.Request) {
	reqUser := getRequester(r, w, app.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	vars := mux.Vars(r)
	targetID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid note ID format", http.StatusBadRequest)
		return
	}

	srcData, myErr := app.IServices.INoteSvc.GetNote(targetID, "", bl.SearchByID, reqUser)
	if myErr == nil {
		app.Configs.LogConfigs.Logger.WriteLog("myErr is nil", slog.LevelError, nil)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	} else if myErr.ErrNum == bl.ErrAccessDenied {
		app.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	} else if myErr.ErrNum == bl.NoSuchNote {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	} else if myErr.ErrNum != bl.Ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	convertedData := converters.ToNoteFullData(srcData)

	w.Header().Set("Content-Type", "application/my_json")
	err = json.NewEncoder(w).Encode(convertedData)

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (app *App) GetAllNotesHandler(w http.ResponseWriter, r *http.Request) {
	reqUser := getRequester(r, w, app.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	queryParams := r.URL.Query()
	srcType := queryParams.Get("type")
	if srcType != "all" && srcType != "open" {
		http.Error(w, "Missing 'type' query parameter", http.StatusBadRequest)
		return
	}
	open := false
	if srcType == "open" {
		open = true
	}

	data, myErr := app.IServices.INoteSvc.GetAllNotes(open, reqUser)
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

	dataConverted := make([]transport_models.NoteInfo, 0, len(data))
	for _, el := range data {
		dataConverted = append(dataConverted, converters.ToNoteInfo(el))
	}

	w.Header().Set("Content-Type", "application/my_json")
	err := json.NewEncoder(w).Encode(dataConverted)

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (app *App) AddNoteHandler(w http.ResponseWriter, r *http.Request) {
	reqUser := getRequester(r, w, app.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	var srcData transport_models.NoteFullData
	err := json.NewDecoder(r.Body).Decode(&srcData)

	if err != nil {
		app.Configs.LogConfigs.Logger.WriteLog("Invalid request payload", slog.LevelError, nil)
		http.Error(w, "Invalid request payload", http.StatusUnprocessableEntity)
		return
	}

	var validate = validator.New()

	if err := validate.Struct(&srcData); err != nil {
		app.Configs.LogConfigs.Logger.WriteLog("Invalid request payload", slog.LevelError, nil)
		http.Error(w, "Invalid request payload", http.StatusUnprocessableEntity)
		return
	}

	convertedData, err := converters.FromNoteFullData(&srcData)
	myErr := app.IServices.INoteSvc.AddNote(&convertedData, reqUser)

	if myErr == nil {
		app.Configs.LogConfigs.Logger.WriteLog("myErr is nil", slog.LevelError, nil)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	} else if myErr.ErrNum == bl.NoSuchTeam {
		app.Configs.LogConfigs.Logger.WriteLog("Team not found", slog.LevelError, nil)
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	} else if myErr.ErrNum != bl.Ok {
		app.Configs.LogConfigs.Logger.WriteLog("Error checking user existence", slog.LevelError, nil)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *App) DeleteNoteHandler(w http.ResponseWriter, r *http.Request) {
	reqUser := getRequester(r, w, app.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	vars := mux.Vars(r)
	targetID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid note ID format", http.StatusBadRequest)
		return
	}

	myErr := app.IServices.INoteSvc.DeleteNote(targetID, reqUser)
	if myErr.ErrNum == bl.ErrAccessDenied {
		app.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	} else if myErr.ErrNum == bl.OperationError {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	} else if myErr.ErrNum != bl.Ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *App) UpdateNoteHandler(w http.ResponseWriter, r *http.Request) {
	reqUser := getRequester(r, w, app.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	var data transport_models.NoteFullData
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
		http.Error(w, "Invalid note ID format", http.StatusBadRequest)
		return
	}

	convData, err := converters.FromNoteFullData(&data)
	convData.Id = targetID
	myErr := app.IServices.INoteSvc.UpdateNote(&convData, reqUser,
		"/tmp/tmp."+convData.Content.TextExt,
		"/tmp/tmp."+convData.Content.ImgExt,
		"/tmp/tmp."+convData.Content.RawExt)

	if myErr.ErrNum == bl.ErrAccessDenied {
		app.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	} else if myErr.ErrNum == bl.OperationError {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	} else if myErr.ErrNum != bl.Ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
