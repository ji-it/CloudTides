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
	dbUser                       string          = "postgres"
	dbPassword                   string          = "Shen1997"
	dbName                       string          = "test"
	adminUser                    string          = "admin"
	adminPwd                     string          = "nimda#sedit"
)

var (
	ContextRoot = context.Background()
)

// Config consists fields to setup the clouTides server
type Config struct {
	Port string `json:"port"`
	DB   string `json:"database"`
}
