package handlersv2

import (
	"net/http"
)

func (h *Handlers) GetNoteHandler(w http.ResponseWriter, r *http.Request) {}

func (h *Handlers) GetAllNotesHandler(w http.ResponseWriter, r *http.Request) {}

func (h *Handlers) AddNoteHandler(w http.ResponseWriter, r *http.Request) {}

func (h *Handlers) DeleteNoteHandler(w http.ResponseWriter, r *http.Request) {}

func (h *Handlers) UpdateNoteHandler(w http.ResponseWriter, r *http.Request) {}

func (h *Handlers) AddNoteToCollectionHandler(w http.ResponseWriter, r *http.Request) {}

func (h *Handlers) AddNoteToSectionHandler(w http.ResponseWriter, r *http.Request) {}
