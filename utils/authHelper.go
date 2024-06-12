package utils

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uint
	jwt.RegisteredClaims
}

var jwtKey = []byte(os.Getenv("SW_JWT_KEY"))

func GenerateJWT(userId uint) (string, error) {

	expirationTime := time.Now().Add(60 * time.Minute)
	// TODO change expiration time
	claims := &Claims{
		UserID: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "splitwise",
			Subject:   "user",
			Audience:  jwt.ClaimStrings{"splitwise_users"},
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateUserLogin(c *gin.Context) (*Claims, bool) {
	tokenString := c.Request.Header.Get("token")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"Auth error": "No authorization token provided"})
		return nil, false
	}

	claims, err := validateJWT(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"Auth error": err.Error()})
		return nil, false
	}
	return claims, true
}

func validateJWT(tokenString string) (*Claims, error) {

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	expTime, err := claims.GetExpirationTime()

	if err != nil {
		return nil, fmt.Errorf("invalid token")
	}

	if expTime.Before(time.Now()) {
		return nil, fmt.Errorf("token expired, please login again")
	}

	return claims, nil
}
