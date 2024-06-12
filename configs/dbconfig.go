package configs

import (
	"errors"
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DbConfig struct {
	db_user string
	db_pwd  string
	db_name string
	db_host string
	db_port string
}

func GetDbConfig() (*DbConfig, error) {
	dbUser := os.Getenv("SW_DB_USER")
	dbPwd := os.Getenv("SW_DB_PASSWORD")
	dbName := os.Getenv("SW_DB_NAME")
	dbHost := os.Getenv("SW_DB_HOST")
	dbPort := os.Getenv("SW_DB_PORT")

	if dbUser == "" {
		return nil, errors.New("SW_DB_USER environment variable is required")
	}
	if dbPwd == "" {
		return nil, errors.New("SW_DB_PASSWORD environment variable is required")
	}
	if dbName == "" {
		return nil, errors.New("SW_DB_NAME environment variable is required")
	}
	if dbHost == "" {
		return nil, errors.New("SW_DB_HOST environment variable is required")
	}
	if dbPort == "" {
		return nil, errors.New("SW_DB_PORT environment variable is required")
	}

	return &DbConfig{
		db_user: dbUser,
		db_pwd:  dbPwd,
		db_name: dbName,
		db_host: dbHost,
		db_port: dbPort,
	}, nil
}

func InitDB() (*gorm.DB, error) {
	config, err := GetDbConfig()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return nil, err
	}

	dsn := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s",
		config.db_host, config.db_user, config.db_name, config.db_pwd, config.db_port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
