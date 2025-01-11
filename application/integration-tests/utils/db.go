package integr_ut

import (
	"database/sql"
	"fmt"
	"time"
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