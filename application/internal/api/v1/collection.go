package handlersv1

import (
	"net/http"
)

func (app *App) GetCollectionHandler(w http.ResponseWriter, r *http.Request) {}

func (app *App) GetAllCollectionsHandler(w http.ResponseWriter, r *http.Request) {}

func (app *App) AddCollectionHandler(w http.ResponseWriter, r *http.Request) {}

func (app *App) DeleteCollectionHandler(w http.ResponseWriter, r *http.Request) {}

func (app *App) UpdateCollectionHandler(w http.ResponseWriter, r *http.Request) {}

func (app *App) GetAllNotesInCollectionHandler(w http.ResponseWriter, r *http.Request) {}

// duplicate to func in note api
// func (h* Handlers) AddNoteToCollectionHandler(w http.ResponseWriter, r *http.Request) {}
