package pubsub

import (
	"context"
	"time"

	"github.com/nats-io/nats.go/jetstream"
)

type ConsumerManager interface {
	CreateOrUpdateConsumer(
		ctx context.Context,
		stream string,
		cfg jetstream.ConsumerConfig,
	) (jetstream.Consumer, error)

	GetOrderedConsumer(
		ctx context.Context,
		stream string,
		cfg jetstream.OrderedConsumerConfig,
	) (jetstream.Consumer, error)

	GetConsumerByStream(
		ctx context.Context,
		stream string,
		consumer string,
	) (jetstream.Consumer, error)

	PauseConsumer(
		ctx context.Context,
		stream string,
		consumer string,
		until time.Time,
	) error

	ResumeConsumer(
		ctx context.Context,
		stream string,
		consumer string,
	) error

	DeleteConsumer(
		ctx context.Context,
		stream string,
		consumer string,
	) error
}

func (jsm *JetStreamManager) CreateOrUpdateConsumer(
	ctx context.Context,
	stream string,
	cfg jetstream.ConsumerConfig,
) (jetstream.Consumer, error) {
	cons, err := jsm.js.CreateOrUpdateConsumer(ctx, stream, cfg)
	if err != nil {
		jsm.logger.Error(ctx, "failed to create or update consumer", "error", err)
		return nil, err
	}
	return cons, nil
}

func (jsm *JetStreamManager) GetOrderedConsumer(
	ctx context.Context,
	stream string,
	cfg jetstream.OrderedConsumerConfig,
) (jetstream.Consumer, error) {
	cons, err := jsm.js.OrderedConsumer(ctx, stream, cfg)
	if err != nil {
		jsm.logger.Error(ctx, "failed to get ordered consumer", "error", err)
		return nil, err
	}
	return cons, nil
}

func (jsm *JetStreamManager) GetConsumerByStream(
	ctx context.Context,
	stream string,
	consumer string,
) (jetstream.Consumer, error) {
	cons, err := jsm.js.Consumer(ctx, stream, consumer)
	if err != nil {
		jsm.logger.Error(ctx, "failed to get consumer by stream", "error", err)
		return nil, err
	}
	return cons, nil
}

func (jsm *JetStreamManager) PauseConsumer(
	ctx context.Context,
	stream string,
	consumer string,
	until time.Time,
) error {
	resp, err := jsm.js.PauseConsumer(ctx, stream, consumer, until)
	if err != nil || !resp.Paused {
		jsm.logger.Error(ctx, "failed to pause consumer", "error", err)
		return err
	}
	return nil
}

func (jsm *JetStreamManager) ResumeConsumer(
	ctx context.Context,
	stream string,
	consumer string,
) error {
	_, err := jsm.js.ResumeConsumer(ctx, stream, consumer)
	if err != nil {
		jsm.logger.Error(ctx, "failed to resume consumer", "error", err)
		return err
	}
	return nil
}

func (jsm *JetStreamManager) DeleteConsumer(
	ctx context.Context,
	stream string,
	consumer string,
) error {
	err := jsm.js.DeleteConsumer(ctx, stream, consumer)
	if err != nil {
		jsm.logger.Error(ctx, "failed to delete consumer", "error", err)
		return err
	}
	return nil
}
