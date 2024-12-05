package handlersv1

import (
	"net/http"
)

func (h *Handlers) GetUserHandler(w http.ResponseWriter, r *http.Request) {}

func (h *Handlers) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {}

func (h *Handlers) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {}

func (h *Handlers) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {}

func (h *Handlers) LoginUserHandler(w http.ResponseWriter, r *http.Request) {}

func (h *Handlers) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {}
