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

func (hs *HandlersStruct) GetSectionHandler(w http.ResponseWriter, r *http.Request) {
	reqUser := getRequester(r, w, hs.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	vars := mux.Vars(r)
	targetID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid section ID format", http.StatusBadRequest)
		return
	}

	srcData, myErr := hs.IServices.ISecSvc.GetSection(targetID, "", reqUser, bl.SearchByID)
	if myErr == nil {
		hs.Configs.LogConfigs.Logger.WriteLog("myErr is nil", slog.LevelError, nil)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	} else if myErr.ErrNum == bl.ErrAccessDenied {
		hs.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	} else if myErr.ErrNum == bl.NoSuchColl {
		http.Error(w, "Section not found", http.StatusNotFound)
		return
	} else if myErr.ErrNum != bl.Ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	team, myErr := hs.IServices.ITeamSvc.GetSectionTeam(srcData.Id, reqUser)
	if myErr == nil {
		hs.Configs.LogConfigs.Logger.WriteLog("myErr is nil", slog.LevelError, nil)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	} else if myErr.ErrNum == bl.ErrAccessDenied {
		hs.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	} else if myErr.ErrNum == bl.NoSuchColl {
		http.Error(w, "Team not found", http.StatusNotFound)
		return
	} else if myErr.ErrNum != bl.Ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	convertedData := converters.ToSectionFullInfo(srcData, team)

	w.Header().Set("Content-Type", "application/my_json")
	err = json.NewEncoder(w).Encode(convertedData)

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (hs *HandlersStruct) GetAllSectionsHandler(w http.ResponseWriter, r *http.Request) {
	reqUser := getRequester(r, w, hs.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	data, myErr := hs.IServices.ISecSvc.GetAllSections(reqUser)
	if myErr == nil {
		hs.Configs.LogConfigs.Logger.WriteLog("myErr is nil", slog.LevelError, nil)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	} else if myErr.ErrNum == bl.ErrAccessDenied {
		hs.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	} else if myErr.ErrNum != bl.Ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	dataConverted := make([]transport_models.SectionInfo, 0, len(data))
	for _, el := range data {
		team, myErr := hs.IServices.ITeamSvc.GetSectionTeam(el.Id, reqUser)

		if myErr == nil {
			hs.Configs.LogConfigs.Logger.WriteLog("myErr is nil", slog.LevelError, nil)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		} else if myErr.ErrNum == bl.ErrAccessDenied {
			hs.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
			http.Error(w, "Insufficient permissions", http.StatusForbidden)
			return
		} else if myErr.ErrNum == bl.NoSuchColl {
			http.Error(w, "Team not found", http.StatusNotFound)
			return
		} else if myErr.ErrNum != bl.Ok {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		dataConverted = append(dataConverted, converters.ToSectionInfo(el, team))
	}

	w.Header().Set("Content-Type", "application/my_json")
	err := json.NewEncoder(w).Encode(dataConverted)

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (hs *HandlersStruct) AddSectionHandler(w http.ResponseWriter, r *http.Request) {
	reqUser := getRequester(r, w, hs.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	var srcData transport_models.SectionInfo
	err := json.NewDecoder(r.Body).Decode(&srcData)

	if err != nil {
		hs.Configs.LogConfigs.Logger.WriteLog("Invalid request payload", slog.LevelError, nil)
		http.Error(w, "Invalid request payload", http.StatusUnprocessableEntity)
		return
	}

	var validate = validator.New()

	if err := validate.Struct(&srcData); err != nil {
		hs.Configs.LogConfigs.Logger.WriteLog("Invalid request payload", slog.LevelError, nil)
		http.Error(w, "Invalid request payload", http.StatusUnprocessableEntity)
		return
	}

	convertedData, teamID := converters.FromSectionInfo(&srcData)

	team, myErr := hs.IServices.ITeamSvc.GetTeam(teamID, "", bl.SearchByID, reqUser)

	if myErr == nil {
		hs.Configs.LogConfigs.Logger.WriteLog("myErr is nil", slog.LevelError, nil)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	} else if myErr.ErrNum == bl.ErrAccessDenied {
		hs.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	} else if myErr.ErrNum == bl.NoSuchColl {
		http.Error(w, "Team not found", http.StatusNotFound)
		return
	} else if myErr.ErrNum != bl.Ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var idStruct sectionIDStruct
	idStruct.SecID, myErr = hs.IServices.ISecSvc.AddSection(&convertedData, team, reqUser)

	if myErr == nil {
		hs.Configs.LogConfigs.Logger.WriteLog("myErr is nil", slog.LevelError, nil)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	} else if myErr.ErrNum == bl.NoSuchTeam {
		hs.Configs.LogConfigs.Logger.WriteLog("Team not found", slog.LevelError, nil)
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	} else if myErr.ErrNum != bl.Ok {
		hs.Configs.LogConfigs.Logger.WriteLog("Error checking user existence", slog.LevelError, nil)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/my_json")
	err = json.NewEncoder(w).Encode(idStruct)

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (hs *HandlersStruct) DeleteSectionHandler(w http.ResponseWriter, r *http.Request) {
	reqUser := getRequester(r, w, hs.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	vars := mux.Vars(r)
	targetID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid section ID format", http.StatusBadRequest)
		return
	}

	myErr := hs.IServices.ISecSvc.DeleteSection(targetID, reqUser)
	if myErr.ErrNum == bl.ErrAccessDenied {
		hs.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	} else if myErr.ErrNum == bl.OperationError {
		http.Error(w, "Section not found", http.StatusNotFound)
		return
	} else if myErr.ErrNum != bl.Ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (hs *HandlersStruct) UpdateSectionHandler(w http.ResponseWriter, r *http.Request) {
	reqUser := getRequester(r, w, hs.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	var data transport_models.SectionInfo
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var validate = validator.New()

	if err := validate.Struct(&data); err != nil {
		hs.Configs.LogConfigs.Logger.WriteLog("Invalid request payload", slog.LevelError, nil)
		http.Error(w, "Invalid request payload", http.StatusUnprocessableEntity)
		return
	}

	vars := mux.Vars(r)
	targetID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid section ID format", http.StatusBadRequest)
		return
	}

	convData, _ := converters.FromSectionInfo(&data)
	convData.Id = targetID
	myErr := hs.IServices.ISecSvc.UpdateSection(&convData, reqUser)

	if myErr.ErrNum == bl.ErrAccessDenied {
		hs.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	} else if myErr.ErrNum == bl.OperationError {
		http.Error(w, "Section not found", http.StatusNotFound)
		return
	} else if myErr.ErrNum != bl.Ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (hs *HandlersStruct) GetAllNotesInSectionHandler(w http.ResponseWriter, r *http.Request) {
	reqUser := getRequester(r, w, hs.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	vars := mux.Vars(r)
	targetID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid section ID format", http.StatusBadRequest)
		return
	}

	data, myErr := hs.IServices.ISecSvc.GetAllNotesInSection(targetID, reqUser)
	if myErr == nil {
		hs.Configs.LogConfigs.Logger.WriteLog("myErr is nil", slog.LevelError, nil)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	} else if myErr.ErrNum == bl.ErrAccessDenied {
		hs.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	} else if myErr.ErrNum == bl.NoSuchColl {
		http.Error(w, "Section not found", http.StatusNotFound)
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
	err = json.NewEncoder(w).Encode(dataConverted)

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (hs *HandlersStruct) AddNoteToSectionHandler(w http.ResponseWriter, r *http.Request) {
	reqUser := getRequester(r, w, hs.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	var data noteIDStruct
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var validate = validator.New()
	if err := validate.Struct(&data); err != nil {
		hs.Configs.LogConfigs.Logger.WriteLog("Invalid request payload", slog.LevelError, nil)
		http.Error(w, "Invalid request payload", http.StatusUnprocessableEntity)
		return
	}

	vars := mux.Vars(r)
	targetID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid section ID format", http.StatusBadRequest)
		return
	}

	section, myErr := hs.IServices.ISecSvc.GetSection(targetID, "", reqUser, bl.SearchByID)
	if myErr == nil {
		hs.Configs.LogConfigs.Logger.WriteLog("myErr is nil", slog.LevelError, nil)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	} else if myErr.ErrNum == bl.ErrAccessDenied {
		hs.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	} else if myErr.ErrNum == bl.NoSuchColl {
		http.Error(w, "Section not found", http.StatusNotFound)
		return
	} else if myErr.ErrNum != bl.Ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	note, myErr := hs.IServices.INoteSvc.GetNote(targetID, "", bl.SearchByID, reqUser)
	if myErr == nil {
		hs.Configs.LogConfigs.Logger.WriteLog("myErr is nil", slog.LevelError, nil)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	} else if myErr.ErrNum == bl.ErrAccessDenied {
		hs.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	} else if myErr.ErrNum == bl.NoSuchColl {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	} else if myErr.ErrNum != bl.Ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	myErr = hs.IServices.ISecSvc.AddNoteToSection(section, note, reqUser)

	if myErr.ErrNum == bl.ErrAccessDenied {
		hs.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	} else if myErr.ErrNum == bl.OperationError {
		http.Error(w, "Section not found", http.StatusNotFound)
		return
	} else if myErr.ErrNum != bl.Ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (hs *HandlersStruct) DeleteNoteFromSectionHandler(w http.ResponseWriter, r *http.Request) {
	reqUser := getRequester(r, w, hs.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	var data noteIDStruct
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var validate = validator.New()
	if err := validate.Struct(&data); err != nil {
		hs.Configs.LogConfigs.Logger.WriteLog("Invalid request payload", slog.LevelError, nil)
		http.Error(w, "Invalid request payload", http.StatusUnprocessableEntity)
		return
	}

	vars := mux.Vars(r)
	targetID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid section ID format", http.StatusBadRequest)
		return
	}

	section, myErr := hs.IServices.ISecSvc.GetSection(targetID, "", reqUser, bl.SearchByID)
	if myErr == nil {
		hs.Configs.LogConfigs.Logger.WriteLog("myErr is nil", slog.LevelError, nil)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	} else if myErr.ErrNum == bl.ErrAccessDenied {
		hs.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	} else if myErr.ErrNum == bl.NoSuchColl {
		http.Error(w, "Section not found", http.StatusNotFound)
		return
	} else if myErr.ErrNum != bl.Ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	note, myErr := hs.IServices.INoteSvc.GetNote(targetID, "", bl.SearchByID, reqUser)
	if myErr == nil {
		hs.Configs.LogConfigs.Logger.WriteLog("myErr is nil", slog.LevelError, nil)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	} else if myErr.ErrNum == bl.ErrAccessDenied {
		hs.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	} else if myErr.ErrNum == bl.NoSuchColl {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	} else if myErr.ErrNum != bl.Ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	myErr = hs.IServices.ISecSvc.DeleteNoteFromSection(section, note, reqUser)

	if myErr.ErrNum == bl.ErrAccessDenied {
		hs.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	} else if myErr.ErrNum == bl.OperationError {
		http.Error(w, "Section not found", http.StatusNotFound)
		return
	} else if myErr.ErrNum != bl.Ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
