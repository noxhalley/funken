package mongodb

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/noxhalley/funken/config"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type MongoDB struct {
	DBName string
	logger *logrus.Entry
	cli    *mongo.Client
}

var (
	mongoInstance *MongoDB
	once          sync.Once
)

func NewOrGetSingleton(ctx context.Context, cfg *config.Config, logger *logrus.Logger) *MongoDB {
	once.Do(func() {
		m, err := initMongo(ctx, cfg, logger)
		if err != nil {
			panic(err)
		}
		mongoInstance = m
	})
	return mongoInstance
}

func initMongo(ctx context.Context, cfg *config.Config, logger *logrus.Logger) (*MongoDB, error) {
	mongoLogger := newMongoLogger(logger)
	loggerOpts := options.
		Logger().
		SetSink(mongoLogger).
		SetComponentLevel(options.LogComponentCommand, options.LogLevelDebug)

	creds := options.Credential{
		AuthSource: cfg.Mongo.AuthSource,
		Username:   cfg.Mongo.Username,
		Password:   cfg.Mongo.Password,
	}

	uri := "mongodb://" + net.JoinHostPort(cfg.Mongo.Hostname, cfg.Mongo.Port)
	clientOpts := options.Client().
		ApplyURI(uri).
		SetAuth(creds).
		SetMaxPoolSize(uint64(cfg.Mongo.PoolSize)).
		SetTimeout(time.Duration(cfg.Mongo.Timeout) * time.Millisecond).
		SetConnectTimeout(time.Duration(cfg.Mongo.ConnTimeout) * time.Millisecond).
		SetLoggerOptions(loggerOpts)

	client, err := mongo.Connect(clientOpts)
	if err != nil {
		return nil, err
	}

	return &MongoDB{
		DBName: cfg.Mongo.Database,
		logger: logger.WithField("service", "mongodb"),
		cli:    client,
	}, nil
}

func (m *MongoDB) Ping(ctx context.Context) error {
	return m.cli.Ping(ctx, readpref.PrimaryPreferred())
}

func (m *MongoDB) Close(ctx context.Context) {
	if err := m.cli.Disconnect(ctx); err != nil {
		m.logger.Errorf("Error while closing MongoDB: %v", err)
	}
}
