package pubsub

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

var (
	ErrInvalidStream   = errors.New("invalid stream")
	ErrInvalidConsumer = errors.New("invalid consumer")
)

type Publisher interface {
	Publish(
		ctx context.Context,
		subject string,
		payload interface{},
		metadata map[string][]string,
		opts ...jetstream.PublishOpt,
	) (*jetstream.PubAck, error)

	PublishAsync(
		ctx context.Context,
		subject string,
		payload interface{},
		metadata map[string][]string,
		opts ...jetstream.PublishOpt,
	) (
		<-chan *jetstream.PubAck,
		<-chan error,
		error,
	)
}

func (jsm *JetStreamManager) Publish(
	ctx context.Context,
	subject string,
	payload interface{},
	metadata map[string][]string,
	opts ...jetstream.PublishOpt,
) (*jetstream.PubAck, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		jsm.logger.Error(ctx, err.Error())
		return nil, err
	}

	msg := nats.Msg{
		Subject: subject,
		Data:    data,
		Header:  metadata,
	}

	pa, err := jsm.js.PublishMsg(ctx, &msg, opts...)
	if err != nil {
		jsm.logger.Error(ctx, "failed to publish message", "error", err)
		return nil, err
	}

	return pa, nil
}

func (jsm *JetStreamManager) PublishAsync(
	ctx context.Context,
	subject string,
	payload interface{},
	metadata map[string][]string,
	opts ...jetstream.PublishOpt,
) (
	<-chan *jetstream.PubAck,
	<-chan error,
	error,
) {

	data, err := json.Marshal(payload)
	if err != nil {
		jsm.logger.Error(ctx, err.Error())
		return nil, nil, err
	}

	msg := nats.Msg{
		Subject: subject,
		Data:    data,
		Header:  metadata,
	}

	paf, err := jsm.js.PublishMsgAsync(&msg, opts...)
	if err != nil {
		jsm.logger.Error(ctx, "failed to publish message async", "error", err)
		return nil, nil, err
	}

	return paf.Ok(), paf.Err(), nil
}
