package main

import (
	"errors"
	"log"
	"splitwise/configs"
	"splitwise/controller"

	"github.com/gin-gonic/gin"
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

func init() {
	configs.InitConfig()
}

func main() {
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

	r.POST("/addUserToGroup", func(c *gin.Context) {
		controller.AddMembersToGroup(db, c)
	})

	r.GET("/groupDetails", func(c *gin.Context) {
		controller.GetGroupDetails(db, c)
	})

	r.DELETE("/group", func(c *gin.Context) {
		controller.DeleteGroup(db, c)
	})

	r.Run()
}
