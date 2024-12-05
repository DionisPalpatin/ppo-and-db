package handlersv2

import (
	"net/http"

	bl "github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/business_logic"
	"github.com/gorilla/mux"
)

type Handlers struct {
	svcs *bl.IServices
	reps *bl.IRepositories
}

func (app *App) InitInterfaces(svcs *bl.IServices, reps *bl.IRepositories) {
	h.svcs = svcs
	h.reps = reps
}

func InitRouter(router **mux.Router) {
	*router = mux.NewRouter()

	(*router).HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	(*router).HandleFunc("/api/login", LoginHandler).Methods("POST")
	(*router).HandleFunc("/api/register", RegisterHandler).Methods("POST")

	(*router).HandleFunc("/api/notes", AddNoteHandler).Methods("POST")

	(*router).HandleFunc("/api/notes/id/{noteId}", DeleteNoteHandler).Methods("DELETE")
	(*router).HandleFunc("/api/notes/name/{noteId}", DeleteNoteHandler).Methods("DELETE")

	(*router).HandleFunc("/api/notes/id/{noteInput}", FindNoteHandler).Methods("GET")
	(*router).HandleFunc("/api/notes/name/{noteInput}", FindNoteHandler).Methods("GET")

	(*router).HandleFunc("/api/notes", GetAllNotesHandler).Methods("GET")

	(*router).HandleFunc("/api/collections/id/{collectionId}/notes", ShowCollectionNotesHandler).Methods("GET")
	(*router).HandleFunc("/api/collections/name/{collectionId}/notes", ShowCollectionNotesHandler).Methods("GET")

	(*router).HandleFunc("/api/notes", GetAllNotesInTeamSectionHandler).Methods("GET")
	/* TODO */ /* fetchNotes */ // conf.(*router).HandleFunc("/api/notes?s=allteams", conf.GetAllNotesHandler).Methods("GET")

	/* check renderAddCollectionForm */
	(*router).HandleFunc("/api/collections", AddCollectionHandler).Methods("POST")

	/* check renderDeleteCollectionForm */
	(*router).HandleFunc("/api/collections/id/{collectionId}", DeleteCollectionHandler).Methods("DELETE")

	/* check fetchCollections */
	(*router).HandleFunc("/api/collections", GetAllCollectionsHandler).Methods("GET")
	/* check fetchCollections */
	(*router).HandleFunc("/api/collections", GetAllCollectionsHandler).Methods("GET")

	/* check renderAddNoteToCollectionForm */
	(*router).HandleFunc("/api/collections/id/{collectionId}/notes/id/{noteId}", AddNoteToCollectionHandler).Methods("POST")
	/* check renderAddNoteToCollectionForm */
	(*router).HandleFunc("/api/collections/name/{collectionId}/notes/id/{noteId}", AddNoteToCollectionHandler).Methods("POST")
	/* check renderAddNoteToCollectionForm */
	(*router).HandleFunc("/api/collections/id/{collectionId}/notes/name/{noteId}", AddNoteToCollectionHandler).Methods("POST")
	/* check renderAddNoteToCollectionForm */
	(*router).HandleFunc("/api/collections/name/{collectionId}/notes/name/{noteId}", AddNoteToCollectionHandler).Methods("POST")

	/* check renderDeleteNoteFromCollectionForm */
	(*router).HandleFunc("/api/collections/{sectionId}/notes/id/{noteId}", DeleteNoteFromCollectionHandler).Methods("DELETE")
	/* check renderDeleteNoteFromCollectionForm */
	(*router).HandleFunc("/api/collections/{sectionId}/notes/name/{noteId}", DeleteNoteFromSectionHandler).Methods("DELETE")
	/* check renderDeleteNoteFromCollectionForm */
	(*router).HandleFunc("/api/collections/{sectionId}/notes/id/{noteId}", DeleteNoteFromCollectionHandler).Methods("DELETE")
	/* check renderDeleteNoteFromCollectionForm */
	(*router).HandleFunc("/api/collections/{sectionId}/notes/name/{noteId}", DeleteNoteFromCollectionHandler).Methods("DELETE")

	/* check renderDeleteUserForm */
	(*router).HandleFunc("/api/users/id/{userId}", DeleteUserHandler).Methods("DELETE")
	/* check renderDeleteUserForm */
	(*router).HandleFunc("/api/users/name/{userId}", DeleteUserHandler).Methods("DELETE")

	/* check renderUpdateUserFioForm */
	(*router).HandleFunc("/api/users/id/{userId}/fio", UpdateUserFioHandler).Methods("PATCH")
	/* check renderUpdateUserFioForm */
	(*router).HandleFunc("/api/users/name/{userId}/fio", UpdateUserFioHandler).Methods("PATCH")

	/* check renderUpdateUserRoleForm */
	(*router).HandleFunc("/api/users/id/{userId}/role", UpdateUserRoleHandler).Methods("PATCH")
	/* check renderUpdateUserRoleForm */
	(*router).HandleFunc("/api/users/name/{userId}/role", UpdateUserRoleHandler).Methods("PATCH")

	/* check renderFindUserForm */
	(*router).HandleFunc("/api/users/id/{userId}", FindUserHandler).Methods("GET")
	/* check renderFindUserForm */
	(*router).HandleFunc("/api/users/name/{userId}", FindUserHandler).Methods("GET")

	/* check fetchUsers */
	(*router).HandleFunc("/api/users", GetAllUsersHandler).Methods("GET")

	/* check renderAddTeamForm */
	(*router).HandleFunc("/api/teams", AddTeamHandler).Methods("POST")

	/* check renderDeleteTeamForm */
	(*router).HandleFunc("/api/teams/id/{teamId}", DeleteTeamHandler).Methods("DELETE")
	/* check renderDeleteTeamForm */
	(*router).HandleFunc("/api/teams/name/{teamId}", DeleteTeamHandler).Methods("DELETE")

	/* check renderFindTeamForm */
	(*router).HandleFunc("/api/teams/id/{teamIdOrName}", FindTeamHandler).Methods("GET")
	/* check renderFindTeamForm */
	(*router).HandleFunc("/api/teams/id/{teamIdOrName}", FindTeamHandler).Methods("GET")

	/* check renderShowTeamMembersForm */
	(*router).HandleFunc("/api/teams/id/{teamId}/members", ShowTeamMembersHandler).Methods("GET")
	/* check renderShowTeamMembersForm */
	(*router).HandleFunc("/api/teams/name/{teamId}/members", ShowTeamMembersHandler).Methods("GET")

	/* check fetchTeams */
	(*router).HandleFunc("/api/teams", GetAllTeamsHandler).Methods("GET")

	/* check renderAddUserToTeamForm */
	(*router).HandleFunc("/api/teams/id/{teamId}/members/id/{userId}", AddUserToTeamHandler).Methods("POST")
	/* check renderAddUserToTeamForm */
	(*router).HandleFunc("/api/teams/name/{teamId}/members/id/{userId}", AddUserToTeamHandler).Methods("POST")
	/* check renderAddUserToTeamForm */
	(*router).HandleFunc("/api/teams/id/{teamId}/members/name/{userId}", AddUserToTeamHandler).Methods("POST")
	/* check renderAddUserToTeamForm */
	(*router).HandleFunc("/api/teams/name/{teamId}/members/name/{userId}", AddUserToTeamHandler).Methods("POST")

	/* check renderDeleteUserFromTeamForm */
	(*router).HandleFunc("/api/teams/id/{teamId}/members/id/{userId}", DeleteUserFromTeamHandler).Methods("DELETE")
	/* check renderDeleteUserFromTeamForm */
	(*router).HandleFunc("/api/teams/id/{teamId}/members/name/{userId}", DeleteUserFromTeamHandler).Methods("DELETE")
	/* check renderDeleteUserFromTeamForm */
	(*router).HandleFunc("/api/teams/name/{teamId}/members/id/{userId}", DeleteUserFromTeamHandler).Methods("DELETE")
	/* check renderDeleteUserFromTeamForm */
	(*router).HandleFunc("/api/teams/name/{teamId}/members/name/{userId}", DeleteUserFromTeamHandler).Methods("DELETE")

	/* check renderAddSectionForm */
	(*router).HandleFunc("/api/sections/id/{teamName}", AddSectionHandler).Methods("POST")
	/* check renderAddSectionForm */
	(*router).HandleFunc("/api/sections/name/{teamName}", AddSectionHandler).Methods("POST")

	/* check renderDeleteSectionForm */
	(*router).HandleFunc("/api/sections/id/{teamName}", DeleteSectionHandler).Methods("DELETE")
	/* check renderDeleteSectionForm */
	(*router).HandleFunc("/api/sections/name/{teamName}", DeleteSectionHandler).Methods("DELETE")

	/* check fetchSections */
	(*router).HandleFunc("/api/sections", GetAllSectionsHandler).Methods("GET")

	/* check renderAddNoteToSectionForm */
	(*router).HandleFunc("/api/sections/{sectionId}/notes/id/{noteId}", AddNoteToSectionHandler).Methods("POST")
	/* check renderAddNoteToSectionForm */
	(*router).HandleFunc("/api/sections/{sectionId}/notes/name/{noteId}", AddNoteToSectionHandler).Methods("POST")

	/* check renderDeleteNoteFromSectionForm */
	(*router).HandleFunc("/api/sections/{sectionId}/notes/id/{noteId}", DeleteNoteFromSectionHandler).Methods("DELETE")
	/* check renderDeleteNoteFromSectionForm */
	(*router).HandleFunc("/api/sections/{sectionId}/notes/name/{noteId}", DeleteNoteFromSectionHandler).Methods("DELETE")

	(*router).HandleFunc("/api/stat", GetFullStatHendler).Methods("GET")

}
