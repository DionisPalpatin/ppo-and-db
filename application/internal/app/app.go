package app

import (
	"fmt"
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/config"
	dapostgres "github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/data_access"
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/database"
	"log/slog"
	//"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/UI/TechUI"
	handlers "github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/api/v2"
	bl "github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/business_logic"
	//"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/database"
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

// Структура для общей конфигурации приложения
type App struct {
	Configs  *config.Configs
	handlers *handlers.HandlersStruct

	IServices *bl.IServices
	IRepos    *bl.IRepositories
}

func initInterfaces(appStruct *App) {
	dbconf := appStruct.Configs.DBConfigs
	logger := appStruct.Configs.LogConfigs

	if dbconf.DriverName == "postgres" {
		appStruct.IRepos = &bl.IRepositories{
			IUsrRepo:  &dapostgres.UserRepository{DbConfigs: dbconf, MyLogger: logger.Logger},
			ISecRepo:  &dapostgres.SectionRepository{DbConfigs: dbconf, MyLogger: logger.Logger},
			INoteRepo: &dapostgres.NoteRepository{DbConfigs: dbconf, MyLogger: logger.Logger},
			IColRepo:  &dapostgres.CollectionRepository{DbConfigs: dbconf, MyLogger: logger.Logger},
			ITeamRepo: &dapostgres.TeamRepository{DbConfigs: dbconf, MyLogger: logger.Logger},
			//IStatRepo: &dapostgres.StatisticRepository{DbConfigs: dbconf, MyLogger: logger.Logger},
		}
	}
	// } else if dbconf.DriverName == "clickhouse" {
	// 	appStruct.IRepos = &bl.IRepositories{
	// 		IUsrRepo:  &daclickhouse.UserRepository{DbConfigs: dbconf, MyLogger: logger.Logger},
	// 		ISecRepo:  &daclickhouse.SectionRepository{DbConfigs: dbconf, MyLogger: logger.Logger},
	// 		INoteRepo: &daclickhouse.NoteRepository{DbConfigs: dbconf, MyLogger: logger.Logger},
	// 		IColRepo:  &daclickhouse.CollectionRepository{DbConfigs: dbconf, MyLogger: logger.Logger},
	// 		ITeamRepo: &daclickhouse.TeamRepository{DbConfigs: dbconf, MyLogger: logger.Logger},
	// 		IStatRepo: &daclickhouse.StatisticRepository{DbConfigs: dbconf, MyLogger: logger.Logger},
	// 	}
	// }

	appStruct.IServices = &bl.IServices{
		IUsrSvc:   &bl.UserService{},
		ISecSvc:   &bl.SectionService{},
		INoteSvc:  &bl.NoteService{},
		IColSvc:   &bl.CollectionService{},
		ITeamSvc:  &bl.TeamService{},
		IOAuthSvc: &bl.OAuthService{},
		//IStatSvc:  &bl.StatService{},
	}

	appStruct.handlers = &handlers.HandlersStruct{}
	appStruct.handlers.IServices = &bl.IServices{
		IUsrSvc:   &bl.UserService{},
		ISecSvc:   &bl.SectionService{},
		INoteSvc:  &bl.NoteService{},
		IColSvc:   &bl.CollectionService{},
		ITeamSvc:  &bl.TeamService{},
		IOAuthSvc: &bl.OAuthService{},
		//IStatSvc:  &bl.StatService{},
	}
	appStruct.handlers.IRepos = &bl.IRepositories{
		IUsrRepo:  &dapostgres.UserRepository{DbConfigs: dbconf, MyLogger: logger.Logger},
		ISecRepo:  &dapostgres.SectionRepository{DbConfigs: dbconf, MyLogger: logger.Logger},
		INoteRepo: &dapostgres.NoteRepository{DbConfigs: dbconf, MyLogger: logger.Logger},
		IColRepo:  &dapostgres.CollectionRepository{DbConfigs: dbconf, MyLogger: logger.Logger},
		ITeamRepo: &dapostgres.TeamRepository{DbConfigs: dbconf, MyLogger: logger.Logger},
		//IStatRepo: &dapostgres.StatisticRepository{DbConfigs: dbconf, MyLogger: logger.Logger},
	}
}

func RunBackend() error {
	configFile := "./config.yaml"

	appStruct := new(App)
	var err error

	appStruct.Configs, err = config.ReadConfig(configFile)
	if err != nil {
		return err
	}
	configs := appStruct.Configs
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

	initInterfaces(appStruct)

	//if appStruct.Configs.Mode == "tech" {
	//	for {
	//		user := TechUI.AuthorizationMenu(configs, appStruct.IRepos, appStruct.IServices)
	//		if user == nil {
	//			break
	//		} else if user.Role == bl.Reader {
	//			TechUI.ReaderMenu(user, configs, appStruct.IRepos, appStruct.IServices)
	//		} else if user.Role == bl.Author {
	//			TechUI.AuthorMenu(user, configs, appStruct.IRepos, appStruct.IServices)
	//		} else if user.Role == bl.Admin {
	//			TechUI.AdminMenu(user, configs, appStruct.IRepos, appStruct.IServices)
	//		}
	//	}
	//} else {
	//	handlers.InitRouter(appStruct.handlers)
	//	port := fmt.Sprintf(":%d", appStruct.Configs.ServerPort)
	//
	//	fs := http.FileServer(http.Dir("./static/"))
	//	appStruct.handlers.Router.PathPrefix("/").Handler(fs)
	//
	//	http.ListenAndServe(port, appStruct.handlers.Router)
	//	appStruct.Configs.LogConfigs.Logger.WriteLog("Server is running on port 8080", slog.LevelInfo, nil)
	//}

	handlers.InitHandlersStructFields(appStruct.handlers, appStruct.Configs, appStruct.IRepos, appStruct.IServices)
	handlers.InitRouter(appStruct.handlers)
	port := fmt.Sprintf(":%d", appStruct.Configs.ServerPort)

	err = appStruct.handlers.Router.Run(port)
	if err != nil {
		appStruct.Configs.LogConfigs.Logger.WriteLog("Fail run server", slog.LevelError, nil)
	}
	appStruct.Configs.LogConfigs.Logger.WriteLog("Server is running on port 8080", slog.LevelInfo, nil)

	return nil
}
