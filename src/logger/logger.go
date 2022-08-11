package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

var loggers = make(map[string]*logger)

type logger struct {
	logrus.Logger
	name string
}

func GetLogger(name string) *logger {
	if logger, ok := loggers[name]; ok {
		return logger
	}
	logger := newLogger(name)
	loggers[name] = logger
	return logger
}

func newLogger(name string) *logger {
	logger := &logger{Logger: *logrus.New(), name: name}
	logger.SetReportCaller(true)
	logger.SetLevel(logrus.DebugLevel)

	fileName := "/tmp/countdown/log/log.txt"
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logger.Error("failed to log to file")
	} else {
		logger.Out = file
	}
	return logger
}
