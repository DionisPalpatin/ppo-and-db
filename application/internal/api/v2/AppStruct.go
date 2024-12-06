package handlersv2

import (
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/config"
	bl "github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/business_logic"
	"github.com/gorilla/mux"
)

type HandlersStruct struct {
	Configs *config.Configs
	Router  *mux.Router

	IServices *bl.IServices
	IRepos    *bl.IRepositories
}
