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
	DB_HOST                      string          = "192.168.0.162"
	DB_USER                      string          = "postgres"
	DB_PASSWORD                  string          = "postgres"
	DB_NAME                      string          = "test"
	adminUser                    string          = "admin"
)

var (
	ContextRoot = context.Background()
)

// Config consists fields to setup the clouTides server
type Config struct {
	Port string `json:"port"`
	DB   string `json:"database"`
}
