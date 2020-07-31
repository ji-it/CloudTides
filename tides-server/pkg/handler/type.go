package handler

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	secretKey  string        = "i7@q3rhti=*m3tokpaf@15qgxh15d8-o#-9l1)ke9)e2ec079)"
	expireTime time.Duration = 24 * time.Hour
	issuer     string        = "CloudTides"
)

// Configuration of Json Web Token
type Claims struct {
	Id uint
	jwt.StandardClaims
}
