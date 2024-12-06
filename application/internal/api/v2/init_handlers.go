package handlersv2

import (
	"net/http"

	"github.com/gorilla/mux"
)

func InitRouter(hs *HandlersStruct) {
	hs.Router = mux.NewRouter()
	router := *hs.Router

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	// -----------------------------------------------------------------------------------------------------------------
	// Authorisation handlers
	// -----------------------------------------------------------------------------------------------------------------
	router.HandleFunc("/login", hs.LoginHandler).Methods("POST")
	router.HandleFunc("/register", hs.RegisterHandler).Methods("POST")

	// -----------------------------------------------------------------------------------------------------------------
	// User handlers
	// -----------------------------------------------------------------------------------------------------------------
	router.HandleFunc("/users", hs.GetAllUsersHandler).Methods("GET")
	router.HandleFunc("/users/{id}", hs.GetUserHandler).Methods("GET")
	router.HandleFunc("/users/{id}", hs.DeleteUserHandler).Methods("DELETE")
	router.HandleFunc("/users/{id}", hs.UpdateUserHandler).Methods("PATCH")

	// -----------------------------------------------------------------------------------------------------------------
	// Team handlers
	// -----------------------------------------------------------------------------------------------------------------
	router.HandleFunc("/teams", hs.GetAllTeamsHandler).Methods("GET")
	router.HandleFunc("/teams", hs.AddTeamHandler).Methods("POST")
	router.HandleFunc("/teams/{id}", hs.GetTeamHandler).Methods("GET")
	router.HandleFunc("/teams/{id}", hs.DeleteTeamHandler).Methods("DELETE")
	router.HandleFunc("/teams/{id}", hs.UpdateTeamHandler).Methods("PATCH")
	router.HandleFunc("/teams/{id}/members", hs.GetTeamMembersHandler).Methods("GET")
	router.HandleFunc("/teams/{id}/members", hs.AddUserToTeamHandler).Methods("POST")
	router.HandleFunc("/teams/{teamID}/members/{userID}", hs.DeleteUserFromTeamHandler).Methods("DELETE")

	// -----------------------------------------------------------------------------------------------------------------
	// Note handlers
	// -----------------------------------------------------------------------------------------------------------------
	router.HandleFunc("/notes", hs.GetAllNotesHandler).Methods("GET")
	router.HandleFunc("/notes", hs.AddNoteHandler).Methods("POST")
	router.HandleFunc("/notes/{id}", hs.GetNoteHandler).Methods("GET")
	router.HandleFunc("/notes/{id}", hs.DeleteNoteHandler).Methods("DELETE")
	router.HandleFunc("/notes/{id}", hs.UpdateNoteHandler).Methods("PATCH")

	// -----------------------------------------------------------------------------------------------------------------
	// Collection handlers
	// -----------------------------------------------------------------------------------------------------------------
	router.HandleFunc("/collections", hs.GetAllCollectionsHandler).Methods("GET")
	router.HandleFunc("/collections", hs.AddCollectionHandler).Methods("POST")
	router.HandleFunc("/collections/users/{id}", hs.GetAllUsersCollectionsHandler).Methods("GET")
	router.HandleFunc("/collections/{id}", hs.GetCollectionHandler).Methods("GET")
	router.HandleFunc("/collections/{id}", hs.DeleteCollectionHandler).Methods("DELETE")
	router.HandleFunc("/collections/{id}", hs.UpdateCollectionHandler).Methods("PATCH")
	router.HandleFunc("/collections/{id}/notes", hs.GetAllNotesInCollectionHandler).Methods("GET")
	router.HandleFunc("/collections/{id}/notes", hs.AddNoteToSectionHandler).Methods("POST")
	router.HandleFunc("/collections/{collID}/members/{noteID}", hs.DeleteNoteFromCollectionHandler).Methods("DELETE")

	// -----------------------------------------------------------------------------------------------------------------
	// Section handlers
	// -----------------------------------------------------------------------------------------------------------------
	router.HandleFunc("/sections", hs.GetAllSectionsHandler).Methods("GET")
	router.HandleFunc("/sections", hs.AddSectionHandler).Methods("POST")
	router.HandleFunc("/sections/{id}", hs.GetSectionHandler).Methods("GET")
	router.HandleFunc("/sections/{id}", hs.DeleteSectionHandler).Methods("DELETE")
	router.HandleFunc("/sections/{id}", hs.UpdateSectionHandler).Methods("PATCH")
	router.HandleFunc("/sections/{id}/notes", hs.GetAllNotesInSectionHandler).Methods("GET")
	router.HandleFunc("/sections/{id}/notes", hs.AddNoteToSectionHandler).Methods("POST")
	router.HandleFunc("/sections/{secID}/members/{noteID}", hs.DeleteNoteFromSectionHandler).Methods("DELETE")
}
