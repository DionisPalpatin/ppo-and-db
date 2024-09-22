package appconfigs

import (
	"github.com/gorilla/mux"

	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/config"
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/business_logic"
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/data_access"
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/models"
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
