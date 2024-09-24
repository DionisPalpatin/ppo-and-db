package app

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/config"
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/UI/TechUI"
	appconfigs "github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/app/AppConfigs"
	bl "github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/business_logic"
	daclickhouse "github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/data_access/clickhouse"
	dapostgres "github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/data_access/postgres"
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/database"
	mylogger "github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/logger"
)

func initDBConnection(config *config.Configs) error {
	logger := config.LogConfigs.Logger
	logger.WriteLog("Init DB connection", slog.LevelInfo, nil)

	err := database.Connect(config.LogConfigs.Logger, config.DBConfigs)
	if err != nil {
		logger.WriteLog("Unable to init db connection", slog.LevelError, nil)
	}

	return err
}

func initInterfaces(conf *appconfigs.AppConfigs) {
	dbconf := conf.Configs.DBConfigs
	logger := conf.Configs.LogConfigs

	if dbconf.DriverName == "postgres" {
		conf.IRepos = &bl.IRepositories{
			IUsrRepo:  &dapostgres.UserRepository{DbConfigs: dbconf, MyLogger: logger.Logger},
			ISecRepo:  &dapostgres.SectionRepository{DbConfigs: dbconf, MyLogger: logger.Logger},
			INoteRepo: &dapostgres.NoteRepository{DbConfigs: dbconf, MyLogger: logger.Logger},
			IColRepo:  &dapostgres.CollectionRepository{DbConfigs: dbconf, MyLogger: logger.Logger},
			ITeamRepo: &dapostgres.TeamRepository{DbConfigs: dbconf, MyLogger: logger.Logger},
		}
	} else if dbconf.DriverName == "clickhouse" {
		conf.IRepos = &bl.IRepositories{
			IUsrRepo:  &daclickhouse.UserRepository{DbConfigs: dbconf, MyLogger: logger.Logger},
			ISecRepo:  &daclickhouse.SectionRepository{DbConfigs: dbconf, MyLogger: logger.Logger},
			INoteRepo: &daclickhouse.NoteRepository{DbConfigs: dbconf, MyLogger: logger.Logger},
			IColRepo:  &daclickhouse.CollectionRepository{DbConfigs: dbconf, MyLogger: logger.Logger},
			ITeamRepo: &daclickhouse.TeamRepository{DbConfigs: dbconf, MyLogger: logger.Logger},
		}
	}

	conf.IServices = &bl.IServices{
		IUsrSvc:   &bl.UserService{},
		ISecSvc:   &bl.SectionService{},
		INoteSvc:  &bl.NoteService{},
		IColSvc:   &bl.CollectionService{},
		ITeamSvc:  &bl.TeamService{},
		IOAuthSvc: &bl.OAuthService{},
	}
}

