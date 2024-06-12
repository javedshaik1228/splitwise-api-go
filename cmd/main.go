package main

import (
	"errors"
	"log"
	"splitwise/configs"
	"splitwise/controller"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {

	cnxStr := configs.GetDbCnxString()
	if cnxStr == "" {
		log.Fatalf("DB connection string is empty")
		return nil, errors.New("Invalid connection string")
	}

	db, err := gorm.Open(postgres.Open(cnxStr), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func main() {

	err := godotenv.Load("configs/.env")
	if err != nil {
		log.Fatalf("error loading .env file")
	}

	db, err := InitDB()
	if err != nil {
		log.Fatal("Error while connecting to db")
	}

	if db == nil {
		log.Fatalf("DB obj is nil")
	}

	r := gin.Default()

	r.POST("/signup", func(c *gin.Context) {
		controller.SignupHandler(db, c)
	})

	r.POST("/login", func(c *gin.Context) {
		controller.LoginHandler(db, c)
	})

	r.POST("/createGroup", func(c *gin.Context) {
		controller.CreateGroup(db, c)
	})

	r.DELETE("/group", func(c *gin.Context) {
		controller.DeleteGroup(db, c)
	})

	r.Run()
}
