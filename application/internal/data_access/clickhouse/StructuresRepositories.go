package postgres

import (
	_ "github.com/lib/pq"

	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/config"
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/logger"
)

// ---------------------------------------------------------------------------------------------------------------------
// Structures
// ---------------------------------------------------------------------------------------------------------------------

type UserRepository struct {
	DbConfigs *config.DBConfigs
	MyLogger  *mylogger.MyLogger
}

type TeamRepository struct {
	DbConfigs *config.DBConfigs
	MyLogger  *mylogger.MyLogger
}

type SectionRepository struct {
	DbConfigs *config.DBConfigs
	MyLogger  *mylogger.MyLogger
}

type NoteRepository struct {
	DbConfigs *config.DBConfigs
	MyLogger  *mylogger.MyLogger
}

type CollectionRepository struct {
	DbConfigs *config.DBConfigs
	MyLogger  *mylogger.MyLogger
}

type Repositories struct {
	UsrRepo  UserRepository
	SecRepo  SectionRepository
	NoteRepo NoteRepository
	ColRepo  CollectionRepository
	TeamRepo TeamRepository
}
