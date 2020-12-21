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
)

// ContextRoot
var (
	ContextRoot = context.Background()
)

// Config consists fields to setup the CloudTides server
type Config struct {
	ServerIP         string `json:"serverIP"`
	ServerPort       string `json:"serverPort"`
	PostgresHost     string `json:"postgresHost"`
	PostgresPort     string `json:"postgresPort"`
	PostgresUser     string `json:"postgresUser"`
	PostgresPassword string `json:"postgresPassword"`
	PostgresDB       string `json:"postgresDB"`
	SecretKey        string `json:"secretKey"`
	AdminUser        string `json:"adminUser"`
	AdminPassword    string `json:"adminPassword"`
}
