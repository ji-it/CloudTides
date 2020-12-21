package config

import (
	"log/syslog"
	"os"
)

// Logger by default
type Logger struct {
	syslog.Writer
}

// LoggingRemoteOpts remote logging options
type LoggingRemoteOpts struct {
	RemoteProtocol string          `json:"protocol"`
	RemoteServer   string          `json:"remote_server"`
	Flag           int             `json:"flag"`
	Priority       syslog.Priority `json:"priority"`
	Tag            string          `json:"tag"`
}

// LoggingLocalOpts local logging options
type LoggingLocalOpts struct {
	Output   *os.File
	Flag     int
	Priority syslog.Priority
	Tag      string
}

func defaultLoggingLocalOpts() (logopts *LoggingLocalOpts) {
	return &LoggingLocalOpts{
		Output:   os.Stdout,
		Flag:     defaultLoggingFlag,
		Priority: defaultLoggingPriority,
		Tag:      defaultLoggingTag,
	}
}
