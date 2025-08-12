package pubsub

import (
	"context"
	"sync"
	"time"

	"github.com/noxhalley/funken/config"
	"github.com/noxhalley/funken/internal/infrastructure/log"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type jetStreamManager struct {
	logger *log.Logger
	js     jetstream.JetStream
	conn   *nats.Conn
}

var (
	once     sync.Once
	instance JetStreamManager
)

func NewOrGetSingleton(cfg *config.Config) JetStreamManager {
	once.Do(func() {
		jsm, err := initJetStream(cfg)
		if err != nil {
			panic(err)
		}
		instance = jsm
	})
	return instance
}

func initJetStream(cfg *config.Config) (JetStreamManager, error) {
	logger := log.With("service", "jetstream_manager")

	conn, err := initNATSConn(logger, cfg)
	if err != nil {
		logger.Error(context.Background(), "Failed to initialize NATS connection", "error", err)
		return nil, err
	}

	opts := []jetstream.JetStreamOpt{
		jetstream.WithDefaultTimeout(time.Duration(cfg.JetStream.Timeout) * time.Second),
		jetstream.WithPublishAsyncTimeout(time.Duration(cfg.JetStream.PublishAsyncTimeout) * time.Second),
		jetstream.WithPublishAsyncMaxPending(cfg.JetStream.PublishAsyncMaxPending),
		jetstream.WithClientTrace(&jetstream.ClientTrace{
			RequestSent: func(subj string, payload []byte) {
				logger.Debug(context.Background(), "JS Request Sent", "subject", subj, "payload", string(payload))
			},
			ResponseReceived: func(subj string, payload []byte, hdr nats.Header) {
				logger.Debug(context.Background(), "JS Response Received", "subject", subj, "payload", string(payload), "header", hdr)
			},
		}),
	}

	js, err := jetstream.NewWithDomain(conn, cfg.JetStream.Domain, opts...)
	if err != nil {
		conn.Close()
		logger.Error(context.Background(), "Failed to create JetStream instance", "error", err)
		return nil, err
	}

	return &jetStreamManager{
		logger: logger,
		conn:   conn,
		js:     js,
	}, nil
}

func initNATSConn(logger *log.Logger, cfg *config.Config) (*nats.Conn, error) {
	opts := []nats.Option{
		nats.Name(cfg.Nats.Name),
		nats.MaxReconnects(cfg.Nats.MaxReconnect),
		nats.ReconnectWait(time.Duration(cfg.Nats.ReconnectWait) * time.Millisecond),
		nats.ReconnectJitter(
			time.Duration(cfg.Nats.ReconnectJitter)*time.Millisecond,
			time.Duration(cfg.Nats.ReconnectJitterTLS)*time.Millisecond,
		),
		nats.Timeout(time.Duration(cfg.Nats.Timeout) * time.Millisecond),
		nats.PingInterval(time.Duration(cfg.Nats.PingInterval) * time.Minute),
		nats.MaxPingsOutstanding(cfg.Nats.MaxPingsOut),
		nats.ClosedHandler(func(c *nats.Conn) {
			logger.Info(context.Background(), "Closed connection to NATS")
		}),
		nats.DisconnectErrHandler(func(c *nats.Conn, err error) {
			logger.Warn(context.Background(), "Disconnected from NATS", "error", err)
		}),
	}

	return nats.Connect(cfg.Nats.Url, opts...)
}
