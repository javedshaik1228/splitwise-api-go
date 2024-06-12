package utils

import (
	"errors"
	"os"
	"splitwise/configs"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uint
	jwt.RegisteredClaims
}

func validateJWT(tokenString string) (*Claims, error) {
	var jwtKey = []byte(os.Getenv("JWT_KEY"))
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

func GenerateJWT(userId uint) (string, error) {
	appConfig := configs.AppConfig

	if appConfig.JwtKey == nil || appConfig.TokenExpiryTime < 0 {
		return "", errors.New("invalid app configs")
	}

	tokenExpiresAt := time.Now().Add(appConfig.TokenExpiryTime)

	claims := &Claims{
		UserID: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "splitwise",
			Subject:   "user",
			Audience:  jwt.ClaimStrings{"splitwise_users"},
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(tokenExpiresAt),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(appConfig.JwtKey)
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
