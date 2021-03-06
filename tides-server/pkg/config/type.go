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
	DB_HOST                      string          = "10.185.143.234"
	DB_PORT                      string          = "30123"
	DB_USER                      string          = "postgres"
	DB_PASSWORD                  string          = "GZzmKrHCBNSMSkkJd2Fm5knqPYhWQEHECJToZwaxgqIFFCo4kb4KOAZDxmGFIbZq"
	DB_NAME                      string          = "tides"
)

var (
	ContextRoot = context.Background()
)

// Config consists fields to setup the clouTides server
type Config struct {
	Port string `json:"port"`
	DB   string `json:"database"`
}
