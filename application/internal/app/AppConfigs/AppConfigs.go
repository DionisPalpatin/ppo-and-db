package appconfigs

import (
	"github.com/gorilla/mux"

	"notebook_app/config"
	"notebook_app/internal/business_logic"
	"notebook_app/internal/data_access"
	"notebook_app/internal/models"
)

// Структура для общей конфигурации приложения
type AppConfigs struct {
	Configs *config.Configs
	Router  mux.Router
	CurUser *models.User

	IServices *bl.IServices
	IRepos    *bl.IRepositories
	Repos     *da.Repositories
}
