package configs

import (
	"errors"
	"fmt"
	"log"
	"os"
)

type DbConfig struct {
	db_user string
	db_pwd  string
	db_name string
	db_host string
	db_port string
}

func getDbConfig() (*DbConfig, error) {
	dbUser := os.Getenv("DB_USER")
	dbPwd := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	if dbUser == "" {
		return nil, errors.New("DB_USER environment variable is required")
	}
	if dbPwd == "" {
		return nil, errors.New("DB_PASSWORD environment variable is required")
	}
	if dbName == "" {
		return nil, errors.New("DB_NAME environment variable is required")
	}
	if dbHost == "" {
		return nil, errors.New("DB_HOST environment variable is required")
	}
	if dbPort == "" {
		return nil, errors.New("DB_PORT environment variable is required")
	}

	return &DbConfig{
		db_user: dbUser,
		db_pwd:  dbPwd,
		db_name: dbName,
		db_host: dbHost,
		db_port: dbPort,
	}, nil
}

func GetDbCnxString() string {
	config, err := getDbConfig()
	if err != nil {
		log.Printf("Error: %v\n", err)
		return ""
	}
	cnxStr := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s",
		config.db_host, config.db_user, config.db_name, config.db_pwd, config.db_port)
	return cnxStr
}
