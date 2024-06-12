package utils

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uint
	jwt.RegisteredClaims
}

var jwtKey = []byte(os.Getenv("JWT_KEY"))

func getExpirationTime() (*time.Time, error) {
	var expTimeStr = os.Getenv("AUTH_EXP_TIME")
	authExpTime, err := strconv.Atoi(expTimeStr)
	if err != nil {
		return nil, errors.New("AUTH_EXP_TIME must be a valid integer")
	}
	expirationTime := time.Now().Add(time.Duration(authExpTime) * time.Minute)
	return &expirationTime, nil
}

func GenerateJWT(userId uint) (string, error) {

	tokenExpTime, err := getExpirationTime()
	if err != nil {
		return "", err
	}

	claims := &Claims{
		UserID: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "splitwise",
			Subject:   "user",
			Audience:  jwt.ClaimStrings{"splitwise_users"},
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(*tokenExpTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateUserLogin(c *gin.Context) (*Claims, error) {
	tokenString := c.Request.Header.Get("token")
	if tokenString == "" {
		return nil, errors.New("no authorization token provided")
	}

	claims, err := validateJWT(tokenString)
	if err != nil {
		return nil, err
	}
	return claims, nil
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
		return nil, errors.New("invalid token")
	}

	expTime, err := claims.GetExpirationTime()

	if err != nil {
		return nil, errors.New("invalid token")
	}

	if expTime.Before(time.Now()) {
		return nil, errors.New("token expired, please login again")
	}

	return claims, nil
}
