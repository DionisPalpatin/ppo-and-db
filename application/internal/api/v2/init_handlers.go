package handlersv2

import (
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/config"
	//_ "github.com/DionisPalpatin/ppo-and-db/tree/master/application/docs"
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

	adapter := func(handler func(http.ResponseWriter, *http.Request)) gin.HandlerFunc {
		return func(c *gin.Context) {
			handler(c.Writer, c.Request)
		}
	}

	// -----------------------------------------------------------------------------------------------------------------
	// Authorisation handlers
	// -----------------------------------------------------------------------------------------------------------------
	router.POST("/login", adapter(hs.LoginHandler))
	router.POST("/register", adapter(hs.RegisterHandler))

	// -----------------------------------------------------------------------------------------------------------------
	// User handlers
	// -----------------------------------------------------------------------------------------------------------------
	router.GET("/users", adapter(hs.GetAllUsersHandler))
	router.GET("/users/:id", adapter(hs.GetUserHandler))
	router.DELETE("/users/:id", adapter(hs.DeleteUserHandler))
	router.PATCH("/users/:id", adapter(hs.UpdateUserHandler))

	// -----------------------------------------------------------------------------------------------------------------
	// Team handlers
	// -----------------------------------------------------------------------------------------------------------------
	router.GET("/teams", adapter(hs.GetAllTeamsHandler))
	router.POST("/teams", adapter(hs.AddTeamHandler))
	router.DELETE("/teams/:teamID/members/:userID", adapter(hs.DeleteUserFromTeamHandler))
	router.GET("/teams/:teamID/members", adapter(hs.GetTeamMembersHandler))
	router.POST("/teams/:teamID/members", adapter(hs.AddUserToTeamHandler))
	router.GET("/teams/:teamID", adapter(hs.GetTeamHandler))
	router.DELETE("/teams/:teamID", adapter(hs.DeleteTeamHandler))
	router.PATCH("/teams/:teamID", adapter(hs.UpdateTeamHandler))

	// -----------------------------------------------------------------------------------------------------------------
	// Note handlers
	// -----------------------------------------------------------------------------------------------------------------
	router.GET("/notes", adapter(hs.GetAllNotesHandler))
	router.POST("/notes", adapter(hs.AddNoteHandler))
	router.GET("/notes/:id", adapter(hs.GetNoteHandler))
	router.DELETE("/notes/:id", adapter(hs.DeleteNoteHandler))
	router.PATCH("/notes/:id", adapter(hs.UpdateNoteHandler))

	// -----------------------------------------------------------------------------------------------------------------
	// Collection handlers
	// -----------------------------------------------------------------------------------------------------------------
	router.GET("/collections", adapter(hs.GetAllCollectionsHandler))
	router.POST("/collections", adapter(hs.AddCollectionHandler))
	router.DELETE("/collections/:collID/members/:noteID", adapter(hs.DeleteNoteFromCollectionHandler))
	router.GET("/collections/:collID/notes", adapter(hs.GetAllNotesInCollectionHandler))
	router.POST("/collections/:collID/notes", adapter(hs.AddNoteToSectionHandler))
	router.GET("/collections/users/:collID", adapter(hs.GetAllUsersCollectionsHandler))
	router.GET("/collections/:collID", adapter(hs.GetCollectionHandler))
	router.DELETE("/collections/:collID", adapter(hs.DeleteCollectionHandler))
	router.PATCH("/collections/:collID", adapter(hs.UpdateCollectionHandler))

	// -----------------------------------------------------------------------------------------------------------------
	// Section handlers
	// -----------------------------------------------------------------------------------------------------------------
	router.GET("/sections", adapter(hs.GetAllSectionsHandler))
	router.POST("/sections", adapter(hs.AddSectionHandler))
	router.DELETE("/sections/:secID/members/:noteID", adapter(hs.DeleteNoteFromSectionHandler))
	router.GET("/sections/:secID/notes", adapter(hs.GetAllNotesInSectionHandler))
	router.POST("/sections/:secID/notes", adapter(hs.AddNoteToSectionHandler))
	router.GET("/sections/:secID", adapter(hs.GetSectionHandler))
	router.DELETE("/sections/:secID", adapter(hs.DeleteSectionHandler))
	router.PATCH("/sections/:secID", adapter(hs.UpdateSectionHandler))

	// -----------------------------------------------------------------------------------------------------------------
	// Static path
	// -----------------------------------------------------------------------------------------------------------------
	//router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Static("/swaggerui", "./docs/swaggerui")
	router.StaticFS("/static", http.Dir("./static/"))
}
