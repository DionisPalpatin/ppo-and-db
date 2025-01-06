package handlersv2

import (
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/api/v2/converters"
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/api/v2/transport_models"
	bl "github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/business_logic"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"strconv"
)

func (hs *HandlersStruct) GetCollectionHandler(c *gin.Context) {
	reqUser := getRequester(c, hs.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	targetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid collection ID format"})
		return
	}

	srcData, myErr := hs.IServices.IColSvc.GetCollection(targetID, "", bl.SearchByID)
	if myErr == nil {
		hs.Configs.LogConfigs.Logger.WriteLog("myErr is nil", slog.LevelError, nil)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	} else if myErr.ErrNum == bl.ErrAccessDenied {
		hs.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	} else if myErr.ErrNum == bl.NoSuchColl {
		c.JSON(http.StatusNotFound, gin.H{"error": "Collection not found"})
		return
	} else if myErr.ErrNum != bl.Ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	convertedData := converters.ToCollectionFullInfo(srcData)

	c.Header("Content-Type", "application/my_json")
	c.JSON(http.StatusOK, convertedData)
}

func (hs *HandlersStruct) GetAllCollectionsHandler(c *gin.Context) {
	reqUser := getRequester(c, hs.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	data, myErr := hs.IServices.IColSvc.GetAllCollections(reqUser)
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

	dataConverted := make([]transport_models.CollectionInfo, 0, len(data))
	for _, el := range data {
		dataConverted = append(dataConverted, converters.ToCollectionInfo(el))
	}

	c.Header("Content-Type", "application/my_json")
	c.JSON(http.StatusOK, dataConverted)
}

func (hs *HandlersStruct) GetAllUsersCollectionsHandler(c *gin.Context) {
	reqUser := getRequester(c, hs.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	data, myErr := hs.IServices.IColSvc.GetAllUsersCollections(reqUser)
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

	dataConverted := make([]transport_models.CollectionInfo, 0, len(data))
	for _, el := range data {
		dataConverted = append(dataConverted, converters.ToCollectionInfo(el))
	}

	c.Header("Content-Type", "application/my_json")
	c.JSON(http.StatusOK, dataConverted)
}

func (hs *HandlersStruct) AddCollectionHandler(c *gin.Context) {
	reqUser := getRequester(c, hs.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	var srcData transport_models.CollectionAddInfo
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

	convertedData := converters.FromCollectionAddInfo(&srcData)
	var idStruct teamIDStruct
	var myErr *bl.MyError
	idStruct.TeamID, myErr = hs.IServices.IColSvc.AddCollection(&convertedData)

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
}

func (hs *HandlersStruct) DeleteCollectionHandler(c *gin.Context) {
	reqUser := getRequester(c, hs.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	targetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid collection ID format"})
		return
	}

	myErr := hs.IServices.IColSvc.DeleteCollection(targetID, reqUser)
	if myErr.ErrNum == bl.ErrAccessDenied {
		hs.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	} else if myErr.ErrNum == bl.OperationError {
		c.JSON(http.StatusNotFound, gin.H{"error": "Collection not found"})
		return
	} else if myErr.ErrNum != bl.Ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Status(http.StatusOK)
}

func (hs *HandlersStruct) UpdateCollectionHandler(c *gin.Context) {
	reqUser := getRequester(c, hs.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	var data transport_models.CollectionAddInfo
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid collection ID format"})
		return
	}

	convData := converters.FromCollectionAddInfo(&data)
	convData.Id = targetID
	myErr := hs.IServices.IColSvc.UpdateCollection(&convData)

	if myErr.ErrNum == bl.ErrAccessDenied {
		hs.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	} else if myErr.ErrNum == bl.OperationError {
		c.JSON(http.StatusNotFound, gin.H{"error": "Collection not found"})
		return
	} else if myErr.ErrNum != bl.Ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Status(http.StatusOK)
}

func (hs *HandlersStruct) GetAllNotesInCollectionHandler(c *gin.Context) {
	reqUser := getRequester(c, hs.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	targetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid collection ID format"})
		return
	}

	collection, myErr := hs.IServices.IColSvc.GetCollection(targetID, "", bl.SearchByID)
	if myErr == nil {
		hs.Configs.LogConfigs.Logger.WriteLog("myErr is nil", slog.LevelError, nil)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	} else if myErr.ErrNum == bl.ErrAccessDenied {
		hs.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	} else if myErr.ErrNum == bl.NoSuchColl {
		c.JSON(http.StatusNotFound, gin.H{"error": "Collection not found"})
		return
	} else if myErr.ErrNum != bl.Ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	data, myErr := hs.IServices.IColSvc.GetAllNotesInCollection(collection)
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

func (hs *HandlersStruct) AddNoteToCollectionHandler(c *gin.Context) {
	reqUser := getRequester(c, hs.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	var data noteIDStruct
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid collection ID format"})
		return
	}

	myErr := hs.IServices.INoteSvc.AddNoteToCollection(data.NoteID, targetID)

	if myErr.ErrNum == bl.ErrAccessDenied {
		hs.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	} else if myErr.ErrNum == bl.OperationError {
		c.JSON(http.StatusNotFound, gin.H{"error": "Collection not found"})
		return
	} else if myErr.ErrNum != bl.Ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Status(http.StatusOK)
}

func (hs *HandlersStruct) DeleteNoteFromCollectionHandler(c *gin.Context) {
	reqUser := getRequester(c, hs.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	noteID, err := strconv.Atoi(c.Param("note_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid collection ID format"})
		return
	}
	collID, err := strconv.Atoi(c.Param("coll_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid collection ID format"})
		return
	}

	myErr := hs.IServices.INoteSvc.DeleteNoteFromCollection(noteID, collID)

	if myErr.ErrNum == bl.ErrAccessDenied {
		hs.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	} else if myErr.ErrNum == bl.OperationError {
		c.JSON(http.StatusNotFound, gin.H{"error": "Collection not found"})
		return
	} else if myErr.ErrNum != bl.Ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Status(http.StatusOK)
}
