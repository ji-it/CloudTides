package logger

import (
	"log"
	"os"
	"sync"
)

var (
	consoleLogger *log.Logger
	onceConsole   sync.Once
)

// GetConsoleLogger returns a pointer to consoleLogger
func GetConsoleLogger() *log.Logger {
	onceConsole.Do(func() {
		consoleLogger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	})
	return consoleLogger
}
