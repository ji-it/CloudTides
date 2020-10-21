package handler

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/vmware/go-vcloud-director/v2/govcd"

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

// Creates a vCD client
func (c *VcdConfig) Client() (*govcd.VCDClient, error) {
	u, err := url.ParseRequestURI(c.Href)
	if err != nil {
		return nil, fmt.Errorf("unable to pass url: %s", err)
	}

	vcdClient := govcd.NewVCDClient(*u, c.Insecure)
	if c.Token != "" {
		_ = vcdClient.SetToken(c.Org, govcd.AuthorizationHeader, c.Token)
	} else {
		_, err := vcdClient.GetAuthResponse(c.User, c.Password, c.Org)
		if err != nil {
			return nil, fmt.Errorf("unable to authenticate: %s", err)
		}
		// fmt.Printf("Token: %s\n", resp.Header[govcd.AuthorizationHeader])
	}
	return vcdClient, nil
}

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
