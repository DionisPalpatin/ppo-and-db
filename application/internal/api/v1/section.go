package handlersv1

import (
	"net/http"
)

func (h *Handlers) GetSectionHandler(w http.ResponseWriter, r *http.Request) {}

func (h *Handlers) GetAllSectionsHandler(w http.ResponseWriter, r *http.Request) {}

func (h *Handlers) AddSectionHandler(w http.ResponseWriter, r *http.Request) {}

func (h *Handlers) DeleteSectionHandler(w http.ResponseWriter, r *http.Request) {}

func (h *Handlers) UpdateSectionHandler(w http.ResponseWriter, r *http.Request) {}

func (h *Handlers) GetAllNotesInSectionHandler(w http.ResponseWriter, r *http.Request) {}

// duplicate to func in note api
// func (h* Handlers) AddNoteToSectionHandler(w http.ResponseWriter, r *http.Request) {}
