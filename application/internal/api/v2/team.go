package handlersv2

import (
	"net/http"
)

func (h *Handlers) GetTeamHandler(w http.ResponseWriter, r *http.Request) {}

func (h *Handlers) GetAllTeamsHandler(w http.ResponseWriter, r *http.Request) {}

func (h *Handlers) AddTeamHandler(w http.ResponseWriter, r *http.Request) {}

func (h *Handlers) DeleteTeamHandler(w http.ResponseWriter, r *http.Request) {}

func (h *Handlers) UpdateTeamHandler(w http.ResponseWriter, r *http.Request) {}

func (h *Handlers) AddUserToTeamHandler(w http.ResponseWriter, r *http.Request) {}

func (h *Handlers) GetTeamMembersHandler(w http.ResponseWriter, r *http.Request) {}
