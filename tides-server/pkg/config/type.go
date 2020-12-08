package config

import (
	"context"
	"log"
	"log/syslog"
)

const (
	defaultLoggingRemoteProtocol string          = "tcp"
	defaultLoggingFlag           int             = log.LstdFlags
	defaultLoggingPriority       syslog.Priority = syslog.LOG_INFO
	defaultLoggingTag            string          = "CloudTides-Server"
	defaultPort                  string          = "80"
	dbHost                       string          = "106.14.190.68"
	dbPort                       string          = "30125"
	dbUser                       string          = "cloudtides"
	dbPassword                   string          = "ca$hc0w"
	dbName                       string          = "template1"
	adminUser                    string          = "admin"
	adminPwd                     string          = "nimda#sedit"
)

var (
	ContextRoot = context.Background()
)

// Config consists fields to setup the cloudTides server
type Config struct {
	Port  string `json:"port"`
	DB    string `json:"database"`
	Debug bool   `json:"debug"`
}
