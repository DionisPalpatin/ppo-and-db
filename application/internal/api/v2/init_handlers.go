package handlersv2

import (
	"net/http"

	"github.com/gorilla/mux"
)

func InitRouter(app *App) {
	app.Router = mux.NewRouter()
	router := *app.Router

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	// -----------------------------------------------------------------------------------------------------------------
	// Authorisation handlers
	// -----------------------------------------------------------------------------------------------------------------
	router.HandleFunc("/login", app.LoginHandler).Methods("POST")
	router.HandleFunc("/register", app.RegisterHandler).Methods("POST")

	// -----------------------------------------------------------------------------------------------------------------
	// User handlers
	// -----------------------------------------------------------------------------------------------------------------
	router.HandleFunc("/users", app.GetAllUsersHandler).Methods("GET")
	router.HandleFunc("/users/{id}", app.GetUserHandler).Methods("GET")
	router.HandleFunc("/users/{id}", app.DeleteUserHandler).Methods("DELETE")
	router.HandleFunc("/users/{id}", app.UpdateUserHandler).Methods("PATCH")

	// -----------------------------------------------------------------------------------------------------------------
	// Team handlers
	// -----------------------------------------------------------------------------------------------------------------
	router.HandleFunc("/teams", app.GetAllTeamsHandler).Methods("GET")
	router.HandleFunc("/teams", app.AddTeamHandler).Methods("POST")
	router.HandleFunc("/teams/{id}", app.GetTeamHandler).Methods("GET")
	router.HandleFunc("/teams/{id}", app.DeleteTeamHandler).Methods("DELETE")
	router.HandleFunc("/teams/{id}", app.UpdateTeamHandler).Methods("PATCH")
	router.HandleFunc("/teams/{id}/members", app.GetTeamMembersHandler).Methods("GET")
	router.HandleFunc("/teams/{id}/members", app.AddUserToTeamHandler).Methods("POST")
	router.HandleFunc("/teams/{teamID}/members/{userID}", app.DeleteUserFromTeamHandler).Methods("DELETE")

	// -----------------------------------------------------------------------------------------------------------------
	// Note handlers
	// -----------------------------------------------------------------------------------------------------------------
	router.HandleFunc("/notes", app.GetAllNotesHandler).Methods("GET")
	router.HandleFunc("/notes", app.AddNoteHandler).Methods("POST")
	router.HandleFunc("/notes/{id}", app.GetNoteHandler).Methods("GET")
	router.HandleFunc("/notes/{id}", app.DeleteNoteHandler).Methods("DELETE")
	router.HandleFunc("/notes/{id}", app.UpdateNoteHandler).Methods("PATCH")

	// -----------------------------------------------------------------------------------------------------------------
	// Collection handlers
	// -----------------------------------------------------------------------------------------------------------------
	router.HandleFunc("/collections", app.GetAllCollectionsHandler).Methods("GET")
	router.HandleFunc("/collections", app.AddCollectionHandler).Methods("POST")
	router.HandleFunc("/collections/users/{id}", app.GetAllUsersCollectionsHandler).Methods("GET")
	router.HandleFunc("/collections/{id}", app.GetCollectionHandler).Methods("GET")
	router.HandleFunc("/collections/{id}", app.DeleteCollectionHandler).Methods("DELETE")
	router.HandleFunc("/collections/{id}", app.UpdateCollectionHandler).Methods("PATCH")
	router.HandleFunc("/collections/{id}/notes", app.GetAllNotesInCollectionHandler).Methods("GET")
	router.HandleFunc("/collections/{id}/notes", app.AddNoteToSectionHandler).Methods("POST")
	router.HandleFunc("/collections/{collID}/members/{noteID}", app.DeleteNoteFromCollectionHandler).Methods("DELETE")

	// -----------------------------------------------------------------------------------------------------------------
	// Section handlers
	// -----------------------------------------------------------------------------------------------------------------
	router.HandleFunc("/sections", app.GetAllSectionsHandler).Methods("GET")
	router.HandleFunc("/sections", app.AddSectionHandler).Methods("POST")
	router.HandleFunc("/sections/{id}", app.GetSectionHandler).Methods("GET")
	router.HandleFunc("/sections/{id}", app.DeleteSectionHandler).Methods("DELETE")
	router.HandleFunc("/sections/{id}", app.UpdateSectionHandler).Methods("PATCH")
	router.HandleFunc("/sections/{id}/notes", app.GetAllNotesInSectionHandler).Methods("GET")
	router.HandleFunc("/sections/{id}/notes", app.AddNoteToSectionHandler).Methods("POST")
	router.HandleFunc("/sections/{secID}/members/{noteID}", app.DeleteNoteFromSectionHandler).Methods("DELETE")
}
