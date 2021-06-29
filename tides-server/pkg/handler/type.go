package handler

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	expireTime      time.Duration = 24 * time.Hour
	issuer          string        = "CloudTides"
	temName         string        = "tides-boinc-attached"
	vmName          string        = "tides-gromacs"
	vappName        string        = "tides-vapp-setup"
	idlePolicy      string        = "{\"cpu\": 0.6, \"ram\": 0.6}"
	thresholdPolicy string        = "{\"cpu\": 0.8, \"ram\": 0.8}"
)

var (
	letters = []rune("abcdefghijklmnopqrstuvwxyz")
	network string
	catalog string
)

// Claims is configuration of Json Web Token
type Claims struct {
	ID uint
	jwt.StandardClaims
}

// VcdConfig is configuration for Vcd Connection
/*type VcdConfig struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Org      string `json:"org"`
	Href     string `json:"href"`
	VDC      string `json:"vdc"`
	Insecure bool   `json:"insecure"`
	Token    string `json:"token"`
}*/
