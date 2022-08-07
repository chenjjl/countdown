package logger

import "github.com/sirupsen/logrus"

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
	return logger
}
