package handlersv2

import (
	"github.com/DionisPalpatin/ppo-and-db/application/internal/config"
	// swaggerfiles "github.com/swaggo/files"
	// ginSwagger "github.com/swaggo/gin-swagger"

	// _ "github.com/DionisPalpatin/ppo-and-db/application/docs"

	"net/http"

	bl "github.com/DionisPalpatin/ppo-and-db/application/internal/business-logic"

	"github.com/gin-gonic/gin"
)

type HandlersStruct struct {
	Configs *config.Configs
	Router  *gin.Engine

	IServices *bl.IServices
	IRepos    *bl.IRepositories
}

func InitHandlersStructFields(hs *HandlersStruct, confs *config.Configs, ireps *bl.IRepositories, iservs *bl.IServices) {
	hs.Configs = confs
	hs.IServices = iservs
	hs.IRepos = ireps
}

func InitRouter(hs *HandlersStruct) {
	hs.Router = gin.Default()
	router := hs.Router

	// -----------------------------------------------------------------------------------------------------------------
	// Authorisation handlers
	// -----------------------------------------------------------------------------------------------------------------
	router.POST("/api/v1/login", hs.LoginHandler)
	router.POST("/api/v1/register", hs.RegisterHandler)

	// -----------------------------------------------------------------------------------------------------------------
	// User handlers
	// -----------------------------------------------------------------------------------------------------------------
	router.GET("/api/v1/users", hs.GetAllUsersHandler)
	router.GET("/api/v1/users/:id", hs.GetUserHandler)
	router.DELETE("/api/v1/users/:id", hs.DeleteUserHandler)
	router.PATCH("/api/v1/users/:id", hs.UpdateUserHandler)

	// -----------------------------------------------------------------------------------------------------------------
	// Team handlers
	// -----------------------------------------------------------------------------------------------------------------
	router.GET("/api/v1/teams", hs.GetAllTeamsHandler)
	router.POST("/api/v1/teams", hs.AddTeamHandler)
	router.DELETE("/api/v1/teams/:teamID/members/:userID", hs.DeleteUserFromTeamHandler)
	router.GET("/api/v1/teams/:teamID/members", hs.GetTeamMembersHandler)
	router.POST("/api/v1/teams/:teamID/members", hs.AddUserToTeamHandler)
	router.GET("/api/v1/teams/:teamID", hs.GetTeamHandler)
	router.DELETE("/api/v1/teams/:teamID", hs.DeleteTeamHandler)
	router.PATCH("/api/v1/teams/:teamID", hs.UpdateTeamHandler)

	// -----------------------------------------------------------------------------------------------------------------
	// Note handlers
	// -----------------------------------------------------------------------------------------------------------------
	router.GET("/api/v1/notes", hs.GetAllNotesHandler)
	router.POST("/api/v1/notes", hs.AddNoteHandler)
	router.GET("/api/v1/notes/:id", hs.GetNoteHandler)
	router.DELETE("/api/v1/notes/:id", hs.DeleteNoteHandler)
	router.PATCH("/api/v1/notes/:id", hs.UpdateNoteHandler)

	// -----------------------------------------------------------------------------------------------------------------
	// Collection handlers
	// -----------------------------------------------------------------------------------------------------------------
	router.GET("/api/v1/collections", hs.GetAllCollectionsHandler)
	router.POST("/api/v1/collections", hs.AddCollectionHandler)
	router.DELETE("/api/v1/collections/:collID/members/:noteID", hs.DeleteNoteFromCollectionHandler)
	router.GET("/api/v1/collections/:collID/notes", hs.GetAllNotesInCollectionHandler)
	// router.POST("/api/v1/collections/:collID/notes", hs.AddNoteToSectionHandler)
	router.GET("/api/v1/collections/users/:collID", hs.GetAllUsersCollectionsHandler)
	router.GET("/api/v1/collections/:collID", hs.GetCollectionHandler)
	router.DELETE("/api/v1/collections/:collID", hs.DeleteCollectionHandler)
	router.PATCH("/api/v1/collections/:collID", hs.UpdateCollectionHandler)

	// -----------------------------------------------------------------------------------------------------------------
	// Section handlers
	// -----------------------------------------------------------------------------------------------------------------
	// router.GET("/api/v1/sections", hs.GetAllSectionsHandler)
	// router.POST("/api/v1/sections", hs.AddSectionHandler)
	// router.DELETE("/api/v1/sections/:secID/members/:noteID", hs.DeleteNoteFromSectionHandler)
	// router.GET("/api/v1/sections/:secID/notes", hs.GetAllNotesInSectionHandler)
	// router.POST("/api/v1/sections/:secID/notes", hs.AddNoteToSectionHandler)
	// router.GET("/api/v1/sections/:secID", hs.GetSectionHandler)
	// router.DELETE("/api/v1/sections/:secID", hs.DeleteSectionHandler)
	// router.PATCH("/api/v1/sections/:secID", hs.UpdateSectionHandler)

	// -----------------------------------------------------------------------------------------------------------------
	// Documentation
	// -----------------------------------------------------------------------------------------------------------------
	router.StaticFile("/api/v1/apifile", "/servers/notebook-app/docs/swagger.json")
	router.StaticFile("/api/v1/docs", "/servers/notebook-app/docs/index.html")

	// -----------------------------------------------------------------------------------------------------------------
	// Static path
	// -----------------------------------------------------------------------------------------------------------------
	router.StaticFS("/api/v1/static", http.Dir("/servers/notebook-app/static/"))
}
