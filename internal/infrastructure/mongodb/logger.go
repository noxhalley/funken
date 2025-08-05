package mongodb

import (
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type mongoLogger struct {
	logger *logrus.Entry
}

func newMongoLogger(logger *logrus.Logger) options.LogSink {
	return &mongoLogger{
		logger: logger.WithField("service", "mongodb"),
	}
}

func (l *mongoLogger) Info(level int, msg string, args ...interface{}) {
	if options.LogLevel(level+1) == options.LogLevelDebug {
		l.logger.Debug(args...)
	} else {
		l.logger.Info(args...)
	}
}

func (l *mongoLogger) Error(err error, msg string, args ...interface{}) {
	l.logger.Errorf(err.Error(), args...)
}
