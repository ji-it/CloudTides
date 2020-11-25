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
	DB_HOST                      string          = "106.14.190.68"
	DB_PORT                      string          = "30125"
	DB_USER                      string          = "cloudtides"
	DB_PASSWORD                  string          = "ca$hc0w"
	DB_NAME                      string          = "template1"
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