func initRouter(conf *appconfigs.AppConfigs) {
	conf.Router = mux.NewRouter()

	conf.Router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	conf.Router.HandleFunc("/api/login", conf.LoginHandler).Methods("POST")
	conf.Router.HandleFunc("/api/register", conf.RegisterHandler).Methods("POST")

	/* check AddNoteForm */
	conf.Router.HandleFunc("/api/notes", conf.AddNoteHandler).Methods("POST")

	/* check renderDeleteNoteForm */
	conf.Router.HandleFunc("/api/notes/id/{noteId}", conf.DeleteNoteHandler).Methods("DELETE")
	/* check renderDeleteNoteForm */
	conf.Router.HandleFunc("/api/notes/name/{noteId}", conf.DeleteNoteHandler).Methods("DELETE")

	/* check renderFindNoteForm  */
	conf.Router.HandleFunc("/api/notes/id/{noteInput}", conf.FindNoteHandler).Methods("GET")
	/* check renderFindNoteForm */
	conf.Router.HandleFunc("/api/notes/name/{noteInput}", conf.FindNoteHandler).Methods("GET")

	/* check fetchNotes */
	conf.Router.HandleFunc("/api/notes", conf.GetAllNotesHandler).Methods("GET")
	/* check fetchNotes */
	conf.Router.HandleFunc("/api/notes", conf.GetAllNotesHandler).Methods("GET")

	/* check renderShowCollectionNotesForm */
	conf.Router.HandleFunc("/api/collections/id/{collectionId}/notes", conf.ShowCollectionNotesHandler).Methods("GET")
	/* check renderShowCollectionNotesForm */
	conf.Router.HandleFunc("/api/collections/name/{collectionId}/notes", conf.ShowCollectionNotesHandler).Methods("GET")

	/* check fetchNotes */
	conf.Router.HandleFunc("/api/notes", conf.GetAllNotesInTeamSectionHandler).Methods("GET")
	/* TODO */ /* fetchNotes */ // conf.Router.HandleFunc("/api/notes?s=allteams", conf.GetAllNotesHandler).Methods("GET")

	/* check renderAddCollectionForm */
	conf.Router.HandleFunc("/api/collections", conf.AddCollectionHandler).Methods("POST")

	/* check renderDeleteCollectionForm */
	conf.Router.HandleFunc("/api/collections/id/{collectionId}", conf.DeleteCollectionHandler).Methods("DELETE")

	/* check fetchCollections */
	conf.Router.HandleFunc("/api/collections", conf.GetAllCollectionsHandler).Methods("GET")
	/* check fetchCollections */
	conf.Router.HandleFunc("/api/collections", conf.GetAllCollectionsHandler).Methods("GET")

	/* check renderAddNoteToCollectionForm */
	conf.Router.HandleFunc("/api/collections/id/{collectionId}/notes/id/{noteId}", conf.AddNoteToCollectionHandler).Methods("POST")
	/* check renderAddNoteToCollectionForm */
	conf.Router.HandleFunc("/api/collections/name/{collectionId}/notes/id/{noteId}", conf.AddNoteToCollectionHandler).Methods("POST")
	/* check renderAddNoteToCollectionForm */
	conf.Router.HandleFunc("/api/collections/id/{collectionId}/notes/name/{noteId}", conf.AddNoteToCollectionHandler).Methods("POST")
	/* check renderAddNoteToCollectionForm */
	conf.Router.HandleFunc("/api/collections/name/{collectionId}/notes/name/{noteId}", conf.AddNoteToCollectionHandler).Methods("POST")

	/* check renderDeleteNoteFromCollectionForm */
	conf.Router.HandleFunc("/api/collections/{sectionId}/notes/id/{noteId}", conf.DeleteNoteFromCollectionHandler).Methods("DELETE")
	/* check renderDeleteNoteFromCollectionForm */
	conf.Router.HandleFunc("/api/collections/{sectionId}/notes/name/{noteId}", conf.DeleteNoteFromSectionHandler).Methods("DELETE")
	/* check renderDeleteNoteFromCollectionForm */
	conf.Router.HandleFunc("/api/collections/{sectionId}/notes/id/{noteId}", conf.DeleteNoteFromCollectionHandler).Methods("DELETE")
	/* check renderDeleteNoteFromCollectionForm */
	conf.Router.HandleFunc("/api/collections/{sectionId}/notes/name/{noteId}", conf.DeleteNoteFromCollectionHandler).Methods("DELETE")

	/* check renderDeleteUserForm */
	conf.Router.HandleFunc("/api/users/id/{userId}", conf.DeleteUserHandler).Methods("DELETE")
	/* check renderDeleteUserForm */
	conf.Router.HandleFunc("/api/users/name/{userId}", conf.DeleteUserHandler).Methods("DELETE")

	/* check renderUpdateUserFioForm */
	conf.Router.HandleFunc("/api/users/id/{userId}/fio", conf.UpdateUserFioHandler).Methods("PATCH")
	/* check renderUpdateUserFioForm */
	conf.Router.HandleFunc("/api/users/name/{userId}/fio", conf.UpdateUserFioHandler).Methods("PATCH")

	/* check renderUpdateUserRoleForm */
	conf.Router.HandleFunc("/api/users/id/{userId}/role", conf.UpdateUserRoleHandler).Methods("PATCH")
	/* check renderUpdateUserRoleForm */
	conf.Router.HandleFunc("/api/users/name/{userId}/role", conf.UpdateUserRoleHandler).Methods("PATCH")

	/* check renderFindUserForm */
	conf.Router.HandleFunc("/api/users/id/{userId}", conf.FindUserHandler).Methods("GET")
	/* check renderFindUserForm */
	conf.Router.HandleFunc("/api/users/name/{userId}", conf.FindUserHandler).Methods("GET")

	/* check fetchUsers */
	conf.Router.HandleFunc("/api/users", conf.GetAllUsersHandler).Methods("GET")

	/* check renderAddTeamForm */
	conf.Router.HandleFunc("/api/teams", conf.AddTeamHandler).Methods("POST")

	/* check renderDeleteTeamForm */
	conf.Router.HandleFunc("/api/teams/id/{teamId}", conf.DeleteTeamHandler).Methods("DELETE")
	/* check renderDeleteTeamForm */
	conf.Router.HandleFunc("/api/teams/name/{teamId}", conf.DeleteTeamHandler).Methods("DELETE")

	/* check renderFindTeamForm */
	conf.Router.HandleFunc("/api/teams/id/{teamIdOrName}", conf.FindTeamHandler).Methods("GET")
	/* check renderFindTeamForm */
	conf.Router.HandleFunc("/api/teams/id/{teamIdOrName}", conf.FindTeamHandler).Methods("GET")

	/* check renderShowTeamMembersForm */
	conf.Router.HandleFunc("/api/teams/id/{teamId}/members", conf.ShowTeamMembersHandler).Methods("GET")
	/* check renderShowTeamMembersForm */
	conf.Router.HandleFunc("/api/teams/name/{teamId}/members", conf.ShowTeamMembersHandler).Methods("GET")

	/* check fetchTeams */
	conf.Router.HandleFunc("/api/teams", conf.GetAllTeamsHandler).Methods("GET")

	/* check renderAddUserToTeamForm */
	conf.Router.HandleFunc("/api/teams/id/{teamId}/members/id/{userId}", conf.AddUserToTeamHandler).Methods("POST")
	/* check renderAddUserToTeamForm */
	conf.Router.HandleFunc("/api/teams/name/{teamId}/members/id/{userId}", conf.AddUserToTeamHandler).Methods("POST")
	/* check renderAddUserToTeamForm */
	conf.Router.HandleFunc("/api/teams/id/{teamId}/members/name/{userId}", conf.AddUserToTeamHandler).Methods("POST")
	/* check renderAddUserToTeamForm */
	conf.Router.HandleFunc("/api/teams/name/{teamId}/members/name/{userId}", conf.AddUserToTeamHandler).Methods("POST")

	/* check renderDeleteUserFromTeamForm */
	conf.Router.HandleFunc("/api/teams/id/{teamId}/members/id/{userId}", conf.DeleteUserFromTeamHandler).Methods("DELETE")
	/* check renderDeleteUserFromTeamForm */
	conf.Router.HandleFunc("/api/teams/id/{teamId}/members/name/{userId}", conf.DeleteUserFromTeamHandler).Methods("DELETE")
	/* check renderDeleteUserFromTeamForm */
	conf.Router.HandleFunc("/api/teams/name/{teamId}/members/id/{userId}", conf.DeleteUserFromTeamHandler).Methods("DELETE")
	/* check renderDeleteUserFromTeamForm */
	conf.Router.HandleFunc("/api/teams/name/{teamId}/members/name/{userId}", conf.DeleteUserFromTeamHandler).Methods("DELETE")

	/* check renderAddSectionForm */
	conf.Router.HandleFunc("/api/sections/id/{teamName}", conf.AddSectionHandler).Methods("POST")
	/* check renderAddSectionForm */
	conf.Router.HandleFunc("/api/sections/name/{teamName}", conf.AddSectionHandler).Methods("POST")

	/* check renderDeleteSectionForm */
	conf.Router.HandleFunc("/api/sections/id/{teamName}", conf.DeleteSectionHandler).Methods("DELETE")
	/* check renderDeleteSectionForm */
	conf.Router.HandleFunc("/api/sections/name/{teamName}", conf.DeleteSectionHandler).Methods("DELETE")

	/* check fetchSections */
	conf.Router.HandleFunc("/api/sections", conf.GetAllSectionsHandler).Methods("GET")

	/* check renderAddNoteToSectionForm */
	conf.Router.HandleFunc("/api/sections/{sectionId}/notes/id/{noteId}", conf.AddNoteToSectionHandler).Methods("POST")
	/* check renderAddNoteToSectionForm */
	conf.Router.HandleFunc("/api/sections/{sectionId}/notes/name/{noteId}", conf.AddNoteToSectionHandler).Methods("POST")

	/* check renderDeleteNoteFromSectionForm */
	conf.Router.HandleFunc("/api/sections/{sectionId}/notes/id/{noteId}", conf.DeleteNoteFromSectionHandler).Methods("DELETE")
	/* check renderDeleteNoteFromSectionForm */
	conf.Router.HandleFunc("/api/sections/{sectionId}/notes/name/{noteId}", conf.DeleteNoteFromSectionHandler).Methods("DELETE")
}

