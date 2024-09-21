package config

import (
	"database/sql"

	"gopkg.in/yaml.v2"

	"fmt"
	"os"

	"notebook_app/internal/logger"
)

func ReadConfig(filePath string) (*Configs, error) {
	// Чтение файла конфигурации
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения файла конфигурации: %w", err)
	}

	// Декодирование YAML-данных
	var config Configs
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("ошибка декодирования YAML: %w", err)
	}

	return &config, nil
}

// Структура для конфигурации PostgreSQL
type DBConfigs struct {
	Host       string `yaml:"psql_host"`
	Port       int    `yaml:"psql_port"`
	User       string `yaml:"psql_user"`
	Password   string `yaml:"psql_password"`
	Name       string `yaml:"psql_db_name"`
	SchemaName string `yaml:"psql_schema_name"`
	DriverName string `yaml:"db_driver_name"`

	DB *sql.DB
}

type LoggerConfigs struct {
	LogLevel string `yaml:"log_level"`
	LogFile  string `yaml:"log_file"`

	Logger *mylogger.MyLogger
}

// Структура для общей конфигурации приложения
type Configs struct {
	Mode           string   `yaml:"ui_mode"`
	DateTimeFormat string   `yaml:"datetime_format"`
	System         int      `yaml:"system"`
	TextTypes      []string `yaml:"text_types"`
	ImageTypes     []string `yaml:"image_types"`

	DBConfigs  *DBConfigs     `yaml:"postgres"`
	LogConfigs *LoggerConfigs `yaml:"logger"`
}
