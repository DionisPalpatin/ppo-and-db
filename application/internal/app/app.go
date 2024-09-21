package app

import (
	"log/slog"

	"notebook_app/config"
	"notebook_app/internal/UI/TechUI"
	"notebook_app/internal/app/AppConfigs"
	bl "notebook_app/internal/business_logic"
	da "notebook_app/internal/data_access"
	"notebook_app/internal/database"
	mylogger "notebook_app/internal/logger"
)

func initDBConnection(config *config.Configs) error {
	logger := config.LogConfigs.Logger
	logger.WriteLog("Init DB connection", slog.LevelInfo, nil)

	err := database.Connect(config.DBConfigs)
	if err != nil {
		logger.WriteLog("Unable to init db connection", slog.LevelError, nil)
	}

	return err
}

func initRepositories(config *appconfigs.AppConfigs) {
	logger := config.Configs.LogConfigs.Logger
	logger.WriteLog("Init services", slog.LevelInfo, nil)

	dbconfigs := config.Configs.DBConfigs
	config.Repos = new(da.Repositories)

	config.Repos.ColRepo.DbConfigs = dbconfigs
	config.Repos.UsrRepo.DbConfigs = dbconfigs
	config.Repos.TeamRepo.DbConfigs = dbconfigs
	config.Repos.NoteRepo.DbConfigs = dbconfigs
	config.Repos.SecRepo.DbConfigs = dbconfigs

	config.Repos.ColRepo.MyLogger = logger
	config.Repos.UsrRepo.MyLogger = logger
	config.Repos.TeamRepo.MyLogger = logger
	config.Repos.NoteRepo.MyLogger = logger
	config.Repos.SecRepo.MyLogger = logger
}

func initInterfaces(conf *appconfigs.AppConfigs) {
	dbconf := conf.Configs.DBConfigs
	logger := conf.Configs.LogConfigs

	conf.IRepos = &bl.IRepositories{
		IUsrRepo:  &da.UserRepository{DbConfigs: dbconf, MyLogger: logger.Logger},
		ISecRepo:  &da.SectionRepository{DbConfigs: dbconf, MyLogger: logger.Logger},
		INoteRepo: &da.NoteRepository{DbConfigs: dbconf, MyLogger: logger.Logger},
		IColRepo:  &da.CollectionRepository{DbConfigs: dbconf, MyLogger: logger.Logger},
		ITeamRepo: &da.TeamRepository{DbConfigs: dbconf, MyLogger: logger.Logger},
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

// func initRouter(conf *appconfigs.AppConfigs) {
// 	app.Router.HandleFunc("/api/notes", app.AddNoteHandler).Methods("POST")
//	app.Router.HandleFunc("/api/notes", app.GetAllNotesHandler).Methods("GET")

//	app.Router.HandleFunc("/api/notes/{id}", app.DeleteNoteHandler).Methods("DELETE")

//	app.Router.HandleFunc("/api/sections/{sectionId}/notes/id/{noteId}", app.AddNoteToSectionHandler).Methods("POST")
//	app.Router.HandleFunc("/api/sections/{sectionId}/notes/id/{noteId}", app.DeleteNoteFromSectionHandler).Methods("DELETE")
//
// }

func RunBackend() error {
	configFile := "G:\\Study\\University\\Database Course Project\\notebook app\\config\\config.yaml"

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

	initRepositories(appConfigs)
	initInterfaces(appConfigs)

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

	return nil
}
