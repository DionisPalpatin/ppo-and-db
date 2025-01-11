package handlersv2

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/DionisPalpatin/ppo-and-db/application/internal/api/converters"
	"github.com/DionisPalpatin/ppo-and-db/application/internal/api/transport_models"
	bl "github.com/DionisPalpatin/ppo-and-db/application/internal/business-logic"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func (hs *HandlersStruct) GetSectionHandler(c *gin.Context) {
	reqUser := getRequester(c, hs.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	targetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid section ID format"})
		return
	}

	srcData, myErr := hs.IServices.ISecSvc.GetSection(targetID, "", reqUser, bl.SearchByID)
	if myErr == nil {
		hs.Configs.LogConfigs.Logger.WriteLog("myErr is nil", slog.LevelError, nil)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	} else if myErr.ErrNum == bl.ErrAccessDenied {
		hs.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	} else if myErr.ErrNum == bl.NoSuchColl {
		c.JSON(http.StatusNotFound, gin.H{"error": "Section not found"})
		return
	} else if myErr.ErrNum != bl.Ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	team, myErr := hs.IServices.ITeamSvc.GetSectionTeam(srcData.Id, reqUser)
	if myErr == nil {
		hs.Configs.LogConfigs.Logger.WriteLog("myErr is nil", slog.LevelError, nil)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	} else if myErr.ErrNum == bl.ErrAccessDenied {
		hs.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	} else if myErr.ErrNum == bl.NoSuchColl {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	} else if myErr.ErrNum != bl.Ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	convertedData := converters.ToSectionFullInfo(srcData, team)

	c.Header("Content-Type", "application/my_json")
	c.JSON(http.StatusOK, convertedData)
}

