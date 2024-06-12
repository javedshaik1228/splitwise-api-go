package controller

import (
	"fmt"
	"net/http"

	"splitwise/models"
	"splitwise/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Credentials struct {
	Username     string `json:"username"`
	Email        string `json:"email"`
	PlainTextPwd string `json:"password"`
}

func SignupHandler(db *gorm.DB, c *gin.Context) {
	var creds Credentials

	if err := c.BindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Binding error": err.Error()})
		return
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(creds.PlainTextPwd), 14)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Internal Error": err.Error()})
	}

	var user models.User
	user.Username = creds.Username
	user.Email = creds.Email
	user.PasswordHash = string(hashedPwd)

	if db.Create(&user).Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Internal error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"Status:": "OK"})
}

func LoginHandler(db *gorm.DB, c *gin.Context) {
	var creds Credentials
	if err := c.BindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Binding error": err.Error()})
		return
	}

	var user models.User
	result := db.Where("username = ?", creds.Username).First(&user)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid username"})
		return
	}
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Internal error": result.Error})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(creds.PlainTextPwd)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"Authentication error": err.Error()})
		return
	}

	accessToken, err := utils.GenerateJWT(user.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Internal error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"AccessToken": &accessToken})
}

func GetUser(db *gorm.DB, id uint) (*models.User, error) {
	var user models.User
	err := db.First(&user, id).Error
	return &user, err
}

func GetUserByUsername(db *gorm.DB, c *gin.Context) (*models.User, error) {
	username := c.Param("username")
	var user models.User
	result := db.Where("username = ?", username).First(&user)
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("user not found")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func GetUserIDfromUsername(db *gorm.DB, username string) (uint, error) {
	var userId uint
	result := db.Model(&models.User{}).Select("UserID").Where("Username =?", username).Scan(&userId)
	if result.RowsAffected == 0 {
		return 0, fmt.Errorf("user not found")
	}
	if result.Error != nil {
		return 0, result.Error
	}
	return userId, nil
}

func GetAllUsers(db *gorm.DB) (*[]models.User, error) {
	var users []models.User
	err := db.Find(&users).Error
	return &users, err
}
