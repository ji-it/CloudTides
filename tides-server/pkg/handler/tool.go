package handler

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"

	"tides-server/pkg/config"
	"tides-server/pkg/logger"
	"tides-server/pkg/models"
)

func ParseUserIdFromToken(req *http.Request) (uint, error) {
	reqToken := req.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	if len(splitToken) < 2 {
		logger.SetLogLevel("ERROR")
		logger.Error("Token not supplied in request")
		return 0, errors.New("Token not supplied in request!")
	}
	stringToken := splitToken[1]

	claims := &Claims{}
	_, err := jwt.ParseWithClaims(
		stringToken,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		},
	)
	if claims.ExpiresAt < time.Now().UTC().Unix() {
		logger.SetLogLevel("ERROR")
		logger.Error("JWT is expired")
		return 0, errors.New("JWT is expired")
	}
	return claims.Id, err
}

func VerifyUser(req *http.Request) bool {
	id, err := ParseUserIdFromToken(req)
	if err != nil {
		return false
	}
	db := config.GetDB()
	var queryUser models.User
	if db.Where("id = ?", id).First(&queryUser).Error != nil {
		return false
	}

	return true
}
