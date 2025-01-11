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

// type SectionRepository struct {
// 	DbConfigs *config.DBConfigs
// 	MyLogger  *mylogger.MyLogger
// }

type NoteRepository struct {
	DbConfigs *config.DBConfigs
	MyLogger  *mylogger.MyLogger
}

type CollectionRepository struct {
	DbConfigs *config.DBConfigs
	MyLogger  *mylogger.MyLogger
}

func NewUserRepository(dbconfigs *config.DBConfigs, MyLogger  *mylogger.MyLogger)UserRepository {
	return UserRepository{DbConfigs: dbconfigs, MyLogger: MyLogger}
}

func NewTeamRepository(dbconfigs *config.DBConfigs, MyLogger  *mylogger.MyLogger) TeamRepository {
	return TeamRepository{DbConfigs: dbconfigs, MyLogger: MyLogger}
}

func NewNoteRepository(dbconfigs *config.DBConfigs, MyLogger  *mylogger.MyLogger) NoteRepository {
	return NoteRepository{DbConfigs: dbconfigs, MyLogger: MyLogger}
}

func NewCollectionRepository(dbconfigs *config.DBConfigs, MyLogger  *mylogger.MyLogger) CollectionRepository {
	return CollectionRepository{DbConfigs: dbconfigs, MyLogger: MyLogger}
}
