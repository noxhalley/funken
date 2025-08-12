package pubsub

import (
	"context"

	"github.com/nats-io/nats.go/jetstream"
)

type StreamManager interface {
	CreateOrUpdateStream(ctx context.Context, cfg jetstream.StreamConfig) (jetstream.Stream, error)

	GetStreamByName(ctx context.Context, name string) (jetstream.Stream, error)

	GetStreamBySubject(ctx context.Context, subj string) (jetstream.Stream, error)

	DeleteStream(ctx context.Context, name string) error
}

func (jsm *JetStreamManager) CreateOrUpdateStream(ctx context.Context, cfg jetstream.StreamConfig) (jetstream.Stream, error) {
	s, err := jsm.js.CreateOrUpdateStream(ctx, cfg)
	if err != nil {
		jsm.logger.Error(ctx, "failed to create or update stream", "error", err)
		return nil, err
	}
	return s, nil
}

func (jsm *JetStreamManager) GetStreamByName(ctx context.Context, name string) (jetstream.Stream, error) {
	s, err := jsm.js.Stream(ctx, name)
	if err != nil {
		jsm.logger.Error(ctx, "failed to find stream by name", "error", err)
		return nil, err
	}
	return s, nil
}

func (jsm *JetStreamManager) GetStreamBySubject(ctx context.Context, subj string) (jetstream.Stream, error) {
	name, err := jsm.js.StreamNameBySubject(ctx, subj)
	if err != nil {
		jsm.logger.Error(ctx, "failed to find stream by subject", "error", err)
		return nil, err
	}
	return jsm.GetStreamByName(ctx, name)
}

func (jsm *JetStreamManager) DeleteStream(ctx context.Context, name string) error {
	if err := jsm.js.DeleteStream(ctx, name); err != nil {
		jsm.logger.Error(ctx, "failed to delete stream", "error", err)
		return err
	}
	return nil
}
