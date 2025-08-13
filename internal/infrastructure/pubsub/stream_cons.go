package pubsub

import (
	"context"
	"time"
)

type StreamConsumerManager interface {
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

	DeleteStream(
		ctx context.Context,
		stream string,
	) error
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

func (jsm *JetStreamManager) DeleteStream(ctx context.Context, name string) error {
	if err := jsm.js.DeleteStream(ctx, name); err != nil {
		jsm.logger.Error(ctx, "failed to delete stream", "error", err)
		return err
	}
	return nil
}
