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

func (hs *HandlersStruct) GetCollectionHandler(w http.ResponseWriter, r *http.Request) {
	reqUser := getRequester(r, w, hs.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	vars := mux.Vars(r)
	targetID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid collection ID format", http.StatusBadRequest)
		return
	}

	srcData, myErr := hs.IServices.IColSvc.GetCollection(targetID, "", bl.SearchByID)
	if myErr == nil {
		hs.Configs.LogConfigs.Logger.WriteLog("myErr is nil", slog.LevelError, nil)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	} else if myErr.ErrNum == bl.ErrAccessDenied {
		hs.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	} else if myErr.ErrNum == bl.NoSuchColl {
		http.Error(w, "Collection not found", http.StatusNotFound)
		return
	} else if myErr.ErrNum != bl.Ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	convertedData := converters.ToCollectionFullInfo(srcData)

	w.Header().Set("Content-Type", "application/my_json")
	err = json.NewEncoder(w).Encode(convertedData)

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (hs *HandlersStruct) GetAllCollectionsHandler(w http.ResponseWriter, r *http.Request) {
	reqUser := getRequester(r, w, hs.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	data, myErr := hs.IServices.IColSvc.GetAllCollections(reqUser)
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

	dataConverted := make([]transport_models.CollectionInfo, 0, len(data))
	for _, el := range data {
		dataConverted = append(dataConverted, converters.ToCollectionInfo(el))
	}

	w.Header().Set("Content-Type", "application/my_json")
	err := json.NewEncoder(w).Encode(dataConverted)

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (hs *HandlersStruct) GetAllUsersCollectionsHandler(w http.ResponseWriter, r *http.Request) {
	reqUser := getRequester(r, w, hs.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	data, myErr := hs.IServices.IColSvc.GetAllUsersCollections(reqUser)
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

	dataConverted := make([]transport_models.CollectionInfo, 0, len(data))
	for _, el := range data {
		dataConverted = append(dataConverted, converters.ToCollectionInfo(el))
	}

	w.Header().Set("Content-Type", "application/my_json")
	err := json.NewEncoder(w).Encode(dataConverted)

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (hs *HandlersStruct) AddCollectionHandler(w http.ResponseWriter, r *http.Request) {
	reqUser := getRequester(r, w, hs.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	var srcData transport_models.CollectionAddInfo
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

	convertedData := converters.FromCollectionAddInfo(&srcData)
	var idStruct teamIDStruct
	var myErr *bl.MyError
	idStruct.TeamID, myErr = hs.IServices.IColSvc.AddCollection(&convertedData)

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

func (hs *HandlersStruct) DeleteCollectionHandler(w http.ResponseWriter, r *http.Request) {
	reqUser := getRequester(r, w, hs.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	vars := mux.Vars(r)
	targetID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid collection ID format", http.StatusBadRequest)
		return
	}

	myErr := hs.IServices.IColSvc.DeleteCollection(targetID, reqUser)
	if myErr.ErrNum == bl.ErrAccessDenied {
		hs.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	} else if myErr.ErrNum == bl.OperationError {
		http.Error(w, "Collection not found", http.StatusNotFound)
		return
	} else if myErr.ErrNum != bl.Ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (hs *HandlersStruct) UpdateCollectionHandler(w http.ResponseWriter, r *http.Request) {
	reqUser := getRequester(r, w, hs.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	var data transport_models.CollectionAddInfo
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
		http.Error(w, "Invalid collection ID format", http.StatusBadRequest)
		return
	}

	convData := converters.FromCollectionAddInfo(&data)
	convData.Id = targetID
	myErr := hs.IServices.IColSvc.UpdateCollection(&convData)

	if myErr.ErrNum == bl.ErrAccessDenied {
		hs.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	} else if myErr.ErrNum == bl.OperationError {
		http.Error(w, "Collection not found", http.StatusNotFound)
		return
	} else if myErr.ErrNum != bl.Ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (hs *HandlersStruct) GetAllNotesInCollectionHandler(w http.ResponseWriter, r *http.Request) {
	reqUser := getRequester(r, w, hs.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	vars := mux.Vars(r)
	targetID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid collection ID format", http.StatusBadRequest)
		return
	}

	collection, myErr := hs.IServices.IColSvc.GetCollection(targetID, "", bl.SearchByID)
	if myErr == nil {
		hs.Configs.LogConfigs.Logger.WriteLog("myErr is nil", slog.LevelError, nil)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	} else if myErr.ErrNum == bl.ErrAccessDenied {
		hs.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	} else if myErr.ErrNum == bl.NoSuchColl {
		http.Error(w, "Collection not found", http.StatusNotFound)
		return
	} else if myErr.ErrNum != bl.Ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	data, myErr := hs.IServices.IColSvc.GetAllNotesInCollection(collection)
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

func (hs *HandlersStruct) AddNoteToCollectionHandler(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "Invalid collection ID format", http.StatusBadRequest)
		return
	}

	myErr := hs.IServices.INoteSvc.AddNoteToCollection(data.NoteID, targetID)

	if myErr.ErrNum == bl.ErrAccessDenied {
		hs.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	} else if myErr.ErrNum == bl.OperationError {
		http.Error(w, "Collection not found", http.StatusNotFound)
		return
	} else if myErr.ErrNum != bl.Ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (hs *HandlersStruct) DeleteNoteFromCollectionHandler(w http.ResponseWriter, r *http.Request) {
	reqUser := getRequester(r, w, hs.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	vars := mux.Vars(r)
	noteID, err := strconv.Atoi(vars["note_id"])
	if err != nil {
		http.Error(w, "Invalid collection ID format", http.StatusBadRequest)
		return
	}
	collID, err := strconv.Atoi(vars["coll_id"])
	if err != nil {
		http.Error(w, "Invalid collection ID format", http.StatusBadRequest)
		return
	}

	myErr := hs.IServices.INoteSvc.DeleteNoteFromCollection(noteID, collID)

	if myErr.ErrNum == bl.ErrAccessDenied {
		hs.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	} else if myErr.ErrNum == bl.OperationError {
		http.Error(w, "Collection not found", http.StatusNotFound)
		return
	} else if myErr.ErrNum != bl.Ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
