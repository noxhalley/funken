package mongodb

import (
	"context"

	"github.com/noxhalley/funken/internal/infrastructure/log"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type mongoLogger struct {
	logger *log.Logger
}

func newMongoLogger() options.LogSink {
	return &mongoLogger{
		logger: log.With("service", "mongodb"),
	}
}

func (l *mongoLogger) Info(level int, msg string, args ...interface{}) {
	if options.LogLevel(level+1) == options.LogLevelDebug {
		l.logger.Debug(context.Background(), msg, args...)
	} else {
		l.logger.Info(context.Background(), msg, args...)
	}
}

func (l *mongoLogger) Error(err error, msg string, args ...interface{}) {
	ctx := log.AddLogValToCtx(context.Background(), "error", err.Error())
	l.logger.Error(ctx, msg, args...)
}
