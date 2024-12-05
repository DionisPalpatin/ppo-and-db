package handlersv1

import (
	"net/http"
)

func (h *Handlers) GetCollectionHandler(w http.ResponseWriter, r *http.Request) {}

func (h *Handlers) GetAllCollectionsHandler(w http.ResponseWriter, r *http.Request) {}

func (h *Handlers) AddCollectionHandler(w http.ResponseWriter, r *http.Request) {}

func (h *Handlers) DeleteCollectionHandler(w http.ResponseWriter, r *http.Request) {}

func (h *Handlers) UpdateCollectionHandler(w http.ResponseWriter, r *http.Request) {}

func (h *Handlers) GetAllNotesInCollectionHandler(w http.ResponseWriter, r *http.Request) {}

// duplicate to func in note api
// func (h* Handlers) AddNoteToCollectionHandler(w http.ResponseWriter, r *http.Request) {}