func RunBackend() error {
	configFile := "./config/config.yaml"

	appConfigs := new(appconfigs.AppConfigs)
	var err error

	appConfigs.Configs, err = config.ReadConfig(configFile)
	if err != nil {
		return err
	}
	configs := appConfigs.Configs
	configs.LogConfigs.Logger = new(mylogger.MyLogger)

	logFile := configs.LogConfigs.LogFile
	logLevel := configs.LogConfigs.LogLevel
	err = configs.LogConfigs.Logger.InitLogger(logFile, logLevel)
	if err != nil {
		return err
	}

	err = initDBConnection(configs)
	if err != nil {
		return err
	}

	// initRepositories(appConfigs)
	initInterfaces(appConfigs)

	if appConfigs.Configs.Mode == "tech" {
		for {
			user := TechUI.AuthorizationMenu(configs, appConfigs.IRepos, appConfigs.IServices)
			if user == nil {
				break
			} else if user.Role == bl.Reader {
				TechUI.ReaderMenu(user, configs, appConfigs.IRepos, appConfigs.IServices)
			} else if user.Role == bl.Author {
				TechUI.AuthorMenu(user, configs, appConfigs.IRepos, appConfigs.IServices)
			} else if user.Role == bl.Admin {
				TechUI.AdminMenu(user, configs, appConfigs.IRepos, appConfigs.IServices)
			}
		}
	} else {
		initRouter(appConfigs)
		port := fmt.Sprintf(":%d", appConfigs.Configs.ServerPort)

		fs := http.FileServer(http.Dir("./static/"))
		appConfigs.Router.PathPrefix("/").Handler(fs)

		http.ListenAndServe(port, appConfigs.Router)
		appConfigs.Configs.LogConfigs.Logger.WriteLog("Server is running on port 8080", slog.LevelInfo, nil)
	}

	return nil
}
