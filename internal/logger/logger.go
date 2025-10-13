// Package logger provides a simple wrapper around logrus for logging.
// In this application, we use logrus for structured logging.
// To simplify logging calls, we use global logger functions and do not test logging output.

package logger

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

func Init() {
	levelStr := os.Getenv("LOG_LEVEL")
	if levelStr == "" {
		levelStr = "info"
	}
	level, _ := logrus.ParseLevel(strings.ToLower(levelStr))
	logrus.SetLevel(level)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
}

var (
	Debug  = logrus.Debug
	Debugf = logrus.Debugf
	Info   = logrus.Info
	Infof  = logrus.Infof
	Warn   = logrus.Warn
	Warnf  = logrus.Warnf
	Error  = logrus.Error
	Errorf = logrus.Errorf
	Fatal  = logrus.Fatal
	Fatalf = logrus.Fatalf
)
