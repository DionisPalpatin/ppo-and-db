package handlersv1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (app *App) GetTeamHandler(w http.ResponseWriter, r *http.Request) {}

func (app *App) GetAllTeamsHandler(w http.ResponseWriter, r *http.Request) {}

func (app *App) AddTeamHandler(w http.ResponseWriter, r *http.Request) {}

func (app *App) DeleteTeamHandler(w http.ResponseWriter, r *http.Request) {}

func (app *App) UpdateTeamHandler(w http.ResponseWriter, r *http.Request) {}

func (app *App) AddUserToTeamHandler(w http.ResponseWriter, r *http.Request) {}

func (app *App) GetTeamMembersHandler(w http.ResponseWriter, r *http.Request) {}
