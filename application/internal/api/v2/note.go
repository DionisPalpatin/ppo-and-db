package handlersv2

import (
	"encoding/json"
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/api/v2/converters"
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/api/v2/transport_models"
	bl "github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/business_logic"
	_ "github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"strconv"
)

func (hs *HandlersStruct) GetNoteHandler(c *gin.Context) {
	hs.Configs.LogConfigs.Logger.WriteLog("Start get note handler", slog.LevelInfo, nil)

	reqUser := getRequester(c, hs.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	hs.Configs.LogConfigs.Logger.WriteLog("Requester is got", slog.LevelInfo, nil)

	targetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID format"})
		return
	}

	hs.Configs.LogConfigs.Logger.WriteLog("Start get note", slog.LevelInfo, nil)

	srcData, myErr := hs.IServices.INoteSvc.GetNote(targetID, "", bl.SearchByID, reqUser)
	if myErr == nil {
		hs.Configs.LogConfigs.Logger.WriteLog("myErr is nil", slog.LevelError, nil)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	} else if myErr.ErrNum == bl.ErrAccessDenied {
		hs.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	} else if myErr.ErrNum == bl.NoSuchNote {
		hs.Configs.LogConfigs.Logger.WriteLog("Note not found", slog.LevelError, nil)
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	} else if myErr.ErrNum != bl.Ok {
		hs.Configs.LogConfigs.Logger.WriteLog("Internal server error with num = "+strconv.Itoa(myErr.ErrNum), slog.LevelError, nil)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	hs.Configs.LogConfigs.Logger.WriteLog("Note is got. Start to send reply", slog.LevelInfo, nil)

	convertedData := converters.ToNoteFullData(srcData)

	c.Header("Content-Type", "application/my_json")
	c.JSON(http.StatusOK, convertedData)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		hs.Configs.LogConfigs.Logger.WriteLog("Internal error: "+err.Error(), slog.LevelError, nil)
		return
	}
}

func (hs *HandlersStruct) GetAllNotesHandler(c *gin.Context) {
	reqUser := getRequester(c, hs.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	queryParams := r.URL.Query()
	srcType := queryParams.Get("type")
	if srcType != "all" && srcType != "open" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'type' query parameter"})
		return
	}
	open := false
	if srcType == "open" {
		open = true
	}

	data, myErr := hs.IServices.INoteSvc.GetAllNotes(open, reqUser)
	if myErr == nil {
		hs.Configs.LogConfigs.Logger.WriteLog("myErr is nil", slog.LevelError, nil)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	} else if myErr.ErrNum == bl.ErrAccessDenied {
		hs.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	} else if myErr.ErrNum != bl.Ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	dataConverted := make([]transport_models.NoteInfo, 0, len(data))
	for _, el := range data {
		dataConverted = append(dataConverted, converters.ToNoteInfo(el))
	}

	c.Header("Content-Type", "application/my_json")
	c.JSON(http.StatusOK, dataConverted)
}

func (hs *HandlersStruct) AddNoteHandler(c *gin.Context) {
	reqUser := getRequester(c, hs.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	var srcData transport_models.NoteFullData
	err := c.ShouldBindJSON(&srcData)

	if err != nil {
		hs.Configs.LogConfigs.Logger.WriteLog("Invalid request payload", slog.LevelError, nil)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid request payload"})
		return
	}

	var validate = validator.New()

	if err := validate.Struct(&srcData); err != nil {
		hs.Configs.LogConfigs.Logger.WriteLog("Invalid request payload", slog.LevelError, nil)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid request payload"})
		return
	}

	convertedData, err := converters.FromNoteFullData(&srcData)
	var idStruct teamIDStruct
	var myErr *bl.MyError
	idStruct.TeamID, myErr = hs.IServices.INoteSvc.AddNote(&convertedData, reqUser)

	if myErr == nil {
		hs.Configs.LogConfigs.Logger.WriteLog("myErr is nil", slog.LevelError, nil)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	} else if myErr.ErrNum == bl.NoSuchTeam {
		hs.Configs.LogConfigs.Logger.WriteLog("Team not found", slog.LevelError, nil)
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	} else if myErr.ErrNum != bl.Ok {
		hs.Configs.LogConfigs.Logger.WriteLog("Error checking user existence", slog.LevelError, nil)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Header("Content-Type", "application/my_json")
	c.JSON(http.StatusOK, idStruct)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
}

func (hs *HandlersStruct) DeleteNoteHandler(c *gin.Context) {
	reqUser := getRequester(c, hs.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	targetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID format"})
		return
	}

	myErr := hs.IServices.INoteSvc.DeleteNote(targetID, reqUser)
	if myErr.ErrNum == bl.ErrAccessDenied {
		hs.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	} else if myErr.ErrNum == bl.OperationError {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	} else if myErr.ErrNum != bl.Ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Status(http.StatusOK)
}

func (hs *HandlersStruct) UpdateNoteHandler(c *gin.Context) {
	reqUser := getRequester(c, hs.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	var data transport_models.NoteFullData
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	var validate = validator.New()
	if err := validate.Struct(&data); err != nil {
		hs.Configs.LogConfigs.Logger.WriteLog("Invalid request payload", slog.LevelError, nil)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid request payload"})
		return
	}

	targetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID format"})
		return
	}

	convData, err := converters.FromNoteFullData(&data)
	convData.Id = targetID
	myErr := hs.IServices.INoteSvc.UpdateNote(&convData, reqUser,
		"/tmp/tmp."+convData.Content.TextExt,
		"/tmp/tmp."+convData.Content.ImgExt,
		"/tmp/tmp."+convData.Content.RawExt)

	if myErr.ErrNum == bl.ErrAccessDenied {
		hs.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	} else if myErr.ErrNum == bl.OperationError {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	} else if myErr.ErrNum != bl.Ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Status(http.StatusOK)
}
