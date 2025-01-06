package handlersv1

import (
	"net/http"
)

func (app *App) GetNoteHandler(w http.ResponseWriter, r *http.Request) {}

func (app *App) GetAllNotesHandler(w http.ResponseWriter, r *http.Request) {}

func (app *App) AddNoteHandler(w http.ResponseWriter, r *http.Request) {}

func (app *App) DeleteNoteHandler(w http.ResponseWriter, r *http.Request) {}

func (app *App) UpdateNoteHandler(w http.ResponseWriter, r *http.Request) {}

func (app *App) AddNoteToCollectionHandler(w http.ResponseWriter, r *http.Request) {}

func (app *App) AddNoteToSectionHandler(w http.ResponseWriter, r *http.Request) {}