func (hs *HandlersStruct) GetAllSectionsHandler(c *gin.Context) {
	reqUser := getRequester(c, hs.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	data, myErr := hs.IServices.ISecSvc.GetAllSections(reqUser)
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

	dataConverted := make([]transport_models.SectionInfo, 0, len(data))
	for _, el := range data {
		team, myErr := hs.IServices.ITeamSvc.GetSectionTeam(el.Id, reqUser)

		if myErr == nil {
			hs.Configs.LogConfigs.Logger.WriteLog("myErr is nil", slog.LevelError, nil)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		} else if myErr.ErrNum == bl.ErrAccessDenied {
			hs.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			return
		} else if myErr.ErrNum == bl.NoSuchColl {
			c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
			return
		} else if myErr.ErrNum != bl.Ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		dataConverted = append(dataConverted, converters.ToSectionInfo(el, team))
	}

	c.Header("Content-Type", "application/my_json")
	c.JSON(http.StatusOK, dataConverted)
}

func (hs *HandlersStruct) AddSectionHandler(c *gin.Context) {
	reqUser := getRequester(c, hs.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	var srcData transport_models.SectionInfo
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

	convertedData, teamID := converters.FromSectionInfo(&srcData)

	team, myErr := hs.IServices.ITeamSvc.GetTeam(teamID, "", bl.SearchByID, reqUser)

	if myErr == nil {
		hs.Configs.LogConfigs.Logger.WriteLog("myErr is nil", slog.LevelError, nil)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	} else if myErr.ErrNum == bl.ErrAccessDenied {
		hs.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	} else if myErr.ErrNum == bl.NoSuchColl {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	} else if myErr.ErrNum != bl.Ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var idStruct sectionIDStruct
	idStruct.SecID, myErr = hs.IServices.ISecSvc.AddSection(&convertedData, team, reqUser)

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

func (hs *HandlersStruct) DeleteSectionHandler(c *gin.Context) {
	reqUser := getRequester(c, hs.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	targetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid section ID format"})
		return
	}

	myErr := hs.IServices.ISecSvc.DeleteSection(targetID, reqUser)
	if myErr.ErrNum == bl.ErrAccessDenied {
		hs.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	} else if myErr.ErrNum == bl.OperationError {
		c.JSON(http.StatusNotFound, gin.H{"error": "Section not found"})
		return
	} else if myErr.ErrNum != bl.Ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Status(http.StatusOK)
}

func (hs *HandlersStruct) UpdateSectionHandler(c *gin.Context) {
	reqUser := getRequester(c, hs.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	var data transport_models.SectionInfo
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid section ID format"})
		return
	}

	convData, _ := converters.FromSectionInfo(&data)
	convData.Id = targetID
	myErr := hs.IServices.ISecSvc.UpdateSection(&convData, reqUser)

	if myErr.ErrNum == bl.ErrAccessDenied {
		hs.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	} else if myErr.ErrNum == bl.OperationError {
		c.JSON(http.StatusNotFound, gin.H{"error": "Section not found"})
		return
	} else if myErr.ErrNum != bl.Ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Status(http.StatusOK)
}

func (hs *HandlersStruct) GetAllNotesInSectionHandler(c *gin.Context) {
	reqUser := getRequester(c, hs.Configs.LogConfigs.Logger)
	if reqUser == nil {
		return
	}

	targetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid section ID format"})
		return
	}

	data, myErr := hs.IServices.ISecSvc.GetAllNotesInSection(targetID, reqUser)
	if myErr == nil {
		hs.Configs.LogConfigs.Logger.WriteLog("myErr is nil", slog.LevelError, nil)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	} else if myErr.ErrNum == bl.ErrAccessDenied {
		hs.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	} else if myErr.ErrNum == bl.NoSuchColl {
		c.JSON(http.StatusNotFound, gin.H{"error": "Section not found"})
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

func (hs *HandlersStruct) AddNoteToSectionHandler(c *gin.Context) {
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid section ID format"})
		return
	}

	section, myErr := hs.IServices.ISecSvc.GetSection(targetID, "", reqUser, bl.SearchByID)
	if myErr == nil {
		hs.Configs.LogConfigs.Logger.WriteLog("myErr is nil", slog.LevelError, nil)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	} else if myErr.ErrNum == bl.ErrAccessDenied {
		hs.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	} else if myErr.ErrNum == bl.NoSuchColl {
		c.JSON(http.StatusNotFound, gin.H{"error": "Section not found"})
		return
	} else if myErr.ErrNum != bl.Ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	note, myErr := hs.IServices.INoteSvc.GetNote(targetID, "", bl.SearchByID, reqUser)
	if myErr == nil {
		hs.Configs.LogConfigs.Logger.WriteLog("myErr is nil", slog.LevelError, nil)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	} else if myErr.ErrNum == bl.ErrAccessDenied {
		hs.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	} else if myErr.ErrNum == bl.NoSuchColl {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	} else if myErr.ErrNum != bl.Ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	myErr = hs.IServices.ISecSvc.AddNoteToSection(section, note, reqUser)

	if myErr.ErrNum == bl.ErrAccessDenied {
		hs.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	} else if myErr.ErrNum == bl.OperationError {
		c.JSON(http.StatusNotFound, gin.H{"error": "Section not found"})
		return
	} else if myErr.ErrNum != bl.Ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Status(http.StatusOK)
}

func (hs *HandlersStruct) DeleteNoteFromSectionHandler(c *gin.Context) {
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid section ID format"})
		return
	}

	section, myErr := hs.IServices.ISecSvc.GetSection(targetID, "", reqUser, bl.SearchByID)
	if myErr == nil {
		hs.Configs.LogConfigs.Logger.WriteLog("myErr is nil", slog.LevelError, nil)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	} else if myErr.ErrNum == bl.ErrAccessDenied {
		hs.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	} else if myErr.ErrNum == bl.NoSuchColl {
		c.JSON(http.StatusNotFound, gin.H{"error": "Section not found"})
		return
	} else if myErr.ErrNum != bl.Ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	note, myErr := hs.IServices.INoteSvc.GetNote(targetID, "", bl.SearchByID, reqUser)
	if myErr == nil {
		hs.Configs.LogConfigs.Logger.WriteLog("myErr is nil", slog.LevelError, nil)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	} else if myErr.ErrNum == bl.ErrAccessDenied {
		hs.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	} else if myErr.ErrNum == bl.NoSuchColl {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	} else if myErr.ErrNum != bl.Ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	myErr = hs.IServices.ISecSvc.DeleteNoteFromSection(section, note, reqUser)

	if myErr.ErrNum == bl.ErrAccessDenied {
		hs.Configs.LogConfigs.Logger.WriteLog("Insufficient permissions", slog.LevelError, nil)
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	} else if myErr.ErrNum == bl.OperationError {
		c.JSON(http.StatusNotFound, gin.H{"error": "Section not found"})
		return
	} else if myErr.ErrNum != bl.Ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Status(http.StatusOK)
}
