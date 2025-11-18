package logger

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

func New() *logrus.Logger {
	log := logrus.New()
	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	log.SetLevel(logrus.InfoLevel)
	return log
}

func SetLevel(log *logrus.Logger, level string) {
	if level == "" {
		return
	}
	parsed, err := logrus.ParseLevel(strings.ToLower(level))
	if err != nil {
		log.WithError(err).Warnf("invalid log level %q, falling back to info", level)
		return
	}
	log.SetLevel(parsed)
}
