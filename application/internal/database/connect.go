package database

import (
	"database/sql"
	"fmt"

	"notebook_app/config"
)

func Connect(configs *config.DBConfigs) error {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		configs.Host, configs.Port, configs.User, configs.Password, configs.Name)

	db, err := sql.Open(configs.DriverName, connectionString)

	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	configs.DB = db
	return nil
}
