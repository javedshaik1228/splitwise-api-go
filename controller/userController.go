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
		SendError(ErrBadRequest, err.Error(), c)
		return
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(creds.PlainTextPwd), 14)
	if err != nil {
		SendError(ErrInternalFailure, err.Error(), c)
	}

	var user models.User
	user.Username = creds.Username
	user.Email = creds.Email
	user.PasswordHash = string(hashedPwd)

	if err := db.Create(&user).Error; err != nil {
		SendError(ErrInternalFailure, err.Error(), c)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"Status:": "OK"})
}

func LoginHandler(db *gorm.DB, c *gin.Context) {
	var creds Credentials
	if err := c.BindJSON(&creds); err != nil {
		SendError(ErrBadRequest, err.Error(), c)
		return
	}

	var user models.User
	result := db.Where("username = ?", creds.Username).First(&user)
	if result.RowsAffected == 0 {
		SendError(ErrNotFound, "Invalid username", c)
		return
	}
	if result.Error != nil {
		SendError(ErrInternalFailure, result.Error.Error(), c)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(creds.PlainTextPwd)); err != nil {
		SendError(ErrUnauthorized, "Password mismatch", c)
		return
	}

	accessToken, err := utils.GenerateJWT(user.UserID)
	if err != nil {
		SendError(ErrInternalFailure, "Unable to create access token", c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"AccessToken": &accessToken})
}

func GetAllUsers(db *gorm.DB) (*[]models.User, error) {
	var users []models.User
	err := db.Find(&users).Error
	return &users, err
}

func getUserIDfromUsername(db *gorm.DB, username string) (uint, error) {
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

func getUsernameFromUserID(db *gorm.DB, userID uint) (string, error) {
	var user models.User
	err := db.Select("username").First(&user, userID).Error
	if err != nil {
		return "", err
	}
	return user.Username, nil
}
