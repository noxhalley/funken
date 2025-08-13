package pubsub

import (
	"context"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/noxhalley/funken/internal/pkg/utils"
)

type Subcriber interface {
	Subscribe(
		ctx context.Context,
		subject string,
		callback func(data interface{}),
		params SubcribeParams,
	) error
}

type SubcribeParams struct {
	Stream        string
	Consumer      string
	FilterSubject string
}

func (jsm *JetStreamManager) Subscribe(
	ctx context.Context,
	subject string,
	msgHandler func(data interface{}) error,
	params SubcribeParams,
) error {
	s, err := jsm.getStream(ctx, subject, params.Stream)
	if err != nil {
		return err
	}

	cons, err := jsm.getConsumer(ctx, s, params.Consumer, params.FilterSubject)
	if err != nil {
		return err
	}

	cc, err := cons.Consume(func(msg jetstream.Msg) {
		if handleErr := msgHandler(msg.Data()); handleErr != nil {
			jsm.logger.Warn(ctx, "failed to handle message", "error", handleErr)
			if nakErr := msg.NakWithDelay(3 * time.Second); nakErr != nil {
				jsm.logger.Error(ctx, "failed to NakWithDelay", "error", nakErr)
				return
			}
		}

		if ackErr := msg.Ack(); ackErr != nil {
			jsm.logger.Warn(ctx, "failed to ack message", "error", err)
		}
	})
	if err != nil {
		jsm.logger.Error(ctx, "failled to consume message", "error", err)
		return err
	}

	<-ctx.Done()
	cc.Stop()
	return ctx.Err()
}

func (jsm *JetStreamManager) getStream(
	ctx context.Context,
	subject string,
	streamInput string,
) (
	s jetstream.Stream,
	err error,
) {
	streamName, err := jsm.js.StreamNameBySubject(ctx, subject)
	if err == jetstream.ErrStreamNotFound {
		err = jsm.checkStreamIsEmpty(ctx, streamInput)
		if err != nil {
			return nil, err
		}

		s, err = jsm.js.Stream(ctx, streamInput)
		if err != nil && err != jetstream.ErrStreamNotFound {
			return nil, err
		}

		if err == jetstream.ErrStreamNotFound {
			s, err = jsm.js.CreateStream(ctx, jetstream.StreamConfig{
				Name:        streamInput,
				Subjects:    []string{subject},
				Storage:     jetstream.FileStorage,
				Replicas:    3,
				Retention:   jetstream.LimitsPolicy,
				MaxAge:      time.Duration(24) * time.Hour, // 1 day
				MaxBytes:    500 * 1024 * 1024,             // 500 MB
				Discard:     jetstream.DiscardOld,
				AllowDirect: true,
				Duplicates:  time.Duration(90) * time.Second, // 90s
			})

			if err != nil {
				jsm.logger.Error(ctx, "failed to create stream", "error", err)
				return nil, err
			}
		} else {
			info, err := s.Info(ctx)
			if err != nil {
				jsm.logger.Error(ctx, "failed to get stream info", "error", err)
				return nil, err
			}

			cfg := info.Config
			cfg.Subjects = append(cfg.Subjects, subject)
			s, err = jsm.js.UpdateStream(ctx, cfg)
			if err != nil {
				jsm.logger.Error(ctx, "failed to update stream", "error", err)
				return nil, err
			}
		}
	}

	if err != nil {
		jsm.logger.Error(ctx, "failed to get stream by subject", "error", err)
		return nil, err
	}

	if s == nil {
		if streamName != streamInput {
			return nil, ErrInvalidStream
		}
		s, err = jsm.js.Stream(ctx, streamInput)
		if err != nil {
			return nil, err
		}
	}

	return s, nil
}

func (jsm *JetStreamManager) getConsumer(
	ctx context.Context,
	s jetstream.Stream,
	consName string,
	filterSubj string,
) (
	cons jetstream.Consumer,
	err error,
) {
	cons, err = s.Consumer(ctx, consName)
	if err == jetstream.ErrConsumerNotFound {
		err = jsm.checkConsumerIsEmpty(ctx, consName)
		if err != nil {
			return nil, err
		}

		cons, err = s.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
			Name:              consName,
			Durable:           consName,
			FilterSubject:     filterSubj,
			DeliverPolicy:     jetstream.DeliverByStartTimePolicy,
			OptStartTime:      utils.ToPtr(time.Now().Add(-5 * time.Minute)),
			AckPolicy:         jetstream.AckExplicitPolicy,
			AckWait:           30 * time.Second,
			ReplayPolicy:      jetstream.ReplayInstantPolicy,
			InactiveThreshold: 10 * time.Minute,
			MaxDeliver:        5,
			BackOff: []time.Duration{
				500 * time.Millisecond,
				1 * time.Second,
				2 * time.Second,
			},
		})

		if err != nil {
			jsm.logger.Error(ctx, "failed to create or update consumer", "error", err)
			return nil, err
		}
	}

	if err != nil {
		jsm.logger.Error(ctx, "failed to get consumer", "error", err)
		return nil, err
	}

	return cons, nil
}

func (jsm *JetStreamManager) checkStreamIsEmpty(ctx context.Context, stream string) error {
	if stream == "" {
		jsm.logger.Error(ctx, "stream name must not be empty")
		return ErrInvalidStream
	}
	return nil
}

func (jsm *JetStreamManager) checkConsumerIsEmpty(ctx context.Context, consumer string) error {
	if consumer == "" {
		jsm.logger.Error(ctx, "consumer must not be empty")
		return ErrInvalidConsumer
	}
	return nil
}
