package logger

import (
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var once sync.Once

var Logger *logrus.Entry

func InitializeLogger(fields logrus.Fields) {
	once.Do(func() {
		Logger = logrus.WithFields(fields)
	})
}

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{TimestampFormat: time.RFC3339Nano, FullTimestamp: true})
	logrus.SetLevel(logrus.InfoLevel)

	if level := os.Getenv("LOG_LEVEL"); level != "" {
		if l, err := strconv.Atoi(level); err == nil {
			logrus.SetLevel(logrus.Level(l))
		}
	}

	env := os.Getenv("ENV")
	Logger = logrus.WithField("env", env)
}
