package configs

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type DbConfiguration struct {
	db_user string
	db_pwd  string
	db_name string
	db_host string
	db_port string
}

type AppConfiguration struct {
	JwtKey          []byte
	TokenExpiryTime time.Duration
}

var DbConfig DbConfiguration
var AppConfig AppConfiguration

func InitConfig() {
	// get configs envfilepath from an env var
	// so that tests can have a different env path independent of main one
	envFilePath := os.Getenv("ENV_FILE_PATH")
	if envFilePath == "" {
		// If ENV_FILE_PATH environment variable is not set, use a default path
		envFilePath = "configs/.env"
	}

	err := godotenv.Load(envFilePath)
	if err != nil {
		log.Fatalf("error loading .env file")
		return
	}
	// parse token expiry time
	var expTimeStr = os.Getenv("AUTH_EXP_TIME")
	tokenExpTime, err := strconv.Atoi(expTimeStr)
	if err != nil {
		log.Fatalf("AUTH_EXP_TIME must be a valid integer")
	}

	// set global variables
	AppConfig.JwtKey = []byte(os.Getenv("JWT_KEY"))
	AppConfig.TokenExpiryTime = time.Duration(tokenExpTime) * time.Minute

	if err := initDbConfig(&DbConfig); err != nil {
		log.Fatalf("Error while loading db config: %v", err.Error())
	}
}

func initDbConfig(DbConfig *DbConfiguration) error {
	dbUser := os.Getenv("DB_USER")
	dbPwd := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	if dbUser == "" {
		return errors.New("DB_USER environment variable is required")
	}
	if dbPwd == "" {
		return errors.New("DB_PASSWORD environment variable is required")
	}
	if dbName == "" {
		return errors.New("DB_NAME environment variable is required")
	}
	if dbHost == "" {
		return errors.New("DB_HOST environment variable is required")
	}
	if dbPort == "" {
		return errors.New("DB_PORT environment variable is required")
	}

	DbConfig.db_user = dbUser
	DbConfig.db_pwd = dbPwd
	DbConfig.db_name = dbName
	DbConfig.db_host = dbHost
	DbConfig.db_port = dbPort

	return nil
}

func GetDbCnxString() string {
	cnxStr := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s",
		DbConfig.db_host, DbConfig.db_user, DbConfig.db_name, DbConfig.db_pwd, DbConfig.db_port)
	return cnxStr
}
