package mongodb

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/noxhalley/funken/config"
	"github.com/noxhalley/funken/internal/infrastructure/log"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type MongoDB struct {
	DBName string
	Client *mongo.Client
	logger *log.Logger
}

var (
	mongoInstance *MongoDB
	once          sync.Once
)

func NewOrGetSingleton(cfg *config.Config) *MongoDB {
	once.Do(func() {
		m, err := initMongo(cfg)
		if err != nil {
			panic(err)
		}
		mongoInstance = m
	})
	return mongoInstance
}

func initMongo(cfg *config.Config) (*MongoDB, error) {
	mongoLogger := newMongoLogger()
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
		logger: log.With("service", "mongodb"),
		Client: client,
	}, nil
}

func (m *MongoDB) Ping(ctx context.Context) error {
	return m.Client.Ping(ctx, readpref.PrimaryPreferred())
}

func (m *MongoDB) Close(ctx context.Context) {
	m.logger.Info(ctx, "Closing MongoDB")
	if err := m.Client.Disconnect(ctx); err != nil {
		m.logger.Error(ctx, "Error while closing MongoDB", "error", err.Error())
	}
}
