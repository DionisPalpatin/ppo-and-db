package data_access

import (
	mylogger "github.com/DionisPalpatin/ppo-and-db/application/internal/logger"
	"github.com/DionisPalpatin/ppo-and-db/application/internal/config"
	_ "github.com/lib/pq"
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

type StatisticRepository struct {
	DbConfigs *config.DBConfigs
	MyLogger  *mylogger.MyLogger
}
