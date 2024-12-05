package app

import (
	"fmt"
	dapostgres "github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/data_access"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/config"
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/UI/TechUI"
	handlers "github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/api/v2"
	app "github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/app/handlers"
	appconfigs "github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/app/handlers"
	bl "github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/business_logic"
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

// Структура для общей конфигурации приложения
type App struct {
	Configs *config.Configs
	Router  *mux.Router

	IServices *bl.IServices
	IRepos    *bl.IRepositories
}

func initInterfaces(conf *app.App) {
	dbconf := conf.Configs.DBConfigs
	logger := conf.Configs.LogConfigs

	if dbconf.DriverName == "postgres" {
		conf.IRepos = &bl.IRepositories{
			IUsrRepo:  &dapostgres.UserRepository{DbConfigs: dbconf, MyLogger: logger.Logger},
			ISecRepo:  &dapostgres.SectionRepository{DbConfigs: dbconf, MyLogger: logger.Logger},
			INoteRepo: &dapostgres.NoteRepository{DbConfigs: dbconf, MyLogger: logger.Logger},
			IColRepo:  &dapostgres.CollectionRepository{DbConfigs: dbconf, MyLogger: logger.Logger},
			ITeamRepo: &dapostgres.TeamRepository{DbConfigs: dbconf, MyLogger: logger.Logger},
			IStatRepo: &dapostgres.StatisticRepository{DbConfigs: dbconf, MyLogger: logger.Logger},
		}
	}
	// } else if dbconf.DriverName == "clickhouse" {
	// 	conf.IRepos = &bl.IRepositories{
	// 		IUsrRepo:  &daclickhouse.UserRepository{DbConfigs: dbconf, MyLogger: logger.Logger},
	// 		ISecRepo:  &daclickhouse.SectionRepository{DbConfigs: dbconf, MyLogger: logger.Logger},
	// 		INoteRepo: &daclickhouse.NoteRepository{DbConfigs: dbconf, MyLogger: logger.Logger},
	// 		IColRepo:  &daclickhouse.CollectionRepository{DbConfigs: dbconf, MyLogger: logger.Logger},
	// 		ITeamRepo: &daclickhouse.TeamRepository{DbConfigs: dbconf, MyLogger: logger.Logger},
	// 		IStatRepo: &daclickhouse.StatisticRepository{DbConfigs: dbconf, MyLogger: logger.Logger},
	// 	}
	// }

	conf.IServices = &bl.IServices{
		IUsrSvc:   &bl.UserService{},
		ISecSvc:   &bl.SectionService{},
		INoteSvc:  &bl.NoteService{},
		IColSvc:   &bl.CollectionService{},
		ITeamSvc:  &bl.TeamService{},
		IOAuthSvc: &bl.OAuthService{},
		IStatSvc:  &bl.StatService{},
	}
}

func RunBackend() error {
	configFile := "./config/config.yaml"

	app := new(App)
	var err error

	app.Configs, err = config.ReadConfig(configFile)
	if err != nil {
		return err
	}
	configs := app.Configs
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

	// initRepositories(app)
	initInterfaces((*appconfigs.App)(app))

	if app.Configs.Mode == "tech" {
		for {
			user := TechUI.AuthorizationMenu(configs, app.IRepos, app.IServices)
			if user == nil {
				break
			} else if user.Role == bl.Reader {
				TechUI.ReaderMenu(user, configs, app.IRepos, app.IServices)
			} else if user.Role == bl.Author {
				TechUI.AuthorMenu(user, configs, app.IRepos, app.IServices)
			} else if user.Role == bl.Admin {
				TechUI.AdminMenu(user, configs, app.IRepos, app.IServices)
			}
		}
	} else {
		handlers.InitRouter(&app.Router)
		port := fmt.Sprintf(":%d", app.Configs.ServerPort)

		fs := http.FileServer(http.Dir("./static/"))
		app.Router.PathPrefix("/").Handler(fs)

		http.ListenAndServe(port, app.Router)
		app.Configs.LogConfigs.Logger.WriteLog("Server is running on port 8080", slog.LevelInfo, nil)
	}

	return nil
}
