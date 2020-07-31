package config

import (
	"log/syslog"
	"os"
)

type Logger struct {
	syslog.Writer
}

type LoggingRemoteOpts struct {
	RemoteProtocol string          `json:"protocol"`
	RemoteServer   string          `json:"remote_server"`
	Flag           int             `json:"flag"`
	Priority       syslog.Priority `json:"priority"`
	Tag            string          `json:"tag"`
}

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
