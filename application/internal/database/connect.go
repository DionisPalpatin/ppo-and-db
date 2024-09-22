package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/config"
	mylogger "github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/logger"
)

func Connect(logger *mylogger.MyLogger, configs *config.DBConfigs) error {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		configs.Host, configs.Port, configs.User, configs.Password, configs.Name)

	db, err := sql.Open(configs.DriverName, connectionString)

	if err != nil {
		return err
	}

	for {
		err := db.Ping()
		if err == nil {
			fmt.Println("Соединение с БД установлено!")
			break
		} else {
			logger.WriteLog("Fail to connect to db, trying again...", slog.LevelInfo, nil)
			time.Sleep(5 * time.Second) // Пауза между попытками
		}
	}

	configs.DB = db
	return nil
}
