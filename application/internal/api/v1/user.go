package handlersv1

import (
	"net/http"
)

func (app *App) GetUserHandler(w http.ResponseWriter, r *http.Request) {}

func (app *App) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {}

func (app *App) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {}

func (app *App) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {}

func (app *App) LoginUserHandler(w http.ResponseWriter, r *http.Request) {}

func (app *App) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {}
