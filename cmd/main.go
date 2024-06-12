package main

import (
	"log"
	"splitwise/configs"
	"splitwise/controller"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := configs.InitDB()
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
