package handlersv2

import (
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/config"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	
	bl "github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/business_logic"
	"net/http"

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
	router.POST("/login", hs.LoginHandler)
	router.POST("/register", hs.RegisterHandler)

	// -----------------------------------------------------------------------------------------------------------------
	// User handlers
	// -----------------------------------------------------------------------------------------------------------------
	router.GET("/users", hs.GetAllUsersHandler)
	router.GET("/users/:id", hs.GetUserHandler)
	router.DELETE("/users/:id", hs.DeleteUserHandler)
	router.PATCH("/users/:id", hs.UpdateUserHandler)

	// -----------------------------------------------------------------------------------------------------------------
	// Team handlers
	// -----------------------------------------------------------------------------------------------------------------
	router.GET("/teams", hs.GetAllTeamsHandler)
	router.POST("/teams", hs.AddTeamHandler)
	router.DELETE("/teams/:teamID/members/:userID", hs.DeleteUserFromTeamHandler)
	router.GET("/teams/:teamID/members", hs.GetTeamMembersHandler)
	router.POST("/teams/:teamID/members", hs.AddUserToTeamHandler)
	router.GET("/teams/:teamID", hs.GetTeamHandler)
	router.DELETE("/teams/:teamID", hs.DeleteTeamHandler)
	router.PATCH("/teams/:teamID", hs.UpdateTeamHandler)

	// -----------------------------------------------------------------------------------------------------------------
	// Note handlers
	// -----------------------------------------------------------------------------------------------------------------
	router.GET("/notes", hs.GetAllNotesHandler)
	router.POST("/notes", hs.AddNoteHandler)
	router.GET("/notes/:id", hs.GetNoteHandler)
	router.DELETE("/notes/:id", hs.DeleteNoteHandler)
	router.PATCH("/notes/:id", hs.UpdateNoteHandler)

	// -----------------------------------------------------------------------------------------------------------------
	// Collection handlers
	// -----------------------------------------------------------------------------------------------------------------
	router.GET("/collections", hs.GetAllCollectionsHandler)
	router.POST("/collections", hs.AddCollectionHandler)
	router.DELETE("/collections/:collID/members/:noteID", hs.DeleteNoteFromCollectionHandler)
	router.GET("/collections/:collID/notes", hs.GetAllNotesInCollectionHandler)
	router.POST("/collections/:collID/notes", hs.AddNoteToSectionHandler)
	router.GET("/collections/users/:collID", hs.GetAllUsersCollectionsHandler)
	router.GET("/collections/:collID", hs.GetCollectionHandler)
	router.DELETE("/collections/:collID", hs.DeleteCollectionHandler)
	router.PATCH("/collections/:collID", hs.UpdateCollectionHandler)

	// -----------------------------------------------------------------------------------------------------------------
	// Section handlers
	// -----------------------------------------------------------------------------------------------------------------
	router.GET("/sections", hs.GetAllSectionsHandler)
	router.POST("/sections", hs.AddSectionHandler)
	router.DELETE("/sections/:secID/members/:noteID", hs.DeleteNoteFromSectionHandler)
	router.GET("/sections/:secID/notes", hs.GetAllNotesInSectionHandler)
	router.POST("/sections/:secID/notes", hs.AddNoteToSectionHandler)
	router.GET("/sections/:secID", hs.GetSectionHandler)
	router.DELETE("/sections/:secID", hs.DeleteSectionHandler)
	router.PATCH("/sections/:secID", hs.UpdateSectionHandler)

	// -----------------------------------------------------------------------------------------------------------------
	// Documentation
	// -----------------------------------------------------------------------------------------------------------------
	router.GET("/docs", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/docs/yaml")))
	router.GET("/docs/yaml", func(c *gin.Context) {
		c.File("docs/openapi.yaml")
	})

	// -----------------------------------------------------------------------------------------------------------------
	// Static path
	// -----------------------------------------------------------------------------------------------------------------
	router.StaticFS("/static", http.Dir("./static/"))
}
