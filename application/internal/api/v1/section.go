package handlersv1

import (
	"net/http"
)

func (app *App) GetSectionHandler(w http.ResponseWriter, r *http.Request) {}

func (app *App) GetAllSectionsHandler(w http.ResponseWriter, r *http.Request) {}

func (app *App) AddSectionHandler(w http.ResponseWriter, r *http.Request) {}

func (app *App) DeleteSectionHandler(w http.ResponseWriter, r *http.Request) {}

func (app *App) UpdateSectionHandler(w http.ResponseWriter, r *http.Request) {}

func (app *App) GetAllNotesInSectionHandler(w http.ResponseWriter, r *http.Request) {}

// duplicate to func in note api
// func (h* Handlers) AddNoteToSectionHandler(w http.ResponseWriter, r *http.Request) {}
