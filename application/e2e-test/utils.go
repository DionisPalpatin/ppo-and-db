package e2etest

import (
	"database/sql"
	"fmt"
	"time"
	"os"
)

type PostgresInfo struct {
	Host     string
	User     string
	Password string
	Port     int
	DBName   string
}

func InitDB(configs PostgresInfo) (*sql.DB, error) {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		configs.Host, configs.Port, configs.User, configs.Password, configs.DBName)

	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		return nil, err
	}

	for {
		err := db.Ping()
		if err == nil {
			fmt.Println("Соединение с БД установлено!")
			break
		} else {
			fmt.Println("Fail to connect to db, trying again...")
			time.Sleep(5 * time.Second)
		}
	}

	return db, nil
}

func IsUnitTestsFailed() bool {
	return os.Getenv("UNIT_SUCCESS") != "1"
}

func IsIntegrationTestsFailed() bool {
	return os.Getenv("INTEGRATION_SUCCESS") != "1"
}